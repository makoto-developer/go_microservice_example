package domain

import (
	"context"

	"github.com/google/uuid"
)

// ShipmentHistoryRepository defines repository interface for ShipmentHistory
type ShipmentHistoryRepository interface {
	// Create creates a new ShipmentHistory
	Create(ctx context.Context, shipment_history *ShipmentHistory) error

	// FindByID finds ShipmentHistory by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ShipmentHistory, error)

	// Update updates ShipmentHistory
	Update(ctx context.Context, shipment_history *ShipmentHistory) error

	// Delete deletes ShipmentHistory by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ShipmentHistory
	List(ctx context.Context, limit, offset int) ([]*ShipmentHistory, error)
}
