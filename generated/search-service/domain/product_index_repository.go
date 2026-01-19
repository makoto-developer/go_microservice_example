package domain

import (
	"context"

	"github.com/google/uuid"
)

// ProductIndexRepository defines repository interface for ProductIndex
type ProductIndexRepository interface {
	// Create creates a new ProductIndex
	Create(ctx context.Context, product_index *ProductIndex) error

	// FindByID finds ProductIndex by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ProductIndex, error)

	// Update updates ProductIndex
	Update(ctx context.Context, product_index *ProductIndex) error

	// Delete deletes ProductIndex by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ProductIndex
	List(ctx context.Context, limit, offset int) ([]*ProductIndex, error)

	// FindByProductId finds ProductIndex by product_id
	FindByProductId(ctx context.Context, product_id uuid.UUID) (*ProductIndex, error)
}
