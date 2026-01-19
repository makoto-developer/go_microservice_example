package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RegisterPaymentMethodInput represents input for RegisterPaymentMethod
type RegisterPaymentMethodInput struct {
	CustomerId uuid.UUID
	StripePaymentMethodId string
	IsDefault bool
}

// RegisterPaymentMethodUsecase defines the interface for RegisterPaymentMethod
type RegisterPaymentMethodUsecase interface {
	Execute(ctx context.Context, input RegisterPaymentMethodInput) error
}

type register_payment_methodUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewRegisterPaymentMethodUsecase creates a new instance
func NewRegisterPaymentMethodUsecase() RegisterPaymentMethodUsecase {
	return &register_payment_methodUsecaseImpl{}
}

// Execute executes RegisterPaymentMethod
func (u *register_payment_methodUsecaseImpl) Execute(ctx context.Context, input RegisterPaymentMethodInput) error {
	// TODO: Implement business logic

	return nil
}
