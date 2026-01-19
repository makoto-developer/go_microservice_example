package domain

import (
	"context"

	"github.com/google/uuid"
)

// ShippingMethodRepository defines repository interface for ShippingMethod
type ShippingMethodRepository interface {
	// Create creates a new ShippingMethod
	Create(ctx context.Context, shipping_method *ShippingMethod) error

	// FindByID finds ShippingMethod by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ShippingMethod, error)

	// Update updates ShippingMethod
	Update(ctx context.Context, shipping_method *ShippingMethod) error

	// Delete deletes ShippingMethod by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ShippingMethod
	List(ctx context.Context, limit, offset int) ([]*ShippingMethod, error)
}
