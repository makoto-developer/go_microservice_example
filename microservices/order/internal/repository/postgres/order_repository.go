package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/order/internal/domain"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `INSERT INTO orders (id, customer_id, order_number, status, total_amount, shipping_fee, address_id, payment_method_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query, order.ID, order.CustomerID, order.OrderNumber, order.Status,
		order.TotalAmount, order.ShippingFee, order.AddressID, order.PaymentMethodID, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	query := `SELECT id, customer_id, order_number, status, total_amount, shipping_fee, address_id, payment_method_id, created_at, updated_at
		FROM orders WHERE id = $1`

	var order domain.Order
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.CustomerID, &order.OrderNumber, &order.Status, &order.TotalAmount,
		&order.ShippingFee, &order.AddressID, &order.PaymentMethodID, &order.CreatedAt, &order.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	return &order, nil
}

func (r *orderRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.Order, error) {
	query := `SELECT id, customer_id, order_number, status, total_amount, shipping_fee, address_id, payment_method_id, created_at, updated_at
		FROM orders WHERE customer_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.CustomerID, &order.OrderNumber, &order.Status, &order.TotalAmount,
			&order.ShippingFee, &order.AddressID, &order.PaymentMethodID, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.OrderStatus) error {
	query := `UPDATE orders SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}
