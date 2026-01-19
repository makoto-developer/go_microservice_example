package domain

import (
	"github.com/google/uuid"
	"time"
)

// ChatRoom represents ChatRoom
type ChatRoom struct {
	Id uuid.UUID `db:"id" json:"id"`
	CustomerId uuid.UUID `db:"customer_id" json:"customer_id"`
	ShopId uuid.UUID `db:"shop_id" json:"shop_id"`
	ProductId *uuid.UUID `db:"product_id" json:"product_id,omitempty"`
	Status RoomStatus `db:"status" json:"status"`
	LastMessageAt *time.Time `db:"last_message_at" json:"last_message_at,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewChatRoom creates a new ChatRoom instance
func NewChatRoom() *ChatRoom {
	return &ChatRoom{}
}
