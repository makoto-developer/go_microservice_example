package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateRefundInput represents input for CreateRefund
type CreateRefundInput struct {
	PaymentId uuid.UUID
	OrderId uuid.UUID
	Amount decimal.Decimal
	Reason string
}

// CreateRefundUsecase defines the interface for CreateRefund
type CreateRefundUsecase interface {
	Execute(ctx context.Context, input CreateRefundInput) error
}

type create_refundUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewCreateRefundUsecase creates a new instance
func NewCreateRefundUsecase() CreateRefundUsecase {
	return &create_refundUsecaseImpl{}
}

// Execute executes CreateRefund
func (u *create_refundUsecaseImpl) Execute(ctx context.Context, input CreateRefundInput) error {
	// TODO: Implement business logic

	return nil
}
