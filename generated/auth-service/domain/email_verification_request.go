package domain

import (
	"time"
)

// EmailVerificationRequest represents EmailVerificationRequest value object
type EmailVerificationRequest struct {
	Email string `db:"email" json:"email"`
	VerificationLink string `db:"verification_link" json:"verification_link"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}

// NewEmailVerificationRequest creates a new EmailVerificationRequest instance
func NewEmailVerificationRequest() *EmailVerificationRequest {
	return &EmailVerificationRequest{}
}
