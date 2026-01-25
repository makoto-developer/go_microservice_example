package usecase

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/repository"
)

// OwnerLoginUsecase handles owner user login
type OwnerLoginUsecase struct {
	ownerUserRepo repository.OwnerUserRepository
	jwtService    *JWTServiceV2
}

// NewOwnerLoginUsecase creates a new owner login usecase
func NewOwnerLoginUsecase(
	ownerUserRepo repository.OwnerUserRepository,
	jwtService *JWTServiceV2,
) *OwnerLoginUsecase {
	return &OwnerLoginUsecase{
		ownerUserRepo: ownerUserRepo,
		jwtService:    jwtService,
	}
}

// Execute performs owner login
func (u *OwnerLoginUsecase) Execute(
	ctx context.Context,
	email, password string,
) (userID, accessToken, refreshToken string, err error) {
	// Find user by email
	user, err := u.ownerUserRepo.FindByEmail(ctx, email)
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
	accessToken, err = u.jwtService.GenerateOwnerAccessToken(user.ID.String(), user.Email)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err = u.jwtService.GenerateOwnerRefreshToken(user.ID.String())
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return user.ID.String(), accessToken, refreshToken, nil
}
