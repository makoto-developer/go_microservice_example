package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SendMessageInput represents input for SendMessage
type SendMessageInput struct {
	RoomId uuid.UUID
	SenderId uuid.UUID
	ReceiverId uuid.UUID
	MessageType MessageType
	Content string
	FileUrl string
	FileName string
	FileSize int
}

// SendMessageUsecase defines the interface for SendMessage
type SendMessageUsecase interface {
	Execute(ctx context.Context, input SendMessageInput) error
}

type send_messageUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewSendMessageUsecase creates a new instance
func NewSendMessageUsecase() SendMessageUsecase {
	return &send_messageUsecaseImpl{}
}

// Execute executes SendMessage
func (u *send_messageUsecaseImpl) Execute(ctx context.Context, input SendMessageInput) error {
	// TODO: Implement business logic

	return nil
}
