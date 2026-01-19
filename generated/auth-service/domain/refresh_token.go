package domain

import (
	"github.com/google/uuid"
	"time"
)

// RefreshToken represents RefreshToken
type RefreshToken struct {
	Id uuid.UUID `db:"id" json:"id"`
	UserId uuid.UUID `db:"user_id" json:"user_id"`
	Token string `db:"token" json:"-"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	Revoked bool `db:"revoked" json:"revoked"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewRefreshToken creates a new RefreshToken instance
func NewRefreshToken() *RefreshToken {
	return &RefreshToken{}
}
