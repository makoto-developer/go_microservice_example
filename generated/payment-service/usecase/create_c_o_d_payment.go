package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateCODPaymentInput represents input for CreateCODPayment
type CreateCODPaymentInput struct {
	OrderId uuid.UUID
	Amount decimal.Decimal
}

// CreateCODPaymentUsecase defines the interface for CreateCODPayment
type CreateCODPaymentUsecase interface {
	Execute(ctx context.Context, input CreateCODPaymentInput) error
}

type create_c_o_d_paymentUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewCreateCODPaymentUsecase creates a new instance
func NewCreateCODPaymentUsecase() CreateCODPaymentUsecase {
	return &create_c_o_d_paymentUsecaseImpl{}
}

// Execute executes CreateCODPayment
func (u *create_c_o_d_paymentUsecaseImpl) Execute(ctx context.Context, input CreateCODPaymentInput) error {
	// TODO: Implement business logic

	return nil
}
