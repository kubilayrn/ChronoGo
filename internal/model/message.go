package model

import (
	"time"

	"github.com/google/uuid"
)

type MessageStatus string

const (
	StatusUnsent MessageStatus = "unsent"
	StatusSent   MessageStatus = "sent"
)

type Message struct {
	ID        int           `json:"id"`
	To        string        `json:"to"`
	Content   string        `json:"content"`
	Status    MessageStatus `json:"status"`
	SentAt    *time.Time    `json:"sent_at,omitempty"`
	MessageID *uuid.UUID    `json:"message_id,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
