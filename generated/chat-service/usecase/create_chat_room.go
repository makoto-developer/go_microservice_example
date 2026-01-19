package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateChatRoomInput represents input for CreateChatRoom
type CreateChatRoomInput struct {
	CustomerId uuid.UUID
	ShopId uuid.UUID
	ProductId uuid.UUID
}

// CreateChatRoomUsecase defines the interface for CreateChatRoom
type CreateChatRoomUsecase interface {
	Execute(ctx context.Context, input CreateChatRoomInput) error
}

type create_chat_roomUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewCreateChatRoomUsecase creates a new instance
func NewCreateChatRoomUsecase() CreateChatRoomUsecase {
	return &create_chat_roomUsecaseImpl{}
}

// Execute executes CreateChatRoom
func (u *create_chat_roomUsecaseImpl) Execute(ctx context.Context, input CreateChatRoomInput) error {
	// TODO: Implement business logic

	return nil
}
