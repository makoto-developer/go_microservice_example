package domain

import (
	"github.com/google/uuid"
	"time"
)

// SearchSuggestion represents SearchSuggestion
type SearchSuggestion struct {
	Id uuid.UUID `db:"id" json:"id"`
	Keyword string `db:"keyword" json:"keyword"`
	SearchCount int `db:"search_count" json:"search_count"`
	Category *string `db:"category" json:"category,omitempty"`
	LastSearchedAt time.Time `db:"last_searched_at" json:"last_searched_at"`
}

// NewSearchSuggestion creates a new SearchSuggestion instance
func NewSearchSuggestion() *SearchSuggestion {
	return &SearchSuggestion{}
}
