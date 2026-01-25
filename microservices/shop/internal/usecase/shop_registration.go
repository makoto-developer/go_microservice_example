package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

// ShopRegistrationInput represents the input for shop registration
type ShopRegistrationInput struct {
	OwnerID       uuid.UUID
	Name          string
	Description   string
	LogoURL       string
	OwnerName     string
	PhoneNumber   string
	BusinessHours string
	ReturnPolicy  string
	CategoryIDs   []uuid.UUID
}

// ShopRegistrationOutput represents the output for shop registration
type ShopRegistrationOutput struct {
	ShopID uuid.UUID
	Status domain.ShopStatus
}

// ShopRegistrationUsecase handles shop registration
type ShopRegistrationUsecase interface {
	Execute(ctx context.Context, input ShopRegistrationInput) (ShopRegistrationOutput, error)
}

type shopRegistrationUsecase struct {
	shopRepo         repository.ShopRepository
	shopCategoryRepo repository.ShopCategoryRepository
}

// NewShopRegistrationUsecase creates a new shop registration usecase
func NewShopRegistrationUsecase(
	shopRepo repository.ShopRepository,
	shopCategoryRepo repository.ShopCategoryRepository,
) ShopRegistrationUsecase {
	return &shopRegistrationUsecase{
		shopRepo:         shopRepo,
		shopCategoryRepo: shopCategoryRepo,
	}
}

func (u *shopRegistrationUsecase) Execute(ctx context.Context, input ShopRegistrationInput) (ShopRegistrationOutput, error) {
	// Validate input
	if err := u.validateInput(input); err != nil {
		return ShopRegistrationOutput{}, err
	}

	// Create shop
	shop := &domain.Shop{
		ID:            uuid.New(),
		OwnerID:       input.OwnerID,
		Name:          input.Name,
		Description:   input.Description,
		LogoURL:       input.LogoURL,
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
		return ShopRegistrationOutput{}, fmt.Errorf("failed to create shop: %w", err)
	}

	// Add categories
	for _, categoryID := range input.CategoryIDs {
		if err := u.shopCategoryRepo.AddCategory(ctx, shop.ID, categoryID); err != nil {
			return ShopRegistrationOutput{}, fmt.Errorf("failed to add category: %w", err)
		}
	}

	return ShopRegistrationOutput{
		ShopID: shop.ID,
		Status: shop.Status,
	}, nil
}

func (u *shopRegistrationUsecase) validateInput(input ShopRegistrationInput) error {
	if input.Name == "" {
		return fmt.Errorf("shop name is required")
	}
	if input.Description == "" {
		return fmt.Errorf("shop description is required")
	}
	if input.OwnerName == "" {
		return fmt.Errorf("owner name is required")
	}
	if input.PhoneNumber == "" {
		return fmt.Errorf("owner phone is required")
	}
	return nil
}
