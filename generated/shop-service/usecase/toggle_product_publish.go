package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type ToggleProductPublishInput struct {
	ProductID uuid.UUID
	Published bool
}

type ToggleProductPublishOutput struct {
	ProductID uuid.UUID
	Published bool
	Message   string
}

type ToggleProductPublishUsecase interface {
	Execute(ctx context.Context, input ToggleProductPublishInput) (*ToggleProductPublishOutput, error)
}

type toggleProductPublishUsecaseImpl struct {
	productRepo domain.ProductRepository
}

func NewToggleProductPublishUsecase(productRepo domain.ProductRepository) ToggleProductPublishUsecase {
	return &toggleProductPublishUsecaseImpl{
		productRepo: productRepo,
	}
}

func (u *toggleProductPublishUsecaseImpl) Execute(ctx context.Context, input ToggleProductPublishInput) (*ToggleProductPublishOutput, error) {
	product, err := u.productRepo.FindByID(ctx, input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	product.Published = input.Published
	product.UpdatedAt = time.Now()

	if err := u.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	message := "Product unpublished successfully"
	if input.Published {
		message = "Product published successfully"
	}

	return &ToggleProductPublishOutput{
		ProductID: product.Id,
		Published: product.Published,
		Message:   message,
	}, nil
}
