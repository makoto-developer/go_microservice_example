package domain

import (
	"context"

	"github.com/google/uuid"
)

// ChatRoomRepository defines repository interface for ChatRoom
type ChatRoomRepository interface {
	// Create creates a new ChatRoom
	Create(ctx context.Context, chat_room *ChatRoom) error

	// FindByID finds ChatRoom by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ChatRoom, error)

	// Update updates ChatRoom
	Update(ctx context.Context, chat_room *ChatRoom) error

	// Delete deletes ChatRoom by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ChatRoom
	List(ctx context.Context, limit, offset int) ([]*ChatRoom, error)
}
