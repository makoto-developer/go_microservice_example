defmodule ShopMallWebWeb.Owner.SalesReportLive do
  @moduledoc """
  加盟店の売上レポート画面。
  期間を指定して日別売上(GetSalesReport)を表示し、CSV エクスポート
  (ExportSalesData)のダウンロード URL を発行する。
  """
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.ShopServiceClient, as: Shops

  @impl true
  def mount(_params, session, socket) do
    owner_id = session["user_id"] || "admin-user-id"
    today = Date.utc_today()
    from = Date.add(today, -30)

    {:ok,
     socket
     |> assign(:owner_id, owner_id)
     |> assign(:shop, nil)
     |> assign(:report_type, "daily")
     |> assign(:date_from, Date.to_iso8601(from))
     |> assign(:date_to, Date.to_iso8601(today))
     |> assign(:rows, [])
     |> assign(:summary, nil)
     |> assign(:csv_url, nil)
     |> assign(:loading, true)
     |> assign(:error, nil)
     |> load_shop_and_report()}
  end

  defp load_shop_and_report(socket) do
    case Shops.get_shops_by_owner(socket.assigns.owner_id) do
      {:ok, %{shops: [shop | _]}} ->
        socket |> assign(:shop, shop) |> load_report()

      {:ok, _} ->
        socket |> assign(:loading, false) |> assign(:error, "店舗がありません")

      {:error, reason} ->
        socket |> assign(:loading, false) |> assign(:error, "店舗の取得に失敗しました: #{reason}")
    end
  end

  defp load_report(socket) do
    %{shop: shop, report_type: type, date_from: from, date_to: to} = socket.assigns

    case Shops.get_sales_report(shop.id, type, from, to) do
      {:ok, response} ->
        socket
        |> assign(:rows, response.report_data)
        |> assign(:summary, response.summary)
        |> assign(:loading, false)
        |> assign(:error, nil)

      {:error, reason} ->
        socket
        |> assign(:loading, false)
        |> assign(:error, "売上レポートの取得に失敗しました: #{reason}")
    end
  end

  @impl true
  def handle_event(
        "set_range",
        %{"date_from" => from, "date_to" => to, "report_type" => type},
        socket
      ) do
    {:noreply,
     socket
     |> assign(:date_from, from)
     |> assign(:date_to, to)
     |> assign(:report_type, type)
     |> assign(:loading, true)
     |> load_report()}
  end

  @impl true
  def handle_event("export_csv", _params, socket) do
    %{shop: shop, date_from: from, date_to: to} = socket.assigns

    case Shops.export_sales_data(shop.id, from, to) do
      {:ok, response} ->
        {:noreply,
         socket
         |> assign(:csv_url, response.csv_url)
         |> put_flash(:info, "CSV を書き出しました")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "CSV エクスポートに失敗しました: #{reason}")}
    end
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
                navigate="/owner/orders"
                class="text-emerald-200 hover:text-white text-sm font-medium"
              >
                受注管理
              </.link>
              <span class="text-white text-sm font-medium">売上レポート</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-5xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between mb-6">
          <h1 class="text-2xl font-bold text-gray-900">
            売上レポート <span :if={@shop} class="text-sm font-normal text-gray-500">({@shop.name})</span>
          </h1>
          <button
            phx-click="export_csv"
            class="px-4 py-2 text-sm font-medium text-emerald-700 border border-emerald-300 rounded-md hover:bg-emerald-50"
          >
            ⬇ CSV エクスポート
          </button>
        </div>

        <%= if @csv_url do %>
          <div class="bg-emerald-50 border border-emerald-300 text-emerald-800 px-4 py-3 rounded mb-4 text-sm">
            ダウンロード: <a href={@csv_url} class="underline font-mono">{@csv_url}</a>
          </div>
        <% end %>

        <form
          phx-submit="set_range"
          class="bg-white shadow rounded-lg p-4 mb-6 flex items-end space-x-3"
        >
          <div>
            <label class="block text-xs font-medium text-gray-500 mb-1">開始日</label>
            <input
              type="date"
              name="date_from"
              value={@date_from}
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-500 mb-1">終了日</label>
            <input
              type="date"
              name="date_to"
              value={@date_to}
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-500 mb-1">集計単位</label>
            <select name="report_type" class="border border-gray-300 rounded-md px-2 py-1.5 text-sm">
              <option value="daily" selected={@report_type == "daily"}>日別</option>
              <option value="monthly" selected={@report_type == "monthly"}>月別</option>
            </select>
          </div>
          <button
            type="submit"
            class="px-4 py-2 text-sm font-medium text-white bg-emerald-600 rounded-md hover:bg-emerald-700"
          >
            表示
          </button>
        </form>

        <%= if @error do %>
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            {@error}
          </div>
        <% end %>

        <%= if @summary do %>
          <div class="grid grid-cols-3 gap-4 mb-6">
            <div class="bg-white shadow rounded-lg p-4 text-center">
              <div class="text-xs text-gray-500">売上合計</div>
              <div class="text-2xl font-bold text-gray-900">
                {format_amount(@summary.total_sales)}
              </div>
            </div>
            <div class="bg-white shadow rounded-lg p-4 text-center">
              <div class="text-xs text-gray-500">注文数</div>
              <div class="text-2xl font-bold text-gray-900">{@summary.total_orders}</div>
            </div>
            <div class="bg-white shadow rounded-lg p-4 text-center">
              <div class="text-xs text-gray-500">平均注文額</div>
              <div class="text-2xl font-bold text-gray-900">
                {format_amount(@summary.average_order_value)}
              </div>
            </div>
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
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">日付</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">売上</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    注文数
                  </th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    平均注文額
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <%= if @rows == [] do %>
                  <tr>
                    <td colspan="4" class="px-4 py-8 text-center text-gray-500">この期間の売上データはありません</td>
                  </tr>
                <% end %>
                <tr :for={row <- @rows}>
                  <td class="px-4 py-2.5 text-sm text-gray-900">{row.date}</td>
                  <td class="px-4 py-2.5 text-sm text-right font-semibold text-gray-900">
                    {format_amount(row.total_sales)}
                  </td>
                  <td class="px-4 py-2.5 text-sm text-right text-gray-600">{row.order_count}</td>
                  <td class="px-4 py-2.5 text-sm text-right text-gray-600">
                    {format_amount(row.average_order_value)}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        <% end %>
      </main>
    </div>
    """
  end
end
