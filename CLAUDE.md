# CLAUDE.md

> Claude가 이 프로젝트에서 작업할 때 따라야 할 컨벤션과 가이드라인입니다.

---

## 프로젝트 개요

**Brf.it**은 AI 코딩 어시스턴트에게 코드를 효율적으로 전달하기 위한 CLI 도구입니다.

- 프로젝트 코드베이스에서 함수 시그니처와 문서만 추출
- AI가 핵심 내용을 빠르게 파악할 수 있는 형태로 변환 (XML, Markdown)
- Tree-sitter 기반 지능형 코드 압축

---

## Go Engineering Conventions

Claude, brf.it 프로젝트를 **Go(Golang)**로 구현할 때 준수해야 할 엄격한 컨벤션입니다. Go의 특성을 살려 빠르고, 작고, 견고한 CLI 도구를 만드는 데 집중하세요.

### 1. Go Idiomatic Style

- **Standard Tooling**: `gofmt`와 `goimports`를 반드시 적용하여 표준 스타일을 준수하세요.
- **Error Handling**: Go의 관례에 따라 에러를 항상 명시적으로 반환(`result, err := ...`)하고 처리하세요. 에러를 무시(`_`)하지 마세요.
- **Effective Names**:
  - 패키지 명은 짧고 간결하게(소문자 단일 단어) 작성합니다.
  - 외부로 노출(Export)할 구조체와 함수는 `PascalCase`, 내부용은 `camelCase`를 사용합니다.

### 2. Architecture & Patterns

- **Interface-First**: `Extractor`, `Parser`, `Formatter`는 반드시 인터페이스로 정의하세요. 이는 나중에 Tree-sitter 외의 다른 파서를 도입하거나 테스트(Mocking)할 때 필수적입니다.
- **Project Layout**: Go 표준 레이아웃을 따릅니다.
  ```
  brf.it/
  ├── cmd/
  │   └── brfit/
  │       └── main.go        # 진입점
  ├── pkg/
  │   ├── scanner/           # 파일 시스템 스캔
  │   ├── parser/            # 파서 인터페이스 및 구현
  │   ├── extractor/         # Signature 추출
  │   └── formatter/         # 출력 포맷터 (XML, Markdown)
  ├── internal/
  │   └── ...                # 외부 호출 불가능한 전용 로직
  ├── assets/
  │   └── wasm/              # Tree-sitter WASM (필요시)
  └── README.md
  ```
- **Composition over Inheritance**: 구조체 임베딩(Embedding)을 활용하여 기능을 조합하세요.

### 3. CGO & Tree-sitter (중요)

- **CGO Handling**: Go에서 Tree-sitter를 쓰려면 CGO 바인딩이 필요합니다. 런타임 의존성을 줄이기 위해 가능한 **정적 빌드(Static Build)**가 가능하도록 설계하세요.
- **Concurrency**: 파일 스캔(Scanner)과 분석(Extractor) 시 Go의 강력한 `Goroutine`과 `Channel`을 활용하여 성능을 극대화하세요. (단, 과도한 고루틴 생성을 방지하기 위해 Worker Pool 패턴 고려)

### 4. CLI & UX

- **Cobra/Viper**: CLI 명령처리는 `spf13/cobra` 라이브러리를 사용합니다.
- **Zero Config**: 사용자가 아무 옵션 없이 실행해도 최적의 결과(sig 모드, stdout 출력)를 내도록 Default 설정을 똑똑하게 잡으세요.

### 5. Documentation & Quality

- **GoDoc**: 모든 Export 된 요소에는 GoDoc 스타일 주석을 필수적으로 작성하세요.
- **Testing**: `go test`를 활용해 테이블 기반 테스트(Table-driven tests)를 작성하세요.
- **Parser Import**: Tree-sitter 파서를 사용하는 테스트에서는 `_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"` blank import로 파서 자동 등록을 트리거해야 합니다.

---

## Commit Convention

이 프로젝트의 모든 커밋 메시지는 아래 규칙을 엄격히 따라야 합니다.

### 기본 형식

```
type: 요약 설명
```

- 요약 설명은 **한글**로 작성하며, 핵심 내용을 간결하게 담습니다.

### Type 종류

| Type       | 설명                                     | 예시                                                  |
| ---------- | ---------------------------------------- | ----------------------------------------------------- |
| `feat`     | 새로운 기능 추가                         | `feat: Tree-sitter를 이용한 Signature 추출 기능 구현` |
| `fix`      | 버그 수정                                | `fix: 대용량 파일 처리 시 메모리 누수 수정`           |
| `docs`     | 문서 수정                                | `docs: README에 설치 가이드 추가`                     |
| `style`    | 코드 포맷팅 (로직 변경 없음)             | `style: gofmt 적용`                                   |
| `refactor` | 코드 리팩토링 (기능 변화 없는 구조 개선) | `refactor: Parser 인터페이스 분리`                    |
| `test`     | 테스트 코드 추가 및 수정                 | `test: Scanner 단위 테스트 작성`                      |
| `chore`    | 빌드, 패키지 매니저, 프로젝트 설정 변경  | `chore: go mod init 실행`                             |

### 작성 규칙

- 커밋 메시지 본문이 필요한 경우, 한 줄을 띄우고 상세 내용을 작성합니다.
- 과거 시제가 아닌 **현재 시제**를 사용하여 "무엇을 하는 커밋인지" 명확히 합니다.

---

## CLI 인터페이스 (계획)

```bash
brfit [path] [options]

Options:
  -m, --mode <mode>       출력 모드: "sig" (기본값)
  -f, --format <format>   출력 포맷: "xml" (기본값) | "md"
  -o, --output <file>     출력 파일 경로 (기본값: stdout)
  -i, --ignore <file>     커스텀 ignore 파일 (기본값: .gitignore)
  --include-hidden        숨김 파일 포함
  --include-body          함수 본문 포함 (기본값: 시그니처만)
  --no-tree               디렉토리 트리 생략
  --no-tokens             토큰 수 계산 비활성화
  --max-size <bytes>      최대 파일 크기 (기본값: 512000 = 500KB)
```

---

## 핵심 기술 노트

### .gitignore 주의사항

바이너리 이름(`brfit`)이 디렉토리(`cmd/brfit`)와 겹치면 `/brfit`처럼 루트 경로로 지정해야 함

### 릴리즈 및 배포

- **GoReleaser ldflags**: `-X main.version` 등은 `main` 패키지 변수만 주입 가능. `cmd/brfit/main.go`에 변수 선언 필수
- **릴리즈 트리거**: `git tag v*` push 시 GitHub Actions 자동 실행 (`.github/workflows/release.yml`)
- **CGO cross-compile**: zig cc 사용 (`.goreleaser.yml` 참조)
- **릴리즈 노트**: 배포 후 GitHub Release (`gh release edit v*`) + `CHANGELOG.md` 둘 다 업데이트 (템플릿: `.github/RELEASE_TEMPLATE.md`)
- **CHANGELOG 형식**: [Keep a Changelog](https://keepachangelog.com/) 형식 사용

### 릴리즈 절차

```bash
# 1. 변경사항 커밋
git add <files> && git commit -m "feat: ..."

# 2. 푸시 및 태그 생성
git push origin main
git tag v0.X.0 && git push origin v0.X.0

# 3. GitHub Actions 완료 대기 (~5분)
gh run list --limit 1

# 4. 릴리즈 노트 업데이트
gh release edit v0.X.0 --notes "$(cat <<'EOF'
## Release Notes
...
EOF
)"

# 5. CHANGELOG.md 업데이트 후 커밋
```

### ScanOptions 기본값 사용

`ScanOptions` 구조체는 부분적으로 설정할 때 설정하지 않은 필드가 zero value가 됩니다. 기본값을 유지하려면 `DefaultScanOptions()` 호출 후 필요한 필드만 수정하세요:

```go
defaultOpts := scanner.DefaultScanOptions()
scanOpts := &scanner.ScanOptions{
    RootPath:            tmpDir,
    SupportedExtensions: defaultOpts.SupportedExtensions,
    MaxFileSize:         defaultOpts.MaxFileSize,
}
```

### 대용량 파일 방어 로직

수만 줄짜리 min.js나 vendor 파일이 포함될 경우 메모리 부하 방지:

```go
const DefaultMaxFileSize = 500 * 1024 // 500KB

if fileSize > maxFileSize {
    log.Printf("Skipping large file: %s (%d bytes)", path, fileSize)
    return nil
}
```

### 출력 포맷 예시 (XML)

```xml
<brfit>
  <metadata>
    <tree>...</tree>
  </metadata>
  <files>
    <file path="src/scanner.go" language="go">
      <signature>func ScanDirectory(root string) ([]FileEntry, error)</signature>
      <doc>Recursively scans directory and returns supported files.</doc>
    </file>
  </files>
</brfit>
```

### 다국어 문서

- 파일명 suffix 방식: `README.ko.md`, `docs/languages/go.ja.md`
- 지원 언어: EN (기본), KO, JA, HI, DE
- 모든 문서 상단에 언어 선택 링크 추가: `🌐 [English](file.md) | [한국어](file.ko.md) | ...`

### 새 언어 추가 체크리스트

새 프로그래밍 언어 지원 추가 시 반드시 확인:

1. `pkg/parser/treesitter/languages/[lang].go` - LanguageQuery 구현
2. `pkg/parser/treesitter/parser.go` - init()에 파서 등록, queries 맵에 추가, isExported(), stripBody() 케이스 추가
3. `pkg/scanner/scanner.go` - DefaultScanOptions()에 확장자 추가
4. `internal/config/config.go` - SupportedExtensions()에 확장자 추가 (CLI에서 사용)
5. `docs/languages/[lang].md` + 다국어 버전 생성
6. `README*.md` Supported Languages 테이블 업데이트
7. `pkg/parser/treesitter/languages/[lang]_test.go` - 단위 테스트 작성
8. `pkg/parser/treesitter/parser_test.go` - TestTreeSitterParserLanguages, TestTreeSitterParserAutoRegistration에 언어 추가, 통합 테스트 함수 추가

### Tree-sitter AST 디버깅

새 쿼리 패턴 작성 시 AST 구조 확인이 필요하면:

```go
// 임시 디버그 코드로 AST 출력
func printTree(node *sitter.Node, code []byte, indent int) {
    fmt.Printf("%s%s\n", strings.Repeat("  ", indent), node.Kind())
    for i := uint(0); i < uint(node.ChildCount()); i++ {
        printTree(node.Child(i), code, indent+1)
    }
}
```

**포인터 반환 타입 주의**: `User* func()` 형태는 declarator가 `pointer_declarator` 안에 중첩됨
