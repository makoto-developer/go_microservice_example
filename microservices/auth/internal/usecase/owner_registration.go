package usecase

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/repository"
)

// OwnerRegistrationUsecase handles owner user registration
type OwnerRegistrationUsecase struct {
	ownerUserRepo repository.OwnerUserRepository
	jwtService    *JWTServiceV2
	emailService  EmailService
}

// NewOwnerRegistrationUsecase creates a new owner registration usecase
func NewOwnerRegistrationUsecase(
	ownerUserRepo repository.OwnerUserRepository,
	jwtService *JWTServiceV2,
	emailService EmailService,
) *OwnerRegistrationUsecase {
	return &OwnerRegistrationUsecase{
		ownerUserRepo: ownerUserRepo,
		jwtService:    jwtService,
		emailService:  emailService,
	}
}

// Execute registers a new owner
func (u *OwnerRegistrationUsecase) Execute(
	ctx context.Context,
	email, password string,
) (userID, accessToken, refreshToken string, err error) {
	// Check if email already exists in owner_users table
	existingUser, err := u.ownerUserRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return "", "", "", fmt.Errorf("email already registered as owner")
	}

	// Validate password strength (minimum 8 characters)
	if len(password) < 8 {
		return "", "", "", fmt.Errorf("password must be at least 8 characters")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create owner user
	user := domain.NewOwnerUser(email, string(hashedPassword))

	// Generate email verification token
	token, err := generateRandomToken(32)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate verification token: %w", err)
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	user.SetEmailVerificationToken(token, expiresAt)

	// Save user to database
	if err := u.ownerUserRepo.Create(ctx, user); err != nil {
		return "", "", "", fmt.Errorf("failed to create owner user: %w", err)
	}

	// Generate JWT tokens
	accessToken, err = u.jwtService.GenerateOwnerAccessToken(user.ID.String(), user.Email)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err = u.jwtService.GenerateOwnerRefreshToken(user.ID.String())
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Send verification email (async, don't block on failure)
	go func() {
		_ = u.emailService.SendVerificationEmail(email, token)
	}()

	return user.ID.String(), accessToken, refreshToken, nil
}
