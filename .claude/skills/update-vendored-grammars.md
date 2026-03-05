---
name: update-vendored-grammars
description: Use when checking or updating vendored tree-sitter grammar C sources from their original repositories - covers version checking, cloning, generating, copying, and testing
---

# Update Vendored Grammars

## Overview

Brf.it는 일부 tree-sitter 문법을 vendor 방식(C 소스 직접 포함)으로 관리한다. 이 스킬은 원본 repo가 갱신되었을 때 vendored C 소스를 동기화하는 절차를 안내한다.

## Vendor 대상 레지스트리

| 언어 | 원본 repo | Vendored 경로 | generate 필요 | 현재 버전 |
|------|-----------|--------------|--------------|----------|
| Kotlin | `fwcd/tree-sitter-kotlin` | `pkg/parser/treesitter/grammars/kotlin/` | No (parser.c 포함) | v0.3.8 |
| Swift | `alex-pinkus/tree-sitter-swift` | `pkg/parser/treesitter/grammars/swift/` | Yes (parser.c 미포함) | v0.7.1 |

## Checklist

| 단계 | 작업 |
|------|------|
| 1 | 원본 repo 최신 버전 확인 |
| 2 | 임시 디렉토리에 clone |
| 3 | parser.c 생성 (필요 시) |
| 4 | C 소스 파일 복사 |
| 5 | binding.go 버전 주석 갱신 |
| 6 | 빌드 및 전체 테스트 |
| 7 | 임시 파일 정리 |

## Step 1: 최신 버전 확인

원본 repo의 최신 릴리스/커밋과 `binding.go` 주석의 버전을 비교한다.

```bash
# 현재 vendored 버전 확인
cat pkg/parser/treesitter/grammars/{lang}/binding.go | head -5

# 원본 repo 최신 릴리스 확인
gh api repos/{org}/tree-sitter-{lang}/releases/latest --jq '.tag_name' 2>/dev/null \
  || gh api repos/{org}/tree-sitter-{lang}/commits?per_page=1 --jq '.[0].sha[:7]'
```

## Step 2: 임시 디렉토리에 clone

```bash
git clone --depth 1 https://github.com/{org}/tree-sitter-{lang} /tmp/tree-sitter-{lang}
# 특정 태그가 있으면:
# git clone --depth 1 --branch {tag} https://github.com/{org}/tree-sitter-{lang} /tmp/tree-sitter-{lang}
```

## Step 3: parser.c 생성

레지스트리의 "generate 필요" 열을 확인한다.

### generate 불필요 (parser.c가 src/에 이미 존재)

바로 Step 4로 진행.

### generate 필요

```bash
cd /tmp/tree-sitter-{lang}
npm install
npx tree-sitter generate
```

`npm install` 시 `tree-sitter generate`가 postinstall 스크립트로 자동 실행되는 경우도 있다. `src/parser.c`가 생성되었는지 확인:

```bash
ls -la /tmp/tree-sitter-{lang}/src/parser.c
```

## Step 4: C 소스 파일 복사

vendored 디렉토리 구조:

```
pkg/parser/treesitter/grammars/{lang}/
├── binding.go
├── parser.c
├── scanner.c          # 있는 경우만
└── tree_sitter/
    ├── alloc.h
    ├── array.h
    └── parser.h
```

복사 명령:

```bash
LANG={lang}
SRC=/tmp/tree-sitter-$LANG/src
DEST=pkg/parser/treesitter/grammars/$LANG

cp $SRC/parser.c $DEST/parser.c
# scanner.c가 있으면 복사
[ -f $SRC/scanner.c ] && cp $SRC/scanner.c $DEST/scanner.c

mkdir -p $DEST/tree_sitter
cp $SRC/tree_sitter/parser.h $DEST/tree_sitter/parser.h
cp $SRC/tree_sitter/alloc.h $DEST/tree_sitter/alloc.h
cp $SRC/tree_sitter/array.h $DEST/tree_sitter/array.h
```

## Step 5: binding.go 버전 주석 갱신

`binding.go`의 주석에서 버전을 갱신한다:

```go
// The C source files (parser.c, scanner.c) are vendored from
// https://github.com/{org}/tree-sitter-{lang} ({new_version}).
```

## Step 6: 빌드 및 전체 테스트

```bash
# 빌드 확인
go build ./...

# 해당 언어 테스트
go test ./pkg/parser/treesitter/languages/ -run {Lang} -v

# 전체 테스트 (회귀 방지)
go test ./...
```

## Step 7: 임시 파일 정리

```bash
rm -rf /tmp/tree-sitter-{lang}
```

## binding.go 템플릿

새 언어를 vendor로 추가할 때 사용:

```go
// Package {lang} provides the tree-sitter grammar for {Lang}.
//
// The C source files (parser.c, scanner.c) are vendored from
// https://github.com/{org}/tree-sitter-{lang} ({version}).
// CGO automatically compiles all .c files in this directory.
package {lang}

// #cgo CFLAGS: -std=c11 -fPIC
// extern const void *tree_sitter_{lang}(void);
import "C"

import "unsafe"

// Language returns the tree-sitter Language for {Lang}.
func Language() unsafe.Pointer {
	return unsafe.Pointer(C.tree_sitter_{lang}())
}
```

**주의사항**:
- `extern` 선언을 사용한다 (`#include "parser.c"` 금지 — CGO가 디렉토리 내 모든 .c를 자동 컴파일하므로 중복 심볼 발생)
- import 경로: `github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/{lang}`

## Common Mistakes

- **generate 누락**: parser.c가 repo에 포함되지 않은 경우 `npx tree-sitter generate` 실행 필요
- **header 복사 누락**: `tree_sitter/` 디렉토리의 `alloc.h`, `array.h`, `parser.h` 3개 모두 복사해야 함
- **#include 사용**: `binding.go`에서 `#include "parser.c"` 대신 `extern` 선언 사용 (CLAUDE.md 참조)
- **go.mod에 외부 의존성 잔존**: vendor 전환 후 `go mod tidy`로 이전 의존성 제거 확인
