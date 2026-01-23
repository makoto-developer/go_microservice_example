package domain

import (
	"time"
	"github.com/google/uuid"
)

type Inventory struct {
	ID               uuid.UUID `db:"id" json:"id"`
	ProductID        uuid.UUID `db:"product_id" json:"product_id"`
	VariationID      *uuid.UUID `db:"variation_id" json:"variation_id,omitempty"`
	ShopID           uuid.UUID `db:"shop_id" json:"shop_id"`
	Quantity         int       `db:"quantity" json:"quantity"`
	ReservedQuantity int       `db:"reserved_quantity" json:"reserved_quantity"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

func (i *Inventory) AvailableQuantity() int {
	return i.Quantity - i.ReservedQuantity
}
