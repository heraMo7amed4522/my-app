package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	pb "chat-services/proto"

	"github.com/OneSignal/onesignal-go-api"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	db         *Database
	userClient UserServiceClient
}

func NewChatServer() *ChatServer {
	userClient, err := NewUserServiceClient()
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}

	return &ChatServer{
		db:         NewDatabase(),
		userClient: userClient,
	}
}
func (s *ChatServer) extractTokenFromContext(ctx context.Context) (string, error) {
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

func (s *ChatServer) validateRequest(ctx context.Context) (*TokenClaims, error) {
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
func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
func (s *ChatServer) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	log.Printf("SendMessage called by sender: %s", req.SenderId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.SendMessageResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.SendMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate sender ID matches token
	if req.SenderId != claims.UserID {
		return &pb.SendMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.SendMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Sender ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if !ValidateUserID(req.SenderId) {
		return &pb.SendMessageResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.SendMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid sender ID format",
					Details:   []string{"Sender ID must be a valid UUID"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if !ValidateMessageContent(req.Content) {
		return &pb.SendMessageResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.SendMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid message content",
					Details:   []string{"Message content cannot be empty and must be less than 4000 characters"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.IsGroup {
		if !ValidateUserID(req.GroupId) {
			return &pb.SendMessageResponse{
				StatusCode: 400,
				Message:    "Bad Request",
				Result: &pb.SendMessageResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      400,
						Message:   "Invalid group ID format",
						Details:   []string{"Group ID must be a valid UUID"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
	} else {
		if !ValidateUserID(req.ReceiverId) {
			return &pb.SendMessageResponse{
				StatusCode: 400,
				Message:    "Bad Request",
				Result: &pb.SendMessageResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      400,
						Message:   "Invalid receiver ID format",
						Details:   []string{"Receiver ID must be a valid UUID"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
	}

	// Generate message ID
	messageID := uuid.New().String()
	now := time.Now()

	// Insert message into database
	query := `
		INSERT INTO chat_messages (
			id, sender_id, receiver_id, group_id, content, message_type, 
			is_group, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`

	var savedID string
	var createdAt time.Time
	err = s.db.DB.QueryRow(
		query,
		messageID,
		req.SenderId,
		nullString(req.ReceiverId),
		nullString(req.GroupId),
		req.Content,
		int32(req.Type),
		req.IsGroup,
		0, // SENT status
		now,
		now,
	).Scan(&savedID, &createdAt)

	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.SendMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.SendMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to save message",
					Details:   []string{"Database operation failed"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Create response message
	savedMessage := &pb.ChatMessage{
		MessageId:   savedID,
		SenderId:    req.SenderId,
		ReceiverId:  req.ReceiverId,
		GroupId:     req.GroupId,
		Content:     req.Content,
		Type:        req.Type,
		Timestamp:   timestamppb.New(time.Now()),
		IsGroup:     req.IsGroup,
		IsRead:      false,
		IsEdited:    false,
		Status:      pb.MessageStatus_SENT,
		ScheduledAt: timestamppb.New(time.Now()),
		EditedAt:    timestamppb.New(time.Now()),
	}

	return &pb.SendMessageResponse{
		StatusCode: 200,
		Message:    "Message sent successfully",
		Result: &pb.SendMessageResponse_SavedMessage{
			SavedMessage: savedMessage,
		},
	}, nil
}

func (s *ChatServer) EditMessage(ctx context.Context, req *pb.EditMessageRequest) (*pb.EditMessageResponse, error) {
	log.Printf("EditMessage called for message: %s by user: %s", req.MessageId, req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.EditMessageResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.EditMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.EditMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.EditMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if !ValidateUserID(req.MessageId) || !ValidateMessageContent(req.NewContent) {
		return &pb.EditMessageResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.EditMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid input",
					Details:   []string{"Invalid message ID or content"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if user owns the message
	var senderID string
	err = s.db.DB.QueryRow("SELECT sender_id FROM chat_messages WHERE id = $1", req.MessageId).Scan(&senderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.EditMessageResponse{
				StatusCode: 404,
				Message:    "Message not found",
				Result: &pb.EditMessageResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      404,
						Message:   "Message not found",
						Details:   []string{"Message does not exist"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
		return &pb.EditMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.EditMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Database error",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if senderID != req.UserId {
		return &pb.EditMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.EditMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Cannot edit message from another user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Update message
	now := time.Now()
	_, err = s.db.DB.Exec(
		"UPDATE chat_messages SET content = $1, is_edited = true, edited_at = $2, updated_at = $3 WHERE id = $4",
		req.NewContent, now, now, req.MessageId,
	)
	if err != nil {
		return &pb.EditMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.EditMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to update message",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Fetch updated message
	updatedMessage, err := s.getMessageByID(req.MessageId)
	if err != nil {
		return &pb.EditMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.EditMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to fetch updated message",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.EditMessageResponse{
		StatusCode: 200,
		Message:    "Message updated successfully",
		Result: &pb.EditMessageResponse_UpdatedMessage{
			UpdatedMessage: updatedMessage,
		},
	}, nil
}

func (s *ChatServer) DeleteMessage(ctx context.Context, req *pb.DeleteMessageRequest) (*pb.DeleteMessageResponse, error) {
	log.Printf("DeleteMessage called for message: %s by user: %s", req.MessageId, req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.DeleteMessageResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.DeleteMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.DeleteMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.DeleteMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if user owns the message
	var senderID string
	err = s.db.DB.QueryRow("SELECT sender_id FROM chat_messages WHERE id = $1", req.MessageId).Scan(&senderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.DeleteMessageResponse{
				StatusCode: 404,
				Message:    "Message not found",
				Result: &pb.DeleteMessageResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      404,
						Message:   "Message not found",
						Details:   []string{"Message does not exist"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
		return &pb.DeleteMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.DeleteMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Database error",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if senderID != req.UserId {
		return &pb.DeleteMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.DeleteMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Cannot delete message from another user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Soft delete message
	_, err = s.db.DB.Exec(
		"UPDATE chat_messages SET content = '[Message deleted]', updated_at = $1 WHERE id = $2",
		time.Now(), req.MessageId,
	)
	if err != nil {
		return &pb.DeleteMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.DeleteMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to delete message",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.DeleteMessageResponse{
		StatusCode: 200,
		Message:    "Message deleted successfully",
		Result: &pb.DeleteMessageResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Message deleted",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (s *ChatServer) GetChatHistory(ctx context.Context, req *pb.GetChatHistoryRequest) (*pb.GetChatHistoryResponse, error) {
	log.Printf("GetChatHistory called by user: %s", req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetChatHistoryResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetChatHistoryResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.GetChatHistoryResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.GetChatHistoryResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	var query string
	var args []interface{}

	if req.IsGroup {
		query = `
			SELECT id, sender_id, COALESCE(receiver_id, ''), COALESCE(group_id, ''), content, 
			       message_type, is_group, is_read, is_edited, status, created_at, updated_at
			FROM chat_messages 
			WHERE group_id = $1 AND is_group = true
			ORDER BY created_at DESC 
			LIMIT $2 OFFSET $3`
		args = []interface{}{req.PeerId, req.Limit, req.Offset}
	} else {
		query = `
			SELECT id, sender_id, COALESCE(receiver_id, ''), COALESCE(group_id, ''), content, 
			       message_type, is_group, is_read, is_edited, status, created_at, updated_at
			FROM chat_messages 
			WHERE ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1))
			  AND is_group = false
			ORDER BY created_at DESC 
			LIMIT $3 OFFSET $4`
		args = []interface{}{req.UserId, req.PeerId, req.Limit, req.Offset}
	}

	rows, err := s.db.DB.Query(query, args...)
	if err != nil {
		return &pb.GetChatHistoryResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.GetChatHistoryResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to fetch chat history",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	defer rows.Close()

	var messages []*pb.ChatMessage
	for rows.Next() {
		var msg pb.ChatMessage
		var receiverID, groupID sql.NullString
		var createdAt, updatedAt time.Time
		var messageType, status int32

		err := rows.Scan(
			&msg.MessageId, &msg.SenderId, &receiverID, &groupID, &msg.Content,
			&messageType, &msg.IsGroup, &msg.IsRead, &msg.IsEdited, &status,
			&createdAt, &updatedAt,
		)
		if err != nil {
			continue
		}

		msg.ReceiverId = receiverID.String
		msg.GroupId = groupID.String
		msg.Type = pb.MessageType(messageType)
		msg.Status = pb.MessageStatus(status)
		msg.Timestamp = timestamppb.New(createdAt)
		msg.EditedAt = timestamppb.New(updatedAt)

		messages = append(messages, &msg)
	}

	return &pb.GetChatHistoryResponse{
		StatusCode: 200,
		Message:    "Chat history retrieved successfully",
		Result: &pb.GetChatHistoryResponse_Messages{
			Messages: &pb.ChatMessageList{
				Messages: messages,
			},
		},
	}, nil
}

func (s *ChatServer) MarkAsRead(ctx context.Context, req *pb.ReadReceiptRequest) (*pb.ReadReceiptResponse, error) {
	log.Printf("MarkAsRead called by user: %s", req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.ReadReceiptResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.ReadReceiptResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.ReadReceiptResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.ReadReceiptResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	var query string
	var args []interface{}

	if req.IsGroup {
		query = "UPDATE chat_messages SET is_read = true, read_at = $1 WHERE group_id = $2 AND sender_id != $3"
		args = []interface{}{time.Now(), req.PeerId, req.UserId}
	} else {
		query = "UPDATE chat_messages SET is_read = true, read_at = $1 WHERE sender_id = $2 AND receiver_id = $3"
		args = []interface{}{time.Now(), req.PeerId, req.UserId}
	}

	_, err = s.db.DB.Exec(query, args...)
	if err != nil {
		return &pb.ReadReceiptResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.ReadReceiptResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to mark messages as read",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.ReadReceiptResponse{
		StatusCode: 200,
		Message:    "Messages marked as read",
		Result: &pb.ReadReceiptResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Messages marked as read",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (s *ChatServer) SearchMessages(ctx context.Context, req *pb.SearchMessagesRequest) (*pb.SearchMessagesResponse, error) {
	log.Printf("SearchMessages called by user: %s", req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.SearchMessagesResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.SearchMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.SearchMessagesResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.SearchMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	var query string
	var args []interface{}

	if req.IsGroup {
		query = `
			SELECT id, sender_id, COALESCE(receiver_id, ''), COALESCE(group_id, ''), content, 
			       message_type, is_group, is_read, is_edited, status, created_at, updated_at
			FROM chat_messages 
			WHERE group_id = $1 AND content ILIKE $2
			ORDER BY created_at DESC 
			LIMIT $3 OFFSET $4`
		args = []interface{}{req.GroupId, "%" + req.Query + "%", req.Limit, req.Offset}
	} else {
		query = `
			SELECT id, sender_id, COALESCE(receiver_id, ''), COALESCE(group_id, ''), content, 
			       message_type, is_group, is_read, is_edited, status, created_at, updated_at
			FROM chat_messages 
			WHERE ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1))
			  AND content ILIKE $3 AND is_group = false
			ORDER BY created_at DESC 
			LIMIT $4 OFFSET $5`
		args = []interface{}{req.UserId, req.PeerId, "%" + req.Query + "%", req.Limit, req.Offset}
	}

	rows, err := s.db.DB.Query(query, args...)
	if err != nil {
		return &pb.SearchMessagesResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.SearchMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to search messages",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	defer rows.Close()

	var messages []*pb.ChatMessage
	for rows.Next() {
		var msg pb.ChatMessage
		var receiverID, groupID sql.NullString
		var createdAt, updatedAt time.Time
		var messageType, status int32

		err := rows.Scan(
			&msg.MessageId, &msg.SenderId, &receiverID, &groupID, &msg.Content,
			&messageType, &msg.IsGroup, &msg.IsRead, &msg.IsEdited, &status,
			&createdAt, &updatedAt,
		)
		if err != nil {
			continue
		}

		msg.ReceiverId = receiverID.String
		msg.GroupId = groupID.String
		msg.Type = pb.MessageType(messageType)
		msg.Status = pb.MessageStatus(status)
		msg.Timestamp = timestamppb.New(createdAt)
		msg.EditedAt = timestamppb.New(updatedAt)

		messages = append(messages, &msg)
	}

	return &pb.SearchMessagesResponse{
		StatusCode: 200,
		Message:    "Search completed successfully",
		Result: &pb.SearchMessagesResponse_Messages{
			Messages: &pb.ChatMessageList{
				Messages: messages,
			},
		},
	}, nil
}

func (s *ChatServer) ForwardMessage(ctx context.Context, req *pb.ForwardMessageRequest) (*pb.ForwardMessageResponse, error) {
	log.Printf("ForwardMessage called by user: %s", req.SenderId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.ForwardMessageResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.ForwardMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.SenderId != claims.UserID {
		return &pb.ForwardMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.ForwardMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Sender ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Get original message
	originalMessage, err := s.getMessageByID(req.MessageId)
	if err != nil {
		return &pb.ForwardMessageResponse{
			StatusCode: 404,
			Message:    "Message not found",
			Result: &pb.ForwardMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      404,
					Message:   "Original message not found",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	var forwardedMessages []*pb.ChatMessage
	now := time.Now()

	// Forward to individual users
	for _, receiverID := range req.ReceiverIds {
		messageID := uuid.New().String()
		_, err = s.db.DB.Exec(`
			INSERT INTO chat_messages (
				id, sender_id, receiver_id, content, message_type, 
				is_group, status, created_at, updated_at, original_message_id, forward_count
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			messageID, req.SenderId, receiverID, originalMessage.Content, originalMessage.Type,
			false, 0, now, now, req.MessageId, originalMessage.ForwardCount+1,
		)
		if err != nil {
			continue
		}

		forwardedMsg := &pb.ChatMessage{
			MessageId:         messageID,
			SenderId:          req.SenderId,
			ReceiverId:        receiverID,
			Content:           originalMessage.Content,
			Type:              originalMessage.Type,
			Timestamp:         timestamppb.New(now),
			IsGroup:           false,
			Status:            pb.MessageStatus_SENT,
			OriginalMessageId: req.MessageId,
			ForwardCount:      originalMessage.ForwardCount + 1,
		}
		forwardedMessages = append(forwardedMessages, forwardedMsg)
	}

	// Forward to groups
	for _, groupID := range req.GroupIds {
		messageID := uuid.New().String()
		_, err = s.db.DB.Exec(`
			INSERT INTO chat_messages (
				id, sender_id, group_id, content, message_type, 
				is_group, status, created_at, updated_at, original_message_id, forward_count
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			messageID, req.SenderId, groupID, originalMessage.Content, originalMessage.Type,
			true, 0, now, now, req.MessageId, originalMessage.ForwardCount+1,
		)
		if err != nil {
			continue
		}

		forwardedMsg := &pb.ChatMessage{
			MessageId:         messageID,
			SenderId:          req.SenderId,
			GroupId:           groupID,
			Content:           originalMessage.Content,
			Type:              originalMessage.Type,
			Timestamp:         timestamppb.New(now),
			IsGroup:           true,
			Status:            pb.MessageStatus_SENT,
			OriginalMessageId: req.MessageId,
			ForwardCount:      originalMessage.ForwardCount + 1,
		}
		forwardedMessages = append(forwardedMessages, forwardedMsg)
	}

	return &pb.ForwardMessageResponse{
		StatusCode: 200,
		Message:    "Message forwarded successfully",
		Result: &pb.ForwardMessageResponse_ForwardedMessages{
			ForwardedMessages: &pb.ChatMessageList{
				Messages: forwardedMessages,
			},
		},
	}, nil
}

func (s *ChatServer) PinMessage(ctx context.Context, req *pb.PinMessageRequest) (*pb.PinMessageResponse, error) {
	log.Printf("PinMessage called for message: %s by user: %s", req.MessageId, req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.PinMessageResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.PinMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.PinMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.PinMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Update message to set pinned status
	now := time.Now()
	_, err = s.db.DB.Exec(
		"UPDATE chat_messages SET is_pinned = true, pinned_at = $1, pinned_by = $2 WHERE id = $3",
		now, req.UserId, req.MessageId,
	)
	if err != nil {
		return &pb.PinMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.PinMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to pin message",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.PinMessageResponse{
		StatusCode: 200,
		Message:    "Message pinned successfully",
		Result: &pb.PinMessageResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Message pinned",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (s *ChatServer) UnpinMessage(ctx context.Context, req *pb.UnpinMessageRequest) (*pb.UnpinMessageResponse, error) {
	log.Printf("UnpinMessage called for message: %s by user: %s", req.MessageId, req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.UnpinMessageResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.UnpinMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.UnpinMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.UnpinMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Update message to remove pinned status
	_, err = s.db.DB.Exec(
		"UPDATE chat_messages SET is_pinned = false, pinned_at = NULL, pinned_by = NULL WHERE id = $1",
		req.MessageId,
	)
	if err != nil {
		return &pb.UnpinMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.UnpinMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to unpin message",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.UnpinMessageResponse{
		StatusCode: 200,
		Message:    "Message unpinned successfully",
		Result: &pb.UnpinMessageResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Message unpinned",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (s *ChatServer) GetPinnedMessages(ctx context.Context, req *pb.GetPinnedMessagesRequest) (*pb.GetPinnedMessagesResponse, error) {
	log.Printf("GetPinnedMessages called for chat: %s", req.ChatId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetPinnedMessagesResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetPinnedMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	var query string
	var args []interface{}

	if req.IsGroup {
		query = `
			SELECT id, sender_id, COALESCE(receiver_id, ''), COALESCE(group_id, ''), content, 
			       message_type, is_group, is_read, is_edited, status, created_at, updated_at,
			       pinned_at, pinned_by
			FROM chat_messages 
			WHERE group_id = $1 AND is_pinned = true
			ORDER BY pinned_at DESC`
		args = []interface{}{req.ChatId}
	} else {
		query = `
			SELECT id, sender_id, COALESCE(receiver_id, ''), COALESCE(group_id, ''), content, 
			       message_type, is_group, is_read, is_edited, status, created_at, updated_at,
			       pinned_at, pinned_by
			FROM chat_messages 
			WHERE ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1))
			  AND is_pinned = true AND is_group = false
			ORDER BY pinned_at DESC`
		args = []interface{}{claims.UserID, req.ChatId}
	}

	rows, err := s.db.DB.Query(query, args...)
	if err != nil {
		return &pb.GetPinnedMessagesResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.GetPinnedMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to fetch pinned messages",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	defer rows.Close()

	var messages []*pb.ChatMessage
	for rows.Next() {
		var msg pb.ChatMessage
		var receiverID, groupID sql.NullString
		var createdAt, updatedAt time.Time
		var pinnedAt sql.NullTime
		var pinnedBy sql.NullString
		var messageType, status int32

		err := rows.Scan(
			&msg.MessageId, &msg.SenderId, &receiverID, &groupID, &msg.Content,
			&messageType, &msg.IsGroup, &msg.IsRead, &msg.IsEdited, &status,
			&createdAt, &updatedAt, &pinnedAt, &pinnedBy,
		)
		if err != nil {
			continue
		}

		msg.ReceiverId = receiverID.String
		msg.GroupId = groupID.String
		msg.Type = pb.MessageType(messageType)
		msg.Status = pb.MessageStatus(status)
		msg.Timestamp = timestamppb.New(createdAt)
		msg.EditedAt = timestamppb.New(updatedAt)
		msg.IsPinned = true
		if pinnedAt.Valid {
			msg.PinnedAt = timestamppb.New(pinnedAt.Time)
		}
		msg.PinnedBy = pinnedBy.String

		messages = append(messages, &msg)
	}

	return &pb.GetPinnedMessagesResponse{
		StatusCode: 200,
		Message:    "Pinned messages retrieved successfully",
		Result: &pb.GetPinnedMessagesResponse_PinnedMessages{
			PinnedMessages: &pb.ChatMessageList{
				Messages: messages,
			},
		},
	}, nil
}

func (s *ChatServer) getMessageByID(messageID string) (*pb.ChatMessage, error) {
	var msg pb.ChatMessage
	var receiverID, groupID sql.NullString
	var createdAt, updatedAt time.Time
	var messageType, status int32

	err := s.db.DB.QueryRow(`
		SELECT id, sender_id, COALESCE(receiver_id, ''), COALESCE(group_id, ''), content, 
		       message_type, is_group, is_read, is_edited, status, created_at, updated_at,
		       COALESCE(forward_count, 0)
		FROM chat_messages WHERE id = $1`, messageID).Scan(
		&msg.MessageId, &msg.SenderId, &receiverID, &groupID, &msg.Content,
		&messageType, &msg.IsGroup, &msg.IsRead, &msg.IsEdited, &status,
		&createdAt, &updatedAt, &msg.ForwardCount,
	)
	if err != nil {
		return nil, err
	}

	msg.ReceiverId = receiverID.String
	msg.GroupId = groupID.String
	msg.Type = pb.MessageType(messageType)
	msg.Status = pb.MessageStatus(status)
	msg.Timestamp = timestamppb.New(createdAt)
	msg.EditedAt = timestamppb.New(updatedAt)

	return &msg, nil
}

func (s *ChatServer) AddLikeMessage(ctx context.Context, req *pb.AddLikeMessageRequest) (*pb.AddLikeMessageResponse, error) {
	log.Printf("AddLikeMessage called for message: %s by user: %s", req.MessageId, req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.AddLikeMessageResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.AddLikeMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.AddLikeMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.AddLikeMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if message exists
	var messageExists bool
	err = s.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM chat_messages WHERE id = $1)", req.MessageId).Scan(&messageExists)
	if err != nil || !messageExists {
		return &pb.AddLikeMessageResponse{
			StatusCode: 404,
			Message:    "Message not found",
			Result: &pb.AddLikeMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      404,
					Message:   "Message not found",
					Details:   []string{"Message does not exist"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Insert or update like
	_, err = s.db.DB.Exec(`
		INSERT INTO message_reactions (message_id, user_id, reaction_type, created_at) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (message_id, user_id) 
		DO UPDATE SET reaction_type = $3, created_at = $4`,
		req.MessageId, req.UserId, int32(pb.ReactionType_LIKE), time.Now())

	if err != nil {
		return &pb.AddLikeMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.AddLikeMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to add like",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.AddLikeMessageResponse{
		StatusCode: 200,
		Message:    "Like added successfully",
		Result: &pb.AddLikeMessageResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Like added",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (s *ChatServer) UpdateLikedMessage(ctx context.Context, req *pb.UpdateLikedMessageRequest) (*pb.UpdateLikedMessageResponse, error) {
	log.Printf("UpdateLikedMessage called for message: %s by user: %s", req.MessageId, req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.UpdateLikedMessageResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.UpdateLikedMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.UpdateLikedMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.UpdateLikedMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Update reaction
	result, err := s.db.DB.Exec(`
		UPDATE message_reactions 
		SET reaction_type = $1, created_at = $2 
		WHERE message_id = $3 AND user_id = $4`,
		int32(req.ReactionType), time.Now(), req.MessageId, req.UserId)

	if err != nil {
		return &pb.UpdateLikedMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.UpdateLikedMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to update reaction",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &pb.UpdateLikedMessageResponse{
			StatusCode: 404,
			Message:    "Reaction not found",
			Result: &pb.UpdateLikedMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      404,
					Message:   "Reaction not found",
					Details:   []string{"No existing reaction to update"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.UpdateLikedMessageResponse{
		StatusCode: 200,
		Message:    "Reaction updated successfully",
		Result: &pb.UpdateLikedMessageResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Reaction updated",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (s *ChatServer) GetLikedMessages(ctx context.Context, req *pb.GetLikedMessagesRequest) (*pb.GetLikedMessagesResponse, error) {
	log.Printf("GetLikedMessages called by user: %s", req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetLikedMessagesResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetLikedMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.GetLikedMessagesResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.GetLikedMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	query := `
		SELECT cm.id, cm.sender_id, COALESCE(cm.receiver_id, ''), COALESCE(cm.group_id, ''), 
		       cm.content, cm.message_type, cm.is_group, cm.is_read, cm.is_edited, 
		       cm.status, cm.created_at, cm.updated_at
		FROM chat_messages cm
		INNER JOIN message_reactions mr ON cm.id = mr.message_id
		WHERE mr.user_id = $1
		ORDER BY mr.created_at DESC`

	rows, err := s.db.DB.Query(query, req.UserId)
	if err != nil {
		return &pb.GetLikedMessagesResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.GetLikedMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to fetch liked messages",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	defer rows.Close()

	var messages []*pb.ChatMessage
	for rows.Next() {
		var msg pb.ChatMessage
		var receiverID, groupID sql.NullString
		var createdAt, updatedAt time.Time
		var messageType, status int32

		err := rows.Scan(
			&msg.MessageId, &msg.SenderId, &receiverID, &groupID, &msg.Content,
			&messageType, &msg.IsGroup, &msg.IsRead, &msg.IsEdited, &status,
			&createdAt, &updatedAt,
		)
		if err != nil {
			continue
		}

		msg.ReceiverId = receiverID.String
		msg.GroupId = groupID.String
		msg.Type = pb.MessageType(messageType)
		msg.Status = pb.MessageStatus(status)
		msg.Timestamp = timestamppb.New(createdAt)
		msg.EditedAt = timestamppb.New(updatedAt)

		messages = append(messages, &msg)
	}

	return &pb.GetLikedMessagesResponse{
		StatusCode: 200,
		Message:    "Liked messages retrieved successfully",
		Result: &pb.GetLikedMessagesResponse_LikedMessages{
			LikedMessages: &pb.ChatMessageList{
				Messages: messages,
			},
		},
	}, nil
}

func (s *ChatServer) GetLastMessages(ctx context.Context, req *pb.GetLastMessagesRequest) (*pb.GetLastMessagesResponse, error) {
	log.Printf("GetLastMessages called by user: %s", req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetLastMessagesResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetLastMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.GetLastMessagesResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.GetLastMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Get last message from each conversation
	query := `
		WITH ranked_messages AS (
			SELECT cm.*, 
			       ROW_NUMBER() OVER (
			           PARTITION BY 
			               CASE 
			                   WHEN cm.is_group THEN cm.group_id
			                   ELSE LEAST(cm.sender_id, cm.receiver_id) || '_' || GREATEST(cm.sender_id, cm.receiver_id)
			               END
			           ORDER BY cm.created_at DESC
			       ) as rn
			FROM chat_messages cm
			WHERE cm.sender_id = $1 OR cm.receiver_id = $1 OR 
			      (cm.is_group AND cm.group_id IN (
			          SELECT group_id FROM group_members WHERE user_id = $1
			      ))
		)
		SELECT id, sender_id, COALESCE(receiver_id, ''), COALESCE(group_id, ''), 
		       content, message_type, is_group, is_read, is_edited, status, 
		       created_at, updated_at
		FROM ranked_messages 
		WHERE rn = 1
		ORDER BY created_at DESC`

	rows, err := s.db.DB.Query(query, req.UserId)
	if err != nil {
		return &pb.GetLastMessagesResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.GetLastMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to fetch last messages",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	defer rows.Close()

	var messages []*pb.ChatMessage
	for rows.Next() {
		var msg pb.ChatMessage
		var receiverID, groupID sql.NullString
		var createdAt, updatedAt time.Time
		var messageType, status int32

		err := rows.Scan(
			&msg.MessageId, &msg.SenderId, &receiverID, &groupID, &msg.Content,
			&messageType, &msg.IsGroup, &msg.IsRead, &msg.IsEdited, &status,
			&createdAt, &updatedAt,
		)
		if err != nil {
			continue
		}

		msg.ReceiverId = receiverID.String
		msg.GroupId = groupID.String
		msg.Type = pb.MessageType(messageType)
		msg.Status = pb.MessageStatus(status)
		msg.Timestamp = timestamppb.New(createdAt)
		msg.EditedAt = timestamppb.New(updatedAt)

		messages = append(messages, &msg)
	}

	return &pb.GetLastMessagesResponse{
		StatusCode: 200,
		Message:    "Last messages retrieved successfully",
		Result: &pb.GetLastMessagesResponse_LastMessages{
			LastMessages: &pb.ChatMessageList{
				Messages: messages,
			},
		},
	}, nil
}

func (s *ChatServer) GetUsersInGroup(ctx context.Context, req *pb.GetUsersInGroupRequest) (*pb.GetUsersInGroupResponse, error) {
	log.Printf("GetUsersInGroup called for group: %s", req.GroupId)

	_, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetUsersInGroupResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetUsersInGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if !ValidateUserID(req.GroupId) {
		return &pb.GetUsersInGroupResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.GetUsersInGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid group ID format",
					Details:   []string{"Group ID must be a valid UUID"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Get users from group_members table
	query := `SELECT user_id FROM group_members WHERE group_id = $1`
	rows, err := s.db.DB.Query(query, req.GroupId)
	if err != nil {
		return &pb.GetUsersInGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.GetUsersInGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to fetch group members",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	defer rows.Close()

	var users []*pb.UserInfo
	for rows.Next() {
		var userID string
		err := rows.Scan(&userID)
		if err != nil {
			continue
		}

		// Mock user data - in production, call user service
		user := &pb.UserInfo{
			Id:        userID,
			FullName:  "User " + userID[:8],
			Email:     userID + "@example.com",
			AvatarUrl: "",
			IsOnline:  false,
		}
		users = append(users, user)
	}

	return &pb.GetUsersInGroupResponse{
		StatusCode: 200,
		Message:    "Group users retrieved successfully",
		Result: &pb.GetUsersInGroupResponse_Users{
			Users: &pb.UserInfoList{
				Users: users,
			},
		},
	}, nil
}

func (s *ChatServer) GetUserStatus(ctx context.Context, req *pb.UserStatusRequest) (*pb.UserStatusResponse, error) {
	log.Printf("GetUserStatus called for %d users", len(req.UserIds))

	_, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.UserStatusResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.UserStatusResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if len(req.UserIds) == 0 {
		return &pb.UserStatusResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.UserStatusResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "User IDs list cannot be empty",
					Details:   []string{"At least one user ID is required"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Mock status data - in production, check presence service
	var statuses []*pb.UserStatus
	for _, userID := range req.UserIds {
		status := &pb.UserStatus{
			UserId:   userID,
			IsOnline: false,                                       // Mock offline status
			LastSeen: timestamppb.New(time.Now().Add(-time.Hour)), // Mock last seen
		}
		statuses = append(statuses, status)
	}

	return &pb.UserStatusResponse{
		StatusCode: 200,
		Message:    "User statuses retrieved successfully",
		Result: &pb.UserStatusResponse_Statuses{
			Statuses: &pb.UserStatusList{
				Statuses: statuses,
			},
		},
	}, nil
}

func (s *ChatServer) CreateGroup(ctx context.Context, req *pb.CreateGroupRequest) (*pb.CreateGroupResponse, error) {
	log.Printf("CreateGroup called by user: %s", req.CreatorId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.CreateGroupResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.CreateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.CreatorId != claims.UserID {
		return &pb.CreateGroupResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.CreateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Creator ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.GroupName == "" {
		return &pb.CreateGroupResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.CreateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Group name is required",
					Details:   []string{"Group name cannot be empty"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	groupID := uuid.New().String()
	now := time.Now()

	// Start transaction
	tx, err := s.db.DB.Begin()
	if err != nil {
		return &pb.CreateGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.CreateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to start transaction",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	defer tx.Rollback()

	// Create group
	_, err = tx.Exec(`
		INSERT INTO chat_groups (id, name, description, avatar_url, creator_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		groupID, req.GroupName, req.Description, req.AvatarUrl, req.CreatorId, now, now)

	if err != nil {
		return &pb.CreateGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.CreateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to create group",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Add creator as owner
	_, err = tx.Exec(`
		INSERT INTO group_members (group_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4)`,
		groupID, req.CreatorId, int32(pb.GroupRole_OWNER), now)

	if err != nil {
		return &pb.CreateGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.CreateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to add creator to group",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Add other members
	for _, memberID := range req.MemberIds {
		if memberID != req.CreatorId { // Don't add creator twice
			_, err = tx.Exec(`
				INSERT INTO group_members (group_id, user_id, role, joined_at)
				VALUES ($1, $2, $3, $4)`,
				groupID, memberID, int32(pb.GroupRole_MEMBER), now)

			if err != nil {
				log.Printf("Failed to add member %s to group: %v", memberID, err)
				// Continue with other members
			}
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return &pb.CreateGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.CreateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to commit transaction",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Create response group info
	groupInfo := &pb.GroupInfo{
		Id:          groupID,
		Name:        req.GroupName,
		Description: req.Description,
		AvatarUrl:   req.AvatarUrl,
		CreatorId:   req.CreatorId,
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
		MaxMembers:  100, // Default max members
		IsPrivate:   false,
	}

	return &pb.CreateGroupResponse{
		StatusCode: 200,
		Message:    "Group created successfully",
		Result: &pb.CreateGroupResponse_Group{
			Group: groupInfo,
		},
	}, nil
}

func (s *ChatServer) JoinGroup(ctx context.Context, req *pb.JoinGroupRequest) (*pb.JoinGroupResponse, error) {
	log.Printf("JoinGroup called by user: %s for group: %s", req.UserId, req.GroupId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.JoinGroupResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.JoinGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.JoinGroupResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.JoinGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if group exists
	var groupExists bool
	err = s.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM chat_groups WHERE id = $1)", req.GroupId).Scan(&groupExists)
	if err != nil || !groupExists {
		return &pb.JoinGroupResponse{
			StatusCode: 404,
			Message:    "Group not found",
			Result: &pb.JoinGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      404,
					Message:   "Group not found",
					Details:   []string{"Group does not exist"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if user is already a member
	var isMember bool
	err = s.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)", req.GroupId, req.UserId).Scan(&isMember)
	if err == nil && isMember {
		return &pb.JoinGroupResponse{
			StatusCode: 409,
			Message:    "Already a member",
			Result: &pb.JoinGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      409,
					Message:   "User is already a member of this group",
					Details:   []string{"Conflict"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Add user to group
	_, err = s.db.DB.Exec(`
		INSERT INTO group_members (group_id, user_id, role, joined_at)
		VALUES ($1, $2, $3, $4)`,
		req.GroupId, req.UserId, int32(pb.GroupRole_MEMBER), time.Now())

	if err != nil {
		return &pb.JoinGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.JoinGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to join group",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.JoinGroupResponse{
		StatusCode: 200,
		Message:    "Joined group successfully",
		Result: &pb.JoinGroupResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Joined group",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (s *ChatServer) LeaveGroup(ctx context.Context, req *pb.LeaveGroupRequest) (*pb.LeaveGroupResponse, error) {
	log.Printf("LeaveGroup called by user: %s for group: %s", req.UserId, req.GroupId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.LeaveGroupResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.LeaveGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.LeaveGroupResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.LeaveGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if user is a member and get their role
	var role int32
	err = s.db.DB.QueryRow("SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2", req.GroupId, req.UserId).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.LeaveGroupResponse{
				StatusCode: 404,
				Message:    "Not a member",
				Result: &pb.LeaveGroupResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      404,
						Message:   "User is not a member of this group",
						Details:   []string{"Not found"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
		return &pb.LeaveGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.LeaveGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Database error",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// If user is owner, check if there are other members to transfer ownership
	if pb.GroupRole(role) == pb.GroupRole_OWNER {
		var memberCount int
		err = s.db.DB.QueryRow("SELECT COUNT(*) FROM group_members WHERE group_id = $1", req.GroupId).Scan(&memberCount)
		if err == nil && memberCount > 1 {
			return &pb.LeaveGroupResponse{
				StatusCode: 400,
				Message:    "Cannot leave group",
				Result: &pb.LeaveGroupResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      400,
						Message:   "Owner cannot leave group with members. Transfer ownership first.",
						Details:   []string{"Bad request"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
	}

	// Remove user from group
	_, err = s.db.DB.Exec("DELETE FROM group_members WHERE group_id = $1 AND user_id = $2", req.GroupId, req.UserId)
	if err != nil {
		return &pb.LeaveGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.LeaveGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to leave group",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.LeaveGroupResponse{
		StatusCode: 200,
		Message:    "Left group successfully",
		Result: &pb.LeaveGroupResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Left group",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (s *ChatServer) UpdateGroup(ctx context.Context, req *pb.UpdateGroupRequest) (*pb.UpdateGroupResponse, error) {
	log.Printf("UpdateGroup called by user: %s for group: %s", req.UserId, req.GroupId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.UpdateGroupResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.UpdateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.UpdateGroupResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.UpdateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if user has admin/owner permissions
	var role int32
	err = s.db.DB.QueryRow("SELECT role FROM group_members WHERE group_id = $1 AND user_id = $2", req.GroupId, req.UserId).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UpdateGroupResponse{
				StatusCode: 403,
				Message:    "Forbidden",
				Result: &pb.UpdateGroupResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      403,
						Message:   "User is not a member of this group",
						Details:   []string{"Access denied"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
		return &pb.UpdateGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.UpdateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Database error",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if pb.GroupRole(role) != pb.GroupRole_ADMIN && pb.GroupRole(role) != pb.GroupRole_OWNER {
		return &pb.UpdateGroupResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.UpdateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Only admins and owners can update group",
					Details:   []string{"Insufficient permissions"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Start transaction
	tx, err := s.db.DB.Begin()
	if err != nil {
		return &pb.UpdateGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.UpdateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to start transaction",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	defer tx.Rollback()

	// Build update query dynamically based on provided fields
	updateFields := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.GroupName != "" {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, req.GroupName)
		argIndex++
	}

	if req.Description != "" {
		updateFields = append(updateFields, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, req.Description)
		argIndex++
	}

	if req.AvatarUrl != "" {
		updateFields = append(updateFields, fmt.Sprintf("avatar_url = $%d", argIndex))
		args = append(args, req.AvatarUrl)
		argIndex++
	}

	if len(updateFields) == 0 {
		return &pb.UpdateGroupResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.UpdateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "No fields to update",
					Details:   []string{"At least one field must be provided"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Add updated_at field
	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, req.GroupId)
	query := fmt.Sprintf("UPDATE chat_groups SET %s WHERE id = $%d",
		string(updateFields[0]), argIndex)
	for i := 1; i < len(updateFields); i++ {
		query = fmt.Sprintf("%s, %s", query, updateFields[i])
	}

	// Execute update
	result, err := tx.Exec(query, args...)
	if err != nil {
		return &pb.UpdateGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.UpdateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to update group",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &pb.UpdateGroupResponse{
			StatusCode: 404,
			Message:    "Group not found",
			Result: &pb.UpdateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      404,
					Message:   "Group not found",
					Details:   []string{"Group does not exist"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return &pb.UpdateGroupResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.UpdateGroupResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to commit transaction",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.UpdateGroupResponse{
		StatusCode: 200,
		Message:    "Group updated successfully",
		Result: &pb.UpdateGroupResponse_Group{
			Group: &pb.GroupInfo{
				Id:          req.GroupId,
				Name:        req.GroupName,
				Description: req.Description,
				AvatarUrl:   req.AvatarUrl,
			},
		},
	}, nil
}

func (s *ChatServer) GetUsersByUserEmail(ctx context.Context, req *pb.GetUsersByUserEmailRequest) (*pb.GetUsersByUserEmailResponse, error) {
	log.Printf("GetUsersByUserEmail called for emails: %v", req.UserEmail)

	_, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetUsersByUserEmailResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetUsersByUserEmailResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if len(req.UserEmail) == 0 {
		return &pb.GetUsersByUserEmailResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.GetUsersByUserEmailResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "At least one email is required",
					Details:   []string{"Email list cannot be empty"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	var users []*pb.UserInfo
	for _, email := range req.UserEmail {
		user, err := s.userClient.GetUserByEmail(ctx, email)
		if err != nil {
			log.Printf("Failed to get user by email %s: %v", email, err)
			continue
		}

		userInfo := &pb.UserInfo{
			Id:       user.Id,
			FullName: user.FullName,
			Email:    user.Email,
		}

		users = append(users, userInfo)
	}

	return &pb.GetUsersByUserEmailResponse{
		StatusCode: 200,
		Message:    "Users retrieved successfully",
		Result: &pb.GetUsersByUserEmailResponse_Users{
			Users: &pb.UserInfoList{
				Users: users,
			},
		},
	}, nil
}

func (s *ChatServer) GetAllGroupsByUserEmail(ctx context.Context, req *pb.GetAllGroupsByUserEmailRequest) (*pb.GetAllGroupsByUserEmailResponse, error) {
	log.Printf("GetAllGroupsByUserEmail called for email: %s", req.UserEmail)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetAllGroupsByUserEmailResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetAllGroupsByUserEmailResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	if req.UserEmail == "" {
		return &pb.GetAllGroupsByUserEmailResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.GetAllGroupsByUserEmailResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "User email is required",
					Details:   []string{"Email cannot be empty"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	user, err := s.userClient.GetUserByEmail(ctx, req.UserEmail)
	if err != nil {
		return &pb.GetAllGroupsByUserEmailResponse{
			StatusCode: 404,
			Message:    "User not found",
			Result: &pb.GetAllGroupsByUserEmailResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      404,
					Message:   "User not found",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if user.Email != claims.Email && claims.Role != "admin" {
		return &pb.GetAllGroupsByUserEmailResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.GetAllGroupsByUserEmailResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Access denied",
					Details:   []string{"You can only access your own groups"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	query := `
		SELECT g.id, g.name, g.description, g.avatar_url, g.creator_id, g.created_at, g.updated_at
		FROM chat_groups g
		INNER JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = $1
		ORDER BY g.created_at DESC`

	rows, err := s.db.DB.Query(query, user.Id)
	if err != nil {
		return &pb.GetAllGroupsByUserEmailResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.GetAllGroupsByUserEmailResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to fetch groups",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	defer rows.Close()

	var groups []*pb.GroupInfo
	for rows.Next() {
		var group pb.GroupInfo
		var createdAt, updatedAt time.Time
		var description, avatarURL sql.NullString

		err := rows.Scan(
			&group.Id,
			&group.Name,
			&description,
			&avatarURL,
			&group.CreatorId,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			log.Printf("Error scanning group row: %v", err)
			continue
		}

		group.Description = description.String
		group.AvatarUrl = avatarURL.String
		group.CreatedAt = timestamppb.New(createdAt)
		group.UpdatedAt = timestamppb.New(updatedAt)

		membersQuery := `
			SELECT user_id, role, joined_at
			FROM group_members
			WHERE group_id = $1`

		memberRows, err := s.db.DB.Query(membersQuery, group.Id)
		if err != nil {
			log.Printf("Error fetching group members: %v", err)
			continue
		}

		var members []*pb.GroupMember
		for memberRows.Next() {
			var member pb.GroupMember
			var joinedAt time.Time
			var role int32

			err := memberRows.Scan(&member.UserId, &role, &joinedAt)
			if err != nil {
				log.Printf("Error scanning member row: %v", err)
				continue
			}

			member.Role = pb.GroupRole(role)
			member.JoinedAt = timestamppb.New(joinedAt)
			members = append(members, &member)
		}
		memberRows.Close()

		group.Members = members
		groups = append(groups, &group)
	}

	// Note: The proto definition shows GroupInfo groups = 3, but it should probably be repeated GroupInfo
	// For now, returning the first group if any exist, or you may need to update the proto
	if len(groups) > 0 {
		return &pb.GetAllGroupsByUserEmailResponse{
			StatusCode: 200,
			Message:    "Groups retrieved successfully",
			Result: &pb.GetAllGroupsByUserEmailResponse_Groups{
				Groups: &pb.GroupInfoList{
					Groups: groups,
				},
			},
		}, nil
	} else {
		return &pb.GetAllGroupsByUserEmailResponse{
			StatusCode: 200,
			Message:    "No groups found",
			Result: &pb.GetAllGroupsByUserEmailResponse_Groups{
				Groups: &pb.GroupInfoList{
					Groups: groups,
				},
			},
		}, nil
	}
}

func (s *ChatServer) InitiateCall(ctx context.Context, req *pb.InitiateCallRequest) (*pb.InitiateCallResponse, error) {
	log.Printf("InitiateCall called by caller: %s", req.CallerId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.InitiateCallResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.InitiateCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate caller ID matches token
	if req.CallerId != claims.UserID {
		return &pb.InitiateCallResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.InitiateCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Caller ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if !ValidateUserID(req.CallerId) {
		return &pb.InitiateCallResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.InitiateCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid caller ID format",
					Details:   []string{"Caller ID must be a valid UUID"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.IsGroup {
		if !ValidateUserID(req.GroupId) {
			return &pb.InitiateCallResponse{
				StatusCode: 400,
				Message:    "Bad Request",
				Result: &pb.InitiateCallResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      400,
						Message:   "Invalid group ID format",
						Details:   []string{"Group ID must be a valid UUID"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
	} else {
		if !ValidateUserID(req.ReceiverId) {
			return &pb.InitiateCallResponse{
				StatusCode: 400,
				Message:    "Bad Request",
				Result: &pb.InitiateCallResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      400,
						Message:   "Invalid receiver ID format",
						Details:   []string{"Receiver ID must be a valid UUID"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
	}

	// Generate call ID
	callID := uuid.New().String()
	now := time.Now()

	// Convert call type to string
	callTypeStr := "CALL_VOICE"
	if req.CallType == pb.CallType_CALL_VIDEO {
		callTypeStr = "CALL_VIDEO"
	}

	// Insert call into database
	query := `
		INSERT INTO chat_calls (
			id, caller_id, receiver_id, group_id, call_type, 
			status, is_group, start_time, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, start_time`

	var savedID string
	var startTime time.Time
	err = s.db.DB.QueryRow(
		query,
		callID,
		req.CallerId,
		nullString(req.ReceiverId),
		nullString(req.GroupId),
		callTypeStr,
		"CALL_INITIATED",
		req.IsGroup,
		now,
		now,
		now,
	).Scan(&savedID, &startTime)

	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.InitiateCallResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.InitiateCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to initiate call",
					Details:   []string{"Database operation failed"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Get participants for group calls
	var participants []string
	if req.IsGroup {
		rows, err := s.db.DB.Query("SELECT user_id FROM group_members WHERE group_id = $1", req.GroupId)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var userID string
				if err := rows.Scan(&userID); err == nil {
					participants = append(participants, userID)
				}
			}
		}
	} else {
		participants = []string{req.CallerId, req.ReceiverId}
	}

	// Create response call info
	callInfo := &pb.CallInfo{
		CallId:       savedID,
		CallerId:     req.CallerId,
		ReceiverId:   req.ReceiverId,
		GroupId:      req.GroupId,
		CallType:     req.CallType,
		Status:       pb.CallStatus_CALL_INITIATED,
		StartTime:    timestamppb.New(startTime),
		IsGroup:      req.IsGroup,
		Participants: participants,
	}

	return &pb.InitiateCallResponse{
		StatusCode: 200,
		Message:    "Call initiated successfully",
		Result: &pb.InitiateCallResponse_Call{
			Call: callInfo,
		},
	}, nil
}

func (s *ChatServer) AcceptCall(ctx context.Context, req *pb.AcceptCallRequest) (*pb.AcceptCallResponse, error) {
	log.Printf("AcceptCall called for call: %s by user: %s", req.CallId, req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.AcceptCallResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.AcceptCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.AcceptCallResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.AcceptCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if !ValidateUserID(req.CallId) || !ValidateUserID(req.UserId) {
		return &pb.AcceptCallResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.AcceptCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid call ID or user ID format",
					Details:   []string{"IDs must be valid UUIDs"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if call exists and user is authorized to accept it
	var callerID, receiverID, groupID, callType, status sql.NullString
	var isGroup bool
	var startTime time.Time
	err = s.db.DB.QueryRow(`
		SELECT caller_id, receiver_id, group_id, call_type, status, is_group, start_time 
		FROM chat_calls WHERE id = $1`, req.CallId).Scan(
		&callerID, &receiverID, &groupID, &callType, &status, &isGroup, &startTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.AcceptCallResponse{
				StatusCode: 404,
				Message:    "Call not found",
				Result: &pb.AcceptCallResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      404,
						Message:   "Call not found",
						Details:   []string{"Call does not exist"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
		return &pb.AcceptCallResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.AcceptCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Database error",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if user is authorized to accept the call
	canAccept := false
	if isGroup {
		// For group calls, check if user is a member of the group
		var count int
		err = s.db.DB.QueryRow("SELECT COUNT(*) FROM group_members WHERE group_id = $1 AND user_id = $2",
			groupID.String, req.UserId).Scan(&count)
		if err == nil && count > 0 {
			canAccept = true
		}
	} else {
		// For direct calls, check if user is the receiver
		if receiverID.String == req.UserId {
			canAccept = true
		}
	}

	if !canAccept {
		return &pb.AcceptCallResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.AcceptCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Not authorized to accept this call",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if call is in a state that can be accepted
	if status.String != "CALL_INITIATED" && status.String != "CALL_RINGING" {
		return &pb.AcceptCallResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.AcceptCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Call cannot be accepted in current state",
					Details:   []string{fmt.Sprintf("Current status: %s", status.String)},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Update call status to accepted
	now := time.Now()
	_, err = s.db.DB.Exec(
		"UPDATE chat_calls SET status = 'CALL_ACCEPTED', updated_at = $1 WHERE id = $2",
		now, req.CallId,
	)
	if err != nil {
		return &pb.AcceptCallResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.AcceptCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to accept call",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Get participants
	var participants []string
	if isGroup {
		rows, err := s.db.DB.Query("SELECT user_id FROM group_members WHERE group_id = $1", groupID.String)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var userID string
				if err := rows.Scan(&userID); err == nil {
					participants = append(participants, userID)
				}
			}
		}
	} else {
		participants = []string{callerID.String, receiverID.String}
	}

	// Convert call type
	callTypeEnum := pb.CallType_CALL_VOICE
	if callType.String == "CALL_VIDEO" {
		callTypeEnum = pb.CallType_CALL_VIDEO
	}

	// Create response call info
	callInfo := &pb.CallInfo{
		CallId:       req.CallId,
		CallerId:     callerID.String,
		ReceiverId:   receiverID.String,
		GroupId:      groupID.String,
		CallType:     callTypeEnum,
		Status:       pb.CallStatus_CALL_ACCEPTED,
		StartTime:    timestamppb.New(startTime),
		IsGroup:      isGroup,
		Participants: participants,
	}

	return &pb.AcceptCallResponse{
		StatusCode: 200,
		Message:    "Call accepted successfully",
		Result: &pb.AcceptCallResponse_Call{
			Call: callInfo,
		},
	}, nil
}

func (s *ChatServer) RejectCall(ctx context.Context, req *pb.RejectCallRequest) (*pb.RejectCallResponse, error) {
	log.Printf("RejectCall called for call: %s by user: %s", req.CallId, req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.RejectCallResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.RejectCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.RejectCallResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.RejectCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if !ValidateUserID(req.CallId) || !ValidateUserID(req.UserId) {
		return &pb.RejectCallResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.RejectCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid call ID or user ID format",
					Details:   []string{"IDs must be valid UUIDs"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if call exists and user is authorized to reject it
	var callerID, receiverID, groupID, status sql.NullString
	var isGroup bool
	err = s.db.DB.QueryRow(`
		SELECT caller_id, receiver_id, group_id, status, is_group 
		FROM chat_calls WHERE id = $1`, req.CallId).Scan(
		&callerID, &receiverID, &groupID, &status, &isGroup)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.RejectCallResponse{
				StatusCode: 404,
				Message:    "Call not found",
				Result: &pb.RejectCallResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      404,
						Message:   "Call not found",
						Details:   []string{"Call does not exist"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
		return &pb.RejectCallResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.RejectCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Database error",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if user is authorized to reject the call
	canReject := false
	if isGroup {
		// For group calls, check if user is a member of the group
		var count int
		err = s.db.DB.QueryRow("SELECT COUNT(*) FROM group_members WHERE group_id = $1 AND user_id = $2",
			groupID.String, req.UserId).Scan(&count)
		if err == nil && count > 0 {
			canReject = true
		}
	} else {
		// For direct calls, check if user is the receiver or caller
		if receiverID.String == req.UserId || callerID.String == req.UserId {
			canReject = true
		}
	}

	if !canReject {
		return &pb.RejectCallResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.RejectCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Not authorized to reject this call",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if call is in a state that can be rejected
	if status.String != "CALL_INITIATED" && status.String != "CALL_RINGING" {
		return &pb.RejectCallResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.RejectCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Call cannot be rejected in current state",
					Details:   []string{fmt.Sprintf("Current status: %s", status.String)},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Update call status to rejected
	now := time.Now()
	_, err = s.db.DB.Exec(
		"UPDATE chat_calls SET status = 'CALL_REJECTED', end_time = $1, updated_at = $2 WHERE id = $3",
		now, now, req.CallId,
	)
	if err != nil {
		return &pb.RejectCallResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.RejectCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to reject call",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.RejectCallResponse{
		StatusCode: 200,
		Message:    "Call rejected successfully",
		Result: &pb.RejectCallResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Call rejected",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (s *ChatServer) EndCall(ctx context.Context, req *pb.EndCallRequest) (*pb.EndCallResponse, error) {
	log.Printf("EndCall called for call: %s by user: %s", req.CallId, req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.EndCallResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.EndCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.EndCallResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.EndCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if !ValidateUserID(req.CallId) || !ValidateUserID(req.UserId) {
		return &pb.EndCallResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.EndCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid call ID or user ID format",
					Details:   []string{"IDs must be valid UUIDs"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if call exists and get call details
	var callerID, receiverID, groupID, callType, status sql.NullString
	var isGroup bool
	var startTime time.Time
	err = s.db.DB.QueryRow(`
		SELECT caller_id, receiver_id, group_id, call_type, status, is_group, start_time 
		FROM chat_calls WHERE id = $1`, req.CallId).Scan(
		&callerID, &receiverID, &groupID, &callType, &status, &isGroup, &startTime)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.EndCallResponse{
				StatusCode: 404,
				Message:    "Call not found",
				Result: &pb.EndCallResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      404,
						Message:   "Call not found",
						Details:   []string{"Call does not exist"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
		return &pb.EndCallResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.EndCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Database error",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if user is authorized to end the call
	canEnd := false
	if isGroup {
		// For group calls, check if user is a member of the group
		var count int
		err = s.db.DB.QueryRow("SELECT COUNT(*) FROM group_members WHERE group_id = $1 AND user_id = $2",
			groupID.String, req.UserId).Scan(&count)
		if err == nil && count > 0 {
			canEnd = true
		}
	} else {
		// For direct calls, check if user is the receiver or caller
		if receiverID.String == req.UserId || callerID.String == req.UserId {
			canEnd = true
		}
	}

	if !canEnd {
		return &pb.EndCallResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.EndCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Not authorized to end this call",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if call is already ended
	if status.String == "CALL_ENDED" || status.String == "CALL_REJECTED" {
		return &pb.EndCallResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.EndCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Call is already ended",
					Details:   []string{fmt.Sprintf("Current status: %s", status.String)},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Calculate duration
	now := time.Now()
	duration := int32(now.Sub(startTime).Seconds())

	// Update call status to ended
	_, err = s.db.DB.Exec(`
		UPDATE chat_calls 
		SET status = 'CALL_ENDED', end_time = $1, duration = $2, updated_at = $3 
		WHERE id = $4`,
		now, duration, now, req.CallId,
	)
	if err != nil {
		return &pb.EndCallResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.EndCallResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to end call",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Get participants
	var participants []string
	if isGroup {
		rows, err := s.db.DB.Query("SELECT user_id FROM group_members WHERE group_id = $1", groupID.String)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var userID string
				if err := rows.Scan(&userID); err == nil {
					participants = append(participants, userID)
				}
			}
		}
	} else {
		participants = []string{callerID.String, receiverID.String}
	}

	// Convert call type
	callTypeEnum := pb.CallType_CALL_VOICE
	if callType.String == "CALL_VIDEO" {
		callTypeEnum = pb.CallType_CALL_VIDEO
	}

	// Create response call info
	callInfo := &pb.CallInfo{
		CallId:       req.CallId,
		CallerId:     callerID.String,
		ReceiverId:   receiverID.String,
		GroupId:      groupID.String,
		CallType:     callTypeEnum,
		Status:       pb.CallStatus_CALL_ENDED,
		StartTime:    timestamppb.New(startTime),
		EndTime:      timestamppb.New(now),
		Duration:     duration,
		IsGroup:      isGroup,
		Participants: participants,
	}

	return &pb.EndCallResponse{
		StatusCode: 200,
		Message:    "Call ended successfully",
		Result: &pb.EndCallResponse_Call{
			Call: callInfo,
		},
	}, nil
}

func (s *ChatServer) GetCallHistory(ctx context.Context, req *pb.GetCallHistoryRequest) (*pb.GetCallHistoryResponse, error) {
	log.Printf("GetCallHistory called by user: %s", req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetCallHistoryResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetCallHistoryResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.GetCallHistoryResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.GetCallHistoryResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if !ValidateUserID(req.UserId) {
		return &pb.GetCallHistoryResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.GetCallHistoryResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid user ID format",
					Details:   []string{"User ID must be a valid UUID"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Set default values for limit and offset
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Query call history for the user
	query := `
		SELECT id, caller_id, COALESCE(receiver_id, ''), COALESCE(group_id, ''), 
		       call_type, status, is_group, start_time, end_time, duration
		FROM chat_calls 
		WHERE (caller_id = $1 OR receiver_id = $1 OR 
		       (is_group = true AND group_id IN (
		           SELECT group_id FROM group_members WHERE user_id = $1
		       )))
		ORDER BY start_time DESC 
		LIMIT $2 OFFSET $3`

	rows, err := s.db.DB.Query(query, req.UserId, limit, offset)
	if err != nil {
		return &pb.GetCallHistoryResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.GetCallHistoryResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to fetch call history",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	defer rows.Close()

	var calls []*pb.CallInfo
	for rows.Next() {
		var call pb.CallInfo
		var receiverID, groupID, callType, status sql.NullString
		var startTime time.Time
		var endTime sql.NullTime
		var duration sql.NullInt32

		err := rows.Scan(
			&call.CallId, &call.CallerId, &receiverID, &groupID,
			&callType, &status, &call.IsGroup, &startTime, &endTime, &duration,
		)
		if err != nil {
			continue
		}

		call.ReceiverId = receiverID.String
		call.GroupId = groupID.String
		call.StartTime = timestamppb.New(startTime)

		if endTime.Valid {
			call.EndTime = timestamppb.New(endTime.Time)
		}
		if duration.Valid {
			call.Duration = duration.Int32
		}

		// Convert call type
		if callType.String == "CALL_VIDEO" {
			call.CallType = pb.CallType_CALL_VIDEO
		} else {
			call.CallType = pb.CallType_CALL_VOICE
		}

		// Convert status
		switch status.String {
		case "CALL_INITIATED":
			call.Status = pb.CallStatus_CALL_INITIATED
		case "CALL_RINGING":
			call.Status = pb.CallStatus_CALL_RINGING
		case "CALL_ACCEPTED":
			call.Status = pb.CallStatus_CALL_ACCEPTED
		case "CALL_REJECTED":
			call.Status = pb.CallStatus_CALL_REJECTED
		case "CALL_ENDED":
			call.Status = pb.CallStatus_CALL_ENDED
		case "CALL_MISSED":
			call.Status = pb.CallStatus_CALL_MISSED
		case "CALL_BUSY":
			call.Status = pb.CallStatus_CALL_BUSY
		default:
			call.Status = pb.CallStatus_CALL_INITIATED
		}

		// Get participants
		var participants []string
		if call.IsGroup {
			pRows, err := s.db.DB.Query("SELECT user_id FROM group_members WHERE group_id = $1", call.GroupId)
			if err == nil {
				for pRows.Next() {
					var userID string
					if err := pRows.Scan(&userID); err == nil {
						participants = append(participants, userID)
					}
				}
				pRows.Close()
			}
		} else {
			participants = []string{call.CallerId, call.ReceiverId}
		}
		call.Participants = participants

		calls = append(calls, &call)
	}

	return &pb.GetCallHistoryResponse{
		StatusCode: 200,
		Message:    "Call history retrieved successfully",
		Result: &pb.GetCallHistoryResponse_Calls{
			Calls: &pb.CallInfoList{
				Calls: calls,
			},
		},
	}, nil
}

func (s *ChatServer) AddNotification(ctx context.Context, req *pb.AddNotificationRequest) (*pb.AddNotificationResponse, error) {
	log.Printf("AddNotification called for user: %s", req.Notification.SenderId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.AddNotificationResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.AddNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if req.Notification == nil {
		return &pb.AddNotificationResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.AddNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Notification data is required",
					Details:   []string{"notification field cannot be empty"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	notification := req.Notification
	if notification.Title == "" || notification.SenderId == "" {
		return &pb.AddNotificationResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.AddNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Title and sender ID are required",
					Details:   []string{"title and senderId fields are mandatory"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Generate notification ID
	notificationID := uuid.New().String()

	// Insert notification into database
	query := `
		INSERT INTO chat_notifications (id, user_id, type, title, content, sender_id, chat_id, group_id, message_id, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = s.db.Exec(query,
		notificationID,                     // $1 - id
		claims.UserID,                      // $2 - user_id from token
		notification.Type.String(),         // $3 - type
		notification.Title,                 // $4 - title
		nullString(notification.Content),   // $5 - content
		(notification.SenderId),            // $6 - sender_id (can be null)
		nullString(notification.ChatId),    // $7 - chat_id
		nullString(""),                     // $8 - group_id (separate from chat_id)
		nullString(notification.MessageId), // $9 - message_id
		notification.IsRead,                // $10 - is_read
		time.Now(),                         // $11 - created_at
	)

	if err != nil {
		log.Printf("Error inserting notification: %v", err)
		return &pb.AddNotificationResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.AddNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to create notification",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	// NEW: Send push notification via OneSignal
	err = s.sendPushNotification(claims.UserID, notification.Title, notification.Content)
	if err != nil {
		log.Printf("Failed to send push notification: %v", err)
		// Don't fail the entire request if push notification fails
	}
	// Return the created notification
	notification.Id = notificationID
	notification.Timestamp = timestamppb.New(time.Now())

	return &pb.AddNotificationResponse{
		StatusCode: 200,
		Message:    "Notification created successfully",
		Result: &pb.AddNotificationResponse_Notification{
			Notification: notification,
		},
	}, nil
}
func (s *ChatServer) sendPushNotification(userID, title, content string) error {
	config := onesignal.NewConfiguration()
	client := onesignal.NewAPIClient(config)
	if client == nil {
		return fmt.Errorf("failed to initialize OneSignal client")
	}
	authCtx := context.WithValue(context.Background(), onesignal.ContextAPIKeys, map[string]onesignal.APIKey{
		"app_key": {Key: os.Getenv("ONESIGNAL_REST_API_KEY")},
	})

	notification := onesignal.NewNotification(os.Getenv("ONESIGNAL_APP_ID"))
	notification.SetAppId(os.Getenv("ONESIGNAL_APP_ID"))

	headings := onesignal.StringMap{}
	headings.Set("en", title)
	notification.SetHeadings(headings)

	contents := onesignal.StringMap{}
	contents.Set("en", content)
	notification.SetContents(contents)

	notification.SetIncludeExternalUserIds([]string{userID})

	resp, _, err := client.DefaultApi.CreateNotification(authCtx).Notification(*notification).Execute()
	if err != nil {
		return fmt.Errorf("OneSignal API error: %v", err)
	}

	log.Printf("OneSignal notification sent, ID: %v", resp.GetId())
	return nil
}

func (s *ChatServer) UpdateNotification(ctx context.Context, req *pb.UpdateNotificationRequest) (*pb.UpdateNotificationResponse, error) {
	log.Printf("UpdateNotification called for notification: %s", req.Notification.Id)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.UpdateNotificationResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.UpdateNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if req.Notification == nil || req.Notification.Id == "" {
		return &pb.UpdateNotificationResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.UpdateNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Notification ID is required",
					Details:   []string{"notification.id field cannot be empty"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	notification := req.Notification

	// Check if notification exists and belongs to user
	var existingUserID string
	checkQuery := `SELECT user_id FROM chat_notifications WHERE id = $1`
	err = s.db.DB.QueryRow(checkQuery, notification.Id).Scan(&existingUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UpdateNotificationResponse{
				StatusCode: 404,
				Message:    "Not Found",
				Result: &pb.UpdateNotificationResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      404,
						Message:   "Notification not found",
						Details:   []string{"notification with given ID does not exist"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
		return &pb.UpdateNotificationResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.UpdateNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to check notification",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check authorization
	if existingUserID != claims.UserID {
		return &pb.UpdateNotificationResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.UpdateNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Access denied",
					Details:   []string{"you can only update your own notifications"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Update notification - Fixed parameter count
	updateQuery := `
		UPDATE chat_notifications 
		SET title = $1, content = $2, is_read = $3
		WHERE id = $4
	`

	_, err = s.db.Exec(updateQuery,
		notification.Id,                    // $1 - id
		claims.UserID,                      // $2 - user_id from token
		notification.Type.String(),         // $3 - type
		notification.Title,                 // $4 - title
		nullString(notification.Content),   // $5 - content
		notification.SenderId,              // $6 - sender_id (can be null)
		nullString(notification.ChatId),    // $7 - chat_id
		nullString(""),                     // $8 - group_id (separate from chat_id)
		nullString(notification.MessageId), // $9 - message_id
		notification.IsRead,                // $10 - is_read
		time.Now(),
	)
	if err != nil {
		log.Printf("Error updating notification: %v", err)
		return &pb.UpdateNotificationResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.UpdateNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to update notification",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.UpdateNotificationResponse{
		StatusCode: 200,
		Message:    "Notification updated successfully",
		Result: &pb.UpdateNotificationResponse_Notification{
			Notification: notification,
		},
	}, nil
}

func (s *ChatServer) GetNotification(ctx context.Context, req *pb.GetNotificationRequest) (*pb.GetNotificationResponse, error) {
	log.Printf("GetNotification called for notification: %s", req.NotificationId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetNotificationResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if req.NotificationId == "" {
		return &pb.GetNotificationResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.GetNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Notification ID is required",
					Details:   []string{"notificationId field cannot be empty"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check authorization if userId is provided
	if req.UserId != "" && req.UserId != claims.UserID {
		return &pb.GetNotificationResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.GetNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Access denied",
					Details:   []string{"you can only access your own notifications"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Query notification
	query := `
		SELECT id, user_id, type, title, content, sender_id, chat_id, group_id, message_id, is_read, created_at
		FROM chat_notifications 
		WHERE id = $1 AND user_id = $2
	`

	var notification pb.NotificationUpdate
	var notificationType string
	var content, senderId, chatId, groupId, messageId sql.NullString
	var createdAt time.Time

	err = s.db.DB.QueryRow(query, req.NotificationId, claims.UserID).Scan(
		&notification.Id,
		&notification.SenderId, // This will be overwritten
		&notificationType,
		&notification.Title,
		&content,
		&senderId,
		&chatId,
		&groupId,
		&messageId,
		&notification.IsRead,
		&createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetNotificationResponse{
				StatusCode: 404,
				Message:    "Not Found",
				Result: &pb.GetNotificationResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      404,
						Message:   "Notification not found",
						Details:   []string{"notification with given ID does not exist or access denied"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
		return &pb.GetNotificationResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.GetNotificationResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to retrieve notification",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Map notification type
	switch notificationType {
	case "NOTIFICATION_MESSAGE":
		notification.Type = pb.NotificationType_NOTIFICATION_MESSAGE
	case "NOTIFICATION_MENTION":
		notification.Type = pb.NotificationType_NOTIFICATION_MENTION
	case "NOTIFICATION_REACTION":
		notification.Type = pb.NotificationType_NOTIFICATION_REACTION
	case "NOTIFICATION_CALL":
		notification.Type = pb.NotificationType_NOTIFICATION_CALL
	case "NOTIFICATION_GROUP_INVITE":
		notification.Type = pb.NotificationType_NOTIFICATION_GROUP_INVITE
	case "NOTIFICATION_SYSTEM":
		notification.Type = pb.NotificationType_NOTIFICATION_SYSTEM
	default:
		notification.Type = pb.NotificationType_NOTIFICATION_MESSAGE
	}

	// Set optional fields
	if content.Valid {
		notification.Content = content.String
	}
	if senderId.Valid {
		notification.SenderId = senderId.String
	}
	if chatId.Valid {
		notification.ChatId = chatId.String
	}
	if messageId.Valid {
		notification.MessageId = messageId.String
	}
	notification.Timestamp = timestamppb.New(createdAt)

	return &pb.GetNotificationResponse{
		StatusCode: 200,
		Message:    "Notification retrieved successfully",
		Result: &pb.GetNotificationResponse_Notification{
			Notification: &notification,
		},
	}, nil
}

func (s *ChatServer) MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadRequest) (*pb.MarkNotificationAsReadResponse, error) {
	log.Printf("MarkNotificationAsRead called for notification: %s", req.NotificationId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.MarkNotificationAsReadResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.MarkNotificationAsReadResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if req.NotificationId == "" {
		return &pb.MarkNotificationAsReadResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.MarkNotificationAsReadResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Notification ID is required",
					Details:   []string{"notificationId field cannot be empty"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check authorization if userId is provided
	if req.UserId != "" && req.UserId != claims.UserID {
		return &pb.MarkNotificationAsReadResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.MarkNotificationAsReadResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Access denied",
					Details:   []string{"you can only mark your own notifications as read"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Update notification as read
	query := `
		UPDATE chat_notifications 
		SET is_read = true 
		WHERE id = $1 AND user_id = $2
	`

	result, err := s.db.DB.Exec(query, req.NotificationId, claims.UserID)
	if err != nil {
		log.Printf("Error marking notification as read: %v", err)
		return &pb.MarkNotificationAsReadResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.MarkNotificationAsReadResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to mark notification as read",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if any rows were affected
	sqlResult, ok := result.(sql.Result)
	if !ok {
		return &pb.MarkNotificationAsReadResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.MarkNotificationAsReadResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to get SQL result",
					Details:   []string{"Invalid result type"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	rowsAffected, err := sqlResult.RowsAffected()
	if err != nil {
		return &pb.MarkNotificationAsReadResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.MarkNotificationAsReadResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to check update result",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if rowsAffected == 0 {
		return &pb.MarkNotificationAsReadResponse{
			StatusCode: 404,
			Message:    "Not Found",
			Result: &pb.MarkNotificationAsReadResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      404,
					Message:   "Notification not found",
					Details:   []string{"notification with given ID does not exist or access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.MarkNotificationAsReadResponse{
		StatusCode: 200,
		Message:    "Notification marked as read successfully",
		Result: &pb.MarkNotificationAsReadResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Notification marked as read",
			},
		},
	}, nil
}

func (s *ChatServer) GetUnreadNotificationCount(ctx context.Context, req *pb.GetUnreadNotificationCountRequest) (*pb.GetUnreadNotificationCountResponse, error) {
	log.Printf("GetUnreadNotificationCount called for user: %s", req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetUnreadNotificationCountResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetUnreadNotificationCountResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check authorization if userId is provided
	if req.UserId != "" && req.UserId != claims.UserID {
		return &pb.GetUnreadNotificationCountResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.GetUnreadNotificationCountResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Access denied",
					Details:   []string{"you can only access your own notification count"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Count unread notifications
	query := `
		SELECT COUNT(*) 
		FROM chat_notifications 
		WHERE user_id = $1 AND is_read = false
	`

	var count int32
	err = s.db.DB.QueryRow(query, claims.UserID).Scan(&count)
	if err != nil {
		log.Printf("Error counting unread notifications: %v", err)
		return &pb.GetUnreadNotificationCountResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.GetUnreadNotificationCountResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to count unread notifications",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.GetUnreadNotificationCountResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Found %d unread notifications", count),
		Result: &pb.GetUnreadNotificationCountResponse_Count{
			Count: count,
		},
	}, nil
}

func (s *ChatServer) AddScheduleMessage(ctx context.Context, req *pb.AddScheduleMessageRequest) (*pb.AddScheduleMessageResponse, error) {
	log.Printf("AddScheduleMessage called by sender: %s", req.SenderId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.AddScheduleMessageResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.AddScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate sender ID matches token
	if req.SenderId != claims.UserID {
		return &pb.AddScheduleMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.AddScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Sender ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate input
	if !ValidateUserID(req.SenderId) {
		return &pb.AddScheduleMessageResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.AddScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid sender ID format",
					Details:   []string{"Sender ID must be a valid UUID"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if !ValidateMessageContent(req.Content) {
		return &pb.AddScheduleMessageResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.AddScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid message content",
					Details:   []string{"Message content cannot be empty and must be less than 4000 characters"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if !ValidateUserID(req.ChatId) {
		return &pb.AddScheduleMessageResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.AddScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid chat ID format",
					Details:   []string{"Chat ID must be a valid UUID"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate scheduled time is in the future
	scheduledTime := req.ScheduledAt.AsTime()
	if scheduledTime.Before(time.Now()) {
		return &pb.AddScheduleMessageResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.AddScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Scheduled time must be in the future",
					Details:   []string{"Cannot schedule messages for past times"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Generate scheduled message ID
	scheduledMessageID := uuid.New().String()
	now := time.Now()

	// Insert scheduled message into database
	query := `
		INSERT INTO scheduled_messages (
			id, chat_id, sender_id, content, message_type, scheduled_time, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	var savedID string
	err = s.db.DB.QueryRow(
		query,
		scheduledMessageID,
		req.ChatId,
		req.SenderId,
		req.Content,
		req.Type.String(),
		scheduledTime,
		now,
	).Scan(&savedID)

	if err != nil {
		log.Printf("Database error: %v", err)
		return &pb.AddScheduleMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.AddScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to schedule message",
					Details:   []string{"Database operation failed"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Handle attachments if provided
	if len(req.Attachments) > 0 {
		for _, attachment := range req.Attachments {
			attachmentQuery := `
				INSERT INTO scheduled_message_attachments (
					scheduled_message_id, file_name, file_type, file_size, file_url, created_at
				) VALUES ($1, $2, $3, $4, $5, $6)`

			_, err = s.db.DB.Exec(
				attachmentQuery,
				savedID,
				"attachment", // Default file name
				"unknown",    // Default file type
				0,            // Default file size
				attachment,
				now,
			)
			if err != nil {
				log.Printf("Failed to save attachment: %v", err)
			}
		}
	}

	return &pb.AddScheduleMessageResponse{
		StatusCode: 200,
		Message:    "Message scheduled successfully",
		Result: &pb.AddScheduleMessageResponse_ScheduledMessageId{
			ScheduledMessageId: savedID,
		},
	}, nil
}

func (s *ChatServer) UpdateScheduleMessage(ctx context.Context, req *pb.UpdateScheduleMessageRequest) (*pb.UpdateScheduleMessageResponse, error) {
	log.Printf("UpdateScheduleMessage called for message: %s by user: %s", req.ScheduledMessageId, req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.UpdateScheduleMessageResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.UpdateScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.UpdateScheduleMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.UpdateScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if !ValidateUserID(req.ScheduledMessageId) || !ValidateMessageContent(req.Content) {
		return &pb.UpdateScheduleMessageResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.UpdateScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Invalid input",
					Details:   []string{"Invalid scheduled message ID or content"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if user owns the scheduled message
	var senderID string
	var isSent bool
	err = s.db.DB.QueryRow("SELECT sender_id, is_sent FROM scheduled_messages WHERE id = $1", req.ScheduledMessageId).Scan(&senderID, &isSent)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UpdateScheduleMessageResponse{
				StatusCode: 404,
				Message:    "Scheduled message not found",
				Result: &pb.UpdateScheduleMessageResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      404,
						Message:   "Scheduled message not found",
						Details:   []string{"Scheduled message does not exist"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
		return &pb.UpdateScheduleMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.UpdateScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Database error",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if senderID != req.UserId {
		return &pb.UpdateScheduleMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.UpdateScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Cannot update scheduled message from another user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if isSent {
		return &pb.UpdateScheduleMessageResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.UpdateScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Cannot update already sent scheduled message",
					Details:   []string{"Message has already been sent"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Validate scheduled time is in the future
	scheduledTime := req.ScheduledAt.AsTime()
	if scheduledTime.Before(time.Now()) {
		return &pb.UpdateScheduleMessageResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.UpdateScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Scheduled time must be in the future",
					Details:   []string{"Cannot schedule messages for past times"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Update scheduled message
	_, err = s.db.DB.Exec(
		"UPDATE scheduled_messages SET content = $1, scheduled_time = $2 WHERE id = $3",
		req.Content, scheduledTime, req.ScheduledMessageId,
	)
	if err != nil {
		return &pb.UpdateScheduleMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.UpdateScheduleMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to update scheduled message",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.UpdateScheduleMessageResponse{
		StatusCode: 200,
		Message:    "Scheduled message updated successfully",
		Result: &pb.UpdateScheduleMessageResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Scheduled message updated",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (s *ChatServer) CancelScheduledMessage(ctx context.Context, req *pb.CancelScheduledMessageRequest) (*pb.CancelScheduledMessageResponse, error) {
	log.Printf("CancelScheduledMessage called for message: %s by user: %s", req.ScheduledMessageId, req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.CancelScheduledMessageResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.CancelScheduledMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.CancelScheduledMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.CancelScheduledMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Check if user owns the scheduled message
	var senderID string
	var isSent bool
	err = s.db.DB.QueryRow("SELECT sender_id, is_sent FROM scheduled_messages WHERE id = $1", req.ScheduledMessageId).Scan(&senderID, &isSent)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.CancelScheduledMessageResponse{
				StatusCode: 404,
				Message:    "Scheduled message not found",
				Result: &pb.CancelScheduledMessageResponse_Error{
					Error: &pb.ErrorMessage{
						Code:      404,
						Message:   "Scheduled message not found",
						Details:   []string{"Scheduled message does not exist"},
						Timestamp: timestamppb.New(time.Now()),
					},
				},
			}, nil
		}
		return &pb.CancelScheduledMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.CancelScheduledMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Database error",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if senderID != req.UserId {
		return &pb.CancelScheduledMessageResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.CancelScheduledMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "Cannot cancel scheduled message from another user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if isSent {
		return &pb.CancelScheduledMessageResponse{
			StatusCode: 400,
			Message:    "Bad Request",
			Result: &pb.CancelScheduledMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      400,
					Message:   "Cannot cancel already sent scheduled message",
					Details:   []string{"Message has already been sent"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Delete scheduled message and its attachments
	_, err = s.db.DB.Exec("DELETE FROM scheduled_message_attachments WHERE scheduled_message_id = $1", req.ScheduledMessageId)
	if err != nil {
		log.Printf("Failed to delete scheduled message attachments: %v", err)
	}

	_, err = s.db.DB.Exec("DELETE FROM scheduled_messages WHERE id = $1", req.ScheduledMessageId)
	if err != nil {
		return &pb.CancelScheduledMessageResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.CancelScheduledMessageResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to cancel scheduled message",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	return &pb.CancelScheduledMessageResponse{
		StatusCode: 200,
		Message:    "Scheduled message cancelled successfully",
		Result: &pb.CancelScheduledMessageResponse_Status{
			Status: &pb.StateMessage{
				StatusCode: 200,
				Message:    "Scheduled message cancelled",
				Timestamp:  timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (s *ChatServer) GetScheduledMessages(ctx context.Context, req *pb.GetScheduledMessagesRequest) (*pb.GetScheduledMessagesResponse, error) {
	log.Printf("GetScheduledMessages called by user: %s", req.UserId)

	claims, err := s.validateRequest(ctx)
	if err != nil {
		return &pb.GetScheduledMessagesResponse{
			StatusCode: 401,
			Message:    "Unauthorized",
			Result: &pb.GetScheduledMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      401,
					Message:   "Token validation failed",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	if req.UserId != claims.UserID {
		return &pb.GetScheduledMessagesResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.GetScheduledMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      403,
					Message:   "User ID does not match authenticated user",
					Details:   []string{"Access denied"},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}

	// Build query based on whether chatId is provided
	var query string
	var args []interface{}

	if req.ChatId != "" {
		query = `
			SELECT id, chat_id, sender_id, content, message_type, scheduled_time, created_at
			FROM scheduled_messages 
			WHERE sender_id = $1 AND chat_id = $2 AND is_sent = false
			ORDER BY scheduled_time ASC`
		args = []interface{}{req.UserId, req.ChatId}
	} else {
		query = `
			SELECT id, chat_id, sender_id, content, message_type, scheduled_time, created_at
			FROM scheduled_messages 
			WHERE sender_id = $1 AND is_sent = false
			ORDER BY scheduled_time ASC`
		args = []interface{}{req.UserId}
	}

	rows, err := s.db.DB.Query(query, args...)
	if err != nil {
		return &pb.GetScheduledMessagesResponse{
			StatusCode: 500,
			Message:    "Internal Server Error",
			Result: &pb.GetScheduledMessagesResponse_Error{
				Error: &pb.ErrorMessage{
					Code:      500,
					Message:   "Failed to fetch scheduled messages",
					Details:   []string{err.Error()},
					Timestamp: timestamppb.New(time.Now()),
				},
			},
		}, nil
	}
	defer rows.Close()

	var scheduledMessages []*pb.ScheduledMessage
	for rows.Next() {
		var msg pb.ScheduledMessage
		var scheduledTime, createdAt time.Time
		var messageType string

		err := rows.Scan(
			&msg.Id, &msg.ChatId, &msg.SenderId, &msg.Content,
			&messageType, &scheduledTime, &createdAt,
		)
		if err != nil {
			continue
		}

		// Convert message type string to enum
		switch messageType {
		case "MSG_TEXT":
			msg.Type = pb.MessageType_MSG_TEXT
		case "MSG_IMAGE":
			msg.Type = pb.MessageType_MSG_IMAGE
		case "MSG_VIDEO":
			msg.Type = pb.MessageType_MSG_VIDEO
		case "MSG_FILE":
			msg.Type = pb.MessageType_MSG_FILE
		case "MSG_AUDIO":
			msg.Type = pb.MessageType_MSG_AUDIO
		default:
			msg.Type = pb.MessageType_MSG_TEXT
		}

		msg.ScheduledAt = timestamppb.New(scheduledTime)
		msg.CreatedAt = timestamppb.New(createdAt)

		// Fetch attachments for this scheduled message
		attachmentQuery := `
			SELECT file_url FROM scheduled_message_attachments 
			WHERE scheduled_message_id = $1`
		attachmentRows, err := s.db.DB.Query(attachmentQuery, msg.Id)
		if err == nil {
			var attachments []string
			for attachmentRows.Next() {
				var fileUrl string
				if err := attachmentRows.Scan(&fileUrl); err == nil {
					attachments = append(attachments, fileUrl)
				}
			}
			attachmentRows.Close()
			msg.Attachments = attachments
		}

		scheduledMessages = append(scheduledMessages, &msg)
	}

	return &pb.GetScheduledMessagesResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Found %d scheduled messages", len(scheduledMessages)),
		Result: &pb.GetScheduledMessagesResponse_Messages{
			Messages: &pb.ScheduledMessageList{
				Messages: scheduledMessages,
			},
		},
	}, nil
}
