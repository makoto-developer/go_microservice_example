package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateEmailTemplateInput represents input for CreateEmailTemplate
type CreateEmailTemplateInput struct {
	TemplateKey string
	SubjectTemplate string
	HtmlTemplate string
	TextTemplate string
	Variables []string
}

// CreateEmailTemplateUsecase defines the interface for CreateEmailTemplate
type CreateEmailTemplateUsecase interface {
	Execute(ctx context.Context, input CreateEmailTemplateInput) error
}

type create_email_templateUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewCreateEmailTemplateUsecase creates a new instance
func NewCreateEmailTemplateUsecase() CreateEmailTemplateUsecase {
	return &create_email_templateUsecaseImpl{}
}

// Execute executes CreateEmailTemplate
func (u *create_email_templateUsecaseImpl) Execute(ctx context.Context, input CreateEmailTemplateInput) error {
	// TODO: Implement business logic

	return nil
}
