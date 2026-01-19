package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type ListProductsInput struct {
	ShopID        uuid.UUID
	Category      string
	PublishedOnly bool
	Limit         int
	Offset        int
}

type ListProductsOutput struct {
	Products []*domain.Product
}

type ListProductsUsecase interface {
	Execute(ctx context.Context, input ListProductsInput) (*ListProductsOutput, error)
}

type listProductsUsecaseImpl struct {
	productRepo domain.ProductRepository
}

func NewListProductsUsecase(
	productRepo domain.ProductRepository,
) ListProductsUsecase {
	return &listProductsUsecaseImpl{
		productRepo: productRepo,
	}
}

func (u *listProductsUsecaseImpl) Execute(ctx context.Context, input ListProductsInput) (*ListProductsOutput, error) {
	products, err := u.productRepo.ListByShopID(
		ctx,
		input.ShopID,
		input.Category,
		input.PublishedOnly,
		input.Limit,
		input.Offset,
	)
	if err != nil {
		return nil, err
	}

	return &ListProductsOutput{
		Products: products,
	}, nil
}
