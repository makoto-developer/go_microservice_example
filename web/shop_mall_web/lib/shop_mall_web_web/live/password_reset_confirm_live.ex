defmodule ShopMallWebWeb.PasswordResetConfirmLive do
  use ShopMallWebWeb, :live_view
  alias AuthService.V1.{AuthService.Stub, PasswordResetRequest}

  @impl true
  def mount(%{"token" => token}, _session, socket) do
    {:ok,
     socket
     |> assign(:token, token)
     |> assign(:new_password, "")
     |> assign(:confirm_password, "")
     |> assign(:message, nil)
     |> assign(:error, nil)
     |> assign(:loading, false)
     |> assign(:success, false)}
  end

  def mount(_params, _session, socket) do
    {:ok,
     socket
     |> assign(:error, "無効なリンクです")
     |> assign(:token, nil)}
  end

  @impl true
  def handle_event("update_new_password", %{"password" => password}, socket) do
    {:noreply, assign(socket, :new_password, password)}
  end

  @impl true
  def handle_event("update_confirm_password", %{"password" => password}, socket) do
    {:noreply, assign(socket, :confirm_password, password)}
  end

  @impl true
  def handle_event(
        "reset_password",
        %{"new_password" => new_password, "confirm_password" => confirm_password},
        socket
      ) do
    cond do
      new_password != confirm_password ->
        {:noreply,
         socket
         |> assign(:error, "パスワードが一致しません")
         |> assign(:message, nil)}

      String.length(new_password) < 8 ->
        {:noreply,
         socket
         |> assign(:error, "パスワードは8文字以上である必要があります")
         |> assign(:message, nil)}

      true ->
        socket = assign(socket, :loading, true)

        case reset_password(socket.assigns.token, new_password) do
          {:ok, _response} ->
            {:noreply,
             socket
             |> assign(:loading, false)
             |> assign(:success, true)
             |> assign(:message, "パスワードがリセットされました。ログイン画面からログインしてください。")
             |> assign(:error, nil)}

          {:error, reason} ->
            {:noreply,
             socket
             |> assign(:loading, false)
             |> assign(:error, "パスワードリセットに失敗しました: #{reason}")
             |> assign(:message, nil)}
        end
    end
  end

  defp reset_password(token, new_password) do
    channel = get_auth_channel()
    request = %PasswordResetRequest{token: token, new_password: new_password}

    case Stub.reset_password(channel, request) do
      {:ok, response} -> {:ok, response}
      {:error, %GRPC.RPCError{message: message}} -> {:error, message}
      {:error, reason} -> {:error, inspect(reason)}
    end
  end

  defp get_auth_channel do
    host = System.get_env("AUTH_SERVICE_HOST", "localhost")
    port = String.to_integer(System.get_env("AUTH_SERVICE_PORT", "22100"))

    # No TLS for development (default is insecure)
    {:ok, channel} = GRPC.Stub.connect("#{host}:#{port}")
    channel
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full space-y-8">
        <div>
          <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
            新しいパスワードを設定
          </h2>
        </div>

        <%= if @message do %>
          <div class="rounded-md bg-green-50 p-4">
            <div class="flex">
              <div class="ml-3">
                <p class="text-sm font-medium text-green-800">
                  <%= @message %>
                </p>
              </div>
            </div>
          </div>
          <div class="text-center mt-4">
            <.link
              href="/auth"
              class="font-medium text-indigo-600 hover:text-indigo-500"
            >
              ログイン画面へ
            </.link>
          </div>
        <% end %>

        <%= if @error do %>
          <div class="rounded-md bg-red-50 p-4">
            <div class="flex">
              <div class="ml-3">
                <p class="text-sm font-medium text-red-800">
                  <%= @error %>
                </p>
              </div>
            </div>
          </div>
        <% end %>

        <%= if @token && !@success do %>
          <form phx-submit="reset_password" class="mt-8 space-y-6">
            <div class="space-y-4">
              <div>
                <label for="new_password" class="sr-only">新しいパスワード</label>
                <input
                  id="new_password"
                  name="new_password"
                  type="password"
                  required
                  value={@new_password}
                  phx-change="update_new_password"
                  class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                  placeholder="新しいパスワード（8文字以上）"
                  disabled={@loading}
                />
              </div>

              <div>
                <label for="confirm_password" class="sr-only">パスワード確認</label>
                <input
                  id="confirm_password"
                  name="confirm_password"
                  type="password"
                  required
                  value={@confirm_password}
                  phx-change="update_confirm_password"
                  class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
                  placeholder="パスワード確認"
                  disabled={@loading}
                />
              </div>
            </div>

            <div>
              <button
                type="submit"
                disabled={@loading}
                class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
              >
                <%= if @loading do %>
                  処理中...
                <% else %>
                  パスワードをリセット
                <% end %>
              </button>
            </div>
          </form>
        <% end %>
      </div>
    </div>
    """
  end
end
