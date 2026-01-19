package domain

import (
	"github.com/google/uuid"
	"time"
)

// Notification represents Notification
type Notification struct {
	Id uuid.UUID `db:"id" json:"id"`
	UserId uuid.UUID `db:"user_id" json:"user_id"`
	NotificationType NotificationType `db:"notification_type" json:"notification_type"`
	Channel NotificationChannel `db:"channel" json:"channel"`
	Title string `db:"title" json:"title"`
	Content text `db:"content" json:"content"`
	Status NotificationStatus `db:"status" json:"status"`
	SentAt *time.Time `db:"sent_at" json:"sent_at,omitempty"`
	ErrorMessage *text `db:"error_message" json:"error_message,omitempty"`
	RetryCount int `db:"retry_count" json:"retry_count"`
	Metadata *map[string]interface{} `db:"metadata" json:"metadata,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewNotification creates a new Notification instance
func NewNotification() *Notification {
	return &Notification{}
}
