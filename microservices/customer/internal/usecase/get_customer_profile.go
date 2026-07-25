package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type GetCustomerProfileInput struct {
	CustomerID uuid.UUID
}

type GetCustomerProfileOutput struct {
	Customer *domain.Customer
}

type GetCustomerProfileUsecase interface {
	Execute(ctx context.Context, input GetCustomerProfileInput) (GetCustomerProfileOutput, error)
}

type getCustomerProfileUsecase struct {
	customerRepo repository.CustomerRepository
}

func NewGetCustomerProfileUsecase(customerRepo repository.CustomerRepository) GetCustomerProfileUsecase {
	return &getCustomerProfileUsecase{customerRepo: customerRepo}
}

func (u *getCustomerProfileUsecase) Execute(ctx context.Context, input GetCustomerProfileInput) (GetCustomerProfileOutput, error) {
	customer, err := u.customerRepo.GetByID(ctx, input.CustomerID)
	if err != nil {
		return GetCustomerProfileOutput{}, err
	}

	return GetCustomerProfileOutput{Customer: customer}, nil
}
