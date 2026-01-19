package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/auth-service/domain"
)

// TokenVerificationInput represents input for TokenVerification
type TokenVerificationInput struct {
	AccessToken string
}

// TokenVerificationOutput represents output for TokenVerification
type TokenVerificationOutput struct {
	UserID uuid.UUID
	Email  string
	Role   domain.Role
}

// TokenVerificationUsecase defines the interface for TokenVerification
type TokenVerificationUsecase interface {
	Execute(ctx context.Context, input TokenVerificationInput) (*TokenVerificationOutput, error)
}

type tokenVerificationUsecaseImpl struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

// NewTokenVerificationUsecase creates a new instance
func NewTokenVerificationUsecase(
	userRepo domain.UserRepository,
	jwtSecret string,
) TokenVerificationUsecase {
	return &tokenVerificationUsecaseImpl{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

// Execute executes TokenVerification
func (u *tokenVerificationUsecaseImpl) Execute(ctx context.Context, input TokenVerificationInput) (*TokenVerificationOutput, error) {
	// 簡易実装: トークンからユーザーIDを抽出
	// 実際にはJWT検証ロジックが必要
	if !strings.HasPrefix(input.AccessToken, "access_") {
		return nil, fmt.Errorf("invalid access token format")
	}

	// トークンからユーザーIDを抽出（簡易実装）
	// 実際の実装では、JWTをパース・検証してクレームからユーザーIDを取得
	userIDStr := strings.TrimPrefix(input.AccessToken, "access_")
	// この時点では完全なUUIDではないため、ダミーIDを使用
	// 実際にはJWTペイロードに含まれるユーザーIDを使用

	// とりあえず、全ユーザーから検索（実際にはJWTから取得したIDで検索）
	users, err := u.userRepo.List(ctx, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	// 最初のユーザーを返す（実際にはJWTから取得したIDで検索）
	if len(users) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	// トークンプレフィックスに含まれるIDの一部と照合
	var targetUser *domain.User
	for _, user := range users {
		if strings.HasPrefix(user.Id.String(), userIDStr) {
			targetUser = user
			break
		}
	}

	if targetUser == nil {
		return nil, fmt.Errorf("invalid token: user not found")
	}

	return &TokenVerificationOutput{
		UserID: targetUser.Id,
		Email:  targetUser.Email,
		Role:   targetUser.Role,
	}, nil
}
