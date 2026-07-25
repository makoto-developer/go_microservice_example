defmodule ShopMallWebWeb.Admin.SearchAdminLive do
  @moduledoc """
  管理者の検索管理画面。
  検索分析(GetSearchAnalytics)と検索インデックスの管理
  (IndexProduct / UpdateProductIndex / DeleteProductIndex /
  IndexShop / DeleteShopIndex)を行う。
  """
  use ShopMallWebWeb, :live_view

  alias SearchService.V1.SearchService.Stub
  alias SearchService.V1, as: PB

  @impl true
  def mount(_params, _session, socket) do
    today = Date.utc_today()

    {:ok,
     socket
     |> assign(:date_from, Date.to_iso8601(Date.add(today, -30)))
     |> assign(:date_to, Date.to_iso8601(today))
     |> assign(:analytics_note, nil)}
  end

  @impl true
  def handle_event("load_analytics", %{"date_from" => from, "date_to" => to}, socket) do
    case call(fn ch ->
           Stub.get_search_analytics(ch, %PB.GetSearchAnalyticsRequest{
             date_from: from,
             date_to: to,
             report_type: "keyword"
           })
         end) do
      {:ok, resp} ->
        {:noreply,
         socket
         |> assign(:date_from, from)
         |> assign(:date_to, to)
         |> assign(:analytics_note, resp.message)}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "分析の取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("index_product", params, socket) do
    request = %PB.IndexProductRequest{
      product_id: params["product_id"],
      product_name: params["product_name"],
      shop_id: params["shop_id"] || "",
      price: params["price"] || "0"
    }

    run(socket, call(fn ch -> Stub.index_product(ch, request) end))
  end

  @impl true
  def handle_event("update_product_index", params, socket) do
    request = %PB.UpdateProductIndexRequest{
      product_id: params["product_id"],
      product_name: params["product_name"],
      price: params["price"] || "0"
    }

    run(socket, call(fn ch -> Stub.update_product_index(ch, request) end))
  end

  @impl true
  def handle_event("delete_product_index", %{"product_id" => product_id}, socket) do
    run(
      socket,
      call(fn ch ->
        Stub.delete_product_index(ch, %PB.DeleteProductIndexRequest{product_id: product_id})
      end)
    )
  end

  @impl true
  def handle_event("index_shop", params, socket) do
    request = %PB.IndexShopRequest{shop_id: params["shop_id"], shop_name: params["shop_name"]}
    run(socket, call(fn ch -> Stub.index_shop(ch, request) end))
  end

  @impl true
  def handle_event("delete_shop_index", %{"shop_id" => shop_id}, socket) do
    run(
      socket,
      call(fn ch -> Stub.delete_shop_index(ch, %PB.DeleteShopIndexRequest{shop_id: shop_id}) end)
    )
  end

  defp run(socket, result) do
    case result do
      {:ok, resp} -> {:noreply, put_flash(socket, :info, resp.message || "実行しました")}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "操作に失敗しました: #{reason}")}
    end
  end

  defp call(fun) do
    host = System.get_env("SEARCH_SERVICE_HOST", "localhost")
    port = System.get_env("SEARCH_SERVICE_PORT", "20110")

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
        {:error, "検索サービスに接続できません: #{inspect(reason)}"}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <nav class="bg-gray-900 shadow-md">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16 items-center">
            <span class="text-xl font-bold text-white">🛠 管理者コンソール</span>
            <div class="flex items-center space-x-4">
              <.link navigate="/admin/orders" class="text-gray-400 hover:text-white text-sm">
                注文分析
              </.link>
              <.link navigate="/admin/reviews" class="text-gray-400 hover:text-white text-sm">
                レビュー管理
              </.link>
              <span class="text-gray-300 text-sm font-medium">検索管理</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-6">検索管理</h1>

        <div class="bg-white shadow rounded-lg p-6 mb-6">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">検索分析</h2>
          <form phx-submit="load_analytics" class="flex items-end space-x-2 mb-3">
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
          <div :if={@analytics_note} class="text-sm text-gray-600">{@analytics_note}</div>
        </div>

        <div class="bg-white shadow rounded-lg p-6 mb-6">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">商品インデックス</h2>
          <form phx-submit="index_product" class="grid grid-cols-3 gap-2 mb-3">
            <input
              type="text"
              name="product_id"
              required
              placeholder="商品ID"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <input
              type="text"
              name="product_name"
              required
              placeholder="商品名"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <input
              type="text"
              name="price"
              placeholder="価格"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <input
              type="text"
              name="shop_id"
              placeholder="店舗ID"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm col-span-2"
            />
            <button
              type="submit"
              class="px-3 py-1.5 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
            >
              登録
            </button>
          </form>
          <form phx-submit="update_product_index" class="grid grid-cols-3 gap-2 mb-3">
            <input
              type="text"
              name="product_id"
              required
              placeholder="商品ID"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <input
              type="text"
              name="product_name"
              required
              placeholder="新しい商品名"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <input
              type="text"
              name="price"
              placeholder="新価格"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <div class="col-span-2"></div>
            <button
              type="submit"
              class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              更新
            </button>
          </form>
          <form phx-submit="delete_product_index" class="flex items-end space-x-2">
            <input
              type="text"
              name="product_id"
              required
              placeholder="商品ID"
              class="flex-1 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <button
              type="submit"
              class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-300 rounded-md hover:bg-red-50"
            >
              インデックス削除
            </button>
          </form>
        </div>

        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">店舗インデックス</h2>
          <form phx-submit="index_shop" class="grid grid-cols-3 gap-2 mb-3">
            <input
              type="text"
              name="shop_id"
              required
              placeholder="店舗ID"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <input
              type="text"
              name="shop_name"
              required
              placeholder="店舗名"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <button
              type="submit"
              class="px-3 py-1.5 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
            >
              登録
            </button>
          </form>
          <form phx-submit="delete_shop_index" class="flex items-end space-x-2">
            <input
              type="text"
              name="shop_id"
              required
              placeholder="店舗ID"
              class="flex-1 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <button
              type="submit"
              class="px-3 py-1.5 text-sm font-medium text-red-600 border border-red-300 rounded-md hover:bg-red-50"
            >
              インデックス削除
            </button>
          </form>
        </div>
      </main>
    </div>
    """
  end
end
