package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type UpdateShopInput struct {
	ShopID        uuid.UUID
	Name          string
	Description   string
	LogoURL       *string
	BusinessHours string
	ReturnPolicy  string
}

type UpdateShopOutput struct {
	ShopID            uuid.UUID
	RequiresReapproval bool
}

type UpdateShopUsecase interface {
	Execute(ctx context.Context, input UpdateShopInput) (UpdateShopOutput, error)
}

type updateShopUsecase struct {
	shopRepo repository.ShopRepository
}

func NewUpdateShopUsecase(shopRepo repository.ShopRepository) UpdateShopUsecase {
	return &updateShopUsecase{shopRepo: shopRepo}
}

func (u *updateShopUsecase) Execute(ctx context.Context, input UpdateShopInput) (UpdateShopOutput, error) {
	shop, err := u.shopRepo.GetByID(ctx, input.ShopID)
	if err != nil {
		return UpdateShopOutput{}, err
	}

	requiresReapproval := shop.Name != input.Name || shop.Description != input.Description

	shop.Name = input.Name
	shop.Description = input.Description
	shop.LogoURL = input.LogoURL
	shop.BusinessHours = input.BusinessHours
	shop.ReturnPolicy = input.ReturnPolicy
	shop.UpdatedAt = time.Now()

	if err := u.shopRepo.Update(ctx, shop); err != nil {
		return UpdateShopOutput{}, err
	}

	if requiresReapproval {
		if err := u.shopRepo.UpdateStatus(ctx, shop.ID, domain.ShopStatusPendingApproval); err != nil {
			return UpdateShopOutput{}, err
		}
	}

	return UpdateShopOutput{
		ShopID:            shop.ID,
		RequiresReapproval: requiresReapproval,
	}, nil
}
