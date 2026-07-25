package domain

import (
	"time"

	"github.com/google/uuid"
)

// CustomerUser represents a customer authentication account
// Separate from OwnerUser - same email can be used for both
type CustomerUser struct {
	ID                         uuid.UUID
	Email                      string
	PasswordHash               string
	EmailVerified              bool
	EmailVerificationToken     string
	EmailVerificationExpiresAt *time.Time
	PasswordResetToken         string
	PasswordResetExpiresAt     *time.Time
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

// NewCustomerUser creates a new customer user
func NewCustomerUser(email, passwordHash string) *CustomerUser {
	now := time.Now()
	return &CustomerUser{
		ID:            uuid.New(),
		Email:         email,
		PasswordHash:  passwordHash,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// SetEmailVerificationToken sets the email verification token and expiry
func (u *CustomerUser) SetEmailVerificationToken(token string, expiresAt time.Time) {
	u.EmailVerificationToken = token
	u.EmailVerificationExpiresAt = &expiresAt
	u.UpdatedAt = time.Now()
}

// SetPasswordResetToken sets the password reset token and expiry
func (u *CustomerUser) SetPasswordResetToken(token string, expiresAt time.Time) {
	u.PasswordResetToken = token
	u.PasswordResetExpiresAt = &expiresAt
	u.UpdatedAt = time.Now()
}

// VerifyEmail marks the email as verified
func (u *CustomerUser) VerifyEmail() {
	u.EmailVerified = true
	u.EmailVerificationToken = ""
	u.EmailVerificationExpiresAt = nil
	u.UpdatedAt = time.Now()
}

// ResetPassword updates the password hash and clears reset token
func (u *CustomerUser) ResetPassword(newPasswordHash string) {
	u.PasswordHash = newPasswordHash
	u.PasswordResetToken = ""
	u.PasswordResetExpiresAt = nil
	u.UpdatedAt = time.Now()
}

// IsEmailVerificationExpired checks if the email verification token is expired
func (u *CustomerUser) IsEmailVerificationExpired() bool {
	if u.EmailVerificationExpiresAt == nil {
		return true
	}
	return time.Now().After(*u.EmailVerificationExpiresAt)
}

// IsPasswordResetExpired checks if the password reset token is expired
func (u *CustomerUser) IsPasswordResetExpired() bool {
	if u.PasswordResetExpiresAt == nil {
		return true
	}
	return time.Now().After(*u.PasswordResetExpiresAt)
}
