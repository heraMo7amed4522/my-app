// interfaces/chat_repository.go
package interfaces

import (
	pb "chats-services/proto"
	"context"
)

type ChatRepository interface {

	SaveMessage(ctx context.Context, message *pb.ChatMessage) (*pb.ChatMessage, error)

	GetChatHistory(ctx context.Context, userID, peerID string, isGroup bool, limit, offset int32) ([]*pb.ChatMessage, error)

	GetMessageByID(ctx context.Context, messageID string) (*pb.ChatMessage, error)
	UpdateMessage(ctx context.Context, messageID, content string) (*pb.ChatMessage, error)

	DeleteMessage(ctx context.Context, messageID string) error

	SearchMessages(ctx context.Context, userID, query, peerID, groupID string, isGroup bool, messageType pb.MessageType, limit, offset int32) ([]*pb.ChatMessage, error)

	MarkAsRead(ctx context.Context, userID, peerID string, isGroup bool) error

	UpdateMessageStatus(ctx context.Context, messageID string, status pb.MessageStatus) error

	AddReaction(ctx context.Context, messageID, userID string, reactionType pb.ReactionType) error
	RemoveReaction(ctx context.Context, messageID, userID string, reactionType pb.ReactionType) error
	GetMessageReactions(ctx context.Context, messageID string) ([]*pb.MessageReaction, error)


	GetUsersByIDs(ctx context.Context, userIDs []string) ([]*pb.UserInfo, error)
	GetUsersInGroup(ctx context.Context, groupID string) ([]*pb.UserInfo, error)


	UpdateUserPresence(ctx context.Context, userID string, status pb.PresenceStatus, customMessage string) error
	GetUserPresence(ctx context.Context, userIDs []string) ([]*pb.UserPresence, error)


	CreateGroup(ctx context.Context, group *pb.GroupInfo) (*pb.GroupInfo, error)
	UpdateGroup(ctx context.Context, groupID string, updates map[string]interface{}) (*pb.GroupInfo, error)
	AddGroupMembers(ctx context.Context, groupID string, memberIDs []string) error
	RemoveGroupMembers(ctx context.Context, groupID string, memberIDs []string) error


	CreateThread(ctx context.Context, parentMessageID, userID, content string) (*pb.ChatMessage, error)
	GetThreadMessages(ctx context.Context, parentMessageID string, limit, offset int32) ([]*pb.ChatMessage, error)


	PinMessage(ctx context.Context, messageID, userID, chatID string, isGroup bool) error
	UnpinMessage(ctx context.Context, messageID, userID, chatID string, isGroup bool) error
	GetPinnedMessages(ctx context.Context, chatID string, isGroup bool) ([]*pb.ChatMessage, error)
}