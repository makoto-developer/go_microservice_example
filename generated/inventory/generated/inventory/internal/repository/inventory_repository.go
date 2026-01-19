package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory/internal/domain"
)

type InventoryRepository interface {
	// GetByID retrieves inventory by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error)

	// GetByProductAndShop retrieves inventory by product and shop
	GetByProductAndShop(ctx context.Context, productID, shopID uuid.UUID) (*domain.Inventory, error)

	// Create creates a new inventory record
	Create(ctx context.Context, inventory *domain.Inventory) error

	// Update updates an existing inventory
	Update(ctx context.Context, inventory *domain.Inventory) error

	// UpdateQuantity updates quantity and reserved_quantity atomically
	UpdateQuantity(ctx context.Context, id uuid.UUID, quantity, reservedQuantity int) error

	// Reserve reserves stock for an order
	Reserve(ctx context.Context, inventoryID uuid.UUID, quantity int, orderID uuid.UUID) error

	// Release releases reserved stock
	Release(ctx context.Context, orderID uuid.UUID) error
}
