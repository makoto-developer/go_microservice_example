# テスト

## テストフレームワーク
- **testify** - アサーションライブラリ
- **gomock** - モックライブラリ（検討中）
- **httptest** - HTTPテスト

## テスト実行コマンド

### すべてのテスト
```bash
# すべてのテストを実行
go test ./...

# カバレッジ付き
go test -cover ./...
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 特定のサービス
```bash
# Auth Serviceのみ
cd generated/auth-service
go test ./...

# Shop Serviceのみ
cd generated/shop-service
go test ./...
```

### Verbose モード
```bash
go test -v ./...
```

## テストの種類

### 1. ユニットテスト
```go
// usecase/user_registration_test.go
func TestUserRegistration_Execute(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    usecase := NewUserRegistration(mockRepo)
    
    // Act
    result, err := usecase.Execute(ctx, input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### 2. 統合テスト
```bash
# 統合テスト実行（要Docker環境）
go test -tags=integration ./tests/integration/...
```

### 3. E2Eテスト
```bash
# E2Eテスト実行（要Docker環境）
go test ./tests/e2e/...
```

## テスト時の注意事項

### データベース
- 統合テストはテスト用データベースを使用
- テストごとにトランザクションロールバック

### モックサービス
- 外部サービス（Stripe, FCM等）はモックを使用
- MailHogでメール送信をテスト

## Makefileコマンド

```bash
make test              # すべてのテスト実行
make test-coverage     # カバレッジ付きテスト
make test-integration  # 統合テスト
make test-e2e          # E2Eテスト
```

## テストカバレッジ目標
- 目標: 80%以上
- 重要なビジネスロジックは100%を目指す
