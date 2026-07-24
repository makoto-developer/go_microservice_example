package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/domain"
)

type RefundRepository interface {
	Create(ctx context.Context, refund *domain.Refund) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Refund, error)
}
