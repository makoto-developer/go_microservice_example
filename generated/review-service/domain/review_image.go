package domain

import (
	"github.com/google/uuid"
	"time"
)

// ReviewImage represents ReviewImage
type ReviewImage struct {
	Id uuid.UUID `db:"id" json:"id"`
	ReviewId uuid.UUID `db:"review_id" json:"review_id"`
	Url string `db:"url" json:"url"`
	Thumbnail200Url string `db:"thumbnail_200_url" json:"thumbnail_200_url"`
	Thumbnail400Url string `db:"thumbnail_400_url" json:"thumbnail_400_url"`
	DisplayOrder int `db:"display_order" json:"display_order"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewReviewImage creates a new ReviewImage instance
func NewReviewImage() *ReviewImage {
	return &ReviewImage{}
}
