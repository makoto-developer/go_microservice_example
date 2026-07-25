defmodule ShopMallWebWeb.Owner.ProductFormLive do
  use ShopMallWebWeb, :live_view

  alias ShopMallWeb.ShopServiceClient, as: Shops

  alias ShopService.V1.{
    ShopService.Stub,
    RegisterProductRequest,
    UpdateProductRequest,
    GetProductRequest,
    ToggleProductPublishRequest,
    ProductVariationInput
  }

  @impl true
  def mount(params, _session, socket) do
    # TODO: セッションから実際のowner_idとshop_idを取得
    owner_id = "admin-user-id"
    shop_id = "11111111-1111-1111-1111-111111111111"

    product_id = Map.get(params, "id")
    mode = if product_id, do: :edit, else: :new

    socket =
      socket
      |> assign(:owner_id, owner_id)
      |> assign(:shop_id, shop_id)
      |> assign(:product_id, product_id)
      |> assign(:mode, mode)
      |> assign(:loading, false)
      |> assign(:error, nil)
      |> assign_form_defaults()
      |> allow_upload(:product_image,
        accept: ~w(.jpg .jpeg .png),
        max_entries: 1,
        max_file_size: 2_000_000
      )

    socket =
      if mode == :edit do
        load_product(socket)
      else
        socket
      end

    {:ok, socket}
  end

  defp assign_form_defaults(socket) do
    socket
    |> assign(:name, "")
    |> assign(:description, "")
    |> assign(:price, "")
    |> assign(:stock_quantity, "")
    |> assign(:category, "")
    |> assign(:published, false)
  end

  defp load_product(socket) do
    product_id = socket.assigns.product_id
    request = %GetProductRequest{product_id: product_id}

    case call_shop_service(:get_product, request) do
      {:ok, response} when not is_nil(response.product) ->
        product = response.product

        socket
        |> assign(:name, product.name)
        |> assign(:description, product.description)
        |> assign(:price, product.price)
        |> assign(:stock_quantity, to_string(product.stock_quantity))
        |> assign(:category, product.category)
        |> assign(:published, product.published)

      {:ok, _response} ->
        socket
        |> assign(:error, "商品が見つかりませんでした")
        |> assign(:loading, false)

      {:error, reason} ->
        socket
        |> assign(:error, "商品の読み込みに失敗しました: #{inspect(reason)}")
        |> assign(:loading, false)
    end
  end

  @impl true
  def handle_event("add_variation", params, socket) do
    variation = %ProductVariationInput{
      sku: params["sku"] || "",
      attribute_name: params["attribute_name"] || "",
      attribute_value: params["attribute_value"] || "",
      price: params["price"] || "0",
      stock_quantity: String.to_integer(params["stock_quantity"] || "0")
    }

    case Shops.manage_variation(socket.assigns.product_id, [variation]) do
      {:ok, response} ->
        {:noreply,
         put_flash(
           socket,
           :info,
           "バリエーションを登録しました(ID: #{Enum.join(response.variation_ids, ", ")})"
         )}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "バリエーション登録に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("validate_upload", _params, socket) do
    {:noreply, socket}
  end

  @impl true
  def handle_event("upload_image", _params, socket) do
    results =
      consume_uploaded_entries(socket, :product_image, fn %{path: path}, _entry ->
        {:ok, File.read!(path)}
      end)

    case results do
      [image_data | _] ->
        case Shops.upload_product_image(socket.assigns.product_id, image_data) do
          {:ok, response} ->
            {:noreply, put_flash(socket, :info, "商品画像をアップロードしました(#{response.url})")}

          {:error, reason} ->
            {:noreply, put_flash(socket, :error, "画像アップロードに失敗しました: #{reason}")}
        end

      [] ->
        {:noreply, put_flash(socket, :error, "画像ファイルを選択してください")}
    end
  end

  @impl true
  def handle_event("validate", %{"product" => product_params}, socket) do
    socket =
      socket
      |> assign(:name, product_params["name"] || "")
      |> assign(:description, product_params["description"] || "")
      |> assign(:price, product_params["price"] || "")
      |> assign(:stock_quantity, product_params["stock_quantity"] || "")
      |> assign(:category, product_params["category"] || "")
      |> assign(:published, product_params["published"] == "true")

    {:noreply, socket}
  end

  @impl true
  def handle_event("save", %{"product" => product_params}, socket) do
    IO.puts("=== SAVE EVENT RECEIVED ===")
    IO.inspect(product_params, label: "Product Params")

    name = product_params["name"] || ""
    description = product_params["description"] || ""
    price = product_params["price"] || "0"
    stock_quantity = String.to_integer(product_params["stock_quantity"] || "0")
    category = product_params["category"] || ""
    published = product_params["published"] == "true"

    IO.puts("=== PRODUCT DATA ===")
    IO.puts("Shop ID: #{socket.assigns.shop_id}")
    IO.puts("Name: #{name}")
    IO.puts("Published: #{published}")

    case socket.assigns.mode do
      :new ->
        request = %RegisterProductRequest{
          shop_id: socket.assigns.shop_id,
          name: name,
          description: description,
          price: price,
          stock_quantity: stock_quantity,
          category: category,
          weight: "",
          size: "",
          jan_code: "",
          tags: []
        }

        IO.puts("=== CALLING REGISTER PRODUCT ===")

        case call_shop_service(:register_product, request) do
          {:ok, response} ->
            IO.puts("=== REGISTER PRODUCT SUCCESS ===")
            IO.inspect(response, label: "Response")

            # 公開設定が有効な場合、ToggleProductPublishを呼ぶ
            if published do
              IO.puts("=== CALLING TOGGLE PUBLISH ===")

              toggle_request = %ToggleProductPublishRequest{
                product_id: response.product_id
              }

              call_shop_service(:toggle_product_publish, toggle_request)
            end

            {:noreply,
             socket
             |> put_flash(:info, "商品を登録しました")
             |> push_navigate(to: "/owner/products")}

          {:error, reason} ->
            IO.puts("=== REGISTER PRODUCT ERROR ===")
            IO.inspect(reason, label: "Error")

            {:noreply,
             socket
             |> put_flash(:error, "登録に失敗しました: #{inspect(reason)}")}
        end

      :edit ->
        request = %UpdateProductRequest{
          product_id: socket.assigns.product_id,
          name: name,
          description: description,
          price: price,
          stock_quantity: stock_quantity,
          category: category,
          weight: "",
          size: "",
          jan_code: ""
        }

        case call_shop_service(:update_product, request) do
          {:ok, _response} ->
            # 公開設定が現在の状態と異なる場合、ToggleProductPublishを呼ぶ
            if published != socket.assigns.published do
              toggle_request = %ToggleProductPublishRequest{
                product_id: socket.assigns.product_id
              }

              call_shop_service(:toggle_product_publish, toggle_request)
            end

            {:noreply,
             socket
             |> put_flash(:info, "商品を更新しました")
             |> push_navigate(to: "/owner/products")}

          {:error, reason} ->
            {:noreply,
             socket
             |> put_flash(:error, "更新に失敗しました: #{inspect(reason)}")}
        end
    end
  end

  defp call_shop_service(:get_product, request) do
    case get_shop_channel() do
      {:ok, channel} ->
        case Stub.get_product(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, error}
        end

      {:error, reason} ->
        {:error, "Shop Serviceに接続できません: #{inspect(reason)}"}
    end
  end

  defp call_shop_service(:register_product, request) do
    case get_shop_channel() do
      {:ok, channel} ->
        case Stub.register_product(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, error}
        end

      {:error, reason} ->
        {:error, "Shop Serviceに接続できません: #{inspect(reason)}"}
    end
  end

  defp call_shop_service(:update_product, request) do
    case get_shop_channel() do
      {:ok, channel} ->
        case Stub.update_product(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, error}
        end

      {:error, reason} ->
        {:error, "Shop Serviceに接続できません: #{inspect(reason)}"}
    end
  end

  defp call_shop_service(:toggle_product_publish, request) do
    case get_shop_channel() do
      {:ok, channel} ->
        case Stub.toggle_product_publish(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, error}
        end

      {:error, reason} ->
        {:error, "Shop Serviceに接続できません: #{inspect(reason)}"}
    end
  end

  defp get_shop_channel do
    host = System.get_env("SHOP_SERVICE_HOST", "localhost")
    port = String.to_integer(System.get_env("SHOP_SERVICE_PORT", "22101"))

    case GRPC.Stub.connect("#{host}:#{port}") do
      {:ok, channel} -> {:ok, channel}
      {:error, reason} -> {:error, reason}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <!-- ナビゲーションバー -->
      <nav class="bg-blue-800 shadow-md">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16">
            <div class="flex items-center">
              <.link
                navigate="/owner/dashboard"
                class="text-2xl font-bold text-white hover:text-gray-200"
              >
                ショップ管理
              </.link>
            </div>
            <div class="flex items-center space-x-4">
              <.link
                navigate="/owner/dashboard"
                class="text-white hover:text-gray-200 px-3 py-2 rounded-md text-sm font-medium"
              >
                ダッシュボード
              </.link>
              <.link
                navigate="/owner/products"
                class="text-white hover:text-gray-200 px-3 py-2 rounded-md text-sm font-medium border-b-2 border-white"
              >
                商品管理
              </.link>
              <.link
                navigate="/dashboard"
                class="text-white hover:text-gray-200 px-3 py-2 rounded-md text-sm font-medium"
              >
                顧客画面へ
              </.link>
            </div>
          </div>
        </div>
      </nav>
      
    <!-- メインコンテンツ -->
      <main class="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <div class="px-4 py-6 sm:px-0">
          <div class="bg-white shadow rounded-lg">
            <div class="px-4 py-5 sm:p-6">
              <h1 class="text-2xl font-bold text-gray-900 mb-6">
                {if @mode == :new, do: "新規商品登録", else: "商品編集"}
              </h1>

              <%= if @error do %>
                <div class="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
                  {@error}
                </div>
              <% end %>

              <form phx-change="validate" phx-submit="save" class="space-y-6">
                <!-- 商品名 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    商品名 <span class="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    name="product[name]"
                    value={@name}
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="例: ワイヤレスイヤホン Pro"
                    required
                  />
                </div>
                
    <!-- 説明 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    説明 <span class="text-red-500">*</span>
                  </label>
                  <textarea
                    name="product[description]"
                    rows="4"
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="商品の説明を入力してください"
                    required
                  ><%= @description %></textarea>
                </div>
                
    <!-- 価格 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    価格（円） <span class="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    name="product[price]"
                    value={@price}
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="例: 29800"
                    required
                  />
                  <p class="mt-1 text-sm text-gray-500">小数点も入力可能です（例: 2980.50）</p>
                </div>
                
    <!-- 在庫数 -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    在庫数 <span class="text-red-500">*</span>
                  </label>
                  <input
                    type="number"
                    name="product[stock_quantity]"
                    value={@stock_quantity}
                    min="0"
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="例: 100"
                    required
                  />
                </div>
                
    <!-- カテゴリ -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    カテゴリ <span class="text-red-500">*</span>
                  </label>
                  <select
                    name="product[category]"
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    required
                  >
                    <option value="">選択してください</option>
                    <option value="electronics" selected={@category == "electronics"}>
                      電化製品
                    </option>
                    <option value="fashion" selected={@category == "fashion"}>
                      ファッション
                    </option>
                    <option value="food" selected={@category == "food"}>
                      食品
                    </option>
                    <option value="books" selected={@category == "books"}>
                      書籍
                    </option>
                    <option value="sports" selected={@category == "sports"}>
                      スポーツ
                    </option>
                    <option value="home" selected={@category == "home"}>
                      ホーム・生活用品
                    </option>
                    <option value="toys" selected={@category == "toys"}>
                      おもちゃ・ホビー
                    </option>
                    <option value="other" selected={@category == "other"}>
                      その他
                    </option>
                  </select>
                </div>
                
    <!-- 公開設定 -->
                <div class="flex items-center">
                  <input
                    type="checkbox"
                    name="product[published]"
                    value="true"
                    checked={@published}
                    class="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <label class="ml-2 block text-sm text-gray-900">
                    この商品を公開する
                  </label>
                </div>
                
    <!-- ボタン -->
                <div class="flex justify-end space-x-3 pt-6 border-t">
                  <.link
                    navigate="/owner/products"
                    class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
                  >
                    キャンセル
                  </.link>
                  <button
                    type="submit"
                    phx-disable-with="処理中..."
                    class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
                  >
                    {if @mode == :new, do: "登録", else: "更新"}
                  </button>
                </div>
              </form>

              <%= if @mode == :edit do %>
                <!-- バリエーション管理(ManageVariation) -->
                <div class="mt-8 border-t pt-6">
                  <h2 class="text-lg font-semibold text-gray-900 mb-3">バリエーションを追加</h2>
                  <form phx-submit="add_variation" class="grid grid-cols-2 gap-3">
                    <div>
                      <label class="block text-xs text-gray-500 mb-1">SKU</label>
                      <input
                        type="text"
                        name="sku"
                        required
                        placeholder="TSHIRT-RED-M"
                        class="w-full border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                      />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-500 mb-1">価格(円)</label>
                      <input
                        type="number"
                        name="price"
                        required
                        placeholder="2980"
                        class="w-full border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                      />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-500 mb-1">属性名</label>
                      <input
                        type="text"
                        name="attribute_name"
                        required
                        placeholder="サイズ"
                        class="w-full border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                      />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-500 mb-1">属性値</label>
                      <input
                        type="text"
                        name="attribute_value"
                        required
                        placeholder="M"
                        class="w-full border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                      />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-500 mb-1">在庫数</label>
                      <input
                        type="number"
                        name="stock_quantity"
                        required
                        placeholder="10"
                        class="w-full border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                      />
                    </div>
                    <div class="flex items-end">
                      <button
                        type="submit"
                        class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                      >
                        バリエーション登録
                      </button>
                    </div>
                  </form>
                </div>
                
    <!-- 商品画像アップロード(UploadProductImage) -->
                <div class="mt-8 border-t pt-6">
                  <h2 class="text-lg font-semibold text-gray-900 mb-3">商品画像</h2>
                  <form phx-submit="upload_image" phx-change="validate_upload" class="space-y-3">
                    <.live_file_input upload={@uploads.product_image} class="text-sm" />
                    <div :for={entry <- @uploads.product_image.entries} class="text-xs text-gray-500">
                      {entry.client_name}({div(entry.client_size, 1024)} KB)
                    </div>
                    <button
                      type="submit"
                      class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                    >
                      アップロード
                    </button>
                  </form>
                </div>
              <% end %>
            </div>
          </div>
        </div>
      </main>
    </div>
    """
  end
end
