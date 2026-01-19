package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ArchiveOldMessagesInput represents input for ArchiveOldMessages
type ArchiveOldMessagesInput struct {
	DaysThreshold int
}

// ArchiveOldMessagesUsecase defines the interface for ArchiveOldMessages
type ArchiveOldMessagesUsecase interface {
	Execute(ctx context.Context, input ArchiveOldMessagesInput) error
}

type archive_old_messagesUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewArchiveOldMessagesUsecase creates a new instance
func NewArchiveOldMessagesUsecase() ArchiveOldMessagesUsecase {
	return &archive_old_messagesUsecaseImpl{}
}

// Execute executes ArchiveOldMessages
func (u *archive_old_messagesUsecaseImpl) Execute(ctx context.Context, input ArchiveOldMessagesInput) error {
	// TODO: Implement business logic

	return nil
}
