package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetChatRoomsInput represents input for GetChatRooms
type GetChatRoomsInput struct {
	UserId uuid.UUID
	UserRole string
	StatusFilter RoomStatus
}

// GetChatRoomsUsecase defines the interface for GetChatRooms
type GetChatRoomsUsecase interface {
	Execute(ctx context.Context, input GetChatRoomsInput) error
}

type get_chat_roomsUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetChatRoomsUsecase creates a new instance
func NewGetChatRoomsUsecase() GetChatRoomsUsecase {
	return &get_chat_roomsUsecaseImpl{}
}

// Execute executes GetChatRooms
func (u *get_chat_roomsUsecaseImpl) Execute(ctx context.Context, input GetChatRoomsInput) error {
	// TODO: Implement business logic

	return nil
}
