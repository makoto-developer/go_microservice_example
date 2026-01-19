package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdateCustomerProfileInput represents input for UpdateCustomerProfile
type UpdateCustomerProfileInput struct {
	CustomerId uuid.UUID
	FirstName string
	LastName string
	Phone string
	BirthDate date
	Gender Gender
}

// UpdateCustomerProfileUsecase defines the interface for UpdateCustomerProfile
type UpdateCustomerProfileUsecase interface {
	Execute(ctx context.Context, input UpdateCustomerProfileInput) error
}

type update_customer_profileUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdateCustomerProfileUsecase creates a new instance
func NewUpdateCustomerProfileUsecase() UpdateCustomerProfileUsecase {
	return &update_customer_profileUsecaseImpl{}
}

// Execute executes UpdateCustomerProfile
func (u *update_customer_profileUsecaseImpl) Execute(ctx context.Context, input UpdateCustomerProfileInput) error {
	// TODO: Implement business logic

	return nil
}
