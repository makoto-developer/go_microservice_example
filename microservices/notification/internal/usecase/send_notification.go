package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/notification/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/notification/internal/repository"
)

type SendNotificationInput struct {
	UserID    uuid.UUID
	Type      domain.NotificationType
	Subject   string
	Message   string
	Recipient string
}

type SendNotificationOutput struct {
	NotificationID uuid.UUID
	Status         domain.NotificationStatus
}

type SendNotificationUsecase interface {
	Execute(ctx context.Context, input SendNotificationInput) (SendNotificationOutput, error)
}

type sendNotificationUsecaseImpl struct {
	notificationRepo repository.NotificationRepository
}

func NewSendNotificationUsecase(notificationRepo repository.NotificationRepository) SendNotificationUsecase {
	return &sendNotificationUsecaseImpl{
		notificationRepo: notificationRepo,
	}
}

func (u *sendNotificationUsecaseImpl) Execute(ctx context.Context, input SendNotificationInput) (SendNotificationOutput, error) {
	notification := domain.NewNotification(
		input.UserID,
		input.Type,
		input.Subject,
		input.Message,
		input.Recipient,
	)

	// Simulate sending (in real implementation, call email/SMS service)
	notification.MarkSent()

	err := u.notificationRepo.Create(ctx, notification)
	if err != nil {
		return SendNotificationOutput{}, err
	}

	return SendNotificationOutput{
		NotificationID: notification.ID,
		Status:         notification.Status,
	}, nil
}
