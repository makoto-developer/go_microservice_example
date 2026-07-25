defmodule ShopMallWebWeb.Owner.InventoryLive do
  @moduledoc """
  加盟店の在庫管理画面。
  商品一覧の在庫をまとめて表示(BulkGetInventory)し、商品ごとに
  在庫詳細(GetInventoryByProduct)・入荷/調整(UpdateInventoryQuantity)・
  棚卸し(RecordStockTaking)・在庫アラート確認(CheckStockAlert)・
  変動履歴/棚卸し履歴の参照・在庫レコードの新規登録(RegisterInventory)を行う。
  """
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.ShopServiceClient, as: Shops
  alias ShopService.V1.{ShopService.Stub, ListProductsRequest}

  alias InventoryService.V1.{
    BulkGetInventoryRequest,
    CheckStockAlertRequest,
    GetInventoryByProductRequest,
    GetInventoryRequest,
    GetInventoryHistoryRequest,
    GetStockTakingHistoryRequest,
    RecordStockTakingRequest,
    RegisterInventoryRequest,
    UpdateInventoryQuantityRequest
  }

  alias InventoryService.V1.InventoryService.Stub, as: InventoryStub

  @impl true
  def mount(_params, session, socket) do
    owner_id = session["user_id"] || "admin-user-id"

    {:ok,
     socket
     |> assign(:owner_id, owner_id)
     |> assign(:shop, nil)
     |> assign(:products, [])
     |> assign(:stocks, %{})
     |> assign(:selected, nil)
     |> assign(:alert, nil)
     |> assign(:history, [])
     |> assign(:stock_takings, [])
     |> assign(:loading, true)
     |> assign(:error, nil)
     |> load_all()}
  end

  defp load_all(socket) do
    case Shops.get_shops_by_owner(socket.assigns.owner_id) do
      {:ok, %{shops: [shop | _]}} ->
        socket |> assign(:shop, shop) |> load_products_and_stocks()

      {:ok, _} ->
        socket |> assign(:loading, false) |> assign(:error, "店舗がありません")

      {:error, reason} ->
        socket |> assign(:loading, false) |> assign(:error, "店舗の取得に失敗しました: #{reason}")
    end
  end

  defp load_products_and_stocks(socket) do
    with {:ok, ch} <- shop_channel(),
         {:ok, resp} <-
           Stub.list_products(ch, %ListProductsRequest{
             shop_id: socket.assigns.shop.id,
             category: "",
             published_only: false,
             limit: 100,
             offset: 0
           }) do
      products = resp.products
      stocks = fetch_stocks(Enum.map(products, & &1.id))

      socket
      |> assign(:products, products)
      |> assign(:stocks, stocks)
      |> assign(:loading, false)
      |> assign(:error, nil)
    else
      {:error, reason} ->
        socket |> assign(:loading, false) |> assign(:error, "商品一覧の取得に失敗しました: #{inspect(reason)}")
    end
  end

  # 商品 ID の一覧から在庫をまとめて引く(在庫レコードが無い商品は含まれない)
  defp fetch_stocks(product_ids) do
    case inventory_call(fn ch ->
           InventoryStub.bulk_get_inventory(ch, %BulkGetInventoryRequest{product_ids: product_ids})
         end) do
      {:ok, resp} -> Map.new(resp.inventories, fn inv -> {inv.product_id, inv} end)
      {:error, _} -> %{}
    end
  end

  @impl true
  def handle_event("select_product", %{"product-id" => product_id}, socket) do
    with {:ok, resp} <-
           inventory_call(fn ch ->
             InventoryStub.get_inventory_by_product(ch, %GetInventoryByProductRequest{
               product_id: product_id
             })
           end) do
      inv = resp.inventory
      alert = fetch_alert(inv.id)
      history = fetch_history(inv.id)
      takings = fetch_stock_takings(socket.assigns.shop.id)

      {:noreply,
       socket
       |> assign(:selected, inv)
       |> assign(:alert, alert)
       |> assign(:history, history)
       |> assign(:stock_takings, takings)}
    else
      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "在庫の取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("refresh_inventory", _params, socket) do
    inv = socket.assigns.selected

    case inventory_call(fn ch ->
           InventoryStub.get_inventory(ch, %GetInventoryRequest{inventory_id: inv.id})
         end) do
      {:ok, resp} ->
        {:noreply,
         socket
         |> assign(:selected, resp.inventory)
         |> assign(:alert, fetch_alert(resp.inventory.id))}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "再取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("close_detail", _params, socket) do
    {:noreply, socket |> assign(:selected, nil) |> assign(:alert, nil)}
  end

  @impl true
  def handle_event("register_inventory", %{"product-id" => product_id, "quantity" => qty}, socket) do
    request = %RegisterInventoryRequest{
      product_id: product_id,
      shop_id: socket.assigns.shop.id,
      initial_quantity: String.to_integer(qty),
      alert_threshold: 10
    }

    case inventory_call(fn ch -> InventoryStub.register_inventory(ch, request) end) do
      {:ok, _resp} ->
        {:noreply,
         socket
         |> put_flash(:info, "在庫レコードを登録しました")
         |> load_products_and_stocks()}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "在庫登録に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("adjust_quantity", %{"change" => change, "reason" => reason}, socket) do
    inv = socket.assigns.selected

    request = %UpdateInventoryQuantityRequest{
      inventory_id: inv.id,
      shop_id: socket.assigns.shop.id,
      change_quantity: String.to_integer(change),
      reason: reason,
      operator: socket.assigns.owner_id
    }

    case inventory_call(fn ch -> InventoryStub.update_inventory_quantity(ch, request) end) do
      {:ok, resp} ->
        {:noreply,
         socket
         |> assign(:selected, resp.inventory)
         |> put_flash(:info, "在庫数を更新しました(現在 #{resp.inventory.quantity})")
         |> load_products_and_stocks()}

      {:error, err} ->
        {:noreply, put_flash(socket, :error, "在庫調整に失敗しました: #{err}")}
    end
  end

  @impl true
  def handle_event("record_stock_taking", %{"actual" => actual, "reason" => reason}, socket) do
    inv = socket.assigns.selected

    request = %RecordStockTakingRequest{
      inventory_id: inv.id,
      shop_id: socket.assigns.shop.id,
      actual_quantity: String.to_integer(actual),
      difference_reason: reason,
      operator: socket.assigns.owner_id
    }

    case inventory_call(fn ch -> InventoryStub.record_stock_taking(ch, request) end) do
      {:ok, resp} ->
        st = resp.stock_taking

        {:noreply,
         socket
         |> put_flash(:info, "棚卸しを記録しました(差分 #{st.difference})")
         |> assign(:selected, nil)
         |> load_products_and_stocks()}

      {:error, err} ->
        {:noreply, put_flash(socket, :error, "棚卸しの記録に失敗しました: #{err}")}
    end
  end

  defp fetch_alert(inventory_id) do
    case inventory_call(fn ch ->
           InventoryStub.check_stock_alert(ch, %CheckStockAlertRequest{inventory_id: inventory_id})
         end) do
      {:ok, resp} -> resp
      {:error, _} -> nil
    end
  end

  defp fetch_history(inventory_id) do
    case inventory_call(fn ch ->
           InventoryStub.get_inventory_history(ch, %GetInventoryHistoryRequest{
             inventory_id: inventory_id,
             page: 1,
             page_size: 20
           })
         end) do
      {:ok, resp} -> resp.history
      {:error, _} -> []
    end
  end

  defp fetch_stock_takings(shop_id) do
    case inventory_call(fn ch ->
           InventoryStub.get_stock_taking_history(ch, %GetStockTakingHistoryRequest{
             shop_id: shop_id
           })
         end) do
      {:ok, resp} -> resp.history
      {:error, _} -> []
    end
  end

  defp inventory_call(fun) do
    host = System.get_env("INVENTORY_SERVICE_HOST", "localhost")
    port = System.get_env("INVENTORY_SERVICE_PORT", "50054")

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
        {:error, "在庫サービスに接続できません: #{inspect(reason)}"}
    end
  end

  defp shop_channel do
    host = System.get_env("SHOP_SERVICE_HOST", "localhost")
    port = String.to_integer(System.get_env("SHOP_SERVICE_PORT", "22101"))
    GRPC.Stub.connect("#{host}:#{port}")
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
              <span class="text-white text-sm font-medium">在庫管理</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-5xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-6">
          在庫管理 <span :if={@shop} class="text-sm font-normal text-gray-500">({@shop.name})</span>
        </h1>

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
        <% else %>
          <div class="bg-white shadow rounded-lg overflow-hidden">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">商品</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    在庫数
                  </th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    引当済み
                  </th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">
                    販売可能
                  </th>
                  <th class="px-4 py-3"></th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <%= if @products == [] do %>
                  <tr>
                    <td colspan="5" class="px-4 py-8 text-center text-gray-500">商品がありません</td>
                  </tr>
                <% end %>
                <tr :for={product <- @products} class="hover:bg-gray-50">
                  <td class="px-4 py-3 text-sm text-gray-900">{product.name}</td>
                  <%= if inv = @stocks[product.id] do %>
                    <td class="px-4 py-3 text-sm text-right font-semibold">{inv.quantity}</td>
                    <td class="px-4 py-3 text-sm text-right text-gray-500">
                      {inv.reserved_quantity}
                    </td>
                    <td class="px-4 py-3 text-sm text-right">
                      <span class={
                        if inv.available_quantity < 10,
                          do: "text-red-600 font-semibold",
                          else: "text-gray-900"
                      }>
                        {inv.available_quantity}
                      </span>
                    </td>
                    <td class="px-4 py-3 text-right">
                      <button
                        phx-click="select_product"
                        phx-value-product-id={product.id}
                        class="text-emerald-700 hover:text-emerald-900 text-sm font-medium"
                      >
                        管理
                      </button>
                    </td>
                  <% else %>
                    <td colspan="3" class="px-4 py-3 text-sm text-center text-gray-400">在庫レコード未登録</td>
                    <td class="px-4 py-3 text-right">
                      <form
                        phx-submit="register_inventory"
                        class="flex justify-end items-center space-x-2"
                      >
                        <input type="hidden" name="product-id" value={product.id} />
                        <input
                          type="number"
                          name="quantity"
                          placeholder="初期数"
                          required
                          class="w-20 border border-gray-300 rounded-md px-2 py-1 text-sm"
                        />
                        <button
                          type="submit"
                          class="text-sm font-medium text-white bg-emerald-600 rounded-md px-3 py-1 hover:bg-emerald-700"
                        >
                          在庫登録
                        </button>
                      </form>
                    </td>
                  <% end %>
                </tr>
              </tbody>
            </table>
          </div>
        <% end %>
        
    <!-- 在庫詳細モーダル -->
        <div
          :if={@selected}
          class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50"
        >
          <div class="bg-white rounded-lg shadow-xl max-w-lg w-full mx-4 p-6 max-h-[85vh] overflow-y-auto">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-lg font-bold text-gray-900">在庫詳細</h2>
              <div class="space-x-2">
                <button
                  phx-click="refresh_inventory"
                  class="text-sm text-emerald-700 hover:text-emerald-900"
                  title="最新の在庫を取得"
                >
                  ↻ 更新
                </button>
                <button phx-click="close_detail" class="text-gray-400 hover:text-gray-600">✕</button>
              </div>
            </div>

            <div class="grid grid-cols-3 gap-3 mb-4 text-center">
              <div class="bg-gray-50 rounded-lg p-3">
                <div class="text-xs text-gray-500">在庫数</div>
                <div class="text-xl font-bold">{@selected.quantity}</div>
              </div>
              <div class="bg-gray-50 rounded-lg p-3">
                <div class="text-xs text-gray-500">引当済み</div>
                <div class="text-xl font-bold">{@selected.reserved_quantity}</div>
              </div>
              <div class="bg-gray-50 rounded-lg p-3">
                <div class="text-xs text-gray-500">販売可能</div>
                <div class="text-xl font-bold">{@selected.available_quantity}</div>
              </div>
            </div>

            <div
              :if={@alert}
              class={"px-3 py-2 rounded mb-4 text-sm " <>
                if(@alert.is_low_stock,
                  do: "bg-red-50 border border-red-300 text-red-700",
                  else: "bg-green-50 border border-green-300 text-green-700"
                )}
            >
              <%= if @alert.is_low_stock do %>
                ⚠ 在庫僅少: 販売可能 {@alert.current_quantity}(閾値 {@alert.alert_threshold})
              <% else %>
                ✓ 在庫は十分です(販売可能 {@alert.current_quantity} / 閾値 {@alert.alert_threshold})
              <% end %>
            </div>

            <h3 class="text-sm font-semibold text-gray-700 mb-2">入荷・調整</h3>
            <form phx-submit="adjust_quantity" class="flex items-end space-x-2 mb-4">
              <div>
                <label class="block text-xs text-gray-500 mb-1">増減数(負で減)</label>
                <input
                  type="number"
                  name="change"
                  required
                  placeholder="+50"
                  class="w-24 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                />
              </div>
              <div class="flex-1">
                <label class="block text-xs text-gray-500 mb-1">理由</label>
                <input
                  type="text"
                  name="reason"
                  placeholder="定期入荷"
                  class="w-full border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                />
              </div>
              <button
                type="submit"
                class="px-3 py-1.5 text-sm font-medium text-white bg-emerald-600 rounded-md hover:bg-emerald-700"
              >
                反映
              </button>
            </form>

            <h3 class="text-sm font-semibold text-gray-700 mb-2">棚卸し</h3>
            <form phx-submit="record_stock_taking" class="flex items-end space-x-2 mb-4">
              <div>
                <label class="block text-xs text-gray-500 mb-1">実数</label>
                <input
                  type="number"
                  name="actual"
                  required
                  placeholder="98"
                  class="w-24 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                />
              </div>
              <div class="flex-1">
                <label class="block text-xs text-gray-500 mb-1">差異理由</label>
                <input
                  type="text"
                  name="reason"
                  placeholder="破損 2 点"
                  class="w-full border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                />
              </div>
              <button
                type="submit"
                class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
              >
                記録
              </button>
            </form>

            <h3 class="text-sm font-semibold text-gray-700 mb-2">変動履歴</h3>
            <%= if @history == [] do %>
              <div class="text-xs text-gray-400 mb-3">履歴はありません</div>
            <% else %>
              <ul class="text-xs text-gray-600 space-y-1 mb-3">
                <li :for={h <- @history}>
                  {h.quantity_before} → {h.quantity_after}({h.reason})
                </li>
              </ul>
            <% end %>

            <h3 class="text-sm font-semibold text-gray-700 mb-2">棚卸し履歴</h3>
            <%= if @stock_takings == [] do %>
              <div class="text-xs text-gray-400">棚卸し記録はありません</div>
            <% else %>
              <ul class="text-xs text-gray-600 space-y-1">
                <li :for={st <- @stock_takings}>
                  実数 {st.actual_quantity}(差分 {st.difference}: {st.difference_reason})
                </li>
              </ul>
            <% end %>
          </div>
        </div>
      </main>
    </div>
    """
  end
end
