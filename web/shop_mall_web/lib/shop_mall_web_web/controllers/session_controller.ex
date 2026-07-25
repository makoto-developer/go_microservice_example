defmodule ShopMallWebWeb.SessionController do
  @moduledoc """
  LiveView からはクッキーセッションに書き込めないため、
  ログイン成功後に署名付きトークン経由でこのコントローラを通し、
  Plug セッションへ user_id を保存する。
  """
  use ShopMallWebWeb, :controller

  # トークンの有効期間(秒)。ログイン直後のリダイレクトにだけ使うので短くて良い
  @max_age 60

  def establish(conn, %{"token" => token}) do
    case Phoenix.Token.verify(ShopMallWebWeb.Endpoint, "user session", token, max_age: @max_age) do
      {:ok, %{user_id: user_id, access_token: access_token, refresh_token: refresh_token}} ->
        conn
        |> put_session(:user_id, user_id)
        |> put_session(:access_token, access_token)
        |> put_session(:refresh_token, refresh_token)
        |> redirect(to: "/dashboard")

      # 旧形式(user_id のみ)のトークンにも対応
      {:ok, user_id} when is_binary(user_id) ->
        conn
        |> put_session(:user_id, user_id)
        |> redirect(to: "/dashboard")

      {:error, _reason} ->
        conn
        |> put_flash(:error, "ログインの有効期限が切れました。もう一度お試しください。")
        |> redirect(to: "/auth")
    end
  end

  def establish(conn, _params), do: redirect(conn, to: "/auth")

  def logout(conn, _params) do
    # 認証サービス側のセッション(リフレッシュトークン)も失効させる(best effort)
    with refresh_token when is_binary(refresh_token) <- get_session(conn, :refresh_token),
         {:ok, channel} <- auth_channel() do
      AuthService.V1.AuthService.Stub.logout(channel, %AuthService.V1.UserLogoutRequest{
        refresh_token: refresh_token
      })
    end

    conn
    |> configure_session(drop: true)
    |> put_flash(:info, "ログアウトしました")
    |> redirect(to: "/auth")
  end

  defp auth_channel do
    host = System.get_env("AUTH_SERVICE_HOST", "localhost")
    port = System.get_env("AUTH_SERVICE_PORT", "22100")
    GRPC.Stub.connect("#{host}:#{port}")
  end
end
