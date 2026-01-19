package domain

import (
	"github.com/google/uuid"
	"time"
)

// SystemSettings represents SystemSettings
type SystemSettings struct {
	Id uuid.UUID `db:"id" json:"id"`
	SettingKey string `db:"setting_key" json:"setting_key"`
	SettingValue string `db:"setting_value" json:"setting_value"`
	SettingType SettingType `db:"setting_type" json:"setting_type"`
	Description *text `db:"description" json:"description,omitempty"`
	UpdatedBy uuid.UUID `db:"updated_by" json:"updated_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewSystemSettings creates a new SystemSettings instance
func NewSystemSettings() *SystemSettings {
	return &SystemSettings{}
}
