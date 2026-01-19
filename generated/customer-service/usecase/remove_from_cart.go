package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RemoveFromCartInput represents input for RemoveFromCart
type RemoveFromCartInput struct {
	CartItemId uuid.UUID
	CustomerId uuid.UUID
}

// RemoveFromCartUsecase defines the interface for RemoveFromCart
type RemoveFromCartUsecase interface {
	Execute(ctx context.Context, input RemoveFromCartInput) error
}

type remove_from_cartUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewRemoveFromCartUsecase creates a new instance
func NewRemoveFromCartUsecase() RemoveFromCartUsecase {
	return &remove_from_cartUsecaseImpl{}
}

// Execute executes RemoveFromCart
func (u *remove_from_cartUsecaseImpl) Execute(ctx context.Context, input RemoveFromCartInput) error {
	// TODO: Implement business logic

	return nil
}
