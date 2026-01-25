# Elasticsearch セットアップガイド

このドキュメントは、Elasticsearchのセットアップと日本語検索用Kuromojiプラグインのインストール手順を説明します。

---

## 基本情報

| 項目 | 値 |
|------|-----|
| **バージョン** | 8.11.0 |
| **ポート** | 22000 (HTTP API), 22001 (Transport) |
| **URL** | http://localhost:22000 |
| **認証** | 無効化済み（開発環境） |

---

## 起動確認

### コンテナ状態確認

```bash
docker ps --filter "name=elasticsearch"
```

期待される出力:
```
NAMES                              STATUS
go_microservice_elasticsearch_dev  Up X minutes (healthy)
```

### ヘルスチェック

```bash
curl http://localhost:22000/_cluster/health
```

期待される出力:
```json
{
  "cluster_name": "docker-cluster",
  "status": "green",
  "number_of_nodes": 1
}
```

---

## Kuromoji プラグインのインストール

日本語の形態素解析を行うために、Kuromojiプラグインをインストールします。

### インストール手順

#### 1. Elasticsearchコンテナに入る

```bash
docker exec -it go_microservice_elasticsearch_dev bash
```

#### 2. プラグインをインストール

```bash
elasticsearch-plugin install analysis-kuromoji
```

実行例:
```
-> Installing analysis-kuromoji
-> Downloading analysis-kuromoji from elastic
[=================================================] 100%  
-> Installed analysis-kuromoji
```

#### 3. コンテナを再起動

```bash
exit
docker restart go_microservice_elasticsearch_dev
```

#### 4. インストール確認

```bash
# プラグイン一覧を確認
docker exec -it go_microservice_elasticsearch_dev elasticsearch-plugin list
```

期待される出力:
```
analysis-kuromoji
```

---

## Kuromoji の使用例

### インデックス作成（日本語対応）

```bash
curl -X PUT "http://localhost:22000/products" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "analysis": {
      "analyzer": {
        "kuromoji_analyzer": {
          "type": "custom",
          "tokenizer": "kuromoji_tokenizer",
          "filter": ["kuromoji_baseform", "kuromoji_part_of_speech", "cjk_width", "stop", "lowercase"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "name": {
        "type": "text",
        "analyzer": "kuromoji_analyzer"
      },
      "description": {
        "type": "text",
        "analyzer": "kuromoji_analyzer"
      },
      "price": {
        "type": "integer"
      }
    }
  }
}
'
```

### ドキュメント追加

```bash
curl -X POST "http://localhost:22000/products/_doc" -H 'Content-Type: application/json' -d'
{
  "name": "オーガニックマンゴー",
  "description": "新鮮な有機栽培のマンゴーです。甘くて美味しい。",
  "price": 1500
}
'
```

### 日本語検索

```bash
curl -X GET "http://localhost:22000/products/_search" -H 'Content-Type: application/json' -d'
{
  "query": {
    "match": {
      "description": "美味しい"
    }
  }
}
'
```

---

## トラブルシューティング

### プラグインインストールエラー

**エラー**: `elasticsearch-plugin: command not found`

**原因**: コンテナ内にelasticsearch-pluginコマンドがない

**解決**:
```bash
# 正しいパスを使用
/usr/share/elasticsearch/bin/elasticsearch-plugin install analysis-kuromoji
```

### コンテナが起動しない

**エラー**: `max virtual memory areas vm.max_map_count [65530] is too low`

**解決** (macOS):
```bash
# Dockerリソース設定を変更
# Docker Desktop → Settings → Resources → Advanced
# Memory: 4GB以上
```

**解決** (Linux):
```bash
sudo sysctl -w vm.max_map_count=262144
echo "vm.max_map_count=262144" | sudo tee -a /etc/sysctl.conf
```

### プラグインが認識されない

**原因**: コンテナ再起動していない

**解決**:
```bash
docker restart go_microservice_elasticsearch_dev

# 再起動後、ヘルスチェック
curl http://localhost:22000/_cluster/health
```

---

## よく使うコマンド

### インデックス一覧

```bash
curl http://localhost:22000/_cat/indices?v
```

### インデックス削除

```bash
curl -X DELETE http://localhost:22000/products
```

### すべてのドキュメントを取得

```bash
curl http://localhost:22000/products/_search?pretty
```

### マッピング確認

```bash
curl http://localhost:22000/products/_mapping?pretty
```

---

## 本番環境への移行

開発環境では認証を無効化していますが、本番環境では以下を設定してください：

### セキュリティ設定

```yaml
# docker-compose.yml (本番環境)
elasticsearch:
  environment:
    - xpack.security.enabled=true
    - ELASTIC_PASSWORD=your_strong_password
```

### HTTPS設定

```yaml
elasticsearch:
  environment:
    - xpack.security.http.ssl.enabled=true
    - xpack.security.transport.ssl.enabled=true
```

---

## 参考リンク

- [Elasticsearch公式ドキュメント](https://www.elastic.co/guide/en/elasticsearch/reference/8.11/index.html)
- [Kuromoji公式ドキュメント](https://www.elastic.co/guide/en/elasticsearch/plugins/8.11/analysis-kuromoji.html)
- [日本語検索のベストプラクティス](https://www.elastic.co/guide/en/elasticsearch/reference/8.11/analysis-lang-analyzer.html#japanese-analyzer)

---

**最終更新**: 2026-01-25
**管理者**: Claude
