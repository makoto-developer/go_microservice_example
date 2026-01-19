package domain

import (
	"context"

	"github.com/google/uuid"
)

// ShipmentRepository defines repository interface for Shipment
type ShipmentRepository interface {
	// Create creates a new Shipment
	Create(ctx context.Context, shipment *Shipment) error

	// FindByID finds Shipment by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Shipment, error)

	// Update updates Shipment
	Update(ctx context.Context, shipment *Shipment) error

	// Delete deletes Shipment by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Shipment
	List(ctx context.Context, limit, offset int) ([]*Shipment, error)

	// FindByOrderId finds Shipment by order_id
	FindByOrderId(ctx context.Context, order_id uuid.UUID) (*Shipment, error)
}
