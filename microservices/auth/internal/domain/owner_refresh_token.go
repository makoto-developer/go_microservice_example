package domain

import (
	"time"

	"github.com/google/uuid"
)

// OwnerRefreshToken represents a refresh token for owner authentication
type OwnerRefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// NewOwnerRefreshToken creates a new owner refresh token
func NewOwnerRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) *OwnerRefreshToken {
	return &OwnerRefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
}

// IsExpired checks if the refresh token is expired
func (t *OwnerRefreshToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
