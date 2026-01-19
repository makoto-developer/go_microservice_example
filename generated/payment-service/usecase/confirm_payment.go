package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ConfirmPaymentInput represents input for ConfirmPayment
type ConfirmPaymentInput struct {
	PaymentId uuid.UUID
	PaymentIntentId string
}

// ConfirmPaymentUsecase defines the interface for ConfirmPayment
type ConfirmPaymentUsecase interface {
	Execute(ctx context.Context, input ConfirmPaymentInput) error
}

type confirm_paymentUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewConfirmPaymentUsecase creates a new instance
func NewConfirmPaymentUsecase() ConfirmPaymentUsecase {
	return &confirm_paymentUsecaseImpl{}
}

// Execute executes ConfirmPayment
func (u *confirm_paymentUsecaseImpl) Execute(ctx context.Context, input ConfirmPaymentInput) error {
	// TODO: Implement business logic

	return nil
}
