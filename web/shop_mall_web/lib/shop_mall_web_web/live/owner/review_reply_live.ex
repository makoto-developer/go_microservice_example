defmodule ShopMallWebWeb.Owner.ReviewReplyLive do
  @moduledoc """
  加盟店のレビュー返信画面。
  商品レビューへの返信の投稿(PostShopReply)・修正(UpdateShopReply)・
  削除(DeleteShopReply)を行う。
  """
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.ReviewServiceClient, as: Reviews
  alias ShopMallWeb.ShopServiceClient, as: Shops

  @impl true
  def mount(_params, session, socket) do
    owner_id = session["user_id"] || "admin-user-id"

    shop =
      case Shops.get_shops_by_owner(owner_id) do
        {:ok, %{shops: [shop | _]}} -> shop
        _ -> nil
      end

    {:ok,
     socket
     |> assign(:shop, shop)
     |> assign(:last_result, nil)}
  end

  @impl true
  def handle_event("post_reply", %{"review_id" => review_id, "content" => content}, socket) do
    run(socket, Reviews.post_shop_reply(review_id, socket.assigns.shop.id, content))
  end

  @impl true
  def handle_event("update_reply", %{"reply_id" => reply_id, "content" => content}, socket) do
    run(socket, Reviews.update_shop_reply(reply_id, socket.assigns.shop.id, content))
  end

  @impl true
  def handle_event("delete_reply", %{"reply_id" => reply_id}, socket) do
    run(socket, Reviews.delete_shop_reply(reply_id, socket.assigns.shop.id))
  end

  defp run(socket, result) do
    case result do
      {:ok, resp} ->
        {:noreply, socket |> assign(:last_result, resp.message) |> put_flash(:info, resp.message)}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "操作に失敗しました: #{reason}")}
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
              <span class="text-white text-sm font-medium">レビュー返信</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-2xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-6">
          レビュー返信 <span :if={@shop} class="text-sm font-normal text-gray-500">({@shop.name})</span>
        </h1>

        <div class="bg-white shadow rounded-lg p-6 mb-6">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">レビューに返信する</h2>
          <form phx-submit="post_reply" class="space-y-3">
            <input
              type="text"
              name="review_id"
              required
              placeholder="レビューID"
              class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <textarea
              name="content"
              rows="3"
              required
              placeholder="ご購入ありがとうございます。…"
              class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
            ></textarea>
            <div class="flex justify-end">
              <button
                type="submit"
                class="px-4 py-2 text-sm font-medium text-white bg-emerald-600 rounded-md hover:bg-emerald-700"
              >
                返信を投稿
              </button>
            </div>
          </form>
        </div>

        <div class="bg-white shadow rounded-lg p-6 mb-6">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">返信を修正する</h2>
          <form phx-submit="update_reply" class="space-y-3">
            <input
              type="text"
              name="reply_id"
              required
              placeholder="返信ID"
              class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <textarea
              name="content"
              rows="2"
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
            ></textarea>
            <div class="flex justify-end">
              <button
                type="submit"
                class="px-4 py-2 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
              >
                修正する
              </button>
            </div>
          </form>
        </div>

        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">返信を削除する</h2>
          <form phx-submit="delete_reply" class="flex items-end space-x-2">
            <input
              type="text"
              name="reply_id"
              required
              placeholder="返信ID"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <button
              type="submit"
              class="px-4 py-2 text-sm font-medium text-red-600 border border-red-300 rounded-md hover:bg-red-50"
            >
              削除
            </button>
          </form>
        </div>

        <div :if={@last_result} class="mt-4 text-sm text-gray-600">結果: {@last_result}</div>
      </main>
    </div>
    """
  end
end
