package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/domain"
)

// CustomerUserRepository defines the interface for customer user persistence
type CustomerUserRepository interface {
	// Create creates a new customer user
	Create(ctx context.Context, user *domain.CustomerUser) error

	// FindByID finds a customer user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*domain.CustomerUser, error)

	// FindByEmail finds a customer user by email
	FindByEmail(ctx context.Context, email string) (*domain.CustomerUser, error)

	// FindByVerificationToken finds a customer user by email verification token
	FindByVerificationToken(ctx context.Context, token string) (*domain.CustomerUser, error)

	// FindByResetToken finds a customer user by password reset token
	FindByResetToken(ctx context.Context, token string) (*domain.CustomerUser, error)

	// Update updates a customer user
	Update(ctx context.Context, user *domain.CustomerUser) error

	// Delete deletes a customer user
	Delete(ctx context.Context, id uuid.UUID) error
}
