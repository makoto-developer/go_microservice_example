package domain

import (
	"context"

	"github.com/google/uuid"
)

// ReviewRepository defines repository interface for Review
type ReviewRepository interface {
	// Create creates a new Review
	Create(ctx context.Context, review *Review) error

	// FindByID finds Review by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Review, error)

	// Update updates Review
	Update(ctx context.Context, review *Review) error

	// Delete deletes Review by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Review
	List(ctx context.Context, limit, offset int) ([]*Review, error)
}
