package domain

import (
	"context"

	"github.com/google/uuid"
)

// NotificationPreferenceRepository defines repository interface for NotificationPreference
type NotificationPreferenceRepository interface {
	// Create creates a new NotificationPreference
	Create(ctx context.Context, notification_preference *NotificationPreference) error

	// FindByID finds NotificationPreference by ID
	FindByID(ctx context.Context, id uuid.UUID) (*NotificationPreference, error)

	// Update updates NotificationPreference
	Update(ctx context.Context, notification_preference *NotificationPreference) error

	// Delete deletes NotificationPreference by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all NotificationPreference
	List(ctx context.Context, limit, offset int) ([]*NotificationPreference, error)

	// FindByUserId finds NotificationPreference by user_id
	FindByUserId(ctx context.Context, user_id uuid.UUID) (*NotificationPreference, error)
}
