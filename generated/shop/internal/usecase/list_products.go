package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type ListProductsInput struct {
	ShopID         uuid.UUID
	IncludeDeleted bool
}

type ListProductsOutput struct {
	Products []*domain.Product
}

type ListProductsUsecase interface {
	Execute(ctx context.Context, input ListProductsInput) (ListProductsOutput, error)
}

type listProductsUsecase struct {
	productRepo repository.ProductRepository
}

func NewListProductsUsecase(productRepo repository.ProductRepository) ListProductsUsecase {
	return &listProductsUsecase{productRepo: productRepo}
}

func (u *listProductsUsecase) Execute(ctx context.Context, input ListProductsInput) (ListProductsOutput, error) {
	products, err := u.productRepo.List(ctx, input.ShopID, input.IncludeDeleted)
	if err != nil {
		return ListProductsOutput{}, err
	}

	return ListProductsOutput{Products: products}, nil
}
