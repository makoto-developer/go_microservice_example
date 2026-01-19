package domain

import (
	"context"

	"github.com/google/uuid"
)

// EmailTemplateRepository defines repository interface for EmailTemplate
type EmailTemplateRepository interface {
	// Create creates a new EmailTemplate
	Create(ctx context.Context, email_template *EmailTemplate) error

	// FindByID finds EmailTemplate by ID
	FindByID(ctx context.Context, id uuid.UUID) (*EmailTemplate, error)

	// Update updates EmailTemplate
	Update(ctx context.Context, email_template *EmailTemplate) error

	// Delete deletes EmailTemplate by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all EmailTemplate
	List(ctx context.Context, limit, offset int) ([]*EmailTemplate, error)

	// FindByTemplateKey finds EmailTemplate by template_key
	FindByTemplateKey(ctx context.Context, template_key string) (*EmailTemplate, error)
}
