package usecase

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/repository"
)

// CustomerLoginUsecase handles customer user login
type CustomerLoginUsecase struct {
	customerUserRepo repository.CustomerUserRepository
	jwtService       *JWTServiceV2
}

// NewCustomerLoginUsecase creates a new customer login usecase
func NewCustomerLoginUsecase(
	customerUserRepo repository.CustomerUserRepository,
	jwtService *JWTServiceV2,
) *CustomerLoginUsecase {
	return &CustomerLoginUsecase{
		customerUserRepo: customerUserRepo,
		jwtService:       jwtService,
	}
}

// Execute performs customer login
func (u *CustomerLoginUsecase) Execute(
	ctx context.Context,
	email, password string,
) (userID, accessToken, refreshToken string, err error) {
	// Find user by email
	user, err := u.customerUserRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return "", "", "", fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", "", fmt.Errorf("invalid email or password")
	}

	// Check if email is verified
	if !user.EmailVerified {
		return "", "", "", fmt.Errorf("email not verified")
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

	return user.ID.String(), accessToken, refreshToken, nil
}
