package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type UpdateProductInput struct {
	ProductID     uuid.UUID
	Name          string
	Description   string
	Price         float64
	Category      string
	StockQuantity int
	Weight        *float64
	Size          *string
	JANCode       *string
}

type UpdateProductOutput struct {
	ProductID uuid.UUID
}

type UpdateProductUsecase interface {
	Execute(ctx context.Context, input UpdateProductInput) (UpdateProductOutput, error)
}

type updateProductUsecase struct {
	productRepo repository.ProductRepository
}

func NewUpdateProductUsecase(productRepo repository.ProductRepository) UpdateProductUsecase {
	return &updateProductUsecase{productRepo: productRepo}
}

func (u *updateProductUsecase) Execute(ctx context.Context, input UpdateProductInput) (UpdateProductOutput, error) {
	product, err := u.productRepo.GetByID(ctx, input.ProductID)
	if err != nil {
		return UpdateProductOutput{}, err
	}

	if input.Name == "" || input.Price <= 0 {
		return UpdateProductOutput{}, domain.ErrInvalidProductData
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Category = input.Category
	product.StockQuantity = input.StockQuantity
	product.Weight = input.Weight
	product.Size = input.Size
	product.JANCode = input.JANCode
	product.UpdatedAt = time.Now()

	if err := u.productRepo.Update(ctx, product); err != nil {
		return UpdateProductOutput{}, err
	}

	return UpdateProductOutput{ProductID: product.ID}, nil
}
