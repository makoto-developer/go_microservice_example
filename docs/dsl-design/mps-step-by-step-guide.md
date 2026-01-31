# MPS 日本語振る舞いDSL - ステップバイステップ実装ガイド

**対象**: MPSで本格実装を行う開発者
**所要時間**: 1-2日
**前提**: JetBrains MPS 2023.3以降がインストール済み

---

## 📋 実装チェックリスト

- [ ] Phase 1: MPSプロジェクト作成（30分）
- [ ] Phase 2: Structure定義（2時間）
- [ ] Phase 3: Editor定義（1時間）
- [ ] Phase 4: Constraints定義（30分）
- [ ] Phase 5: Generator実装 - Go（2時間）
- [ ] Phase 6: Generator実装 - Protobuf（1時間）
- [ ] Phase 7: Generator実装 - TypeScript（1時間）
- [ ] Phase 8: Generator実装 - ドキュメント（30分）
- [ ] Phase 9: サンプル作成・検証（1時間）

---

## Phase 1: MPSプロジェクト作成

### 1.1 JetBrains MPSの起動

1. JetBrains MPSを起動
2. `File > New > Project...`
3. `New Language` を選択

### 1.2 プロジェクト設定

```
Language Name: japanese-behavior-dsl
Namespace: com.example.behaviorDsl
Location: /path/to/mps-workspace/languages/japanese-behavior-dsl

Create directory structure:
✓ Language
✓ Generator
✓ Sandbox
```

### 1.3 プロジェクト構造の確認

```
japanese-behavior-dsl/
├── japanese-behavior-dsl.msd       # 言語記述子
├── languages/
│   └── japanese-behavior-dsl/
│       ├── structure.mps           # 構文構造定義
│       ├── editor.mps              # エディタ定義
│       ├── constraints.mps         # 制約定義
│       └── typesystem.mps          # 型システム（オプション）
├── generator/
│   └── template/
│       ├── go-generator.mps        # Go生成器
│       ├── proto-generator.mps     # Protobuf生成器
│       ├── ts-generator.mps        # TypeScript生成器
│       └── docs-generator.mps      # ドキュメント生成器
└── sandbox/
    └── examples/
        └── user-registration.usecase
```

---

## Phase 2: Structure定義（構文構造）

### 2.1 基本コンセプトの作成

MPSの`Structure`モデルを開き、以下のコンセプトを作成します。

#### 2.1.1 ルートコンセプト: UseCase

`Structure > New Concept > UseCase`

```
Concept: UseCase
Extends: BaseConcept
Implements: INamedConcept

Properties:
  - name: string                    // ユースケース名（例: ユーザー登録）
  - id: string                      // ID（例: UC_AUTH_001）
  - description: string             // 概要
  - owner: string                   // 担当者
  - priority: PriorityLevel         // 優先度（Enum）

Children:
  - preconditions: PreconditionBlock [0..1]
  - input: InputBlock [0..1]
  - output: OutputBlock [0..1]
  - flow: ProcessFlowBlock [0..1]
  - errors: ErrorDefinitionBlock [0..1]
  - testCases: TestCaseBlock [0..1]
  - businessRules: BusinessRuleBlock [0..1]

Root Concept: true
```

**MPSでの操作**:
1. `Structure`タブを開く
2. 右クリック > `New Concept`
3. Name: `UseCase`
4. Properties: `Add Property` で上記プロパティを追加
5. Children: `Add Child Link` で上記子要素を追加
6. `Concept Properties` > `Root Concept` をチェック

#### 2.1.2 Enumの定義: PriorityLevel

```
Enumeration: PriorityLevel
Members:
  - 高 (HIGH)
  - 中 (MEDIUM)
  - 低 (LOW)
```

**MPSでの操作**:
1. `Structure`タブで右クリック > `New Enumeration`
2. Name: `PriorityLevel`
3. `Add Member` で各メンバーを追加

#### 2.1.3 入力ブロック: InputBlock

```
Concept: InputBlock
Children:
  - fields: FieldDefinition [1..*]
```

#### 2.1.4 フィールド定義: FieldDefinition

```
Concept: FieldDefinition
Implements: INamedConcept

Properties:
  - name: string                    // 日本語名（例: メールアドレス）
  - technicalName: string           // 技術名（例: email）
  - description: string             // 説明

Children:
  - type: TypeReference [1..1]
```

#### 2.1.5 型参照: TypeReference

```
Concept: TypeReference
Properties:
  - typeName: PrimitiveType         // Enum

Derived Properties:
  - goType: string (computed)
  - protoType: string (computed)
  - tsType: string (computed)

Enumeration: PrimitiveType
Members:
  - 文字列 (STRING)
  - 整数 (INT)
  - 真偽値 (BOOLEAN)
  - UUID (UUID)
  - 日時 (TIMESTAMP)
  - 小数 (DECIMAL)
```

**Derived Propertyの実装**:

```java
// goType derived property
property string goType {
  get {
    if (this.typeName == PrimitiveType.STRING) {
      return "string";
    } else if (this.typeName == PrimitiveType.INT) {
      return "int64";
    } else if (this.typeName == PrimitiveType.UUID) {
      return "uuid.UUID";
    } else if (this.typeName == PrimitiveType.TIMESTAMP) {
      return "time.Time";
    } else if (this.typeName == PrimitiveType.BOOLEAN) {
      return "bool";
    } else if (this.typeName == PrimitiveType.DECIMAL) {
      return "decimal.Decimal";
    }
    return "interface{}";
  }
}
```

#### 2.1.6 処理フローブロック: ProcessFlowBlock

```
Concept: ProcessFlowBlock
Children:
  - steps: FlowStep [1..*]
```

#### 2.1.7 フローステップ: FlowStep

```
Concept: FlowStep
Properties:
  - stepNumber: int
  - title: string
  - description: string

Children:
  - statements: Statement [0..*]
```

#### 2.1.8 ステートメント（基底クラス）

```
Abstract Concept: Statement
```

#### 2.1.9 条件分岐ステートメント: IfStatement

```
Concept: IfStatement
Extends: Statement

Properties:
  - conditionText: string           // 日本語条件（例: メールアドレスが正しい形式でない）

Children:
  - condition: Expression [1..1]    // 実際の条件式
  - thenStatements: Statement [1..*]
  - elseStatements: Statement [0..*]
```

#### 2.1.10 代入ステートメント: AssignmentStatement

```
Concept: AssignmentStatement
Extends: Statement

Properties:
  - variableName: string
  - description: string

Children:
  - value: Expression [1..1]
```

#### 2.1.11 データベース検索: DatabaseQueryStatement

```
Concept: DatabaseQueryStatement
Extends: Statement

Properties:
  - targetEntity: string            // 対象エンティティ（例: ユーザー）
  - searchField: string             // 検索フィールド（例: メールアドレス）
  - resultVariable: string          // 結果変数名

Children:
  - searchValue: Expression [1..1]
```

#### 2.1.12 エラー返却: ErrorReturnStatement

```
Concept: ErrorReturnStatement
Extends: Statement

Properties:
  - errorMessage: string

References:
  - errorDef: ErrorDefinition [1..1]
```

#### 2.1.13 エラー定義: ErrorDefinition

```
Concept: ErrorDefinition
Implements: INamedConcept

Properties:
  - name: string                    // エラー名（例: メールアドレス重複）
  - technicalName: string           // 技術名（例: EmailDuplicate）
  - code: string                    // エラーコード（例: AUTH_E003）
  - httpStatus: int                 // HTTPステータス
  - message: string                 // メッセージ
```

### 2.2 Structure定義の完成形

すべてのコンセプトを作成後、以下のようなツリー構造になります:

```
Structure
├── UseCase (root)
│   ├── PreconditionBlock
│   ├── InputBlock
│   │   └── FieldDefinition
│   │       └── TypeReference
│   ├── OutputBlock
│   │   └── FieldDefinition
│   ├── ProcessFlowBlock
│   │   └── FlowStep
│   │       └── Statement (abstract)
│   │           ├── IfStatement
│   │           ├── AssignmentStatement
│   │           ├── DatabaseQueryStatement
│   │           ├── ErrorReturnStatement
│   │           └── ...
│   ├── ErrorDefinitionBlock
│   │   └── ErrorDefinition
│   └── TestCaseBlock
└── Enums
    ├── PriorityLevel
    └── PrimitiveType
```

---

## Phase 3: Editor定義（見た目）

### 3.1 UseCaseエディタの作成

`Editor`モデルを開き、`UseCase`のエディタを作成します。

#### 3.1.1 基本レイアウト

```
EditorCell for UseCase:

[Collection: Vertical]
  [Constant: "ユースケース:"] [Property: name] (editable, bold)

  [Collection: Horizontal]
    [Constant: "ID:"] [Property: id] (editable)

  [Collection: Horizontal]
    [Constant: "概要:"] [Property: description] (editable, multiline)

  [Collection: Horizontal]
    [Constant: "優先度:"] [Property: priority] (editable, enum dropdown)

  [Empty Line]

  [Child: preconditions]
  [Child: input]
  [Child: output]
  [Child: flow]
  [Child: errors]
  [Child: testCases]
```

**MPSでの操作**:
1. `Editor`タブを開く
2. 右クリック > `New Editor for Concept > UseCase`
3. `Editor Component`を追加:
   - `Constant Cell`: 固定テキスト（"ユースケース:"等）
   - `Property Cell`: プロパティ参照（name, id等）
   - `Child Cell`: 子要素（input, output等）
   - `Collection Cell`: レイアウトコンテナ（Vertical, Horizontal）

#### 3.1.2 スタイル設定

```
Style for "ユースケース:" constant:
  - font-style: bold
  - text-color: blue

Style for name property:
  - font-size: 14pt
  - font-style: bold

Style for description property:
  - multiline: true
  - min-width: 400px
```

### 3.2 InputBlockエディタ

```
EditorCell for InputBlock:

[Collection: Vertical]
  [Constant: "入力"] [Constant: "{"] (same line)

  [Indent]
    [Collection: Vertical]
      [Child Collection: fields]

  [Constant: "}"]
```

### 3.3 FieldDefinitionエディタ

```
EditorCell for FieldDefinition:

[Collection: Horizontal]
  [Property: name] [Constant: ":"] [Child: type] [Property: description] (quoted)
```

例: `メールアドレス: 文字列 "ユーザーのメールアドレス"`

### 3.4 FlowStepエディタ

```
EditorCell for FlowStep:

[Collection: Vertical]
  [Collection: Horizontal]
    [Constant: "ステップ"][Property: stepNumber][Constant: ":"] [Property: title] (quoted, bold)

  [Indent]
    [Collection: Horizontal]
      [Constant: "説明:"] [Property: description]

    [Empty Line]

    [Child Collection: statements]
```

### 3.5 IfStatementエディタ

```
EditorCell for IfStatement:

[Collection: Vertical]
  [Collection: Horizontal]
    [Constant: "もし"] [Property: conditionText] [Constant: "場合:"] (blue, bold)

  [Indent]
    [Child Collection: thenStatements]

  [Optional: if elseStatements not empty]
    [Collection: Horizontal]
      [Constant: "そうでなければ:"] (blue, bold)

    [Indent]
      [Child Collection: elseStatements]
```

### 3.6 ErrorReturnStatementエディタ

```
EditorCell for ErrorReturnStatement:

[Collection: Horizontal]
  [Constant: "→ エラー"] [Property: errorMessage] (quoted, red) [Constant: "を返す"]
  [Constant: "→ 処理を終了する"]
```

---

## Phase 4: Constraints定義（バリデーション）

### 4.1 UseCase名の制約

`Constraints`モデルを開き、以下を定義:

```kotlin
constraint can be child {
  // UseCaseは最上位のみ（他のコンセプトの子になれない）
  applicable for concept = UseCase

  return node.parent == null;

  error message = "ユースケースは最上位レベルにのみ配置できます";
}

constraint valid name {
  applicable for concept = UseCase
  applicable for property = name

  return node.name != null &&
         node.name.length >= 2 &&
         node.name.length <= 50;

  error message = "ユースケース名は2～50文字である必要があります";
}
```

### 4.2 ErrorDefinitionの制約

```kotlin
constraint valid error code {
  applicable for concept = ErrorDefinition
  applicable for property = code

  return node.code.matches("^[A-Z]+_E[0-9]{3}$");

  error message = "エラーコードは「サービス名_E001」形式である必要があります（例: AUTH_E001）";
}

constraint valid http status {
  applicable for concept = ErrorDefinition
  applicable for property = httpStatus

  int[] validStatuses = {400, 401, 403, 404, 409, 500, 502, 503};
  return Arrays.asList(validStatuses).contains(node.httpStatus);

  error message = "HTTPステータスは有効な値（400, 401, 403, 404, 409, 500等）である必要があります";
}
```

---

## Phase 5: Generator実装 - Go

### 5.1 Generatorの作成

`Generator`モデルを開き、以下を作成:

#### 5.1.1 Mapping Configuration

```
Mapping Configuration: GoGenerator

Input: UseCase
Output: TextGen (Go file)
```

#### 5.1.2 Generator Template

`Generator Template: UseCaseToGo`

```go
// Template Language: Java + TextGen

$TEMPLATE UseCaseToGo(node: UseCase) -> TextGen$

// Code generated by MPS japanese-behavior-dsl. DO NOT EDIT.

package usecase

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    $LOOP node.getRequiredImports()$
        "$IT$"
    $END$
)

// $node.name$ ユースケース
// $node.description$
type $node.getTechnicalName()$Usecase interface {
    Execute(ctx context.Context, input $node.getTechnicalName()$Input) ($node.getTechnicalName()$Output, error)
}

// 入力
type $node.getTechnicalName()$Input struct {
    $LOOP node.input.fields$
        // $IT.description$
        $IT.getTechnicalName()$ $IT.type.goType$
    $END$
}

// 出力
type $node.getTechnicalName()$Output struct {
    $LOOP node.output.fields$
        // $IT.description$
        $IT.getTechnicalName()$ $IT.type.goType$
    $END$
}

type $node.getTechnicalName()$UsecaseImpl struct {
    $LOOP node.getDependencies()$
        $IT.name$ $IT.type$
    $END$
}

func New$node.getTechnicalName()$Usecase(
    $LOOP node.getDependencies()$
        $IT.name$ $IT.type$,
    $END$
) $node.getTechnicalName()$Usecase {
    return &$node.getTechnicalName()$UsecaseImpl{
        $LOOP node.getDependencies()$
            $IT.name$: $IT.name$,
        $END$
    }
}

func (u *$node.getTechnicalName()$UsecaseImpl) Execute(
    ctx context.Context,
    input $node.getTechnicalName()$Input,
) ($node.getTechnicalName()$Output, error) {

    $LOOP node.flow.steps$
        // $IT.title$
        // $IT.description$
        $CALL generateStep(IT)$
    $END$
}

// エラー定義
$LOOP node.errors.definitions$
    $CALL generateErrorType(IT)$
$END$

$ENDTEMPLATE$
```

#### 5.1.3 Macro定義

```java
$MACRO generateStep(step: FlowStep)$
    $LOOP step.statements$
        $CALL generateStatement(IT)$
    $END$
$ENDMACRO$

$MACRO generateStatement(stmt: Statement)$
    $IF stmt instanceof IfStatement$
        if $CALL generateCondition(stmt.condition)$ {
            $LOOP stmt.thenStatements$
                $CALL generateStatement(IT)$
            $END$
        }
        $IF stmt.elseStatements.isNotEmpty()$
            else {
                $LOOP stmt.elseStatements$
                    $CALL generateStatement(IT)$
                $END$
            }
        $END$

    $ELSEIF stmt instanceof AssignmentStatement$
        $stmt.variableName$ := $CALL generateExpression(stmt.value)$

    $ELSEIF stmt instanceof DatabaseQueryStatement$
        $stmt.resultVariable$, err := u.repo.FindBy$stmt.searchField$(ctx, $CALL generateExpression(stmt.searchValue)$)
        if err != nil {
            return $node.getTechnicalName()$Output{}, err
        }

    $ELSEIF stmt instanceof ErrorReturnStatement$
        return $node.getTechnicalName()$Output{}, &$stmt.errorDef.technicalName$Error{
            Code:       "$stmt.errorDef.code$",
            Message:    "$stmt.errorDef.message$",
            HTTPStatus: $stmt.errorDef.httpStatus$,
        }
    $END$
$ENDMACRO$

$MACRO generateErrorType(error: ErrorDefinition)$
// $error.name$
type $error.technicalName$Error struct {
    Code       string
    Message    string
    HTTPStatus int
}

func (e *$error.technicalName$Error) Error() string {
    return e.Message
}

$ENDMACRO$
```

#### 5.1.4 Behavior Methods（ヘルパーメソッド）

`UseCase`コンセプトに以下のBehavior Methodsを追加:

```java
// UseCaseのBehavior Methods

public string getTechnicalName() {
    // "ユーザー登録" -> "UserRegistration"
    return CaseConverter.toPascalCase(this.name);
}

public List<Dependency> getDependencies() {
    List<Dependency> deps = new ArrayList<>();

    // データベース検索があればRepositoryが必要
    if (hasDatab aseQuery()) {
        deps.add(new Dependency("repo", "domain.UserRepository"));
    }

    // JWT生成があればJWTServiceが必要
    if (hasJWTGeneration()) {
        deps.add(new Dependency("jwtService", "JWTService"));
    }

    // イベント発行があればEventPublisherが必要
    if (hasEventPublish()) {
        deps.add(new Dependency("eventPublisher", "EventPublisher"));
    }

    return deps;
}

private boolean hasDatabaseQuery() {
    for (FlowStep step : this.flow.steps) {
        for (Statement stmt : step.statements) {
            if (stmt instanceof DatabaseQueryStatement) {
                return true;
            }
        }
    }
    return false;
}
```

---

## Phase 6: Generator実装 - Protobuf

### 6.1 Protobuf Template

```protobuf
$TEMPLATE UseCaseToProto(node: UseCase) -> TextGen$

// Code generated by MPS japanese-behavior-dsl. DO NOT EDIT.

syntax = "proto3";

package $node.getServiceName()$.v1;

option go_package = "github.com/.../proto/$node.getServiceName()$/v1;$node.getServiceName()$v1";

// $node.name$ リクエスト
message $node.getTechnicalName()$Request {
    $LOOP node.input.fields WITH INDEX$
        // $IT.description$
        $IT.type.protoType$ $IT.technicalName$ = $INDEX + 1$;
    $END$
}

// $node.name$ レスポンス
message $node.getTechnicalName()$Response {
    $LOOP node.output.fields WITH INDEX$
        // $IT.description$
        $IT.type.protoType$ $IT.technicalName$ = $INDEX + 1$;
    $END$
}

// エラーレスポンス
message ErrorResponse {
    string code = 1;
    string message = 2;
    int32 http_status = 3;
}

service $node.getServiceName()$Service {
    // $node.name$
    // $node.description$
    rpc $node.getTechnicalName()$($node.getTechnicalName()$Request) returns ($node.getTechnicalName()$Response);
}

$ENDTEMPLATE$
```

---

## Phase 7: Generator実装 - TypeScript

### 7.1 TypeScript Template

```typescript
$TEMPLATE UseCaseToTypeScript(node: UseCase) -> TextGen$

// Code generated by MPS japanese-behavior-dsl. DO NOT EDIT.

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
 $LOOP node.errors.definitions$
 * @throws $IT.code$ $IT.message$
 $END$
 */
export async function $node.getTechnicalNameCamelCase()$(
    $LOOP node.input.fields$
        $IT.technicalName$: $IT.type.tsType$,
    $END$
): Promise<$node.getTechnicalName()$Response> {
    const response = await fetch('/api/v1/$node.getServiceName()$/$node.getTechnicalNameCamelCase()$', {
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

export interface $node.getTechnicalName()$Response {
    $LOOP node.output.fields$
        /** $IT.description$ */
        $IT.technicalName$: $IT.type.tsType$;
    $END$
}

$ENDTEMPLATE$
```

---

## Phase 8: サンプル作成・検証

### 8.1 Sandboxでサンプル作成

1. `Sandbox`ソリューションを開く
2. 新しいモデルを作成: `user-registration-example`
3. `UseCase`のインスタンスを作成
4. 日本語で仕様を記述

### 8.2 生成の実行

1. MPSメニュー: `Build > Make Project`
2. 生成先ディレクトリを確認: `generated/`
3. 生成されたファイルを確認:
   - `generated/auth/usecase/user_registration.go`
   - `generated/proto/auth/v1/auth.proto`
   - `generated/frontend/src/api/auth.ts`

### 8.3 生成コードの検証

```bash
# Goコードのコンパイル確認
cd generated/auth
go mod init github.com/.../auth
go mod tidy
go build ./...

# Protobufコードの生成
protoc --go_out=. --go-grpc_out=. proto/auth/v1/auth.proto

# TypeScriptの型チェック
cd generated/frontend
npm install
npm run type-check
```

---

## トラブルシューティング

### Q1: Generatorが実行されない

**原因**: Generatorのマッピングルールが正しく設定されていない

**対処**:
1. `Build > Clean`
2. `Build > Rebuild Language japanese-behavior-dsl`
3. `Build > Make Project`

### Q2: 生成されたコードがコンパイルエラー

**原因**: Templateのロジックが間違っている、または型変換が不正

**対処**:
1. Generatorのデバッグモードを有効化
2. `Build > Show Generation Plan`でトレースログを確認
3. Templateを修正して再生成

### Q3: 日本語が文字化けする

**原因**: ファイルエンコーディングがUTF-8でない

**対処**:
1. `File > Settings > Editor > File Encodings`
2. すべてのエンコーディングを`UTF-8`に設定
3. MPSを再起動

### Q4: Behavior Methodsが呼び出せない

**原因**: Behavior Methodsが正しく定義されていない

**対処**:
1. `Behavior`モデルで対象コンセプトを開く
2. `Add Method`で必要なメソッドを追加
3. メソッドの実装を確認

---

## 次のステップ

### 完了したら

- [ ] 生成されたGoコードがコンパイルできることを確認
- [ ] Protobufファイルが正しく生成されることを確認
- [ ] TypeScriptファイルの型が正しいことを確認
- [ ] 他のユースケース（ログイン、パスワードリセット等）でも試す

### さらに拡張する場合

- [ ] フロー図生成器の追加（Mermaid）
- [ ] ER図生成器の追加
- [ ] お客さん向け仕様書生成器の追加
- [ ] テストコード生成器の追加

---

## まとめ

このガイドに従うことで、日本語振る舞いDSLを完全に実装できます。

**重要なポイント**:
1. Structure定義が最も重要（ここが言語の骨格）
2. Editor定義で見た目を整える
3. Generatorで実際のコード生成を行う
4. Behavior Methodsでヘルパーロジックを実装

**所要時間**: 約1-2日

次のドキュメント: `mps-generator-advanced.md`（高度なGenerator実装）
