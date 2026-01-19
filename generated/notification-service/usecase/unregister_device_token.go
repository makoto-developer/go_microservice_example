package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UnregisterDeviceTokenInput represents input for UnregisterDeviceToken
type UnregisterDeviceTokenInput struct {
	DeviceId string
}

// UnregisterDeviceTokenUsecase defines the interface for UnregisterDeviceToken
type UnregisterDeviceTokenUsecase interface {
	Execute(ctx context.Context, input UnregisterDeviceTokenInput) error
}

type unregister_device_tokenUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUnregisterDeviceTokenUsecase creates a new instance
func NewUnregisterDeviceTokenUsecase() UnregisterDeviceTokenUsecase {
	return &unregister_device_tokenUsecaseImpl{}
}

// Execute executes UnregisterDeviceToken
func (u *unregister_device_tokenUsecaseImpl) Execute(ctx context.Context, input UnregisterDeviceTokenInput) error {
	// TODO: Implement business logic

	return nil
}
