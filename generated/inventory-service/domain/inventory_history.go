package domain

import (
	"github.com/google/uuid"
	"time"
)

// InventoryHistory represents InventoryHistory
type InventoryHistory struct {
	Id uuid.UUID `db:"id" json:"id"`
	InventoryId uuid.UUID `db:"inventory_id" json:"inventory_id"`
	ChangeType ChangeType `db:"change_type" json:"change_type"`
	ChangeQuantity int `db:"change_quantity" json:"change_quantity"`
	QuantityBefore int `db:"quantity_before" json:"quantity_before"`
	QuantityAfter int `db:"quantity_after" json:"quantity_after"`
	Reason string `db:"reason" json:"reason"`
	Operator string `db:"operator" json:"operator"`
	OrderId *uuid.UUID `db:"order_id" json:"order_id,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewInventoryHistory creates a new InventoryHistory instance
func NewInventoryHistory() *InventoryHistory {
	return &InventoryHistory{}
}
