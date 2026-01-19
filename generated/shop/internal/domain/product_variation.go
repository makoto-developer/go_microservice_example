package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductVariation struct {
	ID             uuid.UUID `db:"id" json:"id"`
	ProductID      uuid.UUID `db:"product_id" json:"product_id"`
	SKU            string    `db:"sku" json:"sku"`
	AttributeName  string    `db:"attribute_name" json:"attribute_name"`
	AttributeValue string    `db:"attribute_value" json:"attribute_value"`
	Price          float64   `db:"price" json:"price"`
	StockQuantity  int       `db:"stock_quantity" json:"stock_quantity"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

func NewProductVariation(productID uuid.UUID, sku, attrName, attrValue string, price float64, stock int) *ProductVariation {
	now := time.Now()
	return &ProductVariation{
		ID:             uuid.New(),
		ProductID:      productID,
		SKU:            sku,
		AttributeName:  attrName,
		AttributeValue: attrValue,
		Price:          price,
		StockQuantity:  stock,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
