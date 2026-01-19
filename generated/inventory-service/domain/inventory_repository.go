package domain

import (
	"context"

	"github.com/google/uuid"
)

// InventoryRepository defines repository interface for Inventory
type InventoryRepository interface {
	// Create creates a new Inventory
	Create(ctx context.Context, inventory *Inventory) error

	// FindByID finds Inventory by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Inventory, error)

	// Update updates Inventory
	Update(ctx context.Context, inventory *Inventory) error

	// Delete deletes Inventory by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Inventory
	List(ctx context.Context, limit, offset int) ([]*Inventory, error)
}
