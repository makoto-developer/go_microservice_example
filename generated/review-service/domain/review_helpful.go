package domain

import (
	"github.com/google/uuid"
	"time"
)

// ReviewHelpful represents ReviewHelpful
type ReviewHelpful struct {
	Id uuid.UUID `db:"id" json:"id"`
	ReviewId uuid.UUID `db:"review_id" json:"review_id"`
	UserId uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UniqueConstraint (review_id, `db:"unique_constraint" json:"unique_constraint"`
}

// NewReviewHelpful creates a new ReviewHelpful instance
func NewReviewHelpful() *ReviewHelpful {
	return &ReviewHelpful{}
}
