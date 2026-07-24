package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/domain"
)

// PaymentListFilter は管理者・加盟店画面の一覧取得条件。ゼロ値のフィールドは無視される。
type PaymentListFilter struct {
	OrderID  uuid.UUID // uuid.Nil なら全件
	Statuses []domain.PaymentStatus
	Method   domain.PaymentMethod // "" なら全手段
	Page     int                  // 1 始まり(0 は 1 扱い)
	PageSize int                  // 0 は 20 扱い
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error)
	List(ctx context.Context, filter PaymentListFilter) ([]*domain.Payment, int, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus, transactionID string) error
}
