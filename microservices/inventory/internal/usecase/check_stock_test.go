package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/usecase"
)

type mockInventoryRepository struct {
	getByProductIDFunc func(ctx context.Context, productID uuid.UUID, variationID *uuid.UUID) (*domain.Inventory, error)
}

func (m *mockInventoryRepository) GetByProductID(ctx context.Context, productID uuid.UUID, variationID *uuid.UUID) (*domain.Inventory, error) {
	if m.getByProductIDFunc != nil {
		return m.getByProductIDFunc(ctx, productID, variationID)
	}
	return nil, nil
}

func (m *mockInventoryRepository) Create(ctx context.Context, inventory *domain.Inventory) error {
	return nil
}

func (m *mockInventoryRepository) Update(ctx context.Context, inventory *domain.Inventory) error {
	return nil
}

func (m *mockInventoryRepository) UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error {
	return nil
}

func (m *mockInventoryRepository) Reserve(ctx context.Context, id uuid.UUID, quantity int) error {
	return nil
}

func (m *mockInventoryRepository) Release(ctx context.Context, id uuid.UUID, quantity int) error {
	return nil
}

func TestCheckStockUsecase_Success(t *testing.T) {
	productID := uuid.New()
	shopID := uuid.New()

	repo := &mockInventoryRepository{
		getByProductIDFunc: func(ctx context.Context, pid uuid.UUID, _ *uuid.UUID) (*domain.Inventory, error) {
			return &domain.Inventory{
				ID:               uuid.New(),
				ProductID:        productID,
				ShopID:           shopID,
				Quantity:         100,
				ReservedQuantity: 20,
			}, nil
		},
	}

	uc := usecase.NewCheckStockUsecase(repo)

	input := usecase.CheckStockInput{
		ProductID: productID,
		ShopID:    shopID,
		Quantity:  50,
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !output.Available {
		t.Error("expected stock to be available")
	}

	if output.CurrentQuantity != 100 {
		t.Errorf("expected current quantity 100, got %d", output.CurrentQuantity)
	}

	if output.AvailableQuantity != 80 {
		t.Errorf("expected available quantity 80, got %d", output.AvailableQuantity)
	}
}

func TestCheckStockUsecase_InsufficientStock(t *testing.T) {
	productID := uuid.New()
	shopID := uuid.New()

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

	uc := usecase.NewCheckStockUsecase(repo)

	input := usecase.CheckStockInput{
		ProductID: productID,
		ShopID:    shopID,
		Quantity:  50,
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Available {
		t.Error("expected stock to be unavailable")
	}

	if output.AvailableQuantity != 20 {
		t.Errorf("expected available quantity 20, got %d", output.AvailableQuantity)
	}
}
