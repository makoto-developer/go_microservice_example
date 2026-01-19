package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type GetShopInput struct {
	ShopID uuid.UUID
}

type GetShopOutput struct {
	Shop       *domain.Shop
	Categories []*domain.ShopCategory
}

type GetShopUsecase interface {
	Execute(ctx context.Context, input GetShopInput) (GetShopOutput, error)
}

type getShopUsecase struct {
	shopRepo repository.ShopRepository
}

func NewGetShopUsecase(shopRepo repository.ShopRepository) GetShopUsecase {
	return &getShopUsecase{shopRepo: shopRepo}
}

func (u *getShopUsecase) Execute(ctx context.Context, input GetShopInput) (GetShopOutput, error) {
	shop, err := u.shopRepo.GetByID(ctx, input.ShopID)
	if err != nil {
		return GetShopOutput{}, err
	}

	categories, err := u.shopRepo.GetCategories(ctx, input.ShopID)
	if err != nil {
		return GetShopOutput{}, err
	}

	return GetShopOutput{
		Shop:       shop,
		Categories: categories,
	}, nil
}
