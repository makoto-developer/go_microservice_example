package domain

import (
	"context"

	"github.com/google/uuid"
)

// ProductTagRepository defines repository interface for ProductTag
type ProductTagRepository interface {
	// Create creates a new ProductTag
	Create(ctx context.Context, product_tag *ProductTag) error

	// FindByID finds ProductTag by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ProductTag, error)

	// Update updates ProductTag
	Update(ctx context.Context, product_tag *ProductTag) error

	// Delete deletes ProductTag by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ProductTag
	List(ctx context.Context, limit, offset int) ([]*ProductTag, error)
}
