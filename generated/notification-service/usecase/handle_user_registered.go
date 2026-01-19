package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// HandleUserRegisteredInput represents input for HandleUserRegistered
type HandleUserRegisteredInput struct {
	UserId uuid.UUID
	Email string
	VerificationLink string
}

// HandleUserRegisteredUsecase defines the interface for HandleUserRegistered
type HandleUserRegisteredUsecase interface {
	Execute(ctx context.Context, input HandleUserRegisteredInput) error
}

type handle_user_registeredUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewHandleUserRegisteredUsecase creates a new instance
func NewHandleUserRegisteredUsecase() HandleUserRegisteredUsecase {
	return &handle_user_registeredUsecaseImpl{}
}

// Execute executes HandleUserRegistered
func (u *handle_user_registeredUsecaseImpl) Execute(ctx context.Context, input HandleUserRegisteredInput) error {
	// TODO: Implement business logic

	return nil
}
