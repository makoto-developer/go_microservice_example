package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeleteAddressInput represents input for DeleteAddress
type DeleteAddressInput struct {
	AddressId uuid.UUID
	CustomerId uuid.UUID
}

// DeleteAddressUsecase defines the interface for DeleteAddress
type DeleteAddressUsecase interface {
	Execute(ctx context.Context, input DeleteAddressInput) error
}

type delete_addressUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewDeleteAddressUsecase creates a new instance
func NewDeleteAddressUsecase() DeleteAddressUsecase {
	return &delete_addressUsecaseImpl{}
}

// Execute executes DeleteAddress
func (u *delete_addressUsecaseImpl) Execute(ctx context.Context, input DeleteAddressInput) error {
	// TODO: Implement business logic

	return nil
}
