defmodule ShopService.V1.ShopStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "shop_service.v1.ShopStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:SHOP_STATUS_UNSPECIFIED, 0)
  field(:PENDING_APPROVAL, 1)
  field(:APPROVED, 2)
  field(:REJECTED, 3)
  field(:SUSPENDED, 4)
end

defmodule ShopService.V1.OrderStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "shop_service.v1.OrderStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:ORDER_STATUS_UNSPECIFIED, 0)
  field(:RECEIVED, 1)
  field(:PREPARING, 2)
  field(:SHIPPED, 3)
  field(:DELIVERED, 4)
  field(:CANCELLED, 5)
end

defmodule ShopService.V1.Carrier do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "shop_service.v1.Carrier",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:CARRIER_UNSPECIFIED, 0)
  field(:YAMATO, 1)
  field(:SAGAWA, 2)
  field(:JAPAN_POST, 3)
end

defmodule ShopService.V1.Shop do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.Shop",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:owner_id, 2, type: :string, json_name: "ownerId")
  field(:name, 3, type: :string)
  field(:description, 4, type: :string)
  field(:logo_url, 5, type: :string, json_name: "logoUrl")
  field(:owner_name, 6, type: :string, json_name: "ownerName")
  field(:phone_number, 7, type: :string, json_name: "phoneNumber")
  field(:business_hours, 8, type: :string, json_name: "businessHours")
  field(:return_policy, 9, type: :string, json_name: "returnPolicy")
  field(:status, 10, type: ShopService.V1.ShopStatus, enum: true)
  field(:published, 11, type: :bool)
  field(:created_at, 12, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 13, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule ShopService.V1.ShopCategory do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ShopCategory",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:shop_id, 2, type: :string, json_name: "shopId")
  field(:category_name, 3, type: :string, json_name: "categoryName")
  field(:created_at, 4, type: Google.Protobuf.Timestamp, json_name: "createdAt")
end

defmodule ShopService.V1.Product do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.Product",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:shop_id, 2, type: :string, json_name: "shopId")
  field(:name, 3, type: :string)
  field(:description, 4, type: :string)
  field(:price, 5, type: :string)
  field(:category, 6, type: :string)
  field(:stock_quantity, 7, type: :int32, json_name: "stockQuantity")
  field(:weight, 8, type: :string)
  field(:size, 9, type: :string)
  field(:jan_code, 10, type: :string, json_name: "janCode")
  field(:published, 11, type: :bool)
  field(:deleted, 12, type: :bool)
  field(:created_at, 13, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 14, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule ShopService.V1.ProductImage do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ProductImage",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:product_id, 2, type: :string, json_name: "productId")
  field(:url, 3, type: :string)
  field(:display_order, 4, type: :int32, json_name: "displayOrder")
  field(:thumbnail_200_url, 5, type: :string, json_name: "thumbnail200Url")
  field(:thumbnail_400_url, 6, type: :string, json_name: "thumbnail400Url")
  field(:thumbnail_800_url, 7, type: :string, json_name: "thumbnail800Url")
  field(:created_at, 8, type: Google.Protobuf.Timestamp, json_name: "createdAt")
end

defmodule ShopService.V1.ProductTag do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ProductTag",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:product_id, 2, type: :string, json_name: "productId")
  field(:tag_name, 3, type: :string, json_name: "tagName")
  field(:created_at, 4, type: Google.Protobuf.Timestamp, json_name: "createdAt")
end

defmodule ShopService.V1.ProductVariation do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ProductVariation",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:product_id, 2, type: :string, json_name: "productId")
  field(:sku, 3, type: :string)
  field(:attribute_name, 4, type: :string, json_name: "attributeName")
  field(:attribute_value, 5, type: :string, json_name: "attributeValue")
  field(:price, 6, type: :string)
  field(:stock_quantity, 7, type: :int32, json_name: "stockQuantity")
  field(:created_at, 8, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 9, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule ShopService.V1.Order do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.Order",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:shop_id, 2, type: :string, json_name: "shopId")
  field(:customer_id, 3, type: :string, json_name: "customerId")
  field(:order_number, 4, type: :string, json_name: "orderNumber")
  field(:status, 5, type: ShopService.V1.OrderStatus, enum: true)
  field(:total_amount, 6, type: :string, json_name: "totalAmount")
  field(:shipping_address, 7, type: :string, json_name: "shippingAddress")
  field(:payment_method, 8, type: :string, json_name: "paymentMethod")
  field(:tracking_number, 9, type: :string, json_name: "trackingNumber")
  field(:carrier, 10, type: ShopService.V1.Carrier, enum: true)
  field(:created_at, 11, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 12, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule ShopService.V1.OrderItem do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.OrderItem",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:order_id, 2, type: :string, json_name: "orderId")
  field(:product_id, 3, type: :string, json_name: "productId")
  field(:product_name, 4, type: :string, json_name: "productName")
  field(:quantity, 5, type: :int32)
  field(:unit_price, 6, type: :string, json_name: "unitPrice")
  field(:subtotal, 7, type: :string)
  field(:created_at, 8, type: Google.Protobuf.Timestamp, json_name: "createdAt")
end

defmodule ShopService.V1.SalesReport do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.SalesReport",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:shop_id, 2, type: :string, json_name: "shopId")
  field(:date, 3, type: :string)
  field(:total_sales, 4, type: :string, json_name: "totalSales")
  field(:order_count, 5, type: :int32, json_name: "orderCount")
  field(:average_order_value, 6, type: :string, json_name: "averageOrderValue")
  field(:created_at, 7, type: Google.Protobuf.Timestamp, json_name: "createdAt")
end

defmodule ShopService.V1.RegisterShopRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.RegisterShopRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:owner_id, 1, type: :string, json_name: "ownerId")
  field(:name, 2, type: :string)
  field(:description, 3, type: :string)
  field(:logo_url, 4, type: :string, json_name: "logoUrl")
  field(:owner_name, 5, type: :string, json_name: "ownerName")
  field(:phone_number, 6, type: :string, json_name: "phoneNumber")
  field(:business_hours, 7, type: :string, json_name: "businessHours")
  field(:return_policy, 8, type: :string, json_name: "returnPolicy")
  field(:categories, 9, repeated: true, type: :string)
end

defmodule ShopService.V1.UpdateShopRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.UpdateShopRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:name, 2, type: :string)
  field(:description, 3, type: :string)
  field(:logo_url, 4, type: :string, json_name: "logoUrl")
  field(:business_hours, 5, type: :string, json_name: "businessHours")
  field(:return_policy, 6, type: :string, json_name: "returnPolicy")
end

defmodule ShopService.V1.GetShopRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.GetShopRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
end

defmodule ShopService.V1.ToggleShopPublishRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ToggleShopPublishRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:published, 2, type: :bool)
end

defmodule ShopService.V1.RegisterProductRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.RegisterProductRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:name, 2, type: :string)
  field(:description, 3, type: :string)
  field(:price, 4, type: :string)
  field(:category, 5, type: :string)
  field(:stock_quantity, 6, type: :int32, json_name: "stockQuantity")
  field(:weight, 7, type: :string)
  field(:size, 8, type: :string)
  field(:jan_code, 9, type: :string, json_name: "janCode")
  field(:tags, 10, repeated: true, type: :string)
end

defmodule ShopService.V1.UpdateProductRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.UpdateProductRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:name, 2, type: :string)
  field(:description, 3, type: :string)
  field(:price, 4, type: :string)
  field(:category, 5, type: :string)
  field(:stock_quantity, 6, type: :int32, json_name: "stockQuantity")
  field(:weight, 7, type: :string)
  field(:size, 8, type: :string)
  field(:jan_code, 9, type: :string, json_name: "janCode")
end

defmodule ShopService.V1.DeleteProductRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.DeleteProductRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
end

defmodule ShopService.V1.ToggleProductPublishRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ToggleProductPublishRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:published, 2, type: :bool)
end

defmodule ShopService.V1.GetProductRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.GetProductRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
end

defmodule ShopService.V1.ListProductsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ListProductsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:category, 2, type: :string)
  field(:published_only, 3, type: :bool, json_name: "publishedOnly")
  field(:limit, 4, type: :int32)
  field(:offset, 5, type: :int32)
end

defmodule ShopService.V1.UploadProductImageRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.UploadProductImageRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:image_data, 2, type: :bytes, json_name: "imageData")
  field(:display_order, 3, type: :int32, json_name: "displayOrder")
end

defmodule ShopService.V1.ManageVariationRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ManageVariationRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:variations, 2, repeated: true, type: ShopService.V1.ProductVariationInput)
end

defmodule ShopService.V1.ListOrdersRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ListOrdersRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:status, 2, type: ShopService.V1.OrderStatus, enum: true)
  field(:date_from, 3, type: :string, json_name: "dateFrom")
  field(:date_to, 4, type: :string, json_name: "dateTo")
  field(:customer_name, 5, type: :string, json_name: "customerName")
  field(:product_name, 6, type: :string, json_name: "productName")
  field(:sort_by, 7, type: :string, json_name: "sortBy")
  field(:sort_order, 8, type: :string, json_name: "sortOrder")
end

defmodule ShopService.V1.GetOrderDetailRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.GetOrderDetailRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order_id, 1, type: :string, json_name: "orderId")
  field(:shop_id, 2, type: :string, json_name: "shopId")
end

defmodule ShopService.V1.UpdateOrderStatusRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.UpdateOrderStatusRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order_id, 1, type: :string, json_name: "orderId")
  field(:shop_id, 2, type: :string, json_name: "shopId")
  field(:new_status, 3, type: ShopService.V1.OrderStatus, json_name: "newStatus", enum: true)
  field(:tracking_number, 4, type: :string, json_name: "trackingNumber")
  field(:carrier, 5, type: ShopService.V1.Carrier, enum: true)
end

defmodule ShopService.V1.GetSalesReportRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.GetSalesReportRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:report_type, 2, type: :string, json_name: "reportType")
  field(:date_from, 3, type: :string, json_name: "dateFrom")
  field(:date_to, 4, type: :string, json_name: "dateTo")
end

defmodule ShopService.V1.ExportSalesDataRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ExportSalesDataRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:date_from, 2, type: :string, json_name: "dateFrom")
  field(:date_to, 3, type: :string, json_name: "dateTo")
end

defmodule ShopService.V1.ProductVariationInput do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ProductVariationInput",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:sku, 1, type: :string)
  field(:attribute_name, 2, type: :string, json_name: "attributeName")
  field(:attribute_value, 3, type: :string, json_name: "attributeValue")
  field(:price, 4, type: :string)
  field(:stock_quantity, 5, type: :int32, json_name: "stockQuantity")
end

defmodule ShopService.V1.OrderSummary do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.OrderSummary",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:order_number, 2, type: :string, json_name: "orderNumber")
  field(:status, 3, type: ShopService.V1.OrderStatus, enum: true)
  field(:total_amount, 4, type: :string, json_name: "totalAmount")
  field(:created_at, 5, type: Google.Protobuf.Timestamp, json_name: "createdAt")
end

defmodule ShopService.V1.OrderDetail do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.OrderDetail",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order, 1, type: ShopService.V1.Order)
  field(:items, 2, repeated: true, type: ShopService.V1.OrderItem)
end

defmodule ShopService.V1.SalesData do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.SalesData",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:date, 1, type: :string)
  field(:total_sales, 2, type: :string, json_name: "totalSales")
  field(:order_count, 3, type: :int32, json_name: "orderCount")
  field(:average_order_value, 4, type: :string, json_name: "averageOrderValue")
end

defmodule ShopService.V1.SalesSummary do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.SalesSummary",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:total_sales, 1, type: :string, json_name: "totalSales")
  field(:total_orders, 2, type: :int32, json_name: "totalOrders")
  field(:average_order_value, 3, type: :string, json_name: "averageOrderValue")
end

defmodule ShopService.V1.RegisterShopResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.RegisterShopResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:status, 2, type: ShopService.V1.ShopStatus, enum: true)
  field(:message, 3, type: :string)
end

defmodule ShopService.V1.UpdateShopResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.UpdateShopResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:requires_reapproval, 2, type: :bool, json_name: "requiresReapproval")
end

defmodule ShopService.V1.ToggleShopPublishResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ToggleShopPublishResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:published, 2, type: :bool)
end

defmodule ShopService.V1.GetShopResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.GetShopResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop, 1, type: ShopService.V1.Shop)
end

defmodule ShopService.V1.ListShopsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ListShopsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:published_only, 1, type: :bool, json_name: "publishedOnly")
  field(:limit, 2, type: :int32)
  field(:offset, 3, type: :int32)
end

defmodule ShopService.V1.ListShopsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ListShopsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shops, 1, repeated: true, type: ShopService.V1.Shop)
  field(:total_count, 2, type: :int32, json_name: "totalCount")
end

defmodule ShopService.V1.GetShopsByOwnerRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.GetShopsByOwnerRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:owner_id, 1, type: :string, json_name: "ownerId")
end

defmodule ShopService.V1.RegisterProductResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.RegisterProductResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
end

defmodule ShopService.V1.UpdateProductResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.UpdateProductResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
end

defmodule ShopService.V1.DeleteProductResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.DeleteProductResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:deleted, 2, type: :bool)
end

defmodule ShopService.V1.ToggleProductPublishResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ToggleProductPublishResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:published, 2, type: :bool)
end

defmodule ShopService.V1.GetProductResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.GetProductResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product, 1, type: ShopService.V1.Product)
end

defmodule ShopService.V1.ListProductsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ListProductsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:products, 1, repeated: true, type: ShopService.V1.Product)
  field(:total_count, 2, type: :int32, json_name: "totalCount")
end

defmodule ShopService.V1.UploadProductImageResponse.ThumbnailsEntry do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.UploadProductImageResponse.ThumbnailsEntry",
    map: true,
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:key, 1, type: :string)
  field(:value, 2, type: :string)
end

defmodule ShopService.V1.UploadProductImageResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.UploadProductImageResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:image_id, 1, type: :string, json_name: "imageId")
  field(:url, 2, type: :string)

  field(:thumbnails, 3,
    repeated: true,
    type: ShopService.V1.UploadProductImageResponse.ThumbnailsEntry,
    map: true
  )
end

defmodule ShopService.V1.ManageVariationResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ManageVariationResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:variation_ids, 1, repeated: true, type: :string, json_name: "variationIds")
end

defmodule ShopService.V1.ListOrdersResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ListOrdersResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:orders, 1, repeated: true, type: ShopService.V1.OrderSummary)
  field(:total_count, 2, type: :int32, json_name: "totalCount")
end

defmodule ShopService.V1.GetOrderDetailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.GetOrderDetailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order, 1, type: ShopService.V1.OrderDetail)
end

defmodule ShopService.V1.UpdateOrderStatusResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.UpdateOrderStatusResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order_id, 1, type: :string, json_name: "orderId")
  field(:status, 2, type: ShopService.V1.OrderStatus, enum: true)
end

defmodule ShopService.V1.GetSalesReportResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.GetSalesReportResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:report_data, 1, repeated: true, type: ShopService.V1.SalesData, json_name: "reportData")
  field(:summary, 2, type: ShopService.V1.SalesSummary)
end

defmodule ShopService.V1.ExportSalesDataResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shop_service.v1.ExportSalesDataResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:csv_url, 1, type: :string, json_name: "csvUrl")
  field(:expires_at, 2, type: Google.Protobuf.Timestamp, json_name: "expiresAt")
end

defmodule ShopService.V1.ShopService.Service do
  @moduledoc false

  use GRPC.Service, name: "shop_service.v1.ShopService", protoc_gen_elixir_version: "0.16.0"

  rpc(:RegisterShop, ShopService.V1.RegisterShopRequest, ShopService.V1.RegisterShopResponse)

  rpc(:UpdateShop, ShopService.V1.UpdateShopRequest, ShopService.V1.UpdateShopResponse)

  rpc(
    :ToggleShopPublish,
    ShopService.V1.ToggleShopPublishRequest,
    ShopService.V1.ToggleShopPublishResponse
  )

  rpc(:GetShop, ShopService.V1.GetShopRequest, ShopService.V1.GetShopResponse)

  rpc(:ListShops, ShopService.V1.ListShopsRequest, ShopService.V1.ListShopsResponse)

  rpc(:GetShopsByOwner, ShopService.V1.GetShopsByOwnerRequest, ShopService.V1.ListShopsResponse)

  rpc(
    :RegisterProduct,
    ShopService.V1.RegisterProductRequest,
    ShopService.V1.RegisterProductResponse
  )

  rpc(:UpdateProduct, ShopService.V1.UpdateProductRequest, ShopService.V1.UpdateProductResponse)

  rpc(:DeleteProduct, ShopService.V1.DeleteProductRequest, ShopService.V1.DeleteProductResponse)

  rpc(
    :ToggleProductPublish,
    ShopService.V1.ToggleProductPublishRequest,
    ShopService.V1.ToggleProductPublishResponse
  )

  rpc(:GetProduct, ShopService.V1.GetProductRequest, ShopService.V1.GetProductResponse)

  rpc(:ListProducts, ShopService.V1.ListProductsRequest, ShopService.V1.ListProductsResponse)

  rpc(
    :UploadProductImage,
    ShopService.V1.UploadProductImageRequest,
    ShopService.V1.UploadProductImageResponse
  )

  rpc(
    :ManageVariation,
    ShopService.V1.ManageVariationRequest,
    ShopService.V1.ManageVariationResponse
  )

  rpc(:ListOrders, ShopService.V1.ListOrdersRequest, ShopService.V1.ListOrdersResponse)

  rpc(
    :GetOrderDetail,
    ShopService.V1.GetOrderDetailRequest,
    ShopService.V1.GetOrderDetailResponse
  )

  rpc(
    :UpdateOrderStatus,
    ShopService.V1.UpdateOrderStatusRequest,
    ShopService.V1.UpdateOrderStatusResponse
  )

  rpc(
    :GetSalesReport,
    ShopService.V1.GetSalesReportRequest,
    ShopService.V1.GetSalesReportResponse
  )

  rpc(
    :ExportSalesData,
    ShopService.V1.ExportSalesDataRequest,
    ShopService.V1.ExportSalesDataResponse
  )
end

defmodule ShopService.V1.ShopService.Stub do
  @moduledoc false

  use GRPC.Stub, service: ShopService.V1.ShopService.Service
end
