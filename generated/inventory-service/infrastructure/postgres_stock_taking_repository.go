package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory-service/domain"
)

type PostgresStockTakingRepository struct {
	db *sql.DB
}

func NewPostgresStockTakingRepository(db *sql.DB) *PostgresStockTakingRepository {
	return &PostgresStockTakingRepository{db: db}
}

func (r *PostgresStockTakingRepository) Create(ctx context.Context, stockTaking *domain.StockTaking) error {
	query := `
		INSERT INTO stock_takings (
			id, inventory_id, shop_id, system_quantity, actual_quantity,
			difference, difference_reason, operator, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		stockTaking.ID,
		stockTaking.InventoryID,
		stockTaking.ShopID,
		stockTaking.SystemQuantity,
		stockTaking.ActualQuantity,
		stockTaking.Difference,
		stockTaking.DifferenceReason,
		stockTaking.Operator,
		stockTaking.CreatedAt,
	)

	return err
}

func (r *PostgresStockTakingRepository) FindByShopID(
	ctx context.Context,
	shopID uuid.UUID,
	dateFrom, dateTo time.Time,
) ([]*domain.StockTaking, error) {
	query := `
		SELECT id, inventory_id, shop_id, system_quantity, actual_quantity,
			difference, difference_reason, operator, created_at
		FROM stock_takings
		WHERE shop_id = $1
			AND created_at >= $2
			AND created_at <= $3
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, shopID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stockTakings []*domain.StockTaking
	for rows.Next() {
		var st domain.StockTaking
		var differenceReason sql.NullString

		err := rows.Scan(
			&st.ID,
			&st.InventoryID,
			&st.ShopID,
			&st.SystemQuantity,
			&st.ActualQuantity,
			&st.Difference,
			&differenceReason,
			&st.Operator,
			&st.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if differenceReason.Valid {
			st.DifferenceReason = &differenceReason.String
		}

		stockTakings = append(stockTakings, &st)
	}

	return stockTakings, rows.Err()
}
