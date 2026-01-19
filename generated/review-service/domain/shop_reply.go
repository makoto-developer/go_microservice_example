package domain

import (
	"github.com/google/uuid"
	"time"
)

// ShopReply represents ShopReply
type ShopReply struct {
	Id uuid.UUID `db:"id" json:"id"`
	ReviewId uuid.UUID `db:"review_id" json:"review_id"`
	ShopId uuid.UUID `db:"shop_id" json:"shop_id"`
	Content text `db:"content" json:"content"`
	Edited bool `db:"edited" json:"edited"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewShopReply creates a new ShopReply instance
func NewShopReply() *ShopReply {
	return &ShopReply{}
}
