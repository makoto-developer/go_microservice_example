package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/auth-service/domain"
	"golang.org/x/crypto/bcrypt"
)

// UserLoginInput represents input for UserLogin
type UserLoginInput struct {
	Email        string
	Password     string
	KeepLoggedIn bool
}

// UserLoginOutput represents output for UserLogin
type UserLoginOutput struct {
	UserID       uuid.UUID
	AccessToken  string
	RefreshToken string
}

// UserLoginUsecase defines the interface for UserLogin
type UserLoginUsecase interface {
	Execute(ctx context.Context, input UserLoginInput) (*UserLoginOutput, error)
}

type userLoginUsecaseImpl struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	jwtSecret        string
}

// NewUserLoginUsecase creates a new instance
func NewUserLoginUsecase(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	jwtSecret string,
) UserLoginUsecase {
	return &userLoginUsecaseImpl{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
	}
}

// Execute executes UserLogin
func (u *userLoginUsecaseImpl) Execute(ctx context.Context, input UserLoginInput) (*UserLoginOutput, error) {
	// ユーザー検索
	user, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// パスワード検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// メール認証確認
	if !user.EmailVerified {
		return nil, fmt.Errorf("email not verified")
	}

	// JWT トークン生成（簡易実装）
	accessToken := fmt.Sprintf("access_%s", user.Id.String()[:8])
	refreshTokenStr := fmt.Sprintf("refresh_%s", user.Id.String()[:8])

	// Refresh Token の有効期限
	refreshTokenExpires := time.Now().Add(7 * 24 * time.Hour)
	if input.KeepLoggedIn {
		refreshTokenExpires = time.Now().Add(30 * 24 * time.Hour)
	}

	// Refresh Token 保存
	refreshToken := &domain.RefreshToken{
		Id:        uuid.New(),
		UserId:    user.Id,
		Token:     refreshTokenStr,
		ExpiresAt: refreshTokenExpires,
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	if err := u.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &UserLoginOutput{
		UserID:       user.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
	}, nil
}
