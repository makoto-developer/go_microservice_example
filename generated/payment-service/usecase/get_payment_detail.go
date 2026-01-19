package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetPaymentDetailInput represents input for GetPaymentDetail
type GetPaymentDetailInput struct {
	PaymentId uuid.UUID
	RequesterId uuid.UUID
	RequesterRole string
}

// GetPaymentDetailUsecase defines the interface for GetPaymentDetail
type GetPaymentDetailUsecase interface {
	Execute(ctx context.Context, input GetPaymentDetailInput) error
}

type get_payment_detailUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetPaymentDetailUsecase creates a new instance
func NewGetPaymentDetailUsecase() GetPaymentDetailUsecase {
	return &get_payment_detailUsecaseImpl{}
}

// Execute executes GetPaymentDetail
func (u *get_payment_detailUsecaseImpl) Execute(ctx context.Context, input GetPaymentDetailInput) error {
	// TODO: Implement business logic

	return nil
}
