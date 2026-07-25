package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidRating     = errors.New("rating must be between 1 and 5")
	ErrReviewNotEditable = errors.New("review is no longer editable")
)

type Review struct {
	ID            uuid.UUID `db:"id" json:"id"`
	CustomerID    uuid.UUID `db:"customer_id" json:"customer_id"`
	ProductID     uuid.UUID `db:"product_id" json:"product_id"`
	OrderID       uuid.UUID `db:"order_id" json:"order_id"`
	Rating        int       `db:"rating" json:"rating"` // 1-5
	ReviewText    string    `db:"review_text" json:"review_text"`
	EditableUntil time.Time `db:"editable_until" json:"editable_until"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func NewReview(customerID, productID, orderID uuid.UUID, rating int, reviewText string) (*Review, error) {
	if rating < 1 || rating > 5 {
		return nil, ErrInvalidRating
	}

	now := time.Now()
	return &Review{
		ID:            uuid.New(),
		CustomerID:    customerID,
		ProductID:     productID,
		OrderID:       orderID,
		Rating:        rating,
		ReviewText:    reviewText,
		EditableUntil: now.Add(30 * 24 * time.Hour), // 30 days
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (r *Review) CanEdit() bool {
	return time.Now().Before(r.EditableUntil)
}

func (r *Review) Update(rating int, reviewText string) error {
	if !r.CanEdit() {
		return ErrReviewNotEditable
	}
	if rating < 1 || rating > 5 {
		return ErrInvalidRating
	}

	r.Rating = rating
	r.ReviewText = reviewText
	r.UpdatedAt = time.Now()
	return nil
}
