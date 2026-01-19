package domain

import (
	"time"

	"github.com/google/uuid"
)

type Favorite struct {
	ID              uuid.UUID `db:"id" json:"id"`
	CustomerID      uuid.UUID `db:"customer_id" json:"customer_id"`
	ProductID       uuid.UUID `db:"product_id" json:"product_id"`
	NotifyOnRestock bool      `db:"notify_on_restock" json:"notify_on_restock"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

func NewFavorite(customerID, productID uuid.UUID, notifyOnRestock bool) *Favorite {
	return &Favorite{
		ID:              uuid.New(),
		CustomerID:      customerID,
		ProductID:       productID,
		NotifyOnRestock: notifyOnRestock,
		CreatedAt:       time.Now(),
	}
}
