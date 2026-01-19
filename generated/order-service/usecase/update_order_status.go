package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdateOrderStatusInput represents input for UpdateOrderStatus
type UpdateOrderStatusInput struct {
	OrderId uuid.UUID
	NewStatus OrderStatus
	ChangedBy string
	ChangeReason string
	TrackingNumber string
	Carrier string
}

// UpdateOrderStatusUsecase defines the interface for UpdateOrderStatus
type UpdateOrderStatusUsecase interface {
	Execute(ctx context.Context, input UpdateOrderStatusInput) error
}

type update_order_statusUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdateOrderStatusUsecase creates a new instance
func NewUpdateOrderStatusUsecase() UpdateOrderStatusUsecase {
	return &update_order_statusUsecaseImpl{}
}

// Execute executes UpdateOrderStatus
func (u *update_order_statusUsecaseImpl) Execute(ctx context.Context, input UpdateOrderStatusInput) error {
	// TODO: Implement business logic

	return nil
}
