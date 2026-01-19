package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// StopTypingInput represents input for StopTyping
type StopTypingInput struct {
	RoomId uuid.UUID
	UserId uuid.UUID
}

// StopTypingUsecase defines the interface for StopTyping
type StopTypingUsecase interface {
	Execute(ctx context.Context, input StopTypingInput) error
}

type stop_typingUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewStopTypingUsecase creates a new instance
func NewStopTypingUsecase() StopTypingUsecase {
	return &stop_typingUsecaseImpl{}
}

// Execute executes StopTyping
func (u *stop_typingUsecaseImpl) Execute(ctx context.Context, input StopTypingInput) error {
	// TODO: Implement business logic

	return nil
}
