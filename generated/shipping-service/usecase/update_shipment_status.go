package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdateShipmentStatusInput represents input for UpdateShipmentStatus
type UpdateShipmentStatusInput struct {
	ShipmentId uuid.UUID
	NewStatus ShipmentStatus
	Location string
	Description string
	UpdatedBy string
}

// UpdateShipmentStatusUsecase defines the interface for UpdateShipmentStatus
type UpdateShipmentStatusUsecase interface {
	Execute(ctx context.Context, input UpdateShipmentStatusInput) error
}

type update_shipment_statusUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdateShipmentStatusUsecase creates a new instance
func NewUpdateShipmentStatusUsecase() UpdateShipmentStatusUsecase {
	return &update_shipment_statusUsecaseImpl{}
}

// Execute executes UpdateShipmentStatus
func (u *update_shipment_statusUsecaseImpl) Execute(ctx context.Context, input UpdateShipmentStatusInput) error {
	// TODO: Implement business logic

	return nil
}
