package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type ListPaymentMethodsInput struct {
	CustomerID uuid.UUID
}

type ListPaymentMethodsOutput struct {
	PaymentMethods []*domain.PaymentMethod
}

type ListPaymentMethodsUsecase interface {
	Execute(ctx context.Context, input ListPaymentMethodsInput) (ListPaymentMethodsOutput, error)
}

type listPaymentMethodsUsecase struct {
	paymentMethodRepo repository.PaymentMethodRepository
}

func NewListPaymentMethodsUsecase(paymentMethodRepo repository.PaymentMethodRepository) ListPaymentMethodsUsecase {
	return &listPaymentMethodsUsecase{paymentMethodRepo: paymentMethodRepo}
}

func (u *listPaymentMethodsUsecase) Execute(ctx context.Context, input ListPaymentMethodsInput) (ListPaymentMethodsOutput, error) {
	paymentMethods, err := u.paymentMethodRepo.List(ctx, input.CustomerID)
	if err != nil {
		return ListPaymentMethodsOutput{}, err
	}

	return ListPaymentMethodsOutput{PaymentMethods: paymentMethods}, nil
}
