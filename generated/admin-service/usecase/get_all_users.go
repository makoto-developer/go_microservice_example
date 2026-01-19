package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetAllUsersInput represents input for GetAllUsers
type GetAllUsersInput struct {
	RoleFilter string
	StatusFilter string
	EmailVerified bool
	DateFrom date
	DateTo date
	SortBy string
	Page int
	PageSize int
}

// GetAllUsersUsecase defines the interface for GetAllUsers
type GetAllUsersUsecase interface {
	Execute(ctx context.Context, input GetAllUsersInput) error
}

type get_all_usersUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetAllUsersUsecase creates a new instance
func NewGetAllUsersUsecase() GetAllUsersUsecase {
	return &get_all_usersUsecaseImpl{}
}

// Execute executes GetAllUsers
func (u *get_all_usersUsecaseImpl) Execute(ctx context.Context, input GetAllUsersInput) error {
	// TODO: Implement business logic

	return nil
}
