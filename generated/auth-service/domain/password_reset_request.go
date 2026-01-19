package domain

import (
	"time"
)

// PasswordResetRequest represents PasswordResetRequest value object
type PasswordResetRequest struct {
	Email string `db:"email" json:"email"`
	ResetLink string `db:"reset_link" json:"reset_link"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
}

// NewPasswordResetRequest creates a new PasswordResetRequest instance
func NewPasswordResetRequest() *PasswordResetRequest {
	return &PasswordResetRequest{}
}
