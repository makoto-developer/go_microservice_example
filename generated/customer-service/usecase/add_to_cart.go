package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// AddToCartInput represents input for AddToCart
type AddToCartInput struct {
	CustomerId uuid.UUID
	ProductId uuid.UUID
	VariationId uuid.UUID
	Quantity int
}

// AddToCartUsecase defines the interface for AddToCart
type AddToCartUsecase interface {
	Execute(ctx context.Context, input AddToCartInput) error
}

type add_to_cartUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewAddToCartUsecase creates a new instance
func NewAddToCartUsecase() AddToCartUsecase {
	return &add_to_cartUsecaseImpl{}
}

// Execute executes AddToCart
func (u *add_to_cartUsecaseImpl) Execute(ctx context.Context, input AddToCartInput) error {
	// TODO: Implement business logic

	return nil
}
