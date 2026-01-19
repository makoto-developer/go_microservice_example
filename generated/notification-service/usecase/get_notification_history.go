package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetNotificationHistoryInput represents input for GetNotificationHistory
type GetNotificationHistoryInput struct {
	UserId uuid.UUID
	Channel NotificationChannel
	Status NotificationStatus
	DateFrom date
	DateTo date
	Page int
	PageSize int
}

// GetNotificationHistoryUsecase defines the interface for GetNotificationHistory
type GetNotificationHistoryUsecase interface {
	Execute(ctx context.Context, input GetNotificationHistoryInput) error
}

type get_notification_historyUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetNotificationHistoryUsecase creates a new instance
func NewGetNotificationHistoryUsecase() GetNotificationHistoryUsecase {
	return &get_notification_historyUsecaseImpl{}
}

// Execute executes GetNotificationHistory
func (u *get_notification_historyUsecaseImpl) Execute(ctx context.Context, input GetNotificationHistoryInput) error {
	// TODO: Implement business logic

	return nil
}
