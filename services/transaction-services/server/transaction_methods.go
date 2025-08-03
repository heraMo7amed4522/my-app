package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	pb "transaction-services/proto"
)

// GetTransactionByCardID retrieves transactions by card ID
func (s *TransactionServer) GetTransactionByCardID(ctx context.Context, req *pb.GetTransactionByCardIDRequest) (*pb.GetTransactionByCardIDResponse, error) {
	log.Printf("GetTransactionByCardID called with CardID: %s", req.CardID)

	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetTransactionByCardIDResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
		}, nil
	}

	// Verify that the card belongs to the user
	cardQuery := `SELECT user_id FROM cards WHERE id = $1`
	var cardUserID string
	err = s.db.DB.QueryRow(cardQuery, req.CardID).Scan(&cardUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetTransactionByCardIDResponse{
				StatusCode: 404,
				Message:    "Card not found",
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	if cardUserID != claims.UserID {
		return &pb.GetTransactionByCardIDResponse{
			StatusCode: 403,
			Message:    "Forbidden",
		}, nil
	}

	query := `
		SELECT id, user_id, card_id, merchant_id, merchant_name, card_number, 
		       merchant_category, amount, currency, status, created_at, updated_at
		FROM transactions 
		WHERE card_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.DB.Query(query, req.CardID)
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

	return &pb.GetTransactionByCardIDResponse{
		StatusCode:   200,
		Message:      "Transactions retrieved successfully",
		Transactions: transactions,
	}, nil
}

// GetTransactionByStatus retrieves transactions by status
func (s *TransactionServer) GetTransactionByStatus(ctx context.Context, req *pb.GetTransactionByStatusRequest) (*pb.GetTransactionByStatusResponse, error) {
	log.Printf("GetTransactionByStatus called with Status: %s", req.Status)

	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetTransactionByStatusResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
		}, nil
	}

	query := `
		SELECT id, user_id, card_id, merchant_id, merchant_name, card_number, 
		       merchant_category, amount, currency, status, created_at, updated_at
		FROM transactions 
		WHERE status = $1 AND user_id = $2
		ORDER BY created_at DESC
	`

	rows, err := s.db.DB.Query(query, req.Status, claims.UserID)
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

	return &pb.GetTransactionByStatusResponse{
		StatusCode:   200,
		Message:      "Transactions retrieved successfully",
		Transactions: transactions,
	}, nil
}

// GetTransactionByDate retrieves transactions by date
func (s *TransactionServer) GetTransactionByDate(ctx context.Context, req *pb.GetTransactionByDateRequest) (*pb.GetTransactionByDateResponse, error) {
	log.Printf("GetTransactionByDate called with Date: %s", req.Date)

	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetTransactionByDateResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
		}, nil
	}

	// Parse the date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return &pb.GetTransactionByDateResponse{
			StatusCode: 400,
			Message:    "Invalid date format. Use YYYY-MM-DD",
		}, nil
	}

	startOfDay := date
	endOfDay := date.Add(24 * time.Hour)

	query := `
		SELECT id, user_id, card_id, merchant_id, merchant_name, card_number, 
		       merchant_category, amount, currency, status, created_at, updated_at
		FROM transactions 
		WHERE created_at >= $1 AND created_at < $2 AND user_id = $3
		ORDER BY created_at DESC
	`

	rows, err := s.db.DB.Query(query, startOfDay, endOfDay, claims.UserID)
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

	return &pb.GetTransactionByDateResponse{
		StatusCode:   200,
		Message:      "Transactions retrieved successfully",
		Transactions: transactions,
	}, nil
}
