package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/usecase"
)

type mockCartRepository struct {
	addItemFunc      func(ctx context.Context, item *domain.CartItem) error
	getByIDFunc      func(ctx context.Context, id uuid.UUID) (*domain.CartItem, error)
	getByCustomerIDFunc func(ctx context.Context, customerID uuid.UUID) ([]*domain.CartItem, error)
}

func (m *mockCartRepository) AddItem(ctx context.Context, item *domain.CartItem) error {
	if m.addItemFunc != nil {
		return m.addItemFunc(ctx, item)
	}
	return nil
}

func (m *mockCartRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CartItem, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockCartRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.CartItem, error) {
	if m.getByCustomerIDFunc != nil {
		return m.getByCustomerIDFunc(ctx, customerID)
	}
	return nil, nil
}

func (m *mockCartRepository) GetBySessionID(ctx context.Context, sessionID string) ([]*domain.GuestCartItem, error) {
	return nil, nil
}

func (m *mockCartRepository) UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error {
	return nil
}

func (m *mockCartRepository) RemoveItem(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockCartRepository) ClearCart(ctx context.Context, customerID uuid.UUID) error {
	return nil
}

func (m *mockCartRepository) AddGuestItem(ctx context.Context, item *domain.GuestCartItem) error {
	return nil
}

func (m *mockCartRepository) ClearGuestCart(ctx context.Context, sessionID string) error {
	return nil
}

func TestAddToCartUsecase_Success(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()
	quantity := 2

	repo := &mockCartRepository{
		addItemFunc: func(ctx context.Context, item *domain.CartItem) error {
			if item.CustomerID != customerID {
				t.Errorf("expected customer ID %v, got %v", customerID, item.CustomerID)
			}
			if item.ProductID != productID {
				t.Errorf("expected product ID %v, got %v", productID, item.ProductID)
			}
			if item.Quantity != quantity {
				t.Errorf("expected quantity %v, got %v", quantity, item.Quantity)
			}
			return nil
		},
	}

	uc := usecase.NewAddToCartUsecase(repo)

	input := usecase.AddToCartInput{
		CustomerID: customerID,
		ProductID:  productID,
		Quantity:   quantity,
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.TotalQuantity != quantity {
		t.Errorf("expected total quantity %v, got %v", quantity, output.TotalQuantity)
	}
}

func TestAddToCartUsecase_InvalidQuantity(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()

	repo := &mockCartRepository{}
	uc := usecase.NewAddToCartUsecase(repo)

	input := usecase.AddToCartInput{
		CustomerID: customerID,
		ProductID:  productID,
		Quantity:   0,
	}

	_, err := uc.Execute(context.Background(), input)

	if err != domain.ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity, got %v", err)
	}
}

func TestAddToCartUsecase_NegativeQuantity(t *testing.T) {
	customerID := uuid.New()
	productID := uuid.New()

	repo := &mockCartRepository{}
	uc := usecase.NewAddToCartUsecase(repo)

	input := usecase.AddToCartInput{
		CustomerID: customerID,
		ProductID:  productID,
		Quantity:   -1,
	}

	_, err := uc.Execute(context.Background(), input)

	if err != domain.ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity, got %v", err)
	}
}
