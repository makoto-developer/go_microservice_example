package domain

import (
	"time"
	"github.com/google/uuid"
)

// GuestCartItem represents GuestCartItem
type GuestCartItem struct {
	Id uuid.UUID `db:"id" json:"id"`
	SessionId string `db:"session_id" json:"session_id"`
	ProductId uuid.UUID `db:"product_id" json:"product_id"`
	VariationId *uuid.UUID `db:"variation_id" json:"variation_id,omitempty"`
	Quantity int `db:"quantity" json:"quantity"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewGuestCartItem creates a new GuestCartItem instance
func NewGuestCartItem() *GuestCartItem {
	return &GuestCartItem{}
}
