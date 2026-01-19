package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UploadChatFileInput represents input for UploadChatFile
type UploadChatFileInput struct {
	RoomId uuid.UUID
	UserId uuid.UUID
	FileData bytes
	FileName string
}

// UploadChatFileUsecase defines the interface for UploadChatFile
type UploadChatFileUsecase interface {
	Execute(ctx context.Context, input UploadChatFileInput) error
}

type upload_chat_fileUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUploadChatFileUsecase creates a new instance
func NewUploadChatFileUsecase() UploadChatFileUsecase {
	return &upload_chat_fileUsecaseImpl{}
}

// Execute executes UploadChatFile
func (u *upload_chat_fileUsecaseImpl) Execute(ctx context.Context, input UploadChatFileInput) error {
	// TODO: Implement business logic

	return nil
}
