package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CancelOrderInput represents input for CancelOrder
type CancelOrderInput struct {
	OrderId uuid.UUID
	CancelledBy string
	CancelReason CancelReason
	CancelNote string
}

// CancelOrderUsecase defines the interface for CancelOrder
type CancelOrderUsecase interface {
	Execute(ctx context.Context, input CancelOrderInput) error
}

type cancel_orderUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewCancelOrderUsecase creates a new instance
func NewCancelOrderUsecase() CancelOrderUsecase {
	return &cancel_orderUsecaseImpl{}
}

// Execute executes CancelOrder
func (u *cancel_orderUsecaseImpl) Execute(ctx context.Context, input CancelOrderInput) error {
	// TODO: Implement business logic

	return nil
}
