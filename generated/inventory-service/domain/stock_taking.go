package domain

import (
	"github.com/google/uuid"
	"time"
)

// StockTaking represents StockTaking
type StockTaking struct {
	Id uuid.UUID `db:"id" json:"id"`
	InventoryId uuid.UUID `db:"inventory_id" json:"inventory_id"`
	ShopId uuid.UUID `db:"shop_id" json:"shop_id"`
	SystemQuantity int `db:"system_quantity" json:"system_quantity"`
	ActualQuantity int `db:"actual_quantity" json:"actual_quantity"`
	Difference int `db:"difference" json:"difference"`
	DifferenceReason *string `db:"difference_reason" json:"difference_reason,omitempty"`
	Operator string `db:"operator" json:"operator"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewStockTaking creates a new StockTaking instance
func NewStockTaking() *StockTaking {
	return &StockTaking{}
}
