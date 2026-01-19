package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RegisterTrackingNumberInput represents input for RegisterTrackingNumber
type RegisterTrackingNumberInput struct {
	ShipmentId uuid.UUID
	TrackingNumber string
	Carrier Carrier
}

// RegisterTrackingNumberUsecase defines the interface for RegisterTrackingNumber
type RegisterTrackingNumberUsecase interface {
	Execute(ctx context.Context, input RegisterTrackingNumberInput) error
}

type register_tracking_numberUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewRegisterTrackingNumberUsecase creates a new instance
func NewRegisterTrackingNumberUsecase() RegisterTrackingNumberUsecase {
	return &register_tracking_numberUsecaseImpl{}
}

// Execute executes RegisterTrackingNumber
func (u *register_tracking_numberUsecaseImpl) Execute(ctx context.Context, input RegisterTrackingNumberInput) error {
	// TODO: Implement business logic

	return nil
}
