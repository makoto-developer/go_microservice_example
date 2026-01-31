# MPS 日本語振る舞いDSL - クイックスタート

**目的**: 最速でMPS実装を開始する
**所要時間**: 30分で動作するプロトタイプを作成

---

## 🚀 5ステップでプロトタイプ作成

### STEP 1: MPSインストール（10分）

```bash
# 1. JetBrains MPSをダウンロード
# https://www.jetbrains.com/mps/download/

# 2. インストール
# macOS: DMGをマウント、アプリケーションフォルダにドラッグ
# Windows: インストーラーを実行
# Linux: tar.gzを展開、bin/mps.shを実行

# 3. 起動確認
# JetBrains MPSを起動
# "Welcome to JetBrains MPS"が表示されればOK
```

### STEP 2: 最小限のプロジェクト作成（5分）

1. **新規プロジェクト作成**
   ```
   File > New > Project...
   → Select "New Language"
   → Language Name: "japanese-behavior-dsl"
   → Namespace: "com.example.behaviorDsl"
   → Location: /path/to/mps-workspace/languages/japanese-behavior-dsl
   → Create
   ```

2. **プロジェクト構造確認**
   ```
   japanese-behavior-dsl/
   ├── languages/
   │   └── japanese-behavior-dsl/
   │       └── (models will be created here)
   ├── solutions/
   │   └── sandbox/
   └── build.gradle (auto-generated)
   ```

### STEP 3: 最小限のStructure定義（5分）

**目標**: 1つのコンセプトだけ作成

1. `languages/japanese-behavior-dsl`を右クリック
2. `New > Model > Structure`
3. Model名: `structure`

4. Structure作成:
   ```
   右クリック > New > Concept
   Name: UseCase

   Add Property:
     - name: string
     - description: string

   Concept Properties:
     ✓ Root Concept
   ```

### STEP 4: 最小限のEditor定義（5分）

1. `languages/japanese-behavior-dsl`を右クリック
2. `New > Model > Editor`
3. Model名: `editor`

4. Editor作成:
   ```
   右クリック > New > Editor for Concept > UseCase

   EditorCell構成:
   [Collection: Vertical]
     [Constant: "ユースケース:"]
     [Property: name] (editable)
     [Property: description] (editable)
   ```

### STEP 5: Sandboxでテスト（5分）

1. `solutions/sandbox`を右クリック
2. `New > Model`
3. Model名: `examples`

4. サンプル作成:
   ```
   右クリック > New > Root Node > UseCase

   名前: ユーザー登録
   説明: 新規ユーザーを登録する
   ```

5. **確認**:
   - エディタで日本語が表示されていればOK！
   - これが最小限の動作するDSLです

---

## 🎯 次のステップ: Generator追加

### STEP 6: 最小限のGenerator（10分）

1. **Generatorモデル作成**
   ```
   languages/japanese-behavior-dsl を右クリック
   New > Model > Generator
   Model名: generator
   ```

2. **簡単なテンプレート作成**
   ```
   右クリック > New > Template Fragment
   Name: UseCaseToText

   Template:
   ユースケース: $node.name$
   説明: $node.description$
   ```

3. **Mapping Configuration**
   ```
   右クリック > New > Mapping Configuration

   Add Root Mapping Rule:
   Input: UseCase
   Template: UseCaseToText
   Output Path: generated/usecase.txt
   ```

4. **生成実行**
   ```
   Build > Make Project

   確認: generated/usecase.txt が作成される
   内容:
   ユースケース: ユーザー登録
   説明: 新規ユーザーを登録する
   ```

---

## ✅ 動作確認

以下が完了していればOK:

- [ ] MPSが起動する
- [ ] japanese-behavior-dslプロジェクトが作成された
- [ ] UseCaseコンセプトが定義された
- [ ] Sandboxでサンプルが作成できた
- [ ] 生成されたファイルが確認できた

---

## 🔧 よくある問題と解決

### Q1: MPSが起動しない

**原因**: Java環境が不足

**対処**:
```bash
# JDKのインストール確認
java -version

# Java 11以降が必要
# インストールされていない場合:
# https://adoptium.net/ から Java 17をダウンロード
```

### Q2: プロジェクトが作成できない

**原因**: ディレクトリの権限がない

**対処**:
```bash
# ディレクトリに書き込み権限があるか確認
ls -ld /path/to/mps-workspace/languages

# 権限がない場合
chmod 755 /path/to/mps-workspace/languages
```

### Q3: 生成されたファイルが見つからない

**原因**: 出力パスの設定が間違っている

**対処**:
1. `Build > Show Generation Plan`でログを確認
2. Mapping Configurationの出力パスを確認
3. 絶対パスで指定してみる

---

## 📚 次に読むドキュメント

### 基本を学んだ後

1. **詳細な実装**: `mps-step-by-step-guide.md`
   - 完全なStructure定義
   - 複雑なEditor定義
   - 実用的なGenerator実装

2. **Generator詳細**: `mps-generator-advanced.md`（次に作成）
   - Go生成の詳細
   - Protobuf生成の詳細
   - TypeScript生成の詳細

3. **デバッグ方法**: `mps-debugging.md`（次に作成）
   - Generatorのデバッグ
   - パフォーマンス最適化

---

## 🎉 成功の定義

このクイックスタートを完了すると:

- ✅ MPSの基本操作ができる
- ✅ Structure/Editor/Generatorの関係が理解できる
- ✅ 簡単なDSLが作成できる
- ✅ コード生成の仕組みが分かる

次は`mps-step-by-step-guide.md`で本格的な実装に進みましょう！

---

## 🆘 サポート

困ったら:
1. MPSドキュメント: https://www.jetbrains.com/help/mps/
2. MPSフォーラム: https://mps-support.jetbrains.com/
3. プロジェクトIssue: このプロジェクトのGitHub Issues
