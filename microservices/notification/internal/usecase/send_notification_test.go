package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/notification/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/notification/internal/usecase"
)

type mockNotificationRepository struct {
	createFunc func(ctx context.Context, notification *domain.Notification) error
}

func (m *mockNotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, notification)
	}
	return nil
}

func (m *mockNotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	return nil, nil
}

func (m *mockNotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Notification, error) {
	return nil, nil
}

func (m *mockNotificationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.NotificationStatus) error {
	return nil
}

func TestSendNotificationUsecase_Success(t *testing.T) {
	userID := uuid.New()

	repo := &mockNotificationRepository{
		createFunc: func(ctx context.Context, notification *domain.Notification) error {
			if notification.UserID != userID {
				t.Errorf("expected user ID %v, got %v", userID, notification.UserID)
			}
			if notification.Status != domain.NotificationStatusSent {
				t.Errorf("expected status sent, got %v", notification.Status)
			}
			return nil
		},
	}

	uc := usecase.NewSendNotificationUsecase(repo)

	input := usecase.SendNotificationInput{
		UserID:    userID,
		Type:      domain.NotificationTypeEmail,
		Subject:   "Order Confirmation",
		Message:   "Your order has been confirmed",
		Recipient: "user@example.com",
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Status != domain.NotificationStatusSent {
		t.Errorf("expected status sent, got %v", output.Status)
	}
}
