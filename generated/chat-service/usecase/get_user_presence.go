package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetUserPresenceInput represents input for GetUserPresence
type GetUserPresenceInput struct {
	UserId uuid.UUID
}

// GetUserPresenceUsecase defines the interface for GetUserPresence
type GetUserPresenceUsecase interface {
	Execute(ctx context.Context, input GetUserPresenceInput) error
}

type get_user_presenceUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetUserPresenceUsecase creates a new instance
func NewGetUserPresenceUsecase() GetUserPresenceUsecase {
	return &get_user_presenceUsecaseImpl{}
}

// Execute executes GetUserPresence
func (u *get_user_presenceUsecaseImpl) Execute(ctx context.Context, input GetUserPresenceInput) error {
	// TODO: Implement business logic

	return nil
}
