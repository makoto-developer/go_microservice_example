package domain

import (
	"context"

	"github.com/google/uuid"
)

// OrderStatusHistoryRepository defines repository interface for OrderStatusHistory
type OrderStatusHistoryRepository interface {
	// Create creates a new OrderStatusHistory
	Create(ctx context.Context, order_status_history *OrderStatusHistory) error

	// FindByID finds OrderStatusHistory by ID
	FindByID(ctx context.Context, id uuid.UUID) (*OrderStatusHistory, error)

	// Update updates OrderStatusHistory
	Update(ctx context.Context, order_status_history *OrderStatusHistory) error

	// Delete deletes OrderStatusHistory by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all OrderStatusHistory
	List(ctx context.Context, limit, offset int) ([]*OrderStatusHistory, error)
}
