defmodule ShopMallWebWeb.Admin.NotificationCenterLive do
  @moduledoc """
  管理者の通知センター。
  メールテンプレートの作成/更新/プレビュー、一斉メール(SendBulkEmail)、
  プッシュ通知(SendPushNotification)、配信履歴(GetNotificationHistory)、
  再送(ResendNotification)を行う。
  """
  use ShopMallWebWeb, :live_view

  alias NotificationService.V1.NotificationService.Stub
  alias NotificationService.V1, as: PB

  @impl true
  def mount(_params, _session, socket) do
    {:ok,
     socket
     |> assign(:preview_note, nil)
     |> assign(:history_note, nil)}
  end

  @impl true
  def handle_event("create_template", params, socket) do
    request = %PB.CreateEmailTemplateRequest{
      template_key: params["template_key"],
      subject_template: params["subject"],
      html_template: params["body"],
      text_template: params["body"]
    }

    run(socket, call(fn ch -> Stub.create_email_template(ch, request) end), "テンプレートを作成しました")
  end

  @impl true
  def handle_event("update_template", params, socket) do
    request = %PB.UpdateEmailTemplateRequest{
      template_id: params["template_id"],
      subject_template: params["subject"],
      html_template: params["body"],
      text_template: params["body"]
    }

    run(socket, call(fn ch -> Stub.update_email_template(ch, request) end), "テンプレートを更新しました")
  end

  @impl true
  def handle_event("preview_template", %{"template_key" => key}, socket) do
    case call(fn ch ->
           Stub.preview_email_template(ch, %PB.PreviewEmailTemplateRequest{template_key: key})
         end) do
      {:ok, resp} -> {:noreply, assign(socket, :preview_note, resp.message || "プレビューを生成しました")}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "プレビューに失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("send_bulk", %{"user_ids" => user_ids, "template_key" => key}, socket) do
    request = %PB.SendBulkEmailRequest{
      user_ids: String.split(user_ids, ",", trim: true) |> Enum.map(&String.trim/1),
      template_key: key
    }

    run(socket, call(fn ch -> Stub.send_bulk_email(ch, request) end), "一斉メールを送信しました")
  end

  @impl true
  def handle_event("send_push", %{"user_id" => user_id, "template_key" => key}, socket) do
    request = %PB.SendPushNotificationRequest{
      user_id: user_id,
      template_key: key,
      notification_type: :ORDER_CONFIRMED
    }

    run(socket, call(fn ch -> Stub.send_push_notification(ch, request) end), "プッシュ通知を送信しました")
  end

  @impl true
  def handle_event("load_history", %{"user_id" => user_id}, socket) do
    case call(fn ch ->
           Stub.get_notification_history(ch, %PB.GetNotificationHistoryRequest{
             user_id: user_id,
             page: 1,
             page_size: 20
           })
         end) do
      {:ok, resp} -> {:noreply, assign(socket, :history_note, resp.message || "履歴を取得しました")}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "履歴の取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("resend", %{"notification_id" => id}, socket) do
    run(
      socket,
      call(fn ch ->
        Stub.resend_notification(ch, %PB.ResendNotificationRequest{notification_id: id})
      end),
      "再送しました"
    )
  end

  defp run(socket, result, ok_msg) do
    case result do
      {:ok, resp} -> {:noreply, put_flash(socket, :info, Map.get(resp, :message) || ok_msg)}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "操作に失敗しました: #{reason}")}
    end
  end

  defp call(fun) do
    host = System.get_env("NOTIFICATION_SERVICE_HOST", "localhost")
    port = System.get_env("NOTIFICATION_SERVICE_PORT", "20107")

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
        {:error, "通知サービスに接続できません: #{inspect(reason)}"}
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
              <.link navigate="/admin/orders" class="text-gray-400 hover:text-white text-sm">
                注文分析
              </.link>
              <.link navigate="/admin/search" class="text-gray-400 hover:text-white text-sm">
                検索管理
              </.link>
              <span class="text-gray-300 text-sm font-medium">通知センター</span>
            </div>
          </div>
        </div>
      </nav>

      <main class="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-6">通知センター</h1>

        <div class="bg-white shadow rounded-lg p-6 mb-6">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">メールテンプレート</h2>
          <form phx-submit="create_template" class="grid grid-cols-2 gap-2 mb-3">
            <input
              type="text"
              name="template_key"
              required
              placeholder="キー(例: campaign_summer)"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <input
              type="text"
              name="subject"
              required
              placeholder="件名テンプレート"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <textarea
              name="body"
              rows="2"
              required
              placeholder="本文テンプレート({{name}} など)"
              class="col-span-2 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            ></textarea>
            <div class="col-span-2 flex justify-end">
              <button
                type="submit"
                class="px-3 py-1.5 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
              >
                作成
              </button>
            </div>
          </form>
          <form phx-submit="update_template" class="grid grid-cols-2 gap-2 mb-3">
            <input
              type="text"
              name="template_id"
              required
              placeholder="テンプレートID"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <input
              type="text"
              name="subject"
              required
              placeholder="新しい件名"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <textarea
              name="body"
              rows="2"
              required
              placeholder="新しい本文"
              class="col-span-2 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            ></textarea>
            <div class="col-span-2 flex justify-end">
              <button
                type="submit"
                class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
              >
                更新
              </button>
            </div>
          </form>
          <form phx-submit="preview_template" class="flex items-end space-x-2">
            <input
              type="text"
              name="template_key"
              required
              placeholder="キーを指定してプレビュー"
              class="flex-1 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <button
              type="submit"
              class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              プレビュー
            </button>
          </form>
          <div :if={@preview_note} class="text-xs text-gray-600 mt-2">{@preview_note}</div>
        </div>

        <div class="bg-white shadow rounded-lg p-6 mb-6">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">配信</h2>
          <form phx-submit="send_bulk" class="grid grid-cols-3 gap-2 mb-3">
            <input
              type="text"
              name="user_ids"
              required
              placeholder="ユーザーID(カンマ区切り)"
              class="col-span-2 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <input
              type="text"
              name="template_key"
              required
              placeholder="テンプレートキー"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <div class="col-span-3 flex justify-end">
              <button
                type="submit"
                class="px-3 py-1.5 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
              >
                📧 一斉メール送信
              </button>
            </div>
          </form>
          <form phx-submit="send_push" class="grid grid-cols-3 gap-2">
            <input
              type="text"
              name="user_id"
              required
              placeholder="ユーザーID"
              class="col-span-2 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <input
              type="text"
              name="template_key"
              required
              placeholder="テンプレートキー"
              class="border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <div class="col-span-3 flex justify-end">
              <button
                type="submit"
                class="px-3 py-1.5 text-sm font-medium text-white bg-gray-900 rounded-md hover:bg-gray-700"
              >
                🔔 プッシュ通知送信
              </button>
            </div>
          </form>
        </div>

        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-sm font-semibold text-gray-700 mb-3">配信履歴・再送</h2>
          <form phx-submit="load_history" class="flex items-end space-x-2 mb-3">
            <input
              type="text"
              name="user_id"
              required
              placeholder="ユーザーID"
              class="flex-1 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <button
              type="submit"
              class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              履歴を表示
            </button>
          </form>
          <div :if={@history_note} class="text-xs text-gray-600 mb-3">{@history_note}</div>
          <form phx-submit="resend" class="flex items-end space-x-2">
            <input
              type="text"
              name="notification_id"
              required
              placeholder="通知IDを指定して再送"
              class="flex-1 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
            />
            <button
              type="submit"
              class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
            >
              再送
            </button>
          </form>
        </div>
      </main>
    </div>
    """
  end
end
