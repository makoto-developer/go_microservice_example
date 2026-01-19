package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Customer, error)
	Update(ctx context.Context, customer *domain.Customer) error
}
