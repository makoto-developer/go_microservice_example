package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory/internal/repository"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
)

type ReserveStockInput struct {
	ProductID uuid.UUID
	ShopID    uuid.UUID
	Quantity  int
	OrderID   uuid.UUID
}

type ReserveStockOutput struct {
	Reserved bool
}

type ReserveStockUsecase interface {
	Execute(ctx context.Context, input ReserveStockInput) (ReserveStockOutput, error)
}

type reserveStockUsecaseImpl struct {
	inventoryRepo repository.InventoryRepository
}

func NewReserveStockUsecase(inventoryRepo repository.InventoryRepository) ReserveStockUsecase {
	return &reserveStockUsecaseImpl{
		inventoryRepo: inventoryRepo,
	}
}

func (u *reserveStockUsecaseImpl) Execute(ctx context.Context, input ReserveStockInput) (ReserveStockOutput, error) {
	inventory, err := u.inventoryRepo.GetByProductAndShop(ctx, input.ProductID, input.ShopID)
	if err != nil {
		return ReserveStockOutput{}, err
	}

	if !inventory.CanReserve(input.Quantity) {
		return ReserveStockOutput{}, ErrInsufficientStock
	}

	err = u.inventoryRepo.Reserve(ctx, inventory.ID, input.Quantity, input.OrderID)
	if err != nil {
		return ReserveStockOutput{}, err
	}

	return ReserveStockOutput{Reserved: true}, nil
}
