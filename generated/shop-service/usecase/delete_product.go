package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type DeleteProductInput struct {
	ProductID uuid.UUID
}

type DeleteProductOutput struct {
	ProductID uuid.UUID
	Message   string
}

type DeleteProductUsecase interface {
	Execute(ctx context.Context, input DeleteProductInput) (*DeleteProductOutput, error)
}

type deleteProductUsecaseImpl struct {
	productRepo domain.ProductRepository
}

func NewDeleteProductUsecase(productRepo domain.ProductRepository) DeleteProductUsecase {
	return &deleteProductUsecaseImpl{
		productRepo: productRepo,
	}
}

func (u *deleteProductUsecaseImpl) Execute(ctx context.Context, input DeleteProductInput) (*DeleteProductOutput, error) {
	product, err := u.productRepo.FindByID(ctx, input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	product.Deleted = true
	product.UpdatedAt = time.Now()

	if err := u.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to delete product: %w", err)
	}

	return &DeleteProductOutput{
		ProductID: product.Id,
		Message:   "Product deleted successfully",
	}, nil
}
