defmodule ShopMallWebWeb.Owner.PaymentListLive do
  @moduledoc """
  加盟店(店舗オーナー)向けの決済確認画面。
  自店舗宛の決済一覧の確認と、代引き注文の集金確定(配達完了報告)を行う。

  注: このサンプルでは payments テーブルに shop_id が無いため、
  店舗によるサーバー側の絞り込みは行っていない(全件が見える)。
  """
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.PaymentServiceClient, as: Payments
  alias ShopMallWeb.ShippingServiceClient, as: Shipping

  @method_filters [
    {"all", "すべて"},
    {"cod", "代金引換"},
    {"card", "クレジットカード"}
  ]

  @impl true
  def mount(_params, session, socket) do
    {:ok,
     socket
     |> assign(:current_user_id, session["user_id"])
     |> assign(:payments, [])
     |> assign(:total_count, 0)
     |> assign(:loading, true)
     |> assign(:error, nil)
     |> assign(:method_filter, "all")
     |> assign(:method_filters, @method_filters)
     |> assign(:shipment, nil)
     |> assign(:shipment_error, nil)
     |> load_payments()}
  end

  defp load_payments(socket) do
    filters =
      case socket.assigns.method_filter do
        "cod" -> [payment_method: :CASH_ON_DELIVERY]
        "card" -> [payment_method: :CREDIT_CARD]
        _ -> []
      end

    case Payments.list_payments(filters) do
      {:ok, response} ->
        socket
        |> assign(:payments, response.payments)
        |> assign(:total_count, response.total_count)
        |> assign(:loading, false)
        |> assign(:error, nil)

      {:error, reason} ->
        socket
        |> assign(:loading, false)
        |> assign(:error, "決済一覧の取得に失敗しました: #{reason}")
    end
  end

  @impl true
  def handle_event("filter_method", %{"method" => method}, socket) do
    {:noreply,
     socket
     |> assign(:method_filter, method)
     |> assign(:loading, true)
     |> load_payments()}
  end

  @impl true
  def handle_event(
        "confirm_cod",
        %{"payment-id" => payment_id, "order-id" => order_id},
        socket
      ) do
    case Payments.confirm_cod_payment(payment_id, order_id) do
      {:ok, _response} ->
        # 確定後の状態を決済サービスに問い合わせて検証・表示する
        status_note =
          case Payments.get_payment_status(payment_id) do
            {:ok, st} -> "(現在の状態: #{Payments.status_label(st.status)})"
            {:error, _} -> ""
          end

        {:noreply,
         socket
         |> put_flash(:info, "集金を確定しました#{status_note}")
         |> load_payments()}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "集金の確定に失敗しました: #{reason}")}
    end
  end

  # --- 配送(shipping サービス連携) ---

  @impl true
  def handle_event("shipment_detail", _params, socket) do
    case Shipping.get_shipment_detail(socket.assigns.shipment.id) do
      {:ok, resp} ->
        {:noreply,
         socket |> assign(:shipment, resp.shipment) |> put_flash(:info, "最新の配送詳細を取得しました")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "詳細の取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("open_shipment", %{"order-id" => order_id}, socket) do
    case Shipping.get_shipment_by_order(order_id) do
      {:ok, response} ->
        {:noreply, socket |> assign(:shipment, response.shipment) |> assign(:shipment_error, nil)}

      {:error, reason} ->
        {:noreply,
         socket
         |> assign(:shipment, nil)
         |> assign(:shipment_error, "出荷情報が見つかりません: #{reason}")
         |> put_flash(:error, "出荷情報が見つかりません: #{reason}")}
    end
  end

  @impl true
  def handle_event("close_shipment", _params, socket) do
    {:noreply, socket |> assign(:shipment, nil) |> assign(:shipment_error, nil)}
  end

  @impl true
  def handle_event("register_tracking", %{"tracking_number" => tracking_number}, socket) do
    shipment = socket.assigns.shipment

    case Shipping.register_tracking_number(shipment.id, tracking_number) do
      {:ok, _response} ->
        {:noreply,
         socket
         |> put_flash(:info, "追跡番号を登録しました(出荷済み)")
         |> reload_shipment(shipment.order_id)}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "追跡番号の登録に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("mark_in_transit", _params, socket) do
    update_shipment(socket, :SHIPMENT_STATUS_IN_TRANSIT)
  end

  @impl true
  def handle_event("mark_delivered", _params, socket) do
    # 配達完了にすると shipping サービスが代引きの集金確定を payment に通知する
    update_shipment(socket, :SHIPMENT_STATUS_DELIVERED)
  end

  defp update_shipment(socket, new_status) do
    shipment = socket.assigns.shipment

    case Shipping.update_shipment_status(shipment.id, new_status) do
      {:ok, response} ->
        {:noreply,
         socket
         |> put_flash(:info, response.message)
         |> reload_shipment(shipment.order_id)
         |> load_payments()}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "配送状態の更新に失敗しました: #{reason}")}
    end
  end

  defp reload_shipment(socket, order_id) do
    case Shipping.get_shipment_by_order(order_id) do
      {:ok, response} -> assign(socket, :shipment, response.shipment)
      {:error, _} -> assign(socket, :shipment, nil)
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <nav class="bg-emerald-800 shadow-md">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16 items-center">
            <span class="text-xl font-bold text-white">🏪 加盟店ポータル</span>
            <div class="flex items-center space-x-4">
              <.link
                navigate="/owner/dashboard"
                class="text-emerald-200 hover:text-white text-sm font-medium"
              >
                ダッシュボード
              </.link>
              <.link
                navigate="/owner/products"
                class="text-emerald-200 hover:text-white text-sm font-medium"
              >
                商品管理
              </.link>
              <span class="text-white text-sm font-medium">決済確認</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between mb-2">
          <h1 class="text-2xl font-bold text-gray-900">
            決済確認 <span class="text-sm font-normal text-gray-500">({@total_count}件)</span>
          </h1>
          <div class="flex space-x-2">
            <button
              :for={{value, label} <- @method_filters}
              phx-click="filter_method"
              phx-value-method={value}
              class={"px-3 py-1.5 rounded-full text-sm font-medium transition-colors " <>
                if(@method_filter == value,
                  do: "bg-emerald-700 text-white",
                  else: "bg-white text-gray-600 border border-gray-300 hover:bg-gray-100"
                )}
            >
              {label}
            </button>
          </div>
        </div>
        <p class="text-sm text-gray-500 mb-6">
          代金引換の注文は、商品のお届け完了後に「集金確定」を押してください。
        </p>

        <%= if @error do %>
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            {@error}
          </div>
        <% end %>

        <%= if @loading do %>
          <div class="text-center py-12">
            <div class="inline-block animate-spin rounded-full h-10 w-10 border-b-2 border-emerald-700">
            </div>
          </div>
        <% else %>
          <div class="bg-white shadow rounded-lg overflow-hidden">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    注文ID
                  </th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    支払い方法
                  </th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    金額
                  </th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    状態
                  </th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    受注日時
                  </th>
                  <th class="px-4 py-3"></th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <%= if @payments == [] do %>
                  <tr>
                    <td colspan="6" class="px-4 py-8 text-center text-gray-500">
                      該当する決済はありません
                    </td>
                  </tr>
                <% end %>
                <tr :for={payment <- @payments} class="hover:bg-gray-50">
                  <td class="px-4 py-3 text-sm font-mono text-gray-600">
                    {String.slice(payment.order_id, 0, 8)}…
                  </td>
                  <td class="px-4 py-3 text-sm text-gray-900">
                    {Payments.method_label(payment.payment_method)}
                  </td>
                  <td class="px-4 py-3 text-sm text-right font-semibold text-gray-900">
                    {Payments.format_amount(payment.amount)}
                  </td>
                  <td class="px-4 py-3">
                    <span class={"inline-flex px-2 py-0.5 rounded-full text-xs font-medium " <>
                      Payments.status_color(payment.status)}>
                      {Payments.status_label(payment.status)}
                    </span>
                  </td>
                  <td class="px-4 py-3 text-sm text-gray-500">
                    {Payments.format_timestamp(payment.created_at)}
                  </td>
                  <td class="px-4 py-3 text-right whitespace-nowrap space-x-2">
                    <button
                      phx-click="open_shipment"
                      phx-value-order-id={payment.order_id}
                      class="px-3 py-1.5 text-sm font-medium text-emerald-700 border border-emerald-300 rounded-md hover:bg-emerald-50"
                    >
                      配送
                    </button>
                    <button
                      :if={
                        payment.payment_method == :CASH_ON_DELIVERY and
                          payment.status == :PAYMENT_STATUS_PENDING
                      }
                      phx-click="confirm_cod"
                      phx-value-payment-id={payment.id}
                      phx-value-order-id={payment.order_id}
                      class="px-3 py-1.5 text-sm font-medium text-white bg-emerald-600 rounded-md hover:bg-emerald-700"
                    >
                      集金確定
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        <% end %>
        
    <!-- 配送モーダル -->
        <div
          :if={@shipment}
          class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50"
        >
          <div class="bg-white rounded-lg shadow-xl max-w-lg w-full mx-4 p-6">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-bold text-gray-900">配送情報</h2>
              <div class="space-x-2">
                <button
                  phx-click="shipment_detail"
                  class="text-sm text-emerald-700 hover:text-emerald-900"
                >
                  ↻ 最新詳細
                </button>
                <button phx-click="close_shipment" class="text-gray-400 hover:text-gray-600">
                  ✕
                </button>
              </div>
            </div>

            <dl class="grid grid-cols-3 gap-y-3 text-sm mb-4">
              <dt class="text-gray-500">注文ID</dt>
              <dd class="col-span-2 font-mono text-gray-900 break-all">{@shipment.order_id}</dd>
              <dt class="text-gray-500">状態</dt>
              <dd class="col-span-2">
                <span class={"inline-flex px-2 py-0.5 rounded-full text-xs font-medium " <>
                  Shipping.status_color(@shipment.status)}>
                  {Shipping.status_label(@shipment.status)}
                </span>
              </dd>
              <dt class="text-gray-500">配送業者</dt>
              <dd class="col-span-2 text-gray-900">{@shipment.carrier}</dd>
              <dt class="text-gray-500">追跡番号</dt>
              <dd class="col-span-2 font-mono text-gray-900">
                {if @shipment.tracking_number == "", do: "未登録", else: @shipment.tracking_number}
              </dd>
            </dl>

            <form
              :if={@shipment.tracking_number == ""}
              phx-submit="register_tracking"
              class="flex space-x-2 mb-4"
            >
              <input
                type="text"
                name="tracking_number"
                placeholder="追跡番号(例: TRK-1234)"
                required
                class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
              />
              <button
                type="submit"
                class="px-4 py-2 text-sm font-medium text-white bg-emerald-600 rounded-md hover:bg-emerald-700"
              >
                出荷する
              </button>
            </form>

            <div class="flex justify-end space-x-3">
              <button
                :if={@shipment.status in [:SHIPMENT_STATUS_SHIPPED]}
                phx-click="mark_in_transit"
                class="px-4 py-2 text-sm font-medium text-blue-700 border border-blue-300 rounded-md hover:bg-blue-50"
              >
                輸送中にする
              </button>
              <button
                :if={
                  @shipment.status in [
                    :SHIPMENT_STATUS_SHIPPED,
                    :SHIPMENT_STATUS_IN_TRANSIT
                  ]
                }
                phx-click="mark_delivered"
                class="px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700"
              >
                配達完了(代引きは集金確定)
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
    """
  end
end
