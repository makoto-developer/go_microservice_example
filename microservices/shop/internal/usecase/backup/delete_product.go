package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type DeleteProductInput struct {
	ProductID uuid.UUID
}

type DeleteProductOutput struct {
	ProductID uuid.UUID
	Deleted   bool
}

type DeleteProductUsecase interface {
	Execute(ctx context.Context, input DeleteProductInput) (DeleteProductOutput, error)
}

type deleteProductUsecase struct {
	productRepo repository.ProductRepository
}

func NewDeleteProductUsecase(productRepo repository.ProductRepository) DeleteProductUsecase {
	return &deleteProductUsecase{productRepo: productRepo}
}

func (u *deleteProductUsecase) Execute(ctx context.Context, input DeleteProductInput) (DeleteProductOutput, error) {
	if _, err := u.productRepo.GetByID(ctx, input.ProductID); err != nil {
		return DeleteProductOutput{}, err
	}

	if err := u.productRepo.Delete(ctx, input.ProductID); err != nil {
		return DeleteProductOutput{}, err
	}

	return DeleteProductOutput{
		ProductID: input.ProductID,
		Deleted:   true,
	}, nil
}
