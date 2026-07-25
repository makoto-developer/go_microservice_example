defmodule ShopMallWebWeb.ProductDetailLive do
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.CustomerServiceClient, as: Customers
  alias ShopMallWeb.ReviewServiceClient, as: Reviews

  alias ShopService.V1.{
    ShopService.Stub,
    GetProductRequest,
    GetShopRequest
  }

  @impl true
  def mount(%{"id" => product_id}, session, socket) do
    {:ok,
     socket
     |> assign(:current_user_id, session["user_id"])
     |> assign(:product_id, product_id)
     |> assign(:product, nil)
     |> assign(:shop, nil)
     |> assign(:loading, true)
     |> assign(:error, nil)
     |> assign(:quantity, 1)
     |> assign(:payment_method, "credit_card")
     |> assign(:review_summary, nil)
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

  # 住所管理が未実装のため、デモ用の固定住所 ID を使う
  @demo_address_id "00000000-0000-0000-0000-000000000001"

  defp place_order(socket, user_id) do
    product = socket.assigns.product
    cod? = socket.assigns.payment_method == "cash_on_delivery"

    request = %OrderService.V1.CreateOrderRequest{
      customer_id: user_id,
      cart_items: [
        %OrderService.V1.CartItemInput{
          product_id: product.id,
          quantity: socket.assigns.quantity,
          unit_price: price_to_int(product.price)
        }
      ],
      shipping_address_id: @demo_address_id,
      payment_method: if(cod?, do: :CASH_ON_DELIVERY, else: :CREDIT_CARD),
      payment_method_id: if(cod?, do: "", else: "pm_demo_card"),
      shipping_method: "standard"
    }

    with {:ok, channel} <- connect_order_service(),
         {:ok, response} <- OrderService.V1.OrderService.Stub.create_order(channel, request) do
      message =
        if cod?,
          do: "注文が確定しました。お支払いは商品お届け時です(注文ID: #{response.order_id})",
          else: "注文と決済が完了しました(注文ID: #{response.order_id})"

      {:noreply, put_flash(socket, :info, message)}
    else
      {:error, %GRPC.RPCError{message: message}} ->
        {:noreply, put_flash(socket, :error, "注文に失敗しました: #{message}")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "注文に失敗しました: #{inspect(reason)}")}
    end
  end

  defp connect_order_service do
    host = System.get_env("ORDER_SERVICE_HOST", "localhost")
    port = System.get_env("ORDER_SERVICE_PORT", "50055")
    GRPC.Stub.connect("#{host}:#{port}")
  end

  defp price_to_int(price) when is_binary(price) do
    case Decimal.parse(price) do
      {decimal, ""} ->
        decimal |> Decimal.round(0) |> Decimal.to_integer()

      _ ->
        0
    end
  end

  defp price_to_int(_), do: 0

  @impl true
  def handle_event("add_to_cart", _params, socket) do
    case socket.assigns.current_user_id do
      nil ->
        {:noreply, put_flash(socket, :error, "カートを使うにはログインが必要です")}

      user_id ->
        case Customers.add_to_cart(user_id, socket.assigns.product.id, socket.assigns.quantity) do
          {:ok, _} ->
            {:noreply, put_flash(socket, :info, "カートに追加しました")}

          {:error, reason} ->
            {:noreply, put_flash(socket, :error, "カートへの追加に失敗しました: #{reason}")}
        end
    end
  end

  @impl true
  def handle_event("add_to_favorite", _params, socket) do
    case socket.assigns.current_user_id do
      nil ->
        {:noreply, put_flash(socket, :error, "お気に入りにはログインが必要です")}

      user_id ->
        case Customers.add_to_favorite(user_id, socket.assigns.product.id) do
          {:ok, _} ->
            {:noreply, put_flash(socket, :info, "お気に入りに追加しました")}

          {:error, reason} ->
            {:noreply, put_flash(socket, :error, "お気に入りへの追加に失敗しました: #{reason}")}
        end
    end
  end

  @impl true
  def handle_event("post_review", %{"rating" => rating, "review_text" => text}, socket) do
    case socket.assigns.current_user_id do
      nil ->
        {:noreply, put_flash(socket, :error, "レビュー投稿にはログインが必要です")}

      user_id ->
        case Reviews.post_review(%{
               customer_id: user_id,
               product_id: socket.assigns.product.id,
               rating: String.to_integer(rating),
               content: text
             }) do
          {:ok, _} ->
            {:noreply, put_flash(socket, :info, "レビューを投稿しました(公開には承認が必要です)")}

          {:error, reason} ->
            {:noreply, put_flash(socket, :error, "レビュー投稿に失敗しました: #{reason}")}
        end
    end
  end

  @impl true
  def handle_event("load_reviews", _params, socket) do
    rating_note =
      case Reviews.get_product_rating(socket.assigns.product.id) do
        {:ok, resp} -> resp.message
        {:error, reason} -> "評価取得エラー: #{reason}"
      end

    reviews_note =
      case Reviews.get_reviews_by_product(socket.assigns.product.id) do
        {:ok, resp} -> resp.message
        {:error, reason} -> "レビュー取得エラー: #{reason}"
      end

    {:noreply, assign(socket, :review_summary, "#{rating_note} / #{reviews_note}")}
  end

  @impl true
  def handle_event("buy_now", _params, socket) do
    case socket.assigns.current_user_id do
      nil ->
        {:noreply, put_flash(socket, :error, "購入にはログインが必要です")}

      user_id ->
        place_order(socket, user_id)
    end
  end

  @impl true
  def handle_event("select_payment_method", %{"method" => method}, socket)
      when method in ["credit_card", "cash_on_delivery"] do
    {:noreply, assign(socket, :payment_method, method)}
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
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                商品一覧
              </.link>
              <.link
                navigate="/orders"
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                注文履歴
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
                <svg
                  class="h-32 w-32 text-gray-400"
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
              <div class="flex flex-col">
                <!-- ショップ名 -->
                <%= if @shop do %>
                  <p class="text-sm text-gray-500 mb-2">
                    {@shop.name}
                  </p>
                <% end %>

                <h1 class="text-3xl font-bold text-gray-900 mb-4">
                  {@product.name}
                </h1>

                <div class="text-4xl font-bold text-blue-600 mb-6">
                  {format_price(@product.price)}
                </div>
                
    <!-- 在庫情報 -->
                <div class="mb-6">
                  <%= if @product.stock_quantity > 0 do %>
                    <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
                      ✓ 在庫あり（残り{@product.stock_quantity}個）
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
                    {@product.description}
                  </p>
                </div>
                
    <!-- 商品詳細 -->
                <div class="mb-6 border-t pt-6">
                  <h2 class="text-lg font-semibold text-gray-900 mb-3">商品詳細</h2>
                  <dl class="grid grid-cols-2 gap-3 text-sm">
                    <%= if @product.category && @product.category != "" do %>
                      <dt class="text-gray-500">カテゴリー</dt>
                      <dd class="text-gray-900">{@product.category}</dd>
                    <% end %>
                    <%= if @product.weight && @product.weight != "0" do %>
                      <dt class="text-gray-500">重量</dt>
                      <dd class="text-gray-900">{@product.weight}g</dd>
                    <% end %>
                    <%= if @product.size && @product.size != "" do %>
                      <dt class="text-gray-500">サイズ</dt>
                      <dd class="text-gray-900">{@product.size}</dd>
                    <% end %>
                    <%= if @product.jan_code && @product.jan_code != "" do %>
                      <dt class="text-gray-500">JANコード</dt>
                      <dd class="text-gray-900">{@product.jan_code}</dd>
                    <% end %>
                  </dl>
                </div>
                
    <!-- 支払い方法 -->
                <div class="mb-6">
                  <h2 class="text-lg font-semibold text-gray-900 mb-3">お支払い方法</h2>
                  <div class="grid grid-cols-2 gap-3">
                    <button
                      phx-click="select_payment_method"
                      phx-value-method="credit_card"
                      class={"border rounded-lg px-4 py-3 text-sm font-medium text-left transition-colors " <>
                        if(@payment_method == "credit_card",
                          do: "border-blue-600 bg-blue-50 text-blue-700",
                          else: "border-gray-300 text-gray-700 hover:border-gray-400"
                        )}
                    >
                      💳 クレジットカード <span class="block text-xs text-gray-500 mt-1">今すぐ決済されます</span>
                    </button>
                    <button
                      phx-click="select_payment_method"
                      phx-value-method="cash_on_delivery"
                      class={"border rounded-lg px-4 py-3 text-sm font-medium text-left transition-colors " <>
                        if(@payment_method == "cash_on_delivery",
                          do: "border-blue-600 bg-blue-50 text-blue-700",
                          else: "border-gray-300 text-gray-700 hover:border-gray-400"
                        )}
                    >
                      📦 代金引換 <span class="block text-xs text-gray-500 mt-1">お届け時に現金でお支払い</span>
                    </button>
                  </div>
                </div>
                
    <!-- 購入(注文作成 → 決済まで実行) -->
                <div class="mt-auto">
                  <button
                    phx-click="buy_now"
                    disabled={@product.stock_quantity == 0}
                    class={"w-full py-3 px-6 rounded-lg font-semibold text-white transition-colors " <>
                      if(@product.stock_quantity > 0,
                        do: "bg-blue-600 hover:bg-blue-700",
                        else: "bg-gray-400 cursor-not-allowed"
                      )}
                  >
                    <%= if @product.stock_quantity > 0 do %>
                      購入する
                    <% else %>
                      在庫切れ
                    <% end %>
                  </button>

                  <div class="grid grid-cols-2 gap-3 mt-3">
                    <button
                      phx-click="add_to_cart"
                      disabled={@product.stock_quantity == 0}
                      class="py-2.5 px-4 rounded-lg font-medium text-blue-700 border border-blue-300 hover:bg-blue-50 disabled:opacity-50"
                    >
                      🛒 カートに入れる
                    </button>
                    <button
                      phx-click="add_to_favorite"
                      class="py-2.5 px-4 rounded-lg font-medium text-pink-600 border border-pink-300 hover:bg-pink-50"
                    >
                      ♡ お気に入り
                    </button>
                  </div>

                  <div class="mt-6 border-t pt-4">
                    <div class="flex items-center justify-between mb-2">
                      <h2 class="text-sm font-semibold text-gray-900">レビュー</h2>
                      <button
                        phx-click="load_reviews"
                        class="text-xs text-blue-600 hover:text-blue-800"
                      >
                        評価・レビューを読み込む
                      </button>
                    </div>
                    <div :if={@review_summary} class="text-xs text-gray-600 mb-3">
                      {@review_summary}
                    </div>
                    <h2 class="text-sm font-semibold text-gray-900 mb-2">レビューを書く</h2>
                    <form phx-submit="post_review" class="space-y-2">
                      <select
                        name="rating"
                        class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                      >
                        <option :for={n <- 5..1//-1} value={n}>{String.duplicate("★", n)}</option>
                      </select>
                      <textarea
                        name="review_text"
                        rows="2"
                        required
                        placeholder="商品の感想を書いてください"
                        class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                      ></textarea>
                      <div class="flex justify-end">
                        <button
                          type="submit"
                          class="px-4 py-1.5 text-sm font-medium text-white bg-gray-800 rounded-md hover:bg-gray-700"
                        >
                          投稿する
                        </button>
                      </div>
                    </form>
                  </div>

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
