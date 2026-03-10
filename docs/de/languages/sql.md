---
title: SQL
---

# SQL-Unterstützung

[English](../../languages/sql.md) | [한국어](../../ko/languages/sql.md) | [日本語](../../ja/languages/sql.md) | [हिन्दी](../../hi/languages/sql.md) | [Deutsch](sql.md)

## Unterstützte Erweiterungen

- `.sql`

## Grammatik

- [tree-sitter-sql](https://github.com/DerekStride/tree-sitter-sql) v0.3.11 von DerekStride

## Extraktionsziele

| Element | Art | XML-Tag | Beispiel |
|---------|-----|---------|----------|
| Tabelle | `struct` | `<type>` | `CREATE TABLE users (...)` |
| Funktion | `function` | `<function>` | `CREATE FUNCTION get_user(id INT) RETURNS TEXT` |
| View | `type` | `<type>` | `CREATE VIEW active_users` |
| Materialisierter View | `type` | `<type>` | `CREATE MATERIALIZED VIEW stats` |
| Index | `variable` | `<variable>` | `CREATE INDEX idx_name ON users (name)` |
| Trigger | `function` | `<function>` | `CREATE TRIGGER audit_trigger ...` |
| Typ/Enum | `type` | `<type>` | `CREATE TYPE mood AS ENUM (...)` |
| Schema | `namespace` | `<type>` | `CREATE SCHEMA analytics` |
| Sequenz | `variable` | `<variable>` | `CREATE SEQUENCE user_id_seq ...` |
| Tabellenänderung | `type` | `<type>` | `ALTER TABLE users ADD COLUMN ...` |

## Beispiel

### Eingabe

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

### Ausgabe (XML)

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

## Hinweise

### Entfernung des Funktionskörpers

Wenn das `--include-body`-Flag nicht verwendet wird:

- Tabellen: Spaltendefinitionen (Schema) bleiben in der Ausgabe erhalten
- Funktionen/Prozeduren: Der Körper wird entfernt, Rückgabetyp und Sprache bleiben erhalten
- Views: Die `AS SELECT...`-Abfrage wird entfernt, nur die Deklaration bleibt erhalten
- Materialisierte Views: wie bei Views, die Abfrage wird entfernt

### Kommentare

- Einzeilige Kommentare (`-- Kommentar`) werden als Dokumentation extrahiert
- Mehrzeilige Kommentare (`/* Kommentar */`) werden als Dokumentation extrahiert

### Schema-qualifizierte Namen

- Schema-qualifizierte Namen wie `schema.table` werden unterstützt
- Beispiel: `CREATE TABLE analytics.events (...)` wird korrekt extrahiert

### Einschränkungen

- `CREATE PROCEDURE` wird von der Grammatik nicht unterstützt (tree-sitter-sql v0.3.11)
- Nur DDL-Anweisungen werden extrahiert; DML-Anweisungen (`INSERT`, `UPDATE`, `DELETE`, `SELECT`) werden ignoriert
