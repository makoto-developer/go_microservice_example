package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RegisterDeviceTokenInput represents input for RegisterDeviceToken
type RegisterDeviceTokenInput struct {
	UserId uuid.UUID
	DeviceId string
	Platform Platform
	Token string
}

// RegisterDeviceTokenUsecase defines the interface for RegisterDeviceToken
type RegisterDeviceTokenUsecase interface {
	Execute(ctx context.Context, input RegisterDeviceTokenInput) error
}

type register_device_tokenUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewRegisterDeviceTokenUsecase creates a new instance
func NewRegisterDeviceTokenUsecase() RegisterDeviceTokenUsecase {
	return &register_device_tokenUsecaseImpl{}
}

// Execute executes RegisterDeviceToken
func (u *register_device_tokenUsecaseImpl) Execute(ctx context.Context, input RegisterDeviceTokenInput) error {
	// TODO: Implement business logic

	return nil
}
