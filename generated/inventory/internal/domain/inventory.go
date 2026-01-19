package domain

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID               uuid.UUID `db:"id" json:"id"`
	ProductID        uuid.UUID `db:"product_id" json:"product_id"`
	ShopID           uuid.UUID `db:"shop_id" json:"shop_id"`
	Quantity         int       `db:"quantity" json:"quantity"`
	ReservedQuantity int       `db:"reserved_quantity" json:"reserved_quantity"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

func NewInventory(productID, shopID uuid.UUID, quantity int) *Inventory {
	now := time.Now()
	return &Inventory{
		ID:               uuid.New(),
		ProductID:        productID,
		ShopID:           shopID,
		Quantity:         quantity,
		ReservedQuantity: 0,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func (i *Inventory) AvailableQuantity() int {
	return i.Quantity - i.ReservedQuantity
}

func (i *Inventory) CanReserve(quantity int) bool {
	return i.AvailableQuantity() >= quantity
}
