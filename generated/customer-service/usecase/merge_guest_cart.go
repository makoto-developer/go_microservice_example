package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// MergeGuestCartInput represents input for MergeGuestCart
type MergeGuestCartInput struct {
	CustomerId uuid.UUID
	SessionId string
}

// MergeGuestCartUsecase defines the interface for MergeGuestCart
type MergeGuestCartUsecase interface {
	Execute(ctx context.Context, input MergeGuestCartInput) error
}

type merge_guest_cartUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewMergeGuestCartUsecase creates a new instance
func NewMergeGuestCartUsecase() MergeGuestCartUsecase {
	return &merge_guest_cartUsecaseImpl{}
}

// Execute executes MergeGuestCart
func (u *merge_guest_cartUsecaseImpl) Execute(ctx context.Context, input MergeGuestCartInput) error {
	// TODO: Implement business logic

	return nil
}
