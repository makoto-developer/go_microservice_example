package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
	"github.com/shopspring/decimal"
)

type ProductVariationInput struct {
	SKU            string
	AttributeName  string
	AttributeValue string
	Price          decimal.Decimal
	StockQuantity  int
}

type ManageProductVariationInput struct {
	ProductID  uuid.UUID
	Variations []ProductVariationInput
}

type ManageProductVariationOutput struct {
	ProductID uuid.UUID
	Message   string
}

type ManageProductVariationUsecase interface {
	Execute(ctx context.Context, input ManageProductVariationInput) (*ManageProductVariationOutput, error)
}

type manageProductVariationUsecaseImpl struct {
	productVariationRepo domain.ProductVariationRepository
}

func NewManageProductVariationUsecase(productVariationRepo domain.ProductVariationRepository) ManageProductVariationUsecase {
	return &manageProductVariationUsecaseImpl{
		productVariationRepo: productVariationRepo,
	}
}

func (u *manageProductVariationUsecaseImpl) Execute(ctx context.Context, input ManageProductVariationInput) (*ManageProductVariationOutput, error) {
	for _, varInput := range input.Variations {
		variation := &domain.ProductVariation{
			Id:             uuid.New(),
			ProductId:      input.ProductID,
			Sku:            varInput.SKU,
			AttributeName:  varInput.AttributeName,
			AttributeValue: varInput.AttributeValue,
			Price:          varInput.Price,
			StockQuantity:  varInput.StockQuantity,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := u.productVariationRepo.Create(ctx, variation); err != nil {
			return nil, fmt.Errorf("failed to create product variation: %w", err)
		}
	}

	return &ManageProductVariationOutput{
		ProductID: input.ProductID,
		Message:   "Product variations managed successfully",
	}, nil
}
