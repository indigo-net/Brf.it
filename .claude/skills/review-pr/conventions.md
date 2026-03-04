# Brf.it Project Conventions (Review Reference)

> PR 코드 리뷰 시 참조할 프로젝트 컨벤션 요약입니다.

---

## Go Idiomatic Style

- **gofmt / goimports** 필수 적용
- **에러 처리**: 모든 에러를 명시적으로 반환하고 처리. `_`로 에러 무시 금지
  ```go
  result, err := doSomething()
  if err != nil {
      return fmt.Errorf("context: %w", err)
  }
  ```
- **네이밍**: 패키지명은 소문자 단일 단어, Export는 `PascalCase`, 내부는 `camelCase`

---

## Interface-First Architecture

핵심 인터페이스는 반드시 정의 후 구현:

| 인터페이스 | 위치 | 메서드 |
|-----------|------|--------|
| `Scanner` | `pkg/scanner/` | `Scan() (*ScanResult, error)` |
| `Parser` | `pkg/parser/` | `Parse()`, `Languages()` |
| `Extractor` | `pkg/extractor/` | `Extract()` |
| `Formatter` | `pkg/formatter/` | `Format()`, `Name()` |
| `LanguageQuery` | `pkg/parser/treesitter/` | `Language()`, `Query()`, `Captures()`, `KindMapping()` |

---

## Tree-sitter / CGO Patterns

- `github.com/tree-sitter/go-tree-sitter` 바인딩 사용
- **리소스 해제**: `defer parser.Close()`, `defer tree.Close()` 필수
- 언어 쿼리는 `LanguageQuery` 인터페이스로 정의
- 파서 등록: `init()` 함수에서 `parser.RegisterParser("go", NewTreeSitterParser())` 패턴

---

## Registry Pattern

- `sync.RWMutex`로 스레드 안전한 파서 레지스트리
- `DefaultRegistry()` 싱글턴
- `Register()`, `Get()`, `Languages()` 메서드

---

## Default Options Pattern

- `DefaultScanOptions()`, `DefaultConfig()`, `DefaultOptions()` 팩토리 함수로 기본값 제공
- `nil` 전달 시 기본값 사용하도록 설계

---

## Testing Conventions

- **Table-driven tests** 필수:
  ```go
  tests := []struct {
      name      string
      input     Type
      expected  Type
      wantError bool
  }{...}
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) { ... })
  }
  ```
- **Blank import**로 파서 등록: `_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"`
- `go test ./pkg/...` 및 `go build` 통과 필수

---

## Formatter Patterns

- `bytes.Buffer` 기반 문자열 빌딩
- `NewXMLFormatter()`, `NewMarkdownFormatter()` 생성자
- `Name()` 메서드: 소문자 포맷터 식별자 반환
- XML 출력 시 `escapeXML()` 헬퍼 사용

---

## Scanner Patterns

- `DefaultScanOptions()` 기반 옵션 관리
- `filepath.WalkDir()` 디렉토리 순회
- `IsHidden()` 헬퍼로 숨김 파일 감지
- `MaxFileSize` (기본 500KB) 초과 시 스킵 + 로깅
- `go-gitignore` 라이브러리로 .gitignore 지원

---

## Packager Composition

- `Packager` 구조체가 Scanner -> Extractor -> Formatter -> Tokenizer 파이프라인 조율
- `NewDefaultPackager()` 팩토리
- 선택적 컴포넌트(tokenizer)는 graceful fallback

---

## Commit Convention

```
type: 한글 요약 설명
```

| Type | 용도 |
|------|------|
| `feat` | 새 기능 추가 |
| `fix` | 버그 수정 |
| `docs` | 문서 수정 |
| `style` | 코드 포맷팅 (로직 변경 없음) |
| `refactor` | 구조 개선 (기능 변화 없음) |
| `test` | 테스트 코드 |
| `chore` | 빌드/설정 변경 |

- 현재 시제 사용
- 이슈 연결 시: `type: 설명 (#issue)`

---

## GoDoc

- 모든 Export 요소에 GoDoc 스타일 주석 필수
- 패키지 레벨 문서 포함
- 형식: `// FunctionName does X.`

---

## Project Layout

```
brf.it/
├── cmd/brfit/         # CLI 진입점 (Cobra)
├── pkg/
│   ├── scanner/       # 파일 시스템 스캔
│   ├── parser/        # 파서 인터페이스 + Registry
│   │   └── treesitter/  # Tree-sitter 구현
│   ├── extractor/     # Signature 추출
│   └── formatter/     # 출력 포맷터 (XML, Markdown)
├── internal/
│   ├── config/        # CLI 설정
│   └── context/       # Packager 파이프라인
└── assets/wasm/       # Tree-sitter WASM
```
