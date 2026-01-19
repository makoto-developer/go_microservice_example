package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductImage struct {
	ID               uuid.UUID `db:"id" json:"id"`
	ProductID        uuid.UUID `db:"product_id" json:"product_id"`
	URL              string    `db:"url" json:"url"`
	DisplayOrder     int       `db:"display_order" json:"display_order"`
	Thumbnail200URL  string    `db:"thumbnail_200_url" json:"thumbnail_200_url"`
	Thumbnail400URL  string    `db:"thumbnail_400_url" json:"thumbnail_400_url"`
	Thumbnail800URL  string    `db:"thumbnail_800_url" json:"thumbnail_800_url"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}

func NewProductImage(productID uuid.UUID, url string, displayOrder int, thumb200, thumb400, thumb800 string) *ProductImage {
	return &ProductImage{
		ID:              uuid.New(),
		ProductID:       productID,
		URL:             url,
		DisplayOrder:    displayOrder,
		Thumbnail200URL: thumb200,
		Thumbnail400URL: thumb400,
		Thumbnail800URL: thumb800,
		CreatedAt:       time.Now(),
	}
}
