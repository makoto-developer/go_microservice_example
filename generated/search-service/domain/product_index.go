package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// ProductIndex represents ProductIndex
type ProductIndex struct {
	Id uuid.UUID `db:"id" json:"id"`
	ProductId uuid.UUID `db:"product_id" json:"product_id"`
	ProductName string `db:"product_name" json:"product_name"`
	Description text `db:"description" json:"description"`
	Category string `db:"category" json:"category"`
	Tags []string `db:"tags" json:"tags"`
	ShopId uuid.UUID `db:"shop_id" json:"shop_id"`
	ShopName string `db:"shop_name" json:"shop_name"`
	Price decimal.Decimal `db:"price" json:"price"`
	AverageRating decimal.Decimal `db:"average_rating" json:"average_rating"`
	ReviewCount int `db:"review_count" json:"review_count"`
	StockStatus StockStatus `db:"stock_status" json:"stock_status"`
	ImageUrl *string `db:"image_url" json:"image_url,omitempty"`
	IndexedAt time.Time `db:"indexed_at" json:"indexed_at"`
}

// NewProductIndex creates a new ProductIndex instance
func NewProductIndex() *ProductIndex {
	return &ProductIndex{}
}
