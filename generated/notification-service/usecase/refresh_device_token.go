package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RefreshDeviceTokenInput represents input for RefreshDeviceToken
type RefreshDeviceTokenInput struct {
	DeviceId string
	NewToken string
}

// RefreshDeviceTokenUsecase defines the interface for RefreshDeviceToken
type RefreshDeviceTokenUsecase interface {
	Execute(ctx context.Context, input RefreshDeviceTokenInput) error
}

type refresh_device_tokenUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewRefreshDeviceTokenUsecase creates a new instance
func NewRefreshDeviceTokenUsecase() RefreshDeviceTokenUsecase {
	return &refresh_device_tokenUsecaseImpl{}
}

// Execute executes RefreshDeviceToken
func (u *refresh_device_tokenUsecaseImpl) Execute(ctx context.Context, input RefreshDeviceTokenInput) error {
	// TODO: Implement business logic

	return nil
}
