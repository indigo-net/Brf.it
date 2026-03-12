---
title: YAML
---

# YAML 지원

[English](../../languages/yaml.md) | [한국어](yaml.md) | [日本語](../../ja/languages/yaml.md) | [हिन्दी](../../hi/languages/yaml.md) | [Deutsch](../../de/languages/yaml.md)

## 지원 확장자

- `.yaml`
- `.yml`

## 문법

- [tree-sitter-yaml](https://github.com/tree-sitter-grammars/tree-sitter-yaml) v0.7.2 by tree-sitter-grammars

## 추출 대상

| 요소 | 종류 | XML 태그 | 예시 |
|------|------|----------|------|
| 키-값 쌍 | `variable` | `<variable>` | `name: value` |

## 예시

### 입력

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

### 출력 (XML)

```xml
<file path="config.yaml" language="yaml">
  <variable>name: myapp</variable>
  <variable>version: 1.0.0</variable>
  <variable>database:</variable>
  <variable>features:</variable>
</file>
```

## 참고사항

### 본문 제거

`--include-body` 플래그를 사용하지 않을 때:

- 컨테이너 키(중첩 값을 가진 매핑)는 중첩 내용이 제거되고, 키만 표시됩니다
- 스칼라 키-값 쌍은 값이 보존됩니다

`--include-private`를 사용하여 비공개/unexported 심볼 포함.

### 주석

- 단일 행 주석(`# 주석`)은 문서로 추출됩니다

### 임포트

- YAML에는 임포트 시스템이 없습니다. `--include-imports`는 효과가 없습니다

### 제한사항

- 과도한 노이즈를 방지하기 위해 최상위 키만 시그니처로 캡처됩니다
- 앵커와 별칭은 특별히 처리되지 않습니다
