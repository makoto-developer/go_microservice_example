package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type GetProductInput struct {
	ProductID uuid.UUID
}

type GetProductOutput struct {
	Product    *domain.Product
	Images     []*domain.ProductImage
	Tags       []*domain.ProductTag
	Variations []*domain.ProductVariation
}

type GetProductUsecase interface {
	Execute(ctx context.Context, input GetProductInput) (GetProductOutput, error)
}

type getProductUsecase struct {
	productRepo repository.ProductRepository
}

func NewGetProductUsecase(productRepo repository.ProductRepository) GetProductUsecase {
	return &getProductUsecase{productRepo: productRepo}
}

func (u *getProductUsecase) Execute(ctx context.Context, input GetProductInput) (GetProductOutput, error) {
	product, err := u.productRepo.GetByID(ctx, input.ProductID)
	if err != nil {
		return GetProductOutput{}, err
	}

	images, err := u.productRepo.GetImages(ctx, input.ProductID)
	if err != nil {
		return GetProductOutput{}, err
	}

	tags, err := u.productRepo.GetTags(ctx, input.ProductID)
	if err != nil {
		return GetProductOutput{}, err
	}

	variations, err := u.productRepo.GetVariations(ctx, input.ProductID)
	if err != nil {
		return GetProductOutput{}, err
	}

	return GetProductOutput{
		Product:    product,
		Images:     images,
		Tags:       tags,
		Variations: variations,
	}, nil
}
