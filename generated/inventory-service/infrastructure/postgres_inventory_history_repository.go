package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory-service/domain"
)

type PostgresInventoryHistoryRepository struct {
	db *sql.DB
}

func NewPostgresInventoryHistoryRepository(db *sql.DB) *PostgresInventoryHistoryRepository {
	return &PostgresInventoryHistoryRepository{db: db}
}

func (r *PostgresInventoryHistoryRepository) Create(ctx context.Context, history *domain.InventoryHistory) error {
	query := `
		INSERT INTO inventory_history (
			id, inventory_id, change_type, change_quantity,
			quantity_before, quantity_after, reason, operator, order_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		history.ID,
		history.InventoryID,
		history.ChangeType,
		history.ChangeQuantity,
		history.QuantityBefore,
		history.QuantityAfter,
		history.Reason,
		history.Operator,
		history.OrderID,
		history.CreatedAt,
	)

	return err
}

func (r *PostgresInventoryHistoryRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, history *domain.InventoryHistory) error {
	query := `
		INSERT INTO inventory_history (
			id, inventory_id, change_type, change_quantity,
			quantity_before, quantity_after, reason, operator, order_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := tx.ExecContext(ctx, query,
		history.ID,
		history.ChangeType,
		history.ChangeQuantity,
		history.QuantityBefore,
		history.QuantityAfter,
		history.Reason,
		history.Operator,
		history.OrderID,
		history.CreatedAt,
	)

	return err
}

func (r *PostgresInventoryHistoryRepository) FindByInventoryID(
	ctx context.Context,
	inventoryID uuid.UUID,
	dateFrom, dateTo time.Time,
	changeTypes []domain.ChangeType,
	page, pageSize int,
) ([]*domain.InventoryHistory, int, error) {
	// Count total
	countQuery := `
		SELECT COUNT(*)
		FROM inventory_history
		WHERE inventory_id = $1
			AND created_at >= $2
			AND created_at <= $3
			AND ($4::text[] IS NULL OR change_type = ANY($4))
	`

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, inventoryID, dateFrom, dateTo, changeTypes).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch data
	query := `
		SELECT id, inventory_id, change_type, change_quantity,
			quantity_before, quantity_after, reason, operator, order_id, created_at
		FROM inventory_history
		WHERE inventory_id = $1
			AND created_at >= $2
			AND created_at <= $3
			AND ($4::text[] IS NULL OR change_type = ANY($4))
		ORDER BY created_at DESC
		LIMIT $5 OFFSET $6
	`

	offset := (page - 1) * pageSize
	rows, err := r.db.QueryContext(ctx, query, inventoryID, dateFrom, dateTo, changeTypes, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var histories []*domain.InventoryHistory
	for rows.Next() {
		var history domain.InventoryHistory
		var orderID sql.NullString

		err := rows.Scan(
			&history.ID,
			&history.InventoryID,
			&history.ChangeType,
			&history.ChangeQuantity,
			&history.QuantityBefore,
			&history.QuantityAfter,
			&history.Reason,
			&history.Operator,
			&orderID,
			&history.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if orderID.Valid {
			oid, _ := uuid.Parse(orderID.String)
			history.OrderID = &oid
		}

		histories = append(histories, &history)
	}

	return histories, total, rows.Err()
}
