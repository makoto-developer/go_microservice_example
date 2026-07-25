package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type UpdateCustomerProfileInput struct {
	CustomerID uuid.UUID
	FirstName  string
	LastName   string
	Phone      string
	BirthDate  *time.Time
	Gender     *domain.Gender
}

type UpdateCustomerProfileOutput struct {
	Customer *domain.Customer
}

type UpdateCustomerProfileUsecase interface {
	Execute(ctx context.Context, input UpdateCustomerProfileInput) (UpdateCustomerProfileOutput, error)
}

type updateCustomerProfileUsecase struct {
	customerRepo repository.CustomerRepository
}

func NewUpdateCustomerProfileUsecase(customerRepo repository.CustomerRepository) UpdateCustomerProfileUsecase {
	return &updateCustomerProfileUsecase{customerRepo: customerRepo}
}

func (u *updateCustomerProfileUsecase) Execute(ctx context.Context, input UpdateCustomerProfileInput) (UpdateCustomerProfileOutput, error) {
	customer, err := u.customerRepo.GetByID(ctx, input.CustomerID)
	if err != nil {
		return UpdateCustomerProfileOutput{}, err
	}

	customer.FirstName = input.FirstName
	customer.LastName = input.LastName
	customer.PhoneNumber = input.Phone
	customer.BirthDate = input.BirthDate
	if input.Gender != nil {
		genderStr := string(*input.Gender)
		customer.Gender = &genderStr
	}
	customer.UpdatedAt = time.Now()

	if err := u.customerRepo.Update(ctx, customer); err != nil {
		return UpdateCustomerProfileOutput{}, err
	}

	return UpdateCustomerProfileOutput{Customer: customer}, nil
}
