package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UploadChatImageInput represents input for UploadChatImage
type UploadChatImageInput struct {
	RoomId uuid.UUID
	UserId uuid.UUID
	ImageData bytes
}

// UploadChatImageUsecase defines the interface for UploadChatImage
type UploadChatImageUsecase interface {
	Execute(ctx context.Context, input UploadChatImageInput) error
}

type upload_chat_imageUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUploadChatImageUsecase creates a new instance
func NewUploadChatImageUsecase() UploadChatImageUsecase {
	return &upload_chat_imageUsecaseImpl{}
}

// Execute executes UploadChatImage
func (u *upload_chat_imageUsecaseImpl) Execute(ctx context.Context, input UploadChatImageInput) error {
	// TODO: Implement business logic

	return nil
}
