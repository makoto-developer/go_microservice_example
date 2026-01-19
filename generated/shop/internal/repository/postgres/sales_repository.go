package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type salesRepository struct {
	db *sql.DB
}

func NewSalesRepository(db *sql.DB) repository.SalesRepository {
	return &salesRepository{db: db}
}

func (r *salesRepository) Create(ctx context.Context, report *domain.SalesReport) error {
	query := `
		INSERT INTO sales_reports (id, shop_id, date, total_sales, order_count, average_order_value, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		report.ID, report.ShopID, report.Date, report.TotalSales,
		report.OrderCount, report.AverageOrderValue, report.CreatedAt,
	)
	return err
}

func (r *salesRepository) GetByShopAndDate(ctx context.Context, shopID uuid.UUID, date time.Time) (*domain.SalesReport, error) {
	query := `
		SELECT id, shop_id, date, total_sales, order_count, average_order_value, created_at
		FROM sales_reports WHERE shop_id = $1 AND date = $2
	`
	report := &domain.SalesReport{}
	err := r.db.QueryRowContext(ctx, query, shopID, date).Scan(
		&report.ID, &report.ShopID, &report.Date, &report.TotalSales,
		&report.OrderCount, &report.AverageOrderValue, &report.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (r *salesRepository) GetByDateRange(ctx context.Context, shopID uuid.UUID, dateFrom, dateTo time.Time) ([]*domain.SalesReport, error) {
	query := `
		SELECT id, shop_id, date, total_sales, order_count, average_order_value, created_at
		FROM sales_reports WHERE shop_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date ASC
	`
	rows, err := r.db.QueryContext(ctx, query, shopID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*domain.SalesReport
	for rows.Next() {
		report := &domain.SalesReport{}
		if err := rows.Scan(
			&report.ID, &report.ShopID, &report.Date, &report.TotalSales,
			&report.OrderCount, &report.AverageOrderValue, &report.CreatedAt,
		); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, rows.Err()
}

func (r *salesRepository) GenerateReport(ctx context.Context, shopID uuid.UUID, date time.Time) (*domain.SalesReport, error) {
	query := `
		SELECT COUNT(*) as order_count, COALESCE(SUM(total_amount), 0) as total_sales
		FROM orders
		WHERE shop_id = $1 AND DATE(created_at) = $2 AND status != 'CANCELLED'
	`
	var orderCount int
	var totalSales float64
	err := r.db.QueryRowContext(ctx, query, shopID, date).Scan(&orderCount, &totalSales)
	if err != nil {
		return nil, err
	}

	return domain.NewSalesReport(shopID, date, totalSales, orderCount), nil
}
