package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserPresenceRepository defines repository interface for UserPresence
type UserPresenceRepository interface {
	// Create creates a new UserPresence
	Create(ctx context.Context, user_presence *UserPresence) error

	// FindByID finds UserPresence by ID
	FindByID(ctx context.Context, id uuid.UUID) (*UserPresence, error)

	// Update updates UserPresence
	Update(ctx context.Context, user_presence *UserPresence) error

	// Delete deletes UserPresence by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all UserPresence
	List(ctx context.Context, limit, offset int) ([]*UserPresence, error)

	// FindByUserId finds UserPresence by user_id
	FindByUserId(ctx context.Context, user_id uuid.UUID) (*UserPresence, error)
}
