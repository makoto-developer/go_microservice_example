package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ValidateAddressInput represents input for ValidateAddress
type ValidateAddressInput struct {
	PostalCode string
	Prefecture string
	City string
	AddressLine string
}

// ValidateAddressUsecase defines the interface for ValidateAddress
type ValidateAddressUsecase interface {
	Execute(ctx context.Context, input ValidateAddressInput) error
}

type validate_addressUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewValidateAddressUsecase creates a new instance
func NewValidateAddressUsecase() ValidateAddressUsecase {
	return &validate_addressUsecaseImpl{}
}

// Execute executes ValidateAddress
func (u *validate_addressUsecaseImpl) Execute(ctx context.Context, input ValidateAddressInput) error {
	// TODO: Implement business logic

	return nil
}
