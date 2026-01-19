package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ChangeUserRoleInput represents input for ChangeUserRole
type ChangeUserRoleInput struct {
	UserId uuid.UUID
	AdminId uuid.UUID
	NewRole string
	Reason string
}

// ChangeUserRoleUsecase defines the interface for ChangeUserRole
type ChangeUserRoleUsecase interface {
	Execute(ctx context.Context, input ChangeUserRoleInput) error
}

type change_user_roleUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewChangeUserRoleUsecase creates a new instance
func NewChangeUserRoleUsecase() ChangeUserRoleUsecase {
	return &change_user_roleUsecaseImpl{}
}

// Execute executes ChangeUserRole
func (u *change_user_roleUsecaseImpl) Execute(ctx context.Context, input ChangeUserRoleInput) error {
	// TODO: Implement business logic

	return nil
}
