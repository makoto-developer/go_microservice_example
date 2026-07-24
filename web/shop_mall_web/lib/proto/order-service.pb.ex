# microservices/order/proto/order-service.proto から、web が使う CreateOrder 関連のみを
# 手書きで移植した最小スタブ。フィールド番号・型は order-service.proto を正とすること。
defmodule OrderService.V1.PaymentMethod do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "order_service.v1.PaymentMethod",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:PAYMENT_METHOD_UNSPECIFIED, 0)
  field(:CREDIT_CARD, 1)
  field(:CASH_ON_DELIVERY, 2)
end

defmodule OrderService.V1.OrderStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "order_service.v1.OrderStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:ORDER_STATUS_UNSPECIFIED, 0)
  field(:PENDING, 1)
  field(:PAYMENT_PROCESSING, 2)
  field(:PAYMENT_FAILED, 3)
  field(:CONFIRMED, 4)
  field(:PREPARING, 5)
  field(:SHIPPED, 6)
  field(:DELIVERED, 7)
  field(:CANCELLED, 8)
end

defmodule OrderService.V1.CancelReason do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "order_service.v1.CancelReason",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:CANCEL_REASON_UNSPECIFIED, 0)
  field(:CUSTOMER_NO_LONGER_NEEDED, 1)
  field(:CUSTOMER_ORDERED_BY_MISTAKE, 2)
  field(:CUSTOMER_DELIVERY_TIME_ISSUE, 3)
  field(:CUSTOMER_OTHER, 4)
  field(:SHOP_OUT_OF_STOCK, 5)
  field(:SHOP_DEFECTIVE_PRODUCT, 6)
  field(:SHOP_OTHER, 7)
end

defmodule OrderService.V1.Order do
  @moduledoc false

  use Protobuf,
    full_name: "order_service.v1.Order",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:order_number, 2, type: :string, json_name: "orderNumber")
  field(:customer_id, 3, type: :string, json_name: "customerId")
  field(:status, 5, type: OrderService.V1.OrderStatus, enum: true)
  field(:total_amount, 6, type: :string, json_name: "totalAmount")
  field(:shipping_fee, 8, type: :string, json_name: "shippingFee")

  field(:payment_method, 9,
    type: OrderService.V1.PaymentMethod,
    enum: true,
    json_name: "paymentMethod"
  )

  field(:created_at, 18, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 19, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule OrderService.V1.ListOrdersRequest do
  @moduledoc false

  use Protobuf,
    full_name: "order_service.v1.ListOrdersRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:page, 9, type: :int32)
  field(:page_size, 10, type: :int32, json_name: "pageSize")
end

defmodule OrderService.V1.ListOrdersResponse do
  @moduledoc false

  use Protobuf,
    full_name: "order_service.v1.ListOrdersResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:orders, 1, repeated: true, type: OrderService.V1.Order)
  field(:total_count, 2, type: :int32, json_name: "totalCount")
  field(:page, 3, type: :int32)
  field(:page_size, 4, type: :int32, json_name: "pageSize")
end

defmodule OrderService.V1.CancelOrderRequest do
  @moduledoc false

  use Protobuf,
    full_name: "order_service.v1.CancelOrderRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order_id, 1, type: :string, json_name: "orderId")
  field(:cancelled_by, 2, type: :string, json_name: "cancelledBy")

  field(:cancel_reason, 3,
    type: OrderService.V1.CancelReason,
    enum: true,
    json_name: "cancelReason"
  )

  field(:cancel_note, 4, type: :string, json_name: "cancelNote")
end

defmodule OrderService.V1.CancelOrderResponse do
  @moduledoc false

  use Protobuf,
    full_name: "order_service.v1.CancelOrderResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule OrderService.V1.CartItemInput do
  @moduledoc false

  use Protobuf,
    full_name: "order_service.v1.CartItemInput",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:variation_id, 2, type: :string, json_name: "variationId")
  field(:quantity, 3, type: :int32)
  field(:unit_price, 4, type: :int64, json_name: "unitPrice")
end

defmodule OrderService.V1.CreateOrderRequest do
  @moduledoc false

  use Protobuf,
    full_name: "order_service.v1.CreateOrderRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:customer_id, 1, type: :string, json_name: "customerId")
  field(:customer_email, 2, type: :string, json_name: "customerEmail")

  field(:cart_items, 3,
    repeated: true,
    type: OrderService.V1.CartItemInput,
    json_name: "cartItems"
  )

  field(:shipping_address_id, 4, type: :string, json_name: "shippingAddressId")

  field(:payment_method, 5,
    type: OrderService.V1.PaymentMethod,
    enum: true,
    json_name: "paymentMethod"
  )

  field(:payment_method_id, 6, type: :string, json_name: "paymentMethodId")
  field(:shipping_method, 7, type: :string, json_name: "shippingMethod")
  field(:notes, 8, type: :string)
end

defmodule OrderService.V1.CreateOrderResponse do
  @moduledoc false

  use Protobuf,
    full_name: "order_service.v1.CreateOrderResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order_id, 1, type: :string, json_name: "orderId")
  field(:order_number, 2, type: :string, json_name: "orderNumber")
  field(:message, 3, type: :string)
end

defmodule OrderService.V1.OrderService.Service do
  @moduledoc false

  use GRPC.Service, name: "order_service.v1.OrderService", protoc_gen_elixir_version: "0.16.0"

  rpc(:CreateOrder, OrderService.V1.CreateOrderRequest, OrderService.V1.CreateOrderResponse)
  rpc(:ListOrders, OrderService.V1.ListOrdersRequest, OrderService.V1.ListOrdersResponse)
  rpc(:CancelOrder, OrderService.V1.CancelOrderRequest, OrderService.V1.CancelOrderResponse)
end

defmodule OrderService.V1.OrderService.Stub do
  @moduledoc false

  use GRPC.Stub, service: OrderService.V1.OrderService.Service
end
