package usecase

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserType represents the type of user (customer or owner)
type UserType string

const (
	UserTypeCustomer UserType = "customer"
	UserTypeOwner    UserType = "owner"
)

// JWTServiceV2 provides JWT token generation and validation with user type separation
type JWTServiceV2 struct {
	accessTokenSecret  string
	refreshTokenSecret string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// NewJWTServiceV2 creates a new JWT service with user type support
func NewJWTServiceV2(accessSecret, refreshSecret string) *JWTServiceV2 {
	return &JWTServiceV2{
		accessTokenSecret:  accessSecret,
		refreshTokenSecret: refreshSecret,
		accessTokenExpiry:  1 * time.Hour,       // 1 hour for access token
		refreshTokenExpiry: 30 * 24 * time.Hour, // 30 days for refresh token
	}
}

// AccessTokenClaimsV2 contains claims for access tokens with user type
type AccessTokenClaimsV2 struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	UserType UserType `json:"user_type"` // "customer" or "owner"
	jwt.RegisteredClaims
}

// RefreshTokenClaimsV2 contains claims for refresh tokens with user type
type RefreshTokenClaimsV2 struct {
	UserID   string   `json:"user_id"`
	UserType UserType `json:"user_type"` // "customer" or "owner"
	jwt.RegisteredClaims
}

// GenerateCustomerAccessToken generates an access token for a customer
func (s *JWTServiceV2) GenerateCustomerAccessToken(userID, email string) (string, error) {
	return s.generateAccessToken(userID, email, UserTypeCustomer)
}

// GenerateOwnerAccessToken generates an access token for an owner
func (s *JWTServiceV2) GenerateOwnerAccessToken(userID, email string) (string, error) {
	return s.generateAccessToken(userID, email, UserTypeOwner)
}

// generateAccessToken generates an access token with the specified user type
func (s *JWTServiceV2) generateAccessToken(userID, email string, userType UserType) (string, error) {
	claims := AccessTokenClaimsV2{
		UserID:   userID,
		Email:    email,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Audience:  jwt.ClaimStrings{string(userType)}, // "customer" or "owner"
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.accessTokenSecret))
}

// GenerateCustomerRefreshToken generates a refresh token for a customer
func (s *JWTServiceV2) GenerateCustomerRefreshToken(userID string) (string, error) {
	return s.generateRefreshToken(userID, UserTypeCustomer)
}

// GenerateOwnerRefreshToken generates a refresh token for an owner
func (s *JWTServiceV2) GenerateOwnerRefreshToken(userID string) (string, error) {
	return s.generateRefreshToken(userID, UserTypeOwner)
}

// generateRefreshToken generates a refresh token with the specified user type
func (s *JWTServiceV2) generateRefreshToken(userID string, userType UserType) (string, error) {
	claims := RefreshTokenClaimsV2{
		UserID:   userID,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Audience:  jwt.ClaimStrings{string(userType)},
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshTokenSecret))
}

// ValidateAccessToken validates an access token and returns the claims
func (s *JWTServiceV2) ValidateAccessToken(tokenString string) (*AccessTokenClaimsV2, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaimsV2{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.accessTokenSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*AccessTokenClaimsV2); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken validates a refresh token and returns the claims
func (s *JWTServiceV2) ValidateRefreshToken(tokenString string) (*RefreshTokenClaimsV2, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaimsV2{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.refreshTokenSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*RefreshTokenClaimsV2); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// ValidateCustomerToken validates that a token is for a customer
func (s *JWTServiceV2) ValidateCustomerToken(tokenString string) (*AccessTokenClaimsV2, error) {
	claims, err := s.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.UserType != UserTypeCustomer {
		return nil, fmt.Errorf("invalid user type: expected customer, got %s", claims.UserType)
	}

	// Check audience
	if !audienceContains(claims.Audience, string(UserTypeCustomer)) {
		return nil, fmt.Errorf("invalid audience")
	}

	return claims, nil
}

// ValidateOwnerToken validates that a token is for an owner
func (s *JWTServiceV2) ValidateOwnerToken(tokenString string) (*AccessTokenClaimsV2, error) {
	claims, err := s.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.UserType != UserTypeOwner {
		return nil, fmt.Errorf("invalid user type: expected owner, got %s", claims.UserType)
	}

	// Check audience
	if !audienceContains(claims.Audience, string(UserTypeOwner)) {
		return nil, fmt.Errorf("invalid audience")
	}

	return claims, nil
}

// audienceContains は aud クレームに期待する値が含まれるかを返す(jwt v5 には VerifyAudience が無い)。
func audienceContains(aud jwt.ClaimStrings, expected string) bool {
	for _, a := range aud {
		if a == expected {
			return true
		}
	}
	return false
}

// ValidateCustomerRefreshToken は顧客のリフレッシュトークンを検証する。
func (s *JWTServiceV2) ValidateCustomerRefreshToken(tokenString string) (*RefreshTokenClaimsV2, error) {
	claims, err := s.ValidateRefreshToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.UserType != UserTypeCustomer {
		return nil, fmt.Errorf("invalid user type: expected customer, got %s", claims.UserType)
	}
	return claims, nil
}

// ValidateOwnerRefreshToken はオーナーのリフレッシュトークンを検証する。
func (s *JWTServiceV2) ValidateOwnerRefreshToken(tokenString string) (*RefreshTokenClaimsV2, error) {
	claims, err := s.ValidateRefreshToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.UserType != UserTypeOwner {
		return nil, fmt.Errorf("invalid user type: expected owner, got %s", claims.UserType)
	}
	return claims, nil
}
