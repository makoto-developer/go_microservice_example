package usecase

import (
	"context"
	"fmt"

	"github.com/makoto-developer/go_microservice_example/generated/auth-service/domain"
)

// UserLogoutInput represents input for UserLogout
type UserLogoutInput struct {
	RefreshToken string
}

// UserLogoutUsecase defines the interface for UserLogout
type UserLogoutUsecase interface {
	Execute(ctx context.Context, input UserLogoutInput) error
}

type userLogoutUsecaseImpl struct {
	refreshTokenRepo domain.RefreshTokenRepository
}

// NewUserLogoutUsecase creates a new instance
func NewUserLogoutUsecase(
	refreshTokenRepo domain.RefreshTokenRepository,
) UserLogoutUsecase {
	return &userLogoutUsecaseImpl{
		refreshTokenRepo: refreshTokenRepo,
	}
}

// Execute executes UserLogout
func (u *userLogoutUsecaseImpl) Execute(ctx context.Context, input UserLogoutInput) error {
	// Refresh Token 検索
	refreshToken, err := u.refreshTokenRepo.FindByToken(ctx, input.RefreshToken)
	if err != nil {
		return fmt.Errorf("invalid refresh token")
	}

	// トークンを無効化
	refreshToken.Revoked = true

	// 更新
	if err := u.refreshTokenRepo.Update(ctx, refreshToken); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}
