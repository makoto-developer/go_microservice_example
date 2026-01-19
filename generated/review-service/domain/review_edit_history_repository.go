package domain

import (
	"context"

	"github.com/google/uuid"
)

// ReviewEditHistoryRepository defines repository interface for ReviewEditHistory
type ReviewEditHistoryRepository interface {
	// Create creates a new ReviewEditHistory
	Create(ctx context.Context, review_edit_history *ReviewEditHistory) error

	// FindByID finds ReviewEditHistory by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ReviewEditHistory, error)

	// Update updates ReviewEditHistory
	Update(ctx context.Context, review_edit_history *ReviewEditHistory) error

	// Delete deletes ReviewEditHistory by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ReviewEditHistory
	List(ctx context.Context, limit, offset int) ([]*ReviewEditHistory, error)
}
