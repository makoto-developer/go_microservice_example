defmodule ShopMallWebWeb.ProductDetailLive do
  use ShopMallWebWeb, :live_view

  alias ShopService.V1.{
    ShopService.Stub,
    GetProductRequest,
    GetShopRequest
  }

  @impl true
  def mount(%{"id" => product_id}, _session, socket) do
    {:ok,
     socket
     |> assign(:product_id, product_id)
     |> assign(:product, nil)
     |> assign(:shop, nil)
     |> assign(:loading, true)
     |> assign(:error, nil)
     |> assign(:quantity, 1)
     |> load_product()}
  end

  defp load_product(socket) do
    product_id = socket.assigns.product_id
    request = %GetProductRequest{product_id: product_id}

    case call_shop_service(:get_product, request) do
      {:ok, response} when not is_nil(response.product) ->
        shop = load_shop(response.product.shop_id)

        socket
        |> assign(:product, response.product)
        |> assign(:shop, shop)
        |> assign(:loading, false)

      {:ok, _response} ->
        socket
        |> assign(:loading, false)
        |> assign(:error, "商品が見つかりませんでした")

      {:error, reason} ->
        socket
        |> assign(:error, "商品の読み込みに失敗しました: #{inspect(reason)}")
        |> assign(:loading, false)
    end
  end

  defp load_shop(shop_id) do
    case call_shop_service(:get_shop, %GetShopRequest{shop_id: shop_id}) do
      {:ok, response} -> response.shop
      {:error, _} -> nil
    end
  end

  defp call_shop_service(:get_product, request) do
    channel = get_shop_channel()

    case Stub.get_product(channel, request) do
      {:ok, response} -> {:ok, response}
      {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
      error -> {:error, error}
    end
  end

  defp call_shop_service(:get_shop, request) do
    channel = get_shop_channel()

    case Stub.get_shop(channel, request) do
      {:ok, response} -> {:ok, response}
      {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
      error -> {:error, error}
    end
  end

  defp get_shop_channel do
    host = System.get_env("SHOP_SERVICE_HOST", "localhost")
    port = String.to_integer(System.get_env("SHOP_SERVICE_PORT", "22101"))

    {:ok, channel} = GRPC.Stub.connect("#{host}:#{port}")
    channel
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
  def handle_event("add_to_cart", _params, socket) do
    # カート機能は未実装
    {:noreply,
     socket
     |> put_flash(:info, "カート機能は準備中です")}
  end

  @impl true
  def handle_event("update_quantity", %{"quantity" => qty}, socket) do
    quantity =
      case Integer.parse(qty) do
        {n, ""} when n > 0 -> n
        _ -> 1
      end

    {:noreply, assign(socket, :quantity, quantity)}
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <!-- ナビゲーションバー -->
      <nav class="bg-white shadow-md">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16">
            <div class="flex items-center">
              <.link navigate="/dashboard" class="text-2xl font-bold text-gray-900 hover:text-gray-700">
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
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                商品一覧
              </.link>
              <.link
                navigate="/auth"
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
        <%= if @loading do %>
          <div class="text-center py-12">
            <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p class="mt-4 text-gray-600">商品を読み込んでいます...</p>
          </div>
        <% end %>

        <%= if @error do %>
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            <%= @error %>
          </div>
          <.link
            navigate="/products"
            class="text-blue-600 hover:text-blue-800"
          >
            ← 商品一覧に戻る
          </.link>
        <% end %>

        <%= if @product do %>
          <div class="bg-white shadow rounded-lg overflow-hidden">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 p-6">
              <!-- 商品画像 -->
              <div class="aspect-square bg-gradient-to-br from-blue-100 to-purple-100 rounded-lg flex items-center justify-center">
                <svg class="h-32 w-32 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                  />
                </svg>
              </div>

              <!-- 商品情報 -->
              <div class="flex flex-col">
                <!-- ショップ名 -->
                <%= if @shop do %>
                  <p class="text-sm text-gray-500 mb-2">
                    <%= @shop.name %>
                  </p>
                <% end %>

                <h1 class="text-3xl font-bold text-gray-900 mb-4">
                  <%= @product.name %>
                </h1>

                <div class="text-4xl font-bold text-blue-600 mb-6">
                  <%= format_price(@product.price) %>
                </div>

                <!-- 在庫情報 -->
                <div class="mb-6">
                  <%= if @product.stock_quantity > 0 do %>
                    <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
                      ✓ 在庫あり（残り<%= @product.stock_quantity %>個）
                    </span>
                  <% else %>
                    <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-red-100 text-red-800">
                      × 在庫切れ
                    </span>
                  <% end %>
                </div>

                <!-- 商品説明 -->
                <div class="mb-6">
                  <h2 class="text-lg font-semibold text-gray-900 mb-2">商品説明</h2>
                  <p class="text-gray-700 whitespace-pre-line">
                    <%= @product.description %>
                  </p>
                </div>

                <!-- 商品詳細 -->
                <div class="mb-6 border-t pt-6">
                  <h2 class="text-lg font-semibold text-gray-900 mb-3">商品詳細</h2>
                  <dl class="grid grid-cols-2 gap-3 text-sm">
                    <%= if @product.category && @product.category != "" do %>
                      <dt class="text-gray-500">カテゴリー</dt>
                      <dd class="text-gray-900"><%= @product.category %></dd>
                    <% end %>
                    <%= if @product.weight && @product.weight != "0" do %>
                      <dt class="text-gray-500">重量</dt>
                      <dd class="text-gray-900"><%= @product.weight %>g</dd>
                    <% end %>
                    <%= if @product.size && @product.size != "" do %>
                      <dt class="text-gray-500">サイズ</dt>
                      <dd class="text-gray-900"><%= @product.size %></dd>
                    <% end %>
                    <%= if @product.jan_code && @product.jan_code != "" do %>
                      <dt class="text-gray-500">JANコード</dt>
                      <dd class="text-gray-900"><%= @product.jan_code %></dd>
                    <% end %>
                  </dl>
                </div>

                <!-- カートに追加（未実装） -->
                <div class="mt-auto">
                  <button
                    phx-click="add_to_cart"
                    disabled={@product.stock_quantity == 0}
                    class={"w-full py-3 px-6 rounded-lg font-semibold text-white transition-colors " <>
                      if(@product.stock_quantity > 0,
                        do: "bg-blue-600 hover:bg-blue-700",
                        else: "bg-gray-400 cursor-not-allowed"
                      )}
                  >
                    <%= if @product.stock_quantity > 0 do %>
                      カートに追加（準備中）
                    <% else %>
                      在庫切れ
                    <% end %>
                  </button>

                  <.link
                    navigate="/products"
                    class="block text-center mt-4 text-blue-600 hover:text-blue-800"
                  >
                    ← 商品一覧に戻る
                  </.link>
                </div>
              </div>
            </div>
          </div>
        <% end %>
      </main>
    </div>
    """
  end
end
