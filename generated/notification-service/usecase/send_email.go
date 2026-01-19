package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// SendEmailInput represents input for SendEmail
type SendEmailInput struct {
	UserId uuid.UUID
	Email string
	TemplateKey string
	Variables map<string,
	NotificationType NotificationType
}

// SendEmailUsecase defines the interface for SendEmail
type SendEmailUsecase interface {
	Execute(ctx context.Context, input SendEmailInput) error
}

type send_emailUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewSendEmailUsecase creates a new instance
func NewSendEmailUsecase() SendEmailUsecase {
	return &send_emailUsecaseImpl{}
}

// Execute executes SendEmail
func (u *send_emailUsecaseImpl) Execute(ctx context.Context, input SendEmailInput) error {
	// TODO: Implement business logic

	return nil
}
