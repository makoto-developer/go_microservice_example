package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// NormalizeAddressInput represents input for NormalizeAddress
type NormalizeAddressInput struct {
	PostalCode string
	Prefecture string
	City string
	AddressLine string
	Building string
}

// NormalizeAddressUsecase defines the interface for NormalizeAddress
type NormalizeAddressUsecase interface {
	Execute(ctx context.Context, input NormalizeAddressInput) error
}

type normalize_addressUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewNormalizeAddressUsecase creates a new instance
func NewNormalizeAddressUsecase() NormalizeAddressUsecase {
	return &normalize_addressUsecaseImpl{}
}

// Execute executes NormalizeAddress
func (u *normalize_addressUsecaseImpl) Execute(ctx context.Context, input NormalizeAddressInput) error {
	// TODO: Implement business logic

	return nil
}
