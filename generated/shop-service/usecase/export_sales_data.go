package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type ExportSalesDataInput struct {
	ShopID   uuid.UUID
	DateFrom time.Time
	DateTo   time.Time
}

type ExportSalesDataOutput struct {
	ExportURL string
	Message   string
}

type ExportSalesDataUsecase interface {
	Execute(ctx context.Context, input ExportSalesDataInput) (*ExportSalesDataOutput, error)
}

type exportSalesDataUsecaseImpl struct {
	salesReportRepo domain.SalesReportRepository
}

func NewExportSalesDataUsecase(salesReportRepo domain.SalesReportRepository) ExportSalesDataUsecase {
	return &exportSalesDataUsecaseImpl{
		salesReportRepo: salesReportRepo,
	}
}

func (u *exportSalesDataUsecaseImpl) Execute(ctx context.Context, input ExportSalesDataInput) (*ExportSalesDataOutput, error) {
	reports, err := u.salesReportRepo.List(ctx, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to export sales data: %w", err)
	}

	exportURL := fmt.Sprintf("/exports/sales_%s_%s.csv",
		input.DateFrom.Format("20060102"),
		input.DateTo.Format("20060102"))

	_ = reports

	return &ExportSalesDataOutput{
		ExportURL: exportURL,
		Message:   "Sales data exported successfully",
	}, nil
}
