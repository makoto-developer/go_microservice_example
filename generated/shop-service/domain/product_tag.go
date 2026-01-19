package domain

import (
	"time"

	"github.com/google/uuid"
)

// ProductTag represents ProductTag
type ProductTag struct {
	Id        uuid.UUID `db:"id" json:"id"`
	ProductId uuid.UUID `db:"product_id" json:"product_id"`
	TagName   string    `db:"tag_name" json:"tag_name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewProductTag creates a new ProductTag instance
func NewProductTag() *ProductTag {
	return &ProductTag{}
}
