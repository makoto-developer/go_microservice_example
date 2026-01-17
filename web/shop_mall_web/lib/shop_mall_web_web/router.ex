defmodule ShopMallWebWeb.Router do
  use ShopMallWebWeb, :router

  pipeline :browser do
    plug :accepts, ["html"]
    plug :fetch_session
    plug :fetch_live_flash
    plug :put_root_layout, html: {ShopMallWebWeb.Layouts, :root}
    plug :protect_from_forgery
    plug :put_secure_browser_headers
  end

  pipeline :api do
    plug :accepts, ["json"]
  end

  scope "/", ShopMallWebWeb do
    pipe_through :browser

    get "/", PageController, :home
    live "/auth", AuthLive

    # Customer routes
    live "/dashboard", DashboardLive
    live "/products", ProductListLive
    live "/products/:id", ProductDetailLive

    # Owner routes
    live "/owner/dashboard", Owner.DashboardLive
    live "/owner/products", Owner.ProductListLive
    live "/owner/products/new", Owner.ProductFormLive
    live "/owner/products/:id/edit", Owner.ProductFormLive
  end

  # Other scopes may use custom stacks.
  # scope "/api", ShopMallWebWeb do
  #   pipe_through :api
  # end

  # Enable LiveDashboard in development
  if Application.compile_env(:shop_mall_web, :dev_routes) do
    # If you want to use the LiveDashboard in production, you should put
    # it behind authentication and allow only admins to access it.
    # If your application does not have an admins-only section yet,
    # you can use Plug.BasicAuth to set up some basic authentication
    # as long as you are also using SSL (which you should anyway).
    import Phoenix.LiveDashboard.Router

    scope "/dev" do
      pipe_through :browser

      live_dashboard "/dashboard", metrics: ShopMallWebWeb.Telemetry
    end
  end
end
