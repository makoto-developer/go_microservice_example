package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/domain"
)

// CustomerRefreshTokenRepository defines the interface for customer refresh token persistence
type CustomerRefreshTokenRepository interface {
	// Create creates a new customer refresh token
	Create(ctx context.Context, token *domain.CustomerRefreshToken) error

	// FindByToken finds a customer refresh token by token string
	FindByToken(ctx context.Context, token string) (*domain.CustomerRefreshToken, error)

	// FindByUserID finds all customer refresh tokens for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.CustomerRefreshToken, error)

	// Delete deletes a customer refresh token
	Delete(ctx context.Context, token string) error

	// DeleteByUserID deletes all customer refresh tokens for a user
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// DeleteExpired deletes all expired customer refresh tokens
	DeleteExpired(ctx context.Context) error
}
