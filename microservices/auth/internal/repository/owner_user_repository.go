package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/domain"
)

// OwnerUserRepository defines the interface for owner user persistence
type OwnerUserRepository interface {
	// Create creates a new owner user
	Create(ctx context.Context, user *domain.OwnerUser) error

	// FindByID finds an owner user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*domain.OwnerUser, error)

	// FindByEmail finds an owner user by email
	FindByEmail(ctx context.Context, email string) (*domain.OwnerUser, error)

	// FindByVerificationToken finds an owner user by email verification token
	FindByVerificationToken(ctx context.Context, token string) (*domain.OwnerUser, error)

	// FindByResetToken finds an owner user by password reset token
	FindByResetToken(ctx context.Context, token string) (*domain.OwnerUser, error)

	// Update updates an owner user
	Update(ctx context.Context, user *domain.OwnerUser) error

	// Delete deletes an owner user
	Delete(ctx context.Context, id uuid.UUID) error

	// ListPendingVerifications lists owner users pending business verification
	ListPendingVerifications(ctx context.Context, limit, offset int) ([]*domain.OwnerUser, error)
}
