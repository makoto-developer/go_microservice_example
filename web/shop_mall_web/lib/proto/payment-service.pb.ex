# microservices/payment/proto/payment-service.proto から、web(管理者・加盟店画面)が使う
# RPC のみを手書きで移植した最小スタブ。フィールド番号・型は payment-service.proto を正とすること。
defmodule PaymentService.V1.PaymentMethodType do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "payment_service.v1.PaymentMethodType",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:PAYMENT_METHOD_TYPE_UNSPECIFIED, 0)
  field(:CREDIT_CARD, 1)
  field(:CASH_ON_DELIVERY, 2)
end

defmodule PaymentService.V1.PaymentStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "payment_service.v1.PaymentStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:PAYMENT_STATUS_PAYMENT_STATUS_UNSPECIFIED, 0)
  field(:PAYMENT_STATUS_PENDING, 1)
  field(:PAYMENT_STATUS_PROCESSING, 2)
  field(:PAYMENT_STATUS_REQUIRES_AUTHENTICATION, 3)
  field(:PAYMENT_STATUS_SUCCEEDED, 4)
  field(:PAYMENT_STATUS_FAILED, 5)
  field(:PAYMENT_STATUS_REFUNDING, 6)
  field(:PAYMENT_STATUS_REFUNDED, 7)
end

defmodule PaymentService.V1.RefundStatus do
  @moduledoc false

  use Protobuf,
    enum: true,
    full_name: "payment_service.v1.RefundStatus",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:REFUND_STATUS_REFUND_STATUS_UNSPECIFIED, 0)
  field(:REFUND_STATUS_PENDING, 1)
  field(:REFUND_STATUS_PROCESSING, 2)
  field(:REFUND_STATUS_SUCCEEDED, 3)
  field(:REFUND_STATUS_FAILED, 4)
end

defmodule PaymentService.V1.Payment do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.Payment",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:id, 1, type: :string)
  field(:order_id, 2, type: :string, json_name: "orderId")

  field(:payment_method, 3,
    type: PaymentService.V1.PaymentMethodType,
    enum: true,
    json_name: "paymentMethod"
  )

  field(:amount, 4, type: :string)
  field(:currency, 5, type: :string)
  field(:status, 6, type: PaymentService.V1.PaymentStatus, enum: true)
  field(:cod_fee, 13, type: :string, json_name: "codFee")
  field(:created_at, 14, type: Google.Protobuf.Timestamp, json_name: "createdAt")
  field(:updated_at, 15, type: Google.Protobuf.Timestamp, json_name: "updatedAt")
end

defmodule PaymentService.V1.GetPaymentStatusRequest do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.GetPaymentStatusRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:payment_id, 1, type: :string, json_name: "paymentId")
end

defmodule PaymentService.V1.GetPaymentStatusResponse do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.GetPaymentStatusResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:status, 1, type: PaymentService.V1.PaymentStatus, enum: true)
  field(:transaction_id, 2, type: :string, json_name: "transactionId")
end

defmodule PaymentService.V1.ConfirmCODPaymentRequest do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.ConfirmCODPaymentRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:payment_id, 1, type: :string, json_name: "paymentId")
  field(:order_id, 2, type: :string, json_name: "orderId")
end

defmodule PaymentService.V1.ConfirmCODPaymentResponse do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.ConfirmCODPaymentResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:success, 1, type: :bool)
  field(:message, 2, type: :string)
end

defmodule PaymentService.V1.CreateRefundRequest do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.CreateRefundRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:payment_id, 1, type: :string, json_name: "paymentId")
  field(:order_id, 2, type: :string, json_name: "orderId")
  field(:amount, 3, type: :string)
  field(:reason, 4, type: :string)
end

defmodule PaymentService.V1.CreateRefundResponse do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.CreateRefundResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:refund_id, 1, type: :string, json_name: "refundId")
  field(:message, 2, type: :string)
end

defmodule PaymentService.V1.GetRefundStatusRequest do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.GetRefundStatusRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:refund_id, 1, type: :string, json_name: "refundId")
end

defmodule PaymentService.V1.GetRefundStatusResponse do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.GetRefundStatusResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:status, 1, type: PaymentService.V1.RefundStatus, enum: true)
  field(:refund_amount, 2, type: :int64, json_name: "refundAmount")
end

defmodule PaymentService.V1.ListPaymentsRequest do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.ListPaymentsRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:order_id, 1, type: :string, json_name: "orderId")
  field(:customer_id, 2, type: :string, json_name: "customerId")
  field(:shop_id, 3, type: :string, json_name: "shopId")

  field(:status_filter, 4,
    repeated: true,
    type: PaymentService.V1.PaymentStatus,
    enum: true,
    json_name: "statusFilter"
  )

  field(:date_from, 5, type: :string, json_name: "dateFrom")
  field(:date_to, 6, type: :string, json_name: "dateTo")

  field(:payment_method, 7,
    type: PaymentService.V1.PaymentMethodType,
    enum: true,
    json_name: "paymentMethod"
  )

  field(:page, 8, type: :int32)
  field(:page_size, 9, type: :int32, json_name: "pageSize")
end

defmodule PaymentService.V1.ListPaymentsResponse do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.ListPaymentsResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:payments, 1, repeated: true, type: PaymentService.V1.Payment)
  field(:total_count, 2, type: :int32, json_name: "totalCount")
end

defmodule PaymentService.V1.GetPaymentDetailRequest do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.GetPaymentDetailRequest",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:payment_id, 1, type: :string, json_name: "paymentId")
  field(:requester_id, 2, type: :string, json_name: "requesterId")
  field(:requester_role, 3, type: :string, json_name: "requesterRole")
end

defmodule PaymentService.V1.GetPaymentDetailResponse do
  @moduledoc false

  use Protobuf,
    full_name: "payment_service.v1.GetPaymentDetailResponse",
    protoc_gen_elixir_version: "0.16.0",
    syntax: :proto3

  field(:payment, 1, type: PaymentService.V1.Payment)
end

defmodule PaymentService.V1.PaymentService.Service do
  @moduledoc false

  use GRPC.Service,
    name: "payment_service.v1.PaymentService",
    protoc_gen_elixir_version: "0.16.0"

  rpc(
    :GetPaymentStatus,
    PaymentService.V1.GetPaymentStatusRequest,
    PaymentService.V1.GetPaymentStatusResponse
  )

  rpc(
    :ConfirmCODPayment,
    PaymentService.V1.ConfirmCODPaymentRequest,
    PaymentService.V1.ConfirmCODPaymentResponse
  )

  rpc(
    :CreateRefund,
    PaymentService.V1.CreateRefundRequest,
    PaymentService.V1.CreateRefundResponse
  )

  rpc(
    :GetRefundStatus,
    PaymentService.V1.GetRefundStatusRequest,
    PaymentService.V1.GetRefundStatusResponse
  )

  rpc(
    :ListPayments,
    PaymentService.V1.ListPaymentsRequest,
    PaymentService.V1.ListPaymentsResponse
  )

  rpc(
    :GetPaymentDetail,
    PaymentService.V1.GetPaymentDetailRequest,
    PaymentService.V1.GetPaymentDetailResponse
  )
end

defmodule PaymentService.V1.PaymentService.Stub do
  @moduledoc false

  use GRPC.Stub, service: PaymentService.V1.PaymentService.Service
end
