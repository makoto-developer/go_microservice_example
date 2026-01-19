package domain

import (
	"github.com/google/uuid"
	"time"
)

// User represents User
type User struct {
	Id uuid.UUID `db:"id" json:"id"`
	Email string `db:"email" json:"email"`
	PasswordHash string `db:"password_hash" json:"-"`
	Role Role `db:"role" json:"role"`
	EmailVerified bool `db:"email_verified" json:"email_verified"`
	EmailVerificationToken *string `db:"email_verification_token" json:"-"`
	EmailVerificationExpiresAt *time.Time `db:"email_verification_expires_at" json:"email_verification_expires_at,omitempty"`
	PasswordResetToken *string `db:"password_reset_token" json:"-"`
	PasswordResetExpiresAt *time.Time `db:"password_reset_expires_at" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewUser creates a new User instance
func NewUser() *User {
	return &User{}
}
