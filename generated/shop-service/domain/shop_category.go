package domain

import (
	"time"

	"github.com/google/uuid"
)

// ShopCategory represents ShopCategory
type ShopCategory struct {
	Id           uuid.UUID `db:"id" json:"id"`
	ShopId       uuid.UUID `db:"shop_id" json:"shop_id"`
	CategoryName string    `db:"category_name" json:"category_name"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// NewShopCategory creates a new ShopCategory instance
func NewShopCategory() *ShopCategory {
	return &ShopCategory{}
}
