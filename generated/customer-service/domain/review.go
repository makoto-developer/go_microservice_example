package domain

import (
	"github.com/google/uuid"
	"time"
)

// Review represents Review
type Review struct {
	Id uuid.UUID `db:"id" json:"id"`
	CustomerId uuid.UUID `db:"customer_id" json:"customer_id"`
	ProductId uuid.UUID `db:"product_id" json:"product_id"`
	OrderId uuid.UUID `db:"order_id" json:"order_id"`
	Rating int `db:"rating" json:"rating"`
	ReviewText text `db:"review_text" json:"review_text"`
	ImageUrls []string `db:"image_urls" json:"image_urls"`
	EditableUntil time.Time `db:"editable_until" json:"editable_until"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewReview creates a new Review instance
func NewReview() *Review {
	return &Review{}
}
