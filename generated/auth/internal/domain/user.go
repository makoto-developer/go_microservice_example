package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleCustomer  Role = "CUSTOMER"
	RoleShopOwner Role = "SHOP_OWNER"
	RoleAdmin     Role = "ADMIN"
)

type User struct {
	ID                        uuid.UUID
	Email                     string
	PasswordHash              string
	Role                      Role
	EmailVerified             bool
	EmailVerificationToken    string
	EmailVerificationExpiresAt *time.Time
	PasswordResetToken        string
	PasswordResetExpiresAt    *time.Time
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
}

func NewUser(email, passwordHash string, role Role) *User {
	now := time.Now()
	return &User{
		ID:            uuid.New(),
		Email:         email,
		PasswordHash:  passwordHash,
		Role:          role,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
