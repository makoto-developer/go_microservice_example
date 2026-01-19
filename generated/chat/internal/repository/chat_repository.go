package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/chat/internal/domain"
)

type ChatRepository interface {
	// ChatRoom operations
	CreateRoom(ctx context.Context, room *domain.ChatRoom) error
	GetRoomByID(ctx context.Context, id uuid.UUID) (*domain.ChatRoom, error)
	GetRoomByCustomerAndShop(ctx context.Context, customerID, shopID uuid.UUID) (*domain.ChatRoom, error)
	GetRoomsByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.ChatRoom, error)

	// Message operations
	CreateMessage(ctx context.Context, message *domain.Message) error
	GetMessagesByRoomID(ctx context.Context, roomID uuid.UUID) ([]*domain.Message, error)
	MarkMessageRead(ctx context.Context, messageID uuid.UUID) error
}
