#!/bin/bash

set -e

echo "Fixing proto type errors..."

# すべてのprotoファイルを処理
find proto/ -name "*.proto" | while read -r file; do
    echo "Processing $file..."
    
    # date型をstringに置換（コメント付き）
    sed -i.bak 's/\bdate\s\+\([a-zA-Z_][a-zA-Z0-9_]*\)/string \1  \/\/ YYYY-MM-DD format/g' "$file"
    
    # text型をstringに置換
    sed -i.bak 's/\btext\s\+\([a-zA-Z_][a-zA-Z0-9_]*\)/string \1/g' "$file"
    
    # json型をstringに置換（コメント付き）
    sed -i.bak 's/\bjson\s\+\([a-zA-Z_][a-zA-Z0-9_]*\)/string \1  \/\/ JSON format/g' "$file"
    
    # バックアップファイルを削除
    rm -f "$file.bak"
done

echo "Type fixes completed!"
