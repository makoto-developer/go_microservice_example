package domain

import (
	"context"

	"github.com/google/uuid"
)

// ShopIndexRepository defines repository interface for ShopIndex
type ShopIndexRepository interface {
	// Create creates a new ShopIndex
	Create(ctx context.Context, shop_index *ShopIndex) error

	// FindByID finds ShopIndex by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ShopIndex, error)

	// Update updates ShopIndex
	Update(ctx context.Context, shop_index *ShopIndex) error

	// Delete deletes ShopIndex by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ShopIndex
	List(ctx context.Context, limit, offset int) ([]*ShopIndex, error)

	// FindByShopId finds ShopIndex by shop_id
	FindByShopId(ctx context.Context, shop_id uuid.UUID) (*ShopIndex, error)
}
