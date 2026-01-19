package domain

import (
	"github.com/google/uuid"
	"time"
)

// UserPresence represents UserPresence
type UserPresence struct {
	Id uuid.UUID `db:"id" json:"id"`
	UserId uuid.UUID `db:"user_id" json:"user_id"`
	Status PresenceStatus `db:"status" json:"status"`
	LastActiveAt time.Time `db:"last_active_at" json:"last_active_at"`
	DeviceInfo *string `db:"device_info" json:"device_info,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewUserPresence creates a new UserPresence instance
func NewUserPresence() *UserPresence {
	return &UserPresence{}
}
