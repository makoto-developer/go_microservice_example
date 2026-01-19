package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/makoto-developer/go_microservice_example/generated/auth-service/domain"
	"github.com/makoto-developer/go_microservice_example/manual/auth"
)

// PasswordResetRequestInput represents input for PasswordResetRequest
type PasswordResetRequestInput struct {
	Email string
}

// PasswordResetRequestOutput represents output for PasswordResetRequest
type PasswordResetRequestOutput struct {
	ResetToken string
}

// PasswordResetRequestUsecase defines the interface for PasswordResetRequest
type PasswordResetRequestUsecase interface {
	Execute(ctx context.Context, input PasswordResetRequestInput) (*PasswordResetRequestOutput, error)
}

type passwordResetRequestUsecaseImpl struct {
	userRepo    domain.UserRepository
	emailSender auth.EmailSender
}

// NewPasswordResetRequestUsecase creates a new instance
func NewPasswordResetRequestUsecase(
	userRepo domain.UserRepository,
	emailSender auth.EmailSender,
) PasswordResetRequestUsecase {
	return &passwordResetRequestUsecaseImpl{
		userRepo:    userRepo,
		emailSender: emailSender,
	}
}

// Execute executes PasswordResetRequest
func (u *passwordResetRequestUsecaseImpl) Execute(ctx context.Context, input PasswordResetRequestInput) (*PasswordResetRequestOutput, error) {
	// ユーザー検索
	user, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		// セキュリティ上、ユーザーが存在しない場合もエラーを返さない
		return nil, fmt.Errorf("if the email exists, a password reset link has been sent")
	}

	// パスワードリセットトークン生成
	resetToken := generateRandomToken(32)
	resetExpires := time.Now().Add(1 * time.Hour)

	// ユーザー更新
	user.PasswordResetToken = &resetToken
	user.PasswordResetExpiresAt = &resetExpires
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create password reset request: %w", err)
	}

	// パスワードリセットメールを送信
	if err := u.emailSender.SendPasswordResetEmail(ctx, input.Email, resetToken); err != nil {
		// メール送信失敗はログのみ（リセットリクエスト自体は成功）
		fmt.Printf("Warning: Failed to send password reset email to %s: %v\n", input.Email, err)
	}

	return &PasswordResetRequestOutput{
		ResetToken: resetToken,
	}, nil
}
