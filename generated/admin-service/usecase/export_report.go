package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ExportReportInput represents input for ExportReport
type ExportReportInput struct {
	ReportType string
	DateFrom date
	DateTo date
	Format string
}

// ExportReportUsecase defines the interface for ExportReport
type ExportReportUsecase interface {
	Execute(ctx context.Context, input ExportReportInput) error
}

type export_reportUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewExportReportUsecase creates a new instance
func NewExportReportUsecase() ExportReportUsecase {
	return &export_reportUsecaseImpl{}
}

// Execute executes ExportReport
func (u *export_reportUsecaseImpl) Execute(ctx context.Context, input ExportReportInput) error {
	// TODO: Implement business logic

	return nil
}
