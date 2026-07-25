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

	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/repository"
	"github.com/makoto-developer/go_microservice_example/microservices/auth/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/auth/proto/owner_auth/v1"
)

// OwnerAuthHandler implements OwnerAuthService gRPC server
type OwnerAuthHandler struct {
	pb.UnimplementedOwnerAuthServiceServer
	registrationUsecase *usecase.OwnerRegistrationUsecase
	loginUsecase        *usecase.OwnerLoginUsecase
	ownerUserRepo       repository.OwnerUserRepository
	jwtService          *usecase.JWTServiceV2
}

// NewOwnerAuthHandler creates a new owner auth handler
func NewOwnerAuthHandler(
	registrationUsecase *usecase.OwnerRegistrationUsecase,
	loginUsecase *usecase.OwnerLoginUsecase,
	ownerUserRepo repository.OwnerUserRepository,
	jwtService *usecase.JWTServiceV2,
) *OwnerAuthHandler {
	return &OwnerAuthHandler{
		registrationUsecase: registrationUsecase,
		loginUsecase:        loginUsecase,
		ownerUserRepo:       ownerUserRepo,
		jwtService:          jwtService,
	}
}

// Register registers a new shop owner
func (h *OwnerAuthHandler) Register(
	ctx context.Context,
	req *pb.OwnerRegisterRequest,
) (*pb.OwnerRegisterResponse, error) {
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
		if err.Error() == "email already registered as owner" {
			return nil, status.Error(codes.AlreadyExists, "email already registered")
		}
		if err.Error() == "password must be at least 8 characters" {
			return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to register owner: %v", err))
	}

	return &pb.OwnerRegisterResponse{
		UserId:                     userID,
		AccessToken:                accessToken,
		RefreshToken:               refreshToken,
		BusinessVerificationStatus: string(domain.BusinessVerificationPending),
	}, nil
}

// Login authenticates a shop owner
func (h *OwnerAuthHandler) Login(
	ctx context.Context,
	req *pb.OwnerLoginRequest,
) (*pb.OwnerLoginResponse, error) {
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

	// Get user to check business verification status
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid user id")
	}
	user, err := h.ownerUserRepo.FindByID(ctx, parsedUserID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find user: %v", err))
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.OwnerLoginResponse{
		UserId:                     userID,
		AccessToken:                accessToken,
		RefreshToken:               refreshToken,
		BusinessVerified:           user.BusinessVerified,
		BusinessVerificationStatus: string(user.BusinessVerificationStatus),
	}, nil
}

// Logout logs out a shop owner
func (h *OwnerAuthHandler) Logout(
	ctx context.Context,
	req *pb.OwnerLogoutRequest,
) (*pb.OwnerLogoutResponse, error) {
	// Validate request
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	// TODO: Implement token revocation (e.g., add to blacklist in Redis)
	// For now, we just return success as tokens will expire naturally

	return &pb.OwnerLogoutResponse{
		Success: true,
	}, nil
}

// VerifyEmail verifies owner's email address
func (h *OwnerAuthHandler) VerifyEmail(
	ctx context.Context,
	req *pb.OwnerVerifyEmailRequest,
) (*pb.OwnerVerifyEmailResponse, error) {
	// Validate request
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// Find user by verification token
	user, err := h.ownerUserRepo.FindByVerificationToken(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find user: %v", err))
	}
	if user == nil {
		return &pb.OwnerVerifyEmailResponse{
			Success: false,
			Message: "Invalid or expired verification token",
		}, nil
	}

	// Verify email
	user.VerifyEmail()

	// Update user
	if err := h.ownerUserRepo.Update(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update user: %v", err))
	}

	return &pb.OwnerVerifyEmailResponse{
		Success: true,
		Message: "Email verified successfully",
	}, nil
}

// RequestPasswordReset initiates password reset process
func (h *OwnerAuthHandler) RequestPasswordReset(
	ctx context.Context,
	req *pb.OwnerRequestPasswordResetRequest,
) (*pb.OwnerRequestPasswordResetResponse, error) {
	// Validate request
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Find user by email
	user, err := h.ownerUserRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find user: %v", err))
	}

	// Always return success to prevent email enumeration
	if user == nil {
		return &pb.OwnerRequestPasswordResetResponse{
			Success: true,
			Message: "If the email exists, a password reset link has been sent",
		}, nil
	}

	// Generate reset token
	token, err := generateOwnerRandomToken(32)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate reset token")
	}

	// Set reset token
	user.SetPasswordResetToken(token, time.Now().Add(1*time.Hour))

	// Update user
	if err := h.ownerUserRepo.Update(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update user: %v", err))
	}

	// TODO: Send password reset email (async)
	// go h.emailService.SendPasswordResetEmail(user.Email, token)

	return &pb.OwnerRequestPasswordResetResponse{
		Success: true,
		Message: "If the email exists, a password reset link has been sent",
	}, nil
}

// ResetPassword resets owner's password
func (h *OwnerAuthHandler) ResetPassword(
	ctx context.Context,
	req *pb.OwnerResetPasswordRequest,
) (*pb.OwnerResetPasswordResponse, error) {
	// Validate request
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	if req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "new_password is required")
	}
	if len(req.NewPassword) < 8 {
		return &pb.OwnerResetPasswordResponse{
			Success: false,
			Message: "Password must be at least 8 characters",
		}, nil
	}

	// Find user by reset token
	user, err := h.ownerUserRepo.FindByResetToken(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find user: %v", err))
	}
	if user == nil {
		return &pb.OwnerResetPasswordResponse{
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
	if err := h.ownerUserRepo.Update(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update user: %v", err))
	}

	return &pb.OwnerResetPasswordResponse{
		Success: true,
		Message: "Password reset successfully",
	}, nil
}

// RefreshToken refreshes access token
func (h *OwnerAuthHandler) RefreshToken(
	ctx context.Context,
	req *pb.OwnerRefreshTokenRequest,
) (*pb.OwnerRefreshTokenResponse, error) {
	// Validate request
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateOwnerRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	// Find user
	claimUserID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token subject")
	}
	user, err := h.ownerUserRepo.FindByID(ctx, claimUserID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find user: %v", err))
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Generate new tokens
	accessToken, err := h.jwtService.GenerateOwnerAccessToken(user.ID.String(), user.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}

	refreshToken, err := h.jwtService.GenerateOwnerRefreshToken(user.ID.String())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	return &pb.OwnerRefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetBusinessVerificationStatus gets business verification status
func (h *OwnerAuthHandler) GetBusinessVerificationStatus(
	ctx context.Context,
	req *pb.OwnerGetBusinessVerificationStatusRequest,
) (*pb.OwnerGetBusinessVerificationStatusResponse, error) {
	// Validate request
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Find user
	reqUserID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}
	user, err := h.ownerUserRepo.FindByID(ctx, reqUserID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to find user: %v", err))
	}
	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.OwnerGetBusinessVerificationStatusResponse{
		BusinessVerified:           user.BusinessVerified,
		BusinessVerificationStatus: string(user.BusinessVerificationStatus),
	}, nil
}

// generateOwnerRandomToken generates a random hex token
func generateOwnerRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
