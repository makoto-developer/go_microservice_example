package domain

import (
	"time"

	"github.com/google/uuid"
)

// CustomerRefreshToken represents a refresh token for customer authentication
type CustomerRefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// NewCustomerRefreshToken creates a new customer refresh token
func NewCustomerRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) *CustomerRefreshToken {
	return &CustomerRefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}

// IsExpired checks if the refresh token is expired
func (t *CustomerRefreshToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
