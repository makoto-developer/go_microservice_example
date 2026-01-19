package domain

import (
	"context"

	"github.com/google/uuid"
)

// CustomerRepository defines repository interface for Customer
type CustomerRepository interface {
	// Create creates a new Customer
	Create(ctx context.Context, customer *Customer) error

	// FindByID finds Customer by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Customer, error)

	// Update updates Customer
	Update(ctx context.Context, customer *Customer) error

	// Delete deletes Customer by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Customer
	List(ctx context.Context, limit, offset int) ([]*Customer, error)
}
