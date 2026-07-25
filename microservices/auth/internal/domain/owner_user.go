package domain

import (
	"time"

	"github.com/google/uuid"
)

// BusinessVerificationStatus represents the business verification status
type BusinessVerificationStatus string

const (
	BusinessVerificationPending  BusinessVerificationStatus = "pending"
	BusinessVerificationApproved BusinessVerificationStatus = "approved"
	BusinessVerificationRejected BusinessVerificationStatus = "rejected"
)

// OwnerUser represents a shop owner authentication account
// Separate from CustomerUser - same email can be used for both
type OwnerUser struct {
	ID                         uuid.UUID
	Email                      string
	PasswordHash               string
	EmailVerified              bool
	EmailVerificationToken     string
	EmailVerificationExpiresAt *time.Time
	PasswordResetToken         string
	PasswordResetExpiresAt     *time.Time
	BusinessVerified           bool
	BusinessVerificationStatus BusinessVerificationStatus
	CreatedAt                  time.Time
	UpdatedAt                  time.Time
}

// NewOwnerUser creates a new owner user
func NewOwnerUser(email, passwordHash string) *OwnerUser {
	now := time.Now()
	return &OwnerUser{
		ID:                         uuid.New(),
		Email:                      email,
		PasswordHash:               passwordHash,
		EmailVerified:              false,
		BusinessVerified:           false,
		BusinessVerificationStatus: BusinessVerificationPending,
		CreatedAt:                  now,
		UpdatedAt:                  now,
	}
}

// SetEmailVerificationToken sets the email verification token and expiry
func (u *OwnerUser) SetEmailVerificationToken(token string, expiresAt time.Time) {
	u.EmailVerificationToken = token
	u.EmailVerificationExpiresAt = &expiresAt
	u.UpdatedAt = time.Now()
}

// SetPasswordResetToken sets the password reset token and expiry
func (u *OwnerUser) SetPasswordResetToken(token string, expiresAt time.Time) {
	u.PasswordResetToken = token
	u.PasswordResetExpiresAt = &expiresAt
	u.UpdatedAt = time.Now()
}

// VerifyEmail marks the email as verified
func (u *OwnerUser) VerifyEmail() {
	u.EmailVerified = true
	u.EmailVerificationToken = ""
	u.EmailVerificationExpiresAt = nil
	u.UpdatedAt = time.Now()
}

// ResetPassword updates the password hash and clears reset token
func (u *OwnerUser) ResetPassword(newPasswordHash string) {
	u.PasswordHash = newPasswordHash
	u.PasswordResetToken = ""
	u.PasswordResetExpiresAt = nil
	u.UpdatedAt = time.Now()
}

// ApproveBusiness approves the business verification
func (u *OwnerUser) ApproveBusiness() {
	u.BusinessVerified = true
	u.BusinessVerificationStatus = BusinessVerificationApproved
	u.UpdatedAt = time.Now()
}

// RejectBusiness rejects the business verification
func (u *OwnerUser) RejectBusiness() {
	u.BusinessVerified = false
	u.BusinessVerificationStatus = BusinessVerificationRejected
	u.UpdatedAt = time.Now()
}

// IsEmailVerificationExpired checks if the email verification token is expired
func (u *OwnerUser) IsEmailVerificationExpired() bool {
	if u.EmailVerificationExpiresAt == nil {
		return true
	}
	return time.Now().After(*u.EmailVerificationExpiresAt)
}

// IsPasswordResetExpired checks if the password reset token is expired
func (u *OwnerUser) IsPasswordResetExpired() bool {
	if u.PasswordResetExpiresAt == nil {
		return true
	}
	return time.Now().After(*u.PasswordResetExpiresAt)
}

// CanAccessShopFeatures checks if the owner can access shop management features
func (u *OwnerUser) CanAccessShopFeatures() bool {
	return u.EmailVerified && u.BusinessVerified
}
