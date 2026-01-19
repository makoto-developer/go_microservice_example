package domain

import (
	"time"
	"github.com/google/uuid"
)

// Message represents Message
type Message struct {
	Id uuid.UUID `db:"id" json:"id"`
	RoomId uuid.UUID `db:"room_id" json:"room_id"`
	SenderId uuid.UUID `db:"sender_id" json:"sender_id"`
	ReceiverId uuid.UUID `db:"receiver_id" json:"receiver_id"`
	MessageType MessageType `db:"message_type" json:"message_type"`
	Content *text `db:"content" json:"content,omitempty"`
	FileUrl *string `db:"file_url" json:"file_url,omitempty"`
	FileName *string `db:"file_name" json:"file_name,omitempty"`
	FileSize *int `db:"file_size" json:"file_size,omitempty"`
	IsRead bool `db:"is_read" json:"is_read"`
	ReadAt *time.Time `db:"read_at" json:"read_at,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewMessage creates a new Message instance
func NewMessage() *Message {
	return &Message{}
}
