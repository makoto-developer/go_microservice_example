defmodule ShopMallWebWeb.PageController do
  use ShopMallWebWeb, :controller

  def home(conn, _params) do
    render(conn, :home)
  end
end
