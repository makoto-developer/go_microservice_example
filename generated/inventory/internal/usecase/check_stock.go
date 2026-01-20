package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory/internal/repository"
)

type CheckStockInput struct {
	ProductID uuid.UUID
	ShopID    uuid.UUID
	Quantity  int
}

type CheckStockOutput struct {
	Available         bool
	CurrentQuantity   int
	AvailableQuantity int
}

type CheckStockUsecase interface {
	Execute(ctx context.Context, input CheckStockInput) (CheckStockOutput, error)
}

type checkStockUsecaseImpl struct {
	inventoryRepo repository.InventoryRepository
}

func NewCheckStockUsecase(inventoryRepo repository.InventoryRepository) CheckStockUsecase {
	return &checkStockUsecaseImpl{
		inventoryRepo: inventoryRepo,
	}
}

func (u *checkStockUsecaseImpl) Execute(ctx context.Context, input CheckStockInput) (CheckStockOutput, error) {
	inventory, err := u.inventoryRepo.GetByProductID(ctx, input.ProductID, nil)
	if err != nil {
		return CheckStockOutput{}, err
	}

	availableQty := inventory.AvailableQuantity()
	canReserve := inventory.CanReserve(input.Quantity)

	return CheckStockOutput{
		Available:         canReserve,
		CurrentQuantity:   inventory.Quantity,
		AvailableQuantity: availableQty,
	}, nil
}
