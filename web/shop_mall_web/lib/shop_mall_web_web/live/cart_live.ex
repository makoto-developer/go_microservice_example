defmodule ShopMallWebWeb.CartLive do
  @moduledoc """
  カート画面。customer サービスのカート RPC(GetCart / UpdateCartItemQuantity /
  RemoveFromCart / MergeGuestCart)を使用する。
  """
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.CustomerServiceClient, as: Customers

  @impl true
  def mount(_params, session, socket) do
    {:ok,
     socket
     |> assign(:customer_id, session["user_id"])
     |> assign(:items, [])
     |> assign(:total_quantity, 0)
     |> load_cart()}
  end

  defp load_cart(%{assigns: %{customer_id: nil}} = socket), do: socket

  defp load_cart(socket) do
    case Customers.get_cart(socket.assigns.customer_id) do
      {:ok, resp} ->
        socket
        |> assign(:items, resp.cart_items)
        |> assign(:total_quantity, resp.total_quantity)

      {:error, _} ->
        socket
    end
  end

  @impl true
  def handle_event("change_quantity", %{"item-id" => item_id, "quantity" => qty}, socket) do
    case Customers.update_cart_item_quantity(
           item_id,
           socket.assigns.customer_id,
           String.to_integer(qty)
         ) do
      {:ok, _} -> {:noreply, load_cart(socket)}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "数量の変更に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("remove_item", %{"item-id" => item_id}, socket) do
    case Customers.remove_from_cart(item_id, socket.assigns.customer_id) do
      {:ok, _} -> {:noreply, socket |> put_flash(:info, "カートから削除しました") |> load_cart()}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "削除に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("merge_guest_cart", _params, socket) do
    # ログイン前のゲストカート(セッション ID ベース)を取り込む
    case Customers.merge_guest_cart(socket.assigns.customer_id, "guest-session-demo") do
      {:ok, _} ->
        {:noreply, socket |> put_flash(:info, "ゲストカートを統合しました") |> load_cart()}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "統合に失敗しました: #{reason}")}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
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
                navigate="/products"
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                商品一覧
              </.link>
              <span class="text-gray-900 px-3 py-2 text-sm font-semibold">カート</span>
              <.link
                navigate="/mypage"
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                マイページ
              </.link>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between mb-6">
          <h1 class="text-2xl font-bold text-gray-900">
            カート <span class="text-sm font-normal text-gray-500">({@total_quantity} 点)</span>
          </h1>
          <button
            phx-click="merge_guest_cart"
            class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
            title="ログイン前に入れた商品を取り込む"
          >
            ゲストカートを統合
          </button>
        </div>

        <%= if is_nil(@customer_id) do %>
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            カートを見るにはログインが必要です
          </div>
        <% else %>
          <%= if @items == [] do %>
            <div class="bg-white shadow rounded-lg p-12 text-center text-gray-500">
              カートは空です
              <div class="mt-4">
                <.link navigate="/products" class="text-blue-600 hover:text-blue-800 font-medium">
                  商品を見る →
                </.link>
              </div>
            </div>
          <% end %>

          <div class="space-y-3">
            <div
              :for={item <- @items}
              class="bg-white shadow rounded-lg p-4 flex items-center justify-between"
            >
              <.link
                navigate={"/products/#{item.product_id}"}
                class="text-sm text-blue-600 hover:text-blue-800 font-mono"
              >
                {String.slice(item.product_id, 0, 8)}…
              </.link>
              <div class="flex items-center space-x-3">
                <form phx-change="change_quantity" class="flex items-center space-x-2">
                  <input type="hidden" name="item-id" value={item.id} />
                  <select name="quantity" class="border border-gray-300 rounded-md px-2 py-1 text-sm">
                    <option :for={n <- 1..9} value={n} selected={item.quantity == n}>{n}</option>
                  </select>
                </form>
                <button
                  phx-click="remove_item"
                  phx-value-item-id={item.id}
                  class="text-sm text-red-600 hover:text-red-800"
                >
                  削除
                </button>
              </div>
            </div>
          </div>
        <% end %>
      </main>
    </div>
    """
  end
end
