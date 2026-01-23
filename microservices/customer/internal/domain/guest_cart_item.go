package domain

import (
	"time"

	"github.com/google/uuid"
)

// GuestCartItem represents an item in guest's cart
type GuestCartItem struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	SessionID   string     `db:"session_id" json:"session_id"`
	GuestToken  string     `db:"guest_token" json:"guest_token"`
	ProductID   uuid.UUID  `db:"product_id" json:"product_id"`
	VariationID *uuid.UUID `db:"variation_id" json:"variation_id,omitempty"`
	Quantity    int        `db:"quantity" json:"quantity"`
	ExpiresAt   *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}
