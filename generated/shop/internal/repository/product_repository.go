package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdatePublished(ctx context.Context, productID uuid.UUID, published bool) error
	List(ctx context.Context, shopID uuid.UUID, includeDeleted bool) ([]*domain.Product, error)

	AddImage(ctx context.Context, image *domain.ProductImage) error
	GetImages(ctx context.Context, productID uuid.UUID) ([]*domain.ProductImage, error)
	DeleteImage(ctx context.Context, imageID uuid.UUID) error
	CountImages(ctx context.Context, productID uuid.UUID) (int, error)

	AddTag(ctx context.Context, tag *domain.ProductTag) error
	GetTags(ctx context.Context, productID uuid.UUID) ([]*domain.ProductTag, error)
	DeleteTags(ctx context.Context, productID uuid.UUID) error

	CreateVariation(ctx context.Context, variation *domain.ProductVariation) error
	GetVariations(ctx context.Context, productID uuid.UUID) ([]*domain.ProductVariation, error)
	UpdateVariation(ctx context.Context, variation *domain.ProductVariation) error
	DeleteVariations(ctx context.Context, productID uuid.UUID) error
}
