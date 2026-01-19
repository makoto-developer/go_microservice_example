package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/auth-service/domain"
)

// TokenRefreshInput represents input for TokenRefresh
type TokenRefreshInput struct {
	RefreshToken string
}

// TokenRefreshOutput represents output for TokenRefresh
type TokenRefreshOutput struct {
	AccessToken  string
	RefreshToken string
}

// TokenRefreshUsecase defines the interface for TokenRefresh
type TokenRefreshUsecase interface {
	Execute(ctx context.Context, input TokenRefreshInput) (*TokenRefreshOutput, error)
}

type tokenRefreshUsecaseImpl struct {
	refreshTokenRepo domain.RefreshTokenRepository
	userRepo         domain.UserRepository
	jwtSecret        string
}

// NewTokenRefreshUsecase creates a new instance
func NewTokenRefreshUsecase(
	refreshTokenRepo domain.RefreshTokenRepository,
	userRepo domain.UserRepository,
	jwtSecret string,
) TokenRefreshUsecase {
	return &tokenRefreshUsecaseImpl{
		refreshTokenRepo: refreshTokenRepo,
		userRepo:         userRepo,
		jwtSecret:        jwtSecret,
	}
}

// Execute executes TokenRefresh
func (u *tokenRefreshUsecaseImpl) Execute(ctx context.Context, input TokenRefreshInput) (*TokenRefreshOutput, error) {
	// Refresh Token 検索
	refreshToken, err := u.refreshTokenRepo.FindByToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// トークン検証
	if refreshToken.Revoked {
		return nil, fmt.Errorf("refresh token has been revoked")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token has expired")
	}

	// ユーザー検索
	user, err := u.userRepo.FindByID(ctx, refreshToken.UserId)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// 新しいアクセストークン生成
	newAccessToken := fmt.Sprintf("access_%s", user.Id.String()[:8])

	// 新しいリフレッシュトークン生成
	newRefreshTokenStr := fmt.Sprintf("refresh_%s", uuid.New().String()[:8])
	newRefreshTokenExpires := time.Now().Add(7 * 24 * time.Hour)

	newRefreshToken := &domain.RefreshToken{
		Id:        uuid.New(),
		UserId:    user.Id,
		Token:     newRefreshTokenStr,
		ExpiresAt: newRefreshTokenExpires,
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	// 古いトークンを無効化
	refreshToken.Revoked = true
	if err := u.refreshTokenRepo.Update(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	// 新しいトークンを保存
	if err := u.refreshTokenRepo.Create(ctx, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to create new refresh token: %w", err)
	}

	return &TokenRefreshOutput{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshTokenStr,
	}, nil
}
