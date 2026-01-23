package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory/internal/domain"
)

type inventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *inventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) Create(ctx context.Context, inventory *domain.Inventory) error {
	query := `INSERT INTO inventories (id, product_id, variation_id, shop_id, quantity, reserved_quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, inventory.ID, inventory.ProductID, inventory.VariationID,
		inventory.ShopID, inventory.Quantity, inventory.ReservedQuantity, inventory.CreatedAt, inventory.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}
	return nil
}

func (r *inventoryRepository) GetByProductID(ctx context.Context, productID uuid.UUID, variationID *uuid.UUID) (*domain.Inventory, error) {
	query := `SELECT id, product_id, variation_id, shop_id, quantity, reserved_quantity, created_at, updated_at
		FROM inventories WHERE product_id = $1 AND ($2::uuid IS NULL OR variation_id = $2)`
	
	var inv domain.Inventory
	err := r.db.QueryRowContext(ctx, query, productID, variationID).Scan(
		&inv.ID, &inv.ProductID, &inv.VariationID, &inv.ShopID, &inv.Quantity, &inv.ReservedQuantity, &inv.CreatedAt, &inv.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("inventory not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	return &inv, nil
}

func (r *inventoryRepository) Update(ctx context.Context, inventory *domain.Inventory) error {
	query := `UPDATE inventories SET quantity = $2, reserved_quantity = $3, updated_at = $4 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, inventory.ID, inventory.Quantity, inventory.ReservedQuantity, inventory.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update inventory: %w", err)
	}
	return nil
}

func (r *inventoryRepository) UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `UPDATE inventories SET quantity = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, quantity)
	if err != nil {
		return fmt.Errorf("failed to update quantity: %w", err)
	}
	return nil
}

func (r *inventoryRepository) Reserve(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `UPDATE inventories SET reserved_quantity = reserved_quantity + $2, updated_at = NOW()
		WHERE id = $1 AND (quantity - reserved_quantity) >= $2`
	result, err := r.db.ExecContext(ctx, query, id, quantity)
	if err != nil {
		return fmt.Errorf("failed to reserve: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("insufficient inventory")
	}
	return nil
}

func (r *inventoryRepository) Release(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `UPDATE inventories SET reserved_quantity = reserved_quantity - $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, quantity)
	if err != nil {
		return fmt.Errorf("failed to release: %w", err)
	}
	return nil
}
