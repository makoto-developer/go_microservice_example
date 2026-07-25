defmodule ShopMallWebWeb.Admin.AdminPortalLive do
  @moduledoc """
  管理者ポータル。admin サービスの全 28 RPC を使用する:
  ダッシュボード(指標/売上チャート/サービスヘルス)、ユーザー管理、
  店舗管理(承認/却下/停止/再開)、カテゴリ管理、禁止語管理、
  監査ログ、レポート生成/エクスポート、システム設定。
  """
  use ShopMallWebWeb, :live_view

  alias AdminService.V1.AdminService.Stub
  alias AdminService.V1, as: PB

  @impl true
  def mount(_params, session, socket) do
    {:ok,
     socket
     |> assign(:admin_id, session["user_id"] || "admin")
     |> assign(:tab, "dashboard")
     |> assign(:notes, %{})}
  end

  @impl true
  def handle_event("set_tab", %{"tab" => tab}, socket) do
    {:noreply, assign(socket, :tab, tab)}
  end

  # ---- ダッシュボード ----

  @impl true
  def handle_event("load_dashboard", _params, socket) do
    metrics =
      simple(fn ch ->
        Stub.get_dashboard_metrics(ch, %PB.GetDashboardMetricsRequest{
          date: Date.to_iso8601(Date.utc_today())
        })
      end)

    chart =
      simple(fn ch -> Stub.get_sales_chart(ch, %PB.GetSalesChartRequest{group_by: "day"}) end)

    health = simple(fn ch -> Stub.get_service_health(ch, %PB.GetServiceHealthRequest{}) end)

    {:noreply, note(socket, "dashboard", "指標: #{metrics} / 売上: #{chart} / ヘルス: #{health}")}
  end

  # ---- ユーザー管理 ----

  @impl true
  def handle_event("load_users", _params, socket) do
    {:noreply,
     note(
       socket,
       "users",
       simple(fn ch -> Stub.get_all_users(ch, %PB.GetAllUsersRequest{page: 1, page_size: 20}) end)
     )}
  end

  @impl true
  def handle_event("user_detail", %{"user_id" => user_id}, socket) do
    {:noreply,
     note(
       socket,
       "users",
       simple(fn ch -> Stub.get_user_detail(ch, %PB.GetUserDetailRequest{user_id: user_id}) end)
     )}
  end

  @impl true
  def handle_event(
        "user_action",
        %{"action" => action, "user_id" => user_id, "reason" => reason},
        socket
      ) do
    admin_id = socket.assigns.admin_id

    result =
      case action do
        "suspend" ->
          simple(fn ch ->
            Stub.suspend_user(ch, %PB.SuspendUserRequest{
              user_id: user_id,
              admin_id: admin_id,
              reason: reason
            })
          end)

        "activate" ->
          simple(fn ch ->
            Stub.activate_user(ch, %PB.ActivateUserRequest{
              user_id: user_id,
              admin_id: admin_id,
              reason: reason
            })
          end)

        "make_owner" ->
          simple(fn ch ->
            Stub.change_user_role(ch, %PB.ChangeUserRoleRequest{
              user_id: user_id,
              admin_id: admin_id,
              new_role: "owner",
              reason: reason
            })
          end)

        _ ->
          "不明な操作"
      end

    {:noreply, note(socket, "users", result)}
  end

  # ---- 店舗管理 ----

  @impl true
  def handle_event("load_shops", _params, socket) do
    all =
      simple(fn ch -> Stub.get_all_shops(ch, %PB.GetAllShopsRequest{page: 1, page_size: 20}) end)

    pending =
      simple(fn ch ->
        Stub.get_pending_shops(ch, %PB.GetPendingShopsRequest{page: 1, page_size: 20})
      end)

    {:noreply, note(socket, "shops", "全店舗: #{all} / 承認待ち: #{pending}")}
  end

  @impl true
  def handle_event(
        "shop_action",
        %{"action" => action, "shop_id" => shop_id, "reason" => reason},
        socket
      ) do
    admin_id = socket.assigns.admin_id

    result =
      case action do
        "approve" ->
          simple(fn ch ->
            Stub.approve_shop(ch, %PB.ApproveShopRequest{shop_id: shop_id, admin_id: admin_id})
          end)

        "reject" ->
          simple(fn ch ->
            Stub.reject_shop(ch, %PB.RejectShopRequest{
              shop_id: shop_id,
              admin_id: admin_id,
              reason: reason
            })
          end)

        "suspend" ->
          simple(fn ch ->
            Stub.suspend_shop(ch, %PB.SuspendShopRequest{
              shop_id: shop_id,
              admin_id: admin_id,
              reason: reason
            })
          end)

        "activate" ->
          simple(fn ch ->
            Stub.activate_shop(ch, %PB.ActivateShopRequest{
              shop_id: shop_id,
              admin_id: admin_id,
              reason: reason
            })
          end)

        _ ->
          "不明な操作"
      end

    {:noreply, note(socket, "shops", result)}
  end

  # ---- カテゴリ管理 ----

  @impl true
  def handle_event("load_categories", _params, socket) do
    {:noreply,
     note(
       socket,
       "categories",
       simple(fn ch -> Stub.get_categories(ch, %PB.GetCategoriesRequest{}) end)
     )}
  end

  @impl true
  def handle_event("create_category", %{"name" => name}, socket) do
    {:noreply,
     note(
       socket,
       "categories",
       simple(fn ch ->
         Stub.create_category(ch, %PB.CreateCategoryRequest{
           admin_id: socket.assigns.admin_id,
           name: name
         })
       end)
     )}
  end

  @impl true
  def handle_event("update_category", %{"category_id" => id, "name" => name}, socket) do
    {:noreply,
     note(
       socket,
       "categories",
       simple(fn ch ->
         Stub.update_category(ch, %PB.UpdateCategoryRequest{
           admin_id: socket.assigns.admin_id,
           category_id: id,
           name: name
         })
       end)
     )}
  end

  @impl true
  def handle_event("delete_category", %{"category_id" => id}, socket) do
    {:noreply,
     note(
       socket,
       "categories",
       simple(fn ch ->
         Stub.delete_category(ch, %PB.DeleteCategoryRequest{
           admin_id: socket.assigns.admin_id,
           category_id: id
         })
       end)
     )}
  end

  # ---- 禁止語管理 ----

  @impl true
  def handle_event("load_words", _params, socket) do
    {:noreply,
     note(
       socket,
       "words",
       simple(fn ch ->
         Stub.get_forbidden_words(ch, %PB.GetForbiddenWordsRequest{context: "review"})
       end)
     )}
  end

  @impl true
  def handle_event("add_word", %{"word" => word}, socket) do
    {:noreply,
     note(
       socket,
       "words",
       simple(fn ch ->
         Stub.add_forbidden_word(ch, %PB.AddForbiddenWordRequest{
           admin_id: socket.assigns.admin_id,
           word: word,
           context: "review",
           severity: "high"
         })
       end)
     )}
  end

  @impl true
  def handle_event("delete_word", %{"word_id" => word_id}, socket) do
    {:noreply,
     note(
       socket,
       "words",
       simple(fn ch ->
         Stub.delete_forbidden_word(ch, %PB.DeleteForbiddenWordRequest{
           admin_id: socket.assigns.admin_id,
           word_id: word_id
         })
       end)
     )}
  end

  # ---- 監査ログ ----

  @impl true
  def handle_event("load_audit", _params, socket) do
    {:noreply,
     note(
       socket,
       "audit",
       simple(fn ch ->
         Stub.get_audit_logs(ch, %PB.GetAuditLogsRequest{page: 1, page_size: 20})
       end)
     )}
  end

  @impl true
  def handle_event("export_audit", %{"date_from" => from, "date_to" => to}, socket) do
    {:noreply,
     note(
       socket,
       "audit",
       simple(fn ch ->
         Stub.export_audit_logs(ch, %PB.ExportAuditLogsRequest{date_from: from, date_to: to})
       end)
     )}
  end

  # ---- レポート ----

  @impl true
  def handle_event("gen_sales_report", %{"date_from" => from, "date_to" => to}, socket) do
    {:noreply,
     note(
       socket,
       "reports",
       simple(fn ch ->
         Stub.generate_sales_report(ch, %PB.GenerateSalesReportRequest{
           date_from: from,
           date_to: to,
           report_type: "monthly"
         })
       end)
     )}
  end

  @impl true
  def handle_event("gen_user_report", %{"date_from" => from, "date_to" => to}, socket) do
    {:noreply,
     note(
       socket,
       "reports",
       simple(fn ch ->
         Stub.generate_user_report(ch, %PB.GenerateUserReportRequest{date_from: from, date_to: to})
       end)
     )}
  end

  @impl true
  def handle_event("export_report", %{"report_type" => type, "format" => format}, socket) do
    {:noreply,
     note(
       socket,
       "reports",
       simple(fn ch ->
         Stub.export_report(ch, %PB.ExportReportRequest{report_type: type, format: format})
       end)
     )}
  end

  # ---- システム設定 ----

  @impl true
  def handle_event("load_settings", _params, socket) do
    {:noreply,
     note(
       socket,
       "settings",
       simple(fn ch -> Stub.get_system_settings(ch, %PB.GetSystemSettingsRequest{}) end)
     )}
  end

  @impl true
  def handle_event("update_setting", %{"key" => key, "value" => value}, socket) do
    {:noreply,
     note(
       socket,
       "settings",
       simple(fn ch ->
         Stub.update_system_setting(ch, %PB.UpdateSystemSettingRequest{
           admin_id: socket.assigns.admin_id,
           setting_key: key,
           setting_value: value
         })
       end)
     )}
  end

  # ---- 共通 ----

  defp note(socket, key, message) do
    assign(socket, :notes, Map.put(socket.assigns.notes, key, message))
  end

  # RPC を呼び、結果メッセージ(または成功/失敗の文言)を返す
  defp simple(fun) do
    host = System.get_env("ADMIN_SERVICE_HOST", "localhost")
    port = System.get_env("ADMIN_SERVICE_PORT", "20111")

    case GRPC.Stub.connect("#{host}:#{port}") do
      {:ok, channel} ->
        try do
          case fun.(channel) do
            {:ok, response} -> Map.get(response, :message) || "OK"
            {:error, %GRPC.RPCError{message: message}} -> "エラー: #{message}"
            {:error, reason} -> "エラー: #{inspect(reason)}"
          end
        after
          GRPC.Stub.disconnect(channel)
        end

      {:error, reason} ->
        "接続エラー: #{inspect(reason)}"
    end
  end

  defp tab_class(current, tab) do
    if current == tab do
      "px-3 py-1.5 rounded-full text-sm font-medium bg-gray-900 text-white"
    else
      "px-3 py-1.5 rounded-full text-sm font-medium bg-white text-gray-600 border border-gray-300 hover:bg-gray-100"
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen bg-gray-50">
      <nav class="bg-gray-900 shadow-md">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex justify-between h-16 items-center">
            <span class="text-xl font-bold text-white">🛠 管理者コンソール</span>
            <div class="flex items-center space-x-4">
              <.link navigate="/admin/payments" class="text-gray-400 hover:text-white text-sm">
                決済
              </.link>
              <.link navigate="/admin/orders" class="text-gray-400 hover:text-white text-sm">
                注文
              </.link>
              <.link navigate="/admin/reviews" class="text-gray-400 hover:text-white text-sm">
                レビュー
              </.link>
              <.link navigate="/admin/notifications" class="text-gray-400 hover:text-white text-sm">
                通知
              </.link>
              <span class="text-gray-300 text-sm font-medium">ポータル</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-4xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-4">管理者ポータル</h1>

        <div class="flex flex-wrap gap-2 mb-6">
          <button
            :for={
              {tab, label} <- [
                {"dashboard", "ダッシュボード"},
                {"users", "ユーザー"},
                {"shops", "店舗"},
                {"categories", "カテゴリ"},
                {"words", "禁止語"},
                {"audit", "監査ログ"},
                {"reports", "レポート"},
                {"settings", "設定"}
              ]
            }
            phx-click="set_tab"
            phx-value-tab={tab}
            class={tab_class(@tab, tab)}
          >
            {label}
          </button>
        </div>

        <div
          :if={@notes[@tab]}
          class="bg-blue-50 border border-blue-200 text-blue-800 px-4 py-3 rounded mb-4 text-sm"
        >
          {@notes[@tab]}
        </div>

        <%= case @tab do %>
          <% "dashboard" -> %>
            <div class="bg-white shadow rounded-lg p-6">
              <h2 class="text-sm font-semibold text-gray-700 mb-3">今日の運営状況</h2>
              <button
                phx-click="load_dashboard"
                class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
              >
                指標・売上チャート・サービスヘルスを取得
              </button>
            </div>
          <% "users" -> %>
            <div class="bg-white shadow rounded-lg p-6 space-y-4">
              <button
                phx-click="load_users"
                class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
              >
                ユーザー一覧を取得
              </button>
              <form phx-submit="user_detail" class="flex space-x-2">
                <input
                  type="text"
                  name="user_id"
                  required
                  placeholder="ユーザーID"
                  class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
                >
                  詳細
                </button>
              </form>
              <form phx-submit="user_action" class="grid grid-cols-4 gap-2">
                <input
                  type="text"
                  name="user_id"
                  required
                  placeholder="ユーザーID"
                  class="border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <select name="action" class="border border-gray-300 rounded-md px-2 py-2 text-sm">
                  <option value="suspend">停止</option>
                  <option value="activate">再開</option>
                  <option value="make_owner">オーナー権限に変更</option>
                </select>
                <input
                  type="text"
                  name="reason"
                  required
                  placeholder="理由"
                  class="border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
                >
                  実行
                </button>
              </form>
            </div>
          <% "shops" -> %>
            <div class="bg-white shadow rounded-lg p-6 space-y-4">
              <button
                phx-click="load_shops"
                class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
              >
                店舗一覧・承認待ちを取得
              </button>
              <form phx-submit="shop_action" class="grid grid-cols-4 gap-2">
                <input
                  type="text"
                  name="shop_id"
                  required
                  placeholder="店舗ID"
                  class="border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <select name="action" class="border border-gray-300 rounded-md px-2 py-2 text-sm">
                  <option value="approve">承認</option>
                  <option value="reject">却下</option>
                  <option value="suspend">停止</option>
                  <option value="activate">再開</option>
                </select>
                <input
                  type="text"
                  name="reason"
                  placeholder="理由(却下/停止時)"
                  class="border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
                >
                  実行
                </button>
              </form>
            </div>
          <% "categories" -> %>
            <div class="bg-white shadow rounded-lg p-6 space-y-4">
              <button
                phx-click="load_categories"
                class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
              >
                カテゴリ一覧を取得
              </button>
              <form phx-submit="create_category" class="flex space-x-2">
                <input
                  type="text"
                  name="name"
                  required
                  placeholder="新しいカテゴリ名"
                  class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm text-white bg-gray-900 rounded-md hover:bg-gray-700"
                >
                  作成
                </button>
              </form>
              <form phx-submit="update_category" class="flex space-x-2">
                <input
                  type="text"
                  name="category_id"
                  required
                  placeholder="カテゴリID"
                  class="border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <input
                  type="text"
                  name="name"
                  required
                  placeholder="新しい名前"
                  class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
                >
                  更新
                </button>
              </form>
              <form phx-submit="delete_category" class="flex space-x-2">
                <input
                  type="text"
                  name="category_id"
                  required
                  placeholder="カテゴリID"
                  class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm text-red-600 border border-red-300 rounded-md hover:bg-red-50"
                >
                  削除
                </button>
              </form>
            </div>
          <% "words" -> %>
            <div class="bg-white shadow rounded-lg p-6 space-y-4">
              <button
                phx-click="load_words"
                class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
              >
                禁止語一覧を取得
              </button>
              <form phx-submit="add_word" class="flex space-x-2">
                <input
                  type="text"
                  name="word"
                  required
                  placeholder="追加する禁止語"
                  class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm text-white bg-gray-900 rounded-md hover:bg-gray-700"
                >
                  追加
                </button>
              </form>
              <form phx-submit="delete_word" class="flex space-x-2">
                <input
                  type="text"
                  name="word_id"
                  required
                  placeholder="禁止語ID"
                  class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm text-red-600 border border-red-300 rounded-md hover:bg-red-50"
                >
                  削除
                </button>
              </form>
            </div>
          <% "audit" -> %>
            <div class="bg-white shadow rounded-lg p-6 space-y-4">
              <button
                phx-click="load_audit"
                class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
              >
                監査ログを取得
              </button>
              <form phx-submit="export_audit" class="flex items-end space-x-2">
                <input
                  type="date"
                  name="date_from"
                  required
                  class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                />
                <span class="text-gray-400">〜</span>
                <input
                  type="date"
                  name="date_to"
                  required
                  class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
                >
                  ⬇ エクスポート
                </button>
              </form>
            </div>
          <% "reports" -> %>
            <div class="bg-white shadow rounded-lg p-6 space-y-4">
              <form phx-submit="gen_sales_report" class="flex items-end space-x-2">
                <input
                  type="date"
                  name="date_from"
                  required
                  class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                />
                <span class="text-gray-400">〜</span>
                <input
                  type="date"
                  name="date_to"
                  required
                  class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm text-white bg-gray-900 rounded-md hover:bg-gray-700"
                >
                  売上レポート生成
                </button>
              </form>
              <form phx-submit="gen_user_report" class="flex items-end space-x-2">
                <input
                  type="date"
                  name="date_from"
                  required
                  class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                />
                <span class="text-gray-400">〜</span>
                <input
                  type="date"
                  name="date_to"
                  required
                  class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
                >
                  ユーザーレポート生成
                </button>
              </form>
              <form phx-submit="export_report" class="flex items-end space-x-2">
                <select name="report_type" class="border border-gray-300 rounded-md px-2 py-2 text-sm">
                  <option value="sales">売上</option>
                  <option value="users">ユーザー</option>
                </select>
                <select name="format" class="border border-gray-300 rounded-md px-2 py-2 text-sm">
                  <option value="csv">CSV</option>
                  <option value="pdf">PDF</option>
                </select>
                <button
                  type="submit"
                  class="px-4 py-2 text-sm text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
                >
                  ⬇ エクスポート
                </button>
              </form>
            </div>
          <% "settings" -> %>
            <div class="bg-white shadow rounded-lg p-6 space-y-4">
              <button
                phx-click="load_settings"
                class="px-4 py-2 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
              >
                システム設定を取得
              </button>
              <form phx-submit="update_setting" class="flex space-x-2">
                <input
                  type="text"
                  name="key"
                  required
                  placeholder="設定キー(例: maintenance_mode)"
                  class="flex-1 border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <input
                  type="text"
                  name="value"
                  required
                  placeholder="値"
                  class="border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
                <button
                  type="submit"
                  class="px-4 py-2 text-sm text-white bg-gray-900 rounded-md hover:bg-gray-700"
                >
                  更新
                </button>
              </form>
            </div>
          <% _ -> %>
        <% end %>
      </main>
    </div>
    """
  end
end
