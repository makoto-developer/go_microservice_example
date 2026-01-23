package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
)

type SalesRepository interface {
	Create(ctx context.Context, report *domain.SalesReport) error
	GetByShopAndDate(ctx context.Context, shopID uuid.UUID, date time.Time) (*domain.SalesReport, error)
	GetByDateRange(ctx context.Context, shopID uuid.UUID, dateFrom, dateTo time.Time) ([]*domain.SalesReport, error)
	GenerateReport(ctx context.Context, shopID uuid.UUID, date time.Time) (*domain.SalesReport, error)
}
