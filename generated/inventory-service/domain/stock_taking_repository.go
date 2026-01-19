package domain

import (
	"context"

	"github.com/google/uuid"
)

// StockTakingRepository defines repository interface for StockTaking
type StockTakingRepository interface {
	// Create creates a new StockTaking
	Create(ctx context.Context, stock_taking *StockTaking) error

	// FindByID finds StockTaking by ID
	FindByID(ctx context.Context, id uuid.UUID) (*StockTaking, error)

	// Update updates StockTaking
	Update(ctx context.Context, stock_taking *StockTaking) error

	// Delete deletes StockTaking by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all StockTaking
	List(ctx context.Context, limit, offset int) ([]*StockTaking, error)
}
