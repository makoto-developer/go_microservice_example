defmodule ShopMallWebWeb.Owner.ShopRegisterLive do
  use ShopMallWebWeb, :live_view

  alias ShopService.V1.{
    ShopService.Stub,
    RegisterShopRequest
  }

  @impl true
  def mount(_params, session, socket) do
    # TODO: セッションからユーザーID取得
    # 今は仮のUUIDを使用
    owner_id = session["user_id"] || "00000000-0000-0000-0000-000000000000"

    {:ok,
     socket
     |> assign(:owner_id, owner_id)
     |> assign(:error, nil)
     |> assign(:step, 1)
     |> assign(:form_data, %{
       name: "",
       description: "",
       owner_name: "",
       phone_number: "",
       business_hours: "",
       return_policy: "",
       categories: []
     })}
  end

  @impl true
  def handle_event("submit", params, socket) do
    request = %RegisterShopRequest{
      owner_id: socket.assigns.owner_id,
      name: params["name"],
      description: params["description"],
      # TODO: 画像アップロード実装
      logo_url: "",
      owner_name: params["owner_name"],
      phone_number: params["phone_number"],
      business_hours: params["business_hours"],
      return_policy: params["return_policy"],
      # カテゴリー機能は後で実装
      categories: []
    }

    case call_shop_service(request) do
      {:ok, _response} ->
        {:noreply,
         socket
         |> put_flash(:info, "ショップ登録が完了しました！管理者の承認をお待ちください。")
         |> push_navigate(to: "/owner/dashboard")}

      {:error, reason} ->
        {:noreply, assign(socket, :error, "登録失敗: #{reason}")}
    end
  end

  defp call_shop_service(request) do
    channel = get_shop_channel()

    case Stub.register_shop(channel, request) do
      {:ok, response} -> {:ok, response}
      {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
      error -> {:error, "接続エラー: #{inspect(error)}"}
    end
  end

  defp get_shop_channel do
    host = System.get_env("SHOP_SERVICE_HOST", "localhost")
    port = String.to_integer(System.get_env("SHOP_SERVICE_PORT", "22101"))

    {:ok, channel} = GRPC.Stub.connect("#{host}:#{port}")
    channel
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gradient-to-br from-purple-100 to-indigo-200 py-12">
      <div class="max-w-2xl mx-auto">
        <div class="bg-white rounded-lg shadow-lg p-8">
          <div class="text-center mb-8">
            <h1 class="text-3xl font-bold text-gray-800">ショップ登録</h1>
            <p class="text-gray-600 mt-2">ショップモールに出店するための情報を入力してください</p>
          </div>

          <%= if @error do %>
            <div class="mb-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
              {@error}
            </div>
          <% end %>

          <form phx-submit="submit" class="space-y-6">
            <!-- 基本情報 -->
            <div class="border-b border-gray-200 pb-6">
              <h2 class="text-lg font-semibold text-gray-800 mb-4">基本情報</h2>

              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    ショップ名 <span class="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    name="name"
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                    placeholder="例: おしゃれ雑貨店"
                    required
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    ショップ説明 <span class="text-red-500">*</span>
                  </label>
                  <textarea
                    name="description"
                    rows="4"
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                    placeholder="ショップの特徴や取り扱い商品について説明してください"
                    required
                  ></textarea>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    カテゴリー
                  </label>
                  <input
                    type="text"
                    name="categories"
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                    placeholder="例: アパレル, 雑貨, アクセサリー（カンマ区切り）"
                  />
                  <p class="text-xs text-gray-500 mt-1">複数ある場合はカンマで区切ってください</p>
                </div>
              </div>
            </div>
            
    <!-- 運営者情報 -->
            <div class="border-b border-gray-200 pb-6">
              <h2 class="text-lg font-semibold text-gray-800 mb-4">運営者情報</h2>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    運営者名 <span class="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    name="owner_name"
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                    placeholder="例: 山田 太郎"
                    required
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    電話番号 <span class="text-red-500">*</span>
                  </label>
                  <input
                    type="tel"
                    name="phone_number"
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                    placeholder="例: 03-1234-5678"
                    required
                  />
                </div>
              </div>
            </div>
            
    <!-- 営業情報 -->
            <div class="pb-6">
              <h2 class="text-lg font-semibold text-gray-800 mb-4">営業情報</h2>

              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    営業時間
                  </label>
                  <input
                    type="text"
                    name="business_hours"
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                    placeholder="例: 平日 10:00-18:00、土日祝 休み"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">
                    返品ポリシー
                  </label>
                  <textarea
                    name="return_policy"
                    rows="3"
                    class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                    placeholder="例: 商品到着後7日以内、未開封の場合のみ返品可能です。"
                  ></textarea>
                </div>
              </div>
            </div>
            
    <!-- 注意事項 -->
            <div class="bg-yellow-50 border border-yellow-200 rounded-md p-4 mb-6">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
                    <path
                      fill-rule="evenodd"
                      d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </div>
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-yellow-800">審査について</h3>
                  <p class="mt-1 text-sm text-yellow-700">
                    ショップ登録後、管理者による審査があります。審査完了後にショップが公開されます。
                  </p>
                </div>
              </div>
            </div>

            <button
              type="submit"
              class="w-full bg-purple-600 text-white py-3 px-4 rounded-md hover:bg-purple-700 transition-colors font-medium text-lg"
            >
              ショップを登録する
            </button>
          </form>
        </div>
      </div>
    </div>
    """
  end
end
