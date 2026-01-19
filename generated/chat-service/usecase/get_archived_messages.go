package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetArchivedMessagesInput represents input for GetArchivedMessages
type GetArchivedMessagesInput struct {
	RoomId uuid.UUID
	UserId uuid.UUID
	DateFrom date
	DateTo date
}

// GetArchivedMessagesUsecase defines the interface for GetArchivedMessages
type GetArchivedMessagesUsecase interface {
	Execute(ctx context.Context, input GetArchivedMessagesInput) error
}

type get_archived_messagesUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetArchivedMessagesUsecase creates a new instance
func NewGetArchivedMessagesUsecase() GetArchivedMessagesUsecase {
	return &get_archived_messagesUsecaseImpl{}
}

// Execute executes GetArchivedMessages
func (u *get_archived_messagesUsecaseImpl) Execute(ctx context.Context, input GetArchivedMessagesInput) error {
	// TODO: Implement business logic

	return nil
}
