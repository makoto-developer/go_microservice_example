package domain

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	CustomerID  uuid.UUID  `db:"customer_id" json:"customer_id"`
	ProductID   uuid.UUID  `db:"product_id" json:"product_id"`
	VariationID *uuid.UUID `db:"variation_id" json:"variation_id,omitempty"`
	Quantity    int        `db:"quantity" json:"quantity"`
	ExpiresAt   time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

func NewCartItem(customerID, productID uuid.UUID, variationID *uuid.UUID, quantity int) *CartItem {
	now := time.Now()
	return &CartItem{
		ID:          uuid.New(),
		CustomerID:  customerID,
		ProductID:   productID,
		VariationID: variationID,
		Quantity:    quantity,
		ExpiresAt:   now.Add(7 * 24 * time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (c *CartItem) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

type GuestCartItem struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	SessionID   string     `db:"session_id" json:"session_id"`
	ProductID   uuid.UUID  `db:"product_id" json:"product_id"`
	VariationID *uuid.UUID `db:"variation_id" json:"variation_id,omitempty"`
	Quantity    int        `db:"quantity" json:"quantity"`
	ExpiresAt   time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

func NewGuestCartItem(sessionID string, productID uuid.UUID, variationID *uuid.UUID, quantity int) *GuestCartItem {
	now := time.Now()
	return &GuestCartItem{
		ID:          uuid.New(),
		SessionID:   sessionID,
		ProductID:   productID,
		VariationID: variationID,
		Quantity:    quantity,
		ExpiresAt:   now.Add(24 * time.Hour),
		CreatedAt:   now,
	}
}
