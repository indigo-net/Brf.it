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
- **Third-party Grammar 통합**: 외부 tree-sitter 문법은 `go get`으로 사용 불가 (C 소스가 Go 모듈 경계 밖에 위치). **Vendor 방식**으로 C 소스를 `pkg/parser/treesitter/grammars/<lang>/`에 직접 복사하여 통합. Fork 방식 사용 금지.
- **binding.go 작성 시**: CGO는 디렉토리 내 모든 `.c` 파일을 자동 컴파일하므로, `#include "parser.c"` 대신 `extern` 선언 사용 (중복 심볼 방지)
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

### 커밋 금지 파일

- **`docs/plans/`**: 설계 문서는 로컬 작업 파일이므로 **절대 커밋하지 않습니다** (`.gitignore`에 등록됨)

### Plan 파일 컨벤션

설계 문서는 `docs/plans/` 하위에 다음 형식으로 작성합니다:

- **파일명**: `YYYY-MM-DD-<topic>-design.md`
- **형식**:

```markdown
# 제목

> 생성일: YYYY-MM-DD
> 상태: 설계 중 | 설계 완료 | 구현 중 | 완료 | 보류

---

## 개요
[배경, 목적, 기대 효과]

## 설계
[구현 방식, 아키텍처, 주요 결정 사항]

## 작업 항목
[구체적인 구현 단계]

## 검증
[테스트 및 확인 방법]
```

> Claude Code의 내부 plan 파일(`~/.claude/plans/`)과 별개로, 프로젝트 맥락의 설계 문서를 이 위치에 저장합니다.

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
  --include-imports       import/export 문 포함 (기본값: 미포함)
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

# 4. 릴리즈 노트 업데이트 (.github/RELEASE_TEMPLATE.md 참조)
gh release edit v0.X.0 --notes "$(cat <<'EOF'
## 🎉 [릴리즈 제목]
...
EOF
)"

# 5. CHANGELOG.md 업데이트 후 커밋
```

### 릴리즈 노트 양식

모든 릴리즈 노트는 아래 양식을 따릅니다 (상세 템플릿: `.github/RELEASE_TEMPLATE.md`):

```markdown
## 🎉 [기능명/릴리즈 요약]

[1-2문장 설명]

### ✨ New Features
### 🔧 Improvements
### 🐛 Bug Fixes
### 📦 Installation
### 📝 Example

**Full Changelog**: https://github.com/indigo-net/Brf.it/compare/vPREV...vNEW
```

**이모지 규칙**:

| 섹션 | 이모지 |
|------|--------|
| 제목 | 🎉 |
| New Features | ✨ |
| Improvements | 🔧 |
| Bug Fixes | 🐛 |
| Breaking Changes | ⚠️ |
| Deprecated | 🗑️ |
| Installation | 📦 |
| Example | 📝 |
| Documentation | 📚 |

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
      <function>func ScanDirectory(root string) ([]FileEntry, error)</function>
      <doc>Recursively scans directory and returns supported files.</doc>
    </file>
  </files>
</brfit>
```

### 다국어 문서

- 국가별 디렉토리 방식: `docs/ko/README.md`, `docs/ja/languages/go.md`
- 지원 언어: EN (기본), KO, JA, HI, DE
- 영어 원본: `README.md` (루트), `docs/languages/*.md`
- 번역본: `docs/{ko,ja,hi,de}/README.md`, `docs/{ko,ja,hi,de}/languages/*.md`
- 모든 문서 상단에 언어 선택 링크 추가: `🌐 [English](../../README.md) | [한국어](README.md) | ...`

### 새 언어 추가

스킬 참조: `.claude/skills/add-language-support.md` (`/add-language-support`)

- **Vendor 참조 구현**: `pkg/parser/treesitter/grammars/kotlin/` (binding.go + C sources)
- **LanguageQuery 참조**: `pkg/parser/treesitter/languages/kotlin.go` (refineKind 패턴 포함)

### GitHub Issue 기반 워크플로우

**주의**: 외부 저장소 fork, 새 저장소 생성 등 사용자 계정에 영향을 주는 작업은 반드시 **사전 승인** 필요

모든 작업은 이슈 기반으로 진행합니다:

1. **이슈 생성**: `gh issue create --assignee indigo-net --label "enhancement"`
2. **브랜치 생성**: `git checkout -b feat/feature-name` (이슈 번호 제외, 기존 스타일 유지)
3. **커밋**: `git commit -m "feat: 구현 내용 (#123)"` (이슈 번호 괄호로 참조)
4. **PR 생성**: `gh pr create --assignee indigo-net` + `Closes #XXX` in body
5. **머지**: PR 머지 시 이슈 자동 닫힘

**브랜치명 형식**: `{type}/{feature-name}` (예: `feat/github-workflow-setup`)

**이슈/PR assignee**: 기본적으로 `indigo-net` 지정

### SAMPLE 파일 생성

프로젝트 코드베이스 요약 파일 생성:
```bash
brfit . -f md --no-tokens --include-imports --no-tree -o SAMPLE.md
brfit . -f xml --no-tokens --include-imports --no-tree -o SAMPLE.xml
```

**자동 갱신**: `.github/workflows/update-code-package.yml`이 main push 또는 release 이벤트 시 자동 실행 (paths-ignore로 무한 루프 방지)

**주의**: `gh release edit` CLI 명령은 GitHub의 `release edited` 웹훅 이벤트를 트리거하지 않음. 수동 트리거가 필요하면 `gh workflow run "Update SAMPLE"` 사용

### 포매터 isEmpty 판정 패턴

`xml.go`, `markdown.go`에서 빈 파일 판정(`isEmpty`)은 import 필터링 **후** 실제 렌더링된 결과를 기준으로 해야 함. `hasRenderedImports` 변수를 import 블록 밖에 선언하고, 렌더링 후 설정하는 패턴 사용:

```go
hasRenderedImports := false
if ... { // import 렌더링 블록
    hasRenderedImports = len(importLines) > 0
}
isEmpty := len(file.Signatures) == 0 && !hasRenderedImports
```

### README 동기화

README 수정 시 5개 파일을 모두 업데이트해야 함:
- `README.md` (영어 원본, 루트)
- `docs/ko/README.md` (한국어)
- `docs/ja/README.md` (일본어)
- `docs/de/README.md` (독일어)
- `docs/hi/README.md` (힌디어)

**팁**: 구조가 동일하므로 동일한 Edit 패턴을 5개 파일에 적용 가능

### GitHub Pages 문서 사이트 동기화

프로젝트 문서 변경 시 웹사이트(https://indigo-net.github.io/Brf.it/)에도 반영 필요:

| 변경 대상 | 동기화할 웹사이트 파일 |
|-----------|------------------------|
| `README.md` 내용 | `index.md` (랜딩 페이지) |
| CLI 옵션 추가/변경 | `docs/cli-reference.md` |
| 새 언어 지원 추가 | `docs/languages/*.md`, `docs/languages/index.md` |
| 설치 방법 변경 | `index.md`, `docs/getting-started.md` |
| 기능 추가/변경 | `index.md`, `docs/getting-started.md` |

**웹사이트 구조**:
- `index.md` — 랜딩 페이지 (README.md 기반)
- `docs/index.md` — 문서 홈
- `docs/getting-started.md` — 설치 및 빠른 시작
- `docs/cli-reference.md` — CLI 옵션 전체
- `docs/languages/` — 언어별 가이드

**배포**: main 브랜치 push 시 GitHub Actions(`.github/workflows/pages.yml`)가 자동 배포
