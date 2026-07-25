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
	GetInventory(ctx context.Context, id uuid.UUID) (*domain.Inventory, error)
	GetInventoryByProduct(ctx context.Context, productID uuid.UUID, variationID *uuid.UUID) (*domain.Inventory, error)
	RegisterInventory(ctx context.Context, productID uuid.UUID, variationID *uuid.UUID, shopID uuid.UUID, initialQuantity int) (*domain.Inventory, error)
	AdjustQuantity(ctx context.Context, inventoryID uuid.UUID, newQuantity int) (*domain.Inventory, error)
	ReleaseExpiredReservations(ctx context.Context) error
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
		// PENDING の引当のみ在庫に戻す。CONFIRMED/RELEASED/EXPIRED を解放すると
		// 二重解放になり在庫が水増しされるため、冪等性のためにスキップする。
		if res.Status != domain.ReservationStatusPending {
			continue
		}

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

// GetInventory は在庫を ID で取得する。
func (u *inventoryManagementUsecase) GetInventory(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	return u.inventoryRepo.GetByID(ctx, id)
}

// GetInventoryByProduct は商品(+バリエーション)の在庫を取得する。
func (u *inventoryManagementUsecase) GetInventoryByProduct(ctx context.Context, productID uuid.UUID, variationID *uuid.UUID) (*domain.Inventory, error) {
	return u.inventoryRepo.GetByProductID(ctx, productID, variationID)
}

// RegisterInventory は商品の在庫レコードを新規登録する。
func (u *inventoryManagementUsecase) RegisterInventory(ctx context.Context, productID uuid.UUID, variationID *uuid.UUID, shopID uuid.UUID, initialQuantity int) (*domain.Inventory, error) {
	now := time.Now()
	inv := &domain.Inventory{
		ID:          uuid.New(),
		ProductID:   productID,
		VariationID: variationID,
		ShopID:      shopID,
		Quantity:    initialQuantity,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := u.inventoryRepo.Create(ctx, inv); err != nil {
		return nil, fmt.Errorf("failed to register inventory: %w", err)
	}
	return inv, nil
}

// AdjustQuantity は在庫数を絶対値で合わせる(入荷・棚卸し・手動調整)。
func (u *inventoryManagementUsecase) AdjustQuantity(ctx context.Context, inventoryID uuid.UUID, newQuantity int) (*domain.Inventory, error) {
	if newQuantity < 0 {
		return nil, fmt.Errorf("quantity must not be negative")
	}
	if err := u.inventoryRepo.UpdateQuantity(ctx, inventoryID, newQuantity); err != nil {
		return nil, fmt.Errorf("failed to update quantity: %w", err)
	}
	return u.inventoryRepo.GetByID(ctx, inventoryID)
}

// ReleaseExpiredReservations は期限切れの引当を解放する(バッチから呼ばれる)。
// 単にステータスを EXPIRED にするだけでなく、確保していた在庫を実際に戻す。
func (u *inventoryManagementUsecase) ReleaseExpiredReservations(ctx context.Context) error {
	expired, err := u.reservationRepo.GetExpiredPending(ctx)
	if err != nil {
		return fmt.Errorf("failed to get expired reservations: %w", err)
	}

	for _, res := range expired {
		if err := u.inventoryRepo.Release(ctx, res.InventoryID, res.Quantity); err != nil {
			return fmt.Errorf("failed to release expired reservation %s: %w", res.ID, err)
		}
		if err := u.reservationRepo.UpdateStatus(ctx, res.ID, domain.ReservationStatusExpired); err != nil {
			return fmt.Errorf("failed to mark reservation %s expired: %w", res.ID, err)
		}
	}
	return nil
}
