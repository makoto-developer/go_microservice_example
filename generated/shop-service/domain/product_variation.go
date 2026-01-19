package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductVariation represents ProductVariation
type ProductVariation struct {
	Id             uuid.UUID       `db:"id" json:"id"`
	ProductId      uuid.UUID       `db:"product_id" json:"product_id"`
	Sku            string          `db:"sku" json:"sku"`
	AttributeName  string          `db:"attribute_name" json:"attribute_name"`
	AttributeValue string          `db:"attribute_value" json:"attribute_value"`
	Price          decimal.Decimal `db:"price" json:"price"`
	StockQuantity  int             `db:"stock_quantity" json:"stock_quantity"`
	CreatedAt      time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at" json:"updated_at"`
}

// NewProductVariation creates a new ProductVariation instance
func NewProductVariation() *ProductVariation {
	return &ProductVariation{}
}
