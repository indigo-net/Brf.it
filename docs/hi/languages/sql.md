---
title: SQL
---

# SQL समर्थन

[English](../../languages/sql.md) | [한국어](../../ko/languages/sql.md) | [日本語](../../ja/languages/sql.md) | [हिन्दी](sql.md) | [Deutsch](../../de/languages/sql.md)

## समर्थित एक्सटेंशन

- `.sql`

## व्याकरण

- [tree-sitter-sql](https://github.com/DerekStride/tree-sitter-sql) v0.3.11 by DerekStride

## निष्कर्षण लक्ष्य

| तत्व | प्रकार | XML टैग | उदाहरण |
|------|--------|---------|--------|
| टेबल | `struct` | `<type>` | `CREATE TABLE users (...)` |
| फ़ंक्शन | `function` | `<function>` | `CREATE FUNCTION get_user(id INT) RETURNS TEXT` |
| व्यू | `type` | `<type>` | `CREATE VIEW active_users` |
| मटीरियलाइज़्ड व्यू | `type` | `<type>` | `CREATE MATERIALIZED VIEW stats` |
| इंडेक्स | `variable` | `<variable>` | `CREATE INDEX idx_name ON users (name)` |
| ट्रिगर | `function` | `<function>` | `CREATE TRIGGER audit_trigger ...` |
| टाइप/एनम | `type` | `<type>` | `CREATE TYPE mood AS ENUM (...)` |
| स्कीमा | `namespace` | `<type>` | `CREATE SCHEMA analytics` |
| सीक्वेंस | `variable` | `<variable>` | `CREATE SEQUENCE user_id_seq ...` |
| टेबल बदलाव | `type` | `<type>` | `ALTER TABLE users ADD COLUMN ...` |

## उदाहरण

### इनपुट

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

### आउटपुट (XML)

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

## नोट्स

### बॉडी हटाना

जब `--include-body` फ़्लैग का उपयोग नहीं किया जाता:

- टेबल: कॉलम परिभाषाएँ (स्कीमा) आउटपुट में संरक्षित रहती हैं
- फ़ंक्शन/प्रोसीजर: बॉडी हटा दी जाती है, रिटर्न टाइप और भाषा संरक्षित रहती है
- व्यू: `AS SELECT...` क्वेरी हटा दी जाती है, केवल घोषणा रखी जाती है
- मटीरियलाइज़्ड व्यू: व्यू के समान, क्वेरी हटा दी जाती है

### टिप्पणियाँ

- एकल पंक्ति टिप्पणियाँ (`-- टिप्पणी`) दस्तावेज़ के रूप में निष्कर्षित की जाती हैं
- बहु पंक्ति टिप्पणियाँ (`/* टिप्पणी */`) दस्तावेज़ के रूप में निष्कर्षित की जाती हैं

### स्कीमा-योग्य नाम

- `schema.table` जैसे स्कीमा-योग्य नाम समर्थित हैं
- उदाहरण: `CREATE TABLE analytics.events (...)` सही ढंग से निष्कर्षित होता है

### सीमाएँ

- `CREATE PROCEDURE` व्याकरण द्वारा समर्थित नहीं है (tree-sitter-sql v0.3.11)
- केवल DDL स्टेटमेंट निष्कर्षित किए जाते हैं; DML स्टेटमेंट (`INSERT`, `UPDATE`, `DELETE`, `SELECT`) को अनदेखा किया जाता है
