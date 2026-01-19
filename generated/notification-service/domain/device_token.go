package domain

import (
	"github.com/google/uuid"
	"time"
)

// DeviceToken represents DeviceToken
type DeviceToken struct {
	Id uuid.UUID `db:"id" json:"id"`
	UserId uuid.UUID `db:"user_id" json:"user_id"`
	DeviceId string `db:"device_id" json:"device_id"`
	Platform Platform `db:"platform" json:"platform"`
	Token string `db:"token" json:"-"`
	IsActive bool `db:"is_active" json:"is_active"`
	LastUsedAt *time.Time `db:"last_used_at" json:"last_used_at,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewDeviceToken creates a new DeviceToken instance
func NewDeviceToken() *DeviceToken {
	return &DeviceToken{}
}
