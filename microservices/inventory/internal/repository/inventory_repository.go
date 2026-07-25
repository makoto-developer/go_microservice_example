package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/domain"
)

type InventoryRepository interface {
	Create(ctx context.Context, inventory *domain.Inventory) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error)
	GetByProductID(ctx context.Context, productID uuid.UUID, variationID *uuid.UUID) (*domain.Inventory, error)
	Update(ctx context.Context, inventory *domain.Inventory) error
	UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error
	Reserve(ctx context.Context, id uuid.UUID, quantity int) error
	Release(ctx context.Context, id uuid.UUID, quantity int) error
}

type ReservationRepository interface {
	Create(ctx context.Context, reservation *domain.Reservation) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.Reservation, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ReservationStatus) error
	DeleteExpired(ctx context.Context) error
}
