package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdateAddressInput represents input for UpdateAddress
type UpdateAddressInput struct {
	AddressId uuid.UUID
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

// UpdateAddressUsecase defines the interface for UpdateAddress
type UpdateAddressUsecase interface {
	Execute(ctx context.Context, input UpdateAddressInput) error
}

type update_addressUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdateAddressUsecase creates a new instance
func NewUpdateAddressUsecase() UpdateAddressUsecase {
	return &update_addressUsecaseImpl{}
}

// Execute executes UpdateAddress
func (u *update_addressUsecaseImpl) Execute(ctx context.Context, input UpdateAddressInput) error {
	// TODO: Implement business logic

	return nil
}
