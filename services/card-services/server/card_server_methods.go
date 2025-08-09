package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pb "card-services/proto"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/paymentintent"
)

func (s *CardServer) UpdateCardData(ctx context.Context, req *pb.UpdateCardDataRequest) (*pb.UpdateCardDataResponse, error) {
	log.Printf("UpdateCardData called for card ID: %s", req.Id)
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.UpdateCardDataResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.UpdateCardDataResponse_Error{
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
		return &pb.UpdateCardDataResponse{
			StatusCode: 400,
			Message:    "Card ID is required",
			Result: &pb.UpdateCardDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Card ID cannot be empty",
					Details:   []string{"ID field is required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	var existingUserID string
	err = s.db.GetDB().QueryRow("SELECT user_id FROM cards WHERE id = $1", req.Id).Scan(&existingUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UpdateCardDataResponse{
				StatusCode: 404,
				Message:    "Card not found",
				Result: &pb.UpdateCardDataResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Card not found",
						Details:   []string{"No card found with the provided ID"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return &pb.UpdateCardDataResponse{
			StatusCode: 500,
			Message:    "Database error",
			Result: &pb.UpdateCardDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Database query failed",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if existingUserID != claims.UserID {
		return &pb.UpdateCardDataResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.UpdateCardDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      403,
					Message:   "Access denied",
					Details:   []string{"You can only update your own cards"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	updateFields := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.CardHolderName != "" {
		encryptedName, erroe := EncryptData(req.CardHolderName)
		if erroe != nil {
			return &pb.UpdateCardDataResponse{
				StatusCode: 500,
				Message:    "Encryption failed",
				Result: &pb.UpdateCardDataResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      500,
						Message:   "Failed to encrypt cardholder name",
						Details:   []string{erroe.Error()},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		updateFields = append(updateFields, fmt.Sprintf("encrypted_cardholder_name = $%d", argIndex))
		args = append(args, encryptedName)
		argIndex++
	}

	if req.CardNumber != "" {
		if !ValidateCardNumber(req.CardNumber) {
			return &pb.UpdateCardDataResponse{
				StatusCode: 400,
				Message:    "Invalid card number",
				Result: &pb.UpdateCardDataResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      400,
						Message:   "Invalid card number format",
						Details:   []string{"Card number failed validation"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		encryptedNumber, error := EncryptData(req.CardNumber)
		if error != nil {
			return &pb.UpdateCardDataResponse{
				StatusCode: 500,
				Message:    "Encryption failed",
				Result: &pb.UpdateCardDataResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      500,
						Message:   "Failed to encrypt card number",
						Details:   []string{error.Error()},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		maskedNumber := MaskCardNumber(req.CardNumber)
		updateFields = append(updateFields, fmt.Sprintf("encrypted_card_number = $%d", argIndex))
		args = append(args, encryptedNumber)
		argIndex++
		updateFields = append(updateFields, fmt.Sprintf("masked_card_number = $%d", argIndex))
		args = append(args, maskedNumber)
		argIndex++
	}

	if req.Cvv != "" {
		encryptedCVV, erre := EncryptData(req.Cvv)
		if erre != nil {
			return &pb.UpdateCardDataResponse{
				StatusCode: 500,
				Message:    "Encryption failed",
				Result: &pb.UpdateCardDataResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      500,
						Message:   "Failed to encrypt CVV",
						Details:   []string{erre.Error()},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		updateFields = append(updateFields, fmt.Sprintf("encrypted_cvv = $%d", argIndex))
		args = append(args, encryptedCVV)
		argIndex++
	}

	if req.ExpirationDate != "" {
		updateFields = append(updateFields, fmt.Sprintf("expiration_date = $%d", argIndex))
		args = append(args, req.ExpirationDate)
		argIndex++
	}

	if req.Status != "" {
		updateFields = append(updateFields, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, req.Status)
		argIndex++
	}

	if req.CardType != "" {
		updateFields = append(updateFields, fmt.Sprintf("card_type = $%d", argIndex))
		args = append(args, req.CardType)
		argIndex++
	}

	if req.Balance != "" {
		updateFields = append(updateFields, fmt.Sprintf("balance = $%d", argIndex))
		args = append(args, req.Balance)
		argIndex++
	}

	if len(updateFields) == 0 {
		return &pb.UpdateCardDataResponse{
			StatusCode: 400,
			Message:    "No fields to update",
			Result: &pb.UpdateCardDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "At least one field must be provided for update",
					Details:   []string{"No update fields specified"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++
	args = append(args, req.Id)
	query := fmt.Sprintf("UPDATE cards SET %s WHERE id = $%d",
		string(updateFields[0]), argIndex)
	for i := 1; i < len(updateFields); i++ {
		query = fmt.Sprintf("%s, %s", query, updateFields[i])
	}
	_, err = s.db.GetDB().Exec(query, args...)
	if err != nil {
		return &pb.UpdateCardDataResponse{
			StatusCode: 500,
			Message:    "Update failed",
			Result: &pb.UpdateCardDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to update card",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	updatedCard, err := s.GetCardByID(ctx, &pb.GetCardByIDRequest{Id: req.Id})
	if err != nil || updatedCard.Result.(*pb.GetCardByIDResponse_Card) == nil {
		return &pb.UpdateCardDataResponse{
			StatusCode: 500,
			Message:    "Failed to fetch updated card",
			Result: &pb.UpdateCardDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Card updated but failed to retrieve",
					Details:   []string{"Internal error"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	return &pb.UpdateCardDataResponse{
		StatusCode: 200,
		Message:    "Card updated successfully",
		Result: &pb.UpdateCardDataResponse_Card{
			Card: updatedCard.Result.(*pb.GetCardByIDResponse_Card).Card,
		},
	}, nil
}

func (s *CardServer) DeleteCardData(ctx context.Context, req *pb.DeleteCardDataRequest) (*pb.DeleteCardDataResponse, error) {
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.DeleteCardDataResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.DeleteCardDataResponse_Error{
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
		return &pb.DeleteCardDataResponse{
			StatusCode: 400,
			Message:    "Card ID is required",
			Result: &pb.DeleteCardDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Card ID cannot be empty",
					Details:   []string{"ID field is required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	cardResponse, err := s.GetCardByID(ctx, &pb.GetCardByIDRequest{Id: req.Id})
	if err != nil || cardResponse.StatusCode != 200 {
		return &pb.DeleteCardDataResponse{
			StatusCode: 404,
			Message:    "Card not found",
			Result: &pb.DeleteCardDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      404,
					Message:   "Card not found or access denied",
					Details:   []string{"No card found with the provided ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	query := "DELETE FROM cards WHERE id = $1 AND user_id = $2"
	result, err := s.db.GetDB().Exec(query, req.Id, claims.UserID)
	if err != nil {
		return &pb.DeleteCardDataResponse{
			StatusCode: 500,
			Message:    "Delete failed",
			Result: &pb.DeleteCardDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to delete card",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return &pb.DeleteCardDataResponse{
			StatusCode: 404,
			Message:    "Card not found",
			Result: &pb.DeleteCardDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      404,
					Message:   "Card not found or already deleted",
					Details:   []string{"No rows affected"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	return &pb.DeleteCardDataResponse{
		StatusCode: 200,
		Message:    "Card deleted successfully",
		Result: &pb.DeleteCardDataResponse_Card{
			Card: cardResponse.Result.(*pb.GetCardByIDResponse_Card).Card,
		},
	}, nil
}
func (s *CardServer) UpdateCardStatus(ctx context.Context, req *pb.UpdateCardStatusRequest) (*pb.UpdateCardStatusResponse, error) {
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.UpdateCardStatusResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.UpdateCardStatusResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.Id == "" || req.Status == "" {
		return &pb.UpdateCardStatusResponse{
			StatusCode: 400,
			Message:    "Card ID and status are required",
			Result: &pb.UpdateCardStatusResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Both ID and status fields are required",
					Details:   []string{"Missing required fields"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	query := "UPDATE cards SET status = $1, updated_at = $2 WHERE id = $3 AND user_id = $4"
	result, err := s.db.GetDB().Exec(query, req.Status, time.Now(), req.Id, claims.UserID)
	if err != nil {
		return &pb.UpdateCardStatusResponse{
			StatusCode: 500,
			Message:    "Update failed",
			Result: &pb.UpdateCardStatusResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to update card status",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return &pb.UpdateCardStatusResponse{
			StatusCode: 404,
			Message:    "Card not found",
			Result: &pb.UpdateCardStatusResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      404,
					Message:   "Card not found or access denied",
					Details:   []string{"No rows affected"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	updatedCard, err := s.GetCardByID(ctx, &pb.GetCardByIDRequest{Id: req.Id})
	if err != nil || updatedCard.StatusCode != 200 {
		return &pb.UpdateCardStatusResponse{
			StatusCode: 500,
			Message:    "Failed to fetch updated card",
			Result: &pb.UpdateCardStatusResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Status updated but failed to retrieve card",
					Details:   []string{"Internal error"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	return &pb.UpdateCardStatusResponse{
		StatusCode: 200,
		Message:    "Card status updated successfully",
		Result: &pb.UpdateCardStatusResponse_Card{
			Card: updatedCard.Result.(*pb.GetCardByIDResponse_Card).Card,
		},
	}, nil
}

func (s *CardServer) GetCardByUserID(ctx context.Context, req *pb.GetCardByUserIDRequest) (*pb.GetCardByUserIDResponse, error) {
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetCardByUserIDResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetCardByUserIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	if req.UserId != claims.UserID {
		return &pb.GetCardByUserIDResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.GetCardByUserIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      403,
					Message:   "Access denied",
					Details:   []string{"You can only access your own cards"},
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
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := s.db.GetDB().Query(query, claims.UserID)
	if err != nil {
		return &pb.GetCardByUserIDResponse{
			StatusCode: 500,
			Message:    "Database error",
			Result: &pb.GetCardByUserIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to query cards",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	defer rows.Close()
	var cards []*pb.Card
	for rows.Next() {
		var card pb.Card
		var encryptedName, encryptedNumber, encryptedCVV string
		var createdAt, updatedAt time.Time

		error := rows.Scan(
			&card.Id, &card.UserId, &encryptedName, &encryptedNumber,
			&encryptedCVV, &card.CardNumber, &card.ExpirationDate,
			&card.CardType, &card.Status, &card.Balance,
			&createdAt, &updatedAt,
		)
		if error != nil {
			log.Printf("Error scanning card row: %v", err)
			continue
		}
		cardHolderName, error := DecryptData(encryptedName)
		if error != nil {
			log.Printf("Failed to decrypt cardholder name: %v", err)
			cardHolderName = "[ENCRYPTED]"
		}
		card.CardHolderName = cardHolderName
		card.Cvv = "***"
		card.CreateAt = createdAt.Format(time.RFC3339)
		card.UpdateAt = updatedAt.Format(time.RFC3339)
		cards = append(cards, &card)
	}

	if err = rows.Err(); err != nil {
		return &pb.GetCardByUserIDResponse{
			StatusCode: 500,
			Message:    "Database error",
			Result: &pb.GetCardByUserIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Error iterating over cards",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	if len(cards) == 0 {
		return &pb.GetCardByUserIDResponse{
			StatusCode: 404,
			Message:    "No cards found",
			Result: &pb.GetCardByUserIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      404,
					Message:   "No cards found for this user",
					Details:   []string{"User has no cards"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	return &pb.GetCardByUserIDResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Found %d cards", len(cards)),
		Result: &pb.GetCardByUserIDResponse_Card{
			Card: cards[0], // Return the most recent card
		},
	}, nil
}

func (s *CardServer) SearchCard(ctx context.Context, req *pb.SearchCardRequest) (*pb.SearchCardResponse, error) {
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.SearchCardResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.SearchCardResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	if req.CardNumber == "" {
		return &pb.SearchCardResponse{
			StatusCode: 400,
			Message:    "Card number is required",
			Result: &pb.SearchCardResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Card number cannot be empty",
					Details:   []string{"Card number field is required"},
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
		WHERE masked_card_number LIKE $1 AND user_id = $2
	`
	rows, err := s.db.GetDB().Query(query, "%"+req.CardNumber+"%", claims.UserID)
	if err != nil {
		return &pb.SearchCardResponse{
			StatusCode: 500,
			Message:    "Database error",
			Result: &pb.SearchCardResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to search cards",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	defer rows.Close()
	var card *pb.Card
	for rows.Next() {
		var c pb.Card
		var encryptedName, encryptedNumber, encryptedCVV string
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&c.Id, &c.UserId, &encryptedName, &encryptedNumber,
			&encryptedCVV, &c.CardNumber, &c.ExpirationDate,
			&c.CardType, &c.Status, &c.Balance,
			&createdAt, &updatedAt,
		)
		if err != nil {
			log.Printf("Error scanning card row: %v", err)
			continue
		}
		cardHolderName, err := DecryptData(encryptedName)
		if err != nil {
			log.Printf("Failed to decrypt cardholder name: %v", err)
			cardHolderName = "[ENCRYPTED]"
		}
		c.CardHolderName = cardHolderName
		c.Cvv = "***"
		c.CreateAt = createdAt.Format(time.RFC3339)
		c.UpdateAt = updatedAt.Format(time.RFC3339)

		card = &c
		break
	}
	if card == nil {
		return &pb.SearchCardResponse{
			StatusCode: 404,
			Message:    "Card not found",
			Result: &pb.SearchCardResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      404,
					Message:   "No card found matching the search criteria",
					Details:   []string{"Card not found"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.SearchCardResponse{
		StatusCode: 200,
		Message:    "Card found",
		Result: &pb.SearchCardResponse_Cards{
			Cards: card,
		},
	}, nil
}
func (s *CardServer) FindCardByStatus(ctx context.Context, req *pb.FindCardByStatusRequest) (*pb.FindCardByStatusResponse, error) {
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.FindCardByStatusResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.FindCardByStatusResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	if req.Status == "" {
		return &pb.FindCardByStatusResponse{
			StatusCode: 400,
			Message:    "Status is required",
			Result: &pb.FindCardByStatusResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Status field cannot be empty",
					Details:   []string{"Status field is required"},
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
		WHERE status = $1 AND user_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`
	var card pb.Card
	var encryptedName, encryptedNumber, encryptedCVV string
	var createdAt, updatedAt time.Time
	err = s.db.GetDB().QueryRow(query, req.Status, claims.UserID).Scan(
		&card.Id, &card.UserId, &encryptedName, &encryptedNumber,
		&encryptedCVV, &card.CardNumber, &card.ExpirationDate,
		&card.CardType, &card.Status, &card.Balance,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.FindCardByStatusResponse{
				StatusCode: 404,
				Message:    "No cards found with the specified status",
				Result: &pb.FindCardByStatusResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "No cards found",
						Details:   []string{"No cards found with the specified status"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return &pb.FindCardByStatusResponse{
			StatusCode: 500,
			Message:    "Database error",
			Result: &pb.FindCardByStatusResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to query cards",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	cardHolderName, err := DecryptData(encryptedName)
	if err != nil {
		log.Printf("Failed to decrypt cardholder name: %v", err)
		cardHolderName = "[ENCRYPTED]"
	}
	card.CardHolderName = cardHolderName
	card.Cvv = "***"
	card.CreateAt = createdAt.Format(time.RFC3339)
	card.UpdateAt = updatedAt.Format(time.RFC3339)
	return &pb.FindCardByStatusResponse{
		StatusCode: 200,
		Message:    "Card found",
		Result: &pb.FindCardByStatusResponse_Cards{
			Cards: &card,
		},
	}, nil
}

func (s *CardServer) StripePayment(ctx context.Context, req *pb.StripePaymentRequest) (*pb.StripePaymentResponse, error) {
	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.StripePaymentResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.StripePaymentResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	if req.CardId == "" || req.Amount == "" || req.Currency == "" {
		return &pb.StripePaymentResponse{
			StatusCode: 400,
			Message:    "Missing required fields",
			Result: &pb.StripePaymentResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Card ID, amount, and currency are required",
					Details:   []string{"Missing required payment fields"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	cardResponse, err := s.GetCardByID(ctx, &pb.GetCardByIDRequest{Id: req.CardId})
	if err != nil || cardResponse.StatusCode != 200 {
		return &pb.StripePaymentResponse{
			StatusCode: 404,
			Message:    "Card not found",
			Result: &pb.StripePaymentResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      404,
					Message:   "Card not found or access denied",
					Details:   []string{"Invalid card ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	card := cardResponse.Result.(*pb.GetCardByIDResponse_Card).Card
	if card.Status != "ACTIVE" {
		return &pb.StripePaymentResponse{
			StatusCode: 400,
			Message:    "Card is not active",
			Result: &pb.StripePaymentResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Card must be active to process payments",
					Details:   []string{"Inactive card"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	amountFloat, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return &pb.StripePaymentResponse{
			StatusCode: 400,
			Message:    "Invalid amount format",
			Result: &pb.StripePaymentResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Amount must be a valid number",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	amountCents := int64(amountFloat * 100)
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		return &pb.StripePaymentResponse{
			StatusCode: 200,
			Message:    "Payment processed successfully (simulated)",
			Result: &pb.StripePaymentResponse_Card{
				Card: card,
			},
		}, nil
	}
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(req.Currency),
		Metadata: map[string]string{
			"card_id": req.CardId,
			"user_id": claims.UserID,
		},
	}
	_, err = paymentintent.New(params)
	if err != nil {
		return &pb.StripePaymentResponse{
			StatusCode: 500,
			Message:    "Payment processing failed",
			Result: &pb.StripePaymentResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Stripe payment failed",
					Details:   []string{err.Error()},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.StripePaymentResponse{
		StatusCode: 200,
		Message:    "Payment processed successfully",
		Result: &pb.StripePaymentResponse_Card{
			Card: card,
		},
	}, nil
}
