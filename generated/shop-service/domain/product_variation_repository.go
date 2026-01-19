package domain

import (
	"context"

	"github.com/google/uuid"
)

// ProductVariationRepository defines repository interface for ProductVariation
type ProductVariationRepository interface {
	// Create creates a new ProductVariation
	Create(ctx context.Context, product_variation *ProductVariation) error

	// FindByID finds ProductVariation by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ProductVariation, error)

	// Update updates ProductVariation
	Update(ctx context.Context, product_variation *ProductVariation) error

	// Delete deletes ProductVariation by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ProductVariation
	List(ctx context.Context, limit, offset int) ([]*ProductVariation, error)

	// FindBySku finds ProductVariation by sku
	FindBySku(ctx context.Context, sku string) (*ProductVariation, error)
}
