package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/makoto-developer/go_microservice_example/generated/auth/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetUsecase struct {
	userRepo     repository.UserRepository
	emailService *EmailService
}

func NewPasswordResetUsecase(userRepo repository.UserRepository, emailService *EmailService) *PasswordResetUsecase {
	return &PasswordResetUsecase{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

// RequestPasswordReset パスワードリセットをリクエスト（トークン生成＋メール送信）
func (u *PasswordResetUsecase) RequestPasswordReset(ctx context.Context, email string) error {
	// ユーザーを取得
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// セキュリティ上、ユーザーが存在しない場合でもエラーを返さない
		// （メールアドレスの存在を確認されないため）
		return nil
	}

	// リセットトークンを生成
	token, err := generateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// トークンの有効期限（24時間）
	expiresAt := time.Now().Add(24 * time.Hour)

	// トークンをデータベースに保存
	err = u.userRepo.UpdatePasswordResetToken(ctx, user.ID, token, &expiresAt)
	if err != nil {
		return fmt.Errorf("failed to update reset token: %w", err)
	}

	// リセットメールを送信
	err = u.emailService.SendPasswordResetEmail(email, token)
	if err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword パスワードリセットを実行（トークン検証＋パスワード更新）
func (u *PasswordResetUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	// トークンからユーザーを検索
	user, err := u.findUserByResetToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// トークンの有効期限をチェック
	if user.PasswordResetExpiresAt == nil || time.Now().After(*user.PasswordResetExpiresAt) {
		return fmt.Errorf("reset token has expired")
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// パスワードを更新
	err = u.userRepo.UpdatePassword(ctx, user.ID, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// リセットトークンをクリア
	err = u.userRepo.UpdatePasswordResetToken(ctx, user.ID, "", nil)
	if err != nil {
		return fmt.Errorf("failed to clear reset token: %w", err)
	}

	return nil
}

// findUserByResetToken トークンからユーザーを検索
func (u *PasswordResetUsecase) findUserByResetToken(ctx context.Context, token string) (*domain.User, error) {
	return u.userRepo.FindByPasswordResetToken(ctx, token)
}

// generateSecureToken セキュアなランダムトークンを生成
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
