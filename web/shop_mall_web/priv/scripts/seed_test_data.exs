# テストデータ投入スクリプト
# 実行方法: cd /Users/user/work/repositories/github.com/makoto-developer/go_microservice_example/web/shop_mall_web && mix run priv/scripts/seed_test_data.exs

alias ShopService.V1.{
  ShopService.Stub,
  RegisterShopRequest,
  RegisterProductRequest,
  ToggleProductPublishRequest
}

defmodule SeedTestData do
  def run do
    IO.puts("=== テストデータ投入開始 ===\n")

    # 1. ショップを登録
    IO.puts("1. ショップを登録中...")
    shop_id = register_shop()

    if shop_id do
      IO.puts("   ✓ ショップ登録完了: #{shop_id}\n")

      # 2. 商品を登録
      IO.puts("2. 商品を登録中...")
      product_ids = register_products(shop_id)

      IO.puts("   ✓ #{length(product_ids)}件の商品を登録しました\n")

      # 3. 商品を公開
      IO.puts("3. 商品を公開中...")
      publish_products(product_ids)

      IO.puts("   ✓ 全ての商品を公開しました\n")
      IO.puts("=== テストデータ投入完了 ===")
    else
      IO.puts("   ✗ ショップ登録に失敗しました")
    end
  end

  defp register_shop do
    channel = get_shop_channel()

    request = %RegisterShopRequest{
      owner_id: "admin-user-id",
      name: "テクノショップ",
      description: "最新テクノロジー商品を取り扱うショップです",
      logo_url: "",
      owner_name: "山田太郎",
      phone_number: "03-1234-5678",
      business_hours: "10:00-20:00",
      return_policy: "30日間返品可能",
      categories: ["electronics", "gadgets"]
    }

    case Stub.register_shop(channel, request) do
      {:ok, response} ->
        IO.inspect(response, label: "   RegisterShop Response")
        response.shop_id

      {:error, error} ->
        IO.puts("   エラー: #{inspect(error)}")
        nil
    end
  end

  defp register_products(shop_id) do
    channel = get_shop_channel()

    products = [
      %{
        name: "ワイヤレスイヤホン Pro",
        description: "高音質Bluetooth 5.0対応のワイヤレスイヤホン。ノイズキャンセリング機能付き",
        price: "29800",
        category: "electronics",
        stock_quantity: 50
      },
      %{
        name: "スマートウォッチ X1",
        description: "健康管理機能付きスマートウォッチ。心拍数・歩数計測、睡眠モニタリング対応",
        price: "45000",
        category: "electronics",
        stock_quantity: 30
      },
      %{
        name: "カジュアルTシャツ",
        description: "100%コットンの快適な着心地。カラーバリエーション豊富",
        price: "3980",
        category: "fashion",
        stock_quantity: 100
      }
    ]

    Enum.map(products, fn product ->
      request = %RegisterProductRequest{
        shop_id: shop_id,
        name: product.name,
        description: product.description,
        price: product.price,
        stock_quantity: product.stock_quantity,
        category: product.category,
        weight: "",
        size: "",
        jan_code: "",
        tags: []
      }

      case Stub.register_product(channel, request) do
        {:ok, response} ->
          IO.puts("   ✓ #{product.name} を登録しました (ID: #{response.product_id})")
          response.product_id

        {:error, error} ->
          IO.puts("   ✗ #{product.name} の登録に失敗: #{inspect(error)}")
          nil
      end
    end)
    |> Enum.filter(&(&1 != nil))
  end

  defp publish_products(product_ids) do
    channel = get_shop_channel()

    Enum.each(product_ids, fn product_id ->
      request = %ToggleProductPublishRequest{product_id: product_id}

      case Stub.toggle_product_publish(channel, request) do
        {:ok, _response} ->
          IO.puts("   ✓ 商品ID #{product_id} を公開しました")

        {:error, error} ->
          IO.puts("   ✗ 商品ID #{product_id} の公開に失敗: #{inspect(error)}")
      end
    end)
  end

  defp get_shop_channel do
    host = System.get_env("SHOP_SERVICE_HOST", "localhost")
    port = String.to_integer(System.get_env("SHOP_SERVICE_PORT", "20101"))
    {:ok, channel} = GRPC.Stub.connect("#{host}:#{port}")
    channel
  end
end

SeedTestData.run()
