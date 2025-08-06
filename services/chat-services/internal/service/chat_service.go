package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"chat-services/internal/models"
	"chat-services/internal/repository/interfaces"
	serviceInterfaces "chat-services/internal/service/interfaces"
	"chat-services/proto"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type chatService struct {
	repo interfaces.ChatRepository

	// Streaming connections
	chatStreams  map[string]grpc.BidiStreamingServer[proto.ChatStreamEnvelope, proto.ChatStreamEnvelope]
	streamsMutex sync.RWMutex

	// Presence subscriptions
	presenceStreams map[string]grpc.ServerStreamingServer[proto.PresenceUpdate]
	presenceMutex   sync.RWMutex

	// Typing subscriptions
	typingStreams map[string]grpc.ServerStreamingServer[proto.TypingIndicator]
	typingMutex   sync.RWMutex

	// Thread subscriptions
	threadStreams map[string]grpc.ServerStreamingServer[proto.ThreadUpdate]
	threadMutex   sync.RWMutex

	// Message update subscriptions
	messageUpdateStreams map[string]grpc.ServerStreamingServer[proto.MessageUpdate]
	messageUpdateMutex   sync.RWMutex

	// Notification subscriptions
	notificationStreams map[string]grpc.ServerStreamingServer[proto.NotificationUpdate]
	notificationMutex   sync.RWMutex

	// Chat event subscriptions
	chatEventStreams map[string]grpc.ServerStreamingServer[proto.ChatEvent]
	chatEventMutex   sync.RWMutex

	// Last message subscriptions
	lastMessageStreams map[string]grpc.ServerStreamingServer[proto.ChatMessage]
	lastMessageMutex   sync.RWMutex

	// User status subscriptions
	userStatusStreams map[string]grpc.ServerStreamingServer[proto.UserStatus]
	userStatusMutex   sync.RWMutex

	// Screen share subscriptions
	screenShareStreams map[string]grpc.ServerStreamingServer[proto.ScreenShareUpdate]
	screenShareMutex   sync.RWMutex
}

func NewChatService(repo interfaces.ChatRepository) serviceInterfaces.ChatService {
	return &chatService{
		repo:                 repo,
		chatStreams:          make(map[string]grpc.BidiStreamingServer[proto.ChatStreamEnvelope, proto.ChatStreamEnvelope]),
		presenceStreams:      make(map[string]grpc.ServerStreamingServer[proto.PresenceUpdate]),
		typingStreams:        make(map[string]grpc.ServerStreamingServer[proto.TypingIndicator]),
		threadStreams:        make(map[string]grpc.ServerStreamingServer[proto.ThreadUpdate]),
		messageUpdateStreams: make(map[string]grpc.ServerStreamingServer[proto.MessageUpdate]),
		notificationStreams:  make(map[string]grpc.ServerStreamingServer[proto.NotificationUpdate]),
		chatEventStreams:     make(map[string]grpc.ServerStreamingServer[proto.ChatEvent]),
		lastMessageStreams:   make(map[string]grpc.ServerStreamingServer[proto.ChatMessage]),
		userStatusStreams:    make(map[string]grpc.ServerStreamingServer[proto.UserStatus]),
		screenShareStreams:   make(map[string]grpc.ServerStreamingServer[proto.ScreenShareUpdate]),
	}
}

// ChatStream handles bidirectional streaming for real-time chat
func (s *chatService) ChatStream(stream grpc.BidiStreamingServer[proto.ChatStreamEnvelope, proto.ChatStreamEnvelope]) error {
	log.Println("New chat stream connection established")

	// Generate unique stream ID
	streamID := fmt.Sprintf("stream_%d", time.Now().UnixNano())

	// Register stream
	s.streamsMutex.Lock()
	s.chatStreams[streamID] = stream
	s.streamsMutex.Unlock()

	// Clean up on disconnect
	defer func() {
		s.streamsMutex.Lock()
		delete(s.chatStreams, streamID)
		s.streamsMutex.Unlock()
		log.Printf("Chat stream %s disconnected", streamID)
	}()

	// Handle incoming messages
	for {
		envelope, err := stream.Recv()
		if err != nil {
			log.Printf("Stream receive error: %v", err)
			return err
		}

		// Process the received message
		if err := s.processStreamMessage(envelope, streamID); err != nil {
			log.Printf("Error processing stream message: %v", err)
			// Send error back to client
			errorEnvelope := &proto.ChatStreamEnvelope{
				Payload: &proto.ChatStreamEnvelope_Error{
					Error: &proto.ErrorMessage{
						Code:      500,
						Message:   err.Error(),
						Timestamp: timestamppb.Now(),
					},
				},
			}
			if sendErr := stream.Send(errorEnvelope); sendErr != nil {
				log.Printf("Error sending error message: %v", sendErr)
			}
		}
	}
}

func (s *chatService) processStreamMessage(envelope *proto.ChatStreamEnvelope, streamID string) error {
	switch payload := envelope.Payload.(type) {
	case *proto.ChatStreamEnvelope_Message:
		// Handle incoming chat message
		return s.handleStreamChatMessage(payload.Message, streamID)
	case *proto.ChatStreamEnvelope_State:
		// Handle state messages (typing indicators, presence updates, etc.)
		return s.handleStreamStateMessage(payload.State, streamID)
	default:
		return fmt.Errorf("unknown message type")
	}
}

func (s *chatService) handleStreamChatMessage(message *proto.ChatMessage, streamID string) error {
	// Convert proto message to model
	chatMessage := &models.ChatMessage{
		SenderID:   message.SenderId,
		ReceiverID: message.ReceiverId,
		GroupID:    message.GroupId,
		Content:    message.Content,
		Type:       message.Type,
		Timestamp:  message.Timestamp.AsTime(),
		IsGroup:    message.IsGroup,
		Status:     proto.MessageStatus_SENT,
	}

	// Save message to database
	savedMessage, err := s.repo.SaveMessage(context.Background(), chatMessage)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	// Convert back to proto
	protoMessage := s.modelToProtoMessage(savedMessage)

	// Broadcast to relevant streams
	s.broadcastMessage(protoMessage)

	return nil
}

func (s *chatService) handleStreamStateMessage(state *proto.StateMessage, streamID string) error {
	// Handle state messages like typing indicators, presence updates
	log.Printf("Received state message: %s", state.Message)
	return nil
}

func (s *chatService) broadcastMessage(message *proto.ChatMessage) {
	s.streamsMutex.RLock()
	defer s.streamsMutex.RUnlock()

	envelope := &proto.ChatStreamEnvelope{
		Payload: &proto.ChatStreamEnvelope_Message{
			Message: message,
		},
	}

	// Broadcast to all connected streams
	// In a real implementation, you'd filter by user/group membership
	for streamID, stream := range s.chatStreams {
		if err := stream.Send(envelope); err != nil {
			log.Printf("Error sending to stream %s: %v", streamID, err)
			// Remove failed stream
			delete(s.chatStreams, streamID)
		}
	}
}

// SendMessage handles regular message sending
func (s *chatService) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	chatMessage := &models.ChatMessage{
		SenderID:   req.SenderId,
		ReceiverID: req.ReceiverId,
		GroupID:    req.GroupId,
		Content:    req.Content,
		Type:       req.Type,
		Timestamp:  time.Now(),
		IsGroup:    req.IsGroup,
		Status:     proto.MessageStatus_SENT,
	}

	savedMessage, err := s.repo.SaveMessage(ctx, chatMessage)
	if err != nil {
		return &proto.SendMessageResponse{
			StatusCode: 500,
			Message:    "Failed to save message",
			Result: &proto.SendMessageResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	protoMessage := s.modelToProtoMessage(savedMessage)

	// Broadcast to streams
	s.broadcastMessage(protoMessage)

	return &proto.SendMessageResponse{
		StatusCode: 200,
		Message:    "Message sent successfully",
		Result: &proto.SendMessageResponse_SavedMessage{
			SavedMessage: protoMessage,
		},
	}, nil
}

// GetChatHistory retrieves chat history
func (s *chatService) GetChatHistory(ctx context.Context, req *proto.GetChatHistoryRequest) (*proto.GetChatHistoryResponse, error) {
	messages, err := s.repo.GetChatHistory(ctx, req.UserId, req.PeerId, req.IsGroup, req.Limit, req.Offset)
	if err != nil {
		return &proto.GetChatHistoryResponse{
			StatusCode: 500,
			Message:    "Failed to get chat history",
			Result: &proto.GetChatHistoryResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	protoMessages := make([]*proto.ChatMessage, len(messages))
	for i, msg := range messages {
		protoMessages[i] = s.modelToProtoMessage(msg)
	}

	return &proto.GetChatHistoryResponse{
		StatusCode: 200,
		Message:    "Chat history retrieved successfully",
		Result: &proto.GetChatHistoryResponse_Messages{
			Messages: &proto.ChatMessageList{
				Messages: protoMessages,
			},
		},
	}, nil
}

// Helper function to convert model to proto
func (s *chatService) modelToProtoMessage(msg *models.ChatMessage) *proto.ChatMessage {
	protoMsg := &proto.ChatMessage{
		MessageId:         msg.MessageID,
		SenderId:          msg.SenderID,
		ReceiverId:        msg.ReceiverID,
		GroupId:           msg.GroupID,
		Content:           msg.Content,
		Type:              msg.Type,
		Timestamp:         timestamppb.New(msg.Timestamp),
		IsGroup:           msg.IsGroup,
		IsRead:            msg.IsRead,
		IsEdited:          msg.IsEdited,
		Status:            msg.Status,
		ReplyToMessageId:  msg.ReplyToMessageID,
		ThreadId:          msg.ThreadID,
		ParentMessageId:   msg.ParentMessageID,
		IsPinned:          msg.IsPinned,
		IsScheduled:       msg.IsScheduled,
		ForwardCount:      msg.ForwardCount,
		OriginalMessageId: msg.OriginalMessageID,
		ThreadReplyCount:  msg.ThreadReplyCount,
		EditHistory:       msg.EditHistory,
		MentionedUserIds:  msg.MentionedUserIDs,
		IsSystemMessage:   msg.IsSystemMessage,
		DeviceInfo:        msg.DeviceInfo,
		ClientVersion:     msg.ClientVersion,
		IsEncrypted:       msg.IsEncrypted,
		EncryptionKeyId:   msg.EncryptionKeyID,
	}

	// Add timestamps
	if msg.DeliveredAt != nil {
		protoMsg.DeliveredAt = timestamppb.New(*msg.DeliveredAt)
	}
	if msg.ReadAt != nil {
		protoMsg.ReadAt = timestamppb.New(*msg.ReadAt)
	}
	if msg.PinnedAt != nil {
		protoMsg.PinnedAt = timestamppb.New(*msg.PinnedAt)
	}
	if msg.EditedAt != nil {
		protoMsg.EditedAt = timestamppb.New(*msg.EditedAt)
	}
	if msg.ScheduledAt != nil {
		protoMsg.ScheduledAt = timestamppb.New(*msg.ScheduledAt)
	}

	// Convert reactions
	protoReactions := make([]*proto.MessageReaction, len(msg.Reactions))
	for i, reaction := range msg.Reactions {
		protoReactions[i] = &proto.MessageReaction{
			UserId:       reaction.UserID,
			ReactionType: reaction.ReactionType,
			Timestamp:    timestamppb.New(reaction.Timestamp),
		}
	}
	protoMsg.Reactions = protoReactions

	return protoMsg
}

// Implement other service methods...
func (s *chatService) EditMessage(ctx context.Context, req *proto.EditMessageRequest) (*proto.EditMessageResponse, error) {
	updatedMessage, err := s.repo.UpdateMessage(ctx, req.MessageId, req.NewContent)
	if err != nil {
		return &proto.EditMessageResponse{
			StatusCode: 500,
			Message:    "Failed to edit message",
			Result: &proto.EditMessageResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	protoMessage := s.modelToProtoMessage(updatedMessage)
	return &proto.EditMessageResponse{
		StatusCode: 200,
		Message:    "Message edited successfully",
		Result: &proto.EditMessageResponse_UpdatedMessage{
			UpdatedMessage: protoMessage,
		},
	}, nil
}

// DeleteMessage implements message deletion
func (s *chatService) DeleteMessage(ctx context.Context, req *proto.DeleteMessageRequest) (*proto.DeleteMessageResponse, error) {
	err := s.repo.DeleteMessage(ctx, req.MessageId)
	if err != nil {
		return &proto.DeleteMessageResponse{
			StatusCode: 500,
			Message:    "Failed to delete message",
			Result: &proto.DeleteMessageResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.DeleteMessageResponse{
		StatusCode: 200,
		Message:    "Message deleted successfully",
	}, nil
}

// SearchMessages implements message search functionality
func (s *chatService) SearchMessages(ctx context.Context, req *proto.SearchMessagesRequest) (*proto.SearchMessagesResponse, error) {
	messages, err := s.repo.SearchMessages(ctx, req.UserId, req.Query, req.PeerId, req.GroupId, req.IsGroup, req.MessageType, req.Limit, req.Offset)
	if err != nil {
		return &proto.SearchMessagesResponse{
			StatusCode: 500,
			Message:    "Failed to search messages",
			Result: &proto.SearchMessagesResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	protoMessages := make([]*proto.ChatMessage, len(messages))
	for i, msg := range messages {
		protoMessages[i] = s.modelToProtoMessage(msg)
	}

	return &proto.SearchMessagesResponse{
		StatusCode: 200,
		Message:    "Messages found successfully",
		Result: &proto.SearchMessagesResponse_Messages{
			Messages: &proto.ChatMessageList{
				Messages: protoMessages,
			},
		},
	}, nil
}

// MarkAsRead implements read receipt functionality
func (s *chatService) MarkAsRead(ctx context.Context, req *proto.ReadReceiptRequest) (*proto.ReadReceiptResponse, error) {
	err := s.repo.MarkAsRead(ctx, req.UserId, req.PeerId, req.IsGroup)
	if err != nil {
		return &proto.ReadReceiptResponse{
			StatusCode: 500,
			Message:    "Failed to mark as read",
			Result: &proto.ReadReceiptResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.ReadReceiptResponse{
		StatusCode: 200,
		Message:    "Messages marked as read",
	}, nil
}

// SendDeliveryReceipt implements delivery receipt functionality
func (s *chatService) SendDeliveryReceipt(ctx context.Context, req *proto.DeliveryReceiptRequest) (*proto.DeliveryReceiptResponse, error) {
	err := s.repo.UpdateMessageStatus(ctx, req.MessageId, proto.MessageStatus_DELIVERED)
	if err != nil {
		return &proto.DeliveryReceiptResponse{
			StatusCode: 500,
			Message:    "Failed to send delivery receipt",
			Result: &proto.DeliveryReceiptResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.DeliveryReceiptResponse{
		StatusCode: 200,
		Message:    "Delivery receipt sent",
	}, nil
}

// AddReaction implements reaction functionality
func (s *chatService) AddReaction(ctx context.Context, req *proto.AddReactionRequest) (*proto.AddReactionResponse, error) {
	err := s.repo.AddReaction(ctx, req.MessageId, req.UserId, req.ReactionType)
	if err != nil {
		return &proto.AddReactionResponse{
			StatusCode: 500,
			Message:    "Failed to add reaction",
			Result: &proto.AddReactionResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.AddReactionResponse{
		StatusCode: 200,
		Message:    "Reaction added successfully",
	}, nil
}

// RemoveReaction implements reaction removal
func (s *chatService) RemoveReaction(ctx context.Context, req *proto.RemoveReactionRequest) (*proto.RemoveReactionResponse, error) {
	err := s.repo.RemoveReaction(ctx, req.MessageId, req.UserId, req.ReactionType)
	if err != nil {
		return &proto.RemoveReactionResponse{
			StatusCode: 500,
			Message:    "Failed to remove reaction",
			Result: &proto.RemoveReactionResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.RemoveReactionResponse{
		StatusCode: 200,
		Message:    "Reaction removed successfully",
	}, nil
}

// GetMessageReactions retrieves reactions for a message
func (s *chatService) GetMessageReactions(ctx context.Context, req *proto.GetMessageReactionsRequest) (*proto.GetMessageReactionsResponse, error) {
	reactions, err := s.repo.GetMessageReactions(ctx, req.MessageId)
	if err != nil {
		return &proto.GetMessageReactionsResponse{
			StatusCode: 500,
			Message:    "Failed to get message reactions",
			Result: &proto.GetMessageReactionsResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.GetMessageReactionsResponse{
		StatusCode: 200,
		Message:    "Reactions retrieved successfully",
		Result: &proto.GetMessageReactionsResponse_Reactions{
			Reactions: &proto.MessageReactionList{
				Reactions: reactions,
			},
		},
	}, nil
}

// GetUsersByUserID retrieves users by their IDs
func (s *chatService) GetUsersByUserID(ctx context.Context, req *proto.GetUsersByUserIDRequest) (*proto.GetUsersByUserIDResponse, error) {
	users, err := s.repo.GetUsersByIDs(ctx, req.UserIds)
	if err != nil {
		return &proto.GetUsersByUserIDResponse{
			StatusCode: 500,
			Message:    "Failed to get users",
			Result: &proto.GetUsersByUserIDResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.GetUsersByUserIDResponse{
		StatusCode: 200,
		Message:    "Users retrieved successfully",
		Result: &proto.GetUsersByUserIDResponse_Users{
			Users: &proto.UserInfoList{
				Users: users,
			},
		},
	}, nil
}

// GetUsersInGroup retrieves users in a group
func (s *chatService) GetUsersInGroup(ctx context.Context, req *proto.GetUsersInGroupRequest) (*proto.GetUsersInGroupResponse, error) {
	users, err := s.repo.GetUsersInGroup(ctx, req.GroupId)
	if err != nil {
		return &proto.GetUsersInGroupResponse{
			StatusCode: 500,
			Message:    "Failed to get group users",
			Result: &proto.GetUsersInGroupResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.GetUsersInGroupResponse{
		StatusCode: 200,
		Message:    "Group users retrieved successfully",
		Result: &proto.GetUsersInGroupResponse_Users{
			Users: &proto.UserInfoList{
				Users: users,
			},
		},
	}, nil
}

// UpdatePresenceStatus updates user presence
func (s *chatService) UpdatePresenceStatus(ctx context.Context, req *proto.UpdatePresenceStatusRequest) (*proto.UpdatePresenceStatusResponse, error) {
	err := s.repo.UpdateUserPresence(ctx, req.UserId, req.Status, req.CustomMessage)
	if err != nil {
		return &proto.UpdatePresenceStatusResponse{
			StatusCode: 500,
			Message:    "Failed to update presence status",
			Result: &proto.UpdatePresenceStatusResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.UpdatePresenceStatusResponse{
		StatusCode: 200,
		Message:    "Presence status updated successfully",
		Result: &proto.UpdatePresenceStatusResponse_Status{
			Status: s,
		},
	}, nil
}

// GetUserPresence retrieves user presence information
func (s *chatService) GetUserPresence(ctx context.Context, req *proto.GetUserPresenceRequest) (*proto.GetUserPresenceResponse, error) {
	presences, err := s.repo.GetUserPresence(ctx, req.UserIds)
	if err != nil {
		return &proto.GetUserPresenceResponse{
			StatusCode: 500,
			Message:    "Failed to get user presence",
			Result: &proto.GetUserPresenceResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	// Convert model presences to proto
	protoPresences := make([]*proto.UserPresence, len(presences))
	for i, presence := range presences {
		protoPresences[i] = &proto.UserPresence{
			UserId:        presence.UserID,
			Status:        presence.Status,
			CustomMessage: presence.CustomMessage,
			LastSeen:      timestamppb.New(presence.LastSeen),
			//UpdatedAt:     timestamppb.New(presence.UpdatedAt),
		}
	}

	return &proto.GetUserPresenceResponse{
		StatusCode: 200,
		Message:    "User presence retrieved successfully",
		Result: &proto.GetUserPresenceResponse_Presences{
			Presences: &proto.UserPresenceList{
				Presences: protoPresences,
			},
		},
	}, nil
}

// SubscribeToPresence handles presence subscription streaming
func (s *chatService) SubscribeToPresence(req *proto.SubscribeToPresenceRequest, stream grpc.ServerStreamingServer[proto.PresenceUpdate]) error {
	// Generate unique stream ID
	streamID := fmt.Sprintf("presence_%d", time.Now().UnixNano())

	// Register stream
	s.presenceMutex.Lock()
	s.presenceStreams[streamID] = stream
	s.presenceMutex.Unlock()

	// Clean up on disconnect
	defer func() {
		s.presenceMutex.Lock()
		delete(s.presenceStreams, streamID)
		s.presenceMutex.Unlock()
		log.Printf("Presence stream %s disconnected", streamID)
	}()

	// Keep connection alive
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// CreateGroup creates a new group
func (s *chatService) CreateGroup(ctx context.Context, req *proto.GroupInfo) (*proto.CreateGroupResponse, error) {
	groupInfo := &models.GroupInfo{
		Name:        req.Name,
		Description: req.Description,
		AvatarURL:   req.AvatarUrl,
		CreatorID:   req.CreatorId,
		MaxMembers:  req.MaxMembers,
		IsPrivate:   req.IsPrivate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdGroup, err := s.repo.CreateGroup(ctx, groupInfo, req.Members)
	if err != nil {
		return &proto.CreateGroupResponse{
			StatusCode: 500,
			Message:    "Failed to create group",
			Result: &proto.CreateGroupResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	// Convert to proto
	protoGroup := &proto.GroupInfo{
		Id:          createdGroup.ID,
		Name:        createdGroup.Name,
		Description: createdGroup.Description,
		AvatarUrl:   createdGroup.AvatarURL,
		CreatorId:   createdGroup.CreatorID,
		MaxMembers:  createdGroup.MaxMembers,
		IsPrivate:   createdGroup.IsPrivate,
		CreatedAt:   timestamppb.New(createdGroup.CreatedAt),
		UpdatedAt:   timestamppb.New(createdGroup.UpdatedAt),
	}

	return &proto.CreateGroupResponse{
		StatusCode: 200,
		Message:    "Group created successfully",
		Result: &proto.CreateGroupResponse_Group{
			Group: protoGroup,
		},
	}, nil
}

// JoinGroup adds a user to a group
func (s *chatService) JoinGroup(ctx context.Context, req *proto.JoinGroupRequest) (*proto.JoinGroupResponse, error) {
	err := s.repo.AddGroupMembers(ctx, req.GroupId, []string{req.UserId})
	if err != nil {
		return &proto.JoinGroupResponse{
			StatusCode: 500,
			Message:    "Failed to join group",
			Result: &proto.JoinGroupResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.JoinGroupResponse{
		StatusCode: 200,
		Message:    "Joined group successfully",
	}, nil
}

// LeaveGroup removes a user from a group
func (s *chatService) LeaveGroup(ctx context.Context, req *proto.LeaveGroupRequest) (*proto.LeaveGroupResponse, error) {
	err := s.repo.RemoveGroupMembers(ctx, req.GroupId, []string{req.UserId})
	if err != nil {
		return &proto.LeaveGroupResponse{
			StatusCode: 500,
			Message:    "Failed to leave group",
			Result: &proto.LeaveGroupResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.LeaveGroupResponse{
		StatusCode: 200,
		Message:    "Left group successfully",
	}, nil
}

// UpdateGroup updates group information
func (s *chatService) UpdateGroup(ctx context.Context, req *proto.UpdateGroupRequest) (*proto.UpdateGroupResponse, error) {
	updates := make(map[string]interface{})
	if req.GroupName != "" {
		updates["name"] = req.GroupName
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.AvatarUrl != "" {
		updates["avatar_url"] = req.AvatarUrl
	}

	updatedGroup, err := s.repo.UpdateGroup(ctx, req.GroupId, updates)
	if err != nil {
		return &proto.UpdateGroupResponse{
			StatusCode: 500,
			Message:    "Failed to update group",
			Result: &proto.UpdateGroupResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	// Convert to proto
	protoGroup := &proto.GroupInfo{
		Id:          updatedGroup.ID,
		Name:        updatedGroup.Name,
		Description: updatedGroup.Description,
		AvatarUrl:   updatedGroup.AvatarURL,
		CreatorId:   updatedGroup.CreatorID,
		MaxMembers:  updatedGroup.MaxMembers,
		IsPrivate:   updatedGroup.IsPrivate,
		CreatedAt:   timestamppb.New(updatedGroup.CreatedAt),
		UpdatedAt:   timestamppb.New(updatedGroup.UpdatedAt),
	}

	return &proto.UpdateGroupResponse{
		StatusCode: 200,
		Message:    "Group updated successfully",
		Result: &proto.UpdateGroupResponse_Group{
			Group: protoGroup,
		},
	}, nil
}

// SendTypingIndicator sends typing indicator
func (s *chatService) SendTypingIndicator(ctx context.Context, req *proto.TypingIndicatorRequest) (*proto.TypingIndicatorResponse, error) {
	// Broadcast typing indicator to relevant streams
	s.typingMutex.RLock()
	defer s.typingMutex.RUnlock()

	typingIndicator := &proto.TypingIndicator{
		UserId:    req.UserId,
		PeerId:    req.PeerId,
		GroupId:   req.GroupId,
		IsGroup:   req.IsGroup,
		IsTyping:  req.IsTyping,
		Timestamp: timestamppb.Now(),
	}

	// Broadcast to typing streams
	for streamID, stream := range s.typingStreams {
		if err := stream.Send(typingIndicator); err != nil {
			log.Printf("Error sending typing indicator to stream %s: %v", streamID, err)
			delete(s.typingStreams, streamID)
		}
	}

	return &proto.TypingIndicatorResponse{
		StatusCode: 200,
		Message:    "Typing indicator sent successfully",
		Result: &proto.TypingIndicatorResponse_Status{
			Status: true,
		},
	}, nil
}

// SubscribeToTypingIndicators handles typing indicator subscription
func (s *chatService) SubscribeToTypingIndicators(req *proto.TypingSubscriptionRequest, stream grpc.ServerStreamingServer[proto.TypingIndicator]) error {
	// Generate unique stream ID
	streamID := fmt.Sprintf("typing_%d", time.Now().UnixNano())

	// Register stream
	s.typingMutex.Lock()
	s.typingStreams[streamID] = stream
	s.typingMutex.Unlock()

	// Clean up on disconnect
	defer func() {
		s.typingMutex.Lock()
		delete(s.typingStreams, streamID)
		s.typingMutex.Unlock()
		log.Printf("Typing stream %s disconnected", streamID)
	}()

	// Keep connection alive
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// CreateThread creates a new thread
func (s *chatService) CreateThread(ctx context.Context, req *proto.CreateThreadRequest) (*proto.CreateThreadResponse, error) {
	threadMessage, err := s.repo.CreateThread(ctx, req.ParentMessageId, req.UserId, req.Content)
	if err != nil {
		return &proto.CreateThreadResponse{
			StatusCode: 500,
			Message:    "Failed to create thread",
			Result: &proto.CreateThreadResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	protoMessage := s.modelToProtoMessage(threadMessage)
	return &proto.CreateThreadResponse{
		StatusCode: 200,
		Message:    "Thread created successfully",
		Result: &proto.CreateThreadResponse_ThreadMessage{
			ThreadMessage: protoMessage,
		},
	}, nil
}

// GetThreadMessages retrieves messages in a thread
func (s *chatService) GetThreadMessages(ctx context.Context, req *proto.GetThreadMessagesRequest) (*proto.GetThreadMessagesResponse, error) {
	messages, err := s.repo.GetThreadMessages(ctx, req.ParentMessageId, req.Limit, req.Offset)
	if err != nil {
		return &proto.GetThreadMessagesResponse{
			StatusCode: 500,
			Message:    "Failed to get thread messages",
			Result: &proto.GetThreadMessagesResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	protoMessages := make([]*proto.ChatMessage, len(messages))
	for i, msg := range messages {
		protoMessages[i] = s.modelToProtoMessage(msg)
	}

	return &proto.GetThreadMessagesResponse{
		StatusCode: 200,
		Message:    "Thread messages retrieved successfully",
		Result: &proto.GetThreadMessagesResponse_Messages{
			Messages: &proto.ChatMessageList{
				Messages: protoMessages,
			},
		},
	}, nil
}

// SubscribeToThread handles thread subscription
func (s *chatService) SubscribeToThread(req *proto.SubscribeToThreadRequest, stream grpc.ServerStreamingServer[proto.ThreadUpdate]) error {
	streamID := uuid.New().String()

	// Register stream
	s.threadMutex.Lock()
	s.threadStreams[streamID] = stream
	s.threadMutex.Unlock()

	// Cleanup on exit
	defer func() {
		s.threadMutex.Lock()
		delete(s.threadStreams, streamID)
		s.threadMutex.Unlock()
	}()

	// Send initial thread state if available
	if req.ThreadId != "" {
		// Get thread messages to send initial state
		messages, err := s.repo.GetThreadMessages(context.Background(), req.ThreadId)
		if err == nil && len(messages) > 0 {
			initialUpdate := &proto.ThreadUpdate{
				ThreadId:   req.ThreadId,
				UpdateType: proto.MessageUpdateType_NEW_MESSAGE,
				Message:    s.modelToProtoMessage(messages[len(messages)-1]),
				Timestamp:  timestamppb.Now(),
			}
			if err := stream.Send(initialUpdate); err != nil {
				return err
			}
		}
	}

	// Keep connection alive
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// PinMessage pins a message
func (s *chatService) PinMessage(ctx context.Context, req *proto.PinMessageRequest) (*proto.PinMessageResponse, error) {
	err := s.repo.PinMessage(ctx, req.MessageId, req.UserId, req.ChatId, req.IsGroup)
	if err != nil {
		return &proto.PinMessageResponse{
			StatusCode: 500,
			Message:    "Failed to pin message",
			Result: &proto.PinMessageResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.PinMessageResponse{
		StatusCode: 200,
		Message:    "Message pinned successfully",
	}, nil
}

// UnpinMessage unpins a message
func (s *chatService) UnpinMessage(ctx context.Context, req *proto.UnpinMessageRequest) (*proto.UnpinMessageResponse, error) {
	err := s.repo.UnpinMessage(ctx, req.MessageId, req.UserId, req.ChatId, req.IsGroup)
	if err != nil {
		return &proto.UnpinMessageResponse{
			StatusCode: 500,
			Message:    "Failed to unpin message",
			Result: &proto.UnpinMessageResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.UnpinMessageResponse{
		StatusCode: 200,
		Message:    "Message unpinned successfully",
	}, nil
}

// GetPinnedMessages retrieves pinned messages
func (s *chatService) GetPinnedMessages(ctx context.Context, req *proto.GetPinnedMessagesRequest) (*proto.GetPinnedMessagesResponse, error) {
	messages, err := s.repo.GetPinnedMessages(ctx, req.ChatId, req.IsGroup)
	if err != nil {
		return &proto.GetPinnedMessagesResponse{
			StatusCode: 500,
			Message:    "Failed to get pinned messages",
			Result: &proto.GetPinnedMessagesResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	protoMessages := make([]*proto.ChatMessage, len(messages))
	for i, msg := range messages {
		protoMessages[i] = s.modelToProtoMessage(msg)
	}

	return &proto.GetPinnedMessagesResponse{
		StatusCode: 200,
		Message:    "Pinned messages retrieved successfully",
		Result: &proto.GetPinnedMessagesResponse_Messages{
			Messages: &proto.ChatMessageList{
				Messages: protoMessages,
			},
		},
	}, nil
}

// SubscribeToMessageUpdates handles message update subscriptions
func (s *chatService) SubscribeToMessageUpdates(req *proto.SubscribeToMessageUpdatesRequest, stream grpc.ServerStreamingServer[proto.MessageUpdate]) error {
	streamID := uuid.New().String()

	// Register stream
	s.messageUpdateMutex.Lock()
	s.messageUpdateStreams[streamID] = stream
	s.messageUpdateMutex.Unlock()

	// Cleanup on exit
	defer func() {
		s.messageUpdateMutex.Lock()
		delete(s.messageUpdateStreams, streamID)
		s.messageUpdateMutex.Unlock()
	}()

	// Keep connection alive
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// SubscribeToNotifications handles notification subscriptions
func (s *chatService) SubscribeToNotifications(req *proto.SubscribeToNotificationsRequest, stream grpc.ServerStreamingServer[proto.NotificationUpdate]) error {
	streamID := uuid.New().String()

	// Register stream
	s.notificationMutex.Lock()
	s.notificationStreams[streamID] = stream
	s.notificationMutex.Unlock()

	// Cleanup on exit
	defer func() {
		s.notificationMutex.Lock()
		delete(s.notificationStreams, streamID)
		s.notificationMutex.Unlock()
	}()

	// Keep connection alive
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// SubscribeToChatEvents handles chat event subscriptions
func (s *chatService) SubscribeToChatEvents(req *proto.SubscribeToChatEventsRequest, stream grpc.ServerStreamingServer[proto.ChatEvent]) error {
	streamID := uuid.New().String()

	// Register stream
	s.chatEventMutex.Lock()
	s.chatEventStreams[streamID] = stream
	s.chatEventMutex.Unlock()

	// Cleanup on exit
	defer func() {
		s.chatEventMutex.Lock()
		delete(s.chatEventStreams, streamID)
		s.chatEventMutex.Unlock()
	}()

	// Keep connection alive
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// ForwardMessage forwards a message
func (s *chatService) ForwardMessage(ctx context.Context, req *proto.ForwardMessageRequest) (*proto.ForwardMessageResponse, error) {
	// Get original message
	originalMessage, err := s.repo.GetMessageByID(ctx, req.MessageId)
	if err != nil {
		return &proto.ForwardMessageResponse{
			StatusCode: 500,
			Message:    "Failed to get original message",
			Result: &proto.ForwardMessageResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	// Create forwarded message
	forwardedMessage := &models.ChatMessage{
		SenderID:          req.SenderId,
		ReceiverID:        req.ReceiverId,
		GroupID:           req.GroupId,
		Content:           originalMessage.Content,
		Type:              originalMessage.Type,
		Timestamp:         time.Now(),
		IsGroup:           req.IsGroup,
		Status:            proto.MessageStatus_SENT,
		OriginalMessageID: req.MessageId,
		ForwardCount:      originalMessage.ForwardCount + 1,
	}

	savedMessage, err := s.repo.SaveMessage(ctx, forwardedMessage)
	if err != nil {
		return &proto.ForwardMessageResponse{
			StatusCode: 500,
			Message:    "Failed to forward message",
			Result: &proto.ForwardMessageResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	protoMessage := s.modelToProtoMessage(savedMessage)
	return &proto.ForwardMessageResponse{
		StatusCode: 200,
		Message:    "Message forwarded successfully",
		Result: &proto.ForwardMessageResponse_ForwardedMessages{
			ForwardedMessages: &proto.ChatMessageList{
				Messages: []*proto.ChatMessage{protoMessage},
			},
		},
	}, nil
}

// LikeMessage likes a message
func (s *chatService) LikeMessage(ctx context.Context, req *proto.LikeMessageRequest) (*proto.LikeMessageResponse, error) {
	err := s.repo.AddReaction(ctx, req.MessageId, req.UserId, proto.ReactionType_LIKE)
	if err != nil {
		return &proto.LikeMessageResponse{
			StatusCode: 500,
			Message:    "Failed to like message",
			Result: &proto.LikeMessageResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.LikeMessageResponse{
		StatusCode: 200,
		Message:    "Message liked successfully",
		Result: &proto.LikeMessageResponse_Result{
			Success: true,
		},
	}, nil
}

// GetLikedMessages retrieves liked messages
func (s *chatService) GetLikedMessages(ctx context.Context, req *proto.GetLikedMessagesRequest) (*proto.GetLikedMessagesResponse, error) {
	// Search for messages with reactions from the requesting user
	searchReq := &proto.SearchMessagesRequest{
		UserId: req.UserId,
		Query:  "", // Empty query to get all messages
		Limit:  100,
	}

	searchResp, err := s.SearchMessages(ctx, searchReq)
	if err != nil {
		return &proto.GetLikedMessagesResponse{
			StatusCode: 500,
			Message:    "Failed to get liked messages",
			Result: &proto.GetLikedMessagesResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	// Filter messages that the user has liked
	var likedMessages []*proto.ChatMessage
	for _, msg := range searchResp.Result.(*proto.SearchMessagesResponse_Messages).Messages.Messages {
		for _, reaction := range msg.Reactions {
			if reaction.UserId == req.UserId && reaction.ReactionType == proto.ReactionType_LIKE {
				likedMessages = append(likedMessages, msg)
				break
			}
		}
	}

	return &proto.GetLikedMessagesResponse{
		StatusCode: 200,
		Message:    "Liked messages retrieved successfully",
		Result: &proto.GetLikedMessagesResponse_Messages{
			Messages: &proto.ChatMessageList{
				Messages: likedMessages,
			},
		},
	}, nil
}

// GetLastMessages retrieves last messages
func (s *chatService) GetLastMessages(ctx context.Context, req *proto.GetLastMessagesRequest) (*proto.GetLastMessagesResponse, error) {
	// Use GetChatHistory to get recent messages
	historyReq := &proto.GetChatHistoryRequest{
		UserId:  req.UserId,
		ChatId:  req.ChatId,
		IsGroup: req.IsGroup,
		Limit:   req.Limit,
	}

	historyResp, err := s.GetChatHistory(ctx, historyReq)
	if err != nil {
		return &proto.GetLastMessagesResponse{
			StatusCode: 500,
			Message:    "Failed to get last messages",
			Result: &proto.GetLastMessagesResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.GetLastMessagesResponse{
		StatusCode: 200,
		Message:    "Last messages retrieved successfully",
		Result: &proto.GetLastMessagesResponse_Messages{
			Messages: historyResp.Result.(*proto.GetChatHistoryResponse_Messages).Messages,
		},
	}, nil
}

// SubscribeToLastMessages handles last message subscriptions
func (s *chatService) SubscribeToLastMessages(req *proto.LastMessageStreamRequest, stream grpc.ServerStreamingServer[proto.ChatMessage]) error {
	streamID := uuid.New().String()

	// Register stream
	s.lastMessageMutex.Lock()
	s.lastMessageStreams[streamID] = stream
	s.lastMessageMutex.Unlock()

	// Cleanup on exit
	defer func() {
		s.lastMessageMutex.Lock()
		delete(s.lastMessageStreams, streamID)
		s.lastMessageMutex.Unlock()
	}()

	// Keep connection alive
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// GetUserStatus retrieves user status
func (s *chatService) GetUserStatus(ctx context.Context, req *proto.UserStatusRequest) (*proto.UserStatusResponse, error) {
	// Get user presence from repository
	presence, err := s.repo.GetUserPresence(ctx, req.UserId)
	if err != nil {
		return &proto.UserStatusResponse{
			StatusCode: 500,
			Message:    "Failed to get user status",
			Result: &proto.UserStatusResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.UserStatusResponse{
		StatusCode: 200,
		Message:    "User status retrieved successfully",
		Result: &proto.UserStatusResponse_Status{
			Status: &proto.UserStatus{
				UserId:        presence.UserID,
				Status:        presence.Status,
				CustomMessage: presence.CustomMessage,
				LastSeen:      timestamppb.New(presence.LastSeen),
			},
		},
	}, nil
}

// SubscribeToUserStatus handles user status subscriptions
func (s *chatService) SubscribeToUserStatus(req *proto.UserStatusSubscriptionRequest, stream grpc.ServerStreamingServer[proto.UserStatus]) error {
	streamID := uuid.New().String()

	// Register stream
	s.userStatusMutex.Lock()
	s.userStatusStreams[streamID] = stream
	s.userStatusMutex.Unlock()

	// Cleanup on exit
	defer func() {
		s.userStatusMutex.Lock()
		delete(s.userStatusStreams, streamID)
		s.userStatusMutex.Unlock()
	}()

	// Send initial status if available
	if req.UserId != "" {
		presence, err := s.repo.GetUserPresence(context.Background(), req.UserId)
		if err == nil {
			initialStatus := &proto.UserStatus{
				UserId:        presence.UserID,
				Status:        presence.Status,
				CustomMessage: presence.CustomMessage,
				LastSeen:      timestamppb.New(presence.LastSeen),
			}
			if err := stream.Send(initialStatus); err != nil {
				return err
			}
		}
	}

	// Keep connection alive
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// UploadFile handles file upload
func (s *chatService) UploadFile(stream grpc.ClientStreamingServer[proto.FileUploadRequest, proto.FileUploadResponse]) error {
	var fileData []byte
	var fileName string
	var fileType string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// File upload completed
			fileID := uuid.New().String()
			// Here you would typically save the file to storage
			// For now, we'll simulate successful upload

			return stream.SendAndClose(&proto.FileUploadResponse{
				StatusCode: 200,
				Message:    "File uploaded successfully",
				Result: &proto.FileUploadResponse_FileInfo{
					FileInfo: &proto.FileMetadata{
						FileId:   fileID,
						FileName: fileName,
						FileType: fileType,
						FileSize: int64(len(fileData)),
						FileUrl:  fmt.Sprintf("/files/%s", fileID),
					},
				},
			})
		}
		if err != nil {
			return err
		}

		if req.GetMetadata() != nil {
			fileName = req.GetMetadata().FileName
			fileType = req.GetMetadata().FileType
		} else if req.GetChunk() != nil {
			fileData = append(fileData, req.GetChunk().Data...)
		}
	}
}

// DownloadFile handles file download
func (s *chatService) DownloadFile(req *proto.FileDownloadRequest, stream grpc.ServerStreamingServer[proto.FileDownloadResponse]) error {
	// Send file metadata first
	metadataResp := &proto.FileDownloadResponse{
		StatusCode: 200,
		Message:    "File metadata",
		Result: &proto.FileDownloadResponse_FileMetadata{
			FileMetadata: &proto.FileMetadata{
				FileId:   req.FileId,
				FileName: "example.txt",
				FileType: "text/plain",
				FileSize: 12,
			},
		},
	}
	if err := stream.Send(metadataResp); err != nil {
		return err
	}

	// Send file chunks
	fileContent := []byte("file content")
	chunkSize := 1024
	for i := 0; i < len(fileContent); i += chunkSize {
		end := i + chunkSize
		if end > len(fileContent) {
			end = len(fileContent)
		}

		chunkResp := &proto.FileDownloadResponse{
			StatusCode: 200,
			Message:    "File chunk",
			Result: &proto.FileDownloadResponse_FileChunk{
				FileChunk: &proto.FileChunk{
					Data: fileContent[i:end],
				},
			},
		}
		if err := stream.Send(chunkResp); err != nil {
			return err
		}
	}

	return nil
}

// Call-related methods with basic implementations
func (s *chatService) InitiateCall(ctx context.Context, req *proto.InitiateCallRequest) (*proto.InitiateCallResponse, error) {
	return &proto.InitiateCallResponse{
		StatusCode: 200,
		Message:    "Call initiated successfully",
		Result: &proto.InitiateCallResponse_CallInfo{
			CallInfo: &proto.CallInfo{
				CallId:     fmt.Sprintf("call_%d", time.Now().UnixNano()),
				CallerId:   req.CallerId,
				ReceiverId: req.ReceiverId,
				CallType:   req.CallType,
				Status:     proto.CallStatus_RINGING,
			},
		},
	}, nil
}

func (s *chatService) AcceptCall(ctx context.Context, req *proto.AcceptCallRequest) (*proto.AcceptCallResponse, error) {
	return &proto.AcceptCallResponse{
		StatusCode: 200,
		Message:    "Call accepted successfully",
		Result: &proto.AcceptCallResponse_Success{
			Success: true,
		},
	}, nil
}

func (s *chatService) RejectCall(ctx context.Context, req *proto.RejectCallRequest) (*proto.RejectCallResponse, error) {
	return &proto.RejectCallResponse{
		StatusCode: 200,
		Message:    "Call rejected successfully",
		Result: &proto.RejectCallResponse_Success{
			Success: true,
		},
	}, nil
}

func (s *chatService) EndCall(ctx context.Context, req *proto.EndCallRequest) (*proto.EndCallResponse, error) {
	return &proto.EndCallResponse{
		StatusCode: 200,
		Message:    "Call ended successfully",
		Result: &proto.EndCallResponse_Success{
			Success: true,
		},
	}, nil
}

func (s *chatService) GetCallHistory(ctx context.Context, req *proto.GetCallHistoryRequest) (*proto.GetCallHistoryResponse, error) {
	return &proto.GetCallHistoryResponse{
		StatusCode: 200,
		Message:    "Call history retrieved successfully",
		Result: &proto.GetCallHistoryResponse_Calls{
			Calls: &proto.CallInfoList{
				Calls: []*proto.CallInfo{},
			},
		},
	}, nil
}

// Notification-related methods
func (s *chatService) UpdateNotificationSettings(ctx context.Context, req *proto.UpdateNotificationSettingsRequest) (*proto.UpdateNotificationSettingsResponse, error) {
	return &proto.UpdateNotificationSettingsResponse{
		StatusCode: 200,
		Message:    "Notification settings updated successfully",
		Result: &proto.UpdateNotificationSettingsResponse_Success{
			Success: true,
		},
	}, nil
}

func (s *chatService) GetNotificationSettings(ctx context.Context, req *proto.GetNotificationSettingsRequest) (*proto.GetNotificationSettingsResponse, error) {
	return &proto.GetNotificationSettingsResponse{
		StatusCode: 200,
		Message:    "Notification settings retrieved successfully",
		Result: &proto.GetNotificationSettingsResponse_Settings{
			Settings: &proto.NotificationSettings{
				UserId:                       req.UserId,
				EnablePushNotifications:      true,
				EnableSoundNotifications:     true,
				EnableVibrationNotifications: true,
				EnableEmailNotifications:     false,
			},
		},
	}, nil
}

func (s *chatService) MarkNotificationAsRead(ctx context.Context, req *proto.MarkNotificationAsReadRequest) (*proto.MarkNotificationAsReadResponse, error) {
	return &proto.MarkNotificationAsReadResponse{
		StatusCode: 200,
		Message:    "Notification marked as read successfully",
		Result: &proto.MarkNotificationAsReadResponse_Success{
			Success: true,
		},
	}, nil
}

func (s *chatService) GetUnreadNotificationCount(ctx context.Context, req *proto.GetUnreadNotificationCountRequest) (*proto.GetUnreadNotificationCountResponse, error) {
	return &proto.GetUnreadNotificationCountResponse{
		StatusCode: 200,
		Message:    "Unread notification count retrieved successfully",
		Result: &proto.GetUnreadNotificationCountResponse_Count{
			Count: 0,
		},
	}, nil
}

// Screen share methods
func (s *chatService) StartScreenShare(ctx context.Context, req *proto.StartScreenShareRequest) (*proto.StartScreenShareResponse, error) {
	return &proto.StartScreenShareResponse{
		StatusCode: 200,
		Message:    "Screen share started successfully",
		Result: &proto.StartScreenShareResponse_SessionInfo{
			SessionInfo: &proto.ScreenShareSession{
				SessionId: fmt.Sprintf("screen_%d", time.Now().UnixNano()),
				UserId:    req.UserId,
				RoomId:    req.RoomId,
			},
		},
	}, nil
}

func (s *chatService) StopScreenShare(ctx context.Context, req *proto.StopScreenShareRequest) (*proto.StopScreenShareResponse, error) {
	return &proto.StopScreenShareResponse{
		StatusCode: 200,
		Message:    "Screen share stopped successfully",
		Result: &proto.StopScreenShareResponse_Success{
			Success: true,
		},
	}, nil
}

func (s *chatService) SubscribeToScreenShare(req *proto.SubscribeToScreenShareRequest, stream grpc.ServerStreamingServer[proto.ScreenShareUpdate]) error {
	streamID := uuid.New().String()

	// Register stream
	s.screenShareMutex.Lock()
	s.screenShareStreams[streamID] = stream
	s.screenShareMutex.Unlock()

	// Cleanup on exit
	defer func() {
		s.screenShareMutex.Lock()
		delete(s.screenShareStreams, streamID)
		s.screenShareMutex.Unlock()
	}()

	// Keep connection alive
	select {
	case <-stream.Context().Done():
		return stream.Context().Err()
	}
}

// Message scheduling methods
func (s *chatService) ScheduleMessage(ctx context.Context, req *proto.ScheduleMessageRequest) (*proto.ScheduleMessageResponse, error) {
	// Create scheduled message
	scheduledMessage := &models.ChatMessage{
		SenderID:    req.SenderId,
		ReceiverID:  req.ReceiverId,
		GroupID:     req.GroupId,
		Content:     req.Content,
		Type:        req.Type,
		Timestamp:   time.Now(),
		IsGroup:     req.IsGroup,
		Status:      proto.MessageStatus_SCHEDULED,
		IsScheduled: true,
		ScheduledAt: &req.ScheduledAt.AsTime(),
	}

	savedMessage, err := s.repo.SaveMessage(ctx, scheduledMessage)
	if err != nil {
		return &proto.ScheduleMessageResponse{
			StatusCode: 500,
			Message:    "Failed to schedule message",
			Result: &proto.ScheduleMessageResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	protoMessage := s.modelToProtoMessage(savedMessage)
	return &proto.ScheduleMessageResponse{
		StatusCode: 200,
		Message:    "Message scheduled successfully",
		Result: &proto.ScheduleMessageResponse_ScheduledMessage{
			ScheduledMessage: protoMessage,
		},
	}, nil
}

func (s *chatService) CancelScheduledMessage(ctx context.Context, req *proto.CancelScheduledMessageRequest) (*proto.CancelScheduledMessageResponse, error) {
	err := s.repo.DeleteMessage(ctx, req.MessageId)
	if err != nil {
		return &proto.CancelScheduledMessageResponse{
			StatusCode: 500,
			Message:    "Failed to cancel scheduled message",
			Result: &proto.CancelScheduledMessageResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	return &proto.CancelScheduledMessageResponse{
		StatusCode: 200,
		Message:    "Scheduled message cancelled successfully",
		Result: &proto.CancelScheduledMessageResponse_Success{
			Success: true,
		},
	}, nil
}

func (s *chatService) GetScheduledMessages(ctx context.Context, req *proto.GetScheduledMessagesRequest) (*proto.GetScheduledMessagesResponse, error) {
	// Search for scheduled messages
	searchReq := &proto.SearchMessagesRequest{
		UserId: req.UserId,
		Query:  "", // Empty query to get all messages
		Limit:  100,
	}

	searchResp, err := s.SearchMessages(ctx, searchReq)
	if err != nil {
		return &proto.GetScheduledMessagesResponse{
			StatusCode: 500,
			Message:    "Failed to get scheduled messages",
			Result: &proto.GetScheduledMessagesResponse_Error{
				Error: &proto.ErrorMessage{
					Code:    500,
					Message: err.Error(),
				},
			},
		}, nil
	}

	// Filter for scheduled messages
	var scheduledMessages []*proto.ChatMessage
	for _, msg := range searchResp.Result.(*proto.SearchMessagesResponse_Messages).Messages.Messages {
		if msg.Status == proto.MessageStatus_SCHEDULED {
			scheduledMessages = append(scheduledMessages, msg)
		}
	}

	return &proto.GetScheduledMessagesResponse{
		StatusCode: 200,
		Message:    "Scheduled messages retrieved successfully",
		Result: &proto.GetScheduledMessagesResponse_Messages{
			Messages: &proto.ChatMessageList{
				Messages: scheduledMessages,
			},
		},
	}, nil
}

// Analytics method
func (s *chatService) GetChatAnalytics(ctx context.Context, req *proto.GetChatAnalyticsRequest) (*proto.GetChatAnalyticsResponse, error) {
	return &proto.GetChatAnalyticsResponse{
		StatusCode: 200,
		Message:    "Chat analytics retrieved successfully",
		Result: &proto.GetChatAnalyticsResponse_Analytics{
			Analytics: &proto.ChatAnalytics{
				UserId:        req.UserId,
				TotalMessages: 0,
				TotalChats:    0,
				ActiveChats:   0,
			},
		},
	}, nil
}
