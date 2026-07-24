package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/repository"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) repository.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, order_id, amount, payment_method, status, transaction_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		payment.ID, payment.OrderID, payment.Amount, payment.PaymentMethod,
		payment.Status, payment.TransactionID, payment.CreatedAt, payment.UpdatedAt,
	)
	return err
}

func (r *paymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, amount, payment_method, status, transaction_id, created_at, updated_at
		FROM payments WHERE id = $1
	`
	payment := &domain.Payment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID, &payment.OrderID, &payment.Amount, &payment.PaymentMethod,
		&payment.Status, &payment.TransactionID, &payment.CreatedAt, &payment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *paymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, amount, payment_method, status, transaction_id, created_at, updated_at
		FROM payments WHERE order_id = $1
	`
	payment := &domain.Payment{}
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&payment.ID, &payment.OrderID, &payment.Amount, &payment.PaymentMethod,
		&payment.Status, &payment.TransactionID, &payment.CreatedAt, &payment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus, transactionID string) error {
	query := `
		UPDATE payments
		SET status = $1, transaction_id = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, status, transactionID, time.Now(), id)
	return err
}
