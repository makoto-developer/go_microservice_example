package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/usecase"
)

type mockShipmentRepository struct {
	createFunc func(ctx context.Context, shipment *domain.Shipment) error
}

func (m *mockShipmentRepository) Create(ctx context.Context, shipment *domain.Shipment) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, shipment)
	}
	return nil
}

func (m *mockShipmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Shipment, error) {
	return nil, nil
}

func (m *mockShipmentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Shipment, error) {
	return nil, nil
}

func (m *mockShipmentRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Shipment, error) {
	return nil, nil
}

func (m *mockShipmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ShipmentStatus) error {
	return nil
}

func (m *mockShipmentRepository) UpdateTracking(ctx context.Context, id uuid.UUID, trackingNumber string) error {
	return nil
}

func TestCreateShipmentUsecase_Success(t *testing.T) {
	orderID := uuid.New()
	customerID := uuid.New()

	repo := &mockShipmentRepository{
		createFunc: func(ctx context.Context, shipment *domain.Shipment) error {
			if shipment.OrderID != orderID {
				t.Errorf("expected order ID %v, got %v", orderID, shipment.OrderID)
			}
			if shipment.Status != domain.ShipmentStatusPending {
				t.Errorf("expected status pending, got %v", shipment.Status)
			}
			return nil
		},
	}

	uc := usecase.NewCreateShipmentUsecase(repo)

	input := usecase.CreateShipmentInput{
		OrderID:         orderID,
		CustomerID:      customerID,
		ShippingAddress: "123 Main St, City, Country",
		Carrier:         "DHL",
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Status != domain.ShipmentStatusPending {
		t.Errorf("expected status pending, got %v", output.Status)
	}
}
