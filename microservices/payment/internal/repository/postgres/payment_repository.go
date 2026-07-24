package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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

func (r *paymentRepository) List(ctx context.Context, filter repository.PaymentListFilter) ([]*domain.Payment, int, error) {
	where := "WHERE 1=1"
	args := []any{}
	if filter.OrderID != uuid.Nil {
		args = append(args, filter.OrderID)
		where += fmt.Sprintf(" AND order_id = $%d", len(args))
	}
	if len(filter.Statuses) > 0 {
		placeholders := make([]string, 0, len(filter.Statuses))
		for _, s := range filter.Statuses {
			args = append(args, s)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}
		where += " AND status IN (" + strings.Join(placeholders, ", ") + ")"
	}
	if filter.Method != "" {
		args = append(args, filter.Method)
		where += fmt.Sprintf(" AND payment_method = $%d", len(args))
	}

	total := 0
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM payments "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	args = append(args, pageSize, (page-1)*pageSize)
	query := fmt.Sprintf(`
		SELECT id, order_id, amount, payment_method, status, transaction_id, created_at, updated_at
		FROM payments %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	payments := []*domain.Payment{}
	for rows.Next() {
		payment := &domain.Payment{}
		if err := rows.Scan(
			&payment.ID, &payment.OrderID, &payment.Amount, &payment.PaymentMethod,
			&payment.Status, &payment.TransactionID, &payment.CreatedAt, &payment.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		payments = append(payments, payment)
	}
	return payments, total, rows.Err()
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
