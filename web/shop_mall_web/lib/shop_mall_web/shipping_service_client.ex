defmodule ShopMallWeb.ShippingServiceClient do
  @moduledoc """
  配送サービス(shipping service)への gRPC 呼び出しをまとめたクライアント。
  加盟店画面の出荷操作から利用する。
  """

  alias ShippingService.V1.{
    GetShipmentByOrderRequest,
    RegisterTrackingNumberRequest,
    ShippingService.Stub,
    UpdateShipmentStatusRequest
  }

  @doc "注文 ID から出荷情報を取得する。"
  def get_shipment_by_order(order_id) do
    request = %GetShipmentByOrderRequest{order_id: order_id}

    with_channel(fn channel -> Stub.get_shipment_by_order(channel, request) end)
  end

  @doc "追跡番号を登録する(出荷済みになる)。"
  def register_tracking_number(shipment_id, tracking_number) do
    request = %RegisterTrackingNumberRequest{
      shipment_id: shipment_id,
      tracking_number: tracking_number
    }

    with_channel(fn channel -> Stub.register_tracking_number(channel, request) end)
  end

  @doc """
  配送状態を更新する。:SHIPMENT_STATUS_DELIVERED にすると、
  配送サービスが代引き決済の集金確定を決済サービスへ通知する。
  """
  def update_shipment_status(shipment_id, new_status) do
    request = %UpdateShipmentStatusRequest{shipment_id: shipment_id, new_status: new_status}

    with_channel(fn channel -> Stub.update_shipment_status(channel, request) end)
  end

  @doc "出荷 ID から最新の詳細を取得する。"
  def get_shipment_detail(shipment_id) do
    request = %ShippingService.V1.GetShipmentDetailRequest{shipment_id: shipment_id}

    with_channel(fn channel ->
      ShippingService.V1.ShippingService.Stub.get_shipment_detail(channel, request)
    end)
  end

  @doc "住所の妥当性を検証する(郵便番号・都道府県が必須)。"
  def validate_address(postal_code, prefecture, city, address_line) do
    request = %ShippingService.V1.ValidateAddressRequest{
      postal_code: postal_code,
      prefecture: prefecture,
      city: city,
      address_line: address_line
    }

    with_channel(fn channel ->
      ShippingService.V1.ShippingService.Stub.validate_address(channel, request)
    end)
  end

  @doc "住所を正規化した表記に整える。"
  def normalize_address(postal_code, prefecture, city, address_line) do
    request = %ShippingService.V1.NormalizeAddressRequest{
      postal_code: postal_code,
      prefecture: prefecture,
      city: city,
      address_line: address_line
    }

    with_channel(fn channel ->
      ShippingService.V1.ShippingService.Stub.normalize_address(channel, request)
    end)
  end

  defp with_channel(fun) do
    host = System.get_env("SHIPPING_SERVICE_HOST", "localhost")
    port = System.get_env("SHIPPING_SERVICE_PORT", "50057")

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
        {:error, "配送サービスに接続できません: #{inspect(reason)}"}
    end
  end

  # ---- 表示ヘルパー ----

  def status_label(:SHIPMENT_STATUS_PENDING), do: "出荷待ち"
  def status_label(:SHIPMENT_STATUS_PREPARING), do: "梱包中"
  def status_label(:SHIPMENT_STATUS_SHIPPED), do: "出荷済み"
  def status_label(:SHIPMENT_STATUS_IN_TRANSIT), do: "輸送中"
  def status_label(:SHIPMENT_STATUS_DELIVERED), do: "配達完了"
  def status_label(:SHIPMENT_STATUS_FAILED), do: "配送失敗"
  def status_label(_), do: "不明"

  def status_color(:SHIPMENT_STATUS_DELIVERED), do: "bg-green-100 text-green-800"
  def status_color(:SHIPMENT_STATUS_FAILED), do: "bg-red-100 text-red-800"
  def status_color(:SHIPMENT_STATUS_IN_TRANSIT), do: "bg-blue-100 text-blue-800"
  def status_color(:SHIPMENT_STATUS_SHIPPED), do: "bg-blue-100 text-blue-800"
  def status_color(_), do: "bg-yellow-100 text-yellow-800"
end
