package domain

import (
	"time"

	"github.com/google/uuid"
)

// ProductImage represents a product image
type ProductImage struct {
	ID              uuid.UUID `db:"id" json:"id"`
	ProductID       uuid.UUID `db:"product_id" json:"product_id"`
	URL             string    `db:"url" json:"url"`
	ThumbnailURL    string    `db:"thumbnail_url" json:"thumbnail_url"`
	MediumURL       string    `db:"medium_url" json:"medium_url"`
	LargeURL        string    `db:"large_url" json:"large_url"`
	DisplayOrder    int       `db:"display_order" json:"display_order"`
	IsPrimary       bool      `db:"is_primary" json:"is_primary"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}
