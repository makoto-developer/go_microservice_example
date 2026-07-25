package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type AddPaymentMethodInput struct {
	CustomerID            uuid.UUID
	StripePaymentMethodID string
	CardLast4             string
	CardBrand             string
	CardExpMonth          int
	CardExpYear           int
	CardholderName        string
	IsDefault             bool
}

type AddPaymentMethodOutput struct {
	PaymentMethodID uuid.UUID
}

type AddPaymentMethodUsecase interface {
	Execute(ctx context.Context, input AddPaymentMethodInput) (AddPaymentMethodOutput, error)
}

type addPaymentMethodUsecase struct {
	paymentMethodRepo repository.PaymentMethodRepository
}

func NewAddPaymentMethodUsecase(paymentMethodRepo repository.PaymentMethodRepository) AddPaymentMethodUsecase {
	return &addPaymentMethodUsecase{paymentMethodRepo: paymentMethodRepo}
}

func (u *addPaymentMethodUsecase) Execute(ctx context.Context, input AddPaymentMethodInput) (AddPaymentMethodOutput, error) {
	paymentMethod := domain.NewPaymentMethod(
		input.CustomerID,
		input.StripePaymentMethodID,
		input.CardLast4,
		input.CardBrand,
		input.CardExpMonth,
		input.CardExpYear,
		input.CardholderName,
		input.IsDefault,
	)

	if err := u.paymentMethodRepo.Create(ctx, paymentMethod); err != nil {
		return AddPaymentMethodOutput{}, err
	}

	if input.IsDefault {
		if err := u.paymentMethodRepo.SetDefault(ctx, input.CustomerID, paymentMethod.ID); err != nil {
			return AddPaymentMethodOutput{}, err
		}
	}

	return AddPaymentMethodOutput{PaymentMethodID: paymentMethod.ID}, nil
}
