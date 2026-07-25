# microservices/inventory/proto から、web(加盟店の在庫管理画面)が使う
# RPC のみを手書きで移植した最小スタブ。フィールド番号・型は生成コードを正とすること。
defmodule InventoryService.V1.Inventory do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.Inventory",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:product_id, 2, type: :string, json_name: "productId")
  field(:variation_id, 3, type: :string, json_name: "variationId")
  field(:shop_id, 4, type: :string, json_name: "shopId")
  field(:quantity, 5, type: :int32)
  field(:reserved_quantity, 6, type: :int32, json_name: "reservedQuantity")
  field(:available_quantity, 7, type: :int32, json_name: "availableQuantity")
  field(:alert_threshold, 8, type: :int32, json_name: "alertThreshold")
  field(:last_alerted_at, 9, type: Google.Protobuf.Timestamp, json_name: "lastAlertedAt")
  field(:created_at, 10, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 11, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule InventoryService.V1.RegisterInventoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.RegisterInventoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:variation_id, 2, type: :string, json_name: "variationId")
  field(:shop_id, 3, type: :string, json_name: "shopId")
  field(:initial_quantity, 4, type: :int32, json_name: "initialQuantity")
  field(:alert_threshold, 5, type: :int32, json_name: "alertThreshold")
end

defmodule InventoryService.V1.RegisterInventoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.RegisterInventoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:inventory, 1, type: InventoryService.V1.Inventory)
  field(:message, 2, type: :string)
end

defmodule InventoryService.V1.UpdateInventoryQuantityRequest do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.UpdateInventoryQuantityRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:inventory_id, 1, type: :string, json_name: "inventoryId")
  field(:shop_id, 2, type: :string, json_name: "shopId")
  field(:change_quantity, 3, type: :int32, json_name: "changeQuantity")
  field(:change_type, 4, type: :int32, json_name: "changeType")
  field(:reason, 5, type: :string)
  field(:operator, 6, type: :string)
end

defmodule InventoryService.V1.UpdateInventoryQuantityResponse do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.UpdateInventoryQuantityResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:inventory, 1, type: InventoryService.V1.Inventory)
  field(:message, 3, type: :string)
end

defmodule InventoryService.V1.GetInventoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.GetInventoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:inventory_id, 1, type: :string, json_name: "inventoryId")
end

defmodule InventoryService.V1.GetInventoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.GetInventoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:inventory, 1, type: InventoryService.V1.Inventory)
end

defmodule InventoryService.V1.GetInventoryByProductRequest do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.GetInventoryByProductRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_id, 1, type: :string, json_name: "productId")
  field(:variation_id, 2, type: :string, json_name: "variationId")
end

defmodule InventoryService.V1.GetInventoryByProductResponse do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.GetInventoryByProductResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:inventory, 1, type: InventoryService.V1.Inventory)
end

defmodule InventoryService.V1.BulkGetInventoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.BulkGetInventoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:product_ids, 1, repeated: true, type: :string, json_name: "productIds")
end

defmodule InventoryService.V1.BulkGetInventoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.BulkGetInventoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:inventories, 1, repeated: true, type: InventoryService.V1.Inventory)
end

defmodule InventoryService.V1.CheckStockAlertRequest do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.CheckStockAlertRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:inventory_id, 1, type: :string, json_name: "inventoryId")
end

defmodule InventoryService.V1.CheckStockAlertResponse do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.CheckStockAlertResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:is_low_stock, 1, type: :bool, json_name: "isLowStock")
  field(:current_quantity, 2, type: :int32, json_name: "currentQuantity")
  field(:alert_threshold, 3, type: :int32, json_name: "alertThreshold")
end

defmodule InventoryService.V1.InventoryHistory do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.InventoryHistory",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:inventory_id, 2, type: :string, json_name: "inventoryId")
  field(:change_type, 3, type: :int32, json_name: "changeType")
  field(:change_quantity, 4, type: :int32, json_name: "changeQuantity")
  field(:quantity_before, 5, type: :int32, json_name: "quantityBefore")
  field(:quantity_after, 6, type: :int32, json_name: "quantityAfter")
  field(:reason, 7, type: :string)
  field(:operator, 8, type: :string)
  field(:order_id, 9, type: :string, json_name: "orderId")
  field(:created_at, 10, type: Google.Protobuf.Timestamp, json_name: "createdAt")
end

defmodule InventoryService.V1.GetInventoryHistoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.GetInventoryHistoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:inventory_id, 1, type: :string, json_name: "inventoryId")
  field(:date_from, 2, type: :string, json_name: "dateFrom")
  field(:date_to, 3, type: :string, json_name: "dateTo")
  field(:page, 5, type: :int32)
  field(:page_size, 6, type: :int32, json_name: "pageSize")
end

defmodule InventoryService.V1.GetInventoryHistoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.GetInventoryHistoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:history, 1, repeated: true, type: InventoryService.V1.InventoryHistory)
  field(:total_count, 2, type: :int32, json_name: "totalCount")
  field(:page, 3, type: :int32)
  field(:page_size, 4, type: :int32, json_name: "pageSize")
end

defmodule InventoryService.V1.StockTaking do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.StockTaking",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:inventory_id, 2, type: :string, json_name: "inventoryId")
  field(:shop_id, 3, type: :string, json_name: "shopId")
  field(:system_quantity, 4, type: :int32, json_name: "systemQuantity")
  field(:actual_quantity, 5, type: :int32, json_name: "actualQuantity")
  field(:difference, 6, type: :int32)
  field(:difference_reason, 7, type: :string, json_name: "differenceReason")
  field(:operator, 8, type: :string)
  field(:created_at, 9, type: Google.Protobuf.Timestamp, json_name: "createdAt")
end

defmodule InventoryService.V1.RecordStockTakingRequest do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.RecordStockTakingRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:inventory_id, 1, type: :string, json_name: "inventoryId")
  field(:shop_id, 2, type: :string, json_name: "shopId")
  field(:actual_quantity, 3, type: :int32, json_name: "actualQuantity")
  field(:difference_reason, 4, type: :string, json_name: "differenceReason")
  field(:operator, 5, type: :string)
end

defmodule InventoryService.V1.RecordStockTakingResponse do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.RecordStockTakingResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:stock_taking, 1, type: InventoryService.V1.StockTaking, json_name: "stockTaking")
  field(:message, 2, type: :string)
end

defmodule InventoryService.V1.GetStockTakingHistoryRequest do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.GetStockTakingHistoryRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:shop_id, 1, type: :string, json_name: "shopId")
  field(:date_from, 2, type: :string, json_name: "dateFrom")
  field(:date_to, 3, type: :string, json_name: "dateTo")
end

defmodule InventoryService.V1.GetStockTakingHistoryResponse do
  @moduledoc false

  use Protobuf,
    full_name: "inventory_service.v1.GetStockTakingHistoryResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:history, 1, repeated: true, type: InventoryService.V1.StockTaking)
  field(:total_count, 2, type: :int32, json_name: "totalCount")
end

defmodule InventoryService.V1.InventoryService.Service do
  @moduledoc false

  use GRPC.Service,
    name: "inventory_service.v1.InventoryService",
    protoc_gen_elixir_version: "0.16.0"

  rpc(
    :RegisterInventory,
    InventoryService.V1.RegisterInventoryRequest,
    InventoryService.V1.RegisterInventoryResponse
  )

  rpc(
    :UpdateInventoryQuantity,
    InventoryService.V1.UpdateInventoryQuantityRequest,
    InventoryService.V1.UpdateInventoryQuantityResponse
  )

  rpc(
    :GetInventory,
    InventoryService.V1.GetInventoryRequest,
    InventoryService.V1.GetInventoryResponse
  )

  rpc(
    :GetInventoryByProduct,
    InventoryService.V1.GetInventoryByProductRequest,
    InventoryService.V1.GetInventoryByProductResponse
  )

  rpc(
    :BulkGetInventory,
    InventoryService.V1.BulkGetInventoryRequest,
    InventoryService.V1.BulkGetInventoryResponse
  )

  rpc(
    :CheckStockAlert,
    InventoryService.V1.CheckStockAlertRequest,
    InventoryService.V1.CheckStockAlertResponse
  )

  rpc(
    :GetInventoryHistory,
    InventoryService.V1.GetInventoryHistoryRequest,
    InventoryService.V1.GetInventoryHistoryResponse
  )

  rpc(
    :RecordStockTaking,
    InventoryService.V1.RecordStockTakingRequest,
    InventoryService.V1.RecordStockTakingResponse
  )

  rpc(
    :GetStockTakingHistory,
    InventoryService.V1.GetStockTakingHistoryRequest,
    InventoryService.V1.GetStockTakingHistoryResponse
  )
end

defmodule InventoryService.V1.InventoryService.Stub do
  @moduledoc false

  use GRPC.Stub, service: InventoryService.V1.InventoryService.Service
end
