package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/admin/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/admin/internal/usecase"
)

type mockAdminRepository struct {
	createAdminUserFunc func(ctx context.Context, admin *domain.AdminUser) error
}

func (m *mockAdminRepository) CreateAdminUser(ctx context.Context, admin *domain.AdminUser) error {
	if m.createAdminUserFunc != nil {
		return m.createAdminUserFunc(ctx, admin)
	}
	return nil
}

func (m *mockAdminRepository) GetAdminUserByID(ctx context.Context, id uuid.UUID) (*domain.AdminUser, error) {
	return nil, nil
}

func (m *mockAdminRepository) GetAdminUserByEmail(ctx context.Context, email string) (*domain.AdminUser, error) {
	return nil, nil
}

func (m *mockAdminRepository) GetAllAdminUsers(ctx context.Context) ([]*domain.AdminUser, error) {
	return nil, nil
}

func (m *mockAdminRepository) UpdateAdminUser(ctx context.Context, admin *domain.AdminUser) error {
	return nil
}

func (m *mockAdminRepository) CreateAuditLog(ctx context.Context, log *domain.AuditLog) error {
	return nil
}

func (m *mockAdminRepository) GetAuditLogsByAdminID(ctx context.Context, adminID uuid.UUID) ([]*domain.AuditLog, error) {
	return nil, nil
}

func (m *mockAdminRepository) GetAuditLogsByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*domain.AuditLog, error) {
	return nil, nil
}

func TestCreateAdminUserUsecase_Success(t *testing.T) {
	repo := &mockAdminRepository{
		createAdminUserFunc: func(ctx context.Context, admin *domain.AdminUser) error {
			if admin.Email != "admin@example.com" {
				t.Errorf("expected email admin@example.com, got %v", admin.Email)
			}
			if admin.Role != domain.AdminRoleAdmin {
				t.Errorf("expected role admin, got %v", admin.Role)
			}
			return nil
		},
	}

	uc := usecase.NewCreateAdminUserUsecase(repo)

	input := usecase.CreateAdminUserInput{
		Email: "admin@example.com",
		Name:  "Admin User",
		Role:  domain.AdminRoleAdmin,
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Email != "admin@example.com" {
		t.Errorf("expected email admin@example.com, got %v", output.Email)
	}

	if output.Role != domain.AdminRoleAdmin {
		t.Errorf("expected role admin, got %v", output.Role)
	}
}

func TestCreateAdminUserUsecase_InvalidRole(t *testing.T) {
	repo := &mockAdminRepository{}
	uc := usecase.NewCreateAdminUserUsecase(repo)

	input := usecase.CreateAdminUserInput{
		Email: "admin@example.com",
		Name:  "Admin User",
		Role:  "invalid_role",
	}

	_, err := uc.Execute(context.Background(), input)

	if err != domain.ErrInvalidRole {
		t.Errorf("expected ErrInvalidRole, got %v", err)
	}
}
