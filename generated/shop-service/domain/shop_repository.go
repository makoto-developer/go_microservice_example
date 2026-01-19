package domain

import (
	"context"

	"github.com/google/uuid"
)

// ShopRepository defines repository interface for Shop
type ShopRepository interface {
	// Create creates a new Shop
	Create(ctx context.Context, shop *Shop) error

	// FindByID finds Shop by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Shop, error)

	// Update updates Shop
	Update(ctx context.Context, shop *Shop) error

	// Delete deletes Shop by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Shop
	List(ctx context.Context, limit, offset int) ([]*Shop, error)
}
