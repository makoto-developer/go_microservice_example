package domain

import (
	"github.com/shopspring/decimal"
	"time"
	"github.com/google/uuid"
)

// ProductRating represents ProductRating
type ProductRating struct {
	Id uuid.UUID `db:"id" json:"id"`
	ProductId uuid.UUID `db:"product_id" json:"product_id"`
	TotalReviews int `db:"total_reviews" json:"total_reviews"`
	AverageRating decimal.Decimal `db:"average_rating" json:"average_rating"`
	Rating1Count int `db:"rating_1_count" json:"rating_1_count"`
	Rating2Count int `db:"rating_2_count" json:"rating_2_count"`
	Rating3Count int `db:"rating_3_count" json:"rating_3_count"`
	Rating4Count int `db:"rating_4_count" json:"rating_4_count"`
	Rating5Count int `db:"rating_5_count" json:"rating_5_count"`
	ReviewsWithImages int `db:"reviews_with_images" json:"reviews_with_images"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewProductRating creates a new ProductRating instance
func NewProductRating() *ProductRating {
	return &ProductRating{}
}
