package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// MarkMessagesAsReadInput represents input for MarkMessagesAsRead
type MarkMessagesAsReadInput struct {
	RoomId uuid.UUID
	UserId uuid.UUID
	MessageIds []uuid.UUID
}

// MarkMessagesAsReadUsecase defines the interface for MarkMessagesAsRead
type MarkMessagesAsReadUsecase interface {
	Execute(ctx context.Context, input MarkMessagesAsReadInput) error
}

type mark_messages_as_readUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewMarkMessagesAsReadUsecase creates a new instance
func NewMarkMessagesAsReadUsecase() MarkMessagesAsReadUsecase {
	return &mark_messages_as_readUsecaseImpl{}
}

// Execute executes MarkMessagesAsRead
func (u *mark_messages_as_readUsecaseImpl) Execute(ctx context.Context, input MarkMessagesAsReadInput) error {
	// TODO: Implement business logic

	return nil
}
