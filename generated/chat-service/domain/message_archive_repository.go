package domain

import (
	"context"

	"github.com/google/uuid"
)

// MessageArchiveRepository defines repository interface for MessageArchive
type MessageArchiveRepository interface {
	// Create creates a new MessageArchive
	Create(ctx context.Context, message_archive *MessageArchive) error

	// FindByID finds MessageArchive by ID
	FindByID(ctx context.Context, id uuid.UUID) (*MessageArchive, error)

	// Update updates MessageArchive
	Update(ctx context.Context, message_archive *MessageArchive) error

	// Delete deletes MessageArchive by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all MessageArchive
	List(ctx context.Context, limit, offset int) ([]*MessageArchive, error)
}
