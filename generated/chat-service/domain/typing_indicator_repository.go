package domain

import (
	"context"

	"github.com/google/uuid"
)

// TypingIndicatorRepository defines repository interface for TypingIndicator
type TypingIndicatorRepository interface {
	// Create creates a new TypingIndicator
	Create(ctx context.Context, typing_indicator *TypingIndicator) error

	// FindByID finds TypingIndicator by ID
	FindByID(ctx context.Context, id uuid.UUID) (*TypingIndicator, error)

	// Update updates TypingIndicator
	Update(ctx context.Context, typing_indicator *TypingIndicator) error

	// Delete deletes TypingIndicator by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all TypingIndicator
	List(ctx context.Context, limit, offset int) ([]*TypingIndicator, error)
}
