defmodule ShopMallWebWeb.Admin.OrderAnalyticsLive do
  @moduledoc """
  管理者向けの注文分析画面。
  注文統計(GetOrderStatistics)・売れ筋ランキング(GetProductSalesRanking)・
  注文検索(SearchOrders)・CSV エクスポート(ExportOrdersToCSV)を提供する。
  """
  use ShopMallWebWeb, :live_view

  alias OrderService.V1.{
    ExportOrdersToCSVRequest,
    GetOrderStatisticsRequest,
    GetProductSalesRankingRequest,
    OrderService.Stub,
    SearchOrdersRequest
  }

  @impl true
  def mount(_params, _session, socket) do
    today = Date.utc_today()
    from = Date.add(today, -30)

    {:ok,
     socket
     |> assign(:date_from, Date.to_iso8601(from))
     |> assign(:date_to, Date.to_iso8601(today))
     |> assign(:stats, nil)
     |> assign(:rankings, [])
     |> assign(:search_results, nil)
     |> assign(:csv_url, nil)
     |> assign(:error, nil)
     |> load_analytics()}
  end

  defp load_analytics(socket) do
    %{date_from: from, date_to: to} = socket.assigns

    with {:ok, channel} <- connect(),
         {:ok, stats} <-
           Stub.get_order_statistics(channel, %GetOrderStatisticsRequest{
             shop_id: "",
             date_from: from,
             date_to: to,
             group_by: "day"
           }),
         {:ok, ranking} <-
           Stub.get_product_sales_ranking(channel, %GetProductSalesRankingRequest{
             shop_id: "",
             date_from: from,
             date_to: to,
             limit: 10
           }) do
      socket
      |> assign(:stats, stats)
      |> assign(:rankings, ranking.rankings)
      |> assign(:error, nil)
    else
      {:error, reason} ->
        assign(socket, :error, "分析データの取得に失敗しました: #{inspect(reason)}")
    end
  end

  @impl true
  def handle_event("set_range", %{"date_from" => from, "date_to" => to}, socket) do
    {:noreply,
     socket
     |> assign(:date_from, from)
     |> assign(:date_to, to)
     |> load_analytics()}
  end

  @impl true
  def handle_event("search", %{"query" => query}, socket) do
    request = %SearchOrdersRequest{
      order_number: query,
      customer_name: "",
      product_name: "",
      date_from: socket.assigns.date_from,
      date_to: socket.assigns.date_to,
      page: 1,
      page_size: 20
    }

    with {:ok, channel} <- connect(),
         {:ok, response} <- Stub.search_orders(channel, request) do
      {:noreply, assign(socket, :search_results, response)}
    else
      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "検索に失敗しました: #{inspect(reason)}")}
    end
  end

  @impl true
  def handle_event("export_csv", _params, socket) do
    request = %ExportOrdersToCSVRequest{
      shop_id: "",
      date_from: socket.assigns.date_from,
      date_to: socket.assigns.date_to
    }

    with {:ok, channel} <- connect(),
         {:ok, response} <- Stub.export_orders_to_csv(channel, request) do
      note = if response.csv_url == "", do: response.message, else: response.csv_url
      {:noreply, socket |> assign(:csv_url, note) |> put_flash(:info, "エクスポートを実行しました")}
    else
      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "エクスポートに失敗しました: #{inspect(reason)}")}
    end
  end

  defp connect do
    host = System.get_env("ORDER_SERVICE_HOST", "localhost")
    port = System.get_env("ORDER_SERVICE_PORT", "50055")
    GRPC.Stub.connect("#{host}:#{port}")
  end

  defp format_yen(n) when is_integer(n) do
    "¥" <>
      (n
       |> Integer.to_string()
       |> String.reverse()
       |> String.replace(~r/(\d{3})(?=\d)/, "\\1,")
       |> String.reverse())
  end

  defp format_yen(_), do: "¥0"

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <nav class="bg-gray-900 shadow-md">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16 items-center">
            <span class="text-xl font-bold text-white">🛠 管理者コンソール</span>
            <div class="flex items-center space-x-4">
              <.link navigate="/admin/payments" class="text-gray-400 hover:text-white text-sm">
                決済管理
              </.link>
              <span class="text-gray-300 text-sm font-medium">注文分析</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-5xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between mb-6">
          <h1 class="text-2xl font-bold text-gray-900">注文分析</h1>
          <form phx-submit="set_range" class="flex items-end space-x-2">
            <input
              type="date"
              name="date_from"
              value={@date_from}
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <span class="text-gray-400">〜</span>
            <input
              type="date"
              name="date_to"
              value={@date_to}
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <button
              type="submit"
              class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
            >
              集計
            </button>
          </form>
        </div>

        <%= if @error do %>
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            {@error}
          </div>
        <% end %>

        <%= if @stats do %>
          <div class="grid grid-cols-4 gap-4 mb-6">
            <div class="bg-white shadow rounded-lg p-4 text-center">
              <div class="text-xs text-gray-500">注文数</div>
              <div class="text-2xl font-bold text-gray-900">{@stats.total_orders}</div>
            </div>
            <div class="bg-white shadow rounded-lg p-4 text-center">
              <div class="text-xs text-gray-500">売上合計</div>
              <div class="text-2xl font-bold text-gray-900">{format_yen(@stats.total_sales)}</div>
            </div>
            <div class="bg-white shadow rounded-lg p-4 text-center">
              <div class="text-xs text-gray-500">処理中</div>
              <div class="text-2xl font-bold text-yellow-600">{@stats.pending_orders}</div>
            </div>
            <div class="bg-white shadow rounded-lg p-4 text-center">
              <div class="text-xs text-gray-500">完了</div>
              <div class="text-2xl font-bold text-green-600">{@stats.completed_orders}</div>
            </div>
          </div>
        <% end %>

        <div class="grid grid-cols-2 gap-6 mb-6">
          <div class="bg-white shadow rounded-lg p-4">
            <div class="flex items-center justify-between mb-3">
              <h2 class="text-sm font-semibold text-gray-700">売れ筋ランキング</h2>
            </div>
            <%= if @rankings == [] do %>
              <div class="text-sm text-gray-400 py-4 text-center">データがありません</div>
            <% end %>
            <ol class="space-y-1.5">
              <li :for={{r, idx} <- Enum.with_index(@rankings)} class="flex items-center text-sm">
                <span class="w-6 text-gray-400 font-mono">{idx + 1}.</span>
                <span class="flex-1 text-gray-900">{r.product_name}</span>
                <span class="text-gray-500 mr-3">×{r.total_sold}</span>
                <span class="font-medium">{format_yen(r.total_revenue)}</span>
              </li>
            </ol>
          </div>

          <div class="bg-white shadow rounded-lg p-4">
            <div class="flex items-center justify-between mb-3">
              <h2 class="text-sm font-semibold text-gray-700">CSV エクスポート</h2>
              <button
                phx-click="export_csv"
                class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
              >
                ⬇ 注文一覧を書き出す
              </button>
            </div>
            <div :if={@csv_url} class="text-xs font-mono text-gray-600 break-all">{@csv_url}</div>
          </div>
        </div>

        <div class="bg-white shadow rounded-lg p-4">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">注文検索</h2>
          <form phx-submit="search" class="flex space-x-2 mb-4">
            <input
              type="text"
              name="query"
              placeholder="注文番号(例: ORD-1234567890)"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <button
              type="submit"
              class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
            >
              検索
            </button>
          </form>
          <%= if @search_results do %>
            <div class="text-xs text-gray-500 mb-2">{@search_results.total_count} 件</div>
            <table class="min-w-full text-sm">
              <tbody class="divide-y divide-gray-100">
                <tr :for={order <- @search_results.orders}>
                  <td class="py-1.5 font-mono">{order.order_number}</td>
                  <td class="py-1.5 text-gray-500">{order.status}</td>
                  <td class="py-1.5 text-right font-medium">¥{order.total_amount}</td>
                </tr>
              </tbody>
            </table>
          <% end %>
        </div>
      </main>
    </div>
    """
  end
end
