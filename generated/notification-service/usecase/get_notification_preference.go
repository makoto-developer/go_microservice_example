package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetNotificationPreferenceInput represents input for GetNotificationPreference
type GetNotificationPreferenceInput struct {
	UserId uuid.UUID
}

// GetNotificationPreferenceUsecase defines the interface for GetNotificationPreference
type GetNotificationPreferenceUsecase interface {
	Execute(ctx context.Context, input GetNotificationPreferenceInput) error
}

type get_notification_preferenceUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetNotificationPreferenceUsecase creates a new instance
func NewGetNotificationPreferenceUsecase() GetNotificationPreferenceUsecase {
	return &get_notification_preferenceUsecaseImpl{}
}

// Execute executes GetNotificationPreference
func (u *get_notification_preferenceUsecaseImpl) Execute(ctx context.Context, input GetNotificationPreferenceInput) error {
	// TODO: Implement business logic

	return nil
}
