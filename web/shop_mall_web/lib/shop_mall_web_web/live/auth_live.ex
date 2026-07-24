defmodule ShopMallWebWeb.AuthLive do
  use ShopMallWebWeb, :live_view

  alias AuthService.V1.{
    AuthService.Stub,
    UserRegistrationRequest,
    UserLoginRequest,
    Role
  }

  @impl true
  def mount(_params, _session, socket) do
    {:ok,
     socket
     |> assign(:mode, :login)
     |> assign(:error, nil)
     |> assign(:success, nil)}
  end

  @impl true
  def handle_event("toggle_mode", _params, socket) do
    mode = if socket.assigns.mode == :login, do: :register, else: :login
    {:noreply, assign(socket, mode: mode, error: nil, success: nil)}
  end

  @impl true
  def handle_event("submit", %{"email" => email, "password" => password}, socket) do
    IO.puts("=== AUTH EVENT: submit received ===")
    IO.inspect(%{email: email, password: String.length(password), mode: socket.assigns.mode})

    case socket.assigns.mode do
      :login -> handle_login(socket, email, password)
      :register -> handle_register(socket, email, password)
    end
  end

  defp handle_login(socket, email, password) do
    IO.puts("=== HANDLE LOGIN called ===")
    IO.inspect(%{email: email})

    request = %UserLoginRequest{
      email: email,
      password: password
    }

    IO.puts("=== Calling auth service for login ===")

    case call_auth_service(:login, request) do
      {:ok, response} ->
        IO.puts("=== Login SUCCESS ===")
        IO.inspect(response)

        {:noreply,
         socket
         |> put_flash(:info, "ログイン成功！")
         |> push_navigate(to: "/dashboard")}

      {:error, reason} ->
        IO.puts("=== Login FAILED ===")
        IO.inspect(reason)

        {:noreply,
         socket
         |> assign(:error, "ログイン失敗: #{inspect(reason)}")
         |> assign(:success, nil)}
    end
  end

  defp handle_register(socket, email, password) do
    request = %UserRegistrationRequest{
      email: email,
      password: password,
      role: Role.value(:CUSTOMER)
    }

    case call_auth_service(:register, request) do
      {:ok, _response} ->
        {:noreply,
         socket
         |> put_flash(:info, "登録成功！メールを確認してください。ログイン画面に移動します...")
         |> assign(:mode, :login)
         |> assign(:error, nil)}

      {:error, reason} ->
        {:noreply,
         socket
         |> assign(:error, "登録失敗: #{inspect(reason)}")
         |> assign(:success, nil)}
    end
  end

  defp call_auth_service(:login, request) do
    case get_auth_channel() do
      {:ok, channel} ->
        case Stub.login(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, "接続エラー: #{inspect(error)}"}
        end

      {:error, reason} ->
        {:error, "Auth Serviceに接続できません: #{inspect(reason)}"}
    end
  end

  defp call_auth_service(:register, request) do
    case get_auth_channel() do
      {:ok, channel} ->
        case Stub.register(channel, request) do
          {:ok, response} -> {:ok, response}
          {:error, %GRPC.RPCError{} = error} -> {:error, error.message}
          error -> {:error, "接続エラー: #{inspect(error)}"}
        end

      {:error, reason} ->
        {:error, "Auth Serviceに接続できません: #{inspect(reason)}"}
    end
  end

  defp get_auth_channel do
    # Auth Service の接続先（Docker ネットワーク内 or localhost）
    host = System.get_env("AUTH_SERVICE_HOST", "localhost")
    port = String.to_integer(System.get_env("AUTH_SERVICE_PORT", "22100"))

    # No TLS for development (default is insecure)
    case GRPC.Stub.connect("#{host}:#{port}") do
      {:ok, channel} -> {:ok, channel}
      {:error, reason} -> {:error, reason}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen flex items-center justify-center bg-gray-100">
      <div class="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
        <h2 class="text-3xl font-bold text-center text-gray-800 mb-8">
          {if @mode == :login, do: "ログイン", else: "新規登録"}
        </h2>

        <%= if @success do %>
          <div class="mb-4 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
            {@success}
          </div>
        <% end %>

        <%= if @error do %>
          <div class="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
            {@error}
          </div>
        <% end %>

        <form phx-submit="submit" class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              メールアドレス
            </label>
            <input
              type="email"
              name="email"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="email@example.com"
              required
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              パスワード
            </label>
            <input
              type="password"
              name="password"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="••••••••"
              required
            />
          </div>

          <button
            type="submit"
            class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition-colors font-medium"
          >
            {if @mode == :login, do: "ログイン", else: "登録"}
          </button>
        </form>

        <%= if @mode == :login do %>
          <div class="mt-4 text-center">
            <.link
              href="/auth/password-reset"
              class="text-sm text-blue-600 hover:text-blue-800"
            >
              パスワードを忘れた方はこちら
            </.link>
          </div>
        <% end %>

        <div class="mt-6 text-center">
          <button
            phx-click="toggle_mode"
            class="text-blue-600 hover:text-blue-800 text-sm"
          >
            <%= if @mode == :login do %>
              アカウントをお持ちでない方はこちら
            <% else %>
              すでにアカウントをお持ちの方はこちら
            <% end %>
          </button>
        </div>
      </div>
    </div>
    """
  end
end
