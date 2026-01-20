package domain

import (
	"time"

	"github.com/google/uuid"
)

// Favorite represents a customer's favorite product
type Favorite struct {
	ID         uuid.UUID `db:"id" json:"id"`
	CustomerID uuid.UUID `db:"customer_id" json:"customer_id"`
	ProductID  uuid.UUID `db:"product_id" json:"product_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

// NewFavorite creates a new Favorite
func NewFavorite(customerID, productID uuid.UUID) *Favorite {
	return &Favorite{
		ID:         uuid.New(),
		CustomerID: customerID,
		ProductID:  productID,
		CreatedAt:  time.Now(),
	}
}
