package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/auth-service/domain"
	"golang.org/x/crypto/bcrypt"
)

// PasswordChangeInput represents input for PasswordChange
type PasswordChangeInput struct {
	UserId          uuid.UUID
	CurrentPassword string
	NewPassword     string
}

// PasswordChangeUsecase defines the interface for PasswordChange
type PasswordChangeUsecase interface {
	Execute(ctx context.Context, input PasswordChangeInput) error
}

type passwordChangeUsecaseImpl struct {
	userRepo domain.UserRepository
}

// NewPasswordChangeUsecase creates a new instance
func NewPasswordChangeUsecase(
	userRepo domain.UserRepository,
) PasswordChangeUsecase {
	return &passwordChangeUsecaseImpl{
		userRepo: userRepo,
	}
}

// Execute executes PasswordChange
func (u *passwordChangeUsecaseImpl) Execute(ctx context.Context, input PasswordChangeInput) error {
	// ユーザー検索
	user, err := u.userRepo.FindByID(ctx, input.UserId)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// 現在のパスワード検証
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.CurrentPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// 新しいパスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// パスワード更新
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	// 更新
	if err := u.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}
