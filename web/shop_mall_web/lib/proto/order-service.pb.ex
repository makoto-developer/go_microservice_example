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
end

defmodule OrderService.V1.OrderService.Stub do
  @moduledoc false

  use GRPC.Stub, service: OrderService.V1.OrderService.Service
end
