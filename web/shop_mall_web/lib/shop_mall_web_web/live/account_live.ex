defmodule ShopMallWebWeb.AccountLive do
  @moduledoc """
  アカウント設定画面。
  パスワード変更(ChangePassword)・セッション診断(VerifyToken)・
  セッション延長(RefreshToken)を行う。
  """
  use ShopMallWebWeb, :live_view

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
        <% end %>
      </main>
    </div>
    """
  end
end
