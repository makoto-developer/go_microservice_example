package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SearchMessagesInput represents input for SearchMessages
type SearchMessagesInput struct {
	RoomId uuid.UUID
	UserId uuid.UUID
	Keyword string
	DateFrom date
	DateTo date
}

// SearchMessagesUsecase defines the interface for SearchMessages
type SearchMessagesUsecase interface {
	Execute(ctx context.Context, input SearchMessagesInput) error
}

type search_messagesUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewSearchMessagesUsecase creates a new instance
func NewSearchMessagesUsecase() SearchMessagesUsecase {
	return &search_messagesUsecaseImpl{}
}

// Execute executes SearchMessages
func (u *search_messagesUsecaseImpl) Execute(ctx context.Context, input SearchMessagesInput) error {
	// TODO: Implement business logic

	return nil
}
