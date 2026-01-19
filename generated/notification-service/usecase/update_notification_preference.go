package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdateNotificationPreferenceInput represents input for UpdateNotificationPreference
type UpdateNotificationPreferenceInput struct {
	UserId uuid.UUID
	EmailEnabled bool
	PushEnabled bool
	EmailOrderUpdates bool
	EmailShopUpdates bool
	EmailChatMessages bool
	PushOrderUpdates bool
	PushStockRestored bool
	PushCampaigns bool
	PushChatMessages bool
	Frequency NotificationFrequency
	QuietHoursStart time
	QuietHoursEnd time
}

// UpdateNotificationPreferenceUsecase defines the interface for UpdateNotificationPreference
type UpdateNotificationPreferenceUsecase interface {
	Execute(ctx context.Context, input UpdateNotificationPreferenceInput) error
}

type update_notification_preferenceUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdateNotificationPreferenceUsecase creates a new instance
func NewUpdateNotificationPreferenceUsecase() UpdateNotificationPreferenceUsecase {
	return &update_notification_preferenceUsecaseImpl{}
}

// Execute executes UpdateNotificationPreference
func (u *update_notification_preferenceUsecaseImpl) Execute(ctx context.Context, input UpdateNotificationPreferenceInput) error {
	// TODO: Implement business logic

	return nil
}
