package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// GetChatRoomDetailInput represents input for GetChatRoomDetail
type GetChatRoomDetailInput struct {
	RoomId uuid.UUID
	UserId uuid.UUID
}

// GetChatRoomDetailUsecase defines the interface for GetChatRoomDetail
type GetChatRoomDetailUsecase interface {
	Execute(ctx context.Context, input GetChatRoomDetailInput) error
}

type get_chat_room_detailUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewGetChatRoomDetailUsecase creates a new instance
func NewGetChatRoomDetailUsecase() GetChatRoomDetailUsecase {
	return &get_chat_room_detailUsecaseImpl{}
}

// Execute executes GetChatRoomDetail
func (u *get_chat_room_detailUsecaseImpl) Execute(ctx context.Context, input GetChatRoomDetailInput) error {
	// TODO: Implement business logic

	return nil
}
