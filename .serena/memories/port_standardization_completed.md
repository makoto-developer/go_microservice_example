# ポート番号標準化完了

## 実施日時
2026-01-11

## 作業内容

すべてのサービスとインフラのポート番号を20000-21000範囲に統一しました。

### ポート割り当て一覧

#### フロントエンド
- Phoenix Web: **20200** (http://localhost:20200/auth)

#### マイクロサービス (gRPC)
- Auth Service: **20100**
- Shop Service: **20101**
- Customer Service: **20102**
- Inventory Service: **20103**
- Order Service: **20104**
- Payment Service: **20105**
- Shipping Service: **20106**
- Notification Service: **20107**
- Review Service: **20108**
- Chat Service (gRPC): **20109**
- Chat Service (WebSocket): **20110**
- Search Service: **20111**
- Admin Service: **20112**

#### インフラストラクチャ
- PostgreSQL: **20000**
- Redis: **20001**
- RabbitMQ (AMQP): **20002**
- RabbitMQ (Management UI): **20003** (http://localhost:20003)
- MailHog (SMTP): **20004**
- MailHog (UI): **20005** (http://localhost:20005)

#### モックサービス
- Mock Stripe: **20010**
- Mock SendGrid: **20011**
- Mock FCM: **20012**
- Mock Elasticsearch: **20013**
- Mock Carriers: **20014**

### 更新したファイル

1. **Makefile**
   - phoenix, phoenix-bg コマンドでPORT=20200を使用
   - dev コマンドの「Access URLs」セクションを更新
   - status コマンドの「Access URLs」セクションを更新
   - init コマンドの「Next steps」を更新

2. **.env**
   - PHOENIX_PORT=20200を追加（既存の設定を確認）

### 利点

1. **ポート番号の衝突回避**: 20000-21000の専用範囲を使用
2. **順番の整理**: サービス種別ごとに連番で管理
3. **開発環境の明確化**: すべてのサービスが20000番台に統一

### 開発環境起動コマンド

```bash
# 全体起動
make dev

# アクセスURL
# - Phoenix Web: http://localhost:20200/auth
# - RabbitMQ管理画面: http://localhost:20003
# - MailHog UI: http://localhost:20005
```

### 今後の注意点

- 新しいサービスを追加する場合は、20100番台の空き番号（20113以降）を使用
- Docker ComposeファイルやKubernetes manifestsでもこのポート番号を統一すること
