package domain

import (
	"context"

	"github.com/google/uuid"
)

// RefundRepository defines repository interface for Refund
type RefundRepository interface {
	// Create creates a new Refund
	Create(ctx context.Context, refund *Refund) error

	// FindByID finds Refund by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Refund, error)

	// Update updates Refund
	Update(ctx context.Context, refund *Refund) error

	// Delete deletes Refund by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Refund
	List(ctx context.Context, limit, offset int) ([]*Refund, error)
}
