package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
)

// ShopRepository defines the interface for shop data access
type ShopRepository interface {
	Create(ctx context.Context, shop *domain.Shop) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Shop, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Shop, error)
	Update(ctx context.Context, shop *domain.Shop) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ShopStatus) error
	UpdateIsPublic(ctx context.Context, id uuid.UUID, isPublic bool) error
	List(ctx context.Context, limit, offset int) ([]*domain.Shop, error)
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	List(ctx context.Context) ([]*domain.Category, error)
	Create(ctx context.Context, category *domain.Category) error
}

// ShopCategoryRepository defines the interface for shop-category relationship
type ShopCategoryRepository interface {
	AddCategory(ctx context.Context, shopID, categoryID uuid.UUID) error
	RemoveCategory(ctx context.Context, shopID, categoryID uuid.UUID) error
	GetByShopID(ctx context.Context, shopID uuid.UUID) ([]*domain.Category, error)
}
