defmodule ShopMallWebWeb.VerifyEmailLive do
  @moduledoc """
  メールアドレス検証ページ(登録メールのリンクから開く)。
  role パラメータで検証先を振り分ける:
  customer → CustomerAuthService、owner → OwnerAuthService、それ以外 → AuthService。
  """
  use ShopMallWebWeb, :live_view

  @impl true
  def mount(params, _session, socket) do
    token = params["token"] || ""
    role = params["role"] || "user"

    socket = assign(socket, :role, role)

    if token == "" do
      {:ok, assign(socket, :result, {:error, "検証トークンがありません"})}
    else
      {:ok, assign(socket, :result, verify(role, token))}
    end
  end

  defp verify("customer", token) do
    call(fn ch ->
      CustomerAuth.V1.CustomerAuthService.Stub.verify_email(
        ch,
        %CustomerAuth.V1.CustomerVerifyEmailRequest{token: token}
      )
    end)
  end

  defp verify("owner", token) do
    call(fn ch ->
      OwnerAuth.V1.OwnerAuthService.Stub.verify_email(
        ch,
        %OwnerAuth.V1.OwnerVerifyEmailRequest{token: token}
      )
    end)
  end

  defp verify(_other, token) do
    call(fn ch ->
      AuthService.V1.AuthService.Stub.verify_email(
        ch,
        %AuthService.V1.EmailVerificationRequest{token: token}
      )
    end)
  end

  defp call(fun) do
    host = System.get_env("AUTH_SERVICE_HOST", "localhost")
    port = System.get_env("AUTH_SERVICE_PORT", "22100")

    case GRPC.Stub.connect("#{host}:#{port}") do
      {:ok, channel} ->
        case fun.(channel) do
          {:ok, %{success: true}} -> :ok
          {:ok, %{success: false, message: message}} -> {:error, message}
          {:error, %GRPC.RPCError{message: message}} -> {:error, message}
          {:error, reason} -> {:error, inspect(reason)}
        end

      {:error, reason} ->
        {:error, "認証サービスに接続できません: #{inspect(reason)}"}
    end
  end

  @impl true
  def render(assigns) do
    ~H"""
    <div class="min-h-screen flex items-center justify-center bg-gray-100">
      <div class="max-w-md w-full bg-white rounded-lg shadow-lg p-8 text-center">
        <%= case @result do %>
          <% :ok -> %>
            <div class="text-5xl mb-4">✅</div>
            <h2 class="text-2xl font-bold text-gray-800 mb-2">メールアドレスを確認しました</h2>
            <p class="text-gray-600 mb-6">アカウントが有効になりました。ログインしてご利用ください。</p>
            <.link
              navigate={if @role == "owner", do: "/owner/auth", else: "/auth"}
              class="inline-block px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 font-medium"
            >
              ログインへ
            </.link>
          <% {:error, message} -> %>
            <div class="text-5xl mb-4">⚠️</div>
            <h2 class="text-2xl font-bold text-gray-800 mb-2">確認できませんでした</h2>
            <p class="text-gray-600 mb-6">{message}</p>
            <.link navigate="/auth" class="text-blue-600 hover:text-blue-800 text-sm">
              ログイン画面へ戻る
            </.link>
        <% end %>
      </div>
    </div>
    """
  end
end
