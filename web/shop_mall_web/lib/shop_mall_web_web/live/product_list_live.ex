defmodule ShopMallWebWeb.ProductListLive do
  use ShopMallWebWeb, :live_view

  alias ShopService.V1.{
    ShopService.Stub,
    ListProductsRequest,
    GetShopRequest
  }

  @impl true
  def mount(_params, _session, socket) do
    {:ok,
     socket
     |> assign(:products, [])
     |> assign(:shops, %{})
     |> assign(:loading, true)
     |> assign(:error, nil)
     |> load_products()}
  end

  defp load_products(socket) do
    # 全ショップの公開商品を取得
    request = %ListProductsRequest{
      shop_id: "",
      category: "",
      published_only: true,
      limit: 100,
      offset: 0
    }

    case call_shop_service(:list_products, request) do
      {:ok, response} ->
        # 各商品のショップ情報を取得
        shops = load_shops(response.products)

        socket
        |> assign(:products, response.products)
        |> assign(:shops, shops)
        |> assign(:loading, false)

      {:error, reason} ->
        socket
        |> assign(:error, "商品の読み込みに失敗しました: #{inspect(reason)}")
        |> assign(:loading, false)
    end
  end

  defp load_shops(products) do
    products
    |> Enum.map(& &1.shop_id)
    |> Enum.uniq()
    |> Enum.reduce(%{}, fn shop_id, acc ->
      case call_shop_service(:get_shop, %GetShopRequest{shop_id: shop_id}) do
        {:ok, response} -> Map.put(acc, shop_id, response.shop)
        {:error, _} -> acc
      end
    end)
  end

  defp call_shop_service(:list_products, request) do
    case get_shop_channel() do
      {:ok, channel} ->
        case Stub.list_products(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, "接続エラー: #{inspect(error)}"}
        end

      {:error, reason} ->
        {:error, "Shop Serviceに接続できません: #{inspect(reason)}"}
    end
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
        {:error, "Shop Serviceに接続できません: #{inspect(reason)}"}
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

  defp format_price(price) when is_binary(price) do
    case Decimal.parse(price) do
      {decimal, ""} ->
        decimal
        |> Decimal.round(0)
        |> Decimal.to_string()
        |> then(&"¥#{&1}")

      _ ->
        "¥#{price}"
    end
  end

  defp format_price(_), do: "¥0"

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <!-- ナビゲーションバー -->
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
                navigate="/dashboard"
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                ホーム
              </.link>
              <.link
                navigate="/products"
                class="text-blue-600 hover:text-blue-800 px-3 py-2 rounded-md text-sm font-medium border-b-2 border-blue-600"
              >
                商品一覧
              </.link>
              <.link
                href="/session/logout"
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                ログアウト
              </.link>
            </div>
          </div>
        </div>
      </nav>
      
    <!-- メインコンテンツ -->
      <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="px-4 py-6 sm:px-0">
          <h1 class="text-3xl font-bold text-gray-900 mb-6">商品一覧</h1>

          <%= if @loading do %>
            <div class="text-center py-12">
              <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600">
              </div>
              <p class="mt-4 text-gray-600">商品を読み込んでいます...</p>
            </div>
          <% end %>

          <%= if @error do %>
            <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
              {@error}
            </div>
          <% end %>

          <%= if !@loading and @products == [] do %>
            <div class="text-center py-12 bg-white rounded-lg shadow">
              <svg
                class="mx-auto h-12 w-12 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
                />
              </svg>
              <p class="mt-4 text-xl text-gray-600">商品がありません</p>
            </div>
          <% end %>
          
    <!-- 商品グリッド -->
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            <%= for product <- @products do %>
              <.link
                navigate={"/products/#{product.id}"}
                class="bg-white rounded-lg shadow hover:shadow-xl transition-shadow overflow-hidden"
              >
                <!-- 商品画像（プレースホルダー） -->
                <div class="h-48 bg-gradient-to-br from-blue-100 to-purple-100 flex items-center justify-center">
                  <svg
                    class="h-20 w-20 text-gray-400"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                    />
                  </svg>
                </div>
                
    <!-- 商品情報 -->
                <div class="p-4">
                  <h3 class="text-lg font-semibold text-gray-900 mb-2 line-clamp-2">
                    {product.name}
                  </h3>

                  <%= if shop = Map.get(@shops, product.shop_id) do %>
                    <p class="text-sm text-gray-500 mb-2">
                      {shop.name}
                    </p>
                  <% end %>

                  <p class="text-sm text-gray-600 mb-3 line-clamp-2">
                    {product.description}
                  </p>

                  <div class="flex items-center justify-between">
                    <span class="text-2xl font-bold text-blue-600">
                      {format_price(product.price)}
                    </span>

                    <%= if product.stock_quantity > 0 do %>
                      <span class="text-sm text-green-600 font-medium">
                        在庫あり
                      </span>
                    <% else %>
                      <span class="text-sm text-red-600 font-medium">
                        在庫切れ
                      </span>
                    <% end %>
                  </div>
                </div>
              </.link>
            <% end %>
          </div>
        </div>
      </main>
    </div>
    """
  end
end
