package server

import (
	"chat-services/internal/service/interfaces"
	"chat-services/proto"
	"context"

	"google.golang.org/grpc"
)

type ChatServer struct {
	proto.UnimplementedChatServiceServer
	chatService interfaces.ChatService
}

func NewChatServer(chatService interfaces.ChatService) *ChatServer {
	return &ChatServer{
		chatService: chatService,
	}
}

// ChatStream handles bidirectional streaming for real-time chat
func (s *ChatServer) ChatStream(stream grpc.BidiStreamingServer[proto.ChatStreamEnvelope, proto.ChatStreamEnvelope]) error {
	return s.chatService.ChatStream(stream)
}

// Message operations
func (s *ChatServer) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	return s.chatService.SendMessage(ctx, req)
}

func (s *ChatServer) GetChatHistory(ctx context.Context, req *proto.GetChatHistoryRequest) (*proto.GetChatHistoryResponse, error) {
	return s.chatService.GetChatHistory(ctx, req)
}

func (s *ChatServer) EditMessage(ctx context.Context, req *proto.EditMessageRequest) (*proto.EditMessageResponse, error) {
	return s.chatService.EditMessage(ctx, req)
}

func (s *ChatServer) DeleteMessage(ctx context.Context, req *proto.DeleteMessageRequest) (*proto.DeleteMessageResponse, error) {
	return s.chatService.DeleteMessage(ctx, req)
}

func (s *ChatServer) SearchMessages(ctx context.Context, req *proto.SearchMessagesRequest) (*proto.SearchMessagesResponse, error) {
	return s.chatService.SearchMessages(ctx, req)
}

func (s *ChatServer) ForwardMessage(ctx context.Context, req *proto.ForwardMessageRequest) (*proto.ForwardMessageResponse, error) {
	return s.chatService.ForwardMessage(ctx, req)
}

// Message status and interactions
func (s *ChatServer) MarkAsRead(ctx context.Context, req *proto.ReadReceiptRequest) (*proto.ReadReceiptResponse, error) {
	return s.chatService.MarkAsRead(ctx, req)
}

func (s *ChatServer) SendDeliveryReceipt(ctx context.Context, req *proto.DeliveryReceiptRequest) (*proto.DeliveryReceiptResponse, error) {
	return s.chatService.SendDeliveryReceipt(ctx, req)
}

func (s *ChatServer) LikeMessage(ctx context.Context, req *proto.LikeMessageRequest) (*proto.LikeMessageResponse, error) {
	return s.chatService.LikeMessage(ctx, req)
}

func (s *ChatServer) GetLikedMessages(ctx context.Context, req *proto.GetLikedMessagesRequest) (*proto.GetLikedMessagesResponse, error) {
	return s.chatService.GetLikedMessages(ctx, req)
}

func (s *ChatServer) GetLastMessages(ctx context.Context, req *proto.GetLastMessagesRequest) (*proto.GetLastMessagesResponse, error) {
	return s.chatService.GetLastMessages(ctx, req)
}

func (s *ChatServer) SubscribeToLastMessages(req *proto.LastMessageStreamRequest, stream grpc.ServerStreamingServer[proto.ChatMessage]) error {
	return s.chatService.SubscribeToLastMessages(req, stream)
}

// Reactions
func (s *ChatServer) AddReaction(ctx context.Context, req *proto.AddReactionRequest) (*proto.AddReactionResponse, error) {
	return s.chatService.AddReaction(ctx, req)
}

func (s *ChatServer) RemoveReaction(ctx context.Context, req *proto.RemoveReactionRequest) (*proto.RemoveReactionResponse, error) {
	return s.chatService.RemoveReaction(ctx, req)
}

func (s *ChatServer) GetMessageReactions(ctx context.Context, req *proto.GetMessageReactionsRequest) (*proto.GetMessageReactionsResponse, error) {
	return s.chatService.GetMessageReactions(ctx, req)
}

// User operations
func (s *ChatServer) GetUsersByUserID(ctx context.Context, req *proto.GetUsersByUserIDRequest) (*proto.GetUsersByUserIDResponse, error) {
	return s.chatService.GetUsersByUserID(ctx, req)
}

func (s *ChatServer) GetUsersInGroup(ctx context.Context, req *proto.GetUsersInGroupRequest) (*proto.GetUsersInGroupResponse, error) {
	return s.chatService.GetUsersInGroup(ctx, req)
}

func (s *ChatServer) GetUserStatus(ctx context.Context, req *proto.UserStatusRequest) (*proto.UserStatusResponse, error) {
	return s.chatService.GetUserStatus(ctx, req)
}

func (s *ChatServer) SubscribeToUserStatus(req *proto.UserStatusSubscriptionRequest, stream grpc.ServerStreamingServer[proto.UserStatus]) error {
	return s.chatService.SubscribeToUserStatus(req, stream)
}

// Presence
func (s *ChatServer) UpdatePresenceStatus(ctx context.Context, req *proto.UpdatePresenceStatusRequest) (*proto.UpdatePresenceStatusResponse, error) {
	return s.chatService.UpdatePresenceStatus(ctx, req)
}

func (s *ChatServer) GetUserPresence(ctx context.Context, req *proto.GetUserPresenceRequest) (*proto.GetUserPresenceResponse, error) {
	return s.chatService.GetUserPresence(ctx, req)
}

func (s *ChatServer) SubscribeToPresence(req *proto.SubscribeToPresenceRequest, stream grpc.ServerStreamingServer[proto.PresenceUpdate]) error {
	return s.chatService.SubscribeToPresence(req, stream)
}

// Groups
func (s *ChatServer) CreateGroup(ctx context.Context, req *proto.CreateGroupRequest) (*proto.CreateGroupResponse, error) {
	return s.chatService.CreateGroup(ctx, req)
}

func (s *ChatServer) JoinGroup(ctx context.Context, req *proto.JoinGroupRequest) (*proto.JoinGroupResponse, error) {
	return s.chatService.JoinGroup(ctx, req)
}

func (s *ChatServer) LeaveGroup(ctx context.Context, req *proto.LeaveGroupRequest) (*proto.LeaveGroupResponse, error) {
	return s.chatService.LeaveGroup(ctx, req)
}

func (s *ChatServer) UpdateGroup(ctx context.Context, req *proto.UpdateGroupRequest) (*proto.UpdateGroupResponse, error) {
	return s.chatService.UpdateGroup(ctx, req)
}

// Typing indicators
func (s *ChatServer) SendTypingIndicator(ctx context.Context, req *proto.TypingIndicatorRequest) (*proto.TypingIndicatorResponse, error) {
	return s.chatService.SendTypingIndicator(ctx, req)
}

func (s *ChatServer) SubscribeToTypingIndicators(req *proto.TypingSubscriptionRequest, stream grpc.ServerStreamingServer[proto.TypingIndicator]) error {
	return s.chatService.SubscribeToTypingIndicators(req, stream)
}

// Threading
func (s *ChatServer) CreateThread(ctx context.Context, req *proto.CreateThreadRequest) (*proto.CreateThreadResponse, error) {
	return s.chatService.CreateThread(ctx, req)
}

func (s *ChatServer) GetThreadMessages(ctx context.Context, req *proto.GetThreadMessagesRequest) (*proto.GetThreadMessagesResponse, error) {
	return s.chatService.GetThreadMessages(ctx, req)
}

func (s *ChatServer) SubscribeToThread(req *proto.SubscribeToThreadRequest, stream grpc.ServerStreamingServer[proto.ThreadUpdate]) error {
	return s.chatService.SubscribeToThread(req, stream)
}

// Message management
func (s *ChatServer) PinMessage(ctx context.Context, req *proto.PinMessageRequest) (*proto.PinMessageResponse, error) {
	return s.chatService.PinMessage(ctx, req)
}

func (s *ChatServer) UnpinMessage(ctx context.Context, req *proto.UnpinMessageRequest) (*proto.UnpinMessageResponse, error) {
	return s.chatService.UnpinMessage(ctx, req)
}

func (s *ChatServer) GetPinnedMessages(ctx context.Context, req *proto.GetPinnedMessagesRequest) (*proto.GetPinnedMessagesResponse, error) {
	return s.chatService.GetPinnedMessages(ctx, req)
}

// File operations
func (s *ChatServer) UploadFile(stream grpc.ClientStreamingServer[proto.FileUploadRequest, proto.FileUploadResponse]) error {
	return s.chatService.UploadFile(stream)
}

func (s *ChatServer) DownloadFile(req *proto.FileDownloadRequest, stream grpc.ServerStreamingServer[proto.FileDownloadResponse]) error {
	return s.chatService.DownloadFile(req, stream)
}

// Call operations
func (s *ChatServer) InitiateCall(ctx context.Context, req *proto.InitiateCallRequest) (*proto.InitiateCallResponse, error) {
	return s.chatService.InitiateCall(ctx, req)
}

func (s *ChatServer) AcceptCall(ctx context.Context, req *proto.AcceptCallRequest) (*proto.AcceptCallResponse, error) {
	return s.chatService.AcceptCall(ctx, req)
}

func (s *ChatServer) RejectCall(ctx context.Context, req *proto.RejectCallRequest) (*proto.RejectCallResponse, error) {
	return s.chatService.RejectCall(ctx, req)
}

func (s *ChatServer) EndCall(ctx context.Context, req *proto.EndCallRequest) (*proto.EndCallResponse, error) {
	return s.chatService.EndCall(ctx, req)
}

func (s *ChatServer) GetCallHistory(ctx context.Context, req *proto.GetCallHistoryRequest) (*proto.GetCallHistoryResponse, error) {
	return s.chatService.GetCallHistory(ctx, req)
}

// Notification operations
func (s *ChatServer) UpdateNotificationSettings(ctx context.Context, req *proto.UpdateNotificationSettingsRequest) (*proto.UpdateNotificationSettingsResponse, error) {
	return s.chatService.UpdateNotificationSettings(ctx, req)
}

func (s *ChatServer) GetNotificationSettings(ctx context.Context, req *proto.GetNotificationSettingsRequest) (*proto.GetNotificationSettingsResponse, error) {
	return s.chatService.GetNotificationSettings(ctx, req)
}

func (s *ChatServer) MarkNotificationAsRead(ctx context.Context, req *proto.MarkNotificationAsReadRequest) (*proto.MarkNotificationAsReadResponse, error) {
	return s.chatService.MarkNotificationAsRead(ctx, req)
}

func (s *ChatServer) GetUnreadNotificationCount(ctx context.Context, req *proto.GetUnreadNotificationCountRequest) (*proto.GetUnreadNotificationCountResponse, error) {
	return s.chatService.GetUnreadNotificationCount(ctx, req)
}

// Screen sharing
func (s *ChatServer) StartScreenShare(ctx context.Context, req *proto.StartScreenShareRequest) (*proto.StartScreenShareResponse, error) {
	return s.chatService.StartScreenShare(ctx, req)
}

func (s *ChatServer) StopScreenShare(ctx context.Context, req *proto.StopScreenShareRequest) (*proto.StopScreenShareResponse, error) {
	return s.chatService.StopScreenShare(ctx, req)
}

func (s *ChatServer) SubscribeToScreenShare(req *proto.SubscribeToScreenShareRequest, stream grpc.ServerStreamingServer[proto.ScreenShareUpdate]) error {
	return s.chatService.SubscribeToScreenShare(req, stream)
}

// Scheduled messages
func (s *ChatServer) ScheduleMessage(ctx context.Context, req *proto.ScheduleMessageRequest) (*proto.ScheduleMessageResponse, error) {
	return s.chatService.ScheduleMessage(ctx, req)
}

func (s *ChatServer) CancelScheduledMessage(ctx context.Context, req *proto.CancelScheduledMessageRequest) (*proto.CancelScheduledMessageResponse, error) {
	return s.chatService.CancelScheduledMessage(ctx, req)
}

func (s *ChatServer) GetScheduledMessages(ctx context.Context, req *proto.GetScheduledMessagesRequest) (*proto.GetScheduledMessagesResponse, error) {
	return s.chatService.GetScheduledMessages(ctx, req)
}

// Analytics
func (s *ChatServer) GetChatAnalytics(ctx context.Context, req *proto.GetChatAnalyticsRequest) (*proto.GetChatAnalyticsResponse, error) {
	return s.chatService.GetChatAnalytics(ctx, req)
}

// Subscriptions
func (s *ChatServer) SubscribeToMessageUpdates(req *proto.SubscribeToMessageUpdatesRequest, stream grpc.ServerStreamingServer[proto.MessageUpdate]) error {
	return s.chatService.SubscribeToMessageUpdates(req, stream)
}

func (s *ChatServer) SubscribeToNotifications(req *proto.SubscribeToNotificationsRequest, stream grpc.ServerStreamingServer[proto.NotificationUpdate]) error {
	return s.chatService.SubscribeToNotifications(req, stream)
}

func (s *ChatServer) SubscribeToChatEvents(req *proto.SubscribeToChatEventsRequest, stream grpc.ServerStreamingServer[proto.ChatEvent]) error {
	return s.chatService.SubscribeToChatEvents(req, stream)
}
