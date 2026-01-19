package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetUserDetailInput represents input for GetUserDetail
type GetUserDetailInput struct {
	UserId uuid.UUID
}

// GetUserDetailUsecase defines the interface for GetUserDetail
type GetUserDetailUsecase interface {
	Execute(ctx context.Context, input GetUserDetailInput) error
}

type get_user_detailUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetUserDetailUsecase creates a new instance
func NewGetUserDetailUsecase() GetUserDetailUsecase {
	return &get_user_detailUsecaseImpl{}
}

// Execute executes GetUserDetail
func (u *get_user_detailUsecaseImpl) Execute(ctx context.Context, input GetUserDetailInput) error {
	// TODO: Implement business logic

	return nil
}
