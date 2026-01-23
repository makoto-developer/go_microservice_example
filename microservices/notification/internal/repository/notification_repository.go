package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/notification/internal/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Notification, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.NotificationStatus) error
}
