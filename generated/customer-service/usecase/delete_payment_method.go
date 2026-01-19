package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeletePaymentMethodInput represents input for DeletePaymentMethod
type DeletePaymentMethodInput struct {
	PaymentMethodId uuid.UUID
	CustomerId uuid.UUID
}

// DeletePaymentMethodUsecase defines the interface for DeletePaymentMethod
type DeletePaymentMethodUsecase interface {
	Execute(ctx context.Context, input DeletePaymentMethodInput) error
}

type delete_payment_methodUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewDeletePaymentMethodUsecase creates a new instance
func NewDeletePaymentMethodUsecase() DeletePaymentMethodUsecase {
	return &delete_payment_methodUsecaseImpl{}
}

// Execute executes DeletePaymentMethod
func (u *delete_payment_methodUsecaseImpl) Execute(ctx context.Context, input DeletePaymentMethodInput) error {
	// TODO: Implement business logic

	return nil
}
