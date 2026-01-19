package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Product represents Product
type Product struct {
	Id            uuid.UUID        `db:"id" json:"id"`
	ShopId        uuid.UUID        `db:"shop_id" json:"shop_id"`
	Name          string           `db:"name" json:"name"`
	Description   string           `db:"description" json:"description"`
	Price         decimal.Decimal  `db:"price" json:"price"`
	Category      string           `db:"category" json:"category"`
	StockQuantity int              `db:"stock_quantity" json:"stock_quantity"`
	Weight        *decimal.Decimal `db:"weight" json:"weight,omitempty"`
	Size          *string          `db:"size" json:"size,omitempty"`
	JanCode       *string          `db:"jan_code" json:"jan_code,omitempty"`
	Published     bool             `db:"published" json:"published"`
	Deleted       bool             `db:"deleted" json:"deleted"`
	CreatedAt     time.Time        `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time        `db:"updated_at" json:"updated_at"`
}

// NewProduct creates a new Product instance
func NewProduct() *Product {
	return &Product{}
}
