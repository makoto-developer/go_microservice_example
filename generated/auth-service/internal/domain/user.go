package domain

import (
	"time"
)

type UserRole string

const (
	RoleCustomer   UserRole = "customer"
	RoleShopOwner  UserRole = "shop_owner"
	RoleAdmin      UserRole = "admin"
)

type User struct {
	UserID        string
	Email         string
	PasswordHash  string
	Role          UserRole
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type RefreshToken struct {
	TokenID   string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	Revoked   bool
	RevokedAt *time.Time
}

type EmailVerificationToken struct {
	TokenID   string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	Used      bool
	UsedAt    *time.Time
}

type PasswordResetToken struct {
	TokenID   string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	Used      bool
	UsedAt    *time.Time
}

func (r UserRole) IsValid() bool {
	switch r {
	case RoleCustomer, RoleShopOwner, RoleAdmin:
		return true
	}
	return false
}

func (r UserRole) String() string {
	return string(r)
}
