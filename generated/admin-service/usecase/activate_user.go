package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ActivateUserInput represents input for ActivateUser
type ActivateUserInput struct {
	UserId uuid.UUID
	AdminId uuid.UUID
	Reason string
}

// ActivateUserUsecase defines the interface for ActivateUser
type ActivateUserUsecase interface {
	Execute(ctx context.Context, input ActivateUserInput) error
}

type activate_userUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewActivateUserUsecase creates a new instance
func NewActivateUserUsecase() ActivateUserUsecase {
	return &activate_userUsecaseImpl{}
}

// Execute executes ActivateUser
func (u *activate_userUsecaseImpl) Execute(ctx context.Context, input ActivateUserInput) error {
	// TODO: Implement business logic

	return nil
}
