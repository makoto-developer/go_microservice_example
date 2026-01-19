package domain

import (
	"context"

	"github.com/google/uuid"
)

// PopularKeywordRepository defines repository interface for PopularKeyword
type PopularKeywordRepository interface {
	// Create creates a new PopularKeyword
	Create(ctx context.Context, popular_keyword *PopularKeyword) error

	// FindByID finds PopularKeyword by ID
	FindByID(ctx context.Context, id uuid.UUID) (*PopularKeyword, error)

	// Update updates PopularKeyword
	Update(ctx context.Context, popular_keyword *PopularKeyword) error

	// Delete deletes PopularKeyword by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all PopularKeyword
	List(ctx context.Context, limit, offset int) ([]*PopularKeyword, error)

	// FindByKeyword finds PopularKeyword by keyword
	FindByKeyword(ctx context.Context, keyword string) (*PopularKeyword, error)
}
