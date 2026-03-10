---
title: SQL
---

# SQL サポート

[English](../../languages/sql.md) | [한국어](../../ko/languages/sql.md) | [日本語](sql.md) | [हिन्दी](../../hi/languages/sql.md) | [Deutsch](../../de/languages/sql.md)

## 対応拡張子

- `.sql`

## 文法

- [tree-sitter-sql](https://github.com/DerekStride/tree-sitter-sql) v0.3.11 by DerekStride

## 抽出対象

| 要素 | 種類 | XMLタグ | 例 |
|------|------|---------|-----|
| テーブル | `struct` | `<type>` | `CREATE TABLE users (...)` |
| 関数 | `function` | `<function>` | `CREATE FUNCTION get_user(id INT) RETURNS TEXT` |
| ビュー | `type` | `<type>` | `CREATE VIEW active_users` |
| マテリアライズドビュー | `type` | `<type>` | `CREATE MATERIALIZED VIEW stats` |
| インデックス | `variable` | `<variable>` | `CREATE INDEX idx_name ON users (name)` |
| トリガー | `function` | `<function>` | `CREATE TRIGGER audit_trigger ...` |
| 型/列挙型 | `type` | `<type>` | `CREATE TYPE mood AS ENUM (...)` |
| スキーマ | `namespace` | `<type>` | `CREATE SCHEMA analytics` |
| シーケンス | `variable` | `<variable>` | `CREATE SEQUENCE user_id_seq ...` |
| テーブル変更 | `type` | `<type>` | `ALTER TABLE users ADD COLUMN ...` |

## 例

### 入力

```sql
-- User management schema
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE FUNCTION get_user(user_id INT) RETURNS TEXT AS $$
BEGIN
    RETURN 'hello';
END;
$$ LANGUAGE plpgsql;

CREATE VIEW active_users AS
SELECT * FROM users WHERE active = true;

CREATE INDEX idx_users_name ON users (name);
```

### 出力 (XML)

```xml
<file path="schema.sql" language="sql">
  <type>CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
)</type>
  <function>CREATE FUNCTION get_user(user_id INT) RETURNS TEXT LANGUAGE plpgsql</function>
  <type>CREATE VIEW active_users</type>
  <variable>CREATE INDEX idx_users_name ON users (name)</variable>
</file>
```

## 注意事項

### 本文の除去

`--include-body`フラグを使用しない場合：

- テーブル：カラム定義（スキーマ）は出力に保持されます
- 関数/プロシージャ：本文が除去され、戻り値の型と言語が保持されます
- ビュー：`AS SELECT...`クエリが除去され、宣言部のみが保持されます
- マテリアライズドビュー：ビューと同様にクエリが除去されます

### コメント

- 単一行コメント（`-- コメント`）はドキュメントとして抽出されます
- 複数行コメント（`/* コメント */`）はドキュメントとして抽出されます

### スキーマ修飾名

- `schema.table`のようなスキーマ修飾名がサポートされています
- 例：`CREATE TABLE analytics.events (...)`は正しく抽出されます

### 制限事項

- `CREATE PROCEDURE`は文法でサポートされていません（tree-sitter-sql v0.3.11）
- DDL文のみ抽出されます。DML文（`INSERT`、`UPDATE`、`DELETE`、`SELECT`）は無視されます
