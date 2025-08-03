package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	pb "transaction-services/proto"

	"google.golang.org/grpc/metadata"
)

type TransactionServer struct {
	pb.UnimplementedTransactionServiceServer
	db         *Database
	userClient UserServiceClient
}

func NewTransactionServer() *TransactionServer {
	userClient, err := NewUserServiceClient()
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}

	return &TransactionServer{
		db:         NewDatabase(),
		userClient: userClient,
	}
}

// extractTokenFromContext extracts the authorization token from gRPC metadata
func (s *TransactionServer) extractTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata found")
	}

	auth := md.Get("authorization")
	if len(auth) == 0 {
		return "", fmt.Errorf("no authorization header found")
	}

	return auth[0], nil
}

// validateRequest validates the token and returns user claims
func (s *TransactionServer) validateRequest(ctx context.Context) (*TokenClaims, error) {
	token, err := s.extractTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := ValidateTokenMiddleware(s.userClient, token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// GetTransactionByID retrieves a transaction by its ID
func (s *TransactionServer) GetTransactionByID(ctx context.Context, req *pb.GetTransactionByIDRequest) (*pb.GetTransactionByIDResponse, error) {
	log.Printf("GetTransactionByID called with ID: %s", req.Id)

	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetTransactionByIDResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetTransactionByIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.Id == "" {
		return &pb.GetTransactionByIDResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.GetTransactionByIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Transaction ID is required",
					Details:   []string{"ID field cannot be empty"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	query := `
		SELECT id, user_id, card_id, merchant_id, merchant_name, card_number, 
		       merchant_category, amount, currency, status, created_at, updated_at
		FROM transactions 
		WHERE id = $1 AND user_id = $2
	`

	row := s.db.DB.QueryRow(query, req.Id, claims.UserID)

	var transaction pb.Transaction
	var createdAt, updatedAt time.Time

	err = row.Scan(
		&transaction.Id,
		&transaction.UserID,
		&transaction.CardID,
		&transaction.MerchantID,
		&transaction.MerchantName,
		&transaction.CardNumber,
		&transaction.MerchantCategory,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Status,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetTransactionByIDResponse{
				StatusCode: 404,
				Message:    "Transaction not found",
				Result: &pb.GetTransactionByIDResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Transaction not found",
						Details:   []string{"No transaction found with this ID"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	transaction.CreateAt = createdAt.Format(time.RFC3339)
	transaction.UpdateAt = updatedAt.Format(time.RFC3339)

	return &pb.GetTransactionByIDResponse{
		StatusCode: 200,
		Message:    "Transaction retrieved successfully",
		Result: &pb.GetTransactionByIDResponse_Transaction{
			Transaction: &transaction,
		},
	}, nil
}

// GetTransactionByUserID retrieves transactions by user ID
func (s *TransactionServer) GetTransactionByUserID(ctx context.Context, req *pb.GetTransactionByUserIDRequest) (*pb.GetTransactionByUserIDResponse, error) {
	log.Printf("GetTransactionByUserID called with UserID: %s", req.UserID)

	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetTransactionByUserIDResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
		}, nil
	}

	// Ensure user can only access their own transactions
	if claims.UserID != req.UserID {
		return &pb.GetTransactionByUserIDResponse{
			StatusCode: 403,
			Message:    "Forbidden",
		}, nil
	}

	query := `
		SELECT id, user_id, card_id, merchant_id, merchant_name, card_number, 
		       merchant_category, amount, currency, status, created_at, updated_at
		FROM transactions 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.DB.Query(query, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	defer rows.Close()

	var transactions []*pb.Transaction
	for rows.Next() {
		var transaction pb.Transaction
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&transaction.Id,
			&transaction.UserID,
			&transaction.CardID,
			&transaction.MerchantID,
			&transaction.MerchantName,
			&transaction.CardNumber,
			&transaction.MerchantCategory,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		transaction.CreateAt = createdAt.Format(time.RFC3339)
		transaction.UpdateAt = updatedAt.Format(time.RFC3339)
		transactions = append(transactions, &transaction)
	}

	return &pb.GetTransactionByUserIDResponse{
		StatusCode:   200,
		Message:      "Transactions retrieved successfully",
		Transactions: transactions,
	}, nil
}
