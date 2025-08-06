package models

import (
	"chats-services/proto"
	"time"
)

type ChatMessage struct {
	MessageID         string              `json:"message_id" db:"message_id"`
	SenderID          string              `json:"sender_id" db:"sender_id"`
	ReceiverID        string              `json:"receiver_id" db:"receiver_id"`
	GroupID           string              `json:"group_id" db:"group_id"`
	Content           string              `json:"content" db:"content"`
	Type              proto.MessageType   `json:"type" db:"type"`
	Timestamp         time.Time           `json:"timestamp" db:"timestamp"`
	IsGroup           bool                `json:"is_group" db:"is_group"`
	IsRead            bool                `json:"is_read" db:"is_read"`
	IsEdited          bool                `json:"is_edited" db:"is_edited"`
	LikedBy           []string            `json:"liked_by" db:"liked_by"`
	ReplyToMessageID  string              `json:"reply_to_message_id" db:"reply_to_message_id"`
	Attachments       []string            `json:"attachments" db:"attachments"`
	Status            proto.MessageStatus `json:"status" db:"status"`
	DeliveredAt       *time.Time          `json:"delivered_at" db:"delivered_at"`
	ReadAt            *time.Time          `json:"read_at" db:"read_at"`
	Reactions         []MessageReaction   `json:"reactions" db:"reactions"`
	IsPinned          bool                `json:"is_pinned" db:"is_pinned"`
	PinnedAt          *time.Time          `json:"pinned_at" db:"pinned_at"`
	PinnedBy          string              `json:"pinned_by" db:"pinned_by"`
	ForwardCount      int32               `json:"forward_count" db:"forward_count"`
	OriginalMessageID string              `json:"original_message_id" db:"original_message_id"`
	FileMetadata      *FileMetadata       `json:"file_metadata" db:"file_metadata"`
	LocationData      *LocationData       `json:"location_data" db:"location_data"`
	PollData          *PollData           `json:"poll_data" db:"poll_data"`

	// Threading support
	ThreadID         string `json:"thread_id" db:"thread_id"`
	ParentMessageID  string `json:"parent_message_id" db:"parent_message_id"`
	ThreadReplyCount int32  `json:"thread_reply_count" db:"thread_reply_count"`

	// Message editing
	EditedAt    *time.Time `json:"edited_at" db:"edited_at"`
	EditHistory []string   `json:"edit_history" db:"edit_history"`

	// Message scheduling
	IsScheduled bool       `json:"is_scheduled" db:"is_scheduled"`
	ScheduledAt *time.Time `json:"scheduled_at" db:"scheduled_at"`

	// Real-time collaboration
	MentionedUserIDs []string `json:"mentioned_user_ids" db:"mentioned_user_ids"`
	IsSystemMessage  bool     `json:"is_system_message" db:"is_system_message"`

	// Advanced metadata
	DeviceInfo      string `json:"device_info" db:"device_info"`
	ClientVersion   string `json:"client_version" db:"client_version"`
	IsEncrypted     bool   `json:"is_encrypted" db:"is_encrypted"`
	EncryptionKeyID string `json:"encryption_key_id" db:"encryption_key_id"`

	// Database timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type MessageReaction struct {
	UserID       string             `json:"user_id" db:"user_id"`
	ReactionType proto.ReactionType `json:"reaction_type" db:"reaction_type"`
	Timestamp    time.Time          `json:"timestamp" db:"timestamp"`
}
type FileMetadata struct {
	FileName     string `json:"file_name" db:"file_name"`
	FileType     string `json:"file_type" db:"file_type"`
	FileSize     int64  `json:"file_size" db:"file_size"`
	FileURL      string `json:"file_url" db:"file_url"`
	ThumbnailURL string `json:"thumbnail_url" db:"thumbnail_url"`
	Duration     int32  `json:"duration" db:"duration"` // for audio/video files
	Width        int32  `json:"width" db:"width"`       // for images/videos
	Height       int32  `json:"height" db:"height"`     // for images/videos
}

type LocationData struct {
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
	Address   string  `json:"address" db:"address"`
	PlaceName string  `json:"place_name" db:"place_name"`
}
type PollOption struct {
	ID        string   `json:"id" db:"id"`
	Text      string   `json:"text" db:"text"`
	VoterIDs  []string `json:"voter_ids" db:"voter_ids"`
	VoteCount int32    `json:"vote_count" db:"vote_count"`
}
type PollData struct {
	Question             string       `json:"question" db:"question"`
	Options              []PollOption `json:"options" db:"options"`
	AllowMultipleAnswers bool         `json:"allow_multiple_answers" db:"allow_multiple_answers"`
	ExpiresAt            *time.Time   `json:"expires_at" db:"expires_at"`
	IsAnonymous          bool         `json:"is_anonymous" db:"is_anonymous"`
	TotalVotes           int32        `json:"total_votes" db:"total_votes"`
}
type UserPresence struct {
	UserID        string               `json:"user_id" db:"user_id"`
	Status        proto.PresenceStatus `json:"status" db:"status"`
	CustomMessage string               `json:"custom_message" db:"custom_message"`
	LastSeen      time.Time            `json:"last_seen" db:"last_seen"`
	UpdatedAt     time.Time            `json:"updated_at" db:"updated_at"`
}
type GroupMember struct {
	UserID   string          `json:"user_id" db:"user_id"`
	Role     proto.GroupRole `json:"role" db:"role"`
	JoinedAt time.Time       `json:"joined_at" db:"joined_at"`
}
type CallInfo struct {
	CallID       string           `json:"call_id" db:"call_id"`
	CallerID     string           `json:"caller_id" db:"caller_id"`
	ReceiverID   string           `json:"receiver_id" db:"receiver_id"`
	GroupID      string           `json:"group_id" db:"group_id"`
	CallType     proto.CallType   `json:"call_type" db:"call_type"`
	Status       proto.CallStatus `json:"status" db:"status"`
	StartTime    *time.Time       `json:"start_time" db:"start_time"`
	EndTime      *time.Time       `json:"end_time" db:"end_time"`
	Duration     int32            `json:"duration" db:"duration"`
	IsGroup      bool             `json:"is_group" db:"is_group"`
	Participants []string         `json:"participants" db:"participants"`
}
type NotificationSettings struct {
	UserID                       string   `json:"user_id" db:"user_id"`
	EnablePushNotifications      bool     `json:"enable_push_notifications" db:"enable_push_notifications"`
	EnableSoundNotifications     bool     `json:"enable_sound_notifications" db:"enable_sound_notifications"`
	EnableVibrationNotifications bool     `json:"enable_vibration_notifications" db:"enable_vibration_notifications"`
	EnableEmailNotifications     bool     `json:"enable_email_notifications" db:"enable_email_notifications"`
	MuteAllChats                 bool     `json:"mute_all_chats" db:"mute_all_chats"`
	MutedChatIDs                 []string `json:"muted_chat_ids" db:"muted_chat_ids"`
	MutedGroupIDs                []string `json:"muted_group_ids" db:"muted_group_ids"`
	QuietHoursStart              string   `json:"quiet_hours_start" db:"quiet_hours_start"`
	QuietHoursEnd                string   `json:"quiet_hours_end" db:"quiet_hours_end"`
	ShowMessagePreview           bool     `json:"show_message_preview" db:"show_message_preview"`
}
