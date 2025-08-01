package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	pb "card-services/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type CardServer struct {
	pb.UnimplementedCardServiceServer
	db         *Database
	userClient UserServiceClient
}

func NewCardServer() *CardServer {
	userClient, err := NewUserServiceClient()
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}

	return &CardServer{
		db:         NewDatabase(),
		userClient: userClient,
	}
}

// extractTokenFromContext extracts the authorization token from gRPC metadata
func (s *CardServer) extractTokenFromContext(ctx context.Context) (string, error) {
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
func (s *CardServer) validateRequest(ctx context.Context) (*TokenClaims, error) {
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

// GetCardByID retrieves a card by its ID
func (s *CardServer) GetCardByID(ctx context.Context, req *pb.GetCardByIDRequest) (*pb.GetCardByIDResponse, error) {
	log.Printf("GetCardByID called with ID: %s", req.Id)

	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetCardByIDResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetCardByIDResponse_Error{
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
		return &pb.GetCardByIDResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.GetCardByIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Card ID is required",
					Details:   []string{"ID field cannot be empty"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	query := `
		SELECT id, user_id, encrypted_cardholder_name, encrypted_card_number, 
		       encrypted_cvv, masked_card_number, expiration_date, card_type, 
		       status, balance, created_at, updated_at
		FROM cards 
		WHERE id = $1 AND user_id = $2
	`

	var card pb.Card
	var encryptedName, encryptedNumber, encryptedCVV string
	var createdAt, updatedAt time.Time

	err = s.db.GetDB().QueryRow(query, req.Id, claims.UserID).Scan(
		&card.Id, &card.UserId, &encryptedName, &encryptedNumber,
		&encryptedCVV, &card.CardNumber, &card.ExpirationDate,
		&card.CardType, &card.Status, &card.Balance,
		&createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetCardByIDResponse{
				StatusCode: 404,
				Message:    "Card not found",
				Result: &pb.GetCardByIDResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Card not found",
						Details:   []string{"No card found with the provided ID"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return &pb.GetCardByIDResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.GetCardByIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Database error",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Decrypt sensitive data
	cardHolderName, err := DecryptData(encryptedName)
	if err != nil {
		log.Printf("Failed to decrypt cardholder name: %v", err)
		cardHolderName = "[ENCRYPTED]"
	}

	cvv, err := DecryptData(encryptedCVV)
	if err != nil {
		log.Printf("Failed to decrypt CVV: %v", err)
		cvv = "***"
	}

	card.CardHolderName = cardHolderName
	card.Cvv = cvv
	card.CreateAt = createdAt.Format(time.RFC3339)
	card.UpdateAt = updatedAt.Format(time.RFC3339)

	return &pb.GetCardByIDResponse{
		StatusCode: 200,
		Message:    "Card retrieved successfully",
		Result: &pb.GetCardByIDResponse_Card{
			Card: &card,
		},
	}, nil
}

// CreateNewCard creates a new card with encrypted sensitive data
func (s *CardServer) CreateNewCard(ctx context.Context, req *pb.CreateNewCardRequest) (*pb.CreateNewCardResponse, error) {
	log.Printf("CreateNewCard called for user: %s", req.UserId)

	// Validate token
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.CreateNewCardResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.CreateNewCardResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Validate input
	if req.CardNumber == "" || req.CardHolderName == "" || req.ExpirationDate == "" || req.Cvv == "" {
		return &pb.CreateNewCardResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.CreateNewCardResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Missing required fields",
					Details:   []string{"Card number, holder name, expiration date, and CVV are required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Validate card number
	if !ValidateCardNumber(req.CardNumber) {
		return &pb.CreateNewCardResponse{
			StatusCode: 400,
			Message:    "Invalid card number",
			Result: &pb.CreateNewCardResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid card number format",
					Details:   []string{"Card number failed validation"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Encrypt sensitive data
	encryptedName, err := EncryptData(req.CardHolderName)
	if err != nil {
		return &pb.CreateNewCardResponse{
			StatusCode: 500,
			Message:    "Encryption failed",
			Result: &pb.CreateNewCardResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to encrypt cardholder name",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	encryptedNumber, err := EncryptData(req.CardNumber)
	if err != nil {
		return &pb.CreateNewCardResponse{
			StatusCode: 500,
			Message:    "Encryption failed",
			Result: &pb.CreateNewCardResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to encrypt card number",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	encryptedCVV, err := EncryptData(req.Cvv)
	if err != nil {
		return &pb.CreateNewCardResponse{
			StatusCode: 500,
			Message:    "Encryption failed",
			Result: &pb.CreateNewCardResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to encrypt CVV",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Generate card ID and mask card number
	cardID := uuid.New().String()
	maskedNumber := MaskCardNumber(req.CardNumber)

	// Insert into database
	query := `
		INSERT INTO cards (id, user_id, encrypted_cardholder_name, encrypted_card_number, 
		                   encrypted_cvv, masked_card_number, expiration_date, card_type, 
		                   status, balance)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at
	`

	var createdAt, updatedAt time.Time
	err = s.db.GetDB().QueryRow(query, cardID, claims.UserID, encryptedName,
		encryptedNumber, encryptedCVV, maskedNumber, req.ExpirationDate,
		req.CardType, req.Status, req.Balance).Scan(&createdAt, &updatedAt)

	if err != nil {
		return &pb.CreateNewCardResponse{
			StatusCode: 500,
			Message:    "Database error",
			Result: &pb.CreateNewCardResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to create card",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Return the created card (with masked sensitive data)
	card := &pb.Card{
		Id:             cardID,
		UserId:         claims.UserID,
		CardNumber:     maskedNumber,
		CardHolderName: req.CardHolderName,
		ExpirationDate: req.ExpirationDate,
		Cvv:            "***", // Never return real CVV
		Status:         req.Status,
		CardType:       req.CardType,
		Balance:        req.Balance,
		CreateAt:       createdAt.Format(time.RFC3339),
		UpdateAt:       updatedAt.Format(time.RFC3339),
	}

	return &pb.CreateNewCardResponse{
		StatusCode: 201,
		Message:    "Card created successfully",
		Result: &pb.CreateNewCardResponse_Card{
			Card: card,
		},
	}, nil
}