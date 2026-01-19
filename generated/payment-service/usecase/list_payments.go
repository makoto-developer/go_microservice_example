package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ListPaymentsInput represents input for ListPayments
type ListPaymentsInput struct {
	OrderId uuid.UUID
	CustomerId uuid.UUID
	ShopId uuid.UUID
	StatusFilter []PaymentStatus
	DateFrom date
	DateTo date
	PaymentMethod PaymentMethodType
	Page int
	PageSize int
}

// ListPaymentsUsecase defines the interface for ListPayments
type ListPaymentsUsecase interface {
	Execute(ctx context.Context, input ListPaymentsInput) error
}

type list_paymentsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewListPaymentsUsecase creates a new instance
func NewListPaymentsUsecase() ListPaymentsUsecase {
	return &list_paymentsUsecaseImpl{}
}

// Execute executes ListPayments
func (u *list_paymentsUsecaseImpl) Execute(ctx context.Context, input ListPaymentsInput) error {
	// TODO: Implement business logic

	return nil
}
