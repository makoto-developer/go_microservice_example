package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type PostgresSalesReportRepository struct {
	db *sql.DB
}

func NewPostgresSalesReportRepository(db *sql.DB) *PostgresSalesReportRepository {
	return &PostgresSalesReportRepository{db: db}
}

func (r *PostgresSalesReportRepository) Create(ctx context.Context, salesReport *domain.SalesReport) error {
	query := `
		INSERT INTO sales_reports (
			id, shop_id, date, total_sales, order_count,
			average_order_value, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		salesReport.Id,
		salesReport.ShopId,
		salesReport.Date,
		salesReport.TotalSales,
		salesReport.OrderCount,
		salesReport.AverageOrderValue,
		salesReport.CreatedAt,
	)

	return err
}

func (r *PostgresSalesReportRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.SalesReport, error) {
	query := `
		SELECT id, shop_id, date, total_sales, order_count,
			average_order_value, created_at
		FROM sales_reports
		WHERE id = $1
	`

	var salesReport domain.SalesReport

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&salesReport.Id,
		&salesReport.ShopId,
		&salesReport.Date,
		&salesReport.TotalSales,
		&salesReport.OrderCount,
		&salesReport.AverageOrderValue,
		&salesReport.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("sales report not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	return &salesReport, nil
}

func (r *PostgresSalesReportRepository) Update(ctx context.Context, salesReport *domain.SalesReport) error {
	query := `
		UPDATE sales_reports
		SET shop_id = $1,
			date = $2,
			total_sales = $3,
			order_count = $4,
			average_order_value = $5
		WHERE id = $6
	`

	_, err := r.db.ExecContext(ctx, query,
		salesReport.ShopId,
		salesReport.Date,
		salesReport.TotalSales,
		salesReport.OrderCount,
		salesReport.AverageOrderValue,
		salesReport.Id,
	)

	return err
}

func (r *PostgresSalesReportRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM sales_reports WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresSalesReportRepository) List(ctx context.Context, limit, offset int) ([]*domain.SalesReport, error) {
	query := `
		SELECT id, shop_id, date, total_sales, order_count,
			average_order_value, created_at
		FROM sales_reports
		ORDER BY date DESC, created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*domain.SalesReport
	for rows.Next() {
		var report domain.SalesReport

		err := rows.Scan(
			&report.Id,
			&report.ShopId,
			&report.Date,
			&report.TotalSales,
			&report.OrderCount,
			&report.AverageOrderValue,
			&report.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		reports = append(reports, &report)
	}

	return reports, rows.Err()
}

// BeginTx starts a transaction
func (r *PostgresSalesReportRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
