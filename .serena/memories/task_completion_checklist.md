# タスク完了時のチェックリスト

## DSL定義作成時

- [ ] DSL定義を作成（`mps-workspace/solutions/<service>/service.model`）
- [ ] エンティティ定義完了
- [ ] Enum定義完了（必要な場合）
- [ ] ユースケース定義完了
- [ ] gRPCサービス定義完了
- [ ] 依存関係明記
- [ ] イベント定義（必要な場合）
- [ ] DSL定義が100-300行以内
- [ ] コード生成実行（`./scripts/mps-generate.sh <service>`）
- [ ] 生成コードのコンパイル確認

## マイグレーションSQL作成時

- [ ] マイグレーションファイル作成（`scripts/migrations/<service>/001_create_tables.sql`）
- [ ] テーブル定義完了
- [ ] インデックス定義完了
- [ ] 制約定義完了（CHECK, UNIQUE, FOREIGN KEY）
- [ ] トリガー定義（必要な場合）
- [ ] 関数定義（必要な場合）
- [ ] テストデータ挿入（開発環境のみ）

## カスタムロジック実装時

- [ ] カスタムロジックファイル作成（`manual/<service>/`）
- [ ] 外部API連携実装（必要な場合）
- [ ] 複雑なバリデーション実装（必要な場合）
- [ ] ビジネスルール実装（必要な場合）
- [ ] エラーハンドリング実装
- [ ] ユニットテスト作成

## インフラ層実装時

- [ ] Repository実装（`generated/<service>/infrastructure/`）
- [ ] RabbitMQ Event Publisher実装（必要な場合）
- [ ] Redis Cache実装（必要な場合）
- [ ] トランザクション対応（必要な場合）
- [ ] エラーハンドリング実装

## main.go実装時

- [ ] main.go作成（`generated/<service>/main.go`）
- [ ] 環境変数読み込み
- [ ] データベース接続
- [ ] Repository初期化
- [ ] Usecase初期化
- [ ] Handler初期化
- [ ] gRPCサーバー起動
- [ ] Graceful Shutdown実装
- [ ] Signal handling実装

## Go Module管理

- [ ] `go.mod` 作成/更新
- [ ] 必要な依存関係追加
  - `github.com/google/uuid`
  - `github.com/lib/pq` (PostgreSQL)
  - `github.com/streadway/amqp` (RabbitMQ)
  - `github.com/shopspring/decimal` (金額計算)
  - `google.golang.org/grpc`
  - `google.golang.org/protobuf`
- [ ] `go mod tidy` 実行

## Docker統合

- [ ] Dockerfile作成（必要な場合）
- [ ] docker-compose.yml更新（必要な場合）
- [ ] 環境変数設定（`.env`）
- [ ] ポート番号割り当て
- [ ] ヘルスチェック設定

## Git操作

- [ ] 変更ファイル確認（`git status`）
- [ ] ステージング（`git add .`）
- [ ] コミットメッセージ作成
  - Type指定（feat, fix, refactor等）
  - 簡潔な説明
  - Co-Authored-By追加
- [ ] コミット（`git commit`）
- [ ] プッシュ（`git push`）※ユーザー承認必要

## ドキュメント更新

- [ ] プロジェクト状況更新（`docs/PROJECT_STATUS.md`）※存在する場合
- [ ] トークン消費記録
- [ ] 完了したサービスをマーク
- [ ] 次のステップを記載

## 動作確認

- [ ] Docker環境起動（`make up`）
- [ ] サービス起動確認（`make ps`）
- [ ] ログ確認（`make logs`）
- [ ] ヘルスチェック（`make health`）
- [ ] データベース接続確認
- [ ] gRPCサービス確認（`grpcurl -plaintext localhost:<port> list`）

## Serena Memory更新

- [ ] 進捗記録（`mcp__serena__write_memory`）
- [ ] 実装パターン記録（参考情報）
- [ ] トークン消費記録
