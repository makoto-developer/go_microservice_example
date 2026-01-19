package domain

import (
	"context"

	"github.com/google/uuid"
)

// NotificationRepository defines repository interface for Notification
type NotificationRepository interface {
	// Create creates a new Notification
	Create(ctx context.Context, notification *Notification) error

	// FindByID finds Notification by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Notification, error)

	// Update updates Notification
	Update(ctx context.Context, notification *Notification) error

	// Delete deletes Notification by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Notification
	List(ctx context.Context, limit, offset int) ([]*Notification, error)
}
