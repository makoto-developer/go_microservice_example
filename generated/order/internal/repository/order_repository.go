package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/domain"
)

type OrderRepository interface {
	// Create creates a new order with items
	Create(ctx context.Context, order *domain.Order, items []*domain.OrderItem) error

	// GetByID retrieves an order by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)

	// GetByOrderNumber retrieves an order by order number
	GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error)

	// GetByCustomerID retrieves orders for a customer
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Order, error)

	// GetItems retrieves order items for an order
	GetItems(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error)

	// UpdateStatus updates order status
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error

	// Cancel cancels an order
	Cancel(ctx context.Context, id uuid.UUID) error
}
