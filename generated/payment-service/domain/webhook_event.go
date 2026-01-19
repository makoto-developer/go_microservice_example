package domain

import (
	"github.com/google/uuid"
	"time"
)

// WebhookEvent represents WebhookEvent
type WebhookEvent struct {
	Id uuid.UUID `db:"id" json:"id"`
	StripeEventId string `db:"stripe_event_id" json:"stripe_event_id"`
	EventType string `db:"event_type" json:"event_type"`
	Payload map[string]interface{} `db:"payload" json:"payload"`
	Processed bool `db:"processed" json:"processed"`
	ProcessedAt *time.Time `db:"processed_at" json:"processed_at,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewWebhookEvent creates a new WebhookEvent instance
func NewWebhookEvent() *WebhookEvent {
	return &WebhookEvent{}
}
