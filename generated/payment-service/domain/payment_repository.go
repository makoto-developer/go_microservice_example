package domain

import (
	"context"

	"github.com/google/uuid"
)

// PaymentRepository defines repository interface for Payment
type PaymentRepository interface {
	// Create creates a new Payment
	Create(ctx context.Context, payment *Payment) error

	// FindByID finds Payment by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Payment, error)

	// Update updates Payment
	Update(ctx context.Context, payment *Payment) error

	// Delete deletes Payment by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Payment
	List(ctx context.Context, limit, offset int) ([]*Payment, error)
}
