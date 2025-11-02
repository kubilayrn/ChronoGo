package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kubilayrn/ChronoGo/internal/database"
	"github.com/kubilayrn/ChronoGo/internal/model"
)

type MessageRepository struct{}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

func (r *MessageRepository) GetUnsentMessages(ctx context.Context, limit int) ([]model.Message, error) {
	query := `
		SELECT id, "to", content, status, sent_at, message_id, created_at, updated_at
		FROM messages
		WHERE status = 'unsent'
		ORDER BY created_at ASC
		LIMIT $1
	`

	rows, err := database.DB.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query unsent messages: %w", err)
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message
		var sentAt pgtype.Timestamp
		var messageID *uuid.UUID

		err := rows.Scan(
			&msg.ID,
			&msg.To,
			&msg.Content,
			&msg.Status,
			&sentAt,
			&messageID,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		if sentAt.Valid {
			msg.SentAt = &sentAt.Time
		}
		msg.MessageID = messageID

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}

func (r *MessageRepository) UpdateMessageStatus(
	ctx context.Context,
	id int,
	status model.MessageStatus,
	messageID *uuid.UUID,
	sentAt *time.Time,
) error {
	query := `
		UPDATE messages
		SET status = $1, message_id = $2, sent_at = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`

	_, err := database.DB.Exec(ctx, query, status, messageID, sentAt, id)
	if err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	return nil
}

func (r *MessageRepository) GetSentMessages(ctx context.Context) ([]model.Message, error) {
	query := `
		SELECT id, "to", content, status, sent_at, message_id, created_at, updated_at
		FROM messages
		WHERE status = 'sent'
		ORDER BY sent_at DESC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sent messages: %w", err)
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var msg model.Message
		var sentAt pgtype.Timestamp
		var messageID *uuid.UUID

		err := rows.Scan(
			&msg.ID,
			&msg.To,
			&msg.Content,
			&msg.Status,
			&sentAt,
			&messageID,
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		if sentAt.Valid {
			msg.SentAt = &sentAt.Time
		}
		msg.MessageID = messageID

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}
