package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UploadProfileImageInput represents input for UploadProfileImage
type UploadProfileImageInput struct {
	CustomerId uuid.UUID
	ImageData bytes
}

// UploadProfileImageUsecase defines the interface for UploadProfileImage
type UploadProfileImageUsecase interface {
	Execute(ctx context.Context, input UploadProfileImageInput) error
}

type upload_profile_imageUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUploadProfileImageUsecase creates a new instance
func NewUploadProfileImageUsecase() UploadProfileImageUsecase {
	return &upload_profile_imageUsecaseImpl{}
}

// Execute executes UploadProfileImage
func (u *upload_profile_imageUsecaseImpl) Execute(ctx context.Context, input UploadProfileImageInput) error {
	// TODO: Implement business logic

	return nil
}
