package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type RegisterShopInput struct {
	OwnerID       uuid.UUID
	Name          string
	Description   string
	LogoURL       *string
	OwnerName     string
	PhoneNumber   string
	BusinessHours string
	ReturnPolicy  string
	Categories    []string
}

type RegisterShopOutput struct {
	ShopID  uuid.UUID
	Status  domain.ShopStatus
	Message string
}

type RegisterShopUsecase interface {
	Execute(ctx context.Context, input RegisterShopInput) (RegisterShopOutput, error)
}

type registerShopUsecase struct {
	shopRepo repository.ShopRepository
}

func NewRegisterShopUsecase(shopRepo repository.ShopRepository) RegisterShopUsecase {
	return &registerShopUsecase{shopRepo: shopRepo}
}

func (u *registerShopUsecase) Execute(ctx context.Context, input RegisterShopInput) (RegisterShopOutput, error) {
	existingShop, err := u.shopRepo.GetByOwnerID(ctx, input.OwnerID)
	if err == nil && existingShop != nil {
		return RegisterShopOutput{}, domain.ErrShopAlreadyExists
	}

	if input.Name == "" || input.Description == "" {
		return RegisterShopOutput{}, domain.ErrInvalidShopData
	}

	shop := domain.NewShop(
		input.OwnerID, input.Name, input.Description, input.LogoURL,
		input.OwnerName, input.PhoneNumber, input.BusinessHours, input.ReturnPolicy,
	)

	if err := u.shopRepo.Create(ctx, shop); err != nil {
		return RegisterShopOutput{}, err
	}

	for _, categoryName := range input.Categories {
		category := domain.NewShopCategory(shop.ID, categoryName)
		if err := u.shopRepo.AddCategory(ctx, category); err != nil {
			return RegisterShopOutput{}, err
		}
	}

	return RegisterShopOutput{
		ShopID:  shop.ID,
		Status:  shop.Status,
		Message: "Shop registered successfully. Awaiting admin approval.",
	}, nil
}
