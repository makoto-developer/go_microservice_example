package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreatePaymentIntentInput represents input for CreatePaymentIntent
type CreatePaymentIntentInput struct {
	OrderId uuid.UUID
	Amount decimal.Decimal
	Currency string
	PaymentMethodId string
	CustomerId uuid.UUID
}

// CreatePaymentIntentUsecase defines the interface for CreatePaymentIntent
type CreatePaymentIntentUsecase interface {
	Execute(ctx context.Context, input CreatePaymentIntentInput) error
}

type create_payment_intentUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewCreatePaymentIntentUsecase creates a new instance
func NewCreatePaymentIntentUsecase() CreatePaymentIntentUsecase {
	return &create_payment_intentUsecaseImpl{}
}

// Execute executes CreatePaymentIntent
func (u *create_payment_intentUsecaseImpl) Execute(ctx context.Context, input CreatePaymentIntentInput) error {
	// TODO: Implement business logic

	return nil
}
