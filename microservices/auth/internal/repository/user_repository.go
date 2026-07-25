package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByPasswordResetToken(ctx context.Context, token string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	UpdateEmailVerification(ctx context.Context, userID uuid.UUID, verified bool) error
	UpdatePasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt *time.Time) error
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
}
