import Config

# We don't run a server during test. If one is required,
# you can enable the server option below.
config :shop_mall_web, ShopMallWebWeb.Endpoint,
  http: [ip: {127, 0, 0, 1}, port: 4002],
  secret_key_base: "mCXuW8R7NKKoC2Q0q7D/1/wF8qvW4kOpzcI6eTuWGR1zETb8Uk58WzJjYuWBzkpd",
  server: false

# Print only warnings and errors during test
config :logger, level: :warning

# Initialize plugs at runtime for faster test compilation
config :phoenix, :plug_init_mode, :runtime

# Enable helpful, but potentially expensive runtime checks
config :phoenix_live_view,
  enable_expensive_runtime_checks: true
