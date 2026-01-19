package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (
			id, shop_id, customer_id, order_number, status,
			total_amount, shipping_address, payment_method,
			tracking_number, carrier, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		order.Id,
		order.ShopId,
		order.CustomerId,
		order.OrderNumber,
		order.Status,
		order.TotalAmount,
		order.ShippingAddress,
		order.PaymentMethod,
		order.TrackingNumber,
		order.Carrier,
		order.CreatedAt,
		order.UpdatedAt,
	)

	return err
}

func (r *PostgresOrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	query := `
		SELECT id, shop_id, customer_id, order_number, status,
			total_amount, shipping_address, payment_method,
			tracking_number, carrier, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var order domain.Order
	var trackingNumber sql.NullString
	var carrier sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.Id,
		&order.ShopId,
		&order.CustomerId,
		&order.OrderNumber,
		&order.Status,
		&order.TotalAmount,
		&order.ShippingAddress,
		&order.PaymentMethod,
		&trackingNumber,
		&carrier,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	if trackingNumber.Valid {
		order.TrackingNumber = &trackingNumber.String
	}
	if carrier.Valid {
		c := domain.Carrier(carrier.String)
		order.Carrier = &c
	}

	return &order, nil
}

func (r *PostgresOrderRepository) FindByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	query := `
		SELECT id, shop_id, customer_id, order_number, status,
			total_amount, shipping_address, payment_method,
			tracking_number, carrier, created_at, updated_at
		FROM orders
		WHERE order_number = $1
	`

	var order domain.Order
	var trackingNumber sql.NullString
	var carrier sql.NullString

	err := r.db.QueryRowContext(ctx, query, orderNumber).Scan(
		&order.Id,
		&order.ShopId,
		&order.CustomerId,
		&order.OrderNumber,
		&order.Status,
		&order.TotalAmount,
		&order.ShippingAddress,
		&order.PaymentMethod,
		&trackingNumber,
		&carrier,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found with order number: %s", orderNumber)
	}
	if err != nil {
		return nil, err
	}

	if trackingNumber.Valid {
		order.TrackingNumber = &trackingNumber.String
	}
	if carrier.Valid {
		c := domain.Carrier(carrier.String)
		order.Carrier = &c
	}

	return &order, nil
}

func (r *PostgresOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	query := `
		UPDATE orders
		SET shop_id = $1,
			customer_id = $2,
			order_number = $3,
			status = $4,
			total_amount = $5,
			shipping_address = $6,
			payment_method = $7,
			tracking_number = $8,
			carrier = $9,
			updated_at = $10
		WHERE id = $11
	`

	_, err := r.db.ExecContext(ctx, query,
		order.ShopId,
		order.CustomerId,
		order.OrderNumber,
		order.Status,
		order.TotalAmount,
		order.ShippingAddress,
		order.PaymentMethod,
		order.TrackingNumber,
		order.Carrier,
		time.Now(),
		order.Id,
	)

	return err
}

func (r *PostgresOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM orders WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresOrderRepository) List(ctx context.Context, limit, offset int) ([]*domain.Order, error) {
	query := `
		SELECT id, shop_id, customer_id, order_number, status,
			total_amount, shipping_address, payment_method,
			tracking_number, carrier, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		var trackingNumber sql.NullString
		var carrier sql.NullString

		err := rows.Scan(
			&order.Id,
			&order.ShopId,
			&order.CustomerId,
			&order.OrderNumber,
			&order.Status,
			&order.TotalAmount,
			&order.ShippingAddress,
			&order.PaymentMethod,
			&trackingNumber,
			&carrier,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if trackingNumber.Valid {
			order.TrackingNumber = &trackingNumber.String
		}
		if carrier.Valid {
			c := domain.Carrier(carrier.String)
			order.Carrier = &c
		}

		orders = append(orders, &order)
	}

	return orders, rows.Err()
}

// BeginTx starts a transaction
func (r *PostgresOrderRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
