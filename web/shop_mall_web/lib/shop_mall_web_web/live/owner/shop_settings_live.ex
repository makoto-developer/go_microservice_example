defmodule ShopMallWebWeb.Owner.ShopSettingsLive do
  @moduledoc """
  加盟店の店舗設定画面。
  自分の店舗一覧(GetShopsByOwner)から店舗を選び、
  基本情報の更新(UpdateShop)と公開/非公開の切り替え(ToggleShopPublish)を行う。
  """
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.ShopServiceClient, as: Shops

  @impl true
  def mount(_params, session, socket) do
    owner_id = session["user_id"] || "admin-user-id"

    {:ok,
     socket
     |> assign(:owner_id, owner_id)
     |> assign(:shops, [])
     |> assign(:shop, nil)
     |> assign(:loading, true)
     |> assign(:error, nil)
     |> load_shops()}
  end

  defp load_shops(socket) do
    case Shops.get_shops_by_owner(socket.assigns.owner_id) do
      {:ok, response} ->
        shop = socket.assigns.shop || List.first(response.shops)

        socket
        |> assign(:shops, response.shops)
        |> assign(:shop, shop)
        |> assign(:loading, false)
        |> assign(:error, nil)

      {:error, reason} ->
        socket
        |> assign(:loading, false)
        |> assign(:error, "店舗一覧の取得に失敗しました: #{reason}")
    end
  end

  @impl true
  def handle_event("select_shop", %{"shop-id" => shop_id}, socket) do
    shop = Enum.find(socket.assigns.shops, &(&1.id == shop_id))
    {:noreply, assign(socket, :shop, shop)}
  end

  @impl true
  def handle_event("save_shop", params, socket) do
    shop = socket.assigns.shop

    case Shops.update_shop(%{
           shop_id: shop.id,
           name: params["name"],
           description: params["description"],
           logo_url: shop.logo_url,
           business_hours: params["business_hours"],
           return_policy: params["return_policy"]
         }) do
      {:ok, response} ->
        note = if response.requires_reapproval, do: "(再審査が必要です)", else: ""

        {:noreply,
         socket
         |> put_flash(:info, "店舗情報を更新しました#{note}")
         |> load_shops()}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "更新に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("toggle_publish", _params, socket) do
    shop = socket.assigns.shop

    case Shops.toggle_shop_publish(shop.id, !shop.published) do
      {:ok, response} ->
        label = if response.published, do: "公開", else: "非公開"

        {:noreply,
         socket
         |> put_flash(:info, "店舗を#{label}にしました")
         |> load_shops()}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "切り替えに失敗しました: #{reason}")}
    end
  end

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
                navigate="/owner/products"
                class="text-emerald-200 hover:text-white text-sm font-medium"
              >
                商品管理
              </.link>
              <.link
                navigate="/owner/payments"
                class="text-emerald-200 hover:text-white text-sm font-medium"
              >
                決済確認
              </.link>
              <span class="text-white text-sm font-medium">店舗設定</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-6">店舗設定</h1>

        <%= if @error do %>
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            {@error}
          </div>
        <% end %>

        <%= if @loading do %>
          <div class="text-center py-12">
            <div class="inline-block animate-spin rounded-full h-10 w-10 border-b-2 border-emerald-700">
            </div>
          </div>
        <% end %>

        <%= if length(@shops) > 1 do %>
          <div class="mb-4 flex space-x-2">
            <button
              :for={s <- @shops}
              phx-click="select_shop"
              phx-value-shop-id={s.id}
              class={"px-3 py-1.5 rounded-full text-sm font-medium " <>
                if(@shop && @shop.id == s.id,
                  do: "bg-emerald-700 text-white",
                  else: "bg-white text-gray-600 border border-gray-300"
                )}
            >
              {s.name}
            </button>
          </div>
        <% end %>

        <%= if @shop do %>
          <div class="bg-white shadow rounded-lg p-6 mb-6">
            <div class="flex items-center justify-between mb-4">
              <h2 class="text-lg font-semibold text-gray-900">基本情報</h2>
              <div class="flex items-center space-x-3">
                <span class={"inline-flex px-2 py-0.5 rounded-full text-xs font-medium " <>
                  if(@shop.published, do: "bg-green-100 text-green-800", else: "bg-gray-200 text-gray-700")}>
                  {if @shop.published, do: "公開中", else: "非公開"}
                </span>
                <button
                  phx-click="toggle_publish"
                  class="px-3 py-1.5 text-sm font-medium text-emerald-700 border border-emerald-300 rounded-md hover:bg-emerald-50"
                >
                  {if @shop.published, do: "非公開にする", else: "公開する"}
                </button>
              </div>
            </div>

            <form phx-submit="save_shop" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">店舗名</label>
                <input
                  type="text"
                  name="name"
                  value={@shop.name}
                  required
                  class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">店舗紹介</label>
                <textarea
                  name="description"
                  rows="3"
                  class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                >{@shop.description}</textarea>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">営業時間</label>
                <input
                  type="text"
                  name="business_hours"
                  value={@shop.business_hours}
                  placeholder="例: 平日 10:00-18:00"
                  class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">返品ポリシー</label>
                <textarea
                  name="return_policy"
                  rows="2"
                  placeholder="例: 商品到着後7日以内、未開封に限り返品可"
                  class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                >{@shop.return_policy}</textarea>
              </div>
              <div class="flex justify-end">
                <button
                  type="submit"
                  class="px-4 py-2 text-sm font-medium text-white bg-emerald-600 rounded-md hover:bg-emerald-700"
                >
                  保存する
                </button>
              </div>
            </form>
          </div>
        <% end %>

        <%= if not @loading and @shops == [] and is_nil(@error) do %>
          <div class="bg-white shadow rounded-lg p-12 text-center text-gray-500">
            店舗がまだありません
            <div class="mt-4">
              <.link
                navigate="/owner/shop/register"
                class="text-emerald-700 hover:text-emerald-900 font-medium"
              >
                店舗を登録する →
              </.link>
            </div>
          </div>
        <% end %>
      </main>
    </div>
    """
  end
end
