package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/admin/internal/domain"
)

type AdminRepository interface {
	// AdminUser operations
	CreateAdminUser(ctx context.Context, admin *domain.AdminUser) error
	GetAdminUserByID(ctx context.Context, id uuid.UUID) (*domain.AdminUser, error)
	GetAdminUserByEmail(ctx context.Context, email string) (*domain.AdminUser, error)
	GetAllAdminUsers(ctx context.Context) ([]*domain.AdminUser, error)
	UpdateAdminUser(ctx context.Context, admin *domain.AdminUser) error

	// AuditLog operations
	CreateAuditLog(ctx context.Context, log *domain.AuditLog) error
	GetAuditLogsByAdminID(ctx context.Context, adminID uuid.UUID) ([]*domain.AuditLog, error)
	GetAuditLogsByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.AuditLog, error)
}
