package domain

import (
	"context"

	"github.com/google/uuid"
)

// SearchSuggestionRepository defines repository interface for SearchSuggestion
type SearchSuggestionRepository interface {
	// Create creates a new SearchSuggestion
	Create(ctx context.Context, search_suggestion *SearchSuggestion) error

	// FindByID finds SearchSuggestion by ID
	FindByID(ctx context.Context, id uuid.UUID) (*SearchSuggestion, error)

	// Update updates SearchSuggestion
	Update(ctx context.Context, search_suggestion *SearchSuggestion) error

	// Delete deletes SearchSuggestion by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all SearchSuggestion
	List(ctx context.Context, limit, offset int) ([]*SearchSuggestion, error)

	// FindByKeyword finds SearchSuggestion by keyword
	FindByKeyword(ctx context.Context, keyword string) (*SearchSuggestion, error)
}
