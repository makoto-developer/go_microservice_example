package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
)

type ShopRepository interface {
	Create(ctx context.Context, shop *domain.Shop) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Shop, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.Shop, error)
	Update(ctx context.Context, shop *domain.Shop) error
	UpdateStatus(ctx context.Context, shopID uuid.UUID, status domain.ShopStatus) error
	UpdatePublished(ctx context.Context, shopID uuid.UUID, published bool) error

	AddCategory(ctx context.Context, category *domain.ShopCategory) error
	GetCategories(ctx context.Context, shopID uuid.UUID) ([]*domain.ShopCategory, error)
	DeleteCategories(ctx context.Context, shopID uuid.UUID) error
}
