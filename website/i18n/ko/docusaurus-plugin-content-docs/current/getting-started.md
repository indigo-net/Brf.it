---
sidebar_position: 1
title: 시작하기
slug: /
---

# 시작하기

Brf.it는 코드베이스에서 함수 시그니처를 추출하여 AI 어시스턴트에 필요한 컨텍스트를 효율적으로 제공합니다.

## Brf.it란?

Brf.it는 다음 기능을 제공하는 CLI 도구입니다:
- 코드베이스에서 소스 파일 스캔
- Tree-sitter를 사용하여 함수 시그니처와 문서 추출
- AI 친화적인 형식(XML 또는 Markdown)으로 출력
- AI 비용 관리를 위한 토큰 수 계산

## 설치

### macOS (Homebrew)

```bash
brew tap indigo-net/tap
brew install brfit
```

### Linux

```bash
# 최신 릴리스 다운로드
curl -sSL https://github.com/indigo-net/Brf.it/releases/latest/download/brfit-linux-amd64 -o brfit
chmod +x brfit
sudo mv brfit /usr/local/bin/
```

### Windows (Scoop)

```bash
scoop bucket add indigo-net https://github.com/indigo-net/scoop-bucket
scoop install brfit
```

### 소스에서 빌드

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

## 빠른 시작

프로젝트 디렉토리로 이동하여 실행:

```bash
brfit .
```

이 명령은:
1. 현재 디렉토리의 모든 지원 파일을 스캔
2. 함수 시그니처만 추출 (구현 제외)
3. XML 형식으로 stdout에 출력

### 주요 옵션

```bash
# Markdown으로 출력
brfit . -f md

# import/export 문 포함
brfit . --include-imports

# 파일로 저장
brfit . -o context.md

# 토큰 수 계산 건너뛰기
brfit . --no-tokens
```

## 사용 사례

### Claude Code와 함께

```bash
# 컨텍스트 생성 후 클립보드로 복사
brfit . -f md --include-imports | pbcopy
```

그런 다음 Claude Code에 프롬프트와 함께 붙여넣기.

### Cursor와 함께

```bash
# @file 참조용 컨텍스트 파일 생성
brfit ./src -f md -o .cursor/context.md
```

프롬프트에서 `@context.md`로 참조.

### GitHub Copilot과 함께

프로젝트에 컨텍스트 파일 추가:

```bash
brfit . -f md -o AI_CONTEXT.md
```

Copilot이 더 나은 제안을 제공합니다.

## 출력 형식

### XML (기본값)

```xml
<brfit>
  <metadata>
    <tree>src/
├── main.go
└── handler.go</tree>
    <tokens>245</tokens>
  </metadata>
  <files>
    <file path="src/main.go" language="go">
      <function>func main()</function>
      <doc>Entry point for the application</doc>
    </file>
    <file path="src/handler.go" language="go">
      <function>func HandleRequest(ctx context.Context, req Request) (*Response, error)</function>
      <function>func validateInput(req Request) error</function>
    </file>
  </files>
</brfit>
```

### Markdown

```markdown
# Codebase Summary

**Tokens:** 245

## src/main.go (Go)

### Functions

- `func main()` - Entry point for the application

## src/handler.go (Go)

### Functions

- `func HandleRequest(ctx context.Context, req Request) (*Response, error)`
- `func validateInput(req Request) error`
```

## 다음 단계

- [CLI 참조](/docs/cli-reference) - 모든 옵션
- [지원 언어](/docs/languages/) - 언어별 상세 정보
