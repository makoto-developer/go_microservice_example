# 🚀 Quick Start Guide

Go Microserviceを最速で起動する方法

---

## ⚡️ 1分で起動

### ステップ1: データベース起動

```bash
cd infrastructure/docker
docker compose up -d
```

### ステップ2: 全サービス起動

```bash
cd ../..
make up
```

たったこれだけ！ ✨

---

## 📋 主要コマンド

### 起動・停止

```bash
make up          # 全サービス起動
make down        # 全サービス停止
make restart     # 再起動
make status      # 稼働状況確認
```

### ログ確認

```bash
make logs        # 最新ログ表示
make logs-follow # ログをリアルタイム表示
```

### テスト実行

```bash
make test        # 全統合テスト実行
make test-auth   # Auth認証テスト
make test-order  # 注文フローテスト
make test-e2e    # E2Eテスト
```

### ビルド

```bash
make build       # 全サービスビルド
```

### ヘルプ

```bash
make help        # 全コマンド一覧
```

---

## 🎯 動作確認

### サービス稼働確認

```bash
make status
```

期待される出力：
```
✅ Auth Service          Running (PID: 47518, Port: 22100)
✅ Shop Service          Running (PID: 59025, Port:  4000)
✅ Customer Service      Running (PID: 78668, Port: 22102)
...
Summary: 12/12 services running
```

### データベース確認

```bash
make db-status
```

### プロセス確認

```bash
make ps
```

---

## 🧪 テスト実行

### 全テスト（推奨）

```bash
make test
```

実行時間: 約15-20分

### 個別テスト

```bash
# Auth認証テスト（30秒）
make test-auth

# 注文フローテスト（2-3分）
make test-order

# E2Eテスト（5-10分）
make test-e2e
```

---

## 📊 ダッシュボード

### サービスダッシュボード表示

```bash
make dashboard
```

または：
```bash
cat RUNNING_SERVICES_DASHBOARD.md
```

---

## 🔧 トラブルシューティング

### サービスが起動しない場合

```bash
# ログ確認
make logs

# 個別ログ確認
tail -f /tmp/order-service.log
tail -f /tmp/payment-service.log
```

### ポートが使用中の場合

```bash
# ポート使用状況確認
lsof -i :22100-22111

# サービス停止
make down

# 再起動
make up
```

### データベース接続エラー

```bash
# Dockerコンテナ確認
docker ps | grep postgres

# コンテナ再起動
cd infrastructure/docker
docker compose restart

# ログ確認
docker compose logs postgres_auth
```

---

## 📁 ディレクトリ構成

```
go_microservice_example/
├── Makefile                    # 👈 メインコマンド
├── QUICKSTART.md               # このファイル
├── scripts/                    # 自動化スクリプト
│   ├── start_all_services.sh
│   ├── stop_all_services.sh
│   ├── check_all_services.sh
│   └── build_all_services.sh
├── infrastructure/docker/      # Docker設定
│   └── docker-compose.yml
├── microservices/              # サービス実装
├── simple-servers/             # 簡易サービス
└── tests/                      # テストスイート
```

---

## 🎓 次のステップ

### 1. サービス稼働確認

```bash
make status
```

### 2. テスト実行

```bash
make test
```

### 3. ドキュメント確認

- `README.md` - プロジェクト概要
- `RUNNING_SERVICES_DASHBOARD.md` - サービス詳細
- `ALL_SERVICES_STATUS.md` - 実装状況
- `PROJECT_STATUS_FINAL.md` - プロジェクト全体状況

### 4. 開発開始

```bash
# 開発モード（ログ表示付き）
make dev
```

---

## 📞 よくある質問

### Q: 全サービスを一度に起動するには？

```bash
make up
```

### Q: 特定のサービスだけ起動するには？

```bash
cd simple-servers/order
./order-server > /tmp/order-service.log 2>&1 &
```

### Q: ログはどこにありますか？

```bash
# サービスログ
ls /tmp/*-service.log

# Auth Service
cat microservices/auth/auth-server.log

# 全ログ表示
make logs
```

### Q: テストが失敗します

```bash
# サービスが起動しているか確認
make status

# データベースが起動しているか確認
make db-status

# ログでエラー確認
make logs
```

### Q: Shop Service (Phoenix) の起動方法は？

```bash
cd simple-servers/admin  # Phoenixプロジェクト
mix deps.get             # 初回のみ
mix phx.server           # サーバー起動

# ブラウザで確認
open http://localhost:4000
```

---

## 🎉 完了！

これで全12マイクロサービスが稼働しています。

**次のアクション**:
1. `make status` - 稼働確認
2. `make test` - 品質確認
3. `make dashboard` - 詳細確認

---

**Database per Service アーキテクチャが完全に動作しています！** 🚀
