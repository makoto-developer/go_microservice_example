package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateOrderInput represents input for CreateOrder
type CreateOrderInput struct {
	CustomerId uuid.UUID
	CustomerEmail string
	CartItems []CartItemInput
	ShippingAddressId uuid.UUID
	PaymentMethod PaymentMethod
	PaymentMethodId uuid.UUID
	ShippingMethod string
	Notes string
}

// CreateOrderUsecase defines the interface for CreateOrder
type CreateOrderUsecase interface {
	Execute(ctx context.Context, input CreateOrderInput) error
}

type create_orderUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewCreateOrderUsecase creates a new instance
func NewCreateOrderUsecase() CreateOrderUsecase {
	return &create_orderUsecaseImpl{}
}

// Execute executes CreateOrder
func (u *create_orderUsecaseImpl) Execute(ctx context.Context, input CreateOrderInput) error {
	// TODO: Implement business logic

	return nil
}
