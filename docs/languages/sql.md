---
layout: default
title: SQL
parent: Language Guides
nav_order: 18
---

# SQL Support

[English](sql.md) | [한국어](../ko/languages/sql.md) | [日本語](../ja/languages/sql.md) | [हिन्दी](../hi/languages/sql.md) | [Deutsch](../de/languages/sql.md)

## Supported Extensions

- `.sql`

## Grammar

- [tree-sitter-sql](https://github.com/DerekStride/tree-sitter-sql) v0.3.11 by DerekStride

## Extraction Targets

| Element | Kind | XML Tag | Example |
|---------|------|---------|---------|
| Table | `struct` | `<type>` | `CREATE TABLE users (...)` |
| Function | `function` | `<function>` | `CREATE FUNCTION get_user(id INT) RETURNS TEXT` |
| View | `type` | `<type>` | `CREATE VIEW active_users` |
| Materialized View | `type` | `<type>` | `CREATE MATERIALIZED VIEW stats` |
| Index | `variable` | `<variable>` | `CREATE INDEX idx_name ON users (name)` |
| Trigger | `function` | `<function>` | `CREATE TRIGGER audit_trigger ...` |
| Type/Enum | `type` | `<type>` | `CREATE TYPE mood AS ENUM (...)` |
| Schema | `namespace` | `<type>` | `CREATE SCHEMA analytics` |
| Sequence | `variable` | `<variable>` | `CREATE SEQUENCE user_id_seq ...` |
| Alter Table | `type` | `<type>` | `ALTER TABLE users ADD COLUMN ...` |

## Example

### Input

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

### Output (XML)

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

## Notes

### Body Removal

When `--include-body` flag is not used:

- Tables: column definitions (schema) are preserved in the output
- Functions/procedures: body is stripped, return type and language are preserved
- Views: the `AS SELECT...` query is stripped, keeping only the declaration
- Materialized views: same as views, query is stripped

### Comments

- Single-line comments (`-- comment`) are extracted as documentation
- Multi-line comments (`/* comment */`) are extracted as documentation

### Schema-Qualified Names

- Schema-qualified names such as `schema.table` are supported
- Example: `CREATE TABLE analytics.events (...)` is correctly extracted

### Limitations

- `CREATE PROCEDURE` is not supported by the grammar (tree-sitter-sql v0.3.11)
- Only DDL statements are extracted; DML statements (`INSERT`, `UPDATE`, `DELETE`, `SELECT`) are ignored
