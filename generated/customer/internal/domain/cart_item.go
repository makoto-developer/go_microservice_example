package domain

import (
	"time"

	"github.com/google/uuid"
)

// CartItem represents an item in customer's cart
type CartItem struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	CustomerID  uuid.UUID  `db:"customer_id" json:"customer_id"`
	ProductID   uuid.UUID  `db:"product_id" json:"product_id"`
	VariationID *uuid.UUID `db:"variation_id" json:"variation_id,omitempty"`
	Quantity    int        `db:"quantity" json:"quantity"`
	ExpiresAt   *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// NewCartItem creates a new CartItem with 30 days expiration
func NewCartItem(customerID, productID uuid.UUID, variationID *uuid.UUID, quantity int) *CartItem {
	expiresAt := time.Now().Add(30 * 24 * time.Hour) // 30 days
	return &CartItem{
		ID:          uuid.New(),
		CustomerID:  customerID,
		ProductID:   productID,
		VariationID: variationID,
		Quantity:    quantity,
		ExpiresAt:   &expiresAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
