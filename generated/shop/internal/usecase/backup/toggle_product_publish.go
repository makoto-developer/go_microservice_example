package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type ToggleProductPublishInput struct {
	ProductID uuid.UUID
	Published bool
}

type ToggleProductPublishOutput struct {
	ProductID uuid.UUID
	Published bool
}

type ToggleProductPublishUsecase interface {
	Execute(ctx context.Context, input ToggleProductPublishInput) (ToggleProductPublishOutput, error)
}

type toggleProductPublishUsecase struct {
	productRepo repository.ProductRepository
}

func NewToggleProductPublishUsecase(productRepo repository.ProductRepository) ToggleProductPublishUsecase {
	return &toggleProductPublishUsecase{productRepo: productRepo}
}

func (u *toggleProductPublishUsecase) Execute(ctx context.Context, input ToggleProductPublishInput) (ToggleProductPublishOutput, error) {
	if _, err := u.productRepo.GetByID(ctx, input.ProductID); err != nil {
		return ToggleProductPublishOutput{}, err
	}

	if err := u.productRepo.UpdatePublished(ctx, input.ProductID, input.Published); err != nil {
		return ToggleProductPublishOutput{}, err
	}

	return ToggleProductPublishOutput{
		ProductID: input.ProductID,
		Published: input.Published,
	}, nil
}
