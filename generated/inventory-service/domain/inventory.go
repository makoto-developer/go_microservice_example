package domain

import (
	"github.com/google/uuid"
	"time"
)

// Inventory represents Inventory
type Inventory struct {
	Id uuid.UUID `db:"id" json:"id"`
	ProductId uuid.UUID `db:"product_id" json:"product_id"`
	VariationId *uuid.UUID `db:"variation_id" json:"variation_id,omitempty"`
	ShopId uuid.UUID `db:"shop_id" json:"shop_id"`
	Quantity int `db:"quantity" json:"quantity"`
	ReservedQuantity int `db:"reserved_quantity" json:"reserved_quantity"`
	AvailableQuantity int `db:"available_quantity" json:"available_quantity"`
	AlertThreshold int `db:"alert_threshold" json:"alert_threshold"`
	LastAlertedAt *time.Time `db:"last_alerted_at" json:"last_alerted_at,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewInventory creates a new Inventory instance
func NewInventory() *Inventory {
	return &Inventory{}
}
