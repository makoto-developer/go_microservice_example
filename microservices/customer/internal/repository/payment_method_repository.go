package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
)

type PaymentMethodRepository interface {
	Create(ctx context.Context, method *domain.PaymentMethod) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PaymentMethod, error)
	List(ctx context.Context, customerID uuid.UUID) ([]*domain.PaymentMethod, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SetDefault(ctx context.Context, customerID, methodID uuid.UUID) error
}
