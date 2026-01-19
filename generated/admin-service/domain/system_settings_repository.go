package domain

import (
	"context"

	"github.com/google/uuid"
)

// SystemSettingsRepository defines repository interface for SystemSettings
type SystemSettingsRepository interface {
	// Create creates a new SystemSettings
	Create(ctx context.Context, system_settings *SystemSettings) error

	// FindByID finds SystemSettings by ID
	FindByID(ctx context.Context, id uuid.UUID) (*SystemSettings, error)

	// Update updates SystemSettings
	Update(ctx context.Context, system_settings *SystemSettings) error

	// Delete deletes SystemSettings by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all SystemSettings
	List(ctx context.Context, limit, offset int) ([]*SystemSettings, error)

	// FindBySettingKey finds SystemSettings by setting_key
	FindBySettingKey(ctx context.Context, setting_key string) (*SystemSettings, error)
}
