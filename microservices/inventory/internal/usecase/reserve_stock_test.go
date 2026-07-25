package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/usecase"
)

func TestReserveStockUsecase_Success(t *testing.T) {
	productID := uuid.New()
	shopID := uuid.New()
	orderID := uuid.New()
	inventoryID := uuid.New()

	repo := &mockInventoryRepository{
		getByProductIDFunc: func(ctx context.Context, pid uuid.UUID, _ *uuid.UUID) (*domain.Inventory, error) {
			return &domain.Inventory{
				ID:               inventoryID,
				ProductID:        productID,
				ShopID:           shopID,
				Quantity:         100,
				ReservedQuantity: 20,
			}, nil
		},
	}

	uc := usecase.NewReserveStockUsecase(repo)

	input := usecase.ReserveStockInput{
		ProductID: productID,
		ShopID:    shopID,
		Quantity:  50,
		OrderID:   orderID,
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Reserved {
		t.Error("expected stock to be reserved")
	}
}

func TestReserveStockUsecase_InsufficientStock(t *testing.T) {
	productID := uuid.New()
	shopID := uuid.New()
	orderID := uuid.New()

	repo := &mockInventoryRepository{
		getByProductIDFunc: func(ctx context.Context, pid uuid.UUID, _ *uuid.UUID) (*domain.Inventory, error) {
			return &domain.Inventory{
				ID:               uuid.New(),
				ProductID:        productID,
				ShopID:           shopID,
				Quantity:         100,
				ReservedQuantity: 80,
			}, nil
		},
	}

	uc := usecase.NewReserveStockUsecase(repo)

	input := usecase.ReserveStockInput{
		ProductID: productID,
		ShopID:    shopID,
		Quantity:  50,
		OrderID:   orderID,
	}

	_, err := uc.Execute(context.Background(), input)

	if err != usecase.ErrInsufficientStock {
		t.Errorf("expected ErrInsufficientStock, got %v", err)
	}
}
