package domain

import (
	"context"

	"github.com/google/uuid"
)

// CategoryRepository defines repository interface for Category
type CategoryRepository interface {
	// Create creates a new Category
	Create(ctx context.Context, category *Category) error

	// FindByID finds Category by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Category, error)

	// Update updates Category
	Update(ctx context.Context, category *Category) error

	// Delete deletes Category by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Category
	List(ctx context.Context, limit, offset int) ([]*Category, error)
}
