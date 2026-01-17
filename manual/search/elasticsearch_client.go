package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// ElasticsearchClient はElasticsearchクライアント
type ElasticsearchClient struct {
	client *elasticsearch.Client
}

// NewElasticsearchClient はElasticsearchClientを初期化
func NewElasticsearchClient(addresses []string) (*ElasticsearchClient, error) {
	cfg := elasticsearch.Config{
		Addresses: addresses,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	return &ElasticsearchClient{
		client: client,
	}, nil
}

// ProductDocument は商品検索ドキュメント
type ProductDocument struct {
	ProductID   string    `json:"product_id"`
	ShopID      string    `json:"shop_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	Rating      float64   `json:"rating"`
	ReviewCount int       `json:"review_count"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SearchRequest は検索リクエスト
type SearchRequest struct {
	Query      string
	Categories []string
	MinPrice   int64
	MaxPrice   int64
	MinRating  float64
	Tags       []string
	InStock    bool
	SortBy     string // "relevance", "price_asc", "price_desc", "rating", "newest"
	Offset     int
	Limit      int
}

// SearchResponse は検索レスポンス
type SearchResponse struct {
	Total      int64
	Results    []ProductDocument
	Facets     map[string]interface{}
	Took       int64
	Aggregations map[string]interface{}
}

// IndexProduct は商品をインデックス
func (c *ElasticsearchClient) IndexProduct(ctx context.Context, product ProductDocument) error {
	data, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      "products",
		DocumentID: product.ProductID,
		Body:       strings.NewReader(string(data)),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to index product: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.Status())
	}

	return nil
}

// SearchProducts は商品を検索
func (c *ElasticsearchClient) SearchProducts(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	// Elasticsearch クエリを構築
	query := c.buildSearchQuery(req)

	// 検索実行
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex("products"),
		c.client.Search.WithBody(strings.NewReader(buf.String())),
		c.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.Status())
	}

	// レスポンスをパース
	var esResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&esResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 結果を抽出
	hits := esResponse["hits"].(map[string]interface{})
	total := int64(hits["total"].(map[string]interface{})["value"].(float64))
	took := int64(esResponse["took"].(float64))

	results := make([]ProductDocument, 0)
	for _, hit := range hits["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		var product ProductDocument
		sourceJSON, _ := json.Marshal(source)
		json.Unmarshal(sourceJSON, &product)
		results = append(results, product)
	}

	return &SearchResponse{
		Total:   total,
		Results: results,
		Took:    took,
	}, nil
}

// buildSearchQuery は検索クエリを構築
func (c *ElasticsearchClient) buildSearchQuery(req SearchRequest) map[string]interface{} {
	must := make([]interface{}, 0)

	// テキスト検索
	if req.Query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"name^2", "description", "tags"},
				"type":   "best_fields",
			},
		})
	}

	// フィルター条件
	filter := make([]interface{}, 0)

	if len(req.Categories) > 0 {
		filter = append(filter, map[string]interface{}{
			"terms": map[string]interface{}{
				"category": req.Categories,
			},
		})
	}

	if req.MinPrice > 0 || req.MaxPrice > 0 {
		priceRange := make(map[string]interface{})
		if req.MinPrice > 0 {
			priceRange["gte"] = req.MinPrice
		}
		if req.MaxPrice > 0 {
			priceRange["lte"] = req.MaxPrice
		}
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"price": priceRange,
			},
		})
	}

	if req.MinRating > 0 {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"rating": map[string]interface{}{
					"gte": req.MinRating,
				},
			},
		})
	}

	if req.InStock {
		filter = append(filter, map[string]interface{}{
			"range": map[string]interface{}{
				"stock": map[string]interface{}{
					"gt": 0,
				},
			},
		})
	}

	// ソート
	sort := c.buildSort(req.SortBy)

	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   must,
				"filter": filter,
			},
		},
		"sort": sort,
		"from": req.Offset,
		"size": req.Limit,
	}
}

// buildSort はソート条件を構築
func (c *ElasticsearchClient) buildSort(sortBy string) []interface{} {
	switch sortBy {
	case "price_asc":
		return []interface{}{
			map[string]interface{}{"price": "asc"},
		}
	case "price_desc":
		return []interface{}{
			map[string]interface{}{"price": "desc"},
		}
	case "rating":
		return []interface{}{
			map[string]interface{}{"rating": "desc"},
		}
	case "newest":
		return []interface{}{
			map[string]interface{}{"created_at": "desc"},
		}
	default: // relevance
		return []interface{}{
			"_score",
		}
	}
}

// DeleteProduct は商品をインデックスから削除
func (c *ElasticsearchClient) DeleteProduct(ctx context.Context, productID string) error {
	req := esapi.DeleteRequest{
		Index:      "products",
		DocumentID: productID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("elasticsearch error: %s", res.Status())
	}

	return nil
}

// CreateIndex はインデックスを作成
func (c *ElasticsearchClient) CreateIndex(ctx context.Context) error {
	mapping := `{
		"mappings": {
			"properties": {
				"product_id": {"type": "keyword"},
				"shop_id": {"type": "keyword"},
				"name": {
					"type": "text",
					"analyzer": "kuromoji"
				},
				"description": {
					"type": "text",
					"analyzer": "kuromoji"
				},
				"price": {"type": "long"},
				"category": {"type": "keyword"},
				"tags": {"type": "keyword"},
				"rating": {"type": "float"},
				"review_count": {"type": "integer"},
				"stock": {"type": "integer"},
				"created_at": {"type": "date"},
				"updated_at": {"type": "date"}
			}
		},
		"settings": {
			"analysis": {
				"analyzer": {
					"kuromoji": {
						"type": "custom",
						"tokenizer": "kuromoji_tokenizer"
					}
				}
			}
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: "products",
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 400 { // 400 = already exists
		return fmt.Errorf("elasticsearch error: %s", res.Status())
	}

	return nil
}
