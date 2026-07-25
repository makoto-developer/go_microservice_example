defmodule ShopMallWeb.ShopServiceClient do
  @moduledoc """
  ショップサービス(shop service)への gRPC 呼び出しをまとめたクライアント。
  加盟店ポータルの店舗設定・売上レポート・受注管理から利用する。
  """

  alias ShopService.V1.{
    ExportSalesDataRequest,
    GetOrderDetailRequest,
    GetSalesReportRequest,
    GetShopsByOwnerRequest,
    ListOrdersRequest,
    ManageVariationRequest,
    ShopService.Stub,
    ToggleShopPublishRequest,
    UpdateOrderStatusRequest,
    UpdateShopRequest,
    UploadProductImageRequest
  }

  @doc "オーナーが所有する店舗の一覧を取得する。"
  def get_shops_by_owner(owner_id) do
    with_channel(fn ch ->
      Stub.get_shops_by_owner(ch, %GetShopsByOwnerRequest{owner_id: owner_id})
    end)
  end

  @doc "店舗情報(名前・説明・営業時間・返品ポリシー)を更新する。"
  def update_shop(attrs) do
    request = %UpdateShopRequest{
      shop_id: attrs[:shop_id],
      name: attrs[:name] || "",
      description: attrs[:description] || "",
      logo_url: attrs[:logo_url] || "",
      business_hours: attrs[:business_hours] || "",
      return_policy: attrs[:return_policy] || ""
    }

    with_channel(fn ch -> Stub.update_shop(ch, request) end)
  end

  @doc "店舗の公開/非公開を切り替える。"
  def toggle_shop_publish(shop_id, published) do
    request = %ToggleShopPublishRequest{shop_id: shop_id, published: published}
    with_channel(fn ch -> Stub.toggle_shop_publish(ch, request) end)
  end

  @doc "売上レポート(日別データ+サマリ)を取得する。"
  def get_sales_report(shop_id, report_type, date_from, date_to) do
    request = %GetSalesReportRequest{
      shop_id: shop_id,
      report_type: report_type,
      date_from: date_from,
      date_to: date_to
    }

    with_channel(fn ch -> Stub.get_sales_report(ch, request) end)
  end

  @doc "売上データを CSV エクスポートする(ダウンロード URL が返る)。"
  def export_sales_data(shop_id, date_from, date_to) do
    request = %ExportSalesDataRequest{shop_id: shop_id, date_from: date_from, date_to: date_to}
    with_channel(fn ch -> Stub.export_sales_data(ch, request) end)
  end

  @doc "店舗の受注一覧を取得する。"
  def list_orders(shop_id) do
    request = %ListOrdersRequest{
      shop_id: shop_id,
      status: :ORDER_STATUS_UNSPECIFIED,
      date_from: "",
      date_to: "",
      customer_name: "",
      product_name: "",
      sort_by: "created_at",
      sort_order: "desc"
    }

    with_channel(fn ch -> Stub.list_orders(ch, request) end)
  end

  @doc "受注の明細(注文+商品行)を取得する。"
  def get_order_detail(order_id, shop_id) do
    request = %GetOrderDetailRequest{order_id: order_id, shop_id: shop_id}
    with_channel(fn ch -> Stub.get_order_detail(ch, request) end)
  end

  @doc "受注のステータスを進める(発送時は追跡番号と配送業者を添える)。"
  def update_order_status(
        order_id,
        shop_id,
        new_status,
        tracking_number \\ "",
        carrier \\ :CARRIER_UNSPECIFIED
      ) do
    request = %UpdateOrderStatusRequest{
      order_id: order_id,
      shop_id: shop_id,
      new_status: new_status,
      tracking_number: tracking_number,
      carrier: carrier
    }

    with_channel(fn ch -> Stub.update_order_status(ch, request) end)
  end

  @doc "商品バリエーション(SKU・属性・価格・在庫)を登録する。"
  def manage_variation(product_id, variations) do
    request = %ManageVariationRequest{product_id: product_id, variations: variations}
    with_channel(fn ch -> Stub.manage_variation(ch, request) end)
  end

  @doc "商品画像をアップロードする。"
  def upload_product_image(product_id, image_data, display_order \\ 0) do
    request = %UploadProductImageRequest{
      product_id: product_id,
      image_data: image_data,
      display_order: display_order
    }

    with_channel(fn ch -> Stub.upload_product_image(ch, request) end)
  end

  defp with_channel(fun) do
    host = System.get_env("SHOP_SERVICE_HOST", "localhost")
    port = System.get_env("SHOP_SERVICE_PORT", "22101")

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
        {:error, "ショップサービスに接続できません: #{inspect(reason)}"}
    end
  end

  # ---- 表示ヘルパー ----

  def order_status_label(:RECEIVED), do: "受注"
  def order_status_label(:PREPARING), do: "準備中"
  def order_status_label(:SHIPPED), do: "発送済み"
  def order_status_label(:DELIVERED), do: "配達完了"
  def order_status_label(:CANCELLED), do: "キャンセル"
  def order_status_label(_), do: "不明"

  def carrier_label(:YAMATO), do: "ヤマト運輸"
  def carrier_label(:SAGAWA), do: "佐川急便"
  def carrier_label(:JAPAN_POST), do: "日本郵便"
  def carrier_label(_), do: "未指定"
end
