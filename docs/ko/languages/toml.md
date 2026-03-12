---
title: TOML
---

# TOML 지원

[English](../../languages/toml.md) | [한국어](toml.md) | [日本語](../../ja/languages/toml.md) | [हिन्दी](../../hi/languages/toml.md) | [Deutsch](../../de/languages/toml.md)

## 지원 확장자

- `.toml`

## 문법

- [tree-sitter-toml](https://github.com/tree-sitter-grammars/tree-sitter-toml) v0.7.0 by tree-sitter-grammars

## 추출 대상

| 요소 | 종류 | XML 태그 | 예시 |
|------|------|----------|------|
| 테이블 | `namespace` | `<type>` | `[package]` |
| 테이블 배열 | `namespace` | `<type>` | `[[bin]]` |
| 키-값 쌍 | `variable` | `<variable>` | `name = "myapp"` |

## 예시

### 입력

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

### 출력 (XML)

```xml
<file path="config.toml" language="toml">
  <variable>name = "myapp"</variable>
  <variable>version = "1.0.0"</variable>
  <type>[package]</type>
  <type>[[bin]]</type>
</file>
```

## 참고사항

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 테이블 섹션(`[table]`)은 본문이 제거되고, 헤더만 표시됩니다
- 테이블 배열(`[[table]]`)은 본문이 제거되고, 헤더만 표시됩니다
- 최상위 키-값 쌍은 값이 보존됩니다

`--include-private`를 사용하여 비공개/unexported 심볼 포함.

### 주석

- 단일 행 주석(`# 주석`)은 문서로 추출됩니다

### 임포트

- TOML에는 임포트 시스템이 없습니다. `--include-imports`는 효과가 없습니다

### 제한사항

- 인라인 테이블과 인라인 배열은 별도로 추출되지 않습니다
- 점 표기법 키(예: `physical.color = "orange"`)는 하나의 쌍으로 캡처됩니다
