package domain

import (
	"context"

	"github.com/google/uuid"
)

// UserRepository defines repository interface for User
type UserRepository interface {
	// Create creates a new User
	Create(ctx context.Context, user *User) error

	// FindByID finds User by ID
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// Update updates User
	Update(ctx context.Context, user *User) error

	// Delete deletes User by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all User
	List(ctx context.Context, limit, offset int) ([]*User, error)

	// FindByEmail finds User by email
	FindByEmail(ctx context.Context, email string) (*User, error)
}
