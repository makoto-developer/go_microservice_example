package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/usecase"
)

type mockOrderRepository struct {
	createFunc         func(ctx context.Context, order *domain.Order, items []*domain.OrderItem) error
	getByIDFunc        func(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	getByCustomerIDFunc func(ctx context.Context, customerID uuid.UUID) ([]*domain.Order, error)
}

func (m *mockOrderRepository) Create(ctx context.Context, order *domain.Order, items []*domain.OrderItem) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, order, items)
	}
	return nil
}

func (m *mockOrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockOrderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	return nil, nil
}

func (m *mockOrderRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Order, error) {
	if m.getByCustomerIDFunc != nil {
		return m.getByCustomerIDFunc(ctx, customerID)
	}
	return nil, nil
}

func (m *mockOrderRepository) GetItems(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error) {
	return nil, nil
}

func (m *mockOrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	return nil
}

func (m *mockOrderRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestCreateOrderUsecase_Success(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	shopID := uuid.New()

	repo := &mockOrderRepository{
		createFunc: func(ctx context.Context, order *domain.Order, items []*domain.OrderItem) error {
			if order.CustomerID != customerID {
				t.Errorf("expected customer ID %v, got %v", customerID, order.CustomerID)
			}
			if len(items) != 1 {
				t.Errorf("expected 1 item, got %d", len(items))
			}
			return nil
		},
	}

	uc := usecase.NewCreateOrderUsecase(repo)

	input := usecase.CreateOrderInput{
		CustomerID: customerID,
		Items: []usecase.OrderItemInput{
			{
				ProductID: productID,
				ShopID:    shopID,
				Quantity:  2,
				Price:     1000,
			},
		},
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.TotalAmount != 2000 {
		t.Errorf("expected total amount 2000, got %d", output.TotalAmount)
	}

	if output.OrderNumber == "" {
		t.Error("expected order number to be generated")
	}
}

func TestCreateOrderUsecase_MultipleItems(t *testing.T) {
	customerID := uuid.New()

	repo := &mockOrderRepository{
		createFunc: func(ctx context.Context, order *domain.Order, items []*domain.OrderItem) error {
			if len(items) != 3 {
				t.Errorf("expected 3 items, got %d", len(items))
			}
			return nil
		},
	}

	uc := usecase.NewCreateOrderUsecase(repo)

	input := usecase.CreateOrderInput{
		CustomerID: customerID,
		Items: []usecase.OrderItemInput{
			{ProductID: uuid.New(), ShopID: uuid.New(), Quantity: 1, Price: 1000},
			{ProductID: uuid.New(), ShopID: uuid.New(), Quantity: 2, Price: 500},
			{ProductID: uuid.New(), ShopID: uuid.New(), Quantity: 3, Price: 300},
		},
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedTotal := (1 * 1000) + (2 * 500) + (3 * 300)
	if output.TotalAmount != expectedTotal {
		t.Errorf("expected total amount %d, got %d", expectedTotal, output.TotalAmount)
	}
}
