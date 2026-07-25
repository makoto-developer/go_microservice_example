package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
)

type ChatRoom struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	CustomerID    uuid.UUID  `db:"customer_id" json:"customer_id"`
	ShopID        uuid.UUID  `db:"shop_id" json:"shop_id"`
	LastMessage   *string    `db:"last_message" json:"last_message,omitempty"`
	LastMessageAt *time.Time `db:"last_message_at" json:"last_message_at,omitempty"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}

type Message struct {
	ID         uuid.UUID   `db:"id" json:"id"`
	RoomID     uuid.UUID   `db:"room_id" json:"room_id"`
	SenderID   uuid.UUID   `db:"sender_id" json:"sender_id"`
	SenderType string      `db:"sender_type" json:"sender_type"` // "customer" or "shop"
	Type       MessageType `db:"type" json:"type"`
	Content    string      `db:"content" json:"content"`
	ReadAt     *time.Time  `db:"read_at" json:"read_at,omitempty"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
}

func NewChatRoom(customerID, shopID uuid.UUID) *ChatRoom {
	return &ChatRoom{
		ID:         uuid.New(),
		CustomerID: customerID,
		ShopID:     shopID,
		CreatedAt:  time.Now(),
	}
}

func NewMessage(roomID, senderID uuid.UUID, senderType string, msgType MessageType, content string) *Message {
	return &Message{
		ID:         uuid.New(),
		RoomID:     roomID,
		SenderID:   senderID,
		SenderType: senderType,
		Type:       msgType,
		Content:    content,
		CreatedAt:  time.Now(),
	}
}

func (m *Message) MarkRead() {
	now := time.Now()
	m.ReadAt = &now
}

func (m *Message) IsRead() bool {
	return m.ReadAt != nil
}
