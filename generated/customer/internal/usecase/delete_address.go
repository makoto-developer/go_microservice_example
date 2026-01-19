package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type DeleteAddressInput struct {
	CustomerID uuid.UUID
	AddressID  uuid.UUID
}

type DeleteAddressOutput struct {
	Success bool
}

type DeleteAddressUsecase interface {
	Execute(ctx context.Context, input DeleteAddressInput) (DeleteAddressOutput, error)
}

type deleteAddressUsecase struct {
	addressRepo repository.AddressRepository
}

func NewDeleteAddressUsecase(addressRepo repository.AddressRepository) DeleteAddressUsecase {
	return &deleteAddressUsecase{addressRepo: addressRepo}
}

func (u *deleteAddressUsecase) Execute(ctx context.Context, input DeleteAddressInput) (DeleteAddressOutput, error) {
	address, err := u.addressRepo.GetByID(ctx, input.AddressID)
	if err != nil {
		return DeleteAddressOutput{}, err
	}

	if address.CustomerID != input.CustomerID {
		return DeleteAddressOutput{}, domain.ErrAddressNotFound
	}

	if err := u.addressRepo.Delete(ctx, input.AddressID); err != nil {
		return DeleteAddressOutput{}, err
	}

	return DeleteAddressOutput{Success: true}, nil
}
