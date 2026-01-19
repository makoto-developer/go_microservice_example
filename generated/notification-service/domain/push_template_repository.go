package domain

import (
	"context"

	"github.com/google/uuid"
)

// PushTemplateRepository defines repository interface for PushTemplate
type PushTemplateRepository interface {
	// Create creates a new PushTemplate
	Create(ctx context.Context, push_template *PushTemplate) error

	// FindByID finds PushTemplate by ID
	FindByID(ctx context.Context, id uuid.UUID) (*PushTemplate, error)

	// Update updates PushTemplate
	Update(ctx context.Context, push_template *PushTemplate) error

	// Delete deletes PushTemplate by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all PushTemplate
	List(ctx context.Context, limit, offset int) ([]*PushTemplate, error)

	// FindByTemplateKey finds PushTemplate by template_key
	FindByTemplateKey(ctx context.Context, template_key string) (*PushTemplate, error)
}
