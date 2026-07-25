# microservices/shipping/proto/shipping-service.proto から、web(加盟店画面)が使う
# RPC のみを手書きで移植した最小スタブ。フィールド番号・型は shipping-service.proto を正とすること。
defmodule ShippingService.V1.ShipmentStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "shipping_service.v1.ShipmentStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:SHIPMENT_STATUS_UNSPECIFIED, 0)
  field(:SHIPMENT_STATUS_PENDING, 1)
  field(:SHIPMENT_STATUS_PREPARING, 2)
  field(:SHIPMENT_STATUS_SHIPPED, 3)
  field(:SHIPMENT_STATUS_IN_TRANSIT, 4)
  field(:SHIPMENT_STATUS_DELIVERED, 5)
  field(:SHIPMENT_STATUS_FAILED, 6)
end

defmodule ShippingService.V1.Shipment do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.Shipment",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:order_id, 2, type: :string, json_name: "orderId")
  field(:customer_id, 3, type: :string, json_name: "customerId")
  field(:status, 4, type: ShippingService.V1.ShipmentStatus, enum: true)
  field(:tracking_number, 5, type: :string, json_name: "trackingNumber")
  field(:carrier, 6, type: :string)
  field(:shipping_address, 7, type: :string, json_name: "shippingAddress")

  field(:estimated_delivery, 8,
    type: Google.Protobuf.Timestamp,
    json_name: "estimatedDelivery"
  )

  field(:actual_delivery, 9, type: Google.Protobuf.Timestamp, json_name: "actualDelivery")
  field(:created_at, 10, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 11, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule ShippingService.V1.RegisterTrackingNumberRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.RegisterTrackingNumberRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shipment_id, 1, type: :string, json_name: "shipmentId")
  field(:tracking_number, 2, type: :string, json_name: "trackingNumber")
end

defmodule ShippingService.V1.RegisterTrackingNumberResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.RegisterTrackingNumberResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ShippingService.V1.UpdateShipmentStatusRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.UpdateShipmentStatusRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shipment_id, 1, type: :string, json_name: "shipmentId")

  field(:new_status, 2,
    type: ShippingService.V1.ShipmentStatus,
    enum: true,
    json_name: "newStatus"
  )
end

defmodule ShippingService.V1.UpdateShipmentStatusResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.UpdateShipmentStatusResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ShippingService.V1.GetShipmentByOrderRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.GetShipmentByOrderRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order_id, 1, type: :string, json_name: "orderId")
end

defmodule ShippingService.V1.GetShipmentByOrderResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.GetShipmentByOrderResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
  field(:shipment, 3, type: ShippingService.V1.Shipment)
end

defmodule ShippingService.V1.GetShipmentDetailRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.GetShipmentDetailRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shipment_id, 1, type: :string, json_name: "shipmentId")
end

defmodule ShippingService.V1.GetShipmentDetailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.GetShipmentDetailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
  field(:shipment, 3, type: ShippingService.V1.Shipment)
end

defmodule ShippingService.V1.ValidateAddressRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.ValidateAddressRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:postal_code, 1, type: :string, json_name: "postalCode")
  field(:prefecture, 2, type: :string)
  field(:city, 3, type: :string)
  field(:address_line, 4, type: :string, json_name: "addressLine")
end

defmodule ShippingService.V1.ValidateAddressResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.ValidateAddressResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule ShippingService.V1.NormalizeAddressRequest do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.NormalizeAddressRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:postal_code, 1, type: :string, json_name: "postalCode")
  field(:prefecture, 2, type: :string)
  field(:city, 3, type: :string)
  field(:address_line, 4, type: :string, json_name: "addressLine")
end

defmodule ShippingService.V1.NormalizeAddressResponse do
  @moduledoc false

  use Protobuf,
    full_name: "shipping_service.v1.NormalizeAddressResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
  field(:normalized_address, 3, type: :string, json_name: "normalizedAddress")
end

defmodule ShippingService.V1.ShippingService.Service do
  @moduledoc false

  use GRPC.Service,
    name: "shipping_service.v1.ShippingService",
    protoc_gen_elixir_version: "0.16.0"

  rpc(
    :RegisterTrackingNumber,
    ShippingService.V1.RegisterTrackingNumberRequest,
    ShippingService.V1.RegisterTrackingNumberResponse
  )

  rpc(
    :UpdateShipmentStatus,
    ShippingService.V1.UpdateShipmentStatusRequest,
    ShippingService.V1.UpdateShipmentStatusResponse
  )

  rpc(
    :GetShipmentByOrder,
    ShippingService.V1.GetShipmentByOrderRequest,
    ShippingService.V1.GetShipmentByOrderResponse
  )

  rpc(
    :GetShipmentDetail,
    ShippingService.V1.GetShipmentDetailRequest,
    ShippingService.V1.GetShipmentDetailResponse
  )

  rpc(
    :ValidateAddress,
    ShippingService.V1.ValidateAddressRequest,
    ShippingService.V1.ValidateAddressResponse
  )

  rpc(
    :NormalizeAddress,
    ShippingService.V1.NormalizeAddressRequest,
    ShippingService.V1.NormalizeAddressResponse
  )
end

defmodule ShippingService.V1.ShippingService.Stub do
  @moduledoc false

  use GRPC.Stub, service: ShippingService.V1.ShippingService.Service
end
