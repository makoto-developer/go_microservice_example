package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/review/internal/domain"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *domain.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*domain.Review, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Review, error)
	Update(ctx context.Context, review *domain.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
}
