package domain

import (
	"context"

	"github.com/google/uuid"
)

// ShopCategoryRepository defines repository interface for ShopCategory
type ShopCategoryRepository interface {
	// Create creates a new ShopCategory
	Create(ctx context.Context, shop_category *ShopCategory) error

	// FindByID finds ShopCategory by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ShopCategory, error)

	// Update updates ShopCategory
	Update(ctx context.Context, shop_category *ShopCategory) error

	// Delete deletes ShopCategory by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ShopCategory
	List(ctx context.Context, limit, offset int) ([]*ShopCategory, error)
}
