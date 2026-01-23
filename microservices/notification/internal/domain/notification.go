package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string
type NotificationStatus string

const (
	NotificationTypeEmail NotificationType = "email"
	NotificationTypeSMS   NotificationType = "sms"
	NotificationTypePush  NotificationType = "push"
)

const (
	NotificationStatusPending NotificationStatus = "pending"
	NotificationStatusSent    NotificationStatus = "sent"
	NotificationStatusFailed  NotificationStatus = "failed"
)

type Notification struct {
	ID        uuid.UUID          `db:"id" json:"id"`
	UserID    uuid.UUID          `db:"user_id" json:"user_id"`
	Type      NotificationType   `db:"type" json:"type"`
	Status    NotificationStatus `db:"status" json:"status"`
	Subject   string             `db:"subject" json:"subject"`
	Message   string             `db:"message" json:"message"`
	Recipient string             `db:"recipient" json:"recipient"` // email or phone
	SentAt    *time.Time         `db:"sent_at" json:"sent_at,omitempty"`
	CreatedAt time.Time          `db:"created_at" json:"created_at"`
}

func NewNotification(userID uuid.UUID, notifType NotificationType, subject, message, recipient string) *Notification {
	return &Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      notifType,
		Status:    NotificationStatusPending,
		Subject:   subject,
		Message:   message,
		Recipient: recipient,
		CreatedAt: time.Now(),
	}
}

func (n *Notification) MarkSent() {
	now := time.Now()
	n.Status = NotificationStatusSent
	n.SentAt = &now
}

func (n *Notification) MarkFailed() {
	n.Status = NotificationStatusFailed
}
