package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetMessagesInput represents input for GetMessages
type GetMessagesInput struct {
	RoomId uuid.UUID
	UserId uuid.UUID
	Page int
	PageSize int
}

// GetMessagesUsecase defines the interface for GetMessages
type GetMessagesUsecase interface {
	Execute(ctx context.Context, input GetMessagesInput) error
}

type get_messagesUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetMessagesUsecase creates a new instance
func NewGetMessagesUsecase() GetMessagesUsecase {
	return &get_messagesUsecaseImpl{}
}

// Execute executes GetMessages
func (u *get_messagesUsecaseImpl) Execute(ctx context.Context, input GetMessagesInput) error {
	// TODO: Implement business logic

	return nil
}
