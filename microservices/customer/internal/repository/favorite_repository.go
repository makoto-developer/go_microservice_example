package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
)

type FavoriteRepository interface {
	Add(ctx context.Context, favorite *domain.Favorite) error
	Remove(ctx context.Context, customerID, productID uuid.UUID) error
	List(ctx context.Context, customerID uuid.UUID) ([]*domain.Favorite, error)
	Exists(ctx context.Context, customerID, productID uuid.UUID) (bool, error)
}
