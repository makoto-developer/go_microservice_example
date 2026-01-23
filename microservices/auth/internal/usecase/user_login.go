package usecase

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	
	"github.com/makoto-developer/go_microservice_example/generated/auth/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/auth/internal/repository"
)

type UserLoginUsecase struct {
	userRepo   repository.UserRepository
	jwtService *JWTService
}

func NewUserLoginUsecase(userRepo repository.UserRepository, jwtService *JWTService) *UserLoginUsecase {
	return &UserLoginUsecase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (u *UserLoginUsecase) Execute(ctx context.Context, email, password string) (userID, accessToken, refreshToken string, role domain.Role, err error) {
	// Find user by email
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", "", "", fmt.Errorf("invalid email or password")
	}
	
	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", "", "", fmt.Errorf("invalid email or password")
	}
	
	// Generate JWT tokens
	accessToken, err = u.jwtService.GenerateAccessToken(user.ID.String(), string(user.Role))
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	
	refreshToken, err = u.jwtService.GenerateRefreshToken(user.ID.String())
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	
	return user.ID.String(), accessToken, refreshToken, user.Role, nil
}
