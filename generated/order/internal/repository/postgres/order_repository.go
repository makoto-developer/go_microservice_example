package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/repository"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order, items []*domain.OrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO orders (id, customer_id, order_number, status, total_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.ExecContext(ctx, query,
		order.ID, order.CustomerID, order.OrderNumber, order.Status,
		order.TotalAmount, order.CreatedAt, order.UpdatedAt,
	)
	if err != nil {
		return err
	}

	for _, item := range items {
		itemQuery := `
			INSERT INTO order_items (id, order_id, product_id, shop_id, quantity, price, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = tx.ExecContext(ctx, itemQuery,
			item.ID, item.OrderID, item.ProductID, item.ShopID,
			item.Quantity, item.Price, item.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	query := `
		SELECT id, customer_id, order_number, status, total_amount, created_at, updated_at
		FROM orders WHERE id = $1
	`
	order := &domain.Order{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.CustomerID, &order.OrderNumber, &order.Status,
		&order.TotalAmount, &order.CreatedAt, &order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) GetByOrderNumber(ctx context.Context, orderNumber string) (*domain.Order, error) {
	query := `
		SELECT id, customer_id, order_number, status, total_amount, created_at, updated_at
		FROM orders WHERE order_number = $1
	`
	order := &domain.Order{}
	err := r.db.QueryRowContext(ctx, query, orderNumber).Scan(
		&order.ID, &order.CustomerID, &order.OrderNumber, &order.Status,
		&order.TotalAmount, &order.CreatedAt, &order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Order, error) {
	query := `
		SELECT id, customer_id, order_number, status, total_amount, created_at, updated_at
		FROM orders WHERE customer_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		order := &domain.Order{}
		err := rows.Scan(
			&order.ID, &order.CustomerID, &order.OrderNumber, &order.Status,
			&order.TotalAmount, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (r *orderRepository) GetItems(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, shop_id, quantity, price, created_at
		FROM order_items WHERE order_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.OrderItem
	for rows.Next() {
		item := &domain.OrderItem{}
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.ShopID,
			&item.Quantity, &item.Price, &item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *orderRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	return r.UpdateStatus(ctx, id, domain.OrderStatusCancelled)
}
