# MPS 日本語振る舞いDSL 実装ガイド

**対象者**: MPS開発者
**前提知識**: JetBrains MPS基本操作、Generatorの知識

---

## 実装手順サマリー

```
1. MPSプロジェクト作成 (30分)
   ↓
2. Structure定義 (2時間)
   ↓
3. Editor定義 (1時間)
   ↓
4. Constraints定義 (30分)
   ↓
5. Generator実装 (4時間)
   ↓
6. サンプル作成・テスト (1時間)
```

**合計**: 約9時間（1-2日）

---

## STEP 1: MPSプロジェクト作成

### 1.1 JetBrains MPSを起動

1. JetBrains MPSを開く
2. `File > New Project > Language`
3. プロジェクト名: `japanese-behavior-dsl`
4. ロケーション: `/path/to/mps-workspace/languages/japanese-behavior-dsl`

### 1.2 言語構造の作成

プロジェクト内に以下を作成:

```
japanese-behavior-dsl/
├── structure/        # 構文構造定義
├── editor/           # エディタ定義
├── generator/        # コード生成器
├── constraints/      # 制約・バリデーション
└── sandbox/          # サンプル・テスト
```

---

## STEP 2: Structure定義（構文構造）

### 2.1 基本コンセプトの定義

MPSの`Structure`モデルで以下のコンセプトを定義:

#### 2.1.1 ユースケースコンセプト

```
Concept: UseCase
Properties:
  - name: string (ユースケース名)
  - id: string (ID: UC_AUTH_001等)
  - description: string (概要)
  - priority: enum {高, 中, 低}

Children:
  - input: InputDefinition [0..1]
  - output: OutputDefinition [0..1]
  - flow: ProcessFlow [0..1]
  - errors: ErrorDefinition [0..*]
  - testCases: TestCase [0..*]

Implements:
  - INamedConcept
```

#### 2.1.2 入力定義コンセプト

```
Concept: InputDefinition
Children:
  - fields: FieldDefinition [1..*]
```

#### 2.1.3 フィールド定義コンセプト

```
Concept: FieldDefinition
Properties:
  - name: string (フィールド名: "メールアドレス"等)
  - technicalName: string (技術名: "email"等)
  - description: string (説明)

References:
  - type: TypeReference [1..1]
```

#### 2.1.4 型参照コンセプト

```
Concept: TypeReference
Properties:
  - typeName: enum {文字列, 整数, 真偽値, UUID, 日時, 小数}

Derived Properties:
  - goType: string (computed)
    - 文字列 → "string"
    - 整数 → "int64"
    - UUID → "uuid.UUID"

  - protoType: string (computed)
    - 文字列 → "string"
    - 整数 → "int64"
    - UUID → "string"

  - tsType: string (computed)
    - 文字列 → "string"
    - 整数 → "number"
    - UUID → "string"
```

#### 2.1.5 処理フローコンセプト

```
Concept: ProcessFlow
Children:
  - steps: FlowStep [1..*]
```

#### 2.1.6 フローステップコンセプト

```
Concept: FlowStep
Properties:
  - stepNumber: int (ステップ番号)
  - title: string (ステップタイトル)
  - description: string (説明)

Children:
  - statements: Statement [0..*]
```

#### 2.1.7 ステートメント（基底）

```
Abstract Concept: Statement
(すべての処理文の基底クラス)

Concrete Concepts:
  - IfStatement (条件分岐)
  - AssignmentStatement (代入)
  - DatabaseQueryStatement (DB検索)
  - DatabaseSaveStatement (DB保存)
  - TransactionStatement (トランザクション)
  - ErrorReturnStatement (エラー返却)
  - ReturnStatement (結果返却)
```

#### 2.1.8 条件分岐ステートメント

```
Concept: IfStatement extends Statement
Properties:
  - conditionText: string (条件文: "メールアドレスが正しい形式でない")

Children:
  - condition: Expression [1..1]
  - thenBlock: Statement [1..*]
  - elseBlock: Statement [0..*]
```

#### 2.1.9 式（Expression）

```
Abstract Concept: Expression

Concrete Concepts:
  - VariableReference (変数参照)
  - BinaryExpression (二項演算)
  - MethodCallExpression (メソッド呼び出し)
  - LiteralExpression (リテラル値)
```

#### 2.1.10 代入ステートメント

```
Concept: AssignmentStatement extends Statement
Properties:
  - variableName: string (変数名)
  - description: string (説明)

Children:
  - value: Expression [1..1]
```

#### 2.1.11 データベース検索ステートメント

```
Concept: DatabaseQueryStatement extends Statement
Properties:
  - targetEntity: string (対象エンティティ: "ユーザー"等)
  - searchField: string (検索フィールド: "メールアドレス"等)
  - resultVariable: string (結果変数名)

Children:
  - searchValue: Expression [1..1]

Generates:
  Go: repo.FindBy<SearchField>(ctx, <searchValue>)
```

#### 2.1.12 エラー定義コンセプト

```
Concept: ErrorDefinition
Properties:
  - name: string (エラー名: "メールアドレス重複"等)
  - code: string (エラーコード: "AUTH_E003")
  - httpStatus: int (HTTPステータス: 409)
  - message: string (メッセージ)
```

---

## STEP 3: Editor定義（見た目）

### 3.1 UseCaseエディタ

```
EditorCell for UseCase:

┌──────────────────────────────────────────┐
│ ユースケース: [ユーザー登録        ]     │
│ ID: [UC_AUTH_001]                        │
│ 概要: [新規ユーザーを登録し、メール認証  │
│        リンクを送信する              ]   │
│ 優先度: [高 ▼]                           │
│                                          │
│ 入力 {                                   │
│   [フィールド定義...]                    │
│ }                                        │
│                                          │
│ 出力 {                                   │
│   [フィールド定義...]                    │
│ }                                        │
│                                          │
│ 処理フロー {                             │
│   [ステップ...]                          │
│ }                                        │
│                                          │
│ エラー定義 {                             │
│   [エラー...]                            │
│ }                                        │
└──────────────────────────────────────────┘
```

### 3.2 FlowStepエディタ

```
EditorCell for FlowStep:

┌──────────────────────────────────────────┐
│ ステップ1: "入力値の検証"                │
│ 説明: メールアドレスとパスワードの       │
│       形式を確認する                     │
│                                          │
│ 処理:                                    │
│   もし メールアドレスが正しい形式でない  │
│   場合:                                  │
│     → エラー「メールアドレスの形式が不正 │
│              です」を返す                │
│     → 処理を終了する                     │
└──────────────────────────────────────────┘
```

### 3.3 IfStatementエディタ

```
EditorCell for IfStatement:

もし [条件式                          ] 場合:
  [ステートメント...]
そうでなければ:
  [ステートメント...]
```

### 3.4 色分け・スタイル

| 要素 | 色 | 太さ |
|------|---|------|
| キーワード（もし、場合、そうでなければ） | 青 | Bold |
| 変数名 | 緑 | Normal |
| 文字列リテラル | 茶 | Normal |
| コメント | グレー | Italic |
| エラーメッセージ | 赤 | Normal |

---

## STEP 4: Constraints定義（制約）

### 4.1 バリデーション制約

#### UseCase名の制約

```kotlin
constraint can be child {
  // UseCaseは最上位のみ
  node.parent == null
}

constraint valid name {
  // 名前は2文字以上50文字以下
  node.name.length >= 2 && node.name.length <= 50
}
```

#### ErrorDefinitionの制約

```kotlin
constraint valid error code {
  // エラーコードはパターンに従う: SERVICE_E001
  node.code matches pattern "^[A-Z]+_E[0-9]{3}$"
}

constraint valid http status {
  // HTTPステータスは有効な値
  node.httpStatus in [400, 401, 403, 404, 409, 500, 503]
}
```

---

## STEP 5: Generator実装

### 5.1 Generatorの構造

```
generator/
├── templates/
│   ├── go/
│   │   ├── usecase_impl.template
│   │   ├── error_definition.template
│   │   └── test.template
│   ├── proto/
│   │   └── service.template
│   ├── typescript/
│   │   ├── api_client.template
│   │   └── types.template
│   └── docs/
│       ├── flow_diagram.template
│       └── specification.template
└── mapping_rules/
    ├── usecase_to_go.mps
    ├── usecase_to_proto.mps
    ├── usecase_to_ts.mps
    └── usecase_to_docs.mps
```

### 5.2 Go生成テンプレート例

#### UseCaseからGo usecaseへのマッピング

```kotlin
// Generator Mapping Rule
map UseCase -> GoFile {
  filename = "generated/auth/usecase/" + node.technicalName + ".go"

  template {
    package usecase

    import (
      "context"
      $LOOP node.requiredImports$
        "$IT$"
      $END$
    )

    // $node.name$ ユースケース
    // $node.description$
    type $node.technicalName$Usecase interface {
      Execute(ctx context.Context, input $node.technicalName$Input) ($node.technicalName$Output, error)
    }

    type $node.technicalName$UsecaseImpl struct {
      $LOOP node.dependencies$
        $IT.name$ $IT.type$
      $END$
    }

    func New$node.technicalName$Usecase(
      $LOOP node.dependencies$
        $IT.name$ $IT.type$,
      $END$
    ) $node.technicalName$Usecase {
      return &$node.technicalName$UsecaseImpl{
        $LOOP node.dependencies$
          $IT.name$: $IT.name$,
        $END$
      }
    }

    func (u *$node.technicalName$UsecaseImpl) Execute(
      ctx context.Context,
      input $node.technicalName$Input,
    ) ($node.technicalName$Output, error) {
      $LOOP node.flow.steps$
        $CALL generateStep(IT)$
      $END$
    }

    // エラー定義
    $LOOP node.errors$
      $CALL generateError(IT)$
    $END$
  }
}

// ステップ生成
macro generateStep(step: FlowStep) {
  // $step.title$
  $LOOP step.statements$
    $CALL generateStatement(IT)$
  $END$
}

// ステートメント生成
macro generateStatement(stmt: Statement) {
  $IF stmt is IfStatement$
    if $CALL generateCondition(stmt.condition)$ {
      $LOOP stmt.thenBlock$
        $CALL generateStatement(IT)$
      $END$
    }
    $IF stmt.elseBlock != null$
      else {
        $LOOP stmt.elseBlock$
          $CALL generateStatement(IT)$
        $END$
      }
    $END$
  $ELSEIF stmt is AssignmentStatement$
    $stmt.variableName$ := $CALL generateExpression(stmt.value)$
  $ELSEIF stmt is DatabaseQueryStatement$
    $stmt.resultVariable$, err := u.repo.FindBy$stmt.searchField$(ctx, $CALL generateExpression(stmt.searchValue)$)
    if err != nil {
      return $node.technicalName$Output{}, err
    }
  $ELSEIF stmt is ErrorReturnStatement$
    return $node.technicalName$Output{}, Err$stmt.errorName$
  $END$
}
```

### 5.3 Protobuf生成テンプレート例

```kotlin
// Generator Mapping Rule
map UseCase -> ProtoFile {
  filename = "generated/proto/" + node.serviceName + "/v1/" + node.serviceName + ".proto"

  template {
    syntax = "proto3";

    package $node.serviceName$.v1;

    option go_package = "github.com/.../proto/$node.serviceName$/v1;$node.serviceName$v1";

    // $node.name$ リクエスト
    message $node.technicalName$Request {
      $LOOP node.input.fields WITH INDEX$
        // $IT.description$
        $CALL generateProtoType(IT.type)$ $IT.technicalName$ = $INDEX + 1$;
      $END$
    }

    // $node.name$ レスポンス
    message $node.technicalName$Response {
      $LOOP node.output.fields WITH INDEX$
        // $IT.description$
        $CALL generateProtoType(IT.type)$ $IT.technicalName$ = $INDEX + 1$;
      $END$
    }

    service $node.serviceName$Service {
      // $node.name$
      // $node.description$
      rpc $node.technicalName$($node.technicalName$Request) returns ($node.technicalName$Response);
    }
  }
}

macro generateProtoType(type: TypeReference) {
  $IF type.typeName == "文字列"$
    string
  $ELSEIF type.typeName == "整数"$
    int64
  $ELSEIF type.typeName == "UUID"$
    string
  $END$
}
```

### 5.4 TypeScript生成テンプレート例

```kotlin
// Generator Mapping Rule
map UseCase -> TypeScriptFile {
  filename = "generated/frontend/src/api/" + node.serviceName + ".ts"

  template {
    /**
     * $node.name$ API
     *
     * @description $node.description$
     $LOOP node.input.fields$
     * @param $IT.technicalName$ $IT.description$
     $END$
     $LOOP node.output.fields$
     * @returns $IT.technicalName$ $IT.description$
     $END$
     $LOOP node.errors$
     * @throws $IT.code$ $IT.message$
     $END$
     */
    export async function $node.technicalName$(
      $LOOP node.input.fields$
        $IT.technicalName$: $CALL generateTsType(IT.type)$,
      $END$
    ): Promise<$node.technicalName$Response> {
      const response = await fetch('/api/v1/$node.serviceName$/$node.technicalName$', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          $LOOP node.input.fields$
            $IT.technicalName$,
          $END$
        }),
      });

      if (!response.ok) {
        const error: ErrorResponse = await response.json();
        throw new ApiError(error.code, error.message, error.http_status);
      }

      return response.json();
    }

    // 型定義
    export interface $node.technicalName$Response {
      $LOOP node.output.fields$
        /** $IT.description$ */
        $IT.technicalName$: $CALL generateTsType(IT.type)$;
      $END$
    }
  }
}
```

### 5.5 フロー図生成テンプレート例

```kotlin
// Generator Mapping Rule
map UseCase -> MermaidFlowDiagram {
  filename = "generated/docs/flows/" + node.technicalName + "-flow.md"

  template {
    # $node.name$ フロー図

    ## フローチャート

    ```mermaid
    flowchart TD
        Start([開始]) --> Input[入力値受信]

        $LOOP node.flow.steps WITH INDEX$
          $CALL generateFlowStep(IT, INDEX)$
        $END$

        $CALL generateErrorPaths()$

        End([終了])
    ```
  }
}

macro generateFlowStep(step: FlowStep, index: int) {
  Step$index$["$step.title$"]

  $IF index > 0$
    Step$index - 1$ --> Step$index$
  $ELSE$
    Input --> Step$index$
  $END$

  $LOOP step.statements$
    $IF IT is IfStatement && IT contains ErrorReturnStatement$
      Step$index$ -->|エラー| Error$IT.errorCode$["$IT.errorMessage$"]
      Error$IT.errorCode$ --> End
    $END$
  $END$
}
```

---

## STEP 6: サンプル作成・検証

### 6.1 サンプルユースケース作成

`sandbox/auth-service/user-registration.usecase` を作成:

```kotlin
ユースケース: ユーザー登録 {
  ID: "UC_AUTH_001"
  概要: "新規ユーザーを登録し、メール認証リンクを送信する"
  優先度: 高

  入力 {
    メールアドレス: 文字列 "ユーザーのメールアドレス"
    パスワード: 文字列 "ログイン用パスワード（8文字以上）"
    役割: 列挙型 {顧客, ショップオーナー, 管理者}
  }

  出力 {
    ユーザーID: UUID "作成されたユーザーの一意識別子"
    認証トークン: 文字列 "JWT認証トークン"
  }

  処理フロー {
    ステップ1: "入力値の検証" {
      もし メールアドレスが正しい形式でない場合:
        → エラー「メールアドレスの形式が不正です」を返す
    }

    ステップ2: "重複チェック" {
      既存ユーザー = データベースから メールアドレスで検索する

      もし 既存ユーザーが見つかった場合:
        → エラー「このメールアドレスは既に登録されています」を返す
    }
  }

  エラー定義 {
    メールアドレス形式不正 {
      コード: "AUTH_E001"
      HTTPステータス: 400
      メッセージ: "メールアドレスの形式が不正です"
    }

    メールアドレス重複 {
      コード: "AUTH_E003"
      HTTPステータス: 409
      メッセージ: "このメールアドレスは既に登録されています"
    }
  }
}
```

### 6.2 生成実行

1. MPSで `Build > Make Project`
2. 生成先を確認: `generated/auth/`
3. 生成されたファイル:
   - `generated/auth/usecase/user_registration.go`
   - `generated/proto/auth/v1/auth.proto`
   - `generated/frontend/src/api/auth.ts`
   - `generated/docs/flows/user-registration-flow.md`

### 6.3 生成コードの検証

```bash
# Goコードのコンパイル確認
cd generated/auth
go build ./...

# Protobufコードの生成確認
protoc --go_out=. --go-grpc_out=. proto/auth/v1/auth.proto

# TypeScriptの型チェック
cd generated/frontend
npm run type-check
```

---

## トラブルシューティング

### Q1. Generatorが実行されない

**原因**: Generatorのマッピングルールが正しく設定されていない

**対処**:
1. `Build > Clean`
2. `Build > Rebuild Language`
3. `Build > Make Project`

### Q2. 生成コードが期待と異なる

**原因**: テンプレートのロジックが間違っている

**対処**:
1. Generatorのデバッグモードを有効化
2. トレースログを確認
3. テンプレートを修正

### Q3. 日本語が文字化けする

**原因**: ファイルエンコーディングがUTF-8でない

**対処**:
1. MPSのファイルエンコーディング設定を確認
2. `File > Settings > Editor > File Encodings`
3. すべてを`UTF-8`に設定

---

## まとめ

このガイドに従うことで、日本語振る舞いDSLを実装し、以下を自動生成できるようになります:

- ✅ Go実装コード
- ✅ Protobuf定義
- ✅ TypeScript/Reactコード
- ✅ フロー図
- ✅ データ設計
- ✅ お客さん向け仕様書

**推定実装時間**: 1-2日（経験者の場合）

次のステップ: 実際にMPSプロジェクトを作成し、サンプルを試してください。
