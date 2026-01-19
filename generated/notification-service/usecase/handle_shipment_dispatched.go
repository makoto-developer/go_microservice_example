package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// HandleShipmentDispatchedInput represents input for HandleShipmentDispatched
type HandleShipmentDispatchedInput struct {
	OrderId uuid.UUID
	CustomerId uuid.UUID
	TrackingNumber string
	Carrier string
	TrackingUrl string
}

// HandleShipmentDispatchedUsecase defines the interface for HandleShipmentDispatched
type HandleShipmentDispatchedUsecase interface {
	Execute(ctx context.Context, input HandleShipmentDispatchedInput) error
}

type handle_shipment_dispatchedUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewHandleShipmentDispatchedUsecase creates a new instance
func NewHandleShipmentDispatchedUsecase() HandleShipmentDispatchedUsecase {
	return &handle_shipment_dispatchedUsecaseImpl{}
}

// Execute executes HandleShipmentDispatched
func (u *handle_shipment_dispatchedUsecaseImpl) Execute(ctx context.Context, input HandleShipmentDispatchedInput) error {
	// TODO: Implement business logic

	return nil
}
