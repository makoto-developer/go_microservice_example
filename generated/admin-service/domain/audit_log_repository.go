package domain

import (
	"context"

	"github.com/google/uuid"
)

// AuditLogRepository defines repository interface for AuditLog
type AuditLogRepository interface {
	// Create creates a new AuditLog
	Create(ctx context.Context, audit_log *AuditLog) error

	// FindByID finds AuditLog by ID
	FindByID(ctx context.Context, id uuid.UUID) (*AuditLog, error)

	// Update updates AuditLog
	Update(ctx context.Context, audit_log *AuditLog) error

	// Delete deletes AuditLog by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all AuditLog
	List(ctx context.Context, limit, offset int) ([]*AuditLog, error)
}
