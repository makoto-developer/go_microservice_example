package domain

import (
	"github.com/google/uuid"
	"time"
)

// TypingIndicator represents TypingIndicator
type TypingIndicator struct {
	Id uuid.UUID `db:"id" json:"id"`
	RoomId uuid.UUID `db:"room_id" json:"room_id"`
	UserId uuid.UUID `db:"user_id" json:"user_id"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewTypingIndicator creates a new TypingIndicator instance
func NewTypingIndicator() *TypingIndicator {
	return &TypingIndicator{}
}
