package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetByShopID(ctx context.Context, shopID uuid.UUID, includeDeleted bool) ([]*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProductStatus) error
	UpdateIsPublic(ctx context.Context, id uuid.UUID, isPublic bool) error
	UpdateStock(ctx context.Context, id uuid.UUID, stockCount int) error
	List(ctx context.Context, limit, offset int) ([]*domain.Product, error)
}

// ProductVariationRepository defines the interface for product variation data access
type ProductVariationRepository interface {
	Create(ctx context.Context, variation *domain.ProductVariation) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ProductVariation, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*domain.ProductVariation, error)
	Update(ctx context.Context, variation *domain.ProductVariation) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStock(ctx context.Context, id uuid.UUID, stockCount int) error
}

// ProductImageRepository defines the interface for product image data access
type ProductImageRepository interface {
	Create(ctx context.Context, image *domain.ProductImage) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ProductImage, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*domain.ProductImage, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateDisplayOrder(ctx context.Context, id uuid.UUID, order int) error
	SetPrimary(ctx context.Context, productID, imageID uuid.UUID) error
}

// ProductCategoryRepository defines the interface for product category data access
type ProductCategoryRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ProductCategory, error)
	List(ctx context.Context) ([]*domain.ProductCategory, error)
	Create(ctx context.Context, category *domain.ProductCategory) error
}
