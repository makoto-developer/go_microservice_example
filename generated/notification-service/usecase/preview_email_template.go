package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// PreviewEmailTemplateInput represents input for PreviewEmailTemplate
type PreviewEmailTemplateInput struct {
	TemplateKey string
	Variables map<string,
}

// PreviewEmailTemplateUsecase defines the interface for PreviewEmailTemplate
type PreviewEmailTemplateUsecase interface {
	Execute(ctx context.Context, input PreviewEmailTemplateInput) error
}

type preview_email_templateUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewPreviewEmailTemplateUsecase creates a new instance
func NewPreviewEmailTemplateUsecase() PreviewEmailTemplateUsecase {
	return &preview_email_templateUsecaseImpl{}
}

// Execute executes PreviewEmailTemplate
func (u *preview_email_templateUsecaseImpl) Execute(ctx context.Context, input PreviewEmailTemplateInput) error {
	// TODO: Implement business logic

	return nil
}
