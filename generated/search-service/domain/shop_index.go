package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// ShopIndex represents ShopIndex
type ShopIndex struct {
	Id uuid.UUID `db:"id" json:"id"`
	ShopId uuid.UUID `db:"shop_id" json:"shop_id"`
	ShopName string `db:"shop_name" json:"shop_name"`
	Description text `db:"description" json:"description"`
	Categories []string `db:"categories" json:"categories"`
	AverageRating decimal.Decimal `db:"average_rating" json:"average_rating"`
	ProductCount int `db:"product_count" json:"product_count"`
	LogoUrl *string `db:"logo_url" json:"logo_url,omitempty"`
	IndexedAt time.Time `db:"indexed_at" json:"indexed_at"`
}

// NewShopIndex creates a new ShopIndex instance
func NewShopIndex() *ShopIndex {
	return &ShopIndex{}
}
