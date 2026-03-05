---
name: add-language-support
description: Use when adding a new programming language to Brf.it - covers LanguageQuery implementation, parser registration, scanner/config extensions, tests, and multilingual documentation
---

# Add Language Support

## Overview

Brf.it에 새 프로그래밍 언어를 추가하는 전체 체크리스트와 구현 가이드. 기존 8개 언어(Go, TypeScript, JavaScript, Python, C, C++, Java, Rust) 구현 패턴에서 추출.

### 핵심 원칙: 축약은 허용, 누락은 금지

패키징 과정에서 원본 코드의 문법이 **축약**(본문 제거, 시그니처만 추출)되는 것은 허용되지만, 선언이 **완전히 누락**되는 것은 반드시 피해야 한다.

- **축약 (OK)**: `func distance() -> Double { return sqrt(x*x+y*y) }` → `func distance() -> Double`
- **누락 (NG)**: 해당 언어의 특정 선언 유형(예: protocol, extension, subscript)이 쿼리 패턴에서 빠져 출력에 아예 나타나지 않는 것

새 언어의 쿼리 패턴 작성 시 해당 언어가 지원하는 **모든 선언 유형**을 빠짐없이 캡처해야 한다. 실제 코드베이스로 실행 테스트를 수행하여 누락 여부를 반드시 확인할 것.

## Checklist

| 단계 | 파일 | 작업 | 참고 |
|------|------|------|------|
| 1 | `pkg/parser/treesitter/languages/{lang}.go` | LanguageQuery 구현 | ~170줄 |
| 2 | `pkg/parser/treesitter/parser.go` | init() 등록, queries 맵, isExported(), stripBody() | +15줄 |
| 3 | `pkg/scanner/scanner.go` | DefaultScanOptions() 확장자 추가 | +1~3줄 |
| 4 | `internal/config/config.go` | SupportedExtensions() 확장자 추가 | +1~3줄 |
| 5 | `go.mod` / `go.sum` | tree-sitter 파서 의존성 추가 | `go get` |
| 6 | `pkg/parser/treesitter/languages/{lang}_test.go` | 14개 단위 테스트 | ~400줄 |
| 7 | `pkg/parser/treesitter/parser_test.go` | 통합 테스트 3곳 추가 | +30줄 |
| 8 | `docs/languages/{lang}.md` + 4개 번역 | 언어 가이드 (EN, KO, JA, HI, DE) | 5파일 × ~130줄 |
| 9 | `README.md` × 5개 | Supported Languages 테이블 행 추가 | 5파일 × +1줄 |

**예상 규모**: 약 19개 파일, ~1,800줄 (구현+테스트+문서)

## Step 1: LanguageQuery 구현

### 템플릿 (`pkg/parser/treesitter/languages/{lang}.go`)

```go
package languages

import (
    sitter "github.com/tree-sitter/go-tree-sitter"
    tree_sitter_{lang} "github.com/{org}/tree-sitter-{lang}/bindings/go"
)

// {Lang}Query implements LanguageQuery for {Lang} language.
type {Lang}Query struct {
    language *sitter.Language
    query    []byte
}

// New{Lang}Query creates a new {Lang} language query.
func New{Lang}Query() *{Lang}Query {
    return &{Lang}Query{
        language: sitter.NewLanguage(tree_sitter_{lang}.Language()),
        query:    []byte({lang}QueryPattern),
    }
}

// Language returns the {Lang} Tree-sitter language.
func (q *{Lang}Query) Language() *sitter.Language { return q.language }

// Query returns the {Lang} query pattern.
func (q *{Lang}Query) Query() []byte { return q.query }

// Captures returns the capture names for {Lang} queries.
func (q *{Lang}Query) Captures() []string {
    return []string{captureName, captureSignature, captureDoc, captureKind}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *{Lang}Query) KindMapping() map[string]string {
    return map[string]string{
        // AST 노드 타입 → Kind 문자열
    }
}

// ImportQuery returns the {Lang} import query pattern.
func (q *{Lang}Query) ImportQuery() []byte {
    return []byte({lang}ImportQueryPattern)
}

const {lang}ImportQueryPattern = `
; Import 문 캡처
`

const {lang}QueryPattern = `
; 시그니처 캡처 쿼리
`
```

### 5개 필수 인터페이스 메서드

`LanguageQuery` 인터페이스 (`pkg/parser/treesitter/query.go`):
- `Language() *sitter.Language`
- `Query() []byte`
- `ImportQuery() []byte`
- `Captures() []string`
- `KindMapping() map[string]string`

### Kind Mapping 규칙

XML 출력 시 `parser.Signature.Kind` → XML 태그 매핑 (`pkg/formatter/xml.go`의 `kindToTag()`):

| XML 태그 | Kind 값 |
|----------|---------|
| `<function>` | function, method, constructor, destructor, arrow |
| `<type>` | class, interface, type, struct, enum, record, annotation, typedef, namespace, template, impl, trait |
| `<variable>` | variable, field, macro, export |
| `<signature>` | 빈 문자열 또는 알 수 없는 Kind (fallback) |

## Step 2: parser.go 등록

4곳 수정:

```go
// 1. init() - 파서 등록
parser.RegisterParser("{lang}", NewTreeSitterParser())

// 2. queries 맵
"{lang}": languages.New{Lang}Query(),

// 3. isExported() - case 추가
case "{lang}":
    return true  // 또는 언어별 가시성 규칙

// 4. stripBody() - case 추가 + strip{Lang}Body() 함수
case "{lang}":
    return strip{Lang}Body(text, kind)
```

### stripBody 구현 패턴

대부분의 언어는 `{` 기준으로 본문 제거. 제네릭 `<>` 처리가 필요한 언어는 `findXxxBodyStart()` 헬퍼 사용:

```go
func find{Lang}BodyStart(text string) int {
    parenDepth, angleDepth := 0, 0
    for i, ch := range text {
        switch ch {
        case '(': parenDepth++
        case ')': parenDepth--
        case '<': angleDepth++
        case '>':
            if angleDepth > 0 { angleDepth-- }
        case '{':
            if angleDepth == 0 && parenDepth == 0 { return i }
        }
    }
    return -1
}
```

## Step 3-4: Scanner / Config

```go
// pkg/scanner/scanner.go - DefaultScanOptions()
".{ext}": "{lang}",

// internal/config/config.go - SupportedExtensions()
".{ext}": "{lang}",
```

## Step 5: 의존성

```bash
go get github.com/{org}/tree-sitter-{lang}/bindings/go
```

## Step 6: 단위 테스트 (14개)

`pkg/parser/treesitter/languages/{lang}_test.go`:

1. `Test{Lang}QueryLanguage` — `Language() != nil`
2. `Test{Lang}QueryPattern` — 쿼리 컴파일 성공
3. `Test{Lang}QueryImportPattern` — import 쿼리 컴파일 성공
4. `Test{Lang}QueryExtractFunction` — 함수 추출 (일반, async, 수식어)
5. `Test{Lang}QueryExtractTypes` — class/struct/enum 등 타입 추출
6-10. 언어별 특화 테스트 (protocol, extension, properties 등)
11. `Test{Lang}QueryExtractImport` — import 문 추출
12. `Test{Lang}QueryExtractGenerics` — 제네릭 타입 파라미터
13. `Test{Lang}QueryKindMapping` — Kind 매핑 검증
14. `Test{Lang}QueryCaptures` — 캡처 이름 검증

### 테스트 패턴

```go
func Test{Lang}QueryExtractXxx(t *testing.T) {
    parser := sitter.NewParser()
    defer parser.Close()
    lang := sitter.NewLanguage(tree_sitter_{lang}.Language())
    parser.SetLanguage(lang)

    code := []byte(`...`)
    tree := parser.Parse(code, nil)
    defer tree.Close()

    query := New{Lang}Query()
    q, err := sitter.NewQuery(lang, string(query.Query()))
    if err != nil { t.Fatalf("failed to create query: %v", err) }
    defer q.Close()

    qc := sitter.NewQueryCursor()
    defer qc.Close()
    matches := qc.Matches(q, tree.RootNode(), code)

    captureNames := q.CaptureNames()
    foundNames := make(map[string]bool)
    for {
        match := matches.Next()
        if match == nil { break }
        for _, c := range match.Captures {
            if captureNames[c.Index] == "name" {
                foundNames[string(code[c.Node.StartByte():c.Node.EndByte()])] = true
            }
        }
    }

    for _, expected := range []string{...} {
        if !foundNames[expected] {
            t.Errorf("expected to find '%s'", expected)
        }
    }
}
```

## Step 7: 통합 테스트

`pkg/parser/treesitter/parser_test.go`에 3곳:

```go
// 1. TestTreeSitterParserLanguages - expected 슬라이스에 "{lang}" 추가
// 2. TestTreeSitterParserAutoRegistration - 언어 슬라이스에 "{lang}" 추가
// 3. TestTreeSitterParserParse{Lang} - 새 함수 추가 (~50줄)
```

## Step 8-9: 문서

### 언어 가이드 (`docs/languages/{lang}.md`)

- 언어 선택 링크 (EN, KO, JA, HI, DE)
- Supported Extensions
- Extraction Targets 테이블
- Example: Input 코드 → Output XML
- Notes: 언어 특화 참고사항

### 번역본 (4개)

`docs/{ko,ja,hi,de}/languages/{lang}.md` — 동일 구조, 각 언어로 번역

### README 테이블 (5개 파일)

| 파일 | 추가할 행 패턴 |
|------|---------------|
| `README.md` | `\| {Lang} \| \`.{ext}\` \| [{Lang} Guide](docs/languages/{lang}.md) \|` |
| `docs/ko/README.md` | `\| {Lang} \| \`.{ext}\` \| [{Lang} 가이드](languages/{lang}.md) \|` |
| `docs/ja/README.md` | `\| {Lang} \| \`.{ext}\` \| [{Lang}ガイド](languages/{lang}.md) \|` |
| `docs/de/README.md` | `\| {Lang} \| \`.{ext}\` \| [{Lang}-Leitfaden](languages/{lang}.md) \|` |
| `docs/hi/README.md` | `\| {Lang} \| \`.{ext}\` \| [{Lang} गाइड](languages/{lang}.md) \|` |

## Tree-sitter AST 디버깅

새 쿼리 패턴 작성 시 AST 구조 확인이 필요하면:

```go
func printTree(node *sitter.Node, code []byte, indent int) {
    fmt.Printf("%s%s\n", strings.Repeat("  ", indent), node.Kind())
    for i := uint(0); i < uint(node.ChildCount()); i++ {
        printTree(node.Child(i), code, indent+1)
    }
}
```

**포인터 반환 타입 주의**: `User* func()` 형태는 declarator가 `pointer_declarator` 안에 중첩됨

## Import 쿼리 패턴 작성

전체 import 문을 캡처하려면 노드 전체를 캡처:

```scheme
; 경로만: (import_statement source: (string) @import_path)
; 전체 문: (import_statement) @import_path
```

**Go 예외**: `import_spec`은 `"fmt"` 형태라 `cleanImportPath()`에서 `import ` prefix 추가 처리

**cleanImportPath 자동 인식**: `import `, `from `, `#include`, `use `, `extern crate` 프리픽스가 있으면 그대로 반환 (`parser.go`의 `cleanImportPath()`)

## 검증 단계

```bash
# 1. 빌드
go build ./cmd/brfit

# 2. 단위 테스트
go test ./pkg/parser/treesitter/languages/ -run {Lang} -v

# 3. 통합 테스트
go test ./pkg/parser/treesitter/ -run {Lang} -v

# 4. 전체 테스트 (회귀 방지)
go test ./...

# 5. 실행 테스트
echo '...' > /tmp/test.{ext}
brfit /tmp -f xml
brfit /tmp -f md
rm /tmp/test.{ext}

# 6. 커버리지 확인 (누락 방지)
# 해당 언어의 모든 선언 유형을 포함하는 종합 테스트 파일로 실행하여
# 원본 코드의 선언이 빠짐없이 출력되는지 확인한다.
# brfit 출력의 시그니처 수와 원본 코드의 선언 수를 대조할 것.
```

## Common Mistakes

- **선언 유형 누락**: 해당 언어의 선언 유형(함수, 타입, 프로퍼티 등)이 쿼리 패턴에서 빠지면 출력에서 완전히 사라짐. 언어 공식 문서에서 선언 유형 목록을 확인하고, AST 디버깅(`printTree`)으로 모든 노드 타입이 캡처되는지 검증할 것
- **쿼리 패턴 구문 오류**: `sitter.NewQuery()` 컴파일 테스트 먼저 통과시키기
- **cleanImportPath 미처리**: 새 언어의 import 프리픽스가 `cleanImportPath()`에 인식되지 않으면 추가 필요
- **stripBody 미구현**: `parser.go`의 `stripBody()` switch에 case 누락 시 본문이 제거되지 않음
- **확장자 누락**: `scanner.go`와 `config.go` 둘 다 추가해야 함 (하나만 추가하면 CLI/라이브러리 불일치)
- **통합 테스트 누락**: `parser_test.go`의 3곳(Languages, AutoRegistration, Parse{Lang}) 모두 추가
- **문서 누락**: 5개 README + 5개 언어 가이드 = 10개 문서 파일 확인
