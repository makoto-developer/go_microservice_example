package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/repository"
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
	inventoryRepo   repository.InventoryRepository
	reservationRepo repository.ReservationRepository
}

func NewReleaseStockUsecase(
	inventoryRepo repository.InventoryRepository,
	reservationRepo repository.ReservationRepository,
) ReleaseStockUsecase {
	return &releaseStockUsecaseImpl{
		inventoryRepo:   inventoryRepo,
		reservationRepo: reservationRepo,
	}
}

func (u *releaseStockUsecaseImpl) Execute(ctx context.Context, input ReleaseStockInput) (ReleaseStockOutput, error) {
	// Get reservations for this order
	reservations, err := u.reservationRepo.GetByOrderID(ctx, input.OrderID)
	if err != nil {
		return ReleaseStockOutput{}, err
	}

	// Release each reserved inventory
	for _, reservation := range reservations {
		err := u.inventoryRepo.Release(ctx, reservation.InventoryID, reservation.Quantity)
		if err != nil {
			return ReleaseStockOutput{}, err
		}
	}

	return ReleaseStockOutput{Released: true}, nil
}
