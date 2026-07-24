package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/repository"
)

type refundRepository struct {
	db *sql.DB
}

func NewRefundRepository(db *sql.DB) repository.RefundRepository {
	return &refundRepository{db: db}
}

func (r *refundRepository) Create(ctx context.Context, refund *domain.Refund) error {
	query := `
		INSERT INTO refunds (id, payment_id, order_id, amount, reason, status, transaction_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		refund.ID, refund.PaymentID, refund.OrderID, refund.Amount, refund.Reason,
		refund.Status, refund.TransactionID, refund.CreatedAt, refund.UpdatedAt,
	)
	return err
}

func (r *refundRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Refund, error) {
	query := `
		SELECT id, payment_id, order_id, amount, reason, status, transaction_id, created_at, updated_at
		FROM refunds WHERE id = $1
	`
	refund := &domain.Refund{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&refund.ID, &refund.PaymentID, &refund.OrderID, &refund.Amount, &refund.Reason,
		&refund.Status, &refund.TransactionID, &refund.CreatedAt, &refund.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return refund, nil
}
