package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error
}

type OrderItemRepository interface {
	Create(ctx context.Context, item *domain.OrderItem) error
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error)
}
