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
      {:ok, user_id} ->
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
    conn
    |> configure_session(drop: true)
    |> put_flash(:info, "ログアウトしました")
    |> redirect(to: "/auth")
  end
end
