package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SuspendUserInput represents input for SuspendUser
type SuspendUserInput struct {
	UserId uuid.UUID
	AdminId uuid.UUID
	Reason string
}

// SuspendUserUsecase defines the interface for SuspendUser
type SuspendUserUsecase interface {
	Execute(ctx context.Context, input SuspendUserInput) error
}

type suspend_userUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewSuspendUserUsecase creates a new instance
func NewSuspendUserUsecase() SuspendUserUsecase {
	return &suspend_userUsecaseImpl{}
}

// Execute executes SuspendUser
func (u *suspend_userUsecaseImpl) Execute(ctx context.Context, input SuspendUserInput) error {
	// TODO: Implement business logic

	return nil
}
