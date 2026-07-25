defmodule ShopMallWebWeb.SearchLive do
  @moduledoc """
  検索ページ。search サービスの検索系 RPC を使用する:
  商品検索(SearchProducts + RecordSearchHistory)・店舗検索(SearchShops)・
  サジェスト(GetSearchSuggestions)・検索履歴(Get/Delete/ClearAllSearchHistory)・
  人気/急上昇キーワード(GetPopularKeywords / GetTrendingKeywords)。
  """
  use ShopMallWebWeb, :live_view

  alias SearchService.V1.SearchService.Stub
  alias SearchService.V1, as: PB

  @impl true
  def mount(_params, session, socket) do
    {:ok,
     socket
     |> assign(:user_id, session["user_id"] || "guest")
     |> assign(:keyword, "")
     |> assign(:mode, "products")
     |> assign(:result_note, nil)
     |> assign(:suggestions_note, nil)
     |> assign(:history_note, nil)
     |> assign(:popular_note, nil)
     |> assign(:trending_note, nil)
     |> load_keywords()}
  end

  defp load_keywords(socket) do
    popular =
      case call(fn ch ->
             Stub.get_popular_keywords(ch, %PB.GetPopularKeywordsRequest{
               period_type: "weekly",
               limit: 10
             })
           end) do
        {:ok, resp} -> resp.message
        {:error, _} -> nil
      end

    trending =
      case call(fn ch ->
             Stub.get_trending_keywords(ch, %PB.GetTrendingKeywordsRequest{limit: 10})
           end) do
        {:ok, resp} -> resp.message
        {:error, _} -> nil
      end

    socket |> assign(:popular_note, popular) |> assign(:trending_note, trending)
  end

  @impl true
  def handle_event("search", %{"keyword" => keyword, "mode" => mode}, socket) do
    result =
      case mode do
        "shops" ->
          call(fn ch ->
            Stub.search_shops(ch, %PB.SearchShopsRequest{keyword: keyword, page: 1, page_size: 20})
          end)

        _ ->
          call(fn ch ->
            Stub.search_products(ch, %PB.SearchProductsRequest{
              keyword: keyword,
              page: 1,
              page_size: 20
            })
          end)
      end

    # 検索履歴を記録する(best effort)
    call(fn ch ->
      Stub.record_search_history(ch, %PB.RecordSearchHistoryRequest{
        user_id: socket.assigns.user_id,
        keyword: keyword,
        result_count: 0
      })
    end)

    note =
      case result do
        {:ok, resp} -> resp.message
        {:error, reason} -> "検索エラー: #{reason}"
      end

    {:noreply,
     socket |> assign(:keyword, keyword) |> assign(:mode, mode) |> assign(:result_note, note)}
  end

  @impl true
  def handle_event("suggest", %{"keyword" => prefix}, socket) do
    case call(fn ch ->
           Stub.get_search_suggestions(ch, %PB.GetSearchSuggestionsRequest{
             prefix: prefix,
             limit: 5
           })
         end) do
      {:ok, resp} -> {:noreply, assign(socket, :suggestions_note, resp.message)}
      {:error, _} -> {:noreply, socket}
    end
  end

  @impl true
  def handle_event("load_history", _params, socket) do
    case call(fn ch ->
           Stub.get_search_history(ch, %PB.GetSearchHistoryRequest{
             user_id: socket.assigns.user_id,
             limit: 20
           })
         end) do
      {:ok, resp} -> {:noreply, assign(socket, :history_note, resp.message)}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "履歴の取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("delete_history", %{"history_id" => history_id}, socket) do
    case call(fn ch ->
           Stub.delete_search_history(ch, %PB.DeleteSearchHistoryRequest{
             user_id: socket.assigns.user_id,
             history_id: history_id
           })
         end) do
      {:ok, resp} -> {:noreply, put_flash(socket, :info, resp.message || "削除しました")}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "削除に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("clear_history", _params, socket) do
    case call(fn ch ->
           Stub.clear_all_search_history(ch, %PB.ClearAllSearchHistoryRequest{
             user_id: socket.assigns.user_id
           })
         end) do
      {:ok, resp} ->
        {:noreply,
         socket |> assign(:history_note, nil) |> put_flash(:info, resp.message || "履歴を全削除しました")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "全削除に失敗しました: #{reason}")}
    end
  end

  defp call(fun) do
    host = System.get_env("SEARCH_SERVICE_HOST", "localhost")
    port = System.get_env("SEARCH_SERVICE_PORT", "20110")

    case GRPC.Stub.connect("#{host}:#{port}") do
      {:ok, channel} ->
        try do
          case fun.(channel) do
            {:ok, response} -> {:ok, response}
            {:error, %GRPC.RPCError{message: message}} -> {:error, message}
            {:error, reason} -> {:error, inspect(reason)}
          end
        after
          GRPC.Stub.disconnect(channel)
        end

      {:error, reason} ->
        {:error, "検索サービスに接続できません: #{inspect(reason)}"}
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
              <span class="text-gray-900 px-3 py-2 text-sm font-semibold">検索</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-6">検索</h1>

        <form phx-submit="search" phx-change="suggest" class="bg-white shadow rounded-lg p-4 mb-4">
          <div class="flex space-x-2">
            <input
              type="text"
              name="keyword"
              value={@keyword}
              placeholder="キーワードを入力"
              autocomplete="off"
              class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
            />
            <select name="mode" class="border border-gray-300 rounded-md px-2 py-2 text-sm">
              <option value="products" selected={@mode == "products"}>商品</option>
              <option value="shops" selected={@mode == "shops"}>店舗</option>
            </select>
            <button
              type="submit"
              class="px-5 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
            >
              検索
            </button>
          </div>
          <div :if={@suggestions_note} class="text-xs text-gray-500 mt-2">
            候補: {@suggestions_note}
          </div>
        </form>

        <div :if={@result_note} class="bg-white shadow rounded-lg p-4 mb-4 text-sm text-gray-700">
          {@result_note}
        </div>

        <div class="grid grid-cols-2 gap-4 mb-4">
          <div class="bg-white shadow rounded-lg p-4">
            <h2 class="text-sm font-semibold text-gray-700 mb-2">🔥 人気キーワード</h2>
            <div class="text-xs text-gray-600">{@popular_note || "-"}</div>
          </div>
          <div class="bg-white shadow rounded-lg p-4">
            <h2 class="text-sm font-semibold text-gray-700 mb-2">📈 急上昇</h2>
            <div class="text-xs text-gray-600">{@trending_note || "-"}</div>
          </div>
        </div>

        <div class="bg-white shadow rounded-lg p-4">
          <div class="flex items-center justify-between mb-2">
            <h2 class="text-sm font-semibold text-gray-700">検索履歴</h2>
            <div class="space-x-3">
              <button phx-click="load_history" class="text-xs text-blue-600 hover:text-blue-800">
                表示
              </button>
              <button phx-click="clear_history" class="text-xs text-red-600 hover:text-red-800">
                全削除
              </button>
            </div>
          </div>
          <div :if={@history_note} class="text-xs text-gray-600 mb-2">{@history_note}</div>
          <form phx-submit="delete_history" class="flex items-end space-x-2">
            <input
              type="text"
              name="history_id"
              required
              placeholder="履歴IDを指定して削除"
              class="flex-1 border border-gray-300 rounded-md px-2 py-1.5 text-xs"
            />
            <button
              type="submit"
              class="px-3 py-1.5 text-xs font-medium text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              1件削除
            </button>
          </form>
        </div>
      </main>
    </div>
    """
  end
end
