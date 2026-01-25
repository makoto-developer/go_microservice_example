defmodule ShopMallWebWeb.ShopListLive do
  use ShopMallWebWeb, :live_view

  alias ShopService.V1.{
    ShopService.Stub,
    ListShopsRequest
  }

  @impl true
  def mount(_params, _session, socket) do
    socket =
      socket
      |> assign(:shops, [])
      |> assign(:loading, true)
      |> assign(:error, nil)

    send(self(), :load_shops)

    {:ok, socket}
  end

  @impl true
  def handle_info(:load_shops, socket) do
    request = %ListShopsRequest{
      published_only: true,
      limit: 100,
      offset: 0
    }

    socket =
      case call_shop_service(request) do
        {:ok, response} ->
          socket
          |> assign(:shops, response.shops)
          |> assign(:loading, false)

        {:error, reason} ->
          socket
          |> assign(:error, "ショップ一覧の取得に失敗しました: #{reason}")
          |> assign(:loading, false)
      end

    {:noreply, socket}
  end

  defp call_shop_service(request) do
    channel = get_shop_channel()

    case Stub.list_shops(channel, request) do
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
          <div class="flex justify-between items-center">
            <h1 class="text-3xl font-bold text-gray-900">ショップ一覧</h1>
            <nav class="space-x-4">
              <.link href="/" class="text-blue-600 hover:text-blue-800">ホーム</.link>
              <.link href="/products" class="text-blue-600 hover:text-blue-800">商品一覧</.link>
              <.link href="/auth" class="text-blue-600 hover:text-blue-800">ログイン</.link>
            </nav>
          </div>
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

          <%= if length(@shops) == 0 do %>
            <div class="text-center py-12 bg-white rounded-lg shadow">
              <svg class="mx-auto h-16 w-16 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
              <h2 class="mt-4 text-xl font-semibold text-gray-800">ショップがまだありません</h2>
              <p class="mt-2 text-gray-600">ショップモールに出店したい方は<.link href="/owner/auth" class="text-purple-600 hover:text-purple-800 font-medium">オーナー登録</.link>へ</p>
            </div>
          <% else %>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <%= for shop <- @shops do %>
                <.link
                  href={"/shops/#{shop.id}"}
                  class="bg-white rounded-lg shadow hover:shadow-lg transition-shadow overflow-hidden"
                >
                  <div class="h-40 bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center">
                    <%= if shop.logo_url && shop.logo_url != "" do %>
                      <img src={shop.logo_url} alt={shop.name} class="h-24 w-24 object-cover rounded-full" />
                    <% else %>
                      <div class="h-24 w-24 bg-white rounded-full flex items-center justify-center">
                        <span class="text-3xl font-bold text-gray-400">
                          <%= String.first(shop.name || "S") %>
                        </span>
                      </div>
                    <% end %>
                  </div>
                  <div class="p-4">
                    <h3 class="text-lg font-semibold text-gray-800"><%= shop.name %></h3>
                    <p class="mt-1 text-sm text-gray-600 line-clamp-2"><%= shop.description %></p>
                    <div class="mt-3 flex items-center text-sm text-gray-500">
                      <%= if shop.owner_name && shop.owner_name != "" do %>
                        <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                        <span><%= shop.owner_name %></span>
                      <% end %>
                    </div>
                  </div>
                </.link>
              <% end %>
            </div>
          <% end %>
        <% end %>

        <!-- オーナー登録への誘導 -->
        <div class="mt-12 bg-gradient-to-r from-purple-600 to-indigo-600 rounded-lg shadow-lg p-8 text-center text-white">
          <h2 class="text-2xl font-bold">あなたもショップモールに出店しませんか？</h2>
          <p class="mt-2 text-purple-100">簡単な登録で、あなたの商品を全国のお客様に届けられます</p>
          <.link
            href="/owner/auth"
            class="mt-4 inline-block bg-white text-purple-600 px-6 py-3 rounded-md font-medium hover:bg-purple-50 transition-colors"
          >
            オーナー登録はこちら
          </.link>
        </div>
      </main>
    </div>
    """
  end
end
