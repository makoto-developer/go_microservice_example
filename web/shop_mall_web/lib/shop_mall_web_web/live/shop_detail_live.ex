defmodule ShopMallWebWeb.ShopDetailLive do
  use ShopMallWebWeb, :live_view

  alias ShopService.V1.{
    ShopService.Stub,
    GetShopRequest,
    ListProductsRequest
  }

  @impl true
  def mount(%{"id" => shop_id}, _session, socket) do
    socket =
      socket
      |> assign(:shop_id, shop_id)
      |> assign(:shop, nil)
      |> assign(:products, [])
      |> assign(:loading, true)
      |> assign(:error, nil)

    send(self(), :load_data)

    {:ok, socket}
  end

  @impl true
  def handle_info(:load_data, socket) do
    shop_id = socket.assigns.shop_id

    # ショップ情報を取得
    shop_result = call_shop_service(:get_shop, %GetShopRequest{shop_id: shop_id})

    # 商品一覧を取得
    products_result =
      call_shop_service(:list_products, %ListProductsRequest{
        shop_id: shop_id,
        category: "",
        published_only: true,
        limit: 100,
        offset: 0
      })

    socket =
      case {shop_result, products_result} do
        {{:ok, shop_response}, {:ok, products_response}} ->
          socket
          |> assign(:shop, shop_response.shop)
          |> assign(:products, products_response.products)
          |> assign(:loading, false)

        {{:error, reason}, _} ->
          socket
          |> assign(:error, "ショップ情報の取得に失敗しました: #{reason}")
          |> assign(:loading, false)

        {_, {:error, reason}} ->
          socket
          |> assign(:error, "商品情報の取得に失敗しました: #{reason}")
          |> assign(:loading, false)
      end

    {:noreply, socket}
  end

  defp call_shop_service(:get_shop, request) do
    channel = get_shop_channel()

    case Stub.get_shop(channel, request) do
      {:ok, response} -> {:ok, response}
      {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
      error -> {:error, "接続エラー: #{inspect(error)}"}
    end
  end

  defp call_shop_service(:list_products, request) do
    channel = get_shop_channel()

    case Stub.list_products(channel, request) do
      {:ok, response} -> {:ok, response}
      {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
      error -> {:error, "接続エラー: #{inspect(error)}"}
    end
  end

  defp get_shop_channel do
    host = System.get_env("SHOP_SERVICE_HOST", "localhost")
    port = String.to_integer(System.get_env("SHOP_SERVICE_PORT", "22003"))

    {:ok, channel} = GRPC.Stub.connect("#{host}:#{port}")
    channel
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-100">
      <!-- ヘッダー -->
      <header class="bg-white shadow">
        <div class="max-w-7xl mx-auto px-4 py-6">
          <nav class="flex items-center space-x-2 text-sm text-gray-500 mb-4">
            <.link href="/" class="hover:text-blue-600">ホーム</.link>
            <span>/</span>
            <.link href="/shops" class="hover:text-blue-600">ショップ一覧</.link>
            <span>/</span>
            <span class="text-gray-900"><%= if @shop, do: @shop.name, else: "..." %></span>
          </nav>
        </div>
      </header>

      <main class="max-w-7xl mx-auto px-4 py-8">
        <%= if @loading do %>
          <div class="text-center py-12">
            <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <p class="mt-4 text-gray-600">読み込み中...</p>
          </div>
        <% else %>
          <%= if @error do %>
            <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
              <%= @error %>
            </div>
          <% end %>

          <%= if @shop do %>
            <!-- ショップヘッダー -->
            <div class="bg-white rounded-lg shadow overflow-hidden mb-8">
              <div class="h-48 bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center">
                <%= if @shop.logo_url && @shop.logo_url != "" do %>
                  <img src={@shop.logo_url} alt={@shop.name} class="h-32 w-32 object-cover rounded-full border-4 border-white shadow" />
                <% else %>
                  <div class="h-32 w-32 bg-white rounded-full flex items-center justify-center border-4 border-white shadow">
                    <span class="text-5xl font-bold text-gray-400">
                      <%= String.first(@shop.name) %>
                    </span>
                  </div>
                <% end %>
              </div>
              <div class="p-6">
                <h1 class="text-3xl font-bold text-gray-900"><%= @shop.name %></h1>
                <p class="mt-2 text-gray-600"><%= @shop.description %></p>

                <div class="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                  <%= if @shop.owner_name && @shop.owner_name != "" do %>
                    <div class="flex items-center text-gray-500">
                      <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                      </svg>
                      <span>運営者: <%= @shop.owner_name %></span>
                    </div>
                  <% end %>

                  <%= if @shop.business_hours && @shop.business_hours != "" do %>
                    <div class="flex items-center text-gray-500">
                      <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <span><%= @shop.business_hours %></span>
                    </div>
                  <% end %>

                  <%= if @shop.phone_number && @shop.phone_number != "" do %>
                    <div class="flex items-center text-gray-500">
                      <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                      </svg>
                      <span><%= @shop.phone_number %></span>
                    </div>
                  <% end %>
                </div>

                <%= if @shop.return_policy && @shop.return_policy != "" do %>
                  <div class="mt-4 p-4 bg-gray-50 rounded-md">
                    <h3 class="text-sm font-medium text-gray-700">返品ポリシー</h3>
                    <p class="mt-1 text-sm text-gray-600"><%= @shop.return_policy %></p>
                  </div>
                <% end %>
              </div>
            </div>

            <!-- 商品一覧 -->
            <div class="bg-white rounded-lg shadow p-6">
              <h2 class="text-xl font-bold text-gray-900 mb-6">
                商品一覧
                <span class="text-sm font-normal text-gray-500">(<%= length(@products) %>件)</span>
              </h2>

              <%= if length(@products) == 0 do %>
                <div class="text-center py-12">
                  <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                  </svg>
                  <p class="mt-4 text-gray-600">このショップにはまだ商品がありません</p>
                </div>
              <% else %>
                <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                  <%= for product <- @products do %>
                    <.link
                      href={"/products/#{product.id}"}
                      class="group block bg-white border rounded-lg overflow-hidden hover:shadow-lg transition-shadow"
                    >
                      <div class="aspect-square bg-gray-200 flex items-center justify-center">
                        <svg class="h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                      </div>
                      <div class="p-4">
                        <h3 class="text-sm font-medium text-gray-900 group-hover:text-blue-600 line-clamp-2">
                          <%= product.name %>
                        </h3>
                        <p class="mt-1 text-lg font-bold text-gray-900">
                          ¥<%= format_price(product.price) %>
                        </p>
                        <%= if product.stock_quantity <= 0 do %>
                          <span class="mt-1 inline-block px-2 py-1 text-xs bg-red-100 text-red-800 rounded">
                            在庫切れ
                          </span>
                        <% end %>
                      </div>
                    </.link>
                  <% end %>
                </div>
              <% end %>
            </div>
          <% else %>
            <div class="text-center py-12 bg-white rounded-lg shadow">
              <h2 class="text-xl font-semibold text-gray-800">ショップが見つかりません</h2>
              <.link href="/shops" class="mt-4 inline-block text-blue-600 hover:text-blue-800">
                ショップ一覧に戻る
              </.link>
            </div>
          <% end %>
        <% end %>
      </main>
    </div>
    """
  end

  defp format_price(price) when is_binary(price) do
    case Integer.parse(price) do
      {num, _} -> Number.Delimit.number_to_delimited(num, precision: 0)
      :error -> price
    end
  end

  defp format_price(price), do: price
end
