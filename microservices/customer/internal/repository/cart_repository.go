package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
)

type CartRepository interface {
	AddItem(ctx context.Context, item *domain.CartItem) error
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.CartItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.CartItem, error)
	UpdateQuantity(ctx context.Context, cartItemID uuid.UUID, quantity int) error
	RemoveItem(ctx context.Context, id uuid.UUID) error
	ClearCart(ctx context.Context, customerID uuid.UUID) error
	
	AddGuestItem(ctx context.Context, item *domain.GuestCartItem) error
	GetBySessionID(ctx context.Context, sessionID string) ([]*domain.GuestCartItem, error)
	ClearGuestCart(ctx context.Context, sessionID string) error
}
