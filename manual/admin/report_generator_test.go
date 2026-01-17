package admin

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestReportGenerator_GenerateSalesReport_CSV(t *testing.T) {
	generator := NewReportGenerator()
	ctx := context.Background()

	req := SalesReportRequest{
		StartDate: time.Now().AddDate(0, 0, -7),
		EndDate:   time.Now(),
		ShopID:    "shop-123",
		Format:    ReportFormatCSV,
	}

	filename, err := generator.GenerateSalesReport(ctx, req)
	if err != nil {
		t.Fatalf("GenerateSalesReport failed: %v", err)
	}

	if filename == "" {
		t.Error("Expected non-empty filename")
	}

	// ファイルが存在するか確認
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Report file does not exist: %s", filename)
	}

	// クリーンアップ
	os.Remove(filename)
}

func TestReportGenerator_GenerateSalesReport_PDF(t *testing.T) {
	generator := NewReportGenerator()
	ctx := context.Background()

	req := SalesReportRequest{
		StartDate: time.Now().AddDate(0, 0, -7),
		EndDate:   time.Now(),
		ShopID:    "shop-123",
		Format:    ReportFormatPDF,
	}

	filename, err := generator.GenerateSalesReport(ctx, req)
	if err != nil {
		t.Fatalf("GenerateSalesReport failed: %v", err)
	}

	if filename == "" {
		t.Error("Expected non-empty filename")
	}

	// クリーンアップ
	os.Remove(filename)
}

func TestReportGenerator_GenerateUserReport(t *testing.T) {
	generator := NewReportGenerator()
	ctx := context.Background()

	req := UserReportRequest{
		StartDate: time.Now().AddDate(0, 0, -30),
		EndDate:   time.Now(),
		Role:      "CUSTOMER",
		Format:    ReportFormatCSV,
	}

	filename, err := generator.GenerateUserReport(ctx, req)
	if err != nil {
		t.Fatalf("GenerateUserReport failed: %v", err)
	}

	if filename == "" {
		t.Error("Expected non-empty filename")
	}

	// クリーンアップ
	os.Remove(filename)
}

func TestReportGenerator_GenerateOrderReport(t *testing.T) {
	generator := NewReportGenerator()
	ctx := context.Background()

	req := OrderReportRequest{
		StartDate: time.Now().AddDate(0, 0, -7),
		EndDate:   time.Now(),
		Status:    "DELIVERED",
		Format:    ReportFormatCSV,
	}

	filename, err := generator.GenerateOrderReport(ctx, req)
	if err != nil {
		t.Fatalf("GenerateOrderReport failed: %v", err)
	}

	if filename == "" {
		t.Error("Expected non-empty filename")
	}

	// クリーンアップ
	os.Remove(filename)
}
