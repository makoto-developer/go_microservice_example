package domain

import (
	"context"

	"github.com/google/uuid"
)

// RefreshTokenRepository defines repository interface for RefreshToken
type RefreshTokenRepository interface {
	// Create creates a new RefreshToken
	Create(ctx context.Context, refresh_token *RefreshToken) error

	// FindByID finds RefreshToken by ID
	FindByID(ctx context.Context, id uuid.UUID) (*RefreshToken, error)

	// Update updates RefreshToken
	Update(ctx context.Context, refresh_token *RefreshToken) error

	// Delete deletes RefreshToken by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all RefreshToken
	List(ctx context.Context, limit, offset int) ([]*RefreshToken, error)

	// FindByToken finds RefreshToken by token
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
}
