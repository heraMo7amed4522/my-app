package interfaces

import (
	"context"

	"chat-services/internal/models"
	"chat-services/proto"
)

type ChatRepository interface {
	// Message operations
	SaveMessage(ctx context.Context, message *models.ChatMessage) (*models.ChatMessage, error)
	GetChatHistory(ctx context.Context, userID, peerID string, isGroup bool, limit, offset int32) ([]*models.ChatMessage, error)
	GetMessageByID(ctx context.Context, messageID string) (*models.ChatMessage, error)
	UpdateMessage(ctx context.Context, messageID string, content string) (*models.ChatMessage, error)
	DeleteMessage(ctx context.Context, messageID string) error
	SearchMessages(ctx context.Context, userID, query, peerID, groupID string, isGroup bool, messageType proto.MessageType, limit, offset int32) ([]*models.ChatMessage, error)

	// Message status operations
	MarkAsRead(ctx context.Context, userID, peerID string, isGroup bool) error
	UpdateMessageStatus(ctx context.Context, messageID string, status proto.MessageStatus) error

	// Reaction operations
	AddReaction(ctx context.Context, messageID, userID string, reactionType proto.ReactionType) error
	RemoveReaction(ctx context.Context, messageID, userID string, reactionType proto.ReactionType) error
	GetMessageReactions(ctx context.Context, messageID string) ([]*proto.MessageReaction, error)

	// User operations
	GetUsersByIDs(ctx context.Context, userIDs []string) ([]*proto.UserInfo, error)
	GetUsersInGroup(ctx context.Context, groupID string) ([]*proto.UserInfo, error)

	// Presence operations
	UpdateUserPresence(ctx context.Context, userID string, status proto.PresenceStatus, customMessage string) error
	GetUserPresence(ctx context.Context, userIDs []string) ([]*models.UserPresence, error)

	// Group operations
	CreateGroup(ctx context.Context, group *models.GroupInfo, memberIDs []string) (*models.GroupInfo, error)
	UpdateGroup(ctx context.Context, groupID string, updates map[string]interface{}) (*models.GroupInfo, error)
	AddGroupMembers(ctx context.Context, groupID string, memberIDs []string) error
	RemoveGroupMembers(ctx context.Context, groupID string, memberIDs []string) error

	// Thread operations
	CreateThread(ctx context.Context, parentMessageID, userID, content string) (*models.ChatMessage, error)
	GetThreadMessages(ctx context.Context, parentMessageID string, limit, offset int32) ([]*models.ChatMessage, error)

	// Pinning operations
	PinMessage(ctx context.Context, messageID, userID, chatID string, isGroup bool) error
	UnpinMessage(ctx context.Context, messageID, userID, chatID string, isGroup bool) error
	GetPinnedMessages(ctx context.Context, chatID string, isGroup bool) ([]*models.ChatMessage, error)
}