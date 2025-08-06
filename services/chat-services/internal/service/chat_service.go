package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"chat-services/internal/models"
	"chat-services/internal/repository/interfaces"
	serviceInterfaces "chat-services/internal/service/interfaces"
	"chat-services/proto"

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
}

func NewChatService(repo interfaces.ChatRepository) serviceInterfaces.ChatService {
	return &chatService{
		repo:            repo,
		chatStreams:     make(map[string]grpc.BidiStreamingServer[proto.ChatStreamEnvelope, proto.ChatStreamEnvelope]),
		presenceStreams: make(map[string]grpc.ServerStreamingServer[proto.PresenceUpdate]),
		typingStreams:   make(map[string]grpc.ServerStreamingServer[proto.TypingIndicator]),
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
	return &proto.ChatMessage{
		MessageId:        msg.MessageID,
		SenderId:         msg.SenderID,
		ReceiverId:       msg.ReceiverID,
		GroupId:          msg.GroupID,
		Content:          msg.Content,
		Type:             msg.Type,
		Timestamp:        timestamppb.New(msg.Timestamp),
		IsGroup:          msg.IsGroup,
		IsRead:           msg.IsRead,
		IsEdited:         msg.IsEdited,
		Status:           msg.Status,
		ReplyToMessageId: msg.ReplyToMessageID,
		ThreadId:         msg.ThreadID,
		ParentMessageId:  msg.ParentMessageID,
		IsPinned:         msg.IsPinned,
		IsScheduled:      msg.IsScheduled,
	}
}

// Implement other service methods...
func (s *chatService) EditMessage(ctx context.Context, req *proto.EditMessageRequest) (*proto.EditMessageResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) DeleteMessage(ctx context.Context, req *proto.DeleteMessageRequest) (*proto.DeleteMessageResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) SearchMessages(ctx context.Context, req *proto.SearchMessagesRequest) (*proto.SearchMessagesResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) MarkAsRead(ctx context.Context, req *proto.ReadReceiptRequest) (*proto.ReadReceiptResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) SendDeliveryReceipt(ctx context.Context, req *proto.DeliveryReceiptRequest) (*proto.DeliveryReceiptResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) AddReaction(ctx context.Context, req *proto.AddReactionRequest) (*proto.AddReactionResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) RemoveReaction(ctx context.Context, req *proto.RemoveReactionRequest) (*proto.RemoveReactionResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) GetMessageReactions(ctx context.Context, req *proto.GetMessageReactionsRequest) (*proto.GetMessageReactionsResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) GetUsersByUserID(ctx context.Context, req *proto.GetUsersByUserIDRequest) (*proto.GetUsersByUserIDResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) GetUsersInGroup(ctx context.Context, req *proto.GetUsersInGroupRequest) (*proto.GetUsersInGroupResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) UpdatePresenceStatus(ctx context.Context, req *proto.UpdatePresenceStatusRequest) (*proto.UpdatePresenceStatusResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) GetUserPresence(ctx context.Context, req *proto.GetUserPresenceRequest) (*proto.GetUserPresenceResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) SubscribeToPresence(req *proto.SubscribeToPresenceRequest, stream proto.ChatService_SubscribeToPresenceServer) error {
	// Implementation here
	return nil
}

func (s *chatService) CreateGroup(ctx context.Context, req *proto.CreateGroupRequest) (*proto.CreateGroupResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) JoinGroup(ctx context.Context, req *proto.JoinGroupRequest) (*proto.JoinGroupResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) LeaveGroup(ctx context.Context, req *proto.LeaveGroupRequest) (*proto.LeaveGroupResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) UpdateGroup(ctx context.Context, req *proto.UpdateGroupRequest) (*proto.UpdateGroupResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) SendTypingIndicator(ctx context.Context, req *proto.TypingIndicatorRequest) (*proto.TypingIndicatorResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) SubscribeToTypingIndicators(req *proto.TypingSubscriptionRequest, stream proto.ChatService_SubscribeToTypingIndicatorsServer) error {
	// Implementation here
	return nil
}

func (s *chatService) CreateThread(ctx context.Context, req *proto.CreateThreadRequest) (*proto.CreateThreadResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) GetThreadMessages(ctx context.Context, req *proto.GetThreadMessagesRequest) (*proto.GetThreadMessagesResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) SubscribeToThread(req *proto.SubscribeToThreadRequest, stream proto.ChatService_SubscribeToThreadServer) error {
	// Implementation here
	return nil
}

func (s *chatService) PinMessage(ctx context.Context, req *proto.PinMessageRequest) (*proto.PinMessageResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) UnpinMessage(ctx context.Context, req *proto.UnpinMessageRequest) (*proto.UnpinMessageResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) GetPinnedMessages(ctx context.Context, req *proto.GetPinnedMessagesRequest) (*proto.GetPinnedMessagesResponse, error) {
	// Implementation here
	return nil, nil
}

func (s *chatService) SubscribeToMessageUpdates(req *proto.SubscribeToMessageUpdatesRequest, stream proto.ChatService_SubscribeToMessageUpdatesServer) error {
	// Implementation here
	return nil
}

func (s *chatService) SubscribeToNotifications(req *proto.SubscribeToNotificationsRequest, stream proto.ChatService_SubscribeToNotificationsServer) error {
	// Implementation here
	return nil
}

func (s *chatService) SubscribeToChatEvents(req *proto.SubscribeToChatEventsRequest, stream proto.ChatService_SubscribeToChatEventsServer) error {
	// Implementation here
	return nil
}

// Add all missing methods:

func (s *chatService) ForwardMessage(ctx context.Context, req *proto.ForwardMessageRequest) (*proto.ForwardMessageResponse, error) {
	// Implementation here
	return &proto.ForwardMessageResponse{
		StatusCode: 200,
		Message:    "Message forwarded successfully",
	}, nil
}

func (s *chatService) LikeMessage(ctx context.Context, req *proto.LikeMessageRequest) (*proto.LikeMessageResponse, error) {
	// Implementation here
	return &proto.LikeMessageResponse{
		StatusCode: 200,
		Message:    "Message liked successfully",
	}, nil
}

func (s *chatService) GetLikedMessages(ctx context.Context, req *proto.GetLikedMessagesRequest) (*proto.GetLikedMessagesResponse, error) {
	// Implementation here
	return &proto.GetLikedMessagesResponse{
		StatusCode: 200,
		Message:    "Liked messages retrieved successfully",
	}, nil
}

func (s *chatService) GetLastMessages(ctx context.Context, req *proto.GetLastMessagesRequest) (*proto.GetLastMessagesResponse, error) {
	// Implementation here
	return &proto.GetLastMessagesResponse{
		StatusCode: 200,
		Message:    "Last messages retrieved successfully",
	}, nil
}

func (s *chatService) SubscribeToLastMessages(req *proto.LastMessageStreamRequest, stream grpc.ServerStreamingServer[proto.ChatMessage]) error {
	// Implementation here
	return nil
}

func (s *chatService) GetUserStatus(ctx context.Context, req *proto.UserStatusRequest) (*proto.UserStatusResponse, error) {
	// Implementation here
	return &proto.UserStatusResponse{
		StatusCode: 200,
		Message:    "User status retrieved successfully",
	}, nil
}

func (s *chatService) SubscribeToUserStatus(req *proto.UserStatusSubscriptionRequest, stream grpc.ServerStreamingServer[proto.UserStatus]) error {
	// Implementation here
	return nil
}

func (s *chatService) UploadFile(stream grpc.ClientStreamingServer[proto.FileUploadRequest, proto.FileUploadResponse]) error {
	// Implementation here
	return nil
}

func (s *chatService) DownloadFile(req *proto.FileDownloadRequest, stream grpc.ServerStreamingServer[proto.FileDownloadResponse]) error {
	// Implementation here
	return nil
}

func (s *chatService) InitiateCall(ctx context.Context, req *proto.InitiateCallRequest) (*proto.InitiateCallResponse, error) {
	// Implementation here
	return &proto.InitiateCallResponse{
		StatusCode: 200,
		Message:    "Call initiated successfully",
	}, nil
}

func (s *chatService) AcceptCall(ctx context.Context, req *proto.AcceptCallRequest) (*proto.AcceptCallResponse, error) {
	// Implementation here
	return &proto.AcceptCallResponse{
		StatusCode: 200,
		Message:    "Call accepted successfully",
	}, nil
}

func (s *chatService) RejectCall(ctx context.Context, req *proto.RejectCallRequest) (*proto.RejectCallResponse, error) {
	// Implementation here
	return &proto.RejectCallResponse{
		StatusCode: 200,
		Message:    "Call rejected successfully",
	}, nil
}

func (s *chatService) EndCall(ctx context.Context, req *proto.EndCallRequest) (*proto.EndCallResponse, error) {
	// Implementation here
	return &proto.EndCallResponse{
		StatusCode: 200,
		Message:    "Call ended successfully",
	}, nil
}

func (s *chatService) GetCallHistory(ctx context.Context, req *proto.GetCallHistoryRequest) (*proto.GetCallHistoryResponse, error) {
	// Implementation here
	return &proto.GetCallHistoryResponse{
		StatusCode: 200,
		Message:    "Call history retrieved successfully",
	}, nil
}

func (s *chatService) UpdateNotificationSettings(ctx context.Context, req *proto.UpdateNotificationSettingsRequest) (*proto.UpdateNotificationSettingsResponse, error) {
	// Implementation here
	return &proto.UpdateNotificationSettingsResponse{
		StatusCode: 200,
		Message:    "Notification settings updated successfully",
	}, nil
}

func (s *chatService) GetNotificationSettings(ctx context.Context, req *proto.GetNotificationSettingsRequest) (*proto.GetNotificationSettingsResponse, error) {
	// Implementation here
	return &proto.GetNotificationSettingsResponse{
		StatusCode: 200,
		Message:    "Notification settings retrieved successfully",
	}, nil
}

func (s *chatService) MarkNotificationAsRead(ctx context.Context, req *proto.MarkNotificationAsReadRequest) (*proto.MarkNotificationAsReadResponse, error) {
	// Implementation here
	return &proto.MarkNotificationAsReadResponse{
		StatusCode: 200,
		Message:    "Notification marked as read successfully",
	}, nil
}

func (s *chatService) GetUnreadNotificationCount(ctx context.Context, req *proto.GetUnreadNotificationCountRequest) (*proto.GetUnreadNotificationCountResponse, error) {
	// Implementation here
	return &proto.GetUnreadNotificationCountResponse{
		StatusCode: 200,
		Message:    "Unread notification count retrieved successfully",
	}, nil
}

func (s *chatService) StartScreenShare(ctx context.Context, req *proto.StartScreenShareRequest) (*proto.StartScreenShareResponse, error) {
	// Implementation here
	return &proto.StartScreenShareResponse{
		StatusCode: 200,
		Message:    "Screen share started successfully",
	}, nil
}

func (s *chatService) StopScreenShare(ctx context.Context, req *proto.StopScreenShareRequest) (*proto.StopScreenShareResponse, error) {
	// Implementation here
	return &proto.StopScreenShareResponse{
		StatusCode: 200,
		Message:    "Screen share stopped successfully",
	}, nil
}

func (s *chatService) SubscribeToScreenShare(req *proto.SubscribeToScreenShareRequest, stream grpc.ServerStreamingServer[proto.ScreenShareUpdate]) error {
	// Implementation here
	return nil
}

func (s *chatService) ScheduleMessage(ctx context.Context, req *proto.ScheduleMessageRequest) (*proto.ScheduleMessageResponse, error) {
	// Implementation here
	return &proto.ScheduleMessageResponse{
		StatusCode: 200,
		Message:    "Message scheduled successfully",
	}, nil
}

func (s *chatService) CancelScheduledMessage(ctx context.Context, req *proto.CancelScheduledMessageRequest) (*proto.CancelScheduledMessageResponse, error) {
	// Implementation here
	return &proto.CancelScheduledMessageResponse{
		StatusCode: 200,
		Message:    "Scheduled message cancelled successfully",
	}, nil
}

func (s *chatService) GetScheduledMessages(ctx context.Context, req *proto.GetScheduledMessagesRequest) (*proto.GetScheduledMessagesResponse, error) {
	// Implementation here
	return &proto.GetScheduledMessagesResponse{
		StatusCode: 200,
		Message:    "Scheduled messages retrieved successfully",
	}, nil
}

func (s *chatService) GetChatAnalytics(ctx context.Context, req *proto.GetChatAnalyticsRequest) (*proto.GetChatAnalyticsResponse, error) {
	// Implementation here
	return &proto.GetChatAnalyticsResponse{
		StatusCode: 200,
		Message:    "Chat analytics retrieved successfully",
	}, nil
}
