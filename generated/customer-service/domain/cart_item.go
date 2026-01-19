package domain

import (
	"time"
	"github.com/google/uuid"
)

// CartItem represents CartItem
type CartItem struct {
	Id uuid.UUID `db:"id" json:"id"`
	CustomerId uuid.UUID `db:"customer_id" json:"customer_id"`
	ProductId uuid.UUID `db:"product_id" json:"product_id"`
	VariationId *uuid.UUID `db:"variation_id" json:"variation_id,omitempty"`
	Quantity int `db:"quantity" json:"quantity"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewCartItem creates a new CartItem instance
func NewCartItem() *CartItem {
	return &CartItem{}
}
