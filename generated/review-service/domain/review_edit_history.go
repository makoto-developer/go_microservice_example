package domain

import (
	"github.com/google/uuid"
	"time"
)

// ReviewEditHistory represents ReviewEditHistory
type ReviewEditHistory struct {
	Id uuid.UUID `db:"id" json:"id"`
	ReviewId uuid.UUID `db:"review_id" json:"review_id"`
	OldRating int `db:"old_rating" json:"old_rating"`
	NewRating int `db:"new_rating" json:"new_rating"`
	OldTitle string `db:"old_title" json:"old_title"`
	NewTitle string `db:"new_title" json:"new_title"`
	OldContent text `db:"old_content" json:"old_content"`
	NewContent text `db:"new_content" json:"new_content"`
	EditedAt time.Time `db:"edited_at" json:"edited_at"`
}

// NewReviewEditHistory creates a new ReviewEditHistory instance
func NewReviewEditHistory() *ReviewEditHistory {
	return &ReviewEditHistory{}
}
