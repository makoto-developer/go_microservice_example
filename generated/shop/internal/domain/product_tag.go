package domain

import (
	"time"

	"github.com/google/uuid"
)

type ProductTag struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`
	TagName   string    `db:"tag_name" json:"tag_name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewProductTag(productID uuid.UUID, tagName string) *ProductTag {
	return &ProductTag{
		ID:        uuid.New(),
		ProductID: productID,
		TagName:   tagName,
		CreatedAt: time.Now(),
	}
}
