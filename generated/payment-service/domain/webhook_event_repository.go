package domain

import (
	"context"

	"github.com/google/uuid"
)

// WebhookEventRepository defines repository interface for WebhookEvent
type WebhookEventRepository interface {
	// Create creates a new WebhookEvent
	Create(ctx context.Context, webhook_event *WebhookEvent) error

	// FindByID finds WebhookEvent by ID
	FindByID(ctx context.Context, id uuid.UUID) (*WebhookEvent, error)

	// Update updates WebhookEvent
	Update(ctx context.Context, webhook_event *WebhookEvent) error

	// Delete deletes WebhookEvent by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all WebhookEvent
	List(ctx context.Context, limit, offset int) ([]*WebhookEvent, error)

	// FindByStripeEventId finds WebhookEvent by stripe_event_id
	FindByStripeEventId(ctx context.Context, stripe_event_id string) (*WebhookEvent, error)
}
