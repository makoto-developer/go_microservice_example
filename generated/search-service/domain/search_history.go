package domain

import (
	"github.com/google/uuid"
	"time"
)

// SearchHistory represents SearchHistory
type SearchHistory struct {
	Id uuid.UUID `db:"id" json:"id"`
	UserId uuid.UUID `db:"user_id" json:"user_id"`
	Keyword string `db:"keyword" json:"keyword"`
	ResultCount int `db:"result_count" json:"result_count"`
	ClickedProductId *uuid.UUID `db:"clicked_product_id" json:"clicked_product_id,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewSearchHistory creates a new SearchHistory instance
func NewSearchHistory() *SearchHistory {
	return &SearchHistory{}
}
