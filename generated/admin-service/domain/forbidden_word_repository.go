package domain

import (
	"context"

	"github.com/google/uuid"
)

// ForbiddenWordRepository defines repository interface for ForbiddenWord
type ForbiddenWordRepository interface {
	// Create creates a new ForbiddenWord
	Create(ctx context.Context, forbidden_word *ForbiddenWord) error

	// FindByID finds ForbiddenWord by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ForbiddenWord, error)

	// Update updates ForbiddenWord
	Update(ctx context.Context, forbidden_word *ForbiddenWord) error

	// Delete deletes ForbiddenWord by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ForbiddenWord
	List(ctx context.Context, limit, offset int) ([]*ForbiddenWord, error)
}
