defmodule ShopMallWebWeb.Owner.AuthLive do
  use ShopMallWebWeb, :live_view

  alias AuthService.V1.{
    AuthService.Stub,
    UserRegistrationRequest,
    UserLoginRequest,
    Role
  }

  @impl true
  def mount(_params, session, socket) do
    {:ok,
     socket
     |> assign(:mode, :login)
     |> assign(:error, nil)
     |> assign(:success, nil)
     |> assign(:current_user, session["current_user"])}
  end

  @impl true
  def handle_event("toggle_mode", _params, socket) do
    mode = if socket.assigns.mode == :login, do: :register, else: :login
    {:noreply, assign(socket, mode: mode, error: nil, success: nil)}
  end

  @impl true
  def handle_event("submit", %{"email" => email, "password" => password}, socket) do
    case socket.assigns.mode do
      :login -> handle_login(socket, email, password)
      :register -> handle_register(socket, email, password)
    end
  end

  defp handle_login(socket, email, password) do
    request = %UserLoginRequest{
      email: email,
      password: password
    }

    case call_auth_service(:login, request) do
      {:ok, _response} ->
        # TODO: セッションにユーザー情報を保存
        {:noreply,
         socket
         |> put_flash(:info, "ログイン成功！")
         |> push_navigate(to: "/owner/dashboard")}

      {:error, reason} ->
        {:noreply,
         socket
         |> assign(:error, "ログイン失敗: #{reason}")
         |> assign(:success, nil)}
    end
  end

  defp handle_register(socket, email, password) do
    # オーナー登録はSHOP_OWNERロールで登録
    request = %UserRegistrationRequest{
      email: email,
      password: password,
      role: Role.value(:SHOP_OWNER)
    }

    case call_auth_service(:register, request) do
      {:ok, _response} ->
        # 登録成功後、ショップ登録画面へ
        {:noreply,
         socket
         |> put_flash(:info, "オーナー登録成功！続いてショップ情報を登録してください。")
         |> push_navigate(to: "/owner/shop/register")}

      {:error, reason} ->
        {:noreply,
         socket
         |> assign(:error, "登録失敗: #{reason}")
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
    <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-100 to-indigo-200">
      <div class="max-w-md w-full bg-white rounded-lg shadow-lg p-8">
        <div class="text-center mb-6">
          <span class="inline-block px-3 py-1 bg-purple-100 text-purple-800 rounded-full text-sm font-medium">
            ショップオーナー向け
          </span>
        </div>

        <h2 class="text-3xl font-bold text-center text-gray-800 mb-8">
          {if @mode == :login, do: "オーナーログイン", else: "オーナー登録"}
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
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent"
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
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              placeholder="••••••••"
              minlength="8"
              required
            />
          </div>

          <button
            type="submit"
            class="w-full bg-purple-600 text-white py-2 px-4 rounded-md hover:bg-purple-700 transition-colors font-medium"
          >
            {if @mode == :login, do: "ログイン", else: "オーナー登録"}
          </button>
        </form>

        <%= if @mode == :login do %>
          <div class="mt-4 text-center">
            <.link
              href="/auth/password-reset"
              class="text-sm text-purple-600 hover:text-purple-800"
            >
              パスワードを忘れた方はこちら
            </.link>
          </div>
        <% end %>

        <div class="mt-6 text-center">
          <button
            phx-click="toggle_mode"
            class="text-purple-600 hover:text-purple-800 text-sm"
          >
            <%= if @mode == :login do %>
              新規オーナー登録はこちら
            <% else %>
              すでにオーナーアカウントをお持ちの方
            <% end %>
          </button>
        </div>

        <div class="mt-8 pt-6 border-t border-gray-200 text-center">
          <p class="text-sm text-gray-500 mb-2">お買い物をされる方は</p>
          <.link
            href="/auth"
            class="text-blue-600 hover:text-blue-800 text-sm font-medium"
          >
            カスタマーログインへ →
          </.link>
        </div>
      </div>
    </div>
    """
  end
end
