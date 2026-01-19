package domain

import (
	"context"

	"github.com/google/uuid"
)

// CartItemRepository defines repository interface for CartItem
type CartItemRepository interface {
	// Create creates a new CartItem
	Create(ctx context.Context, cart_item *CartItem) error

	// FindByID finds CartItem by ID
	FindByID(ctx context.Context, id uuid.UUID) (*CartItem, error)

	// Update updates CartItem
	Update(ctx context.Context, cart_item *CartItem) error

	// Delete deletes CartItem by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all CartItem
	List(ctx context.Context, limit, offset int) ([]*CartItem, error)
}
