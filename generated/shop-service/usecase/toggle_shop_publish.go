package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type ToggleShopPublishInput struct {
	ShopID    uuid.UUID
	Published bool
}

type ToggleShopPublishOutput struct {
	ShopID    uuid.UUID
	Published bool
	Message   string
}

type ToggleShopPublishUsecase interface {
	Execute(ctx context.Context, input ToggleShopPublishInput) (*ToggleShopPublishOutput, error)
}

type toggleShopPublishUsecaseImpl struct {
	shopRepo domain.ShopRepository
}

func NewToggleShopPublishUsecase(shopRepo domain.ShopRepository) ToggleShopPublishUsecase {
	return &toggleShopPublishUsecaseImpl{
		shopRepo: shopRepo,
	}
}

func (u *toggleShopPublishUsecaseImpl) Execute(ctx context.Context, input ToggleShopPublishInput) (*ToggleShopPublishOutput, error) {
	shop, err := u.shopRepo.FindByID(ctx, input.ShopID)
	if err != nil {
		return nil, fmt.Errorf("failed to find shop: %w", err)
	}

	shop.Published = input.Published
	shop.UpdatedAt = time.Now()

	if err := u.shopRepo.Update(ctx, shop); err != nil {
		return nil, fmt.Errorf("failed to update shop: %w", err)
	}

	message := "Shop unpublished successfully"
	if input.Published {
		message = "Shop published successfully"
	}

	return &ToggleShopPublishOutput{
		ShopID:    shop.Id,
		Published: shop.Published,
		Message:   message,
	}, nil
}
