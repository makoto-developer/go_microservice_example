package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetOrderDetailInput represents input for GetOrderDetail
type GetOrderDetailInput struct {
	OrderId uuid.UUID
	CustomerId uuid.UUID
}

// GetOrderDetailUsecase defines the interface for GetOrderDetail
type GetOrderDetailUsecase interface {
	Execute(ctx context.Context, input GetOrderDetailInput) error
}

type get_order_detailUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetOrderDetailUsecase creates a new instance
func NewGetOrderDetailUsecase() GetOrderDetailUsecase {
	return &get_order_detailUsecaseImpl{}
}

// Execute executes GetOrderDetail
func (u *get_order_detailUsecaseImpl) Execute(ctx context.Context, input GetOrderDetailInput) error {
	// TODO: Implement business logic

	return nil
}
