package domain

import (
	"context"

	"github.com/google/uuid"
)

// DeviceTokenRepository defines repository interface for DeviceToken
type DeviceTokenRepository interface {
	// Create creates a new DeviceToken
	Create(ctx context.Context, device_token *DeviceToken) error

	// FindByID finds DeviceToken by ID
	FindByID(ctx context.Context, id uuid.UUID) (*DeviceToken, error)

	// Update updates DeviceToken
	Update(ctx context.Context, device_token *DeviceToken) error

	// Delete deletes DeviceToken by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all DeviceToken
	List(ctx context.Context, limit, offset int) ([]*DeviceToken, error)

	// FindByDeviceId finds DeviceToken by device_id
	FindByDeviceId(ctx context.Context, device_id string) (*DeviceToken, error)
}
