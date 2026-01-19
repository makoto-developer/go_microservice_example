package domain

import (
	"context"

	"github.com/google/uuid"
)

// PaymentHistoryRepository defines repository interface for PaymentHistory
type PaymentHistoryRepository interface {
	// Create creates a new PaymentHistory
	Create(ctx context.Context, payment_history *PaymentHistory) error

	// FindByID finds PaymentHistory by ID
	FindByID(ctx context.Context, id uuid.UUID) (*PaymentHistory, error)

	// Update updates PaymentHistory
	Update(ctx context.Context, payment_history *PaymentHistory) error

	// Delete deletes PaymentHistory by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all PaymentHistory
	List(ctx context.Context, limit, offset int) ([]*PaymentHistory, error)
}
