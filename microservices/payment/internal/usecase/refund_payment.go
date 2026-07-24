package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/repository"
)

type RefundPaymentInput struct {
	PaymentID uuid.UUID // uuid.Nil の場合は OrderID から決済を引く
	OrderID   uuid.UUID
	Amount    int // 0 = 全額返金
	Reason    string
}

type RefundPaymentOutput struct {
	RefundID uuid.UUID
	Status   domain.RefundStatus
}

type RefundPaymentUsecase interface {
	Execute(ctx context.Context, input RefundPaymentInput) (RefundPaymentOutput, error)
}

type refundPaymentUsecaseImpl struct {
	paymentRepo repository.PaymentRepository
	refundRepo  repository.RefundRepository
}

func NewRefundPaymentUsecase(
	paymentRepo repository.PaymentRepository,
	refundRepo repository.RefundRepository,
) RefundPaymentUsecase {
	return &refundPaymentUsecaseImpl{paymentRepo: paymentRepo, refundRepo: refundRepo}
}

func (u *refundPaymentUsecaseImpl) Execute(ctx context.Context, input RefundPaymentInput) (RefundPaymentOutput, error) {
	payment, err := u.resolvePayment(ctx, input)
	if err != nil {
		return RefundPaymentOutput{}, err
	}

	amount := input.Amount
	if amount == 0 {
		amount = payment.Amount // 全額返金
	}
	refund, err := domain.NewRefund(payment, amount, input.Reason)
	if err != nil {
		return RefundPaymentOutput{}, err
	}

	// 決済プロバイダへの返金要求を模擬(実運用では Stripe 等を呼ぶ)
	refund.Complete("RF-" + uuid.New().String()[:8])

	if err := u.refundRepo.Create(ctx, refund); err != nil {
		return RefundPaymentOutput{}, err
	}
	if err := payment.Refund(); err != nil {
		return RefundPaymentOutput{}, err
	}
	if err := u.paymentRepo.UpdateStatus(ctx, payment.ID, payment.Status, payment.TransactionID); err != nil {
		return RefundPaymentOutput{}, err
	}

	return RefundPaymentOutput{RefundID: refund.ID, Status: refund.Status}, nil
}

func (u *refundPaymentUsecaseImpl) resolvePayment(ctx context.Context, input RefundPaymentInput) (*domain.Payment, error) {
	var payment *domain.Payment
	var err error
	switch {
	case input.PaymentID != uuid.Nil:
		payment, err = u.paymentRepo.GetByID(ctx, input.PaymentID)
	case input.OrderID != uuid.Nil:
		payment, err = u.paymentRepo.GetByOrderID(ctx, input.OrderID)
	default:
		return nil, domain.ErrPaymentNotFound
	}
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, domain.ErrPaymentNotFound
	}
	return payment, nil
}
