package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/makoto-developer/go_microservice_example/generated/auth-service/domain"
	"golang.org/x/crypto/bcrypt"
)

// PasswordResetInput represents input for PasswordReset
type PasswordResetInput struct {
	Token       string
	NewPassword string
}

// PasswordResetUsecase defines the interface for PasswordReset
type PasswordResetUsecase interface {
	Execute(ctx context.Context, input PasswordResetInput) error
}

type passwordResetUsecaseImpl struct {
	userRepo domain.UserRepository
}

// NewPasswordResetUsecase creates a new instance
func NewPasswordResetUsecase(
	userRepo domain.UserRepository,
) PasswordResetUsecase {
	return &passwordResetUsecaseImpl{
		userRepo: userRepo,
	}
}

// Execute executes PasswordReset
func (u *passwordResetUsecaseImpl) Execute(ctx context.Context, input PasswordResetInput) error {
	// 全ユーザーを取得してトークンで検索（実際にはFindByPasswordResetTokenメソッドが必要）
	users, err := u.userRepo.List(ctx, 1000, 0)
	if err != nil {
		return fmt.Errorf("failed to search user: %w", err)
	}

	var targetUser *domain.User
	for _, user := range users {
		if user.PasswordResetToken != nil && *user.PasswordResetToken == input.Token {
			targetUser = user
			break
		}
	}

	if targetUser == nil {
		return fmt.Errorf("invalid password reset token")
	}

	// トークン有効期限確認
	if targetUser.PasswordResetExpiresAt != nil && targetUser.PasswordResetExpiresAt.Before(time.Now()) {
		return fmt.Errorf("password reset token has expired")
	}

	// 新しいパスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// パスワード更新
	targetUser.PasswordHash = string(hashedPassword)
	targetUser.PasswordResetToken = nil
	targetUser.PasswordResetExpiresAt = nil
	targetUser.UpdatedAt = time.Now()

	// 更新
	if err := u.userRepo.Update(ctx, targetUser); err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	return nil
}
