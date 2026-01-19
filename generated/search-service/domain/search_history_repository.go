package domain

import (
	"context"

	"github.com/google/uuid"
)

// SearchHistoryRepository defines repository interface for SearchHistory
type SearchHistoryRepository interface {
	// Create creates a new SearchHistory
	Create(ctx context.Context, search_history *SearchHistory) error

	// FindByID finds SearchHistory by ID
	FindByID(ctx context.Context, id uuid.UUID) (*SearchHistory, error)

	// Update updates SearchHistory
	Update(ctx context.Context, search_history *SearchHistory) error

	// Delete deletes SearchHistory by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all SearchHistory
	List(ctx context.Context, limit, offset int) ([]*SearchHistory, error)
}
