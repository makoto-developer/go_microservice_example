package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type UpdateShopInfoInput struct {
	ShopID        uuid.UUID
	Name          *string
	Description   *string
	LogoURL       *string
	PhoneNumber   *string
	BusinessHours *string
	ReturnPolicy  *string
}

type UpdateShopInfoOutput struct {
	ShopID  uuid.UUID
	Message string
}

type UpdateShopInfoUsecase interface {
	Execute(ctx context.Context, input UpdateShopInfoInput) (*UpdateShopInfoOutput, error)
}

type updateShopInfoUsecaseImpl struct {
	shopRepo domain.ShopRepository
}

func NewUpdateShopInfoUsecase(shopRepo domain.ShopRepository) UpdateShopInfoUsecase {
	return &updateShopInfoUsecaseImpl{
		shopRepo: shopRepo,
	}
}

func (u *updateShopInfoUsecaseImpl) Execute(ctx context.Context, input UpdateShopInfoInput) (*UpdateShopInfoOutput, error) {
	shop, err := u.shopRepo.FindByID(ctx, input.ShopID)
	if err != nil {
		return nil, fmt.Errorf("failed to find shop: %w", err)
	}

	if input.Name != nil {
		shop.Name = *input.Name
	}
	if input.Description != nil {
		shop.Description = *input.Description
	}
	if input.LogoURL != nil {
		shop.LogoUrl = input.LogoURL
	}
	if input.PhoneNumber != nil {
		shop.PhoneNumber = *input.PhoneNumber
	}
	if input.BusinessHours != nil {
		shop.BusinessHours = *input.BusinessHours
	}
	if input.ReturnPolicy != nil {
		shop.ReturnPolicy = *input.ReturnPolicy
	}

	shop.UpdatedAt = time.Now()

	if err := u.shopRepo.Update(ctx, shop); err != nil {
		return nil, fmt.Errorf("failed to update shop: %w", err)
	}

	return &UpdateShopInfoOutput{
		ShopID:  shop.Id,
		Message: "Shop information updated successfully",
	}, nil
}
