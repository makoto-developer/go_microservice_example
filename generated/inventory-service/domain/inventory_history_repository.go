package domain

import (
	"context"

	"github.com/google/uuid"
)

// InventoryHistoryRepository defines repository interface for InventoryHistory
type InventoryHistoryRepository interface {
	// Create creates a new InventoryHistory
	Create(ctx context.Context, inventory_history *InventoryHistory) error

	// FindByID finds InventoryHistory by ID
	FindByID(ctx context.Context, id uuid.UUID) (*InventoryHistory, error)

	// Update updates InventoryHistory
	Update(ctx context.Context, inventory_history *InventoryHistory) error

	// Delete deletes InventoryHistory by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all InventoryHistory
	List(ctx context.Context, limit, offset int) ([]*InventoryHistory, error)
}
