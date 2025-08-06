package interfaces

import (
	"context"

	"chat-services/proto"

	"google.golang.org/grpc"
)

type ChatService interface {
	ChatStream(stream grpc.BidiStreamingServer[proto.ChatStreamEnvelope, proto.ChatStreamEnvelope]) error

	SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error)
	GetChatHistory(ctx context.Context, req *proto.GetChatHistoryRequest) (*proto.GetChatHistoryResponse, error)

	EditMessage(ctx context.Context, req *proto.EditMessageRequest) (*proto.EditMessageResponse, error)
	DeleteMessage(ctx context.Context, req *proto.DeleteMessageRequest) (*proto.DeleteMessageResponse, error)
	SearchMessages(ctx context.Context, req *proto.SearchMessagesRequest) (*proto.SearchMessagesResponse, error)
	ForwardMessage(ctx context.Context, req *proto.ForwardMessageRequest) (*proto.ForwardMessageResponse, error)

	MarkAsRead(ctx context.Context, req *proto.ReadReceiptRequest) (*proto.ReadReceiptResponse, error)
	SendDeliveryReceipt(ctx context.Context, req *proto.DeliveryReceiptRequest) (*proto.DeliveryReceiptResponse, error)
	LikeMessage(ctx context.Context, req *proto.LikeMessageRequest) (*proto.LikeMessageResponse, error)
	GetLikedMessages(ctx context.Context, req *proto.GetLikedMessagesRequest) (*proto.GetLikedMessagesResponse, error)
	GetLastMessages(ctx context.Context, req *proto.GetLastMessagesRequest) (*proto.GetLastMessagesResponse, error)
	SubscribeToLastMessages(req *proto.LastMessageStreamRequest, stream grpc.ServerStreamingServer[proto.ChatMessage]) error

	AddReaction(ctx context.Context, req *proto.AddReactionRequest) (*proto.AddReactionResponse, error)
	RemoveReaction(ctx context.Context, req *proto.RemoveReactionRequest) (*proto.RemoveReactionResponse, error)
	GetMessageReactions(ctx context.Context, req *proto.GetMessageReactionsRequest) (*proto.GetMessageReactionsResponse, error)

	GetUsersByUserID(ctx context.Context, req *proto.GetUsersByUserIDRequest) (*proto.GetUsersByUserIDResponse, error)
	GetUsersInGroup(ctx context.Context, req *proto.GetUsersInGroupRequest) (*proto.GetUsersInGroupResponse, error)
	GetUserStatus(ctx context.Context, req *proto.UserStatusRequest) (*proto.UserStatusResponse, error)
	SubscribeToUserStatus(req *proto.UserStatusSubscriptionRequest, stream grpc.ServerStreamingServer[proto.UserStatus]) error

	UpdatePresenceStatus(ctx context.Context, req *proto.UpdatePresenceStatusRequest) (*proto.UpdatePresenceStatusResponse, error)
	GetUserPresence(ctx context.Context, req *proto.GetUserPresenceRequest) (*proto.GetUserPresenceResponse, error)
	SubscribeToPresence(req *proto.SubscribeToPresenceRequest, stream grpc.ServerStreamingServer[proto.PresenceUpdate]) error

	CreateGroup(ctx context.Context, req *proto.CreateGroupRequest) (*proto.CreateGroupResponse, error)
	JoinGroup(ctx context.Context, req *proto.JoinGroupRequest) (*proto.JoinGroupResponse, error)
	LeaveGroup(ctx context.Context, req *proto.LeaveGroupRequest) (*proto.LeaveGroupResponse, error)
	UpdateGroup(ctx context.Context, req *proto.UpdateGroupRequest) (*proto.UpdateGroupResponse, error)

	SendTypingIndicator(ctx context.Context, req *proto.TypingIndicatorRequest) (*proto.TypingIndicatorResponse, error)
	SubscribeToTypingIndicators(req *proto.TypingSubscriptionRequest, stream grpc.ServerStreamingServer[proto.TypingIndicator]) error

	CreateThread(ctx context.Context, req *proto.CreateThreadRequest) (*proto.CreateThreadResponse, error)
	GetThreadMessages(ctx context.Context, req *proto.GetThreadMessagesRequest) (*proto.GetThreadMessagesResponse, error)
	SubscribeToThread(req *proto.SubscribeToThreadRequest, stream grpc.ServerStreamingServer[proto.ThreadUpdate]) error

	PinMessage(ctx context.Context, req *proto.PinMessageRequest) (*proto.PinMessageResponse, error)
	UnpinMessage(ctx context.Context, req *proto.UnpinMessageRequest) (*proto.UnpinMessageResponse, error)
	GetPinnedMessages(ctx context.Context, req *proto.GetPinnedMessagesRequest) (*proto.GetPinnedMessagesResponse, error)

	UploadFile(stream grpc.ClientStreamingServer[proto.FileUploadRequest, proto.FileUploadResponse]) error
	DownloadFile(req *proto.FileDownloadRequest, stream grpc.ServerStreamingServer[proto.FileDownloadResponse]) error

	InitiateCall(ctx context.Context, req *proto.InitiateCallRequest) (*proto.InitiateCallResponse, error)
	AcceptCall(ctx context.Context, req *proto.AcceptCallRequest) (*proto.AcceptCallResponse, error)
	RejectCall(ctx context.Context, req *proto.RejectCallRequest) (*proto.RejectCallResponse, error)
	EndCall(ctx context.Context, req *proto.EndCallRequest) (*proto.EndCallResponse, error)
	GetCallHistory(ctx context.Context, req *proto.GetCallHistoryRequest) (*proto.GetCallHistoryResponse, error)

	UpdateNotificationSettings(ctx context.Context, req *proto.UpdateNotificationSettingsRequest) (*proto.UpdateNotificationSettingsResponse, error)
	GetNotificationSettings(ctx context.Context, req *proto.GetNotificationSettingsRequest) (*proto.GetNotificationSettingsResponse, error)
	MarkNotificationAsRead(ctx context.Context, req *proto.MarkNotificationAsReadRequest) (*proto.MarkNotificationAsReadResponse, error)
	GetUnreadNotificationCount(ctx context.Context, req *proto.GetUnreadNotificationCountRequest) (*proto.GetUnreadNotificationCountResponse, error)

	StartScreenShare(ctx context.Context, req *proto.StartScreenShareRequest) (*proto.StartScreenShareResponse, error)
	StopScreenShare(ctx context.Context, req *proto.StopScreenShareRequest) (*proto.StopScreenShareResponse, error)
	SubscribeToScreenShare(req *proto.SubscribeToScreenShareRequest, stream grpc.ServerStreamingServer[proto.ScreenShareUpdate]) error

	ScheduleMessage(ctx context.Context, req *proto.ScheduleMessageRequest) (*proto.ScheduleMessageResponse, error)
	CancelScheduledMessage(ctx context.Context, req *proto.CancelScheduledMessageRequest) (*proto.CancelScheduledMessageResponse, error)
	GetScheduledMessages(ctx context.Context, req *proto.GetScheduledMessagesRequest) (*proto.GetScheduledMessagesResponse, error)

	GetChatAnalytics(ctx context.Context, req *proto.GetChatAnalyticsRequest) (*proto.GetChatAnalyticsResponse, error)

	SubscribeToMessageUpdates(req *proto.SubscribeToMessageUpdatesRequest, stream grpc.ServerStreamingServer[proto.MessageUpdate]) error
	SubscribeToNotifications(req *proto.SubscribeToNotificationsRequest, stream grpc.ServerStreamingServer[proto.NotificationUpdate]) error
	SubscribeToChatEvents(req *proto.SubscribeToChatEventsRequest, stream grpc.ServerStreamingServer[proto.ChatEvent]) error
}
