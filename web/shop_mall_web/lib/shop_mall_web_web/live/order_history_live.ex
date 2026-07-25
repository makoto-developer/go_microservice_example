defmodule ShopMallWebWeb.OrderHistoryLive do
  @moduledoc """
  顧客向けの注文履歴画面。
  自分の注文一覧の確認と、発送前の注文のキャンセル(決済済みなら自動で返金される)を行う。
  """
  use ShopMallWebWeb, :live_view

  alias OrderService.V1.{
    CancelOrderRequest,
    CreateReorderRequest,
    GetOrderDetailRequest,
    GetOrderStatusHistoryRequest,
    ListOrdersRequest,
    OrderService.Stub
  }

  @cancel_reasons [
    {"CUSTOMER_NO_LONGER_NEEDED", "不要になった"},
    {"CUSTOMER_ORDERED_BY_MISTAKE", "誤って注文した"},
    {"CUSTOMER_DELIVERY_TIME_ISSUE", "お届け時期が合わない"},
    {"CUSTOMER_OTHER", "その他"}
  ]

  # キャンセルできるのは発送前まで
  @cancellable_statuses [:PENDING, :PAYMENT_PROCESSING, :CONFIRMED, :PREPARING]

  @impl true
  def mount(_params, session, socket) do
    {:ok,
     socket
     |> assign(:current_user_id, session["user_id"])
     |> assign(:orders, [])
     |> assign(:loading, true)
     |> assign(:error, nil)
     |> assign(:cancel_target, nil)
     |> assign(:cancel_reasons, @cancel_reasons)
     |> assign(:detail, nil)
     |> assign(:history, [])
     |> load_orders()}
  end

  defp load_orders(%{assigns: %{current_user_id: nil}} = socket) do
    socket
    |> assign(:loading, false)
    |> assign(:error, "注文履歴を見るにはログインが必要です")
  end

  defp load_orders(socket) do
    request = %ListOrdersRequest{customer_id: socket.assigns.current_user_id}

    with {:ok, channel} <- connect_order_service(),
         {:ok, response} <- Stub.list_orders(channel, request) do
      socket
      |> assign(:orders, response.orders)
      |> assign(:loading, false)
      |> assign(:error, nil)
    else
      {:error, %GRPC.RPCError{message: message}} ->
        socket |> assign(:loading, false) |> assign(:error, "注文履歴の取得に失敗しました: #{message}")

      {:error, reason} ->
        socket |> assign(:loading, false) |> assign(:error, "注文履歴の取得に失敗しました: #{inspect(reason)}")
    end
  end

  defp connect_order_service do
    host = System.get_env("ORDER_SERVICE_HOST", "localhost")
    port = System.get_env("ORDER_SERVICE_PORT", "50055")
    GRPC.Stub.connect("#{host}:#{port}")
  end

  @impl true
  def handle_event("show_detail", %{"order-id" => order_id}, socket) do
    with {:ok, channel} <- connect_order_service(),
         {:ok, detail} <-
           Stub.get_order_detail(channel, %GetOrderDetailRequest{
             order_id: order_id,
             requester_id: socket.assigns.current_user_id || "",
             requester_role: "customer"
           }),
         {:ok, hist} <-
           Stub.get_order_status_history(channel, %GetOrderStatusHistoryRequest{
             order_id: order_id
           }) do
      {:noreply, socket |> assign(:detail, detail.order) |> assign(:history, hist.history)}
    else
      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "詳細の取得に失敗しました: #{inspect(reason)}")}
    end
  end

  @impl true
  def handle_event("close_detail", _params, socket) do
    {:noreply, socket |> assign(:detail, nil) |> assign(:history, [])}
  end

  @impl true
  def handle_event("reorder", %{"order-id" => order_id}, socket) do
    request = %CreateReorderRequest{
      customer_id: socket.assigns.current_user_id || "",
      original_order_id: order_id
    }

    with {:ok, channel} <- connect_order_service(),
         {:ok, response} <- Stub.create_reorder(channel, request) do
      {:noreply,
       socket
       |> put_flash(:info, "同じ内容で再注文しました(注文ID: #{response.order_id})")
       |> load_orders()}
    else
      {:error, %GRPC.RPCError{message: message}} ->
        {:noreply, put_flash(socket, :error, "再注文に失敗しました: #{message}")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "再注文に失敗しました: #{inspect(reason)}")}
    end
  end

  @impl true
  def handle_event("open_cancel", %{"order-id" => order_id}, socket) do
    {:noreply, assign(socket, :cancel_target, order_id)}
  end

  @impl true
  def handle_event("close_cancel", _params, socket) do
    {:noreply, assign(socket, :cancel_target, nil)}
  end

  @impl true
  def handle_event("execute_cancel", %{"reason" => reason}, socket) do
    request = %CancelOrderRequest{
      order_id: socket.assigns.cancel_target,
      cancelled_by: socket.assigns.current_user_id || "",
      cancel_reason: safe_reason(reason),
      cancel_note: ""
    }

    with {:ok, channel} <- connect_order_service(),
         {:ok, _response} <- Stub.cancel_order(channel, request) do
      {:noreply,
       socket
       |> assign(:cancel_target, nil)
       |> put_flash(:info, "注文をキャンセルしました。お支払い済みの場合は自動的に返金されます。")
       |> load_orders()}
    else
      {:error, %GRPC.RPCError{message: message}} ->
        {:noreply,
         socket
         |> assign(:cancel_target, nil)
         |> put_flash(:error, "キャンセルに失敗しました: #{message}")}

      {:error, reason} ->
        {:noreply,
         socket
         |> assign(:cancel_target, nil)
         |> put_flash(:error, "キャンセルに失敗しました: #{inspect(reason)}")}
    end
  end

  defp safe_reason(reason) do
    valid = Enum.map(@cancel_reasons, fn {value, _} -> value end)
    if reason in valid, do: String.to_existing_atom(reason), else: :CUSTOMER_OTHER
  end

  defp cancellable?(order), do: order.status in @cancellable_statuses

  defp status_label(:PENDING), do: "受付中"
  defp status_label(:PAYMENT_PROCESSING), do: "支払い済み"
  defp status_label(:PAYMENT_FAILED), do: "支払い失敗"
  defp status_label(:CONFIRMED), do: "注文確定(代引き)"
  defp status_label(:PREPARING), do: "発送準備中"
  defp status_label(:SHIPPED), do: "発送済み"
  defp status_label(:DELIVERED), do: "お届け済み"
  defp status_label(:CANCELLED), do: "キャンセル済み"
  defp status_label(_), do: "不明"

  defp status_color(:CANCELLED), do: "bg-gray-200 text-gray-700"
  defp status_color(:PAYMENT_FAILED), do: "bg-red-100 text-red-800"
  defp status_color(:DELIVERED), do: "bg-green-100 text-green-800"
  defp status_color(:SHIPPED), do: "bg-blue-100 text-blue-800"
  defp status_color(_), do: "bg-yellow-100 text-yellow-800"

  defp format_amount(amount) when is_binary(amount) do
    case Integer.parse(amount) do
      {n, _} ->
        formatted =
          n
          |> Integer.to_string()
          |> String.reverse()
          |> String.replace(~r/(\d{3})(?=\d)/, "\\1,")
          |> String.reverse()

        "¥" <> formatted

      :error ->
        "¥" <> amount
    end
  end

  defp format_amount(_), do: "¥0"

  defp format_timestamp(%Google.Protobuf.Timestamp{seconds: seconds}) do
    # JST(+9:00 固定)で表示
    seconds
    |> Kernel.+(9 * 3600)
    |> DateTime.from_unix!()
    |> Calendar.strftime("%Y-%m-%d %H:%M")
  end

  defp format_timestamp(_), do: "-"

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <nav class="bg-white shadow-md">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16">
            <div class="flex items-center">
              <.link
                navigate="/dashboard"
                class="text-2xl font-bold text-gray-900 hover:text-gray-700"
              >
                オンラインショップモール
              </.link>
            </div>
            <div class="flex items-center space-x-4">
              <.link
                navigate="/products"
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                商品一覧
              </.link>
              <span class="text-gray-900 px-3 py-2 text-sm font-semibold">注文履歴</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-5xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-6">注文履歴</h1>

        <%= if @error do %>
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            {@error}
          </div>
        <% end %>

        <%= if @loading do %>
          <div class="text-center py-12">
            <div class="inline-block animate-spin rounded-full h-10 w-10 border-b-2 border-blue-600">
            </div>
          </div>
        <% else %>
          <%= if @orders == [] and is_nil(@error) do %>
            <div class="bg-white shadow rounded-lg p-12 text-center text-gray-500">
              注文はまだありません
              <div class="mt-4">
                <.link navigate="/products" class="text-blue-600 hover:text-blue-800 font-medium">
                  商品を見る →
                </.link>
              </div>
            </div>
          <% end %>

          <div class="space-y-4">
            <div :for={order <- @orders} class="bg-white shadow rounded-lg p-5">
              <div class="flex items-center justify-between flex-wrap gap-3">
                <div>
                  <div class="flex items-center space-x-3">
                    <span class="font-mono text-sm text-gray-500">{order.order_number}</span>
                    <span class={"inline-flex px-2 py-0.5 rounded-full text-xs font-medium " <>
                      status_color(order.status)}>
                      {status_label(order.status)}
                    </span>
                  </div>
                  <div class="text-sm text-gray-500 mt-1">
                    注文日時: {format_timestamp(order.created_at)}
                  </div>
                </div>
                <div class="text-right">
                  <div class="text-xl font-bold text-gray-900">
                    {format_amount(order.total_amount)}
                  </div>
                  <div class="text-xs text-gray-500">
                    (送料 {format_amount(order.shipping_fee)} 含む)
                  </div>
                </div>
                <div class="space-x-2">
                  <button
                    phx-click="show_detail"
                    phx-value-order-id={order.id}
                    class="px-3 py-1.5 text-sm font-medium text-blue-600 border border-blue-300 rounded-md hover:bg-blue-50"
                  >
                    詳細
                  </button>
                  <button
                    :if={order.status not in [:CANCELLED]}
                    phx-click="reorder"
                    phx-value-order-id={order.id}
                    class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
                  >
                    再注文
                  </button>
                  <button
                    :if={cancellable?(order)}
                    phx-click="open_cancel"
                    phx-value-order-id={order.id}
                    class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-300 rounded-md hover:bg-red-50"
                  >
                    キャンセル
                  </button>
                </div>
              </div>
            </div>
          </div>
        <% end %>
        
    <!-- 注文詳細モーダル(GetOrderDetail + GetOrderStatusHistory) -->
        <div
          :if={@detail}
          class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50"
        >
          <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-bold text-gray-900">注文詳細 {@detail.order_number}</h2>
              <button phx-click="close_detail" class="text-gray-400 hover:text-gray-600">✕</button>
            </div>
            <dl class="grid grid-cols-3 gap-y-2 text-sm mb-4">
              <dt class="text-gray-500">状態</dt>
              <dd class="col-span-2">
                <span class={"inline-flex px-2 py-0.5 rounded-full text-xs font-medium " <> status_color(@detail.status)}>
                  {status_label(@detail.status)}
                </span>
              </dd>
              <dt class="text-gray-500">合計</dt>
              <dd class="col-span-2 font-semibold">{format_amount(@detail.total_amount)}</dd>
              <dt class="text-gray-500">送料</dt>
              <dd class="col-span-2">{format_amount(@detail.shipping_fee)}</dd>
              <dt class="text-gray-500">注文日時</dt>
              <dd class="col-span-2">{format_timestamp(@detail.created_at)}</dd>
            </dl>
            <h3 class="text-sm font-semibold text-gray-700 mb-2">ステータス履歴</h3>
            <div class="flex items-center flex-wrap gap-1 text-xs">
              <%= for {st, idx} <- Enum.with_index(@history) do %>
                <span :if={idx > 0} class="text-gray-400">→</span>
                <span class={"inline-flex px-2 py-0.5 rounded-full font-medium " <> status_color(st)}>
                  {status_label(st)}
                </span>
              <% end %>
            </div>
          </div>
        </div>
        
    <!-- キャンセル確認モーダル -->
        <div
          :if={@cancel_target}
          class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50"
        >
          <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
            <h2 class="text-lg font-bold text-gray-900 mb-2">注文のキャンセル</h2>
            <p class="text-sm text-gray-600 mb-4">
              この注文をキャンセルします。お支払い済みの場合は全額返金されます。
            </p>
            <form phx-submit="execute_cancel">
              <label class="block text-sm font-medium text-gray-700 mb-1">キャンセル理由</label>
              <select
                name="reason"
                class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm mb-4"
              >
                <option :for={{value, label} <- @cancel_reasons} value={value}>{label}</option>
              </select>
              <div class="flex justify-end space-x-3">
                <button
                  type="button"
                  phx-click="close_cancel"
                  class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
                >
                  戻る
                </button>
                <button
                  type="submit"
                  class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700"
                >
                  キャンセルを確定
                </button>
              </div>
            </form>
          </div>
        </div>
      </main>
    </div>
    """
  end
end
