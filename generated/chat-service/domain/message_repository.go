package domain

import (
	"context"

	"github.com/google/uuid"
)

// MessageRepository defines repository interface for Message
type MessageRepository interface {
	// Create creates a new Message
	Create(ctx context.Context, message *Message) error

	// FindByID finds Message by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Message, error)

	// Update updates Message
	Update(ctx context.Context, message *Message) error

	// Delete deletes Message by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Message
	List(ctx context.Context, limit, offset int) ([]*Message, error)
}
