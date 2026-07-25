defmodule ShopMallWebWeb.MyPageLive do
  @moduledoc """
  顧客のマイページ。タブでプロフィール・住所・支払い方法・お気に入り・
  レビュー・注文履歴を切り替える(すべて customer サービスの RPC を使用)。
  """
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.CustomerServiceClient, as: Customers

  @impl true
  def mount(_params, session, socket) do
    customer_id = session["user_id"]

    {:ok,
     socket
     |> assign(:customer_id, customer_id)
     |> assign(:tab, "profile")
     |> assign(:profile, nil)
     |> assign(:address, nil)
     |> assign(:postal_hint, nil)
     |> assign(:payment_method, nil)
     |> assign(:favorites, [])
     |> assign(:reviews, [])
     |> assign(:orders, [])
     |> allow_upload(:profile_image,
       accept: ~w(.jpg .jpeg .png),
       max_entries: 1,
       max_file_size: 2_000_000
     )
     |> load_tab("profile")}
  end

  defp load_tab(%{assigns: %{customer_id: nil}} = socket, _tab), do: socket

  defp load_tab(socket, "profile") do
    case Customers.get_profile(socket.assigns.customer_id) do
      {:ok, resp} -> assign(socket, :profile, resp.customer)
      {:error, _} -> socket
    end
  end

  defp load_tab(socket, "favorites") do
    case Customers.get_favorites(socket.assigns.customer_id) do
      {:ok, resp} -> assign(socket, :favorites, resp.favorites)
      {:error, _} -> assign(socket, :favorites, [])
    end
  end

  defp load_tab(socket, "reviews") do
    case Customers.get_my_reviews(socket.assigns.customer_id) do
      {:ok, resp} -> assign(socket, :reviews, resp.reviews)
      {:error, _} -> assign(socket, :reviews, [])
    end
  end

  defp load_tab(socket, "orders") do
    case Customers.get_order_history(socket.assigns.customer_id) do
      {:ok, resp} -> assign(socket, :orders, resp.orders)
      {:error, _} -> assign(socket, :orders, [])
    end
  end

  defp load_tab(socket, _), do: socket

  @impl true
  def handle_event("set_tab", %{"tab" => tab}, socket) do
    {:noreply, socket |> assign(:tab, tab) |> load_tab(tab)}
  end

  # ---- プロフィール ----

  @impl true
  def handle_event("save_profile", params, socket) do
    case Customers.update_profile(%{
           customer_id: socket.assigns.customer_id,
           first_name: params["first_name"],
           last_name: params["last_name"],
           phone: params["phone"],
           birth_date: params["birth_date"]
         }) do
      {:ok, resp} ->
        {:noreply, socket |> assign(:profile, resp.customer) |> put_flash(:info, "プロフィールを更新しました")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "更新に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("validate_upload", _params, socket), do: {:noreply, socket}

  @impl true
  def handle_event("upload_profile_image", _params, socket) do
    results =
      consume_uploaded_entries(socket, :profile_image, fn %{path: path}, _entry ->
        {:ok, File.read!(path)}
      end)

    case results do
      [data | _] ->
        case Customers.upload_profile_image(socket.assigns.customer_id, data) do
          {:ok, resp} ->
            {:noreply, put_flash(socket, :info, "画像をアップロードしました(#{resp.image_url})")}

          {:error, reason} ->
            {:noreply, put_flash(socket, :error, "アップロードに失敗しました: #{reason}")}
        end

      [] ->
        {:noreply, put_flash(socket, :error, "画像を選択してください")}
    end
  end

  # ---- 住所 ----

  @impl true
  def handle_event("search_postal", %{"postal_code" => code}, socket) do
    case Customers.search_postal_code(code) do
      {:ok, resp} ->
        {:noreply, assign(socket, :postal_hint, resp)}

      {:error, reason} ->
        {:noreply,
         socket |> assign(:postal_hint, nil) |> put_flash(:error, "住所が見つかりません: #{reason}")}
    end
  end

  @impl true
  def handle_event("save_address", params, socket) do
    attrs = %{
      customer_id: socket.assigns.customer_id,
      address_name: params["address_name"],
      postal_code: params["postal_code"],
      prefecture: params["prefecture"],
      city: params["city"],
      address_line1: params["address_line1"],
      recipient_name: params["recipient_name"],
      recipient_phone: params["recipient_phone"]
    }

    result =
      if socket.assigns.address do
        Customers.update_address(Map.put(attrs, :address_id, socket.assigns.address.id))
      else
        Customers.register_address(attrs)
      end

    case result do
      {:ok, resp} ->
        {:noreply,
         socket
         |> assign(:address, resp.address)
         |> put_flash(:info, "お届け先を保存しました")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "保存に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("delete_address", _params, socket) do
    case Customers.delete_address(socket.assigns.address.id, socket.assigns.customer_id) do
      {:ok, _} ->
        {:noreply, socket |> assign(:address, nil) |> put_flash(:info, "お届け先を削除しました")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "削除に失敗しました: #{reason}")}
    end
  end

  # ---- 支払い方法 ----

  @impl true
  def handle_event("register_payment", %{"pm_id" => pm_id}, socket) do
    case Customers.register_payment_method(socket.assigns.customer_id, pm_id) do
      {:ok, resp} ->
        {:noreply,
         socket
         |> assign(:payment_method, resp.payment_method)
         |> put_flash(:info, "支払い方法を登録しました")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "登録に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("delete_payment", _params, socket) do
    case Customers.delete_payment_method(
           socket.assigns.payment_method.id,
           socket.assigns.customer_id
         ) do
      {:ok, _} ->
        {:noreply, socket |> assign(:payment_method, nil) |> put_flash(:info, "支払い方法を削除しました")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "削除に失敗しました: #{reason}")}
    end
  end

  # ---- お気に入り ----

  @impl true
  def handle_event("remove_favorite", %{"favorite-id" => id}, socket) do
    case Customers.remove_from_favorite(id, socket.assigns.customer_id) do
      {:ok, _} ->
        {:noreply, socket |> put_flash(:info, "お気に入りから削除しました") |> load_tab("favorites")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "削除に失敗しました: #{reason}")}
    end
  end

  # ---- レビュー ----

  @impl true
  def handle_event("update_review", params, socket) do
    case Customers.update_review(%{
           review_id: params["review_id"],
           customer_id: socket.assigns.customer_id,
           rating: String.to_integer(params["rating"]),
           review_text: params["review_text"]
         }) do
      {:ok, _} ->
        {:noreply, socket |> put_flash(:info, "レビューを更新しました") |> load_tab("reviews")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "更新に失敗しました: #{reason}")}
    end
  end

  # ---- 注文履歴 ----

  @impl true
  def handle_event("order_detail", %{"order-id" => order_id}, socket) do
    case Customers.get_order_detail(order_id, socket.assigns.customer_id) do
      {:ok, resp} ->
        o = resp.order
        {:noreply, put_flash(socket, :info, "#{o.order_number}: #{o.status}(¥#{o.total_amount})")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "詳細の取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("cancel_order", %{"order-id" => order_id}, socket) do
    case Customers.request_order_cancel(order_id, socket.assigns.customer_id, "customer request") do
      {:ok, resp} ->
        {:noreply, socket |> put_flash(:info, resp.message) |> load_tab("orders")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "キャンセルに失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("reorder", %{"order-id" => order_id}, socket) do
    case Customers.reorder_from_history(order_id, socket.assigns.customer_id) do
      {:ok, resp} ->
        {:noreply, socket |> put_flash(:info, resp.message) |> load_tab("orders")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "再注文に失敗しました: #{reason}")}
    end
  end

  defp tab_class(current, tab) do
    if current == tab do
      "px-3 py-1.5 rounded-full text-sm font-medium bg-blue-600 text-white"
    else
      "px-3 py-1.5 rounded-full text-sm font-medium bg-white text-gray-600 border border-gray-300 hover:bg-gray-50"
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
                navigate="/cart"
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                カート
              </.link>
              <span class="text-gray-900 px-3 py-2 text-sm font-semibold">マイページ</span>
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

      <main class="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-4">マイページ</h1>

        <%= if is_nil(@customer_id) do %>
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
            マイページにはログインが必要です
          </div>
        <% else %>
          <div class="flex flex-wrap gap-2 mb-6">
            <button
              :for={
                {tab, label} <- [
                  {"profile", "プロフィール"},
                  {"address", "お届け先"},
                  {"payment", "支払い方法"},
                  {"favorites", "お気に入り"},
                  {"reviews", "レビュー"},
                  {"orders", "注文履歴"}
                ]
              }
              phx-click="set_tab"
              phx-value-tab={tab}
              class={tab_class(@tab, tab)}
            >
              {label}
            </button>
          </div>

          <%= case @tab do %>
            <% "profile" -> %>
              <div class="bg-white shadow rounded-lg p-6">
                <h2 class="text-lg font-semibold text-gray-900 mb-4">プロフィール</h2>
                <form phx-submit="save_profile" class="grid grid-cols-2 gap-4 mb-6">
                  <div>
                    <label class="block text-xs text-gray-500 mb-1">姓</label>
                    <input
                      type="text"
                      name="last_name"
                      value={@profile && @profile.last_name}
                      class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                    />
                  </div>
                  <div>
                    <label class="block text-xs text-gray-500 mb-1">名</label>
                    <input
                      type="text"
                      name="first_name"
                      value={@profile && @profile.first_name}
                      class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                    />
                  </div>
                  <div>
                    <label class="block text-xs text-gray-500 mb-1">電話番号</label>
                    <input
                      type="tel"
                      name="phone"
                      value={@profile && @profile.phone}
                      class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                    />
                  </div>
                  <div>
                    <label class="block text-xs text-gray-500 mb-1">生年月日</label>
                    <input
                      type="date"
                      name="birth_date"
                      value={@profile && @profile.birth_date}
                      class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                    />
                  </div>
                  <div class="col-span-2 flex justify-end">
                    <button
                      type="submit"
                      class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                    >
                      保存
                    </button>
                  </div>
                </form>

                <h3 class="text-sm font-semibold text-gray-700 mb-2">プロフィール画像</h3>
                <form
                  phx-submit="upload_profile_image"
                  phx-change="validate_upload"
                  class="flex items-center space-x-3"
                >
                  <.live_file_input upload={@uploads.profile_image} class="text-sm" />
                  <button
                    type="submit"
                    class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
                  >
                    アップロード
                  </button>
                </form>
              </div>
            <% "address" -> %>
              <div class="bg-white shadow rounded-lg p-6">
                <div class="flex items-center justify-between mb-4">
                  <h2 class="text-lg font-semibold text-gray-900">お届け先</h2>
                  <button
                    :if={@address}
                    phx-click="delete_address"
                    class="text-sm text-red-600 hover:text-red-800"
                  >
                    削除
                  </button>
                </div>
                <form phx-submit="save_address" class="grid grid-cols-2 gap-3">
                  <div>
                    <label class="block text-xs text-gray-500 mb-1">郵便番号</label>
                    <div class="flex space-x-2">
                      <input
                        type="text"
                        name="postal_code"
                        required
                        placeholder="1000001"
                        value={@address && @address.postal_code}
                        class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
                      />
                    </div>
                  </div>
                  <div class="flex items-end">
                    <button
                      type="button"
                      phx-click={Phoenix.LiveView.JS.dispatch("submit", to: "#postal-form")}
                      class="hidden"
                    >
                    </button>
                  </div>
                  <div>
                    <label class="block text-xs text-gray-500 mb-1">都道府県</label>
                    <input
                      type="text"
                      name="prefecture"
                      required
                      value={
                        (@postal_hint && @postal_hint.prefecture) || (@address && @address.prefecture)
                      }
                      class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                    />
                  </div>
                  <div>
                    <label class="block text-xs text-gray-500 mb-1">市区町村</label>
                    <input
                      type="text"
                      name="city"
                      required
                      value={(@postal_hint && @postal_hint.city) || (@address && @address.city)}
                      class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                    />
                  </div>
                  <div class="col-span-2">
                    <label class="block text-xs text-gray-500 mb-1">番地・建物名</label>
                    <input
                      type="text"
                      name="address_line1"
                      required
                      value={
                        (@postal_hint && @postal_hint.address_line1) ||
                          (@address && @address.address_line1)
                      }
                      class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                    />
                  </div>
                  <div>
                    <label class="block text-xs text-gray-500 mb-1">宛名</label>
                    <input
                      type="text"
                      name="recipient_name"
                      required
                      value={@address && @address.recipient_name}
                      class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                    />
                  </div>
                  <div>
                    <label class="block text-xs text-gray-500 mb-1">電話番号</label>
                    <input
                      type="tel"
                      name="recipient_phone"
                      value={@address && @address.recipient_phone}
                      class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                    />
                  </div>
                  <input type="hidden" name="address_name" value="自宅" />
                  <div class="col-span-2 flex justify-between items-center">
                    <button
                      type="button"
                      phx-click="search_postal"
                      phx-value-postal_code={(@address && @address.postal_code) || ""}
                      class="text-sm text-blue-600 hover:text-blue-800"
                      title="郵便番号から住所を補完(登録済みの郵便番号を使用)"
                    >
                      🔍 郵便番号から住所検索
                    </button>
                    <button
                      type="submit"
                      class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                    >
                      {if @address, do: "更新", else: "登録"}
                    </button>
                  </div>
                </form>
              </div>
            <% "payment" -> %>
              <div class="bg-white shadow rounded-lg p-6">
                <h2 class="text-lg font-semibold text-gray-900 mb-4">支払い方法</h2>
                <%= if @payment_method do %>
                  <div class="flex items-center justify-between border border-gray-200 rounded-lg p-4 mb-4">
                    <div class="text-sm">
                      <span class="font-semibold">{@payment_method.card_brand}</span>
                      **** {@payment_method.card_last4}
                    </div>
                    <button phx-click="delete_payment" class="text-sm text-red-600 hover:text-red-800">
                      削除
                    </button>
                  </div>
                <% end %>
                <form phx-submit="register_payment" class="flex items-end space-x-2">
                  <div class="flex-1">
                    <label class="block text-xs text-gray-500 mb-1">
                      カードトークン(デモ: pm_demo_visa など)
                    </label>
                    <input
                      type="text"
                      name="pm_id"
                      required
                      placeholder="pm_demo_visa"
                      class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                    />
                  </div>
                  <button
                    type="submit"
                    class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                  >
                    カードを登録
                  </button>
                </form>
              </div>
            <% "favorites" -> %>
              <div class="bg-white shadow rounded-lg p-6">
                <h2 class="text-lg font-semibold text-gray-900 mb-4">お気に入り</h2>
                <%= if @favorites == [] do %>
                  <div class="text-sm text-gray-400 py-4 text-center">お気に入りはまだありません</div>
                <% end %>
                <ul class="divide-y divide-gray-100">
                  <li :for={fav <- @favorites} class="py-3 flex items-center justify-between">
                    <.link
                      navigate={"/products/#{fav.product_id}"}
                      class="text-sm text-blue-600 hover:text-blue-800 font-mono"
                    >
                      {String.slice(fav.product_id, 0, 8)}…
                    </.link>
                    <button
                      phx-click="remove_favorite"
                      phx-value-favorite-id={fav.id}
                      class="text-sm text-red-600 hover:text-red-800"
                    >
                      削除
                    </button>
                  </li>
                </ul>
              </div>
            <% "reviews" -> %>
              <div class="bg-white shadow rounded-lg p-6">
                <h2 class="text-lg font-semibold text-gray-900 mb-4">投稿したレビュー</h2>
                <%= if @reviews == [] do %>
                  <div class="text-sm text-gray-400 py-4 text-center">レビューはまだありません</div>
                <% end %>
                <div :for={review <- @reviews} class="border border-gray-200 rounded-lg p-4 mb-3">
                  <form phx-submit="update_review" class="space-y-2">
                    <input type="hidden" name="review_id" value={review.id} />
                    <div class="flex items-center space-x-2">
                      <select
                        name="rating"
                        class="border border-gray-300 rounded-md px-2 py-1 text-sm"
                      >
                        <option :for={n <- 5..1//-1} value={n} selected={review.rating == n}>
                          {String.duplicate("★", n)}
                        </option>
                      </select>
                      <span class="text-xs text-gray-400 font-mono">
                        {String.slice(review.product_id, 0, 8)}…
                      </span>
                    </div>
                    <textarea
                      name="review_text"
                      rows="2"
                      class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                    >{review.review_text}</textarea>
                    <div class="flex justify-end">
                      <button
                        type="submit"
                        class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
                      >
                        更新
                      </button>
                    </div>
                  </form>
                </div>
              </div>
            <% "orders" -> %>
              <div class="bg-white shadow rounded-lg p-6">
                <h2 class="text-lg font-semibold text-gray-900 mb-4">注文履歴(customer サービス経由)</h2>
                <%= if @orders == [] do %>
                  <div class="text-sm text-gray-400 py-4 text-center">注文はまだありません</div>
                <% end %>
                <ul class="divide-y divide-gray-100">
                  <li :for={order <- @orders} class="py-3 flex items-center justify-between">
                    <div class="text-sm">
                      <span class="font-mono text-gray-600">{order.order_number}</span>
                      <span class="ml-2 text-gray-500">{order.status}</span>
                      <span class="ml-2 font-semibold">¥{order.total_amount}</span>
                    </div>
                    <div class="space-x-2">
                      <button
                        phx-click="order_detail"
                        phx-value-order-id={order.order_id}
                        class="text-sm text-blue-600 hover:text-blue-800"
                      >
                        詳細
                      </button>
                      <button
                        phx-click="reorder"
                        phx-value-order-id={order.order_id}
                        class="text-sm text-gray-600 hover:text-gray-800"
                      >
                        再注文
                      </button>
                      <button
                        phx-click="cancel_order"
                        phx-value-order-id={order.order_id}
                        class="text-sm text-red-600 hover:text-red-800"
                      >
                        キャンセル
                      </button>
                    </div>
                  </li>
                </ul>
              </div>
            <% _ -> %>
          <% end %>
        <% end %>
      </main>
    </div>
    """
  end
end
