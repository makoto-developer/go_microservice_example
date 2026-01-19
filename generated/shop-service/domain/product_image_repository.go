package domain

import (
	"context"

	"github.com/google/uuid"
)

// ProductImageRepository defines repository interface for ProductImage
type ProductImageRepository interface {
	// Create creates a new ProductImage
	Create(ctx context.Context, product_image *ProductImage) error

	// FindByID finds ProductImage by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ProductImage, error)

	// Update updates ProductImage
	Update(ctx context.Context, product_image *ProductImage) error

	// Delete deletes ProductImage by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ProductImage
	List(ctx context.Context, limit, offset int) ([]*ProductImage, error)
}
