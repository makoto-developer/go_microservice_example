package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/inventory/internal/repository"
)

type inventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) repository.InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	query := `
		SELECT id, product_id, shop_id, quantity, reserved_quantity, created_at, updated_at
		FROM inventories
		WHERE id = $1
	`

	var inventory domain.Inventory
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.ShopID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (r *inventoryRepository) GetByProductAndShop(ctx context.Context, productID, shopID uuid.UUID) (*domain.Inventory, error) {
	query := `
		SELECT id, product_id, shop_id, quantity, reserved_quantity, created_at, updated_at
		FROM inventories
		WHERE product_id = $1 AND shop_id = $2
	`

	var inventory domain.Inventory
	err := r.db.QueryRowContext(ctx, query, productID, shopID).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.ShopID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &inventory, nil
}

func (r *inventoryRepository) Create(ctx context.Context, inventory *domain.Inventory) error {
	query := `
		INSERT INTO inventories (id, product_id, shop_id, quantity, reserved_quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		inventory.ID,
		inventory.ProductID,
		inventory.ShopID,
		inventory.Quantity,
		inventory.ReservedQuantity,
		inventory.CreatedAt,
		inventory.UpdatedAt,
	)

	return err
}

func (r *inventoryRepository) Update(ctx context.Context, inventory *domain.Inventory) error {
	query := `
		UPDATE inventories
		SET quantity = $1, reserved_quantity = $2, updated_at = $3
		WHERE id = $4
	`

	inventory.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query,
		inventory.Quantity,
		inventory.ReservedQuantity,
		inventory.UpdatedAt,
		inventory.ID,
	)

	return err
}

func (r *inventoryRepository) UpdateQuantity(ctx context.Context, id uuid.UUID, quantity, reservedQuantity int) error {
	query := `
		UPDATE inventories
		SET quantity = $1, reserved_quantity = $2, updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, quantity, reservedQuantity, time.Now(), id)
	return err
}

func (r *inventoryRepository) Reserve(ctx context.Context, inventoryID uuid.UUID, quantity int, orderID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update inventory
	query := `
		UPDATE inventories
		SET reserved_quantity = reserved_quantity + $1, updated_at = $2
		WHERE id = $3 AND (quantity - reserved_quantity) >= $1
	`

	result, err := tx.ExecContext(ctx, query, quantity, time.Now(), inventoryID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	// Create reservation record
	reservationQuery := `
		INSERT INTO stock_reservations (id, inventory_id, order_id, quantity, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	_, err = tx.ExecContext(ctx, reservationQuery,
		uuid.New(),
		inventoryID,
		orderID,
		quantity,
		"reserved",
		now,
		now,
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *inventoryRepository) Release(ctx context.Context, orderID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get reservation
	var inventoryID uuid.UUID
	var quantity int
	query := `
		SELECT inventory_id, quantity
		FROM stock_reservations
		WHERE order_id = $1 AND status = 'reserved'
	`

	err = tx.QueryRowContext(ctx, query, orderID).Scan(&inventoryID, &quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	// Update inventory
	updateQuery := `
		UPDATE inventories
		SET reserved_quantity = reserved_quantity - $1, updated_at = $2
		WHERE id = $3
	`

	_, err = tx.ExecContext(ctx, updateQuery, quantity, time.Now(), inventoryID)
	if err != nil {
		return err
	}

	// Update reservation status
	reservationQuery := `
		UPDATE stock_reservations
		SET status = 'released', updated_at = $1
		WHERE order_id = $2 AND status = 'reserved'
	`

	_, err = tx.ExecContext(ctx, reservationQuery, time.Now(), orderID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
