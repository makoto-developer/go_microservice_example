package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type UpdateAddressInput struct {
	CustomerID     uuid.UUID
	AddressID      uuid.UUID
	AddressName    string
	PostalCode     string
	Prefecture     string
	City           string
	AddressLine1   string
	AddressLine2   *string
	RecipientName  string
	RecipientPhone string
}

type UpdateAddressOutput struct {
	Address *domain.Address
}

type UpdateAddressUsecase interface {
	Execute(ctx context.Context, input UpdateAddressInput) (UpdateAddressOutput, error)
}

type updateAddressUsecase struct {
	addressRepo repository.AddressRepository
}

func NewUpdateAddressUsecase(addressRepo repository.AddressRepository) UpdateAddressUsecase {
	return &updateAddressUsecase{addressRepo: addressRepo}
}

func (u *updateAddressUsecase) Execute(ctx context.Context, input UpdateAddressInput) (UpdateAddressOutput, error) {
	address, err := u.addressRepo.GetByID(ctx, input.AddressID)
	if err != nil {
		return UpdateAddressOutput{}, err
	}

	if address.CustomerID != input.CustomerID {
		return UpdateAddressOutput{}, domain.ErrAddressNotFound
	}

	address.Label = input.AddressName
	address.PostalCode = input.PostalCode
	address.Prefecture = input.Prefecture
	address.City = input.City
	address.AddressLine1 = input.AddressLine1
	address.AddressLine2 = input.AddressLine2
	address.RecipientName = input.RecipientName
	address.RecipientPhone = input.RecipientPhone
	address.UpdatedAt = time.Now()

	if err := u.addressRepo.Update(ctx, address); err != nil {
		return UpdateAddressOutput{}, err
	}

	return UpdateAddressOutput{Address: address}, nil
}
