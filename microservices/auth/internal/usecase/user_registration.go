package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/repository"
)

type UserRegistrationUsecase struct {
	userRepo   repository.UserRepository
	jwtService *JWTService
}

func NewUserRegistrationUsecase(userRepo repository.UserRepository, jwtService *JWTService) *UserRegistrationUsecase {
	return &UserRegistrationUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (u *UserRegistrationUsecase) Execute(ctx context.Context, email, password string, role domain.Role) (userID, accessToken, refreshToken string, err error) {
	// Check if email already exists
	existingUser, _ := u.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return "", "", "", fmt.Errorf("email already exists")
	}

	// Validate password strength
	if len(password) < 8 {
		return "", "", "", fmt.Errorf("password must be at least 8 characters")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := domain.NewUser(email, string(hashedPassword), role)

	// Generate email verification token
	token, err := generateRandomToken(32)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate verification token: %w", err)
	}
	user.EmailVerificationToken = token
	expiresAt := time.Now().Add(24 * time.Hour)
	user.EmailVerificationExpiresAt = &expiresAt

	// Save user
	if err := u.userRepo.Create(ctx, user); err != nil {
		return "", "", "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT tokens
	accessToken, err = u.jwtService.GenerateAccessToken(user.ID.String(), string(user.Role))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err = u.jwtService.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return user.ID.String(), accessToken, refreshToken, nil
}

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
