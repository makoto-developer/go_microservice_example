package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type GetSalesReportInput struct {
	ShopID     uuid.UUID
	ReportType string
	DateFrom   time.Time
	DateTo     time.Time
}

type SalesData struct {
	Date              time.Time
	TotalSales        float64
	OrderCount        int
	AverageOrderValue float64
}

type SalesSummary struct {
	TotalSales        float64
	TotalOrders       int
	AverageOrderValue float64
}

type GetSalesReportOutput struct {
	ReportData []SalesData
	Summary    SalesSummary
}

type GetSalesReportUsecase interface {
	Execute(ctx context.Context, input GetSalesReportInput) (GetSalesReportOutput, error)
}

type getSalesReportUsecase struct {
	shopRepo  repository.ShopRepository
	salesRepo repository.SalesRepository
}

func NewGetSalesReportUsecase(shopRepo repository.ShopRepository, salesRepo repository.SalesRepository) GetSalesReportUsecase {
	return &getSalesReportUsecase{
		shopRepo:  shopRepo,
		salesRepo: salesRepo,
	}
}

func (u *getSalesReportUsecase) Execute(ctx context.Context, input GetSalesReportInput) (GetSalesReportOutput, error) {
	if _, err := u.shopRepo.GetByID(ctx, input.ShopID); err != nil {
		return GetSalesReportOutput{}, err
	}

	if input.DateFrom.After(input.DateTo) {
		return GetSalesReportOutput{}, domain.ErrInvalidDateRange
	}

	reports, err := u.salesRepo.GetByDateRange(ctx, input.ShopID, input.DateFrom, input.DateTo)
	if err != nil {
		return GetSalesReportOutput{}, err
	}

	if len(reports) == 0 {
		return GetSalesReportOutput{}, domain.ErrNoDataFound
	}

	var reportData []SalesData
	var totalSales float64
	var totalOrders int

	for _, report := range reports {
		reportData = append(reportData, SalesData{
			Date:              report.Date,
			TotalSales:        report.TotalSales,
			OrderCount:        report.OrderCount,
			AverageOrderValue: report.AverageOrderValue,
		})
		totalSales += report.TotalSales
		totalOrders += report.OrderCount
	}

	avgOrderValue := 0.0
	if totalOrders > 0 {
		avgOrderValue = totalSales / float64(totalOrders)
	}

	return GetSalesReportOutput{
		ReportData: reportData,
		Summary: SalesSummary{
			TotalSales:        totalSales,
			TotalOrders:       totalOrders,
			AverageOrderValue: avgOrderValue,
		},
	}, nil
}
