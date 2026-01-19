package domain

import (
	"context"

	"github.com/google/uuid"
)

// FavoriteRepository defines repository interface for Favorite
type FavoriteRepository interface {
	// Create creates a new Favorite
	Create(ctx context.Context, favorite *Favorite) error

	// FindByID finds Favorite by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Favorite, error)

	// Update updates Favorite
	Update(ctx context.Context, favorite *Favorite) error

	// Delete deletes Favorite by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Favorite
	List(ctx context.Context, limit, offset int) ([]*Favorite, error)
}
