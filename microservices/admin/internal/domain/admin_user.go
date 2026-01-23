package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AdminRole string

const (
	AdminRoleSuperAdmin AdminRole = "super_admin"
	AdminRoleAdmin      AdminRole = "admin"
	AdminRoleModerator  AdminRole = "moderator"
)

var (
	ErrInvalidRole = errors.New("invalid admin role")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
)

type AdminUser struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	Role      AdminRole `db:"role" json:"role"`
	Active    bool      `db:"active" json:"active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type AuditLog struct {
	ID        uuid.UUID `db:"id" json:"id"`
	AdminID   uuid.UUID `db:"admin_id" json:"admin_id"`
	Action    string    `db:"action" json:"action"`
	EntityType string   `db:"entity_type" json:"entity_type"`
	EntityID  uuid.UUID `db:"entity_id" json:"entity_id"`
	Details   string    `db:"details" json:"details"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func NewAdminUser(email, name string, role AdminRole) (*AdminUser, error) {
	if role != AdminRoleSuperAdmin && role != AdminRoleAdmin && role != AdminRoleModerator {
		return nil, ErrInvalidRole
	}

	now := time.Now()
	return &AdminUser{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		Role:      role,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func NewAuditLog(adminID uuid.UUID, action, entityType string, entityID uuid.UUID, details string) *AuditLog {
	return &AuditLog{
		ID:         uuid.New(),
		AdminID:    adminID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    details,
		CreatedAt:  time.Now(),
	}
}

func (a *AdminUser) CanManageUsers() bool {
	return a.Role == AdminRoleSuperAdmin || a.Role == AdminRoleAdmin
}

func (a *AdminUser) CanModerateContent() bool {
	return a.Active && (a.Role == AdminRoleSuperAdmin || a.Role == AdminRoleAdmin || a.Role == AdminRoleModerator)
}

func (a *AdminUser) Deactivate() {
	a.Active = false
	a.UpdatedAt = time.Now()
}

func (a *AdminUser) Activate() {
	a.Active = true
	a.UpdatedAt = time.Now()
}
