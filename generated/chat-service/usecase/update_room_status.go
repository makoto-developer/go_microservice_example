package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdateRoomStatusInput represents input for UpdateRoomStatus
type UpdateRoomStatusInput struct {
	RoomId uuid.UUID
	UserId uuid.UUID
	NewStatus RoomStatus
}

// UpdateRoomStatusUsecase defines the interface for UpdateRoomStatus
type UpdateRoomStatusUsecase interface {
	Execute(ctx context.Context, input UpdateRoomStatusInput) error
}

type update_room_statusUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdateRoomStatusUsecase creates a new instance
func NewUpdateRoomStatusUsecase() UpdateRoomStatusUsecase {
	return &update_room_statusUsecaseImpl{}
}

// Execute executes UpdateRoomStatus
func (u *update_room_statusUsecaseImpl) Execute(ctx context.Context, input UpdateRoomStatusInput) error {
	// TODO: Implement business logic

	return nil
}
