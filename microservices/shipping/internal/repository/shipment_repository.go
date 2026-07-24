package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/domain"
)

type ShipmentRepository interface {
	Create(ctx context.Context, shipment *domain.Shipment) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Shipment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Shipment, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Shipment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ShipmentStatus) error
	UpdateTracking(ctx context.Context, id uuid.UUID, trackingNumber string) error
}
