package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `
		INSERT INTO orders (id, shop_id, customer_id, order_number, status, total_amount,
		                   shipping_address, payment_method, tracking_number, carrier, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		order.ID, order.ShopID, order.CustomerID, order.OrderNumber, order.Status, order.TotalAmount,
		order.ShippingAddress, order.PaymentMethod, order.TrackingNumber, order.Carrier,
		order.CreatedAt, order.UpdatedAt,
	)
	return err
}

func (r *orderRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	query := `
		SELECT id, shop_id, customer_id, order_number, status, total_amount, shipping_address,
		       payment_method, tracking_number, carrier, created_at, updated_at
		FROM orders WHERE id = $1
	`
	order := &domain.Order{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.ShopID, &order.CustomerID, &order.OrderNumber, &order.Status,
		&order.TotalAmount, &order.ShippingAddress, &order.PaymentMethod, &order.TrackingNumber,
		&order.Carrier, &order.CreatedAt, &order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrOrderNotFound
	}
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *orderRepository) Update(ctx context.Context, order *domain.Order) error {
	query := `
		UPDATE orders SET status = $2, tracking_number = $3, carrier = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		order.ID, order.Status, order.TrackingNumber, order.Carrier, order.UpdatedAt,
	)
	return err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus) error {
	query := `UPDATE orders SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, orderID, status)
	return err
}

func (r *orderRepository) List(ctx context.Context, filter repository.OrderFilter) ([]*domain.Order, error) {
	query := `
		SELECT id, shop_id, customer_id, order_number, status, total_amount, shipping_address,
		       payment_method, tracking_number, carrier, created_at, updated_at
		FROM orders WHERE shop_id = $1
	`
	args := []interface{}{filter.ShopID}
	argCount := 1

	if filter.Status != nil {
		argCount++
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filter.Status)
	}

	if filter.DateFrom != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, *filter.DateFrom)
	}

	if filter.DateTo != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, *filter.DateTo)
	}

	sortBy := "created_at"
	if filter.SortBy != "" {
		sortBy = filter.SortBy
	}
	sortOrder := "DESC"
	if filter.SortOrder != "" && strings.ToUpper(filter.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o := &domain.Order{}
		if err := rows.Scan(
			&o.ID, &o.ShopID, &o.CustomerID, &o.OrderNumber, &o.Status, &o.TotalAmount,
			&o.ShippingAddress, &o.PaymentMethod, &o.TrackingNumber, &o.Carrier,
			&o.CreatedAt, &o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *orderRepository) AddItem(ctx context.Context, item *domain.OrderItem) error {
	query := `
		INSERT INTO order_items (id, order_id, product_id, product_name, quantity, unit_price, subtotal, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.OrderID, item.ProductID, item.ProductName,
		item.Quantity, item.UnitPrice, item.Subtotal, item.CreatedAt,
	)
	return err
}

func (r *orderRepository) GetItems(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, product_name, quantity, unit_price, subtotal, created_at
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
		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.ProductName,
			&item.Quantity, &item.UnitPrice, &item.Subtotal, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
