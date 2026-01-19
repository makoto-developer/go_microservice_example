package domain

import (
	"context"

	"github.com/google/uuid"
)

// ReviewImageRepository defines repository interface for ReviewImage
type ReviewImageRepository interface {
	// Create creates a new ReviewImage
	Create(ctx context.Context, review_image *ReviewImage) error

	// FindByID finds ReviewImage by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ReviewImage, error)

	// Update updates ReviewImage
	Update(ctx context.Context, review_image *ReviewImage) error

	// Delete deletes ReviewImage by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ReviewImage
	List(ctx context.Context, limit, offset int) ([]*ReviewImage, error)
}
