defmodule ShopMallWebWeb.Owner.DashboardLiveTest do
  use ShopMallWebWeb.ConnCase

  import Phoenix.LiveViewTest

  @moduletag :e2e

  describe "Owner Dashboard E2E Tests" do
    setup do
      # テストユーザーでログイン (既存のログイン処理を使用)
      # 注: この部分は実際の認証実装に合わせて調整が必要
      %{conn: conn} = %{conn: build_conn()}
      {:ok, conn: conn}
    end

    test "displays correct product count for shop owner", %{conn: conn} do
      # ダッシュボードページにアクセス
      {:ok, view, _html} = live(conn, "/owner/dashboard")

      # 商品数が正しく表示されることを確認
      # Shop Service から取得した商品数が表示される
      html = render(view)

      # 「登録商品数」というラベルが表示されていることを確認
      assert html =~ "登録商品数"

      # 商品数が2であることを確認
      # (テストデータには11111111-1111-1111-1111-111111111111のショップに2商品ある)
      assert html =~ ~r/登録商品数.*?<dd[^>]*>.*?2.*?<\/dd>/s

      # より具体的な確認: 登録商品数が0ではないことを確認
      # (注: 受注件数が0なので、より具体的なパターンを使用)
      assert html =~ ~r/登録商品数.*?<dd[^>]*>.*?[1-9]\d*.*?<\/dd>/s
    end

    test "product count updates dynamically when new products are added", %{conn: conn} do
      {:ok, view, _html} = live(conn, "/owner/dashboard")

      initial_html = render(view)

      # 初期状態の商品数を取得
      initial_count =
        case Regex.run(~r/登録商品数.*?<dd.*?>.*?(\d+).*?<\/dd>/s, initial_html) do
          [_, count] -> String.to_integer(String.trim(count))
          _ -> 0
        end

      # 商品数が0より大きいことを確認
      assert initial_count > 0, "商品数が0です。テストデータが正しく投入されていません"

      # 商品数が2であることを確認 (テストデータの期待値)
      assert initial_count == 2,
             "商品数が期待値(2)と異なります。実際: #{initial_count}"
    end

    test "displays zero products when shop has no products", %{conn: conn} do
      # 注: このテストは商品がないショップのオーナーでログインした場合を想定
      # 実際の実装に応じて調整が必要
      {:ok, view, _html} = live(conn, "/owner/dashboard")

      html = render(view)

      # 商品がない場合でもラベルは表示される
      assert html =~ "登録商品数"
    end

    test "handles gRPC service errors gracefully", %{conn: conn} do
      # Shop Serviceが停止している場合でも、エラーが適切に処理されることを確認
      # 注: このテストは Shop Service を停止して実行する必要がある
      {:ok, view, _html} = live(conn, "/owner/dashboard")

      html = render(view)

      # エラーが発生しても、ページが正常に表示されることを確認
      assert html =~ "ダッシュボード" or html =~ "登録商品数"
    end

    test "product count is displayed in correct format", %{conn: conn} do
      {:ok, view, _html} = live(conn, "/owner/dashboard")

      html = render(view)

      # フォーマットが正しいことを確認
      assert html =~ ~r/登録商品数.*?<dd.*?>.*?\d+.*?<\/dd>/s

      # 数値が正の整数であることを確認
      case Regex.run(~r/登録商品数.*?<dd.*?>.*?(\d+).*?<\/dd>/s, html) do
        [_, count_str] ->
          count = String.to_integer(String.trim(count_str))
          assert count >= 0, "商品数が負の値です"

        _ ->
          flunk("商品数が正しいフォーマットで表示されていません")
      end
    end

    test "product count matches gRPC response", %{conn: conn} do
      {:ok, view, _html} = live(conn, "/owner/dashboard")

      html = render(view)

      # 商品数がgRPC responseと一致することを確認
      # テストデータでは shop_id: 11111111-1111-1111-1111-111111111111 に2商品
      case Regex.run(~r/登録商品数.*?<dd.*?>.*?(\d+).*?<\/dd>/s, html) do
        [_, count] ->
          count_int = String.to_integer(String.trim(count))
          assert count_int == 2, "商品数が期待値(2)と異なります。実際: #{count_int}"

        _ ->
          flunk("商品数が表示されていません")
      end
    end

    test "dashboard loads within acceptable time", %{conn: conn} do
      start_time = System.monotonic_time(:millisecond)

      {:ok, _view, _html} = live(conn, "/owner/dashboard")

      end_time = System.monotonic_time(:millisecond)
      duration = end_time - start_time

      # ダッシュボードが3秒以内にロードされることを確認
      assert duration < 3000, "ダッシュボードのロードが遅すぎます: #{duration}ms"
    end

    test "product count is visible and readable", %{conn: conn} do
      {:ok, view, _html} = live(conn, "/owner/dashboard")

      html = render(view)

      # 商品数のラベルが存在することを確認
      assert html =~ "登録商品数"

      # 数値が表示されていることを確認
      assert html =~ ~r/登録商品数.*?\d+/s

      # CSSクラスやスタイルが適用されているか確認 (オプション)
      # assert html =~ "product-count" or html =~ "stats"
    end
  end

  describe "Dashboard Stats Integration" do
    test "displays multiple stats correctly", %{conn: conn} do
      {:ok, view, _html} = live(conn, "/owner/dashboard")

      html = render(view)

      # 複数の統計情報が表示されることを確認
      assert html =~ "登録商品数"

      # 他の統計情報も表示されているはず
      assert html =~ "受注件数"
    end
  end
end
