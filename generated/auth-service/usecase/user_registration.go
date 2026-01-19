package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/auth-service/domain"
	"github.com/makoto-developer/go_microservice_example/manual/auth"
	"golang.org/x/crypto/bcrypt"
)

// UserRegistrationInput represents input for UserRegistration
type UserRegistrationInput struct {
	Email    string
	Password string
	Role     domain.Role
}

// UserRegistrationOutput represents output for UserRegistration
type UserRegistrationOutput struct {
	UserID                 uuid.UUID
	AccessToken            string
	RefreshToken           string
	EmailVerificationToken string
}

// UserRegistrationUsecase defines the interface for UserRegistration
type UserRegistrationUsecase interface {
	Execute(ctx context.Context, input UserRegistrationInput) (*UserRegistrationOutput, error)
}

type userRegistrationUsecaseImpl struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	jwtSecret        string
	emailSender      auth.EmailSender
}

// NewUserRegistrationUsecase creates a new instance
func NewUserRegistrationUsecase(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	jwtSecret string,
	emailSender auth.EmailSender,
) UserRegistrationUsecase {
	return &userRegistrationUsecaseImpl{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
		emailSender:      emailSender,
	}
}

// Execute executes UserRegistration
func (u *userRegistrationUsecaseImpl) Execute(ctx context.Context, input UserRegistrationInput) (*UserRegistrationOutput, error) {
	// メールアドレス重複チェック
	existing, _ := u.userRepo.FindByEmail(ctx, input.Email)
	if existing != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// パスワードハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// メール認証トークン生成
	verificationToken := generateRandomToken(32)
	verificationExpires := time.Now().Add(24 * time.Hour)

	// ユーザー作成
	user := &domain.User{
		Id:                         uuid.New(),
		Email:                      input.Email,
		PasswordHash:               string(hashedPassword),
		Role:                       input.Role,
		EmailVerified:              false,
		EmailVerificationToken:     &verificationToken,
		EmailVerificationExpiresAt: &verificationExpires,
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// JWT トークン生成（簡易実装）
	accessToken := fmt.Sprintf("access_%s", user.Id.String()[:8])
	refreshTokenStr := fmt.Sprintf("refresh_%s", user.Id.String()[:8])

	// Refresh Token 保存
	refreshTokenExpires := time.Now().Add(7 * 24 * time.Hour)
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

	// メール認証メールを送信
	if err := u.emailSender.SendVerificationEmail(ctx, input.Email, verificationToken); err != nil {
		// メール送信失敗はログのみ（ユーザー登録自体は成功）
		fmt.Printf("Warning: Failed to send verification email to %s: %v\n", input.Email, err)
	}

	return &UserRegistrationOutput{
		UserID:                 user.Id,
		AccessToken:            accessToken,
		RefreshToken:           refreshTokenStr,
		EmailVerificationToken: verificationToken,
	}, nil
}

func generateRandomToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}
