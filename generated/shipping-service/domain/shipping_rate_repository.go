package domain

import (
	"context"

	"github.com/google/uuid"
)

// ShippingRateRepository defines repository interface for ShippingRate
type ShippingRateRepository interface {
	// Create creates a new ShippingRate
	Create(ctx context.Context, shipping_rate *ShippingRate) error

	// FindByID finds ShippingRate by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ShippingRate, error)

	// Update updates ShippingRate
	Update(ctx context.Context, shipping_rate *ShippingRate) error

	// Delete deletes ShippingRate by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ShippingRate
	List(ctx context.Context, limit, offset int) ([]*ShippingRate, error)
}
