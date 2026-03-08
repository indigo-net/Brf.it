---
layout: default
title: Bash/Shell
parent: 언어 가이드
nav_order: 13
---

# Bash/Shell 지원

[English](../../languages/bash.md) | [한국어](bash.md) | [日本語](../../ja/languages/bash.md) | [हिन्दी](../../hi/languages/bash.md) | [Deutsch](../../de/languages/bash.md)

## 지원 확장자

- `.sh`
- `.bash`

## 추출 대상

| 요소 | 종류 | 예시 |
|------|------|------|
| 함수 | `function` | `function greet { ... }` |
| 함수 | `function` | `greet() { ... }` |
| 변수 할당 | `variable` | `NAME="value"` |
| 선언 | `variable` | `declare VERBOSE` |
| 지역 변수 | `variable` | `local count=0` |
| 읽기 전용 변수 | `variable` | `readonly VERSION="1.0"` |
| 주석 | `doc` | `# 설명` |
| source 문 | `import` | `source /path/to/lib.sh` |
| dot 문 | `import` | `. ./config.sh` |

## 예시

### 입력

```bash
#!/bin/bash

# 설정
CONFIG_PATH="/etc/myapp"
VERSION="1.0.0"
declare VERBOSE=false

# 애플리케이션 배포
function deploy {
    local app_name="$1"
    echo "Deploying $app_name"
}

# 프로젝트 빌드
build() {
    echo "Building..."
}

source ./utils.sh
. ./config.sh
```

### 출력 (XML)

```xml
<file path="deploy.sh" language="bash">
  <variable kind="variable" line="4">
    <name>CONFIG_PATH</name>
    <text>CONFIG_PATH="/etc/myapp"</text>
  </variable>
  <variable kind="variable" line="5">
    <name>VERSION</name>
    <text>VERSION="1.0.0"</text>
  </variable>
  <variable kind="variable" line="6">
    <name>VERBOSE</name>
    <text>declare VERBOSE=false</text>
  </variable>
  <function kind="function" line="9">
    <name>deploy</name>
    <text>function deploy</text>
  </function>
  <function kind="function" line="15">
    <name>build</name>
    <text>build()</text>
  </function>
</file>
```

## 참고사항

### 가시성

- 모든 선언이 추출됩니다 (Bash는 접근 제어자가 없음)
- 함수 내 `local` 변수도 파싱 시점에 선언되어 있으면 추출됩니다

### 함수 문법

Bash는 두 가지 함수 선언 스타일을 지원합니다:

- `function 이름 { ... }` - `function` 키워드 사용
- `이름() { ... }` - 괄호 사용

두 스타일 모두 `function` 종류로 추출됩니다.

### 본문 제거

`--include-body` 플래그를 사용하지 않으면:

- 함수: 여는 중괄호 `{` 이후 본문 제거
- 변수: 첫 번째 줄만 유지 (여러 줄 할당 처리)

### 가져오기 추출

- `source` 및 `.` 명령이 `--include-imports` 플래그로 추출됩니다
- 따옴표 있음/없음 경로 모두 지원합니다

### 문서 주석

- `#`으로 시작하는 셸 주석이 추출됩니다
- 셔뱅 라인(`#!/bin/bash`)은 주석으로 처리되지 않습니다

### 제한사항

- 중첩 함수가 지원됩니다
- 함수 본문의 here-document도 올바르게 처리됩니다
- 복잡한 변수 확장은 시그니처에 그대로 보존됩니다
