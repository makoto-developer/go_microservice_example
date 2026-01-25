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

// CustomerRegistrationUsecase handles customer user registration
type CustomerRegistrationUsecase struct {
	customerUserRepo repository.CustomerUserRepository
	jwtService       *JWTServiceV2
	emailService     EmailService
}

// EmailService defines the interface for sending emails
type EmailService interface {
	SendVerificationEmail(email, token string) error
}

// NewCustomerRegistrationUsecase creates a new customer registration usecase
func NewCustomerRegistrationUsecase(
	customerUserRepo repository.CustomerUserRepository,
	jwtService *JWTServiceV2,
	emailService EmailService,
) *CustomerRegistrationUsecase {
	return &CustomerRegistrationUsecase{
		customerUserRepo: customerUserRepo,
		jwtService:       jwtService,
		emailService:     emailService,
	}
}

// Execute registers a new customer
func (u *CustomerRegistrationUsecase) Execute(
	ctx context.Context,
	email, password string,
) (userID, accessToken, refreshToken string, err error) {
	// Check if email already exists in customer_users table
	existingUser, err := u.customerUserRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return "", "", "", fmt.Errorf("email already registered as customer")
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

	// Create customer user
	user := domain.NewCustomerUser(email, string(hashedPassword))

	// Generate email verification token
	token, err := generateRandomToken(32)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate verification token: %w", err)
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	user.SetEmailVerificationToken(token, expiresAt)

	// Save user to database
	if err := u.customerUserRepo.Create(ctx, user); err != nil {
		return "", "", "", fmt.Errorf("failed to create customer user: %w", err)
	}

	// Generate JWT tokens
	accessToken, err = u.jwtService.GenerateCustomerAccessToken(user.ID.String(), user.Email)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err = u.jwtService.GenerateCustomerRefreshToken(user.ID.String())
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Send verification email (async, don't block on failure)
	go func() {
		_ = u.emailService.SendVerificationEmail(email, token)
	}()

	return user.ID.String(), accessToken, refreshToken, nil
}

// generateRandomToken generates a random hex token
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
