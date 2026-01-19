package domain

import (
	"context"

	"github.com/google/uuid"
)

// PaymentMethodRepository defines repository interface for PaymentMethod
type PaymentMethodRepository interface {
	// Create creates a new PaymentMethod
	Create(ctx context.Context, payment_method *PaymentMethod) error

	// FindByID finds PaymentMethod by ID
	FindByID(ctx context.Context, id uuid.UUID) (*PaymentMethod, error)

	// Update updates PaymentMethod
	Update(ctx context.Context, payment_method *PaymentMethod) error

	// Delete deletes PaymentMethod by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all PaymentMethod
	List(ctx context.Context, limit, offset int) ([]*PaymentMethod, error)
}
