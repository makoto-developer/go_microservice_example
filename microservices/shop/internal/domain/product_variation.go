package domain

import (
	"time"

	"github.com/google/uuid"
)

// ProductVariation represents a product variation (e.g., size, color)
type ProductVariation struct {
	ID         uuid.UUID `db:"id" json:"id"`
	ProductID  uuid.UUID `db:"product_id" json:"product_id"`
	SKU        string    `db:"sku" json:"sku"`
	Name       string    `db:"name" json:"name"`
	Attributes map[string]string `db:"attributes" json:"attributes"`
	Price      int64     `db:"price" json:"price"`
	StockCount int       `db:"stock_count" json:"stock_count"`
	IsActive   bool      `db:"is_active" json:"is_active"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}
