package domain

import (
	"context"

	"github.com/google/uuid"
)

// ProductRatingRepository defines repository interface for ProductRating
type ProductRatingRepository interface {
	// Create creates a new ProductRating
	Create(ctx context.Context, product_rating *ProductRating) error

	// FindByID finds ProductRating by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ProductRating, error)

	// Update updates ProductRating
	Update(ctx context.Context, product_rating *ProductRating) error

	// Delete deletes ProductRating by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ProductRating
	List(ctx context.Context, limit, offset int) ([]*ProductRating, error)

	// FindByProductId finds ProductRating by product_id
	FindByProductId(ctx context.Context, product_id uuid.UUID) (*ProductRating, error)
}
