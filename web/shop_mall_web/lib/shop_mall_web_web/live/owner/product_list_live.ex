defmodule ShopMallWebWeb.Owner.ProductListLive do
  use ShopMallWebWeb, :live_view

  alias ShopService.V1.{
    ShopService.Stub,
    ListProductsRequest,
    GetShopRequest,
    DeleteProductRequest,
    ToggleProductPublishRequest
  }

  @impl true
  def mount(_params, _session, socket) do
    # TODO: セッションから実際のowner_idとshop_idを取得
    owner_id = "admin-user-id"
    shop_id = "11111111-1111-1111-1111-111111111111"

    {:ok,
     socket
     |> assign(:owner_id, owner_id)
     |> assign(:shop_id, shop_id)
     |> assign(:shop, nil)
     |> assign(:products, [])
     |> assign(:loading, true)
     |> assign(:error, nil)
     |> load_data()}
  end

  defp load_data(socket) do
    shop_id = socket.assigns.shop_id

    # ショップ情報を取得
    shop =
      case call_shop_service(:get_shop, %GetShopRequest{shop_id: shop_id}) do
        {:ok, response} -> response.shop
        {:error, _} -> nil
      end

    # 商品一覧を取得（公開・非公開両方）
    products =
      case call_shop_service(:list_products, %ListProductsRequest{
             shop_id: shop_id,
             category: "",
             published_only: false,
             limit: 1000,
             offset: 0
           }) do
        {:ok, response} -> response.products
        {:error, _} -> []
      end

    socket
    |> assign(:shop, shop)
    |> assign(:products, products)
    |> assign(:loading, false)
  end

  @impl true
  def handle_event("delete_product", %{"product-id" => product_id}, socket) do
    request = %DeleteProductRequest{product_id: product_id}

    case call_shop_service(:delete_product, request) do
      {:ok, _response} ->
        {:noreply,
         socket
         |> put_flash(:info, "商品を削除しました")
         |> assign(:loading, true)
         |> load_data()}

      {:error, reason} ->
        {:noreply,
         socket
         |> put_flash(:error, "削除に失敗しました: #{inspect(reason)}")}
    end
  end

  @impl true
  def handle_event("toggle_publish", %{"product-id" => product_id}, socket) do
    IO.puts("=== TOGGLE PUBLISH EVENT RECEIVED ===")
    IO.inspect(product_id, label: "Product ID")

    request = %ToggleProductPublishRequest{product_id: product_id}

    case call_shop_service(:toggle_product_publish, request) do
      {:ok, response} ->
        IO.puts("=== TOGGLE PUBLISH SUCCESS ===")
        IO.inspect(response, label: "Response")

        {:noreply,
         socket
         |> put_flash(:info, "公開状態を変更しました")
         |> assign(:loading, true)
         |> load_data()}

      {:error, reason} ->
        IO.puts("=== TOGGLE PUBLISH ERROR ===")
        IO.inspect(reason, label: "Error")

        {:noreply,
         socket
         |> put_flash(:error, "変更に失敗しました: #{inspect(reason)}")}
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

  defp call_shop_service(:list_products, request) do
    case get_shop_channel() do
      {:ok, channel} ->
        case Stub.list_products(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, error}
        end

      {:error, reason} ->
        {:error, "Shop Serviceに接続できません: #{inspect(reason)}"}
    end
  end

  defp call_shop_service(:delete_product, request) do
    case get_shop_channel() do
      {:ok, channel} ->
        case Stub.delete_product(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, error}
        end

      {:error, reason} ->
        {:error, "Shop Serviceに接続できません: #{inspect(reason)}"}
    end
  end

  defp call_shop_service(:toggle_product_publish, request) do
    case get_shop_channel() do
      {:ok, channel} ->
        case Stub.toggle_product_publish(channel, request) do
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
                class="text-white hover:text-gray-200 px-3 py-2 rounded-md text-sm font-medium"
              >
                ダッシュボード
              </.link>
              <.link
                navigate="/owner/products"
                class="text-white hover:text-gray-200 px-3 py-2 rounded-md text-sm font-medium border-b-2 border-white"
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
            <!-- ヘッダー -->
            <div class="mb-6 flex justify-between items-center">
              <div>
                <h1 class="text-3xl font-bold text-gray-900">商品管理</h1>
                <%= if @shop do %>
                  <p class="mt-1 text-sm text-gray-600">{@shop.name}</p>
                <% end %>
              </div>
              <.link
                navigate="/owner/products/new"
                class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
              >
                <svg
                  class="-ml-1 mr-2 h-5 w-5"
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
                新規商品登録
              </.link>
            </div>
            
    <!-- 商品一覧 -->
            <%= if Enum.empty?(@products) do %>
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
                    d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
                  />
                </svg>
                <h3 class="mt-2 text-sm font-medium text-gray-900">商品がありません</h3>
                <p class="mt-1 text-sm text-gray-500">新規商品を登録してください。</p>
                <div class="mt-6">
                  <.link
                    navigate="/owner/products/new"
                    class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
                  >
                    <svg
                      class="-ml-1 mr-2 h-5 w-5"
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
                    新規商品登録
                  </.link>
                </div>
              </div>
            <% else %>
              <div class="bg-white shadow overflow-hidden rounded-lg">
                <ul class="divide-y divide-gray-200">
                  <%= for product <- @products do %>
                    <li class="p-6 hover:bg-gray-50">
                      <div class="flex items-center space-x-4">
                        <!-- 商品画像 -->
                        <div class="flex-shrink-0">
                          <div class="h-20 w-20 bg-gray-200 rounded-md flex items-center justify-center">
                            <svg
                              class="h-10 w-10 text-gray-400"
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
                        </div>
                        
    <!-- 商品情報 -->
                        <div class="flex-1 min-w-0">
                          <div class="flex items-center space-x-3">
                            <p class="text-lg font-medium text-gray-900 truncate">
                              {product.name}
                            </p>
                            <%= if product.published do %>
                              <span class="px-2 py-1 text-xs font-medium rounded-full bg-green-100 text-green-800">
                                公開中
                              </span>
                            <% else %>
                              <span class="px-2 py-1 text-xs font-medium rounded-full bg-gray-100 text-gray-800">
                                非公開
                              </span>
                            <% end %>
                          </div>
                          <p class="mt-1 text-sm text-gray-500 line-clamp-2">
                            {product.description}
                          </p>
                          <div class="mt-2 flex items-center space-x-4 text-sm text-gray-500">
                            <span class="font-medium text-gray-900">
                              {format_price(product.price)}
                            </span>
                            <span>在庫: {product.stock_quantity}個</span>
                            <span>カテゴリ: {product.category}</span>
                          </div>
                        </div>
                        
    <!-- アクション -->
                        <div class="flex-shrink-0 flex space-x-2">
                          <.link
                            navigate={"/owner/products/#{product.id}/edit"}
                            class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                          >
                            編集
                          </.link>
                          <button
                            phx-click="toggle_publish"
                            phx-value-product-id={product.id}
                            class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                          >
                            {if product.published, do: "非公開", else: "公開"}
                          </button>
                          <button
                            phx-click="delete_product"
                            phx-value-product-id={product.id}
                            data-confirm="本当に削除しますか？"
                            class="inline-flex items-center px-3 py-2 border border-red-300 shadow-sm text-sm font-medium rounded-md text-red-700 bg-white hover:bg-red-50"
                          >
                            削除
                          </button>
                        </div>
                      </div>
                    </li>
                  <% end %>
                </ul>
              </div>
            <% end %>
          <% end %>
        </div>
      </main>
    </div>
    """
  end
end
