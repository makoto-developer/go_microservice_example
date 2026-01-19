package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdatePresenceInput represents input for UpdatePresence
type UpdatePresenceInput struct {
	UserId uuid.UUID
	Status PresenceStatus
	DeviceInfo string
}

// UpdatePresenceUsecase defines the interface for UpdatePresence
type UpdatePresenceUsecase interface {
	Execute(ctx context.Context, input UpdatePresenceInput) error
}

type update_presenceUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdatePresenceUsecase creates a new instance
func NewUpdatePresenceUsecase() UpdatePresenceUsecase {
	return &update_presenceUsecaseImpl{}
}

// Execute executes UpdatePresence
func (u *update_presenceUsecaseImpl) Execute(ctx context.Context, input UpdatePresenceInput) error {
	// TODO: Implement business logic

	return nil
}
