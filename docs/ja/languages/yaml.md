---
title: YAML
---

# YAML サポート

[English](../../languages/yaml.md) | [한국어](../../ko/languages/yaml.md) | [日本語](yaml.md) | [हिन्दी](../../hi/languages/yaml.md) | [Deutsch](../../de/languages/yaml.md)

## 対応拡張子

- `.yaml`
- `.yml`

## 文法

- [tree-sitter-yaml](https://github.com/tree-sitter-grammars/tree-sitter-yaml) v0.7.2 by tree-sitter-grammars

## 抽出対象

| 要素 | 種類 | XMLタグ | 例 |
|------|------|---------|-----|
| キー・バリューペア | `variable` | `<variable>` | `name: value` |

## 例

### 入力

```yaml
# Application configuration
name: myapp
version: 1.0.0

database:
  host: localhost
  port: 5432

features:
  - logging
  - metrics
```

### 出力 (XML)

```xml
<file path="config.yaml" language="yaml">
  <variable>name: myapp</variable>
  <variable>version: 1.0.0</variable>
  <variable>database:</variable>
  <variable>features:</variable>
</file>
```

## 注意事項

### 本文の除去

`--include-body`フラグを使用しない場合：

- コンテナキー（ネストされた値を持つマッピング）はネストされた内容が除去され、キーのみが表示されます
- スカラーキー・バリューペアは値が保持されます

`--include-private`を使用して非公開/unexportedシンボルを含める。

### コメント

- 単一行コメント（`# コメント`）はドキュメントとして抽出されます

### インポート

- YAMLにはインポートシステムがありません。`--include-imports`は効果がありません

### 制限事項

- 過度なノイズを防ぐため、トップレベルのキーのみがシグネチャとしてキャプチャされます
- アンカーとエイリアスは特別に処理されません
