package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
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
	Execute(ctx context.Context, input RegisterShopInput) (*RegisterShopOutput, error)
}

type registerShopUsecaseImpl struct {
	shopRepo         domain.ShopRepository
	shopCategoryRepo domain.ShopCategoryRepository
}

func NewRegisterShopUsecase(
	shopRepo domain.ShopRepository,
	shopCategoryRepo domain.ShopCategoryRepository,
) RegisterShopUsecase {
	return &registerShopUsecaseImpl{
		shopRepo:         shopRepo,
		shopCategoryRepo: shopCategoryRepo,
	}
}

func (u *registerShopUsecaseImpl) Execute(ctx context.Context, input RegisterShopInput) (*RegisterShopOutput, error) {
	shop := &domain.Shop{
		Id:            uuid.New(),
		OwnerId:       input.OwnerID,
		Name:          input.Name,
		Description:   input.Description,
		LogoUrl:       input.LogoURL,
		OwnerName:     input.OwnerName,
		PhoneNumber:   input.PhoneNumber,
		BusinessHours: input.BusinessHours,
		ReturnPolicy:  input.ReturnPolicy,
		Status:        domain.ShopStatusPendingApproval,
		Published:     false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := u.shopRepo.Create(ctx, shop); err != nil {
		return nil, fmt.Errorf("failed to create shop: %w", err)
	}

	for _, categoryName := range input.Categories {
		category := &domain.ShopCategory{
			Id:           uuid.New(),
			ShopId:       shop.Id,
			CategoryName: categoryName,
			CreatedAt:    time.Now(),
		}
		if err := u.shopCategoryRepo.Create(ctx, category); err != nil {
			return nil, fmt.Errorf("failed to create category: %w", err)
		}
	}

	return &RegisterShopOutput{
		ShopID:  shop.Id,
		Status:  shop.Status,
		Message: "Shop registration submitted for approval",
	}, nil
}
