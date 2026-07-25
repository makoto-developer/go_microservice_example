package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) repository.CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddItem(ctx context.Context, item *domain.CartItem) error {
	query := `
		INSERT INTO cart_items (id, customer_id, product_id, variation_id, quantity, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.CustomerID, item.ProductID, item.VariationID,
		item.Quantity, item.ExpiresAt, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *cartRepository) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.CartItem, error) {
	query := `
		SELECT id, customer_id, product_id, variation_id, quantity, expires_at, created_at, updated_at
		FROM cart_items WHERE customer_id = $1 AND expires_at > NOW()
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.CartItem
	for rows.Next() {
		item := &domain.CartItem{}
		if err := rows.Scan(
			&item.ID, &item.CustomerID, &item.ProductID, &item.VariationID,
			&item.Quantity, &item.ExpiresAt, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *cartRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.CartItem, error) {
	query := `
		SELECT id, customer_id, product_id, variation_id, quantity, expires_at, created_at, updated_at
		FROM cart_items WHERE id = $1
	`
	item := &domain.CartItem{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.CustomerID, &item.ProductID, &item.VariationID,
		&item.Quantity, &item.ExpiresAt, &item.CreatedAt, &item.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrCartItemNotFound
	}
	return item, err
}

func (r *cartRepository) UpdateQuantity(ctx context.Context, cartItemID uuid.UUID, quantity int) error {
	query := `UPDATE cart_items SET quantity = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, cartItemID, quantity)
	return err
}

func (r *cartRepository) RemoveItem(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *cartRepository) ClearCart(ctx context.Context, customerID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE customer_id = $1`
	_, err := r.db.ExecContext(ctx, query, customerID)
	return err
}

func (r *cartRepository) AddGuestItem(ctx context.Context, item *domain.GuestCartItem) error {
	query := `
		INSERT INTO guest_cart_items (id, session_id, product_id, variation_id, quantity, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.SessionID, item.ProductID, item.VariationID,
		item.Quantity, item.ExpiresAt, item.CreatedAt,
	)
	return err
}

func (r *cartRepository) GetBySessionID(ctx context.Context, sessionID string) ([]*domain.GuestCartItem, error) {
	query := `
		SELECT id, session_id, product_id, variation_id, quantity, expires_at, created_at
		FROM guest_cart_items WHERE session_id = $1 AND expires_at > NOW()
	`
	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.GuestCartItem
	for rows.Next() {
		item := &domain.GuestCartItem{}
		if err := rows.Scan(
			&item.ID, &item.SessionID, &item.ProductID, &item.VariationID,
			&item.Quantity, &item.ExpiresAt, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *cartRepository) ClearGuestCart(ctx context.Context, sessionID string) error {
	query := `DELETE FROM guest_cart_items WHERE session_id = $1`
	_, err := r.db.ExecContext(ctx, query, sessionID)
	return err
}
