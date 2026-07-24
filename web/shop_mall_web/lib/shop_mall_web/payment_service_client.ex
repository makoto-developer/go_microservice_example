defmodule ShopMallWeb.PaymentServiceClient do
  @moduledoc """
  決済サービス(payment service)への gRPC 呼び出しをまとめたクライアント。
  管理者画面・加盟店画面から利用する。
  """

  alias PaymentService.V1.{
    ConfirmCODPaymentRequest,
    CreateRefundRequest,
    GetPaymentDetailRequest,
    GetPaymentStatusRequest,
    GetRefundStatusRequest,
    ListPaymentsRequest,
    PaymentService.Stub
  }

  @doc "決済一覧を取得する。filters は ListPaymentsRequest のフィールドをキーワードで渡す。"
  def list_payments(filters \\ []) do
    request = struct!(ListPaymentsRequest, Keyword.merge([page: 1, page_size: 50], filters))

    with_channel(fn channel -> Stub.list_payments(channel, request) end)
  end

  @doc "決済の詳細を取得する。"
  def get_payment_detail(payment_id, requester_id, requester_role) do
    request = %GetPaymentDetailRequest{
      payment_id: payment_id,
      requester_id: requester_id,
      requester_role: requester_role
    }

    with_channel(fn channel -> Stub.get_payment_detail(channel, request) end)
  end

  @doc "決済を返金する。amount が空文字なら全額返金。"
  def create_refund(payment_id, amount \\ "", reason) do
    request = %CreateRefundRequest{
      payment_id: payment_id,
      amount: amount,
      reason: reason
    }

    with_channel(fn channel -> Stub.create_refund(channel, request) end)
  end

  @doc "代引き決済の集金を確定する(配達完了時)。"
  def confirm_cod_payment(payment_id, order_id) do
    request = %ConfirmCODPaymentRequest{payment_id: payment_id, order_id: order_id}

    with_channel(fn channel -> Stub.confirm_cod_payment(channel, request) end)
  end

  @doc "決済の現在状態を取得する(操作後の検証用)。"
  def get_payment_status(payment_id) do
    request = %GetPaymentStatusRequest{payment_id: payment_id}

    with_channel(fn channel -> Stub.get_payment_status(channel, request) end)
  end

  @doc "返金の状態を取得する。"
  def get_refund_status(refund_id) do
    request = %GetRefundStatusRequest{refund_id: refund_id}

    with_channel(fn channel -> Stub.get_refund_status(channel, request) end)
  end

  defp with_channel(fun) do
    host = System.get_env("PAYMENT_SERVICE_HOST", "localhost")
    port = System.get_env("PAYMENT_SERVICE_PORT", "50056")

    case GRPC.Stub.connect("#{host}:#{port}") do
      {:ok, channel} ->
        try do
          case fun.(channel) do
            {:ok, response} -> {:ok, response}
            {:error, %GRPC.RPCError{message: message}} -> {:error, message}
            {:error, reason} -> {:error, inspect(reason)}
          end
        after
          GRPC.Stub.disconnect(channel)
        end

      {:error, reason} ->
        {:error, "決済サービスに接続できません: #{inspect(reason)}"}
    end
  end

  # ---- 表示ヘルパー ----

  def status_label(:PAYMENT_STATUS_PENDING), do: "支払い待ち"
  def status_label(:PAYMENT_STATUS_PROCESSING), do: "処理中"
  def status_label(:PAYMENT_STATUS_REQUIRES_AUTHENTICATION), do: "要認証"
  def status_label(:PAYMENT_STATUS_SUCCEEDED), do: "支払い済み"
  def status_label(:PAYMENT_STATUS_FAILED), do: "失敗"
  def status_label(:PAYMENT_STATUS_REFUNDING), do: "返金中"
  def status_label(:PAYMENT_STATUS_REFUNDED), do: "返金済み"
  def status_label(_), do: "不明"

  def status_color(:PAYMENT_STATUS_SUCCEEDED), do: "bg-green-100 text-green-800"
  def status_color(:PAYMENT_STATUS_PENDING), do: "bg-yellow-100 text-yellow-800"
  def status_color(:PAYMENT_STATUS_FAILED), do: "bg-red-100 text-red-800"
  def status_color(:PAYMENT_STATUS_REFUNDED), do: "bg-gray-200 text-gray-700"
  def status_color(_), do: "bg-gray-100 text-gray-600"

  def method_label(:CREDIT_CARD), do: "クレジットカード"
  def method_label(:CASH_ON_DELIVERY), do: "代金引換"
  def method_label(_), do: "不明"

  def format_amount(amount) when is_binary(amount) do
    case Integer.parse(amount) do
      {n, _} -> "¥" <> delimit(n)
      :error -> "¥" <> amount
    end
  end

  def format_amount(_), do: "¥0"

  defp delimit(num) do
    num
    |> Integer.to_string()
    |> String.reverse()
    |> String.replace(~r/(\d{3})(?=\d)/, "\\1,")
    |> String.reverse()
  end

  def format_timestamp(%Google.Protobuf.Timestamp{seconds: seconds}) do
    # タイムゾーン DB に依存せず JST(+9:00 固定)で表示する
    seconds
    |> Kernel.+(9 * 3600)
    |> DateTime.from_unix!()
    |> Calendar.strftime("%Y-%m-%d %H:%M")
  end

  def format_timestamp(_), do: "-"
end
