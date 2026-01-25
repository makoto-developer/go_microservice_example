package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/domain"
)

// OwnerRefreshTokenRepository defines the interface for owner refresh token persistence
type OwnerRefreshTokenRepository interface {
	// Create creates a new owner refresh token
	Create(ctx context.Context, token *domain.OwnerRefreshToken) error

	// FindByToken finds an owner refresh token by token string
	FindByToken(ctx context.Context, token string) (*domain.OwnerRefreshToken, error)

	// FindByUserID finds all owner refresh tokens for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.OwnerRefreshToken, error)

	// Delete deletes an owner refresh token
	Delete(ctx context.Context, token string) error

	// DeleteByUserID deletes all owner refresh tokens for a user
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// DeleteExpired deletes all expired owner refresh tokens
	DeleteExpired(ctx context.Context) error
}
