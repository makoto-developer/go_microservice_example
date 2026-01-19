package domain

import (
	"github.com/google/uuid"
	"time"
)

// Review represents Review
type Review struct {
	Id uuid.UUID `db:"id" json:"id"`
	ProductId uuid.UUID `db:"product_id" json:"product_id"`
	OrderId uuid.UUID `db:"order_id" json:"order_id"`
	CustomerId uuid.UUID `db:"customer_id" json:"customer_id"`
	Nickname *string `db:"nickname" json:"nickname,omitempty"`
	Rating int `db:"rating" json:"rating"`
	Title string `db:"title" json:"title"`
	Content text `db:"content" json:"content"`
	ImageUrls []string `db:"image_urls" json:"image_urls"`
	Status ReviewStatus `db:"status" json:"status"`
	RejectionReason *string `db:"rejection_reason" json:"rejection_reason,omitempty"`
	EditableUntil time.Time `db:"editable_until" json:"editable_until"`
	Edited bool `db:"edited" json:"edited"`
	Deleted bool `db:"deleted" json:"deleted"`
	HelpfulCount int `db:"helpful_count" json:"helpful_count"`
	ReportedCount int `db:"reported_count" json:"reported_count"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewReview creates a new Review instance
func NewReview() *Review {
	return &Review{}
}
