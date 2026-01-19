package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdateSystemSettingInput represents input for UpdateSystemSetting
type UpdateSystemSettingInput struct {
	AdminId uuid.UUID
	SettingKey string
	SettingValue string
}

// UpdateSystemSettingUsecase defines the interface for UpdateSystemSetting
type UpdateSystemSettingUsecase interface {
	Execute(ctx context.Context, input UpdateSystemSettingInput) error
}

type update_system_settingUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdateSystemSettingUsecase creates a new instance
func NewUpdateSystemSettingUsecase() UpdateSystemSettingUsecase {
	return &update_system_settingUsecaseImpl{}
}

// Execute executes UpdateSystemSetting
func (u *update_system_settingUsecaseImpl) Execute(ctx context.Context, input UpdateSystemSettingInput) error {
	// TODO: Implement business logic

	return nil
}
