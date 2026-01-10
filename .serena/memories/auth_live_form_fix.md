# AuthLive フォーム修正完了

## 実施日時
2026-01-11

## 問題点

パスワード入力時に意図しないリダイレクトが発生していた。

### 原因

`phx-change` イベントが入力フィールドに設定されていたため、1文字入力するたびにサーバーとの通信が発生し、意図しない動作を引き起こしていた。

```elixir
# 問題のあったコード
<input
  type="password"
  phx-change="update_password"
  value={@password}
  ...
/>
```

## 修正内容

### 変更点

1. **`phx-change` イベントを削除**
   - `phx-change="update_email"` を削除
   - `phx-change="update_password"` を削除

2. **`name` 属性を追加**
   - `name="email"` を追加
   - `name="password"` を追加

3. **イベントハンドラーを変更**
   - `update_email` ハンドラーを削除
   - `update_password` ハンドラーを削除
   - `handle_event("submit", params, socket)` でフォームの値を直接取得

4. **state を簡素化**
   - `email` と `password` の state を削除（不要になったため）

### 修正後のコード

```elixir
# フォーム送信時のみ値を取得
@impl true
def handle_event("submit", %{"email" => email, "password" => password}, socket) do
  case socket.assigns.mode do
    :login -> handle_login(socket, email, password)
    :register -> handle_register(socket, email, password)
  end
end

# 入力フィールド
<input
  type="email"
  name="email"
  class="..."
  placeholder="email@example.com"
  required
/>

<input
  type="password"
  name="password"
  class="..."
  placeholder="••••••••"
  required
/>
```

## 効果

1. **パスワード入力時のリダイレクト解消**: 入力中に意図しないイベントが発火しなくなった
2. **パフォーマンス向上**: 1文字入力するたびのサーバー通信がなくなり、フォーム送信時のみ通信
3. **コードの簡素化**: 不要な state とイベントハンドラーが削除され、コードが読みやすくなった

## 動作確認

- Phoenixサーバー再起動: ✅ 成功
- ポート20200で起動: ✅ 確認
- LiveView マウント: ✅ 正常

## アクセス方法

```bash
# ブラウザでアクセス
http://localhost:20200/auth
```

登録画面でパスワードを入力しても、意図しないリダイレクトは発生しなくなりました。
