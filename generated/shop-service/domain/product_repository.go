package domain

import (
	"context"

	"github.com/google/uuid"
)

// ProductRepository defines repository interface for Product
type ProductRepository interface {
	// Create creates a new Product
	Create(ctx context.Context, product *Product) error

	// FindByID finds Product by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Product, error)

	// Update updates Product
	Update(ctx context.Context, product *Product) error

	// Delete deletes Product by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Product
	List(ctx context.Context, limit, offset int) ([]*Product, error)

	// ListByShopID lists Products by shop ID with filters
	ListByShopID(ctx context.Context, shopID uuid.UUID, category string, publishedOnly bool, limit, offset int) ([]*Product, error)
}
