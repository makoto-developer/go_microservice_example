package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type PostgresOrderItemRepository struct {
	db *sql.DB
}

func NewPostgresOrderItemRepository(db *sql.DB) *PostgresOrderItemRepository {
	return &PostgresOrderItemRepository{db: db}
}

func (r *PostgresOrderItemRepository) Create(ctx context.Context, orderItem *domain.OrderItem) error {
	query := `
		INSERT INTO order_items (
			id, order_id, product_id, product_name,
			quantity, unit_price, subtotal, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		orderItem.Id,
		orderItem.OrderId,
		orderItem.ProductId,
		orderItem.ProductName,
		orderItem.Quantity,
		orderItem.UnitPrice,
		orderItem.Subtotal,
		orderItem.CreatedAt,
	)

	return err
}

func (r *PostgresOrderItemRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, product_name,
			quantity, unit_price, subtotal, created_at
		FROM order_items
		WHERE id = $1
	`

	var orderItem domain.OrderItem

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&orderItem.Id,
		&orderItem.OrderId,
		&orderItem.ProductId,
		&orderItem.ProductName,
		&orderItem.Quantity,
		&orderItem.UnitPrice,
		&orderItem.Subtotal,
		&orderItem.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order item not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	return &orderItem, nil
}

func (r *PostgresOrderItemRepository) Update(ctx context.Context, orderItem *domain.OrderItem) error {
	query := `
		UPDATE order_items
		SET order_id = $1,
			product_id = $2,
			product_name = $3,
			quantity = $4,
			unit_price = $5,
			subtotal = $6
		WHERE id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		orderItem.OrderId,
		orderItem.ProductId,
		orderItem.ProductName,
		orderItem.Quantity,
		orderItem.UnitPrice,
		orderItem.Subtotal,
		orderItem.Id,
	)

	return err
}

func (r *PostgresOrderItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM order_items WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresOrderItemRepository) List(ctx context.Context, limit, offset int) ([]*domain.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, product_name,
			quantity, unit_price, subtotal, created_at
		FROM order_items
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []*domain.OrderItem
	for rows.Next() {
		var orderItem domain.OrderItem

		err := rows.Scan(
			&orderItem.Id,
			&orderItem.OrderId,
			&orderItem.ProductId,
			&orderItem.ProductName,
			&orderItem.Quantity,
			&orderItem.UnitPrice,
			&orderItem.Subtotal,
			&orderItem.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, &orderItem)
	}

	return orderItems, rows.Err()
}

// BeginTx starts a transaction
func (r *PostgresOrderItemRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
