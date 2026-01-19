package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory-service/domain"
)

type PostgresInventoryRepository struct {
	db *sql.DB
}

func NewPostgresInventoryRepository(db *sql.DB) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{db: db}
}

func (r *PostgresInventoryRepository) Create(ctx context.Context, inventory *domain.Inventory) error {
	query := `
		INSERT INTO inventories (
			id, product_id, variation_id, shop_id,
			quantity, reserved_quantity, available_quantity, alert_threshold,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		inventory.ID,
		inventory.ProductID,
		inventory.VariationID,
		inventory.ShopID,
		inventory.Quantity,
		inventory.ReservedQuantity,
		inventory.AvailableQuantity,
		inventory.AlertThreshold,
		inventory.CreatedAt,
		inventory.UpdatedAt,
	)

	return err
}

func (r *PostgresInventoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	query := `
		SELECT id, product_id, variation_id, shop_id,
			quantity, reserved_quantity, available_quantity, alert_threshold,
			last_alerted_at, created_at, updated_at
		FROM inventories
		WHERE id = $1
	`

	var inventory domain.Inventory
	var variationID sql.NullString
	var lastAlertedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&variationID,
		&inventory.ShopID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.AvailableQuantity,
		&inventory.AlertThreshold,
		&lastAlertedAt,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("inventory not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	if variationID.Valid {
		vid, _ := uuid.Parse(variationID.String)
		inventory.VariationID = &vid
	}
	if lastAlertedAt.Valid {
		inventory.LastAlertedAt = &lastAlertedAt.Time
	}

	return &inventory, nil
}

func (r *PostgresInventoryRepository) FindByProductID(ctx context.Context, productID, variationID uuid.UUID) (*domain.Inventory, error) {
	query := `
		SELECT id, product_id, variation_id, shop_id,
			quantity, reserved_quantity, available_quantity, alert_threshold,
			last_alerted_at, created_at, updated_at
		FROM inventories
		WHERE product_id = $1 AND (variation_id = $2 OR (variation_id IS NULL AND $2 = '00000000-0000-0000-0000-000000000000'))
	`

	var inventory domain.Inventory
	var variationIDNull sql.NullString
	var lastAlertedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, productID, variationID).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&variationIDNull,
		&inventory.ShopID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.AvailableQuantity,
		&inventory.AlertThreshold,
		&lastAlertedAt,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("inventory not found for product: %s", productID)
	}
	if err != nil {
		return nil, err
	}

	if variationIDNull.Valid {
		vid, _ := uuid.Parse(variationIDNull.String)
		inventory.VariationID = &vid
	}
	if lastAlertedAt.Valid {
		inventory.LastAlertedAt = &lastAlertedAt.Time
	}

	return &inventory, nil
}

func (r *PostgresInventoryRepository) FindByProductIDs(ctx context.Context, productIDs []uuid.UUID) ([]*domain.Inventory, error) {
	if len(productIDs) == 0 {
		return []*domain.Inventory{}, nil
	}

	query := `
		SELECT id, product_id, variation_id, shop_id,
			quantity, reserved_quantity, available_quantity, alert_threshold,
			last_alerted_at, created_at, updated_at
		FROM inventories
		WHERE product_id = ANY($1)
	`

	rows, err := r.db.QueryContext(ctx, query, productIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []*domain.Inventory
	for rows.Next() {
		var inventory domain.Inventory
		var variationID sql.NullString
		var lastAlertedAt sql.NullTime

		err := rows.Scan(
			&inventory.ID,
			&inventory.ProductID,
			&variationID,
			&inventory.ShopID,
			&inventory.Quantity,
			&inventory.ReservedQuantity,
			&inventory.AvailableQuantity,
			&inventory.AlertThreshold,
			&lastAlertedAt,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if variationID.Valid {
			vid, _ := uuid.Parse(variationID.String)
			inventory.VariationID = &vid
		}
		if lastAlertedAt.Valid {
			inventory.LastAlertedAt = &lastAlertedAt.Time
		}

		inventories = append(inventories, &inventory)
	}

	return inventories, rows.Err()
}

func (r *PostgresInventoryRepository) Update(ctx context.Context, inventory *domain.Inventory) error {
	query := `
		UPDATE inventories
		SET quantity = $1,
			reserved_quantity = $2,
			available_quantity = $3,
			alert_threshold = $4,
			last_alerted_at = $5,
			updated_at = $6
		WHERE id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		inventory.Quantity,
		inventory.ReservedQuantity,
		inventory.AvailableQuantity,
		inventory.AlertThreshold,
		inventory.LastAlertedAt,
		time.Now(),
		inventory.ID,
	)

	return err
}

func (r *PostgresInventoryRepository) UpdateWithTx(ctx context.Context, tx *sql.Tx, inventory *domain.Inventory) error {
	query := `
		UPDATE inventories
		SET quantity = $1,
			reserved_quantity = $2,
			available_quantity = $3,
			alert_threshold = $4,
			last_alerted_at = $5,
			updated_at = $6
		WHERE id = $7
	`

	_, err := tx.ExecContext(ctx, query,
		inventory.Quantity,
		inventory.ReservedQuantity,
		inventory.AvailableQuantity,
		inventory.AlertThreshold,
		inventory.LastAlertedAt,
		time.Now(),
		inventory.ID,
	)

	return err
}

func (r *PostgresInventoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM inventories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// トランザクション管理
func (r *PostgresInventoryRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
