package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/makoto-developer/go_microservice_example/generated/auth-service/domain"
)

// EmailVerificationInput represents input for EmailVerification
type EmailVerificationInput struct {
	Token string
}

// EmailVerificationUsecase defines the interface for EmailVerification
type EmailVerificationUsecase interface {
	Execute(ctx context.Context, input EmailVerificationInput) error
}

type emailVerificationUsecaseImpl struct {
	userRepo domain.UserRepository
}

// NewEmailVerificationUsecase creates a new instance
func NewEmailVerificationUsecase(
	userRepo domain.UserRepository,
) EmailVerificationUsecase {
	return &emailVerificationUsecaseImpl{
		userRepo: userRepo,
	}
}

// Execute executes EmailVerification
func (u *emailVerificationUsecaseImpl) Execute(ctx context.Context, input EmailVerificationInput) error {
	// 全ユーザーを取得してトークンで検索（実際にはFindByVerificationTokenメソッドが必要）
	users, err := u.userRepo.List(ctx, 1000, 0)
	if err != nil {
		return fmt.Errorf("failed to search user: %w", err)
	}

	var targetUser *domain.User
	for _, user := range users {
		if user.EmailVerificationToken != nil && *user.EmailVerificationToken == input.Token {
			targetUser = user
			break
		}
	}

	if targetUser == nil {
		return fmt.Errorf("invalid verification token")
	}

	// トークン有効期限確認
	if targetUser.EmailVerificationExpiresAt != nil && targetUser.EmailVerificationExpiresAt.Before(time.Now()) {
		return fmt.Errorf("verification token has expired")
	}

	// メール認証完了
	targetUser.EmailVerified = true
	targetUser.EmailVerificationToken = nil
	targetUser.EmailVerificationExpiresAt = nil
	targetUser.UpdatedAt = time.Now()

	// 更新
	if err := u.userRepo.Update(ctx, targetUser); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	return nil
}
