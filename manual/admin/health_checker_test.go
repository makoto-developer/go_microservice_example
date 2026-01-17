package admin

import (
	"context"
	"testing"
)

func TestHealthChecker_CheckAllServices(t *testing.T) {
	checker := NewHealthChecker()
	ctx := context.Background()

	results, err := checker.CheckAllServices(ctx)
	if err != nil {
		t.Fatalf("CheckAllServices failed: %v", err)
	}

	if len(results) != 12 {
		t.Errorf("Expected 12 services, got %d", len(results))
	}

	for _, result := range results {
		if result.ServiceName == "" {
			t.Error("Expected non-empty service name")
		}

		if result.Status != HealthStatusHealthy {
			t.Errorf("Expected healthy status for %s, got %s", result.ServiceName, result.Status)
		}
	}
}

func TestHealthChecker_CheckDatabase(t *testing.T) {
	checker := NewHealthChecker()
	ctx := context.Background()

	health, err := checker.CheckDatabase(ctx, "auth-service")
	if err != nil {
		t.Fatalf("CheckDatabase failed: %v", err)
	}

	if health.Status != HealthStatusHealthy {
		t.Errorf("Expected healthy status, got %s", health.Status)
	}

	if health.ServiceName != "auth-service-db" {
		t.Errorf("Expected service name auth-service-db, got %s", health.ServiceName)
	}
}

func TestHealthChecker_CheckCache(t *testing.T) {
	checker := NewHealthChecker()
	ctx := context.Background()

	health, err := checker.CheckCache(ctx, "auth-service")
	if err != nil {
		t.Fatalf("CheckCache failed: %v", err)
	}

	if health.Status != HealthStatusHealthy {
		t.Errorf("Expected healthy status, got %s", health.Status)
	}
}

func TestHealthChecker_CheckMessageQueue(t *testing.T) {
	checker := NewHealthChecker()
	ctx := context.Background()

	health, err := checker.CheckMessageQueue(ctx)
	if err != nil {
		t.Fatalf("CheckMessageQueue failed: %v", err)
	}

	if health.Status != HealthStatusHealthy {
		t.Errorf("Expected healthy status, got %s", health.Status)
	}

	if health.ServiceName != "rabbitmq" {
		t.Errorf("Expected service name rabbitmq, got %s", health.ServiceName)
	}
}

func TestHealthChecker_GetSystemMetrics(t *testing.T) {
	checker := NewHealthChecker()
	ctx := context.Background()

	metrics, err := checker.GetSystemMetrics(ctx)
	if err != nil {
		t.Fatalf("GetSystemMetrics failed: %v", err)
	}

	requiredMetrics := []string{
		"cpu_usage",
		"memory_usage",
		"disk_usage",
		"network_rx",
		"network_tx",
		"uptime",
	}

	for _, metric := range requiredMetrics {
		if _, exists := metrics[metric]; !exists {
			t.Errorf("Expected metric %s to exist", metric)
		}
	}
}
