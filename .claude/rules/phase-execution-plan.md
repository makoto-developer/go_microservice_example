# Phase別実行計画

このドキュメントは、12サービスの開発を4つのPhaseに分けて実行するための詳細計画を定義します。

---

## 全体構成

| Phase | サービス数 | サービス名 | トークン見積 | 並行開発 |
|-------|----------|-----------|------------|---------|
| Phase 1 | 2 | Auth, Shop | 6,300 | 可能 |
| Phase 2 | 4 | Customer, Inventory, Order, Payment | 12,600 | 部分的 |
| Phase 3 | 3 | Shipping, Notification, Review | 9,450 | 可能 |
| Phase 4 | 3 | Chat, Search, Admin | 9,450 | 可能 |

**合計**: 37,800トークン（200,000トークンの18.9%）

---

## Phase 1: 基盤サービス ✅ 完了

### 完成済みサービス

#### 1. Auth Service
- DSL定義: 240行
- ステータス: ✅ 完了

#### 2. Shop Service
- DSL定義: 388行
- ステータス: ✅ 完了

### トークン消費実績
- 見積: 6,300トークン
- 実績: 約3,000トークン（DSL作成のみ）

---

## Phase 2: コアサービス（次のステップ）

### サービス依存関係

```
Auth Service (完成)
    ↓
Customer Service → Order Service → Payment Service
    ↓                  ↓
Inventory Service -----┘
    ↑
Shop Service (完成)
```

### 実行順序

#### Step 1: 独立サービス（並行実行可能）
```
Customer Service  ←→  Inventory Service
```
**依存関係**: Auth Service（完成済み）のみ

**並行開発理由**:
- Customer Serviceは顧客情報管理のみ
- Inventory Serviceは在庫管理のみ
- 相互依存なし

#### Step 2: Order Service
**依存関係**: Customer Service, Inventory Service, Shop Service

**実装内容**:
- 注文作成（Customer参照）
- 在庫引き当て（Inventory連携）
- ショップへの通知（Shop連携）

#### Step 3: Payment Service
**依存関係**: Order Service

**実装内容**:
- 決済処理
- 注文ステータス更新
- Saga パターン実装

### 各サービスの詳細計画

---

#### 2.1 Customer Service

**要件定義**: `docs/requirements/03_customer_service.md`

**DSL定義内容**:
```kotlin
microservice CustomerService {
  entity Customer {
    id: UUID
    user_id: UUID  // Auth Service参照
    name: string
    phone: string
    addresses: list<Address>
  }

  entity Address {
    id: UUID
    customer_id: UUID
    postal_code: string
    prefecture: string
    city: string
    address1: string
    address2: string
    is_default: boolean
  }

  usecase RegisterCustomer { ... }
  usecase UpdateCustomerInfo { ... }
  usecase ManageAddress { ... }
}
```

**推定行数**: 150-200行

**トークン消費**: ~1,500

---

#### 2.2 Inventory Service

**要件定義**: `docs/requirements/04_inventory_service.md`

**DSL定義内容**:
```kotlin
microservice InventoryService {
  entity Inventory {
    id: UUID
    product_id: UUID  // Shop Service参照
    shop_id: UUID
    quantity: int
    reserved_quantity: int
  }

  entity Reservation {
    id: UUID
    inventory_id: UUID
    order_id: UUID
    quantity: int
    status: ReservationStatus
    expires_at: timestamp
  }

  usecase ReserveStock { ... }
  usecase ReleaseStock { ... }
  usecase UpdateInventory { ... }

  events {
    publish StockReserved { ... }
    publish StockReleased { ... }
  }
}
```

**推定行数**: 180-220行

**トークン消費**: ~2,000

---

#### 2.3 Order Service

**要件定義**: `docs/requirements/05_order_service.md`

**DSL定義内容**:
```kotlin
microservice OrderService {
  entity Order {
    id: UUID
    customer_id: UUID
    order_number: string
    status: OrderStatus
    total_amount: decimal
  }

  entity OrderItem {
    id: UUID
    order_id: UUID
    product_id: UUID
    quantity: int
    price: decimal
  }

  usecase CreateOrder {
    // Saga パターン実装
    // 1. 在庫引き当て（Inventory）
    // 2. 注文作成
    // 3. 決済準備（Payment）
  }

  usecase CancelOrder {
    // Compensating Transaction
    // 1. 決済キャンセル
    // 2. 在庫解放
    // 3. 注文キャンセル
  }

  events {
    publish OrderCreated { ... }
    publish OrderCancelled { ... }
    subscribe PaymentCompleted { ... }
  }
}
```

**推定行数**: 250-300行

**トークン消費**: ~3,000

---

#### 2.4 Payment Service

**要件定義**: `docs/requirements/06_payment_service.md`

**DSL定義内容**:
```kotlin
microservice PaymentService {
  entity Payment {
    id: UUID
    order_id: UUID
    amount: decimal
    status: PaymentStatus
    payment_method: PaymentMethod
  }

  entity Transaction {
    id: UUID
    payment_id: UUID
    transaction_id: string  // 外部決済システムID
    status: TransactionStatus
  }

  usecase ProcessPayment {
    // 外部API連携（Stripe等）
    // manual/payment/ にカスタムロジック実装
  }

  usecase RefundPayment { ... }

  events {
    publish PaymentCompleted { ... }
    publish PaymentFailed { ... }
    publish PaymentRefunded { ... }
  }
}
```

**推定行数**: 200-250行

**トークン消費**: ~2,500

**カスタムロジック**:
- `manual/payment/stripe_client.go` - Stripe API連携
- `manual/payment/payment_validator.go` - 決済検証

---

### Phase 2 実行手順

#### Step 1: Customer & Inventory（並行実行）

```bash
# Customer Service
1. cat docs/requirements/03_customer_service.md
2. DSL定義作成（150-200行）
3. ./scripts/mps-generate.sh customer-service

# Inventory Service（並行）
1. cat docs/requirements/04_inventory_service.md
2. DSL定義作成（180-220行）
3. ./scripts/mps-generate.sh inventory-service
```

**チェックポイント**:
- Customer Serviceのエンティティ定義完了
- Inventory Serviceのイベント定義完了
- コンパイル成功

#### Step 2: Order Service

```bash
1. cat docs/requirements/05_order_service.md
2. DSL定義作成（250-300行）
   - Customer, Inventory依存関係を明記
   - Saga パターン定義
3. ./scripts/mps-generate.sh order-service
```

**チェックポイント**:
- Saga フロー定義完了
- イベント発行/購読定義完了
- コンパイル成功

#### Step 3: Payment Service

```bash
1. cat docs/requirements/06_payment_service.md
2. DSL定義作成（200-250行）
3. ./scripts/mps-generate.sh payment-service
4. カスタムロジック実装
   - manual/payment/stripe_client.go
   - manual/payment/payment_validator.go
```

**チェックポイント**:
- 外部API連携インターフェース定義完了
- イベント定義完了
- カスタムロジック実装完了

---

## Phase 3: 拡張サービス1

### サービス構成

| サービス | 依存関係 | 並行実行 |
|---------|---------|---------|
| Shipping Service | Order | - |
| Notification Service | すべて | ✅ 可能 |
| Review Service | Customer, Shop | ✅ 可能 |

### 実行順序

#### Step 1: Notification & Review（並行実行）
```
Notification Service  ←→  Review Service
```

#### Step 2: Shipping Service
```
Shipping Service  ←  Order Service
```

### トークン見積
- Notification: ~2,000
- Review: ~2,500
- Shipping: ~3,000
- **合計**: 7,500トークン

---

## Phase 4: 拡張サービス2

### サービス構成

| サービス | 依存関係 | 並行実行 |
|---------|---------|---------|
| Chat Service | Customer, Shop | ✅ 可能 |
| Search Service | Shop, Product | ✅ 可能 |
| Admin Service | すべて | - |

### 実行順序

#### Step 1: Chat & Search（並行実行）
```
Chat Service  ←→  Search Service
```

#### Step 2: Admin Service
```
Admin Service  ←  すべてのサービス
```

### トークン見積
- Chat: ~3,000
- Search: ~3,500
- Admin: ~4,000
- **合計**: 10,500トークン

---

## 並行開発の実行戦略

### 並行実行可能な組み合わせ

#### Phase 2
```
[Customer Service] + [Inventory Service]  →  並行実行
↓
[Order Service]  →  順次実行
↓
[Payment Service]  →  順次実行
```

#### Phase 3
```
[Notification Service] + [Review Service]  →  並行実行
↓
[Shipping Service]  →  順次実行
```

#### Phase 4
```
[Chat Service] + [Search Service]  →  並行実行
↓
[Admin Service]  →  順次実行
```

### 並行実行のメリット

| 実行方法 | 所要時間（相対） | トークン効率 |
|---------|---------------|------------|
| 順次実行 | 100% | 同じ |
| 並行実行 | 50-60% | 同じ |

**時間短縮効果**: 40-50%削減

---

## 品質チェックリスト

### 各サービス完成時の確認事項

#### 1. DSL定義
- [ ] エンティティ定義完了
- [ ] Enum定義完了
- [ ] ユースケース定義完了
- [ ] gRPC定義完了
- [ ] 依存関係明記
- [ ] イベント定義（必要な場合）

#### 2. コード生成
- [ ] `./scripts/mps-generate.sh <service>` 実行成功
- [ ] 生成コードのディレクトリ構造確認
- [ ] `go.mod` 生成確認

#### 3. カスタムロジック（必要な場合）
- [ ] `manual/<service>/` にファイル作成
- [ ] 外部API連携実装
- [ ] 複雑なバリデーション実装

#### 4. ドキュメント
- [ ] `docs/PROJECT_STATUS.md` 更新
- [ ] トークン消費記録

---

## トラブルシューティング

### 依存関係エラー

**症状**: サービス間の依存関係が不明確

**対処**:
1. `docs/requirements/README.md` のサービス依存図を確認
2. DSL定義で `dependencies` セクションに明記
3. イベント駆動の場合は `events` セクションで定義

### トークン超過の危険

**症状**: トークン消費が見積を大幅に超過

**対処**:
1. 生成コードを読まない（DSL定義のみ）
2. 要件定義は機能要件のみ読む
3. 同じファイルを複数回読まない
4. Serena Memoryに要約を保存

### DSL定義の肥大化

**症状**: 1サービスのDSL定義が400行超え

**対処**:
1. エンティティを分割
2. ユースケースを整理
3. 冗長なコメントを削除
4. 目標: 100-300行以内

---

## 進捗追跡

### Phase完了の定義

各Phaseは以下の条件をすべて満たした場合に完了とする：

1. **DSL定義作成**: すべてのサービスの `.model` ファイル作成
2. **コード生成**: すべてのサービスで `./scripts/mps-generate.sh` 実行成功
3. **カスタムロジック**: 必要なサービスで `manual/` 実装完了
4. **ドキュメント更新**: `docs/PROJECT_STATUS.md` 更新

### トークン消費記録

各Phase完了時に以下を記録：

```markdown
## Phase X完了

- **見積トークン**: X,XXX
- **実績トークン**: Y,YYY
- **削減率**: ZZ%
- **所要時間**: N分
```

---

## まとめ

### Phase 2-4の実行計画

1. **Phase 2**: Customer/Inventory並行 → Order → Payment（12,600トークン）
2. **Phase 3**: Notification/Review並行 → Shipping（7,500トークン）
3. **Phase 4**: Chat/Search並行 → Admin（10,500トークン）

### 全体目標

- **トークン合計**: 37,800トークン
- **削減率**: 90%（手動実装180,000トークンから）
- **1セッション完了**: 可能

この計画に従うことで、効率的かつ正確に12サービスを開発できます。
