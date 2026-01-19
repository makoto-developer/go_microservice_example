package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
)

type OrderFilter struct {
	ShopID       uuid.UUID
	Status       *domain.OrderStatus
	DateFrom     *time.Time
	DateTo       *time.Time
	CustomerName *string
	ProductName  *string
	SortBy       string
	SortOrder    string
}

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error
	List(ctx context.Context, filter OrderFilter) ([]*domain.Order, error)

	AddItem(ctx context.Context, item *domain.OrderItem) error
	GetItems(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error)
}
