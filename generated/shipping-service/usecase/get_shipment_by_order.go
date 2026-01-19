package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetShipmentByOrderInput represents input for GetShipmentByOrder
type GetShipmentByOrderInput struct {
	OrderId uuid.UUID
}

// GetShipmentByOrderUsecase defines the interface for GetShipmentByOrder
type GetShipmentByOrderUsecase interface {
	Execute(ctx context.Context, input GetShipmentByOrderInput) error
}

type get_shipment_by_orderUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetShipmentByOrderUsecase creates a new instance
func NewGetShipmentByOrderUsecase() GetShipmentByOrderUsecase {
	return &get_shipment_by_orderUsecaseImpl{}
}

// Execute executes GetShipmentByOrder
func (u *get_shipment_by_orderUsecaseImpl) Execute(ctx context.Context, input GetShipmentByOrderInput) error {
	// TODO: Implement business logic

	return nil
}
