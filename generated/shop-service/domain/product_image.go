package domain

import (
	"time"

	"github.com/google/uuid"
)

// ProductImage represents ProductImage
type ProductImage struct {
	Id              uuid.UUID `db:"id" json:"id"`
	ProductId       uuid.UUID `db:"product_id" json:"product_id"`
	Url             string    `db:"url" json:"url"`
	DisplayOrder    int       `db:"display_order" json:"display_order"`
	Thumbnail200Url string    `db:"thumbnail_200_url" json:"thumbnail_200_url"`
	Thumbnail400Url string    `db:"thumbnail_400_url" json:"thumbnail_400_url"`
	Thumbnail800Url string    `db:"thumbnail_800_url" json:"thumbnail_800_url"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

// NewProductImage creates a new ProductImage instance
func NewProductImage() *ProductImage {
	return &ProductImage{}
}
