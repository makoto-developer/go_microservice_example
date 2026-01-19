package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type paymentMethodRepository struct {
	db *sql.DB
}

func NewPaymentMethodRepository(db *sql.DB) repository.PaymentMethodRepository {
	return &paymentMethodRepository{db: db}
}

func (r *paymentMethodRepository) Create(ctx context.Context, paymentMethod *domain.PaymentMethod) error {
	query := `
		INSERT INTO payment_methods (id, customer_id, stripe_payment_method_id, card_last4, card_brand, card_exp_month, card_exp_year, cardholder_name, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		paymentMethod.ID,
		paymentMethod.CustomerID,
		paymentMethod.StripePaymentMethodID,
		paymentMethod.CardLast4,
		paymentMethod.CardBrand,
		paymentMethod.CardExpMonth,
		paymentMethod.CardExpYear,
		paymentMethod.CardholderName,
		paymentMethod.IsDefault,
		paymentMethod.CreatedAt,
		paymentMethod.UpdatedAt,
	)
	return err
}

func (r *paymentMethodRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM payment_methods WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *paymentMethodRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PaymentMethod, error) {
	query := `
		SELECT id, customer_id, stripe_payment_method_id, card_last4, card_brand, card_exp_month, card_exp_year, cardholder_name, is_default, created_at, updated_at
		FROM payment_methods WHERE id = $1
	`
	pm := &domain.PaymentMethod{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&pm.ID,
		&pm.CustomerID,
		&pm.StripePaymentMethodID,
		&pm.CardLast4,
		&pm.CardBrand,
		&pm.CardExpMonth,
		&pm.CardExpYear,
		&pm.CardholderName,
		&pm.IsDefault,
		&pm.CreatedAt,
		&pm.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrPaymentMethodNotFound
	}
	return pm, err
}

func (r *paymentMethodRepository) List(ctx context.Context, customerID uuid.UUID) ([]*domain.PaymentMethod, error) {
	query := `
		SELECT id, customer_id, stripe_payment_method_id, card_last4, card_brand, card_exp_month, card_exp_year, cardholder_name, is_default, created_at, updated_at
		FROM payment_methods WHERE customer_id = $1
		ORDER BY is_default DESC, created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paymentMethods []*domain.PaymentMethod
	for rows.Next() {
		pm := &domain.PaymentMethod{}
		if err := rows.Scan(
			&pm.ID,
			&pm.CustomerID,
			&pm.StripePaymentMethodID,
			&pm.CardLast4,
			&pm.CardBrand,
			&pm.CardExpMonth,
			&pm.CardExpYear,
			&pm.CardholderName,
			&pm.IsDefault,
			&pm.CreatedAt,
			&pm.UpdatedAt,
		); err != nil {
			return nil, err
		}
		paymentMethods = append(paymentMethods, pm)
	}
	return paymentMethods, rows.Err()
}

func (r *paymentMethodRepository) SetDefault(ctx context.Context, customerID, paymentMethodID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `UPDATE payment_methods SET is_default = false WHERE customer_id = $1`, customerID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE payment_methods SET is_default = true WHERE id = $1`, paymentMethodID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
