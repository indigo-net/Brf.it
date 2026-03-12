---
title: TOML
---

# TOML サポート

[English](../../languages/toml.md) | [한국어](../../ko/languages/toml.md) | [日本語](toml.md) | [हिन्दी](../../hi/languages/toml.md) | [Deutsch](../../de/languages/toml.md)

## 対応拡張子

- `.toml`

## 文法

- [tree-sitter-toml](https://github.com/tree-sitter-grammars/tree-sitter-toml) v0.7.0 by tree-sitter-grammars

## 抽出対象

| 要素 | 種類 | XMLタグ | 例 |
|------|------|---------|-----|
| テーブル | `namespace` | `<type>` | `[package]` |
| テーブル配列 | `namespace` | `<type>` | `[[bin]]` |
| キー・バリューペア | `variable` | `<variable>` | `name = "myapp"` |

## 例

### 入力

```toml
# Project configuration
name = "myapp"
version = "1.0.0"

[package]
authors = ["Alice"]
edition = "2024"

[[bin]]
name = "cli"
path = "src/main.rs"
```

### 出力 (XML)

```xml
<file path="config.toml" language="toml">
  <variable>name = "myapp"</variable>
  <variable>version = "1.0.0"</variable>
  <type>[package]</type>
  <type>[[bin]]</type>
</file>
```

## 注意事項

### 本文の除去

`--include-body`フラグを使用しない場合：

- テーブルセクション（`[table]`）は本文が除去され、ヘッダーのみが表示されます
- テーブル配列（`[[table]]`）は本文が除去され、ヘッダーのみが表示されます
- トップレベルのキー・バリューペアは値が保持されます

`--include-private`を使用して非公開/unexportedシンボルを含める。

### コメント

- 単一行コメント（`# コメント`）はドキュメントとして抽出されます

### インポート

- TOMLにはインポートシステムがありません。`--include-imports`は効果がありません

### 制限事項

- インラインテーブルとインライン配列は個別に抽出されません
- ドットキー（例：`physical.color = "orange"`）は単一のペアとしてキャプチャされます
