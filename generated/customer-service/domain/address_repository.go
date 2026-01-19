package domain

import (
	"context"

	"github.com/google/uuid"
)

// AddressRepository defines repository interface for Address
type AddressRepository interface {
	// Create creates a new Address
	Create(ctx context.Context, address *Address) error

	// FindByID finds Address by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Address, error)

	// Update updates Address
	Update(ctx context.Context, address *Address) error

	// Delete deletes Address by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Address
	List(ctx context.Context, limit, offset int) ([]*Address, error)
}
