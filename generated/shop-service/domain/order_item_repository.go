package domain

import (
	"context"

	"github.com/google/uuid"
)

// OrderItemRepository defines repository interface for OrderItem
type OrderItemRepository interface {
	// Create creates a new OrderItem
	Create(ctx context.Context, order_item *OrderItem) error

	// FindByID finds OrderItem by ID
	FindByID(ctx context.Context, id uuid.UUID) (*OrderItem, error)

	// Update updates OrderItem
	Update(ctx context.Context, order_item *OrderItem) error

	// Delete deletes OrderItem by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all OrderItem
	List(ctx context.Context, limit, offset int) ([]*OrderItem, error)
}
