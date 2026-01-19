package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type ProductVariationInput struct {
	SKU            string
	AttributeName  string
	AttributeValue string
	Price          float64
	StockQuantity  int
}

type ManageVariationInput struct {
	ProductID  uuid.UUID
	Variations []ProductVariationInput
}

type ManageVariationOutput struct {
	VariationIDs []uuid.UUID
}

type ManageVariationUsecase interface {
	Execute(ctx context.Context, input ManageVariationInput) (ManageVariationOutput, error)
}

type manageVariationUsecase struct {
	productRepo repository.ProductRepository
}

func NewManageVariationUsecase(productRepo repository.ProductRepository) ManageVariationUsecase {
	return &manageVariationUsecase{productRepo: productRepo}
}

func (u *manageVariationUsecase) Execute(ctx context.Context, input ManageVariationInput) (ManageVariationOutput, error) {
	if _, err := u.productRepo.GetByID(ctx, input.ProductID); err != nil {
		return ManageVariationOutput{}, err
	}

	if err := u.productRepo.DeleteVariations(ctx, input.ProductID); err != nil {
		return ManageVariationOutput{}, err
	}

	var variationIDs []uuid.UUID
	for _, v := range input.Variations {
		variation := domain.NewProductVariation(
			input.ProductID, v.SKU, v.AttributeName, v.AttributeValue, v.Price, v.StockQuantity,
		)
		if err := u.productRepo.CreateVariation(ctx, variation); err != nil {
			return ManageVariationOutput{}, err
		}
		variationIDs = append(variationIDs, variation.ID)
	}

	return ManageVariationOutput{VariationIDs: variationIDs}, nil
}
