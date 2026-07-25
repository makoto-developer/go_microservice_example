package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/repository"
	"time"
)

type ReserveInventoryInput struct {
	ProductID   uuid.UUID
	VariationID *uuid.UUID
	OrderID     uuid.UUID
	Quantity    int
}

type InventoryManagementUsecase interface {
	ReserveInventory(ctx context.Context, input ReserveInventoryInput) (uuid.UUID, error)
	ReleaseInventory(ctx context.Context, orderID uuid.UUID) error
	ConfirmInventory(ctx context.Context, orderID uuid.UUID) error
	UpdateInventoryQuantity(ctx context.Context, productID uuid.UUID, variationID *uuid.UUID, quantity int) error
}

type inventoryManagementUsecase struct {
	inventoryRepo   repository.InventoryRepository
	reservationRepo repository.ReservationRepository
}

func NewInventoryManagementUsecase(
	inventoryRepo repository.InventoryRepository,
	reservationRepo repository.ReservationRepository,
) InventoryManagementUsecase {
	return &inventoryManagementUsecase{
		inventoryRepo:   inventoryRepo,
		reservationRepo: reservationRepo,
	}
}

func (u *inventoryManagementUsecase) ReserveInventory(ctx context.Context, input ReserveInventoryInput) (uuid.UUID, error) {
	inventory, err := u.inventoryRepo.GetByProductID(ctx, input.ProductID, input.VariationID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("inventory not found: %w", err)
	}

	if err := u.inventoryRepo.Reserve(ctx, inventory.ID, input.Quantity); err != nil {
		return uuid.Nil, fmt.Errorf("failed to reserve: %w", err)
	}

	reservation := &domain.Reservation{
		ID:          uuid.New(),
		InventoryID: inventory.ID,
		OrderID:     input.OrderID,
		Quantity:    input.Quantity,
		Status:      domain.ReservationStatusPending,
		ExpiresAt:   time.Now().Add(30 * time.Minute),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.reservationRepo.Create(ctx, reservation); err != nil {
		return uuid.Nil, fmt.Errorf("failed to create reservation: %w", err)
	}

	return reservation.ID, nil
}

func (u *inventoryManagementUsecase) ReleaseInventory(ctx context.Context, orderID uuid.UUID) error {
	reservations, err := u.reservationRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get reservations: %w", err)
	}

	for _, res := range reservations {
		if err := u.inventoryRepo.Release(ctx, res.InventoryID, res.Quantity); err != nil {
			return fmt.Errorf("failed to release: %w", err)
		}

		if err := u.reservationRepo.UpdateStatus(ctx, res.ID, domain.ReservationStatusReleased); err != nil {
			return fmt.Errorf("failed to update status: %w", err)
		}
	}

	return nil
}

func (u *inventoryManagementUsecase) ConfirmInventory(ctx context.Context, orderID uuid.UUID) error {
	reservations, err := u.reservationRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get reservations: %w", err)
	}

	for _, res := range reservations {
		if err := u.reservationRepo.UpdateStatus(ctx, res.ID, domain.ReservationStatusConfirmed); err != nil {
			return fmt.Errorf("failed to confirm: %w", err)
		}
	}

	return nil
}

func (u *inventoryManagementUsecase) UpdateInventoryQuantity(ctx context.Context, productID uuid.UUID, variationID *uuid.UUID, quantity int) error {
	inventory, err := u.inventoryRepo.GetByProductID(ctx, productID, variationID)
	if err != nil {
		return fmt.Errorf("inventory not found: %w", err)
	}

	if err := u.inventoryRepo.UpdateQuantity(ctx, inventory.ID, quantity); err != nil {
		return fmt.Errorf("failed to update quantity: %w", err)
	}

	return nil
}
