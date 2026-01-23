package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/order/internal/domain"
)

type orderItemRepository struct {
	db *sql.DB
}

func NewOrderItemRepository(db *sql.DB) *orderItemRepository {
	return &orderItemRepository{db: db}
}

func (r *orderItemRepository) Create(ctx context.Context, item *domain.OrderItem) error {
	query := `INSERT INTO order_items (id, order_id, product_id, variation_id, quantity, unit_price, subtotal, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, item.ID, item.OrderID, item.ProductID, item.VariationID,
		item.Quantity, item.UnitPrice, item.Subtotal, item.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}
	return nil
}

func (r *orderItemRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.OrderItem, error) {
	query := `SELECT id, order_id, product_id, variation_id, quantity, unit_price, subtotal, created_at
		FROM order_items WHERE order_id = $1`
	
	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	defer rows.Close()
	
	var items []*domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.VariationID,
			&item.Quantity, &item.UnitPrice, &item.Subtotal, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}
