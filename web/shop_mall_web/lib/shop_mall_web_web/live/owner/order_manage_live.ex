defmodule ShopMallWebWeb.Owner.OrderManageLive do
  @moduledoc """
  加盟店の受注管理画面。
  受注一覧(ListOrders)→ 明細表示(GetOrderDetail)→
  ステータス更新(UpdateOrderStatus: 準備中/発送済み+追跡番号/配達完了)を行う。
  """
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.ShopServiceClient, as: Shops

  @carriers [
    {"YAMATO", "ヤマト運輸"},
    {"SAGAWA", "佐川急便"},
    {"JAPAN_POST", "日本郵便"}
  ]

  @impl true
  def mount(_params, session, socket) do
    owner_id = session["user_id"] || "admin-user-id"

    {:ok,
     socket
     |> assign(:owner_id, owner_id)
     |> assign(:shop, nil)
     |> assign(:orders, [])
     |> assign(:detail, nil)
     |> assign(:carriers, @carriers)
     |> assign(:loading, true)
     |> assign(:error, nil)
     |> load_shop_and_orders()}
  end

  defp load_shop_and_orders(socket) do
    case Shops.get_shops_by_owner(socket.assigns.owner_id) do
      {:ok, %{shops: [shop | _]}} ->
        socket |> assign(:shop, shop) |> load_orders()

      {:ok, _} ->
        socket |> assign(:loading, false) |> assign(:error, "店舗がありません")

      {:error, reason} ->
        socket |> assign(:loading, false) |> assign(:error, "店舗の取得に失敗しました: #{reason}")
    end
  end

  defp load_orders(socket) do
    case Shops.list_orders(socket.assigns.shop.id) do
      {:ok, response} ->
        socket
        |> assign(:orders, response.orders)
        |> assign(:loading, false)
        |> assign(:error, nil)

      {:error, reason} ->
        socket
        |> assign(:loading, false)
        |> assign(:error, "受注一覧の取得に失敗しました: #{reason}")
    end
  end

  @impl true
  def handle_event("show_detail", %{"order-id" => order_id}, socket) do
    case Shops.get_order_detail(order_id, socket.assigns.shop.id) do
      {:ok, response} ->
        {:noreply, assign(socket, :detail, response.order)}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "明細の取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("close_detail", _params, socket) do
    {:noreply, assign(socket, :detail, nil)}
  end

  @impl true
  def handle_event("mark_preparing", %{"order-id" => order_id}, socket) do
    update_status(socket, order_id, :PREPARING, "", :CARRIER_UNSPECIFIED)
  end

  @impl true
  def handle_event(
        "ship_order",
        %{"order_id" => order_id, "tracking_number" => tn, "carrier" => carrier},
        socket
      ) do
    update_status(socket, order_id, :SHIPPED, tn, safe_carrier(carrier))
  end

  @impl true
  def handle_event("mark_delivered", %{"order-id" => order_id}, socket) do
    update_status(socket, order_id, :DELIVERED, "", :CARRIER_UNSPECIFIED)
  end

  defp update_status(socket, order_id, status, tracking_number, carrier) do
    case Shops.update_order_status(
           order_id,
           socket.assigns.shop.id,
           status,
           tracking_number,
           carrier
         ) do
      {:ok, response} ->
        {:noreply,
         socket
         |> assign(:detail, nil)
         |> put_flash(:info, "受注を「#{Shops.order_status_label(response.status)}」にしました")
         |> load_orders()}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "ステータス更新に失敗しました: #{reason}")}
    end
  end

  defp safe_carrier(value) do
    valid = Enum.map(@carriers, fn {v, _} -> v end)
    if value in valid, do: String.to_existing_atom(value), else: :CARRIER_UNSPECIFIED
  end

  defp format_amount(v) when is_binary(v) do
    case Integer.parse(v) do
      {n, _} ->
        "¥" <>
          (n
           |> Integer.to_string()
           |> String.reverse()
           |> String.replace(~r/(\d{3})(?=\d)/, "\\1,")
           |> String.reverse())

      :error ->
        "¥" <> v
    end
  end

  defp format_amount(_), do: "¥0"

  defp format_ts(%Google.Protobuf.Timestamp{seconds: s}) do
    s |> Kernel.+(9 * 3600) |> DateTime.from_unix!() |> Calendar.strftime("%Y-%m-%d %H:%M")
  end

  defp format_ts(_), do: "-"

  defp status_color(:DELIVERED), do: "bg-green-100 text-green-800"
  defp status_color(:CANCELLED), do: "bg-gray-200 text-gray-700"
  defp status_color(:SHIPPED), do: "bg-blue-100 text-blue-800"
  defp status_color(_), do: "bg-yellow-100 text-yellow-800"

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
                navigate="/owner/sales"
                class="text-emerald-200 hover:text-white text-sm font-medium"
              >
                売上レポート
              </.link>
              <span class="text-white text-sm font-medium">受注管理</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-5xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-6">
          受注管理 <span :if={@shop} class="text-sm font-normal text-gray-500">({@shop.name})</span>
        </h1>

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
                    注文番号
                  </th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">状態</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">金額</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">
                    受注日時
                  </th>
                  <th class="px-4 py-3"></th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <%= if @orders == [] do %>
                  <tr>
                    <td colspan="5" class="px-4 py-8 text-center text-gray-500">受注はまだありません</td>
                  </tr>
                <% end %>
                <tr :for={order <- @orders} class="hover:bg-gray-50">
                  <td class="px-4 py-3 text-sm font-mono text-gray-700">{order.order_number}</td>
                  <td class="px-4 py-3">
                    <span class={"inline-flex px-2 py-0.5 rounded-full text-xs font-medium " <> status_color(order.status)}>
                      {Shops.order_status_label(order.status)}
                    </span>
                  </td>
                  <td class="px-4 py-3 text-sm text-right font-semibold text-gray-900">
                    {format_amount(order.total_amount)}
                  </td>
                  <td class="px-4 py-3 text-sm text-gray-500">{format_ts(order.created_at)}</td>
                  <td class="px-4 py-3 text-right whitespace-nowrap space-x-2">
                    <button
                      phx-click="show_detail"
                      phx-value-order-id={order.id}
                      class="text-emerald-700 hover:text-emerald-900 text-sm font-medium"
                    >
                      明細
                    </button>
                    <button
                      :if={order.status == :RECEIVED}
                      phx-click="mark_preparing"
                      phx-value-order-id={order.id}
                      class="px-3 py-1.5 text-sm font-medium text-white bg-emerald-600 rounded-md hover:bg-emerald-700"
                    >
                      準備開始
                    </button>
                    <button
                      :if={order.status == :SHIPPED}
                      phx-click="mark_delivered"
                      phx-value-order-id={order.id}
                      class="px-3 py-1.5 text-sm font-medium text-green-700 border border-green-300 rounded-md hover:bg-green-50"
                    >
                      配達完了
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        <% end %>
        
    <!-- 明細モーダル -->
        <div
          :if={@detail}
          class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50"
        >
          <div class="bg-white rounded-lg shadow-xl max-w-lg w-full mx-4 p-6 max-h-[80vh] overflow-y-auto">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-bold text-gray-900">
                受注明細 {@detail.order && @detail.order.order_number}
              </h2>
              <button phx-click="close_detail" class="text-gray-400 hover:text-gray-600">✕</button>
            </div>

            <%= if @detail.order do %>
              <dl class="grid grid-cols-3 gap-y-2 text-sm mb-4">
                <dt class="text-gray-500">状態</dt>
                <dd class="col-span-2">
                  <span class={"inline-flex px-2 py-0.5 rounded-full text-xs font-medium " <> status_color(@detail.order.status)}>
                    {Shops.order_status_label(@detail.order.status)}
                  </span>
                </dd>
                <dt class="text-gray-500">合計</dt>
                <dd class="col-span-2 font-semibold">{format_amount(@detail.order.total_amount)}</dd>
                <dt class="text-gray-500">お届け先</dt>
                <dd class="col-span-2">{@detail.order.shipping_address}</dd>
                <dt :if={@detail.order.tracking_number != ""} class="text-gray-500">追跡番号</dt>
                <dd :if={@detail.order.tracking_number != ""} class="col-span-2 font-mono">
                  {@detail.order.tracking_number}({Shops.carrier_label(@detail.order.carrier)})
                </dd>
              </dl>
            <% end %>

            <h3 class="text-sm font-semibold text-gray-700 mb-2">商品</h3>
            <table class="min-w-full text-sm mb-4">
              <tbody class="divide-y divide-gray-100">
                <tr :for={item <- @detail.items}>
                  <td class="py-1.5">{item.product_name}</td>
                  <td class="py-1.5 text-right text-gray-500">×{item.quantity}</td>
                  <td class="py-1.5 text-right font-medium">{format_amount(item.subtotal)}</td>
                </tr>
              </tbody>
            </table>

            <%= if @detail.order && @detail.order.status in [:RECEIVED, :PREPARING] do %>
              <h3 class="text-sm font-semibold text-gray-700 mb-2">発送する</h3>
              <form phx-submit="ship_order" class="flex items-end space-x-2">
                <input type="hidden" name="order_id" value={@detail.order.id} />
                <div class="flex-1">
                  <label class="block text-xs text-gray-500 mb-1">追跡番号</label>
                  <input
                    type="text"
                    name="tracking_number"
                    required
                    placeholder="TRK-1234"
                    class="w-full border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                  />
                </div>
                <div>
                  <label class="block text-xs text-gray-500 mb-1">配送業者</label>
                  <select name="carrier" class="border border-gray-300 rounded-md px-2 py-1.5 text-sm">
                    <option :for={{v, label} <- @carriers} value={v}>{label}</option>
                  </select>
                </div>
                <button
                  type="submit"
                  class="px-4 py-2 text-sm font-medium text-white bg-emerald-600 rounded-md hover:bg-emerald-700"
                >
                  発送済みにする
                </button>
              </form>
            <% end %>
          </div>
        </div>
      </main>
    </div>
    """
  end
end
