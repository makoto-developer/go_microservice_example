package domain

import (
	"github.com/google/uuid"
	"time"
)

// PopularKeyword represents PopularKeyword
type PopularKeyword struct {
	Id uuid.UUID `db:"id" json:"id"`
	Keyword string `db:"keyword" json:"keyword"`
	SearchCount int `db:"search_count" json:"search_count"`
	Rank int `db:"rank" json:"rank"`
	Category *string `db:"category" json:"category,omitempty"`
	PeriodType PeriodType `db:"period_type" json:"period_type"`
	PeriodDate date `db:"period_date" json:"period_date"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewPopularKeyword creates a new PopularKeyword instance
func NewPopularKeyword() *PopularKeyword {
	return &PopularKeyword{}
}
