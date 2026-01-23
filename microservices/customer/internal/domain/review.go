package domain

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID            uuid.UUID `db:"id" json:"id"`
	CustomerID    uuid.UUID `db:"customer_id" json:"customer_id"`
	ProductID     uuid.UUID `db:"product_id" json:"product_id"`
	OrderID       uuid.UUID `db:"order_id" json:"order_id"`
	Rating        int       `db:"rating" json:"rating"`
	ReviewText    string    `db:"review_text" json:"review_text"`
	EditableUntil time.Time `db:"editable_until" json:"editable_until"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func NewReview(customerID, productID, orderID uuid.UUID, rating int, reviewText string) *Review {
	now := time.Now()
	return &Review{
		ID:            uuid.New(),
		CustomerID:    customerID,
		ProductID:     productID,
		OrderID:       orderID,
		Rating:        rating,
		ReviewText:    reviewText,
		EditableUntil: now.Add(30 * 24 * time.Hour),
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (r *Review) IsEditable() bool {
	return time.Now().Before(r.EditableUntil)
}

func (r *Review) IsValidRating() bool {
	return r.Rating >= 1 && r.Rating <= 5
}

type ReviewImage struct {
	ID        uuid.UUID `db:"id" json:"id"`
	ReviewID  uuid.UUID `db:"review_id" json:"review_id"`
	ImageURL  string    `db:"image_url" json:"image_url"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewReviewImage(reviewID uuid.UUID, imageURL string) *ReviewImage {
	return &ReviewImage{
		ID:        uuid.New(),
		ReviewID:  reviewID,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
	}
}
