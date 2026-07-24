package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/repository"
)

type CreateCODPaymentInput struct {
	OrderID uuid.UUID
	Amount  int
}

type CreateCODPaymentOutput struct {
	PaymentID uuid.UUID
	Status    domain.PaymentStatus
}

type ConfirmCODPaymentInput struct {
	PaymentID uuid.UUID
	OrderID   uuid.UUID
}

// CODPaymentUsecase は代金引換(現金)の決済フローを扱う。
// Create で「支払い待ち」の決済を作り、配達完了時の Confirm で入金確定にする。
type CODPaymentUsecase interface {
	Create(ctx context.Context, input CreateCODPaymentInput) (CreateCODPaymentOutput, error)
	Confirm(ctx context.Context, input ConfirmCODPaymentInput) error
}

type codPaymentUsecaseImpl struct {
	paymentRepo repository.PaymentRepository
}

func NewCODPaymentUsecase(paymentRepo repository.PaymentRepository) CODPaymentUsecase {
	return &codPaymentUsecaseImpl{paymentRepo: paymentRepo}
}

func (u *codPaymentUsecaseImpl) Create(ctx context.Context, input CreateCODPaymentInput) (CreateCODPaymentOutput, error) {
	if input.Amount <= 0 {
		return CreateCODPaymentOutput{}, domain.ErrInvalidAmount
	}

	existing, err := u.paymentRepo.GetByOrderID(ctx, input.OrderID)
	if err != nil {
		return CreateCODPaymentOutput{}, err
	}
	if existing != nil {
		return CreateCODPaymentOutput{}, domain.ErrPaymentAlreadyProcessed
	}

	// 代引きは配達完了まで pending のまま
	payment := domain.NewPayment(input.OrderID, input.Amount, domain.PaymentMethodCashOnDelivery)
	if err := u.paymentRepo.Create(ctx, payment); err != nil {
		return CreateCODPaymentOutput{}, err
	}

	return CreateCODPaymentOutput{PaymentID: payment.ID, Status: payment.Status}, nil
}

func (u *codPaymentUsecaseImpl) Confirm(ctx context.Context, input ConfirmCODPaymentInput) error {
	payment, err := u.paymentRepo.GetByID(ctx, input.PaymentID)
	if err != nil {
		return err
	}
	if payment == nil {
		return domain.ErrPaymentNotFound
	}
	if payment.OrderID != input.OrderID {
		return domain.ErrOrderMismatch
	}
	if payment.PaymentMethod != domain.PaymentMethodCashOnDelivery {
		return domain.ErrNotCODPayment
	}

	// 配達員の集金を入金として確定する(実運用では配送サービスからの通知で呼ばれる想定)
	transactionID := "COD-" + uuid.New().String()[:8]
	if err := payment.Complete(transactionID); err != nil {
		return err
	}
	return u.paymentRepo.UpdateStatus(ctx, payment.ID, payment.Status, payment.TransactionID)
}
