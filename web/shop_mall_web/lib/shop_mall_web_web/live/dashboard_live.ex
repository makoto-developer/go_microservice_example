defmodule ShopMallWebWeb.DashboardLive do
  use ShopMallWebWeb, :live_view

  @impl true
  def mount(_params, session, socket) do
    # セッションからユーザーIDを取得（後で認証実装）
    user_id = session["user_id"] || "guest"

    {:ok,
     socket
     |> assign(:user_id, user_id)
     |> assign(:page_title, "ダッシュボード")}
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
              <h1 class="text-2xl font-bold text-gray-900">
                オンラインショップモール
              </h1>
            </div>
            <div class="flex items-center space-x-4">
              <.link
                navigate="/products"
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                商品一覧
              </.link>
              <.link
                navigate="/auth"
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
        <div class="px-4 py-6 sm:px-0">
          <!-- ウェルカムメッセージ -->
          <div class="bg-white overflow-hidden shadow rounded-lg mb-6">
            <div class="px-4 py-5 sm:p-6">
              <h2 class="text-3xl font-bold text-gray-900 mb-2">
                ようこそ！
              </h2>
              <p class="text-gray-600">
                オンラインショップモールへようこそ。お買い物をお楽しみください。
              </p>
            </div>
          </div>

          <!-- カード一覧 -->
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <!-- 商品一覧カード -->
            <div class="bg-white overflow-hidden shadow rounded-lg hover:shadow-lg transition-shadow">
              <div class="px-4 py-5 sm:p-6">
                <div class="flex items-center">
                  <div class="flex-shrink-0 bg-blue-500 rounded-md p-3">
                    <svg
                      class="h-6 w-6 text-white"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
                      />
                    </svg>
                  </div>
                  <div class="ml-5 w-0 flex-1">
                    <dl>
                      <dt class="text-sm font-medium text-gray-500 truncate">
                        商品を探す
                      </dt>
                      <dd class="text-lg font-semibold text-gray-900">
                        商品一覧を見る
                      </dd>
                    </dl>
                  </div>
                </div>
                <div class="mt-5">
                  <.link
                    navigate="/products"
                    class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
                  >
                    商品を見る →
                  </.link>
                </div>
              </div>
            </div>

            <!-- 注文履歴カード（未実装） -->
            <div class="bg-white overflow-hidden shadow rounded-lg opacity-60">
              <div class="px-4 py-5 sm:p-6">
                <div class="flex items-center">
                  <div class="flex-shrink-0 bg-green-500 rounded-md p-3">
                    <svg
                      class="h-6 w-6 text-white"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                      />
                    </svg>
                  </div>
                  <div class="ml-5 w-0 flex-1">
                    <dl>
                      <dt class="text-sm font-medium text-gray-500 truncate">
                        注文履歴
                      </dt>
                      <dd class="text-lg font-semibold text-gray-900">
                        準備中
                      </dd>
                    </dl>
                  </div>
                </div>
                <div class="mt-5">
                  <button
                    disabled
                    class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-400 cursor-not-allowed"
                  >
                    近日公開
                  </button>
                </div>
              </div>
            </div>

            <!-- カート（未実装） -->
            <div class="bg-white overflow-hidden shadow rounded-lg opacity-60">
              <div class="px-4 py-5 sm:p-6">
                <div class="flex items-center">
                  <div class="flex-shrink-0 bg-purple-500 rounded-md p-3">
                    <svg
                      class="h-6 w-6 text-white"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
                      />
                    </svg>
                  </div>
                  <div class="ml-5 w-0 flex-1">
                    <dl>
                      <dt class="text-sm font-medium text-gray-500 truncate">
                        カート
                      </dt>
                      <dd class="text-lg font-semibold text-gray-900">
                        準備中
                      </dd>
                    </dl>
                  </div>
                </div>
                <div class="mt-5">
                  <button
                    disabled
                    class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-gray-400 cursor-not-allowed"
                  >
                    近日公開
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
    """
  end
end
