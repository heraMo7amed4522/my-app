package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"chat-services/internal/models"
	"chat-services/internal/repository/interfaces"
	"chat-services/proto"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) interfaces.ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) SaveMessage(ctx context.Context, message *models.ChatMessage) (*models.ChatMessage, error) {
	if message.MessageID == "" {
		message.MessageID = uuid.New().String()
	}
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()

	// Marshal JSON fields
	likedByJSON, _ := json.Marshal(message.LikedBy)
	attachmentsJSON, _ := json.Marshal(message.Attachments)
	reactionsJSON, _ := json.Marshal(message.Reactions)
	fileMetadataJSON, _ := json.Marshal(message.FileMetadata)
	locationDataJSON, _ := json.Marshal(message.LocationData)
	pollDataJSON, _ := json.Marshal(message.PollData)
	editHistoryJSON, _ := json.Marshal(message.EditHistory)
	mentionedUserIDsJSON, _ := json.Marshal(message.MentionedUserIDs)

	query := `
		INSERT INTO chat_messages (
			id, sender_id, receiver_id, group_id, content, message_type, timestamp,
			is_group, is_read, is_edited, liked_by, reply_to_message_id, attachments,
			status, delivered_at, read_at, reactions, is_pinned, pinned_at, pinned_by,
			forward_count, original_message_id, file_metadata, location_data, poll_data,
			thread_id, parent_message_id, thread_reply_count, edited_at, edit_history,
			is_scheduled, scheduled_at, mentioned_user_ids, is_system_message,
			device_info, client_version, is_encrypted, encryption_key_id,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37, $38, $39, $40
		)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		message.MessageID, message.SenderID, message.ReceiverID, message.GroupID,
		message.Content, message.Type, message.Timestamp, message.IsGroup,
		message.IsRead, message.IsEdited, likedByJSON, message.ReplyToMessageID,
		attachmentsJSON, message.Status, message.DeliveredAt, message.ReadAt,
		reactionsJSON, message.IsPinned, message.PinnedAt, message.PinnedBy,
		message.ForwardCount, message.OriginalMessageID, fileMetadataJSON,
		locationDataJSON, pollDataJSON, message.ThreadID, message.ParentMessageID,
		message.ThreadReplyCount, message.EditedAt, editHistoryJSON,
		message.IsScheduled, message.ScheduledAt, mentionedUserIDsJSON,
		message.IsSystemMessage, message.DeviceInfo, message.ClientVersion,
		message.IsEncrypted, message.EncryptionKeyID, message.CreatedAt, message.UpdatedAt,
	).Scan(&message.MessageID, &message.CreatedAt, &message.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	return message, nil
}

func (r *chatRepository) GetChatHistory(ctx context.Context, userID, peerID string, isGroup bool, limit, offset int32) ([]*models.ChatMessage, error) {
	var query string
	var args []interface{}

	if isGroup {
		query = `
			SELECT id, sender_id, receiver_id, group_id, content, message_type,
				   timestamp, is_group, is_read, is_edited, liked_by, reply_to_message_id,
				   attachments, status, delivered_at, read_at, reactions, is_pinned,
				   pinned_at, pinned_by, forward_count, original_message_id, file_metadata,
				   location_data, poll_data, thread_id, parent_message_id, thread_reply_count,
				   edited_at, edit_history, is_scheduled, scheduled_at, mentioned_user_ids,
				   is_system_message, device_info, client_version, is_encrypted, encryption_key_id,
				   created_at, updated_at
			FROM chat_messages
			WHERE group_id = $1 AND is_group = true
			ORDER BY timestamp DESC
			LIMIT $2 OFFSET $3
		`
		args = []interface{}{peerID, limit, offset}
	} else {
		query = `
			SELECT id, sender_id, receiver_id, group_id, content, message_type,
				   timestamp, is_group, is_read, is_edited, liked_by, reply_to_message_id,
				   attachments, status, delivered_at, read_at, reactions, is_pinned,
				   pinned_at, pinned_by, forward_count, original_message_id, file_metadata,
				   location_data, poll_data, thread_id, parent_message_id, thread_reply_count,
				   edited_at, edit_history, is_scheduled, scheduled_at, mentioned_user_ids,
				   is_system_message, device_info, client_version, is_encrypted, encryption_key_id,
				   created_at, updated_at
			FROM chat_messages
			WHERE ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1))
				AND is_group = false
			ORDER BY timestamp DESC
			LIMIT $3 OFFSET $4
		`
		args = []interface{}{userID, peerID, limit, offset}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat history: %w", err)
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		message := &models.ChatMessage{}
		var likedByJSON, attachmentsJSON, reactionsJSON, fileMetadataJSON, locationDataJSON, pollDataJSON, editHistoryJSON, mentionedUserIDsJSON []byte

		err := rows.Scan(
			&message.MessageID, &message.SenderID, &message.ReceiverID, &message.GroupID,
			&message.Content, &message.Type, &message.Timestamp, &message.IsGroup,
			&message.IsRead, &message.IsEdited, &likedByJSON, &message.ReplyToMessageID,
			&attachmentsJSON, &message.Status, &message.DeliveredAt, &message.ReadAt,
			&reactionsJSON, &message.IsPinned, &message.PinnedAt, &message.PinnedBy,
			&message.ForwardCount, &message.OriginalMessageID, &fileMetadataJSON,
			&locationDataJSON, &pollDataJSON, &message.ThreadID, &message.ParentMessageID,
			&message.ThreadReplyCount, &message.EditedAt, &editHistoryJSON, &message.IsScheduled,
			&message.ScheduledAt, &mentionedUserIDsJSON, &message.IsSystemMessage,
			&message.DeviceInfo, &message.ClientVersion, &message.IsEncrypted,
			&message.EncryptionKeyID, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Unmarshal JSON fields
		if len(likedByJSON) > 0 {
			json.Unmarshal(likedByJSON, &message.LikedBy)
		}
		if len(attachmentsJSON) > 0 {
			json.Unmarshal(attachmentsJSON, &message.Attachments)
		}
		if len(reactionsJSON) > 0 {
			json.Unmarshal(reactionsJSON, &message.Reactions)
		}
		if len(fileMetadataJSON) > 0 {
			json.Unmarshal(fileMetadataJSON, &message.FileMetadata)
		}
		if len(locationDataJSON) > 0 {
			json.Unmarshal(locationDataJSON, &message.LocationData)
		}
		if len(pollDataJSON) > 0 {
			json.Unmarshal(pollDataJSON, &message.PollData)
		}
		if len(editHistoryJSON) > 0 {
			json.Unmarshal(editHistoryJSON, &message.EditHistory)
		}
		if len(mentionedUserIDsJSON) > 0 {
			json.Unmarshal(mentionedUserIDsJSON, &message.MentionedUserIDs)
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (r *chatRepository) GetMessageByID(ctx context.Context, messageID string) (*models.ChatMessage, error) {
	query := `
		SELECT id, sender_id, receiver_id, group_id, content, message_type,
			   timestamp, is_group, is_read, is_edited, liked_by, reply_to_message_id,
			   attachments, status, delivered_at, read_at, reactions, is_pinned,
			   pinned_at, pinned_by, forward_count, original_message_id, file_metadata,
			   location_data, poll_data, thread_id, parent_message_id, thread_reply_count,
			   edited_at, edit_history, is_scheduled, scheduled_at, mentioned_user_ids,
			   is_system_message, device_info, client_version, is_encrypted, encryption_key_id,
			   created_at, updated_at
		FROM chat_messages
		WHERE id = $1
	`

	message := &models.ChatMessage{}
	var likedByJSON, attachmentsJSON, reactionsJSON, fileMetadataJSON, locationDataJSON, pollDataJSON, editHistoryJSON, mentionedUserIDsJSON []byte

	err := r.db.QueryRowContext(ctx, query, messageID).Scan(
		&message.MessageID, &message.SenderID, &message.ReceiverID, &message.GroupID,
		&message.Content, &message.Type, &message.Timestamp, &message.IsGroup,
		&message.IsRead, &message.IsEdited, &likedByJSON, &message.ReplyToMessageID,
		&attachmentsJSON, &message.Status, &message.DeliveredAt, &message.ReadAt,
		&reactionsJSON, &message.IsPinned, &message.PinnedAt, &message.PinnedBy,
		&message.ForwardCount, &message.OriginalMessageID, &fileMetadataJSON,
		&locationDataJSON, &pollDataJSON, &message.ThreadID, &message.ParentMessageID,
		&message.ThreadReplyCount, &message.EditedAt, &editHistoryJSON, &message.IsScheduled,
		&message.ScheduledAt, &mentionedUserIDsJSON, &message.IsSystemMessage,
		&message.DeviceInfo, &message.ClientVersion, &message.IsEncrypted,
		&message.EncryptionKeyID, &message.CreatedAt, &message.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	// Unmarshal JSON fields
	if len(likedByJSON) > 0 {
		json.Unmarshal(likedByJSON, &message.LikedBy)
	}
	if len(attachmentsJSON) > 0 {
		json.Unmarshal(attachmentsJSON, &message.Attachments)
	}
	if len(reactionsJSON) > 0 {
		json.Unmarshal(reactionsJSON, &message.Reactions)
	}
	if len(fileMetadataJSON) > 0 {
		json.Unmarshal(fileMetadataJSON, &message.FileMetadata)
	}
	if len(locationDataJSON) > 0 {
		json.Unmarshal(locationDataJSON, &message.LocationData)
	}
	if len(pollDataJSON) > 0 {
		json.Unmarshal(pollDataJSON, &message.PollData)
	}
	if len(editHistoryJSON) > 0 {
		json.Unmarshal(editHistoryJSON, &message.EditHistory)
	}
	if len(mentionedUserIDsJSON) > 0 {
		json.Unmarshal(mentionedUserIDsJSON, &message.MentionedUserIDs)
	}

	return message, nil
}

func (r *chatRepository) UpdateMessage(ctx context.Context, messageID string, content string) (*models.ChatMessage, error) {
	now := time.Now()
	query := `
		UPDATE chat_messages 
		SET content = $1, is_edited = true, edited_at = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, sender_id, receiver_id, group_id, content, message_type,
				  timestamp, is_group, is_read, is_edited, liked_by, reply_to_message_id,
				  attachments, status, delivered_at, read_at, reactions, is_pinned,
				  pinned_at, pinned_by, forward_count, original_message_id, file_metadata,
				  location_data, poll_data, thread_id, parent_message_id, thread_reply_count,
				  edited_at, edit_history, is_scheduled, scheduled_at, mentioned_user_ids,
				  is_system_message, device_info, client_version, is_encrypted, encryption_key_id,
				  created_at, updated_at
	`

	message := &models.ChatMessage{}
	var likedByJSON, attachmentsJSON, reactionsJSON, fileMetadataJSON, locationDataJSON, pollDataJSON, editHistoryJSON, mentionedUserIDsJSON []byte

	err := r.db.QueryRowContext(ctx, query, content, now, now, messageID).Scan(
		&message.MessageID, &message.SenderID, &message.ReceiverID, &message.GroupID,
		&message.Content, &message.Type, &message.Timestamp, &message.IsGroup,
		&message.IsRead, &message.IsEdited, &likedByJSON, &message.ReplyToMessageID,
		&attachmentsJSON, &message.Status, &message.DeliveredAt, &message.ReadAt,
		&reactionsJSON, &message.IsPinned, &message.PinnedAt, &message.PinnedBy,
		&message.ForwardCount, &message.OriginalMessageID, &fileMetadataJSON,
		&locationDataJSON, &pollDataJSON, &message.ThreadID, &message.ParentMessageID,
		&message.ThreadReplyCount, &message.EditedAt, &editHistoryJSON, &message.IsScheduled,
		&message.ScheduledAt, &mentionedUserIDsJSON, &message.IsSystemMessage,
		&message.DeviceInfo, &message.ClientVersion, &message.IsEncrypted,
		&message.EncryptionKeyID, &message.CreatedAt, &message.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to update message: %w", err)
	}

	// Unmarshal JSON fields
	if len(likedByJSON) > 0 {
		json.Unmarshal(likedByJSON, &message.LikedBy)
	}
	if len(attachmentsJSON) > 0 {
		json.Unmarshal(attachmentsJSON, &message.Attachments)
	}
	if len(reactionsJSON) > 0 {
		json.Unmarshal(reactionsJSON, &message.Reactions)
	}
	if len(fileMetadataJSON) > 0 {
		json.Unmarshal(fileMetadataJSON, &message.FileMetadata)
	}
	if len(locationDataJSON) > 0 {
		json.Unmarshal(locationDataJSON, &message.LocationData)
	}
	if len(pollDataJSON) > 0 {
		json.Unmarshal(pollDataJSON, &message.PollData)
	}
	if len(editHistoryJSON) > 0 {
		json.Unmarshal(editHistoryJSON, &message.EditHistory)
	}
	if len(mentionedUserIDsJSON) > 0 {
		json.Unmarshal(mentionedUserIDsJSON, &message.MentionedUserIDs)
	}

	return message, nil
}

func (r *chatRepository) SearchMessages(ctx context.Context, userID, query, peerID, groupID string, isGroup bool, messageType proto.MessageType, limit, offset int32) ([]*models.ChatMessage, error) {
	var sqlQuery strings.Builder
	var args []interface{}
	argIndex := 1

	sqlQuery.WriteString(`
		SELECT id, sender_id, receiver_id, group_id, content, message_type,
			   timestamp, is_group, is_read, is_edited, liked_by, reply_to_message_id,
			   attachments, status, delivered_at, read_at, reactions, is_pinned,
			   pinned_at, pinned_by, forward_count, original_message_id, file_metadata,
			   location_data, poll_data, thread_id, parent_message_id, thread_reply_count,
			   edited_at, edit_history, is_scheduled, scheduled_at, mentioned_user_ids,
			   is_system_message, device_info, client_version, is_encrypted, encryption_key_id,
			   created_at, updated_at
		FROM chat_messages
		WHERE content ILIKE $` + fmt.Sprintf("%d", argIndex))
	args = append(args, "%"+query+"%")
	argIndex++

	if isGroup && groupID != "" {
		sqlQuery.WriteString(fmt.Sprintf(" AND group_id = $%d AND is_group = true", argIndex))
		args = append(args, groupID)
		argIndex++
	} else if !isGroup && peerID != "" {
		sqlQuery.WriteString(fmt.Sprintf(" AND ((sender_id = $%d AND receiver_id = $%d) OR (sender_id = $%d AND receiver_id = $%d)) AND is_group = false", argIndex, argIndex+1, argIndex+1, argIndex))
		args = append(args, userID, peerID)
		argIndex += 2
	}

	if messageType != proto.MessageType_MSG_TEXT {
		sqlQuery.WriteString(fmt.Sprintf(" AND message_type = $%d", argIndex))
		args = append(args, messageType)
		argIndex++
	}

	sqlQuery.WriteString(" ORDER BY timestamp DESC")
	sqlQuery.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1))
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, sqlQuery.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		message := &models.ChatMessage{}
		var likedByJSON, attachmentsJSON, reactionsJSON, fileMetadataJSON, locationDataJSON, pollDataJSON, editHistoryJSON, mentionedUserIDsJSON []byte

		err := rows.Scan(
			&message.MessageID, &message.SenderID, &message.ReceiverID, &message.GroupID,
			&message.Content, &message.Type, &message.Timestamp, &message.IsGroup,
			&message.IsRead, &message.IsEdited, &likedByJSON, &message.ReplyToMessageID,
			&attachmentsJSON, &message.Status, &message.DeliveredAt, &message.ReadAt,
			&reactionsJSON, &message.IsPinned, &message.PinnedAt, &message.PinnedBy,
			&message.ForwardCount, &message.OriginalMessageID, &fileMetadataJSON,
			&locationDataJSON, &pollDataJSON, &message.ThreadID, &message.ParentMessageID,
			&message.ThreadReplyCount, &message.EditedAt, &editHistoryJSON, &message.IsScheduled,
			&message.ScheduledAt, &mentionedUserIDsJSON, &message.IsSystemMessage,
			&message.DeviceInfo, &message.ClientVersion, &message.IsEncrypted,
			&message.EncryptionKeyID, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Unmarshal JSON fields
		if len(likedByJSON) > 0 {
			json.Unmarshal(likedByJSON, &message.LikedBy)
		}
		if len(attachmentsJSON) > 0 {
			json.Unmarshal(attachmentsJSON, &message.Attachments)
		}
		if len(reactionsJSON) > 0 {
			json.Unmarshal(reactionsJSON, &message.Reactions)
		}
		if len(fileMetadataJSON) > 0 {
			json.Unmarshal(fileMetadataJSON, &message.FileMetadata)
		}
		if len(locationDataJSON) > 0 {
			json.Unmarshal(locationDataJSON, &message.LocationData)
		}
		if len(pollDataJSON) > 0 {
			json.Unmarshal(pollDataJSON, &message.PollData)
		}
		if len(editHistoryJSON) > 0 {
			json.Unmarshal(editHistoryJSON, &message.EditHistory)
		}
		if len(mentionedUserIDsJSON) > 0 {
			json.Unmarshal(mentionedUserIDsJSON, &message.MentionedUserIDs)
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (r *chatRepository) MarkAsRead(ctx context.Context, userID, peerID string, isGroup bool) error {
	var query string
	var args []interface{}

	if isGroup {
		query = `UPDATE chat_messages SET is_read = true WHERE group_id = $1 AND sender_id != $2 AND is_group = true`
		args = []interface{}{peerID, userID}
	} else {
		query = `UPDATE chat_messages SET is_read = true WHERE sender_id = $1 AND receiver_id = $2 AND is_group = false`
		args = []interface{}{peerID, userID}
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	return nil
}

func (r *chatRepository) UpdateMessageStatus(ctx context.Context, messageID string, status proto.MessageStatus) error {
	query := `UPDATE chat_messages SET status = $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), messageID)
	if err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}

func (r *chatRepository) AddReaction(ctx context.Context, messageID, userID string, reactionType proto.ReactionType) error {
	query := `
		INSERT INTO message_reactions (id, message_id, user_id, reaction_type, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (message_id, user_id, reaction_type) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, uuid.New().String(), messageID, userID, reactionType, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	return nil
}

func (r *chatRepository) RemoveReaction(ctx context.Context, messageID, userID string, reactionType proto.ReactionType) error {
	query := `DELETE FROM message_reactions WHERE message_id = $1 AND user_id = $2 AND reaction_type = $3`

	result, err := r.db.ExecContext(ctx, query, messageID, userID, reactionType)
	if err != nil {
		return fmt.Errorf("failed to remove reaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reaction not found")
	}

	return nil
}

func (r *chatRepository) GetMessageReactions(ctx context.Context, messageID string) ([]*proto.MessageReaction, error) {
	// Validate input
	if messageID == "" {
		return nil, fmt.Errorf("messageID cannot be empty")
	}

	query := `
		SELECT user_id, reaction_type, created_at
		FROM message_reactions
		WHERE message_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to query message reactions for messageID %s: %w", messageID, err)
	}
	defer rows.Close()

	var reactions []*proto.MessageReaction
	for rows.Next() {
		reaction := &proto.MessageReaction{}
		var createdAt time.Time

		err := rows.Scan(
			&reaction.UserId,
			&reaction.ReactionType,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reaction row: %w", err)
		}

		// Convert timestamp
		reaction.Timestamp = timestamppb.New(createdAt)

		reactions = append(reactions, reaction)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return reactions, nil
}

// Fix GetUsersByIDs function around line 369
func (r *chatRepository) GetUsersByIDs(ctx context.Context, userIDs []string) ([]*proto.UserInfo, error) {
	if len(userIDs) == 0 {
		return []*proto.UserInfo{}, nil
	}

	query := `
		SELECT id, full_name, email, avatar_url
		FROM users
		WHERE id = ANY($1)
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(userIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to get users by IDs: %w", err)
	}
	defer rows.Close()

	var users []*proto.UserInfo
	for rows.Next() {
		user := &proto.UserInfo{}
		var avatarUrl sql.NullString

		err := rows.Scan(&user.Id, &user.FullName, &user.Email, &avatarUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if avatarUrl.Valid {
			user.AvatarUrl = avatarUrl.String
		}
		// Set isOnline to false by default - you might want to check this from a presence table
		user.IsOnline = false

		users = append(users, user)
	}

	return users, nil
}

// Fix GetUsersInGroup function around line 409
func (r *chatRepository) GetUsersInGroup(ctx context.Context, groupID string) ([]*proto.UserInfo, error) {
	query := `
		SELECT u.id, u.full_name, u.email, u.avatar_url
		FROM users u
		INNER JOIN group_members gm ON u.id = gm.user_id
		WHERE gm.group_id = $1 AND gm.is_active = true
		ORDER BY u.full_name
	`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users in group: %w", err)
	}
	defer rows.Close()

	var users []*proto.UserInfo
	for rows.Next() {
		user := &proto.UserInfo{}
		var avatarUrl sql.NullString

		err := rows.Scan(&user.Id, &user.FullName, &user.Email, &avatarUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if avatarUrl.Valid {
			user.AvatarUrl = avatarUrl.String
		}
		// Set isOnline to false by default - you might want to check this from a presence table
		user.IsOnline = false

		users = append(users, user)
	}

	return users, nil
}

func (r *chatRepository) UpdateUserPresence(ctx context.Context, userID string, status proto.PresenceStatus, customMessage string) error {
	query := `
		INSERT INTO user_presence (user_id, status, custom_message, last_updated)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET
			status = EXCLUDED.status,
			custom_message = EXCLUDED.custom_message,
			last_updated = EXCLUDED.last_updated
	`

	_, err := r.db.ExecContext(ctx, query, userID, status, customMessage, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update user presence: %w", err)
	}

	return nil
}

func (r *chatRepository) GetUserPresence(ctx context.Context, userIDs []string) ([]*models.UserPresence, error) {
	if len(userIDs) == 0 {
		return []*models.UserPresence{}, nil
	}

	query := `
		SELECT user_id, status, custom_message, last_updated
		FROM user_presence
		WHERE user_id = ANY($1)
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(userIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to get user presence: %w", err)
	}
	defer rows.Close()

	var presences []*models.UserPresence
	for rows.Next() {
		presence := &models.UserPresence{}
		err := rows.Scan(&presence.UserID, &presence.Status, &presence.CustomMessage, &presence.LastSeen)
		if err != nil {
			return nil, fmt.Errorf("failed to scan presence: %w", err)
		}
		presences = append(presences, presence)
	}

	return presences, nil
}

func (r *chatRepository) CreateGroup(ctx context.Context, group *models.GroupInfo, memberIDs []string) (*models.GroupInfo, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if group.ID == "" {
		group.ID = uuid.New().String()
	}
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	// Create group
	groupQuery := `
		INSERT INTO groups (id, name, description, avatar_url, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err = tx.QueryRowContext(ctx, groupQuery, group.ID, group.Name, group.Description,
		group.AvatarURL, group.CreatorID, group.CreatedAt, group.UpdatedAt).Scan(
		&group.ID, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// Add members
	for _, memberID := range memberIDs {
		memberQuery := `
			INSERT INTO group_members (id, group_id, user_id, role, joined_at, is_active)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		role := "member"
		if memberID == group.CreatorID {
			role = "admin"
		}

		_, err = tx.ExecContext(ctx, memberQuery, uuid.New().String(), group.ID, memberID, role, time.Now(), true)
		if err != nil {
			return nil, fmt.Errorf("failed to add group member: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return group, nil
}

func (r *chatRepository) UpdateGroup(ctx context.Context, groupID string, updates map[string]interface{}) (*models.GroupInfo, error) {
	if len(updates) == 0 {
		return r.getGroupByID(ctx, groupID)
	}

	var setParts []string
	var args []interface{}
	argIndex := 1

	for field, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argIndex))
		args = append(args, value)
		argIndex++
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	args = append(args, groupID)

	query := fmt.Sprintf(`
		UPDATE groups 
		SET %s
		WHERE id = $%d
		RETURNING id, name, description, avatar_url, created_by, created_at, updated_at
	`, strings.Join(setParts, ", "), argIndex)

	group := &models.GroupInfo{}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&group.ID, &group.Name, &group.Description, &group.AvatarURL,
		&group.CreatorID, &group.CreatedAt, &group.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("failed to update group: %w", err)
	}

	return group, nil
}

func (r *chatRepository) getGroupByID(ctx context.Context, groupID string) (*models.GroupInfo, error) {
	query := `
		SELECT id, name, description, avatar_url, created_by, created_at, updated_at
		FROM groups
		WHERE id = $1
	`

	group := &models.GroupInfo{}
	err := r.db.QueryRowContext(ctx, query, groupID).Scan(
		&group.ID, &group.Name, &group.Description, &group.AvatarURL,
		&group.CreatorID, &group.CreatedAt, &group.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	return group, nil
}

func (r *chatRepository) AddGroupMembers(ctx context.Context, groupID string, memberIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, memberID := range memberIDs {
		query := `
			INSERT INTO group_members (id, group_id, user_id, role, joined_at, is_active)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (group_id, user_id) DO UPDATE SET
				is_active = true,
				joined_at = EXCLUDED.joined_at
		`

		_, err = tx.ExecContext(ctx, query, uuid.New().String(), groupID, memberID, "member", time.Now(), true)
		if err != nil {
			return fmt.Errorf("failed to add group member %s: %w", memberID, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *chatRepository) RemoveGroupMembers(ctx context.Context, groupID string, memberIDs []string) error {
	query := `UPDATE group_members SET is_active = false WHERE group_id = $1 AND user_id = ANY($2)`

	_, err := r.db.ExecContext(ctx, query, groupID, pq.Array(memberIDs))
	if err != nil {
		return fmt.Errorf("failed to remove group members: %w", err)
	}

	return nil
}

func (r *chatRepository) CreateThread(ctx context.Context, parentMessageID, userID, content string) (*models.ChatMessage, error) {
	// First, get the parent message to inherit some properties
	parentMessage, err := r.GetMessageByID(ctx, parentMessageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent message: %w", err)
	}

	threadMessage := &models.ChatMessage{
		MessageID:       uuid.New().String(),
		SenderID:        userID,
		ReceiverID:      parentMessage.ReceiverID,
		GroupID:         parentMessage.GroupID,
		Content:         content,
		Type:            proto.MessageType_MSG_TEXT,
		Timestamp:       time.Now(),
		IsGroup:         parentMessage.IsGroup,
		ParentMessageID: parentMessageID,
		ThreadID:        parentMessageID, // Use parent message ID as thread ID
		Status:          proto.MessageStatus_SENT,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return r.SaveMessage(ctx, threadMessage)
}

func (r *chatRepository) GetThreadMessages(ctx context.Context, parentMessageID string, limit, offset int32) ([]*models.ChatMessage, error) {
	query := `
		SELECT id, sender_id, receiver_id, group_id, content, message_type,
			   timestamp, is_group, is_read, is_edited, status, reply_to_message_id,
			   thread_id, parent_message_id, is_pinned, is_scheduled, scheduled_at,
			   created_at, updated_at
		FROM chat_messages
		WHERE thread_id = $1
		ORDER BY timestamp ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, parentMessageID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		message := &models.ChatMessage{}
		err := rows.Scan(
			&message.MessageID, &message.SenderID, &message.ReceiverID, &message.GroupID,
			&message.Content, &message.Type, &message.Timestamp, &message.IsGroup,
			&message.IsRead, &message.IsEdited, &message.Status, &message.ReplyToMessageID,
			&message.ThreadID, &message.ParentMessageID, &message.IsPinned,
			&message.IsScheduled, &message.ScheduledAt, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan thread message: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *chatRepository) PinMessage(ctx context.Context, messageID, userID, chatID string, isGroup bool) error {
	query := `
		UPDATE chat_messages 
		SET is_pinned = true, updated_at = $1
		WHERE id = $2 AND (
			(is_group = $3 AND group_id = $4) OR
			(is_group = false AND ((sender_id = $5 AND receiver_id = $4) OR (sender_id = $4 AND receiver_id = $5)))
		)
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), messageID, isGroup, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to pin message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found or access denied")
	}

	return nil
}

func (r *chatRepository) UnpinMessage(ctx context.Context, messageID, userID, chatID string, isGroup bool) error {
	query := `
		UPDATE chat_messages 
		SET is_pinned = false, updated_at = $1
		WHERE id = $2 AND (
			(is_group = $3 AND group_id = $4) OR
			(is_group = false AND ((sender_id = $5 AND receiver_id = $4) OR (sender_id = $4 AND receiver_id = $5)))
		)
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), messageID, isGroup, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to unpin message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found or access denied")
	}

	return nil
}

func (r *chatRepository) GetPinnedMessages(ctx context.Context, chatID string, isGroup bool) ([]*models.ChatMessage, error) {
	var query string
	var args []interface{}

	if isGroup {
		query = `
			SELECT id, sender_id, receiver_id, group_id, content, message_type,
				   timestamp, is_group, is_read, is_edited, status, reply_to_message_id,
				   thread_id, parent_message_id, is_pinned, is_scheduled, scheduled_at,
				   created_at, updated_at
			FROM chat_messages
			WHERE group_id = $1 AND is_group = true AND is_pinned = true
			ORDER BY timestamp DESC
		`
		args = []interface{}{chatID}
	} else {
		query = `
			SELECT id, sender_id, receiver_id, group_id, content, message_type,
				   timestamp, is_group, is_read, is_edited, status, reply_to_message_id,
				   thread_id, parent_message_id, is_pinned, is_scheduled, scheduled_at,
				   created_at, updated_at
			FROM chat_messages
			WHERE ((sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1))
				AND is_group = false AND is_pinned = true
			ORDER BY timestamp DESC
		`
		// For direct messages, chatID should be in format "userID1:userID2"
		userIDs := strings.Split(chatID, ":")
		if len(userIDs) != 2 {
			return nil, fmt.Errorf("invalid chat ID format for direct message")
		}
		args = []interface{}{userIDs[0], userIDs[1]}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get pinned messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		message := &models.ChatMessage{}
		err := rows.Scan(
			&message.MessageID, &message.SenderID, &message.ReceiverID, &message.GroupID,
			&message.Content, &message.Type, &message.Timestamp, &message.IsGroup,
			&message.IsRead, &message.IsEdited, &message.Status, &message.ReplyToMessageID,
			&message.ThreadID, &message.ParentMessageID, &message.IsPinned,
			&message.IsScheduled, &message.ScheduledAt, &message.CreatedAt, &message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pinned message: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *chatRepository) DeleteMessage(ctx context.Context, messageID string) error {
	query := `DELETE FROM chat_messages WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, messageID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}
