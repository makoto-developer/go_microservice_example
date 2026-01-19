package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CalculateShippingFeeInput represents input for CalculateShippingFee
type CalculateShippingFeeInput struct {
	ShopId uuid.UUID
	Prefecture string
	TotalWeightKg decimal.Decimal
	TotalSizeCm int
	Subtotal decimal.Decimal
}

// CalculateShippingFeeUsecase defines the interface for CalculateShippingFee
type CalculateShippingFeeUsecase interface {
	Execute(ctx context.Context, input CalculateShippingFeeInput) error
}

type calculate_shipping_feeUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewCalculateShippingFeeUsecase creates a new instance
func NewCalculateShippingFeeUsecase() CalculateShippingFeeUsecase {
	return &calculate_shipping_feeUsecaseImpl{}
}

// Execute executes CalculateShippingFee
func (u *calculate_shipping_feeUsecaseImpl) Execute(ctx context.Context, input CalculateShippingFeeInput) error {
	// TODO: Implement business logic

	return nil
}
