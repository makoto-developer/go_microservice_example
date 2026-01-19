package domain

import (
	"context"

	"github.com/google/uuid"
)

// GuestCartItemRepository defines repository interface for GuestCartItem
type GuestCartItemRepository interface {
	// Create creates a new GuestCartItem
	Create(ctx context.Context, guest_cart_item *GuestCartItem) error

	// FindByID finds GuestCartItem by ID
	FindByID(ctx context.Context, id uuid.UUID) (*GuestCartItem, error)

	// Update updates GuestCartItem
	Update(ctx context.Context, guest_cart_item *GuestCartItem) error

	// Delete deletes GuestCartItem by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all GuestCartItem
	List(ctx context.Context, limit, offset int) ([]*GuestCartItem, error)
}
