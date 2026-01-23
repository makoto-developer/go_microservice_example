package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory/internal/repository"
)

type ReleaseStockInput struct {
	OrderID uuid.UUID
}

type ReleaseStockOutput struct {
	Released bool
}

type ReleaseStockUsecase interface {
	Execute(ctx context.Context, input ReleaseStockInput) (ReleaseStockOutput, error)
}

type releaseStockUsecaseImpl struct {
	inventoryRepo repository.InventoryRepository
}

func NewReleaseStockUsecase(inventoryRepo repository.InventoryRepository) ReleaseStockUsecase {
	return &releaseStockUsecaseImpl{
		inventoryRepo: inventoryRepo,
	}
}

func (u *releaseStockUsecaseImpl) Execute(ctx context.Context, input ReleaseStockInput) (ReleaseStockOutput, error) {
	err := u.inventoryRepo.Release(ctx, input.OrderID)
	if err != nil {
		return ReleaseStockOutput{}, err
	}

	return ReleaseStockOutput{Released: true}, nil
}
