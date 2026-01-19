package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// StartTypingInput represents input for StartTyping
type StartTypingInput struct {
	RoomId uuid.UUID
	UserId uuid.UUID
}

// StartTypingUsecase defines the interface for StartTyping
type StartTypingUsecase interface {
	Execute(ctx context.Context, input StartTypingInput) error
}

type start_typingUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewStartTypingUsecase creates a new instance
func NewStartTypingUsecase() StartTypingUsecase {
	return &start_typingUsecaseImpl{}
}

// Execute executes StartTyping
func (u *start_typingUsecaseImpl) Execute(ctx context.Context, input StartTypingInput) error {
	// TODO: Implement business logic

	return nil
}
