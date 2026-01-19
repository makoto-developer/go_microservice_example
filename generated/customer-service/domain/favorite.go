package domain

import (
	"github.com/google/uuid"
	"time"
)

// Favorite represents Favorite
type Favorite struct {
	Id uuid.UUID `db:"id" json:"id"`
	CustomerId uuid.UUID `db:"customer_id" json:"customer_id"`
	ProductId uuid.UUID `db:"product_id" json:"product_id"`
	NotifyOnRestock bool `db:"notify_on_restock" json:"notify_on_restock"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewFavorite creates a new Favorite instance
func NewFavorite() *Favorite {
	return &Favorite{}
}
