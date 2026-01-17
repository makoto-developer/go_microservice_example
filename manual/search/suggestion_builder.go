package search

import (
	"context"
	"strings"
)

// SuggestionBuilder は検索サジェストビルダー
type SuggestionBuilder struct {
	esClient *ElasticsearchClient
}

// NewSuggestionBuilder はSuggestionBuilderを初期化
func NewSuggestionBuilder(esClient *ElasticsearchClient) *SuggestionBuilder {
	return &SuggestionBuilder{
		esClient: esClient,
	}
}

// Suggestion はサジェスト結果
type Suggestion struct {
	Text  string
	Score float64
}

// GetSuggestions はサジェストを取得
func (s *SuggestionBuilder) GetSuggestions(ctx context.Context, prefix string, limit int) ([]Suggestion, error) {
	// Elasticsearchのcompletion suggesterを使用
	// 簡易実装: 前方一致検索
	req := SearchRequest{
		Query:  prefix,
		Limit:  limit,
		SortBy: "relevance",
	}

	resp, err := s.esClient.SearchProducts(ctx, req)
	if err != nil {
		return nil, err
	}

	suggestions := make([]Suggestion, 0)
	seen := make(map[string]bool)

	for _, product := range resp.Results {
		// 商品名をキーワードに分割してサジェスト候補を作成
		words := strings.Fields(product.Name)
		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(prefix)) && !seen[word] {
				suggestions = append(suggestions, Suggestion{
					Text:  word,
					Score: product.Rating,
				})
				seen[word] = true

				if len(suggestions) >= limit {
					return suggestions, nil
				}
			}
		}
	}

	return suggestions, nil
}

// PopularKeyword は人気キーワード
type PopularKeyword struct {
	Keyword    string
	SearchCount int64
	TrendScore  float64
}

// GetPopularKeywords は人気キーワードを取得
func (s *SuggestionBuilder) GetPopularKeywords(ctx context.Context, limit int) ([]PopularKeyword, error) {
	// 簡易実装: ハードコードされた人気キーワード
	// 実際は検索ログから集計
	keywords := []PopularKeyword{
		{Keyword: "スマートフォン", SearchCount: 1000, TrendScore: 0.9},
		{Keyword: "ノートパソコン", SearchCount: 800, TrendScore: 0.85},
		{Keyword: "ワイヤレスイヤホン", SearchCount: 750, TrendScore: 0.8},
		{Keyword: "腕時計", SearchCount: 600, TrendScore: 0.75},
		{Keyword: "カメラ", SearchCount: 500, TrendScore: 0.7},
	}

	if limit > len(keywords) {
		limit = len(keywords)
	}

	return keywords[:limit], nil
}
