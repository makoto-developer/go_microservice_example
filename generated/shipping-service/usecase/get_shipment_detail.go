package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetShipmentDetailInput represents input for GetShipmentDetail
type GetShipmentDetailInput struct {
	ShipmentId uuid.UUID
}

// GetShipmentDetailUsecase defines the interface for GetShipmentDetail
type GetShipmentDetailUsecase interface {
	Execute(ctx context.Context, input GetShipmentDetailInput) error
}

type get_shipment_detailUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetShipmentDetailUsecase creates a new instance
func NewGetShipmentDetailUsecase() GetShipmentDetailUsecase {
	return &get_shipment_detailUsecaseImpl{}
}

// Execute executes GetShipmentDetail
func (u *get_shipment_detailUsecaseImpl) Execute(ctx context.Context, input GetShipmentDetailInput) error {
	// TODO: Implement business logic

	return nil
}
