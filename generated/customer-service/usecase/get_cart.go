package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetCartInput represents input for GetCart
type GetCartInput struct {
	CustomerId uuid.UUID
}

// GetCartUsecase defines the interface for GetCart
type GetCartUsecase interface {
	Execute(ctx context.Context, input GetCartInput) error
}

type get_cartUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetCartUsecase creates a new instance
func NewGetCartUsecase() GetCartUsecase {
	return &get_cartUsecaseImpl{}
}

// Execute executes GetCart
func (u *get_cartUsecaseImpl) Execute(ctx context.Context, input GetCartInput) error {
	// TODO: Implement business logic

	return nil
}
