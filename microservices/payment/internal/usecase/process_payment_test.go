package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/repository"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/usecase"
)

type mockPaymentRepository struct {
	createFunc func(ctx context.Context, payment *domain.Payment) error
}

func (m *mockPaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, payment)
	}
	return nil
}

func (m *mockPaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	return nil, nil
}

func (m *mockPaymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	return nil, nil
}

func (m *mockPaymentRepository) List(ctx context.Context, filter repository.PaymentListFilter) ([]*domain.Payment, int, error) {
	return nil, 0, nil
}

func (m *mockPaymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus, transactionID string) error {
	return nil
}

func TestProcessPaymentUsecase_Success(t *testing.T) {
	orderID := uuid.New()

	repo := &mockPaymentRepository{
		createFunc: func(ctx context.Context, payment *domain.Payment) error {
			if payment.OrderID != orderID {
				t.Errorf("expected order ID %v, got %v", orderID, payment.OrderID)
			}
			if payment.Status != domain.PaymentStatusCompleted {
				t.Errorf("expected status completed, got %v", payment.Status)
			}
			return nil
		},
	}

	uc := usecase.NewProcessPaymentUsecase(repo)

	input := usecase.ProcessPaymentInput{
		OrderID:       orderID,
		Amount:        5000,
		PaymentMethod: domain.PaymentMethodCreditCard,
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Status != domain.PaymentStatusCompleted {
		t.Errorf("expected status completed, got %v", output.Status)
	}

	if output.TransactionID == "" {
		t.Error("expected transaction ID to be generated")
	}
}

func TestProcessPaymentUsecase_InvalidAmount(t *testing.T) {
	repo := &mockPaymentRepository{}
	uc := usecase.NewProcessPaymentUsecase(repo)

	input := usecase.ProcessPaymentInput{
		OrderID:       uuid.New(),
		Amount:        0,
		PaymentMethod: domain.PaymentMethodCreditCard,
	}

	_, err := uc.Execute(context.Background(), input)

	if err != domain.ErrInvalidAmount {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}
