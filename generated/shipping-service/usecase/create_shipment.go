package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateShipmentInput represents input for CreateShipment
type CreateShipmentInput struct {
	OrderId uuid.UUID
	ShippingMethodId uuid.UUID
	ShippingAddress string
	RecipientName string
	RecipientPhone string
	DeliveryDate date
	DeliveryTimeSlot TimeSlot
	DeliveryOption DeliveryOption
}

// CreateShipmentUsecase defines the interface for CreateShipment
type CreateShipmentUsecase interface {
	Execute(ctx context.Context, input CreateShipmentInput) error
}

type create_shipmentUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewCreateShipmentUsecase creates a new instance
func NewCreateShipmentUsecase() CreateShipmentUsecase {
	return &create_shipmentUsecaseImpl{}
}

// Execute executes CreateShipment
func (u *create_shipmentUsecaseImpl) Execute(ctx context.Context, input CreateShipmentInput) error {
	// TODO: Implement business logic

	return nil
}
