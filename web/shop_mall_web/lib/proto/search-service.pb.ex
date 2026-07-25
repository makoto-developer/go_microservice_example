# microservices/search/proto から pb2ex.py で自動生成した最小スタブ。
defmodule SearchService.V1.PeriodType do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "search_service.v1.PeriodType",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:PERIOD_TYPE_UNSPECIFIED, 0)
  field(:HOURLY, 1)
  field(:DAILY, 2)
  field(:WEEKLY, 3)
  field(:MONTHLY, 4)
end

defmodule SearchService.V1.SortBy do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "search_service.v1.SortBy",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:SORT_BY_UNSPECIFIED, 0)
  field(:RELEVANCE, 1)
  field(:PRICE_ASC, 2)
  field(:PRICE_DESC, 3)
  field(:RATING_DESC, 4)
  field(:NEWEST, 5)
  field(:REVIEW_COUNT_DESC, 6)
end

defmodule SearchService.V1.StockStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "search_service.v1.StockStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:STOCK_STATUS_UNSPECIFIED, 0)
  field(:IN_STOCK, 1)
  field(:LOW_STOCK, 2)
  field(:OUT_OF_STOCK, 3)
end

defmodule SearchService.V1.ClearAllSearchHistoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.ClearAllSearchHistoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
end

defmodule SearchService.V1.ClearAllSearchHistoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.ClearAllSearchHistoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.DeleteProductIndexRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.DeleteProductIndexRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
end

defmodule SearchService.V1.DeleteProductIndexResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.DeleteProductIndexResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.DeleteSearchHistoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.DeleteSearchHistoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:history_id, 2, type: :string, json_name: "historyId")
end

defmodule SearchService.V1.DeleteSearchHistoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.DeleteSearchHistoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.DeleteShopIndexRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.DeleteShopIndexRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
end

defmodule SearchService.V1.DeleteShopIndexResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.DeleteShopIndexResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.GetPopularKeywordsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.GetPopularKeywordsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:period_type, 1, type: SearchService.V1.PeriodType, enum: true, json_name: "periodType")
  field(:category, 2, type: :string)
  field(:limit, 3, type: :int32)
end

defmodule SearchService.V1.GetPopularKeywordsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.GetPopularKeywordsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.GetSearchAnalyticsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.GetSearchAnalyticsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:date_from, 1, type: :string, json_name: "dateFrom")
  field(:date_to, 2, type: :string, json_name: "dateTo")
  field(:report_type, 3, type: :string, json_name: "reportType")
end

defmodule SearchService.V1.GetSearchAnalyticsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.GetSearchAnalyticsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.GetSearchHistoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.GetSearchHistoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:limit, 2, type: :int32)
end

defmodule SearchService.V1.GetSearchHistoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.GetSearchHistoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.GetSearchSuggestionsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.GetSearchSuggestionsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:prefix, 1, type: :string)
  field(:limit, 2, type: :int32)
end

defmodule SearchService.V1.GetSearchSuggestionsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.GetSearchSuggestionsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.GetTrendingKeywordsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.GetTrendingKeywordsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:limit, 1, type: :int32)
end

defmodule SearchService.V1.GetTrendingKeywordsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.GetTrendingKeywordsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.IndexProductRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.IndexProductRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:product_name, 2, type: :string, json_name: "productName")
  field(:description, 3, type: :string)
  field(:category, 4, type: :string)
  field(:tags, 5, repeated: true, type: :string)
  field(:shop_id, 6, type: :string, json_name: "shopId")
  field(:shop_name, 7, type: :string, json_name: "shopName")
  field(:price, 8, type: :string)
  field(:average_rating, 9, type: :string, json_name: "averageRating")
  field(:review_count, 10, type: :int32, json_name: "reviewCount")

  field(:stock_status, 11,
    type: SearchService.V1.StockStatus,
    enum: true,
    json_name: "stockStatus"
  )

  field(:image_url, 12, type: :string, json_name: "imageUrl")
end

defmodule SearchService.V1.IndexProductResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.IndexProductResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.IndexShopRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.IndexShopRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:shop_name, 2, type: :string, json_name: "shopName")
  field(:description, 3, type: :string)
  field(:categories, 4, repeated: true, type: :string)
  field(:average_rating, 5, type: :string, json_name: "averageRating")
  field(:product_count, 6, type: :int32, json_name: "productCount")
  field(:logo_url, 7, type: :string, json_name: "logoUrl")
end

defmodule SearchService.V1.IndexShopResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.IndexShopResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.RecordSearchHistoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.RecordSearchHistoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:user_id, 1, type: :string, json_name: "userId")
  field(:keyword, 2, type: :string)
  field(:result_count, 3, type: :int32, json_name: "resultCount")
  field(:clicked_product_id, 4, type: :string, json_name: "clickedProductId")
end

defmodule SearchService.V1.RecordSearchHistoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.RecordSearchHistoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.SearchProductsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.SearchProductsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:keyword, 1, type: :string)
  field(:category, 2, type: :string)
  field(:min_price, 3, type: :string, json_name: "minPrice")
  field(:max_price, 4, type: :string, json_name: "maxPrice")
  field(:min_rating, 5, type: :string, json_name: "minRating")

  field(:stock_status, 6,
    type: SearchService.V1.StockStatus,
    enum: true,
    json_name: "stockStatus"
  )

  field(:shop_id, 7, type: :string, json_name: "shopId")
  field(:sort_by, 8, type: SearchService.V1.SortBy, enum: true, json_name: "sortBy")
  field(:page, 9, type: :int32)
  field(:page_size, 10, type: :int32, json_name: "pageSize")
end

defmodule SearchService.V1.SearchProductsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.SearchProductsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.SearchShopsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.SearchShopsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:keyword, 1, type: :string)
  field(:category, 2, type: :string)
  field(:min_rating, 3, type: :string, json_name: "minRating")
  field(:sort_by, 4, type: :string, json_name: "sortBy")
  field(:page, 5, type: :int32)
  field(:page_size, 6, type: :int32, json_name: "pageSize")
end

defmodule SearchService.V1.SearchShopsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.SearchShopsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.UpdateProductIndexRequest do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.UpdateProductIndexRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:product_name, 2, type: :string, json_name: "productName")
  field(:description, 3, type: :string)
  field(:price, 4, type: :string)
  field(:average_rating, 5, type: :string, json_name: "averageRating")
  field(:review_count, 6, type: :int32, json_name: "reviewCount")

  field(:stock_status, 7,
    type: SearchService.V1.StockStatus,
    enum: true,
    json_name: "stockStatus"
  )
end

defmodule SearchService.V1.UpdateProductIndexResponse do
  @moduledoc false

  use Protobuf,
    full_name: "search_service.v1.UpdateProductIndexResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule SearchService.V1.SearchService.Service do
  @moduledoc false

  use GRPC.Service,
    name: "search_service.v1.SearchService",
    protoc_gen_elixir_version: "0.16.0"

  rpc(
    :SearchProducts,
    SearchService.V1.SearchProductsRequest,
    SearchService.V1.SearchProductsResponse
  )

  rpc(
    :GetSearchSuggestions,
    SearchService.V1.GetSearchSuggestionsRequest,
    SearchService.V1.GetSearchSuggestionsResponse
  )

  rpc(:SearchShops, SearchService.V1.SearchShopsRequest, SearchService.V1.SearchShopsResponse)

  rpc(
    :RecordSearchHistory,
    SearchService.V1.RecordSearchHistoryRequest,
    SearchService.V1.RecordSearchHistoryResponse
  )

  rpc(
    :GetSearchHistory,
    SearchService.V1.GetSearchHistoryRequest,
    SearchService.V1.GetSearchHistoryResponse
  )

  rpc(
    :DeleteSearchHistory,
    SearchService.V1.DeleteSearchHistoryRequest,
    SearchService.V1.DeleteSearchHistoryResponse
  )

  rpc(
    :ClearAllSearchHistory,
    SearchService.V1.ClearAllSearchHistoryRequest,
    SearchService.V1.ClearAllSearchHistoryResponse
  )

  rpc(
    :GetPopularKeywords,
    SearchService.V1.GetPopularKeywordsRequest,
    SearchService.V1.GetPopularKeywordsResponse
  )

  rpc(
    :GetTrendingKeywords,
    SearchService.V1.GetTrendingKeywordsRequest,
    SearchService.V1.GetTrendingKeywordsResponse
  )

  rpc(:IndexProduct, SearchService.V1.IndexProductRequest, SearchService.V1.IndexProductResponse)

  rpc(
    :UpdateProductIndex,
    SearchService.V1.UpdateProductIndexRequest,
    SearchService.V1.UpdateProductIndexResponse
  )

  rpc(
    :DeleteProductIndex,
    SearchService.V1.DeleteProductIndexRequest,
    SearchService.V1.DeleteProductIndexResponse
  )

  rpc(:IndexShop, SearchService.V1.IndexShopRequest, SearchService.V1.IndexShopResponse)

  rpc(
    :DeleteShopIndex,
    SearchService.V1.DeleteShopIndexRequest,
    SearchService.V1.DeleteShopIndexResponse
  )

  rpc(
    :GetSearchAnalytics,
    SearchService.V1.GetSearchAnalyticsRequest,
    SearchService.V1.GetSearchAnalyticsResponse
  )
end

defmodule SearchService.V1.SearchService.Stub do
  @moduledoc false

  use GRPC.Stub, service: SearchService.V1.SearchService.Service
end
