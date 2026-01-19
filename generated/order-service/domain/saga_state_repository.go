package domain

import (
	"context"

	"github.com/google/uuid"
)

// SagaStateRepository defines repository interface for SagaState
type SagaStateRepository interface {
	// Create creates a new SagaState
	Create(ctx context.Context, saga_state *SagaState) error

	// FindByID finds SagaState by ID
	FindByID(ctx context.Context, id uuid.UUID) (*SagaState, error)

	// Update updates SagaState
	Update(ctx context.Context, saga_state *SagaState) error

	// Delete deletes SagaState by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all SagaState
	List(ctx context.Context, limit, offset int) ([]*SagaState, error)
}
