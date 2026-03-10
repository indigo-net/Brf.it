---
title: SQL
---

# SQL 지원

[English](../../languages/sql.md) | [한국어](sql.md) | [日本語](../../ja/languages/sql.md) | [हिन्दी](../../hi/languages/sql.md) | [Deutsch](../../de/languages/sql.md)

## 지원 확장자

- `.sql`

## 문법

- [tree-sitter-sql](https://github.com/DerekStride/tree-sitter-sql) v0.3.11 by DerekStride

## 추출 대상

| 요소 | 종류 | XML 태그 | 예시 |
|------|------|----------|------|
| 테이블 | `struct` | `<type>` | `CREATE TABLE users (...)` |
| 함수 | `function` | `<function>` | `CREATE FUNCTION get_user(id INT) RETURNS TEXT` |
| 뷰 | `type` | `<type>` | `CREATE VIEW active_users` |
| 구체화된 뷰 | `type` | `<type>` | `CREATE MATERIALIZED VIEW stats` |
| 인덱스 | `variable` | `<variable>` | `CREATE INDEX idx_name ON users (name)` |
| 트리거 | `function` | `<function>` | `CREATE TRIGGER audit_trigger ...` |
| 타입/열거형 | `type` | `<type>` | `CREATE TYPE mood AS ENUM (...)` |
| 스키마 | `namespace` | `<type>` | `CREATE SCHEMA analytics` |
| 시퀀스 | `variable` | `<variable>` | `CREATE SEQUENCE user_id_seq ...` |
| 테이블 변경 | `type` | `<type>` | `ALTER TABLE users ADD COLUMN ...` |

## 예시

### 입력

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

### 출력 (XML)

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

## 참고사항

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 테이블: 컬럼 정의(스키마)는 출력에 보존됩니다
- 함수/프로시저: 본문이 제거되고, 반환 타입과 언어가 보존됩니다
- 뷰: `AS SELECT...` 쿼리가 제거되고, 선언부만 유지됩니다
- 구체화된 뷰: 뷰와 동일하게 쿼리가 제거됩니다

### 주석

- 단일 행 주석(`-- 주석`)은 문서로 추출됩니다
- 여러 행 주석(`/* 주석 */`)은 문서로 추출됩니다

### 스키마 한정 이름

- `schema.table`과 같은 스키마 한정 이름이 지원됩니다
- 예시: `CREATE TABLE analytics.events (...)`가 올바르게 추출됩니다

### 제한사항

- `CREATE PROCEDURE`는 문법에서 지원되지 않습니다 (tree-sitter-sql v0.3.11)
- DDL 문만 추출됩니다; DML 문(`INSERT`, `UPDATE`, `DELETE`, `SELECT`)은 무시됩니다
