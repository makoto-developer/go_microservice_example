package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ResendNotificationInput represents input for ResendNotification
type ResendNotificationInput struct {
	NotificationId uuid.UUID
}

// ResendNotificationUsecase defines the interface for ResendNotification
type ResendNotificationUsecase interface {
	Execute(ctx context.Context, input ResendNotificationInput) error
}

type resend_notificationUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewResendNotificationUsecase creates a new instance
func NewResendNotificationUsecase() ResendNotificationUsecase {
	return &resend_notificationUsecaseImpl{}
}

// Execute executes ResendNotification
func (u *resend_notificationUsecaseImpl) Execute(ctx context.Context, input ResendNotificationInput) error {
	// TODO: Implement business logic

	return nil
}
