package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type UpdatePaymentMethodInput struct {
	CustomerID      uuid.UUID
	PaymentMethodID uuid.UUID
	IsDefault       bool
}

type UpdatePaymentMethodOutput struct {
	Success bool
}

type UpdatePaymentMethodUsecase interface {
	Execute(ctx context.Context, input UpdatePaymentMethodInput) (UpdatePaymentMethodOutput, error)
}

type updatePaymentMethodUsecase struct {
	paymentMethodRepo repository.PaymentMethodRepository
}

func NewUpdatePaymentMethodUsecase(paymentMethodRepo repository.PaymentMethodRepository) UpdatePaymentMethodUsecase {
	return &updatePaymentMethodUsecase{paymentMethodRepo: paymentMethodRepo}
}

func (u *updatePaymentMethodUsecase) Execute(ctx context.Context, input UpdatePaymentMethodInput) (UpdatePaymentMethodOutput, error) {
	paymentMethod, err := u.paymentMethodRepo.GetByID(ctx, input.PaymentMethodID)
	if err != nil {
		return UpdatePaymentMethodOutput{}, err
	}

	if paymentMethod.CustomerID != input.CustomerID {
		return UpdatePaymentMethodOutput{}, domain.ErrPaymentMethodNotFound
	}

	if input.IsDefault {
		if err := u.paymentMethodRepo.SetDefault(ctx, input.CustomerID, paymentMethod.ID); err != nil {
			return UpdatePaymentMethodOutput{}, err
		}
	}

	return UpdatePaymentMethodOutput{Success: true}, nil
}
