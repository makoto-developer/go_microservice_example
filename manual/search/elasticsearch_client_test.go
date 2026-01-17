package search

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Note: これらのテストは実際のElasticsearchが必要
// docker-compose upでElasticsearchを起動してからテスト実行

func TestElasticsearchClient_IndexProduct(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client, err := NewElasticsearchClient([]string{"http://localhost:9200"})
	if err != nil {
		t.Fatalf("NewElasticsearchClient failed: %v", err)
	}

	ctx := context.Background()

	// インデックス作成
	if err := client.CreateIndex(ctx); err != nil {
		t.Logf("CreateIndex warning (may already exist): %v", err)
	}

	product := ProductDocument{
		ProductID:   uuid.New().String(),
		ShopID:      uuid.New().String(),
		Name:        "テスト商品",
		Description: "これはテスト用の商品です",
		Price:       1000,
		Category:    "electronics",
		Tags:        []string{"test", "sample"},
		Rating:      4.5,
		ReviewCount: 10,
		Stock:       100,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = client.IndexProduct(ctx, product)
	if err != nil {
		t.Fatalf("IndexProduct failed: %v", err)
	}
}

func TestElasticsearchClient_SearchProducts(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client, err := NewElasticsearchClient([]string{"http://localhost:9200"})
	if err != nil {
		t.Fatalf("NewElasticsearchClient failed: %v", err)
	}

	ctx := context.Background()

	req := SearchRequest{
		Query:  "テスト",
		Offset: 0,
		Limit:  10,
		SortBy: "relevance",
	}

	resp, err := client.SearchProducts(ctx, req)
	if err != nil {
		t.Fatalf("SearchProducts failed: %v", err)
	}

	if resp.Total < 0 {
		t.Error("Expected non-negative total")
	}
}

func TestElasticsearchClient_DeleteProduct(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	client, err := NewElasticsearchClient([]string{"http://localhost:9200"})
	if err != nil {
		t.Fatalf("NewElasticsearchClient failed: %v", err)
	}

	ctx := context.Background()

	productID := uuid.New().String()
	err = client.DeleteProduct(ctx, productID)
	if err != nil {
		t.Fatalf("DeleteProduct failed: %v", err)
	}
}
