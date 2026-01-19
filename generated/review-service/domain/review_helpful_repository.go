package domain

import (
	"context"

	"github.com/google/uuid"
)

// ReviewHelpfulRepository defines repository interface for ReviewHelpful
type ReviewHelpfulRepository interface {
	// Create creates a new ReviewHelpful
	Create(ctx context.Context, review_helpful *ReviewHelpful) error

	// FindByID finds ReviewHelpful by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ReviewHelpful, error)

	// Update updates ReviewHelpful
	Update(ctx context.Context, review_helpful *ReviewHelpful) error

	// Delete deletes ReviewHelpful by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ReviewHelpful
	List(ctx context.Context, limit, offset int) ([]*ReviewHelpful, error)
}
