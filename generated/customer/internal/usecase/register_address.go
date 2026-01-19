package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type RegisterAddressInput struct {
	CustomerID     uuid.UUID
	AddressName    string
	PostalCode     string
	Prefecture     string
	City           string
	AddressLine1   string
	AddressLine2   *string
	RecipientName  string
	RecipientPhone string
	IsDefault      bool
}

type RegisterAddressOutput struct {
	AddressID uuid.UUID
}

type RegisterAddressUsecase interface {
	Execute(ctx context.Context, input RegisterAddressInput) (RegisterAddressOutput, error)
}

type registerAddressUsecase struct {
	customerRepo repository.CustomerRepository
	addressRepo  repository.AddressRepository
}

func NewRegisterAddressUsecase(customerRepo repository.CustomerRepository, addressRepo repository.AddressRepository) RegisterAddressUsecase {
	return &registerAddressUsecase{
		customerRepo: customerRepo,
		addressRepo:  addressRepo,
	}
}

func (u *registerAddressUsecase) Execute(ctx context.Context, input RegisterAddressInput) (RegisterAddressOutput, error) {
	if _, err := u.customerRepo.GetByID(ctx, input.CustomerID); err != nil {
		return RegisterAddressOutput{}, err
	}

	address := domain.NewAddress(
		input.CustomerID, input.AddressName, input.PostalCode, input.Prefecture,
		input.City, input.AddressLine1, input.AddressLine2, input.RecipientName,
		input.RecipientPhone, input.IsDefault,
	)

	if err := u.addressRepo.Create(ctx, address); err != nil {
		return RegisterAddressOutput{}, err
	}

	if input.IsDefault {
		if err := u.addressRepo.SetDefault(ctx, input.CustomerID, address.ID); err != nil {
			return RegisterAddressOutput{}, err
		}
	}

	return RegisterAddressOutput{AddressID: address.ID}, nil
}
