defmodule ShopMallWebWeb.AccountLive do
  @moduledoc """
  アカウント設定画面。
  パスワード変更(ChangePassword)・セッション診断(VerifyToken)・
  セッション延長(RefreshToken)を行う。
  """
  use ShopMallWebWeb, :live_view

  alias NotificationService.V1.NotificationService.Stub, as: NotificationStub
  alias NotificationService.V1, as: NPB

  alias AuthService.V1.{
    AuthService.Stub,
    PasswordChangeRequest,
    TokenRefreshRequest,
    TokenVerificationRequest
  }

  @impl true
  def mount(_params, session, socket) do
    {:ok,
     socket
     |> assign(:user_id, session["user_id"])
     |> assign(:access_token, session["access_token"])
     |> assign(:refresh_token, session["refresh_token"])
     |> assign(:token_info, nil)}
  end

  @impl true
  def handle_event("change_password", %{"current" => current, "new" => new_pw}, socket) do
    request = %PasswordChangeRequest{
      user_id: socket.assigns.user_id || "",
      current_password: current,
      new_password: new_pw
    }

    with {:ok, channel} <- auth_channel(),
         {:ok, response} <- Stub.change_password(channel, request) do
      if response.success do
        {:noreply, put_flash(socket, :info, "パスワードを変更しました")}
      else
        {:noreply, put_flash(socket, :error, "変更できませんでした: #{response.message}")}
      end
    else
      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "パスワード変更に失敗しました: #{inspect(reason)}")}
    end
  end

  @impl true
  def handle_event("verify_token", _params, socket) do
    request = %TokenVerificationRequest{access_token: socket.assigns.access_token || ""}

    with {:ok, channel} <- auth_channel(),
         {:ok, response} <- Stub.verify_token(channel, request) do
      {:noreply, assign(socket, :token_info, response)}
    else
      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "トークン検証に失敗しました: #{inspect(reason)}")}
    end
  end

  @impl true
  def handle_event("refresh_session", _params, socket) do
    request = %TokenRefreshRequest{refresh_token: socket.assigns.refresh_token || ""}

    with {:ok, channel} <- auth_channel(),
         {:ok, response} <- Stub.refresh_token(channel, request) do
      {:noreply,
       socket
       |> assign(:access_token, response.access_token)
       |> assign(:refresh_token, response.refresh_token)
       |> put_flash(:info, "セッションを延長しました(新しいトークンを発行)")}
    else
      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "セッション延長に失敗しました: #{inspect(reason)}")}
    end
  end

  @impl true
  def handle_event("load_preference", _params, socket) do
    case notification_call(fn ch ->
           NotificationStub.get_notification_preference(ch, %NPB.GetNotificationPreferenceRequest{
             user_id: socket.assigns.user_id || ""
           })
         end) do
      {:ok, resp} ->
        {:noreply, put_flash(socket, :info, "通知設定: #{Map.get(resp, :message) || "取得しました"}")}

      {:error, reason} ->
        {:noreply, put_flash(socket, :error, "通知設定の取得に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("save_preference", params, socket) do
    request = %NPB.UpdateNotificationPreferenceRequest{
      user_id: socket.assigns.user_id || "",
      email_enabled: params["email_enabled"] == "on",
      push_enabled: params["push_enabled"] == "on",
      email_order_updates: params["email_enabled"] == "on",
      push_order_updates: params["push_enabled"] == "on"
    }

    case notification_call(fn ch ->
           NotificationStub.update_notification_preference(ch, request)
         end) do
      {:ok, _} -> {:noreply, put_flash(socket, :info, "通知設定を保存しました")}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "保存に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("register_device", %{"token" => token}, socket) do
    request = %NPB.RegisterDeviceTokenRequest{
      user_id: socket.assigns.user_id || "",
      device_id: "web-#{socket.assigns.user_id}",
      platform: "web",
      token: token
    }

    case notification_call(fn ch -> NotificationStub.register_device_token(ch, request) end) do
      {:ok, _} -> {:noreply, put_flash(socket, :info, "デバイスを登録しました")}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "登録に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("refresh_device", %{"token" => token}, socket) do
    request = %NPB.RefreshDeviceTokenRequest{
      device_id: "web-#{socket.assigns.user_id}",
      new_token: token
    }

    case notification_call(fn ch -> NotificationStub.refresh_device_token(ch, request) end) do
      {:ok, _} -> {:noreply, put_flash(socket, :info, "デバイストークンを更新しました")}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "更新に失敗しました: #{reason}")}
    end
  end

  @impl true
  def handle_event("unregister_device", _params, socket) do
    request = %NPB.UnregisterDeviceTokenRequest{device_id: "web-#{socket.assigns.user_id}"}

    case notification_call(fn ch -> NotificationStub.unregister_device_token(ch, request) end) do
      {:ok, _} -> {:noreply, put_flash(socket, :info, "デバイスを解除しました")}
      {:error, reason} -> {:noreply, put_flash(socket, :error, "解除に失敗しました: #{reason}")}
    end
  end

  defp notification_call(fun) do
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

  defp auth_channel do
    host = System.get_env("AUTH_SERVICE_HOST", "localhost")
    port = System.get_env("AUTH_SERVICE_PORT", "22100")
    GRPC.Stub.connect("#{host}:#{port}")
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
                navigate="/orders"
                class="text-gray-700 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
              >
                注文履歴
              </.link>
              <span class="text-gray-900 px-3 py-2 text-sm font-semibold">アカウント</span>
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

      <main class="max-w-2xl mx-auto py-6 sm:px-6 lg:px-8">
        <h1 class="text-2xl font-bold text-gray-900 mb-6">アカウント設定</h1>

        <%= if is_nil(@user_id) do %>
          <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-6">
            アカウント設定にはログインが必要です
          </div>
        <% else %>
          <div class="bg-white shadow rounded-lg p-6 mb-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">パスワード変更</h2>
            <form phx-submit="change_password" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">現在のパスワード</label>
                <input
                  type="password"
                  name="current"
                  required
                  class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">新しいパスワード</label>
                <input
                  type="password"
                  name="new"
                  required
                  minlength="8"
                  class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm"
                />
              </div>
              <div class="flex justify-end">
                <button
                  type="submit"
                  class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                >
                  変更する
                </button>
              </div>
            </form>
          </div>

          <div class="bg-white shadow rounded-lg p-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">セッション</h2>
            <div class="flex space-x-3 mb-4">
              <button
                phx-click="verify_token"
                class="px-4 py-2 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
              >
                ログイン状態を確認
              </button>
              <button
                phx-click="refresh_session"
                class="px-4 py-2 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
              >
                セッションを延長
              </button>
            </div>
            <%= if @token_info do %>
              <div class={"px-3 py-2 rounded text-sm " <>
                if(@token_info.valid,
                  do: "bg-green-50 border border-green-300 text-green-700",
                  else: "bg-red-50 border border-red-300 text-red-700"
                )}>
                <%= if @token_info.valid do %>
                  ✓ 有効なセッションです(ユーザー: {@token_info.user_id}、ロール: {@token_info.role})
                <% else %>
                  ✕ セッションが無効です。再ログインしてください
                <% end %>
              </div>
            <% end %>
          </div>

          <div class="bg-white shadow rounded-lg p-6 mt-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">通知設定</h2>
            <div class="flex items-center justify-between mb-3">
              <button phx-click="load_preference" class="text-sm text-blue-600 hover:text-blue-800">
                現在の設定を確認
              </button>
            </div>
            <form phx-submit="save_preference" class="space-y-3 mb-6">
              <label class="flex items-center space-x-2 text-sm text-gray-700">
                <input type="checkbox" name="email_enabled" checked /> <span>メール通知を受け取る</span>
              </label>
              <label class="flex items-center space-x-2 text-sm text-gray-700">
                <input type="checkbox" name="push_enabled" /> <span>プッシュ通知を受け取る</span>
              </label>
              <div class="flex justify-end">
                <button
                  type="submit"
                  class="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                >
                  保存
                </button>
              </div>
            </form>

            <h3 class="text-sm font-semibold text-gray-700 mb-2">プッシュ通知デバイス</h3>
            <form phx-submit="register_device" class="flex items-end space-x-2 mb-2">
              <input
                type="text"
                name="token"
                required
                placeholder="デバイストークン"
                class="flex-1 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
              />
              <button
                type="submit"
                class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
              >
                登録
              </button>
            </form>
            <form phx-submit="refresh_device" class="flex items-end space-x-2 mb-2">
              <input
                type="text"
                name="token"
                required
                placeholder="新しいトークン"
                class="flex-1 border border-gray-300 rounded-md px-2 py-1.5 text-sm"
              />
              <button
                type="submit"
                class="px-3 py-1.5 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50"
              >
                トークン更新
              </button>
            </form>
            <div class="flex justify-end">
              <button phx-click="unregister_device" class="text-sm text-red-600 hover:text-red-800">
                このデバイスを解除
              </button>
            </div>
          </div>
        <% end %>
      </main>
    </div>
    """
  end
end
