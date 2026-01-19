package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
	"github.com/shopspring/decimal"
)

type UpdateProductInput struct {
	ProductID     uuid.UUID
	Name          *string
	Description   *string
	Price         *decimal.Decimal
	Category      *string
	StockQuantity *int
	Weight        *decimal.Decimal
	Size          *string
	JanCode       *string
}

type UpdateProductOutput struct {
	ProductID uuid.UUID
	Message   string
}

type UpdateProductUsecase interface {
	Execute(ctx context.Context, input UpdateProductInput) (*UpdateProductOutput, error)
}

type updateProductUsecaseImpl struct {
	productRepo domain.ProductRepository
}

func NewUpdateProductUsecase(productRepo domain.ProductRepository) UpdateProductUsecase {
	return &updateProductUsecaseImpl{
		productRepo: productRepo,
	}
}

func (u *updateProductUsecaseImpl) Execute(ctx context.Context, input UpdateProductInput) (*UpdateProductOutput, error) {
	product, err := u.productRepo.FindByID(ctx, input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.Category != nil {
		product.Category = *input.Category
	}
	if input.StockQuantity != nil {
		product.StockQuantity = *input.StockQuantity
	}
	if input.Weight != nil {
		product.Weight = input.Weight
	}
	if input.Size != nil {
		product.Size = input.Size
	}
	if input.JanCode != nil {
		product.JanCode = input.JanCode
	}

	product.UpdatedAt = time.Now()

	if err := u.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &UpdateProductOutput{
		ProductID: product.Id,
		Message:   "Product updated successfully",
	}, nil
}
