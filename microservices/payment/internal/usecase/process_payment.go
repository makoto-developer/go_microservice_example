package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/repository"
)

type ProcessPaymentInput struct {
	OrderID       uuid.UUID
	Amount        int
	PaymentMethod domain.PaymentMethod
}

type ProcessPaymentOutput struct {
	PaymentID     uuid.UUID
	Status        domain.PaymentStatus
	TransactionID string
}

type ProcessPaymentUsecase interface {
	Execute(ctx context.Context, input ProcessPaymentInput) (ProcessPaymentOutput, error)
}

type processPaymentUsecaseImpl struct {
	paymentRepo repository.PaymentRepository
}

func NewProcessPaymentUsecase(paymentRepo repository.PaymentRepository) ProcessPaymentUsecase {
	return &processPaymentUsecaseImpl{
		paymentRepo: paymentRepo,
	}
}

func (u *processPaymentUsecaseImpl) Execute(ctx context.Context, input ProcessPaymentInput) (ProcessPaymentOutput, error) {
	if input.Amount <= 0 {
		return ProcessPaymentOutput{}, domain.ErrInvalidAmount
	}

	// 冪等性: 同一注文への決済が既に存在する場合は二重課金せず既存の結果を返す。
	// (Saga リトライやクライアント再送で Execute が複数回呼ばれても安全にする)
	if existing, err := u.paymentRepo.GetByOrderID(ctx, input.OrderID); err != nil {
		return ProcessPaymentOutput{}, err
	} else if existing != nil {
		return ProcessPaymentOutput{
			PaymentID:     existing.ID,
			Status:        existing.Status,
			TransactionID: existing.TransactionID,
		}, nil
	}

	payment := domain.NewPayment(input.OrderID, input.Amount, input.PaymentMethod)

	// Simulate payment processing (in real implementation, call Stripe/etc)
	transactionID := uuid.New().String()[:8]
	err := payment.Complete(transactionID)
	if err != nil {
		return ProcessPaymentOutput{}, err
	}

	err = u.paymentRepo.Create(ctx, payment)
	if err != nil {
		return ProcessPaymentOutput{}, err
	}

	return ProcessPaymentOutput{
		PaymentID:     payment.ID,
		Status:        payment.Status,
		TransactionID: payment.TransactionID,
	}, nil
}
