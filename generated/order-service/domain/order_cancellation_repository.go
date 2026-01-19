package domain

import (
	"context"

	"github.com/google/uuid"
)

// OrderCancellationRepository defines repository interface for OrderCancellation
type OrderCancellationRepository interface {
	// Create creates a new OrderCancellation
	Create(ctx context.Context, order_cancellation *OrderCancellation) error

	// FindByID finds OrderCancellation by ID
	FindByID(ctx context.Context, id uuid.UUID) (*OrderCancellation, error)

	// Update updates OrderCancellation
	Update(ctx context.Context, order_cancellation *OrderCancellation) error

	// Delete deletes OrderCancellation by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all OrderCancellation
	List(ctx context.Context, limit, offset int) ([]*OrderCancellation, error)
}
