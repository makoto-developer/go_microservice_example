package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
)

type AddressRepository interface {
	Create(ctx context.Context, address *domain.Address) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error)
	List(ctx context.Context, customerID uuid.UUID) ([]*domain.Address, error)
	Update(ctx context.Context, address *domain.Address) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetDefault(ctx context.Context, customerID uuid.UUID) (*domain.Address, error)
	SetDefault(ctx context.Context, customerID, addressID uuid.UUID) error
}
