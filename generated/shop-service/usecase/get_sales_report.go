package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type GetSalesReportInput struct {
	ShopID     uuid.UUID
	ReportType string
	DateFrom   time.Time
	DateTo     time.Time
}

type GetSalesReportOutput struct {
	Reports []*domain.SalesReport
	Total   int
}

type GetSalesReportUsecase interface {
	Execute(ctx context.Context, input GetSalesReportInput) (*GetSalesReportOutput, error)
}

type getSalesReportUsecaseImpl struct {
	salesReportRepo domain.SalesReportRepository
}

func NewGetSalesReportUsecase(salesReportRepo domain.SalesReportRepository) GetSalesReportUsecase {
	return &getSalesReportUsecaseImpl{
		salesReportRepo: salesReportRepo,
	}
}

func (u *getSalesReportUsecaseImpl) Execute(ctx context.Context, input GetSalesReportInput) (*GetSalesReportOutput, error) {
	reports, err := u.salesReportRepo.List(ctx, 100, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales report: %w", err)
	}

	return &GetSalesReportOutput{
		Reports: reports,
		Total:   len(reports),
	}, nil
}
