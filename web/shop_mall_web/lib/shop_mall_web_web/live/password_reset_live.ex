defmodule ShopMallWebWeb.PasswordResetLive do
  use ShopMallWebWeb, :live_view
  alias AuthService.V1.{AuthService.Stub, PasswordResetRequestRequest}

  @impl true
  def mount(_params, _session, socket) do
    {:ok,
     socket
     |> assign(:email, "")
     |> assign(:message, nil)
     |> assign(:error, nil)
     |> assign(:loading, false)}
  end

  @impl true
  def handle_event("update_email", %{"email" => email}, socket) do
    {:noreply, assign(socket, :email, email)}
  end

  @impl true
  def handle_event("request_reset", %{"email" => email}, socket) do
    socket = assign(socket, :loading, true)

    case request_password_reset(email) do
      {:ok, _response} ->
        {:noreply,
         socket
         |> assign(:loading, false)
         |> assign(:message, "パスワードリセットのメールを送信しました。メールをご確認ください。")
         |> assign(:error, nil)
         |> assign(:email, "")}

      {:error, reason} ->
        {:noreply,
         socket
         |> assign(:loading, false)
         |> assign(:error, "エラーが発生しました: #{reason}")
         |> assign(:message, nil)}
    end
  end

  defp request_password_reset(email) do
    channel = get_auth_channel()
    request = %PasswordResetRequestRequest{email: email}

    case Stub.request_password_reset(channel, request) do
      {:ok, response} -> {:ok, response}
      {:error, %GRPC.RPCError{message: message}} -> {:error, message}
      {:error, reason} -> {:error, inspect(reason)}
    end
  end

  defp get_auth_channel do
    host = System.get_env("AUTH_SERVICE_HOST", "localhost")
    port = String.to_integer(System.get_env("AUTH_SERVICE_PORT", "22000"))

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
            パスワードをリセット
          </h2>
          <p class="mt-2 text-center text-sm text-gray-600">
            登録したメールアドレスを入力してください
          </p>
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

        <form phx-submit="request_reset" class="mt-8 space-y-6">
          <div>
            <label for="email" class="sr-only">メールアドレス</label>
            <input
              id="email"
              name="email"
              type="email"
              required
              value={@email}
              phx-change="update_email"
              class="appearance-none rounded-md relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="メールアドレス"
              disabled={@loading}
            />
          </div>

          <div>
            <button
              type="submit"
              disabled={@loading}
              class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
            >
              <%= if @loading do %>
                送信中...
              <% else %>
                リセットメールを送信
              <% end %>
            </button>
          </div>

          <div class="text-center">
            <.link
              href="/auth"
              class="font-medium text-indigo-600 hover:text-indigo-500"
            >
              ログイン画面に戻る
            </.link>
          </div>
        </form>
      </div>
    </div>
    """
  end
end
