package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	pb "wallet-services/proto"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WalletServer struct {
	pb.UnimplementedWalletServiceServer
	db         *Database
	userClient UserServiceClient
}

func NewWalletServer() *WalletServer {
	userClient, err := NewUserServiceClient()
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}

	return &WalletServer{
		db:         NewDatabase(),
		userClient: userClient,
	}
}

// extractTokenFromContext extracts the authorization token from gRPC metadata
func (s *WalletServer) extractTokenFromContext(ctx context.Context) (string, error) {
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
func (s *WalletServer) validateRequest(ctx context.Context) (*TokenClaims, error) {
	token, err := s.extractTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to extract token: %v", err)
	}

	claims, err := s.userClient.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %v", err)
	}

	return claims, nil
}

// GetWalletByUserID retrieves a wallet by user ID
func (s *WalletServer) GetWalletByUserID(ctx context.Context, req *pb.GetWalletByUserIDRequest) (*pb.GetWalletByUserIDResponse, error) {
	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetWalletByUserIDResponse{
			StatusCode: 401,
			Message:    "Unauthorized: " + err.Error(),
			Result:     &pb.GetWalletByUserIDResponse_Error{Error: &pb.ErrorDetails{Code: 401, Message: err.Error()}},
		}, nil
	}

	// Validate request
	if err := ValidateGetWalletByUserIDRequest(req); err != nil {
		return &pb.GetWalletByUserIDResponse{
			StatusCode: 400,
			Message:    "Bad Request: " + err.Error(),
			Result:     &pb.GetWalletByUserIDResponse_Error{Error: &pb.ErrorDetails{Code: 400, Message: err.Error()}},
		}, nil
	}

	// Check if user can access this wallet
	if claims.UserID != req.UserID {
		return &pb.GetWalletByUserIDResponse{
			StatusCode: 403,
			Message:    "Forbidden: Cannot access another user's wallet",
			Result:     &pb.GetWalletByUserIDResponse_Error{Error: &pb.ErrorDetails{Code: 403, Message: "Forbidden"}},
		}, nil
	}

	// Query database
	query := `SELECT id, user_id, currency, balance, status, created_at, updated_at FROM wallets WHERE user_id = $1`
	row := s.db.DB.QueryRow(query, req.UserID)

	var wallet pb.Wallet
	var encryptedBalance string
	var createdAt, updatedAt time.Time

	err = row.Scan(&wallet.Id, &wallet.UserID, &wallet.Currency, &encryptedBalance, &wallet.Status, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetWalletByUserIDResponse{
				StatusCode: 404,
				Message:    "Wallet not found",
				Result:     &pb.GetWalletByUserIDResponse_Error{Error: &pb.ErrorDetails{Code: 404, Message: "Wallet not found"}},
			}, nil
		}
		return &pb.GetWalletByUserIDResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.GetWalletByUserIDResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Database error"}},
		}, nil
	}

	// Decrypt balance
	balance, err := DecryptData(encryptedBalance)
	if err != nil {
		log.Printf("Failed to decrypt balance for wallet %s: %v", wallet.Id, err)
		return &pb.GetWalletByUserIDResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.GetWalletByUserIDResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Decryption error"}},
		}, nil
	}

	wallet.Balance = balance
	wallet.CreatedAt = timestamppb.New(createdAt)
	wallet.UpdatedAt = timestamppb.New(updatedAt)

	return &pb.GetWalletByUserIDResponse{
		StatusCode: 200,
		Message:    "Wallet retrieved successfully",
		Result:     &pb.GetWalletByUserIDResponse_Wallet{Wallet: &wallet},
	}, nil
}

// CreateWallet creates a new wallet for a user
func (s *WalletServer) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {
	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.CreateWalletResponse{
			StatusCode: 401,
			Message:    "Unauthorized: " + err.Error(),
			Result:     &pb.CreateWalletResponse_Error{Error: &pb.ErrorDetails{Code: 401, Message: err.Error()}},
		}, nil
	}

	// Validate request
	if err := ValidateCreateWalletRequest(req); err != nil {
		return &pb.CreateWalletResponse{
			StatusCode: 400,
			Message:    "Bad Request: " + err.Error(),
			Result:     &pb.CreateWalletResponse_Error{Error: &pb.ErrorDetails{Code: 400, Message: err.Error()}},
		}, nil
	}

	// Check if user can create wallet for this user ID
	if claims.UserID != req.UserID {
		return &pb.CreateWalletResponse{
			StatusCode: 403,
			Message:    "Forbidden: Cannot create wallet for another user",
			Result:     &pb.CreateWalletResponse_Error{Error: &pb.ErrorDetails{Code: 403, Message: "Forbidden"}},
		}, nil
	}

	// Check if wallet already exists
	checkQuery := `SELECT id FROM wallets WHERE user_id = $1 AND currency = $2`
	var existingID string
	err = s.db.DB.QueryRow(checkQuery, req.UserID, req.Currency).Scan(&existingID)
	if err == nil {
		return &pb.CreateWalletResponse{
			StatusCode: 409,
			Message:    "Wallet already exists for this user and currency",
			Result:     &pb.CreateWalletResponse_Error{Error: &pb.ErrorDetails{Code: 409, Message: "Wallet already exists"}},
		}, nil
	}

	// Encrypt initial balance (0.00)
	encryptedBalance, err := EncryptData("0.00")
	if err != nil {
		return &pb.CreateWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.CreateWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Encryption error"}},
		}, nil
	}

	// Insert new wallet
	insertQuery := `
		INSERT INTO wallets (user_id, currency, balance, status, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id, created_at, updated_at`

	now := time.Now()
	var walletID string
	var createdAt, updatedAt time.Time

	err = s.db.DB.QueryRow(insertQuery, req.UserID, req.Currency, encryptedBalance, "ACTIVE", now, now).Scan(&walletID, &createdAt, &updatedAt)
	if err != nil {
		log.Printf("Failed to create wallet: %v", err)
		return &pb.CreateWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.CreateWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Database error"}},
		}, nil
	}

	wallet := &pb.Wallet{
		Id:        walletID,
		UserID:    req.UserID,
		Currency:  req.Currency,
		Balance:   "0.00",
		Status:    "ACTIVE",
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}

	return &pb.CreateWalletResponse{
		StatusCode: 201,
		Message:    "Wallet created successfully",
		Result:     &pb.CreateWalletResponse_Wallet{Wallet: wallet},
	}, nil
}

// FundWallet adds funds to a user's wallet
func (s *WalletServer) FundWallet(ctx context.Context, req *pb.FundWalletRequest) (*pb.FundWalletResponse, error) {
	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.FundWalletResponse{
			StatusCode: 401,
			Message:    "Unauthorized: " + err.Error(),
			Result:     &pb.FundWalletResponse_Error{Error: &pb.ErrorDetails{Code: 401, Message: err.Error()}},
		}, nil
	}

	// Validate request
	if err := ValidateFundWalletRequest(req); err != nil {
		return &pb.FundWalletResponse{
			StatusCode: 400,
			Message:    "Bad Request: " + err.Error(),
			Result:     &pb.FundWalletResponse_Error{Error: &pb.ErrorDetails{Code: 400, Message: err.Error()}},
		}, nil
	}

	// Check if user can fund this wallet
	if claims.UserID != req.UserID {
		return &pb.FundWalletResponse{
			StatusCode: 403,
			Message:    "Forbidden: Cannot fund another user's wallet",
			Result:     &pb.FundWalletResponse_Error{Error: &pb.ErrorDetails{Code: 403, Message: "Forbidden"}},
		}, nil
	}

	// Start transaction
	tx, err := s.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return &pb.FundWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.FundWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Transaction error"}},
		}, nil
	}
	defer tx.Rollback()

	// Get current wallet
	query := `SELECT id, balance FROM wallets WHERE user_id = $1 AND currency = $2 AND status = 'ACTIVE' FOR UPDATE`
	var walletID, encryptedBalance string
	err = tx.QueryRow(query, req.UserID, req.Currency).Scan(&walletID, &encryptedBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.FundWalletResponse{
				StatusCode: 404,
				Message:    "Wallet not found or inactive",
				Result:     &pb.FundWalletResponse_Error{Error: &pb.ErrorDetails{Code: 404, Message: "Wallet not found"}},
			}, nil
		}
		return &pb.FundWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.FundWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Database error"}},
		}, nil
	}

	// Decrypt current balance
	currentBalance, err := DecryptData(encryptedBalance)
	if err != nil {
		return &pb.FundWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.FundWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Decryption error"}},
		}, nil
	}

	// Calculate new balance
	currentVal, _ := strconv.ParseFloat(currentBalance, 64)
	fundAmount, _ := strconv.ParseFloat(req.Amount, 64)
	newBalance := currentVal + fundAmount
	newBalanceStr := fmt.Sprintf("%.2f", newBalance)

	// Encrypt new balance
	encryptedNewBalance, err := EncryptData(newBalanceStr)
	if err != nil {
		return &pb.FundWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.FundWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Encryption error"}},
		}, nil
	}

	// Update wallet balance
	updateQuery := `UPDATE wallets SET balance = $1, updated_at = $2 WHERE id = $3`
	now := time.Now()
	_, err = tx.Exec(updateQuery, encryptedNewBalance, now, walletID)
	if err != nil {
		return &pb.FundWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.FundWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Update error"}},
		}, nil
	}

	// Record wallet transaction
	transactionQuery := `
		INSERT INTO wallet_transactions (wallet_id, type, amount, currency, description, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(transactionQuery, walletID, "FUND", req.Amount, req.Currency, "Wallet funding", now)
	if err != nil {
		log.Printf("Failed to record wallet transaction: %v", err)
		// Continue without failing the main operation
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return &pb.FundWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.FundWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Commit error"}},
		}, nil
	}

	// Get updated wallet
	updatedWallet, err := s.getWalletByID(walletID)
	if err != nil {
		return &pb.FundWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.FundWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Fetch error"}},
		}, nil
	}

	return &pb.FundWalletResponse{
		StatusCode: 200,
		Message:    "Wallet funded successfully",
		Result:     &pb.FundWalletResponse_Wallet{Wallet: updatedWallet},
	}, nil
}

// DeductFromWallet deducts funds from a user's wallet
func (s *WalletServer) DeductFromWallet(ctx context.Context, req *pb.DeductFromWalletRequest) (*pb.DeductFromWalletResponse, error) {
	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.DeductFromWalletResponse{
			StatusCode: 401,
			Message:    "Unauthorized: " + err.Error(),
			Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 401, Message: err.Error()}},
		}, nil
	}

	// Validate request
	if err := ValidateDeductFromWalletRequest(req); err != nil {
		return &pb.DeductFromWalletResponse{
			StatusCode: 400,
			Message:    "Bad Request: " + err.Error(),
			Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 400, Message: err.Error()}},
		}, nil
	}

	// Check if user can deduct from this wallet
	if claims.UserID != req.UserID {
		return &pb.DeductFromWalletResponse{
			StatusCode: 403,
			Message:    "Forbidden: Cannot deduct from another user's wallet",
			Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 403, Message: "Forbidden"}},
		}, nil
	}

	// Start transaction
	tx, err := s.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return &pb.DeductFromWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Transaction error"}},
		}, nil
	}
	defer tx.Rollback()

	// Get current wallet
	query := `SELECT id, balance FROM wallets WHERE user_id = $1 AND currency = $2 AND status = 'ACTIVE' FOR UPDATE`
	var walletID, encryptedBalance string
	err = tx.QueryRow(query, req.UserID, req.Currency).Scan(&walletID, &encryptedBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.DeductFromWalletResponse{
				StatusCode: 404,
				Message:    "Wallet not found or inactive",
				Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 404, Message: "Wallet not found"}},
			}, nil
		}
		return &pb.DeductFromWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Database error"}},
		}, nil
	}

	// Decrypt current balance
	currentBalance, err := DecryptData(encryptedBalance)
	if err != nil {
		return &pb.DeductFromWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Decryption error"}},
		}, nil
	}

	// Check sufficient balance
	currentVal, _ := strconv.ParseFloat(currentBalance, 64)
	deductAmount, _ := strconv.ParseFloat(req.Amount, 64)
	if currentVal < deductAmount {
		return &pb.DeductFromWalletResponse{
			StatusCode: 400,
			Message:    "Insufficient balance",
			Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 400, Message: "Insufficient balance"}},
		}, nil
	}

	// Calculate new balance
	newBalance := currentVal - deductAmount
	newBalanceStr := fmt.Sprintf("%.2f", newBalance)

	// Encrypt new balance
	encryptedNewBalance, err := EncryptData(newBalanceStr)
	if err != nil {
		return &pb.DeductFromWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Encryption error"}},
		}, nil
	}

	// Update wallet balance
	updateQuery := `UPDATE wallets SET balance = $1, updated_at = $2 WHERE id = $3`
	now := time.Now()
	_, err = tx.Exec(updateQuery, encryptedNewBalance, now, walletID)
	if err != nil {
		return &pb.DeductFromWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Update error"}},
		}, nil
	}

	// Record wallet transaction
	transactionQuery := `
		INSERT INTO wallet_transactions (wallet_id, type, amount, currency, description, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(transactionQuery, walletID, "DEDUCT", req.Amount, req.Currency, "Wallet deduction", now)
	if err != nil {
		log.Printf("Failed to record wallet transaction: %v", err)
		// Continue without failing the main operation
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return &pb.DeductFromWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Commit error"}},
		}, nil
	}

	// Get updated wallet
	updatedWallet, err := s.getWalletByID(walletID)
	if err != nil {
		return &pb.DeductFromWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.DeductFromWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Fetch error"}},
		}, nil
	}

	return &pb.DeductFromWalletResponse{
		StatusCode: 200,
		Message:    "Amount deducted successfully",
		Result:     &pb.DeductFromWalletResponse_Wallet{Wallet: updatedWallet},
	}, nil
}

// GetWalletBalance retrieves the balance of a user's wallet
func (s *WalletServer) GetWalletBalance(ctx context.Context, req *pb.GetWalletBalanceRequest) (*pb.GetWalletBalanceResponse, error) {
	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetWalletBalanceResponse{
			StatusCode: 401,
			Message:    "Unauthorized: " + err.Error(),
			Result:     &pb.GetWalletBalanceResponse_Error{Error: &pb.ErrorDetails{Code: 401, Message: err.Error()}},
		}, nil
	}

	// Validate request
	if err := ValidateGetWalletBalanceRequest(req); err != nil {
		return &pb.GetWalletBalanceResponse{
			StatusCode: 400,
			Message:    "Bad Request: " + err.Error(),
			Result:     &pb.GetWalletBalanceResponse_Error{Error: &pb.ErrorDetails{Code: 400, Message: err.Error()}},
		}, nil
	}

	// Check if user can access this wallet
	if claims.UserID != req.UserID {
		return &pb.GetWalletBalanceResponse{
			StatusCode: 403,
			Message:    "Forbidden: Cannot access another user's wallet",
			Result:     &pb.GetWalletBalanceResponse_Error{Error: &pb.ErrorDetails{Code: 403, Message: "Forbidden"}},
		}, nil
	}

	// Get wallet by user ID (default currency USD)
	query := `SELECT id, user_id, currency, balance, status, created_at, updated_at FROM wallets WHERE user_id = $1 AND status = 'ACTIVE' LIMIT 1`
	row := s.db.DB.QueryRow(query, req.UserID)

	var wallet pb.Wallet
	var encryptedBalance string
	var createdAt, updatedAt time.Time

	err = row.Scan(&wallet.Id, &wallet.UserID, &wallet.Currency, &encryptedBalance, &wallet.Status, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetWalletBalanceResponse{
				StatusCode: 404,
				Message:    "Wallet not found",
				Result:     &pb.GetWalletBalanceResponse_Error{Error: &pb.ErrorDetails{Code: 404, Message: "Wallet not found"}},
			}, nil
		}
		return &pb.GetWalletBalanceResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.GetWalletBalanceResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Database error"}},
		}, nil
	}

	// Decrypt balance
	balance, err := DecryptData(encryptedBalance)
	if err != nil {
		return &pb.GetWalletBalanceResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.GetWalletBalanceResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Decryption error"}},
		}, nil
	}

	wallet.Balance = balance
	wallet.CreatedAt = timestamppb.New(createdAt)
	wallet.UpdatedAt = timestamppb.New(updatedAt)

	return &pb.GetWalletBalanceResponse{
		StatusCode: 200,
		Message:    "Wallet balance retrieved successfully",
		Result:     &pb.GetWalletBalanceResponse_Wallet{Wallet: &wallet},
	}, nil
}

// GetWalletTransactions retrieves wallet transactions for a user
func (s *WalletServer) GetWalletTransactions(ctx context.Context, req *pb.GetWalletTransactionsRequest) (*pb.GetWalletTransactionsResponse, error) {
	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetWalletTransactionsResponse{
			StatusCode: 401,
			Message:    "Unauthorized: " + err.Error(),
		}, nil
	}

	// Validate request
	if err := ValidateGetWalletTransactionsRequest(req); err != nil {
		return &pb.GetWalletTransactionsResponse{
			StatusCode: 400,
			Message:    "Bad Request: " + err.Error(),
		}, nil
	}

	// Check if user can access these transactions
	if claims.UserID != req.UserID {
		return &pb.GetWalletTransactionsResponse{
			StatusCode: 403,
			Message:    "Forbidden: Cannot access another user's transactions",
		}, nil
	}

	// Get wallet ID first
	walletQuery := `SELECT id FROM wallets WHERE user_id = $1 AND status = 'ACTIVE' LIMIT 1`
	var walletID string
	err = s.db.DB.QueryRow(walletQuery, req.UserID).Scan(&walletID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetWalletTransactionsResponse{
				StatusCode: 404,
				Message:    "Wallet not found",
			}, nil
		}
		return &pb.GetWalletTransactionsResponse{
			StatusCode: 500,
			Message:    "Internal server error",
		}, nil
	}

	// Get wallet transactions
	transactionQuery := `
		SELECT id, wallet_id, type, amount, currency, description, created_at 
		FROM wallet_transactions 
		WHERE wallet_id = $1 
		ORDER BY created_at DESC 
		LIMIT 100`

	rows, err := s.db.DB.Query(transactionQuery, walletID)
	if err != nil {
		return &pb.GetWalletTransactionsResponse{
			StatusCode: 500,
			Message:    "Internal server error",
		}, nil
	}
	defer rows.Close()

	var transactions []*pb.WalletTransaction
	for rows.Next() {
		var transaction pb.WalletTransaction
		var createdAt time.Time

		err := rows.Scan(
			&transaction.Id,
			&transaction.WalletID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Description,
			&createdAt,
		)
		if err != nil {
			log.Printf("Error scanning transaction: %v", err)
			continue
		}

		transaction.CreatedAt = timestamppb.New(createdAt)
		transactions = append(transactions, &transaction)
	}

	return &pb.GetWalletTransactionsResponse{
		StatusCode:   200,
		Message:      "Wallet transactions retrieved successfully",
		Transactions: transactions,
	}, nil
}

// RefundToWallet refunds money to a wallet from a transaction
func (s *WalletServer) RefundToWallet(ctx context.Context, req *pb.RefundToWalletRequest) (*pb.RefundToWalletResponse, error) {
	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.RefundToWalletResponse{
			StatusCode: 401,
			Message:    "Unauthorized: " + err.Error(),
			Result:     &pb.RefundToWalletResponse_Error{Error: &pb.ErrorDetails{Code: 401, Message: err.Error()}},
		}, nil
	}

	// Validate request
	if err := ValidateRefundToWalletRequest(req); err != nil {
		return &pb.RefundToWalletResponse{
			StatusCode: 400,
			Message:    "Bad Request: " + err.Error(),
			Result:     &pb.RefundToWalletResponse_Error{Error: &pb.ErrorDetails{Code: 400, Message: err.Error()}},
		}, nil
	}

	// For now, we'll implement a simple refund that adds money to the user's default wallet
	// In a real implementation, you'd verify the transaction exists and belongs to the user

	// Get user's default wallet (first active wallet)
	walletQuery := `SELECT id, user_id, currency, balance FROM wallets WHERE user_id = $1 AND status = 'ACTIVE' LIMIT 1`
	var walletID, userID, currency, encryptedBalance string
	err = s.db.DB.QueryRow(walletQuery, claims.UserID).Scan(&walletID, &userID, &currency, &encryptedBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.RefundToWalletResponse{
				StatusCode: 404,
				Message:    "Wallet not found",
				Result:     &pb.RefundToWalletResponse_Error{Error: &pb.ErrorDetails{Code: 404, Message: "Wallet not found"}},
			}, nil
		}
		return &pb.RefundToWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.RefundToWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Database error"}},
		}, nil
	}

	// Start transaction
	tx, err := s.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return &pb.RefundToWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.RefundToWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Transaction error"}},
		}, nil
	}
	defer tx.Rollback()

	// Decrypt current balance
	currentBalance, err := DecryptData(encryptedBalance)
	if err != nil {
		return &pb.RefundToWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.RefundToWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Decryption error"}},
		}, nil
	}

	// Calculate new balance
	currentVal, _ := strconv.ParseFloat(currentBalance, 64)
	refundAmount, _ := strconv.ParseFloat(req.Amount, 64)
	newBalance := currentVal + refundAmount
	newBalanceStr := fmt.Sprintf("%.2f", newBalance)

	// Encrypt new balance
	encryptedNewBalance, err := EncryptData(newBalanceStr)
	if err != nil {
		return &pb.RefundToWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.RefundToWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Encryption error"}},
		}, nil
	}

	// Update wallet balance
	updateQuery := `UPDATE wallets SET balance = $1, updated_at = $2 WHERE id = $3`
	now := time.Now()
	_, err = tx.Exec(updateQuery, encryptedNewBalance, now, walletID)
	if err != nil {
		return &pb.RefundToWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.RefundToWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Update error"}},
		}, nil
	}

	// Record wallet transaction
	transactionQuery := `
		INSERT INTO wallet_transactions (wallet_id, type, amount, currency, description, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	description := fmt.Sprintf("Refund from transaction %s", req.TransactionID)
	_, err = tx.Exec(transactionQuery, walletID, "REFUND", req.Amount, currency, description, now)
	if err != nil {
		log.Printf("Failed to record wallet transaction: %v", err)
		// Continue without failing the main operation
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return &pb.RefundToWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.RefundToWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Commit error"}},
		}, nil
	}

	// Get updated wallet
	updatedWallet, err := s.getWalletByID(walletID)
	if err != nil {
		return &pb.RefundToWalletResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result:     &pb.RefundToWalletResponse_Error{Error: &pb.ErrorDetails{Code: 500, Message: "Fetch error"}},
		}, nil
	}

	return &pb.RefundToWalletResponse{
		StatusCode: 200,
		Message:    "Refund processed successfully",
		Result:     &pb.RefundToWalletResponse_Wallet{Wallet: updatedWallet},
	}, nil
}

// getWalletByID is a helper function to retrieve a wallet by ID
func (s *WalletServer) getWalletByID(walletID string) (*pb.Wallet, error) {
	query := `SELECT id, user_id, currency, balance, status, created_at, updated_at FROM wallets WHERE id = $1`
	row := s.db.DB.QueryRow(query, walletID)

	var wallet pb.Wallet
	var encryptedBalance string
	var createdAt, updatedAt time.Time

	err := row.Scan(&wallet.Id, &wallet.UserID, &wallet.Currency, &encryptedBalance, &wallet.Status, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	// Decrypt balance
	balance, err := DecryptData(encryptedBalance)
	if err != nil {
		return nil, err
	}

	wallet.Balance = balance
	wallet.CreatedAt = timestamppb.New(createdAt)
	wallet.UpdatedAt = timestamppb.New(updatedAt)

	return &wallet, nil
}