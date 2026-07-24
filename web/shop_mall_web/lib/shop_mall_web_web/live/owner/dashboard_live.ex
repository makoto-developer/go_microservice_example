defmodule ShopMallWebWeb.Owner.DashboardLive do
  use ShopMallWebWeb, :live_view

  alias ShopService.V1.{
    ShopService.Stub,
    GetShopRequest,
    ListProductsRequest,
    ListOrdersRequest
  }

  @impl true
  def mount(_params, _session, socket) do
    # TODO: セッションから実際のowner_idを取得
    # 今はテストデータのショップオーナーIDを使用
    owner_id = "admin-user-id"
    shop_id = "11111111-1111-1111-1111-111111111111"

    {:ok,
     socket
     |> assign(:owner_id, owner_id)
     |> assign(:shop_id, shop_id)
     |> assign(:shop, nil)
     |> assign(:product_count, 0)
     |> assign(:order_count, 0)
     |> assign(:loading, true)
     |> load_dashboard_data()}
  end

  defp load_dashboard_data(socket) do
    shop_id = socket.assigns.shop_id

    # ショップ情報を取得
    shop =
      case call_shop_service(:get_shop, %GetShopRequest{shop_id: shop_id}) do
        {:ok, response} -> response.shop
        {:error, _} -> nil
      end

    # 商品数を取得
    product_count =
      case call_shop_service(:list_products, %ListProductsRequest{
             shop_id: shop_id,
             category: "",
             published_only: false,
             limit: 1000,
             offset: 0
           }) do
        {:ok, response} ->
          IO.inspect(response.products, label: "Products from gRPC")
          length(response.products)

        {:error, reason} ->
          IO.inspect(reason, label: "Error fetching products")
          0
      end

    # 注文数を取得
    order_count =
      case call_shop_service(:list_orders, %ListOrdersRequest{
             shop_id: shop_id,
             status: :ORDER_STATUS_UNSPECIFIED,
             date_from: "",
             date_to: "",
             customer_name: "",
             product_name: "",
             sort_by: "",
             sort_order: ""
           }) do
        {:ok, response} -> response.total_count
        {:error, _} -> 0
      end

    socket
    |> assign(:shop, shop)
    |> assign(:product_count, product_count)
    |> assign(:order_count, order_count)
    |> assign(:loading, false)
  end

  defp call_shop_service(:get_shop, request) do
    case get_shop_channel() do
      {:ok, channel} ->
        case Stub.get_shop(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, error}
        end

      {:error, reason} ->
        {:error, "接続エラー: #{inspect(reason)}"}
    end
  end

  defp call_shop_service(:list_products, request) do
    case get_shop_channel() do
      {:ok, channel} ->
        case Stub.list_products(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, error}
        end

      {:error, reason} ->
        {:error, "接続エラー: #{inspect(reason)}"}
    end
  end

  defp call_shop_service(:list_orders, request) do
    case get_shop_channel() do
      {:ok, channel} ->
        case Stub.list_orders(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, error}
        end

      {:error, reason} ->
        {:error, "接続エラー: #{inspect(reason)}"}
    end
  end

  defp get_shop_channel do
    host = System.get_env("SHOP_SERVICE_HOST", "localhost")
    port = String.to_integer(System.get_env("SHOP_SERVICE_PORT", "22101"))

    case GRPC.Stub.connect("#{host}:#{port}") do
      {:ok, channel} -> {:ok, channel}
      {:error, reason} -> {:error, reason}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <!-- ナビゲーションバー -->
      <nav class="bg-blue-800 shadow-md">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16">
            <div class="flex items-center">
              <.link
                navigate="/owner/dashboard"
                class="text-2xl font-bold text-white hover:text-gray-200"
              >
                ショップ管理
              </.link>
            </div>
            <div class="flex items-center space-x-4">
              <.link
                navigate="/owner/dashboard"
                class="text-white hover:text-gray-200 px-3 py-2 rounded-md text-sm font-medium border-b-2 border-white"
              >
                ダッシュボード
              </.link>
              <.link
                navigate="/owner/products"
                class="text-white hover:text-gray-200 px-3 py-2 rounded-md text-sm font-medium"
              >
                商品管理
              </.link>
              <.link
                navigate="/dashboard"
                class="text-white hover:text-gray-200 px-3 py-2 rounded-md text-sm font-medium"
              >
                顧客画面へ
              </.link>
            </div>
          </div>
        </div>
      </nav>
      
    <!-- メインコンテンツ -->
      <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="px-4 py-6 sm:px-0">
          <%= if @loading do %>
            <div class="text-center py-12">
              <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600">
              </div>
              <p class="mt-4 text-gray-600">読み込み中...</p>
            </div>
          <% else %>
            <!-- ショップ情報 -->
            <%= if @shop do %>
              <div class="bg-white overflow-hidden shadow rounded-lg mb-6">
                <div class="px-4 py-5 sm:p-6">
                  <h2 class="text-2xl font-bold text-gray-900 mb-2">
                    {@shop.name}
                  </h2>
                  <p class="text-gray-600">
                    {@shop.description}
                  </p>
                  <div class="mt-4 flex items-center space-x-4">
                    <span class={"px-3 py-1 rounded-full text-sm font-medium " <>
                      case @shop.status do
                        :APPROVED -> "bg-green-100 text-green-800"
                        :PENDING_APPROVAL -> "bg-yellow-100 text-yellow-800"
                        :REJECTED -> "bg-red-100 text-red-800"
                        :SUSPENDED -> "bg-gray-100 text-gray-800"
                        _ -> "bg-gray-100 text-gray-800"
                      end}>
                      {case @shop.status do
                        :APPROVED -> "承認済み"
                        :PENDING_APPROVAL -> "承認待ち"
                        :REJECTED -> "却下"
                        :SUSPENDED -> "停止中"
                        _ -> "不明"
                      end}
                    </span>
                    <%= if @shop.published do %>
                      <span class="px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
                        公開中
                      </span>
                    <% else %>
                      <span class="px-3 py-1 rounded-full text-sm font-medium bg-gray-100 text-gray-800">
                        非公開
                      </span>
                    <% end %>
                  </div>
                </div>
              </div>
            <% end %>
            
    <!-- 統計カード -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
              <!-- 商品数 -->
              <div class="bg-white overflow-hidden shadow rounded-lg">
                <div class="px-4 py-5 sm:p-6">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 bg-blue-500 rounded-md p-3">
                      <svg
                        class="h-6 w-6 text-white"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
                        />
                      </svg>
                    </div>
                    <div class="ml-5 w-0 flex-1">
                      <dl>
                        <dt class="text-sm font-medium text-gray-500 truncate">
                          登録商品数
                        </dt>
                        <dd class="text-3xl font-semibold text-gray-900">
                          {@product_count}
                        </dd>
                      </dl>
                    </div>
                  </div>
                  <div class="mt-5">
                    <.link
                      navigate="/owner/products"
                      class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
                    >
                      商品管理 →
                    </.link>
                  </div>
                </div>
              </div>
              
    <!-- 注文数 -->
              <div class="bg-white overflow-hidden shadow rounded-lg">
                <div class="px-4 py-5 sm:p-6">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 bg-green-500 rounded-md p-3">
                      <svg
                        class="h-6 w-6 text-white"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                        />
                      </svg>
                    </div>
                    <div class="ml-5 w-0 flex-1">
                      <dl>
                        <dt class="text-sm font-medium text-gray-500 truncate">
                          受注件数
                        </dt>
                        <dd class="text-3xl font-semibold text-gray-900">
                          {@order_count}
                        </dd>
                      </dl>
                    </div>
                  </div>
                  <div class="mt-5">
                    <button
                      disabled
                      class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-400 cursor-not-allowed"
                    >
                      準備中
                    </button>
                  </div>
                </div>
              </div>
              
    <!-- 売上 -->
              <div class="bg-white overflow-hidden shadow rounded-lg">
                <div class="px-4 py-5 sm:p-6">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 bg-purple-500 rounded-md p-3">
                      <svg
                        class="h-6 w-6 text-white"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                      </svg>
                    </div>
                    <div class="ml-5 w-0 flex-1">
                      <dl>
                        <dt class="text-sm font-medium text-gray-500 truncate">
                          今月の売上
                        </dt>
                        <dd class="text-3xl font-semibold text-gray-900">
                          準備中
                        </dd>
                      </dl>
                    </div>
                  </div>
                  <div class="mt-5">
                    <button
                      disabled
                      class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-400 cursor-not-allowed"
                    >
                      準備中
                    </button>
                  </div>
                </div>
              </div>
            </div>
            
    <!-- クイックアクション -->
            <div class="bg-white shadow rounded-lg">
              <div class="px-4 py-5 sm:p-6">
                <h3 class="text-lg font-medium text-gray-900 mb-4">クイックアクション</h3>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <.link
                    navigate="/owner/products/new"
                    class="flex items-center p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors"
                  >
                    <svg
                      class="h-8 w-8 text-gray-400"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 4v16m8-8H4"
                      />
                    </svg>
                    <span class="ml-3 text-gray-700 font-medium">新しい商品を登録</span>
                  </.link>

                  <.link
                    navigate="/owner/products"
                    class="flex items-center p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors"
                  >
                    <svg
                      class="h-8 w-8 text-gray-400"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                      />
                    </svg>
                    <span class="ml-3 text-gray-700 font-medium">商品を編集</span>
                  </.link>
                </div>
              </div>
            </div>
          <% end %>
        </div>
      </main>
    </div>
    """
  end
end
