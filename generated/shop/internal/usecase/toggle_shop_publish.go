package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type ToggleShopPublishInput struct {
	ShopID    uuid.UUID
	Published bool
}

type ToggleShopPublishOutput struct {
	ShopID    uuid.UUID
	Published bool
}

type ToggleShopPublishUsecase interface {
	Execute(ctx context.Context, input ToggleShopPublishInput) (ToggleShopPublishOutput, error)
}

type toggleShopPublishUsecase struct {
	shopRepo repository.ShopRepository
}

func NewToggleShopPublishUsecase(shopRepo repository.ShopRepository) ToggleShopPublishUsecase {
	return &toggleShopPublishUsecase{shopRepo: shopRepo}
}

func (u *toggleShopPublishUsecase) Execute(ctx context.Context, input ToggleShopPublishInput) (ToggleShopPublishOutput, error) {
	shop, err := u.shopRepo.GetByID(ctx, input.ShopID)
	if err != nil {
		return ToggleShopPublishOutput{}, err
	}

	if input.Published && !shop.CanPublish() {
		return ToggleShopPublishOutput{}, domain.ErrShopNotApproved
	}

	if err := u.shopRepo.UpdatePublished(ctx, shop.ID, input.Published); err != nil {
		return ToggleShopPublishOutput{}, err
	}

	return ToggleShopPublishOutput{
		ShopID:    shop.ID,
		Published: input.Published,
	}, nil
}
