package grpc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/repository"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/auth/proto/customer_auth/v1"
)

// CustomerAuthHandler implements CustomerAuthService gRPC server
type CustomerAuthHandler struct {
	pb.UnimplementedCustomerAuthServiceServer
	registrationUsecase *usecase.CustomerRegistrationUsecase
	loginUsecase        *usecase.CustomerLoginUsecase
	customerUserRepo    repository.CustomerUserRepository
	jwtService          *usecase.JWTServiceV2
}

// NewCustomerAuthHandler creates a new customer auth handler
func NewCustomerAuthHandler(
	registrationUsecase *usecase.CustomerRegistrationUsecase,
	loginUsecase *usecase.CustomerLoginUsecase,
	customerUserRepo repository.CustomerUserRepository,
	jwtService *usecase.JWTServiceV2,
) *CustomerAuthHandler {
	return &CustomerAuthHandler{
		registrationUsecase: registrationUsecase,
		loginUsecase:        loginUsecase,
		customerUserRepo:    customerUserRepo,
		jwtService:          jwtService,
	}
}

// Register registers a new customer
func (h *CustomerAuthHandler) Register(
	ctx context.Context,
	req *pb.CustomerRegisterRequest,
) (*pb.CustomerRegisterResponse, error) {
	// Validate request
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Execute registration
	userID, accessToken, refreshToken, err := h.registrationUsecase.Execute(
		ctx,
		req.Email,
		req.Password,
	)
	if err != nil {
		// Check for specific errors
		if err.Error() == "email already registered as customer" {
			return nil, status.Error(codes.AlreadyExists, "email already registered")
		}
		if err.Error() == "password must be at least 8 characters" {
			return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to register customer: %v", err))
	}

	return &pb.CustomerRegisterResponse{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login authenticates a customer
func (h *CustomerAuthHandler) Login(
	ctx context.Context,
	req *pb.CustomerLoginRequest,
) (*pb.CustomerLoginResponse, error) {
	// Validate request
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Execute login
	userID, accessToken, refreshToken, err := h.loginUsecase.Execute(
		ctx,
		req.Email,
		req.Password,
	)
	if err != nil {
		// Check for specific errors
		if err.Error() == "invalid email or password" {
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		if err.Error() == "email not verified" {
			return nil, status.Error(codes.FailedPrecondition, "email not verified")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to login: %v", err))
	}

	return &pb.CustomerLoginResponse{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout logs out a customer
func (h *CustomerAuthHandler) Logout(
	ctx context.Context,
	req *pb.CustomerLogoutRequest,
) (*pb.CustomerLogoutResponse, error) {
	// Validate request
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	// TODO: Implement token revocation (e.g., add to blacklist in Redis)
	// For now, we just return success as tokens will expire naturally

	return &pb.CustomerLogoutResponse{
		Success: true,
	}, nil
}

// VerifyEmail verifies customer's email address
func (h *CustomerAuthHandler) VerifyEmail(
	ctx context.Context,
	req *pb.CustomerVerifyEmailRequest,
) (*pb.CustomerVerifyEmailResponse, error) {
	// Validate request
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Find user by verification token
	user, err := h.customerUserRepo.FindByVerificationToken(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find user: %v", err))
	}
	if user == nil {
		return &pb.CustomerVerifyEmailResponse{
			Success: false,
			Message: "Invalid or expired verification token",
		}, nil
	}

	// Verify email
	user.VerifyEmail()

	// Update user
	if err := h.customerUserRepo.Update(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update user: %v", err))
	}

	return &pb.CustomerVerifyEmailResponse{
		Success: true,
		Message: "Email verified successfully",
	}, nil
}

// RequestPasswordReset initiates password reset process
func (h *CustomerAuthHandler) RequestPasswordReset(
	ctx context.Context,
	req *pb.CustomerRequestPasswordResetRequest,
) (*pb.CustomerRequestPasswordResetResponse, error) {
	// Validate request
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Find user by email
	user, err := h.customerUserRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find user: %v", err))
	}

	// Always return success to prevent email enumeration
	if user == nil {
		return &pb.CustomerRequestPasswordResetResponse{
			Success: true,
			Message: "If the email exists, a password reset link has been sent",
		}, nil
	}

	// Generate reset token
	token, err := generateRandomToken(32)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate reset token")
	}

	// Set reset token
	user.SetPasswordResetToken(token, time.Now().Add(1*time.Hour))

	// Update user
	if err := h.customerUserRepo.Update(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update user: %v", err))
	}

	// TODO: Send password reset email (async)
	// go h.emailService.SendPasswordResetEmail(user.Email, token)

	return &pb.CustomerRequestPasswordResetResponse{
		Success: true,
		Message: "If the email exists, a password reset link has been sent",
	}, nil
}

// ResetPassword resets customer's password
func (h *CustomerAuthHandler) ResetPassword(
	ctx context.Context,
	req *pb.CustomerResetPasswordRequest,
) (*pb.CustomerResetPasswordResponse, error) {
	// Validate request
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new_password is required")
	}
	if len(req.NewPassword) < 8 {
		return &pb.CustomerResetPasswordResponse{
			Success: false,
			Message: "Password must be at least 8 characters",
		}, nil
	}

	// Find user by reset token
	user, err := h.customerUserRepo.FindByResetToken(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find user: %v", err))
	}
	if user == nil {
		return &pb.CustomerResetPasswordResponse{
			Success: false,
			Message: "Invalid or expired reset token",
		}, nil
	}

	// Reset password(bcrypt でハッシュ化して保存)
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}
	user.ResetPassword(string(hashed))

	// Update user
	if err := h.customerUserRepo.Update(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update user: %v", err))
	}

	return &pb.CustomerResetPasswordResponse{
		Success: true,
		Message: "Password reset successfully",
	}, nil
}

// RefreshToken refreshes access token
func (h *CustomerAuthHandler) RefreshToken(
	ctx context.Context,
	req *pb.CustomerRefreshTokenRequest,
) (*pb.CustomerRefreshTokenResponse, error) {
	// Validate request
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateCustomerRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	// Find user
	claimUserID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token subject")
	}
	user, err := h.customerUserRepo.FindByID(ctx, claimUserID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find user: %v", err))
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Generate new tokens
	accessToken, err := h.jwtService.GenerateCustomerAccessToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}

	refreshToken, err := h.jwtService.GenerateCustomerRefreshToken(user.ID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	return &pb.CustomerRefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateRandomToken generates a random hex token
func generateRandomToken(length int) (string, error) {
	// This is duplicated from usecase, consider extracting to a shared utility
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
