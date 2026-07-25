package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/admin/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/admin/internal/repository"
)

type CreateAdminUserInput struct {
	Email string
	Name  string
	Role  domain.AdminRole
}

type CreateAdminUserOutput struct {
	AdminID uuid.UUID
	Email   string
	Role    domain.AdminRole
}

type CreateAdminUserUsecase interface {
	Execute(ctx context.Context, input CreateAdminUserInput) (CreateAdminUserOutput, error)
}

type createAdminUserUsecaseImpl struct {
	adminRepo repository.AdminRepository
}

func NewCreateAdminUserUsecase(adminRepo repository.AdminRepository) CreateAdminUserUsecase {
	return &createAdminUserUsecaseImpl{
		adminRepo: adminRepo,
	}
}

func (u *createAdminUserUsecaseImpl) Execute(ctx context.Context, input CreateAdminUserInput) (CreateAdminUserOutput, error) {
	admin, err := domain.NewAdminUser(input.Email, input.Name, input.Role)
	if err != nil {
		return CreateAdminUserOutput{}, err
	}

	err = u.adminRepo.CreateAdminUser(ctx, admin)
	if err != nil {
		return CreateAdminUserOutput{}, err
	}

	return CreateAdminUserOutput{
		AdminID: admin.ID,
		Email:   admin.Email,
		Role:    admin.Role,
	}, nil
}
