package domain

import (
	"time"

	"github.com/google/uuid"
)

type ShopCategory struct {
	ID           uuid.UUID `db:"id" json:"id"`
	ShopID       uuid.UUID `db:"shop_id" json:"shop_id"`
	CategoryName string    `db:"category_name" json:"category_name"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

func NewShopCategory(shopID uuid.UUID, categoryName string) *ShopCategory {
	return &ShopCategory{
		ID:           uuid.New(),
		ShopID:       shopID,
		CategoryName: categoryName,
		CreatedAt:    time.Now(),
	}
}
