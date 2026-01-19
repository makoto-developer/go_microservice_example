package domain

import (
	"context"

	"github.com/google/uuid"
)

// OrderRepository defines repository interface for Order
type OrderRepository interface {
	// Create creates a new Order
	Create(ctx context.Context, order *Order) error

	// FindByID finds Order by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Order, error)

	// Update updates Order
	Update(ctx context.Context, order *Order) error

	// Delete deletes Order by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Order
	List(ctx context.Context, limit, offset int) ([]*Order, error)

	// FindByOrderNumber finds Order by order_number
	FindByOrderNumber(ctx context.Context, order_number string) (*Order, error)
}
