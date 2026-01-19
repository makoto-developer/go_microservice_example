package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SendBulkEmailInput represents input for SendBulkEmail
type SendBulkEmailInput struct {
	UserIds []uuid.UUID
	TemplateKey string
	Variables map<string,
}

// SendBulkEmailUsecase defines the interface for SendBulkEmail
type SendBulkEmailUsecase interface {
	Execute(ctx context.Context, input SendBulkEmailInput) error
}

type send_bulk_emailUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewSendBulkEmailUsecase creates a new instance
func NewSendBulkEmailUsecase() SendBulkEmailUsecase {
	return &send_bulk_emailUsecaseImpl{}
}

// Execute executes SendBulkEmail
func (u *send_bulk_emailUsecaseImpl) Execute(ctx context.Context, input SendBulkEmailInput) error {
	// TODO: Implement business logic

	return nil
}
