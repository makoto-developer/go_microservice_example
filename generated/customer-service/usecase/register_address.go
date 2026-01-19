package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// RegisterAddressInput represents input for RegisterAddress
type RegisterAddressInput struct {
	CustomerId uuid.UUID
	AddressName string
	PostalCode string
	Prefecture string
	City string
	AddressLine1 string
	AddressLine2 string
	RecipientName string
	RecipientPhone string
	IsDefault bool
}

// RegisterAddressUsecase defines the interface for RegisterAddress
type RegisterAddressUsecase interface {
	Execute(ctx context.Context, input RegisterAddressInput) error
}

type register_addressUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewRegisterAddressUsecase creates a new instance
func NewRegisterAddressUsecase() RegisterAddressUsecase {
	return &register_addressUsecaseImpl{}
}

// Execute executes RegisterAddress
func (u *register_addressUsecaseImpl) Execute(ctx context.Context, input RegisterAddressInput) error {
	// TODO: Implement business logic

	return nil
}
