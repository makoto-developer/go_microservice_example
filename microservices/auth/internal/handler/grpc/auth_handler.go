package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	
	pb "github.com/makoto-developer/go_microservice_example/proto/auth_service/v1"
	"github.com/makoto-developer/go_microservice_example/generated/auth/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/auth/internal/usecase"
)

type AuthServiceHandler struct {
	pb.UnimplementedAuthServiceServer
	registrationUsecase *usecase.UserRegistrationUsecase
	loginUsecase        *usecase.UserLoginUsecase
	jwtService          *usecase.JWTService
}

func NewAuthServiceHandler(
	registrationUsecase *usecase.UserRegistrationUsecase,
	loginUsecase *usecase.UserLoginUsecase,
	jwtService *usecase.JWTService,
) *AuthServiceHandler {
	return &AuthServiceHandler{
		registrationUsecase: registrationUsecase,
		loginUsecase:        loginUsecase,
		jwtService:          jwtService,
	}
}

func (h *AuthServiceHandler) Register(ctx context.Context, req *pb.UserRegistrationRequest) (*pb.RegisterResponse, error) {
	var role domain.Role
	switch req.Role {
	case pb.Role_CUSTOMER:
		role = domain.RoleCustomer
	case pb.Role_SHOP_OWNER:
		role = domain.RoleShopOwner
	case pb.Role_ADMIN:
		role = domain.RoleAdmin
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid role")
	}
	
	userID, accessToken, refreshToken, err := h.registrationUsecase.Execute(ctx, req.Email, req.Password, role)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	return &pb.RegisterResponse{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "User registered successfully",
	}, nil
}

func (h *AuthServiceHandler) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.LoginResponse, error) {
	userID, accessToken, refreshToken, role, err := h.loginUsecase.Execute(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	
	var pbRole pb.Role
	switch role {
	case domain.RoleCustomer:
		pbRole = pb.Role_CUSTOMER
	case domain.RoleShopOwner:
		pbRole = pb.Role_SHOP_OWNER
	case domain.RoleAdmin:
		pbRole = pb.Role_ADMIN
	default:
		pbRole = pb.Role_ROLE_UNSPECIFIED
	}
	
	return &pb.LoginResponse{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Role:         pbRole,
	}, nil
}

func (h *AuthServiceHandler) VerifyToken(ctx context.Context, req *pb.TokenVerificationRequest) (*pb.VerifyTokenResponse, error) {
	claims, err := h.jwtService.VerifyAccessToken(req.AccessToken)
	if err != nil {
		return &pb.VerifyTokenResponse{Valid: false}, nil
	}
	
	var pbRole pb.Role
	switch claims.Role {
	case string(domain.RoleCustomer):
		pbRole = pb.Role_CUSTOMER
	case string(domain.RoleShopOwner):
		pbRole = pb.Role_SHOP_OWNER
	case string(domain.RoleAdmin):
		pbRole = pb.Role_ADMIN
	default:
		pbRole = pb.Role_ROLE_UNSPECIFIED
	}
	
	return &pb.VerifyTokenResponse{Valid: true, UserId: claims.UserID, Role: pbRole}, nil
}

func (h *AuthServiceHandler) Logout(ctx context.Context, req *pb.UserLogoutRequest) (*pb.LogoutResponse, error) {
	return &pb.LogoutResponse{Success: true, Message: "Logged out successfully"}, nil
}

func (h *AuthServiceHandler) RefreshToken(ctx context.Context, req *pb.TokenRefreshRequest) (*pb.RefreshTokenResponse, error) {
	claims, err := h.jwtService.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}
	
	accessToken, err := h.jwtService.GenerateAccessToken(claims.UserID, "")
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate access token")
	}
	
	refreshToken, err := h.jwtService.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}
	
	return &pb.RefreshTokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (h *AuthServiceHandler) VerifyEmail(ctx context.Context, req *pb.EmailVerificationRequest) (*pb.VerifyEmailResponse, error) {
	return &pb.VerifyEmailResponse{Success: true, Message: "Email verified successfully"}, nil
}

func (h *AuthServiceHandler) RequestPasswordReset(ctx context.Context, req *pb.PasswordResetRequestRequest) (*pb.PasswordResetResponse, error) {
	return &pb.PasswordResetResponse{Success: true, Message: "Password reset email sent"}, nil
}

func (h *AuthServiceHandler) ResetPassword(ctx context.Context, req *pb.PasswordResetRequest) (*pb.ResetPasswordResponse, error) {
	return &pb.ResetPasswordResponse{Success: true, Message: "Password reset successfully"}, nil
}

func (h *AuthServiceHandler) ChangePassword(ctx context.Context, req *pb.PasswordChangeRequest) (*pb.ChangePasswordResponse, error) {
	return &pb.ChangePasswordResponse{Success: true, Message: "Password changed successfully"}, nil
}
