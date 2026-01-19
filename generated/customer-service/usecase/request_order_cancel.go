package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RequestOrderCancelInput represents input for RequestOrderCancel
type RequestOrderCancelInput struct {
	OrderId uuid.UUID
	CustomerId uuid.UUID
	CancelReason string
}

// RequestOrderCancelUsecase defines the interface for RequestOrderCancel
type RequestOrderCancelUsecase interface {
	Execute(ctx context.Context, input RequestOrderCancelInput) error
}

type request_order_cancelUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewRequestOrderCancelUsecase creates a new instance
func NewRequestOrderCancelUsecase() RequestOrderCancelUsecase {
	return &request_order_cancelUsecaseImpl{}
}

// Execute executes RequestOrderCancel
func (u *request_order_cancelUsecaseImpl) Execute(ctx context.Context, input RequestOrderCancelInput) error {
	// TODO: Implement business logic

	return nil
}
