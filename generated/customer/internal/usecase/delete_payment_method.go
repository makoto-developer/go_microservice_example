package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type DeletePaymentMethodInput struct {
	CustomerID      uuid.UUID
	PaymentMethodID uuid.UUID
}

type DeletePaymentMethodOutput struct {
	Success bool
}

type DeletePaymentMethodUsecase interface {
	Execute(ctx context.Context, input DeletePaymentMethodInput) (DeletePaymentMethodOutput, error)
}

type deletePaymentMethodUsecase struct {
	paymentMethodRepo repository.PaymentMethodRepository
}

func NewDeletePaymentMethodUsecase(paymentMethodRepo repository.PaymentMethodRepository) DeletePaymentMethodUsecase {
	return &deletePaymentMethodUsecase{paymentMethodRepo: paymentMethodRepo}
}

func (u *deletePaymentMethodUsecase) Execute(ctx context.Context, input DeletePaymentMethodInput) (DeletePaymentMethodOutput, error) {
	paymentMethod, err := u.paymentMethodRepo.GetByID(ctx, input.PaymentMethodID)
	if err != nil {
		return DeletePaymentMethodOutput{}, err
	}

	if paymentMethod.CustomerID != input.CustomerID {
		return DeletePaymentMethodOutput{}, domain.ErrPaymentMethodNotFound
	}

	if err := u.paymentMethodRepo.Delete(ctx, input.PaymentMethodID); err != nil {
		return DeletePaymentMethodOutput{}, err
	}

	return DeletePaymentMethodOutput{Success: true}, nil
}
