# Brf.it

_"Don't just feed raw code. Brief it for your AI."_

---

## (한국어)

### Brf.it이란?

Brf.it은 AI 코딩 어시스턴트에게 코드를 더 효율적으로 전달하기 위한 CLI 도구입니다.

프로젝트 코드베이스에서 함수 시그니처와 문서만 추출하여, AI가 핵심 내용을 빠르게 파악할 수 있는 형태로 변환합니다. 불필요한 구현 디테일을 제거해 토큰 사용량을 크게 줄일 수 있습니다.

### 왜 필요한가요?

AI에게 큰 프로젝트의 코드를 전달할 때, 모든 파일을 그대로 복사하면:

- 컨텍스트가 너무 길어집니다
- 토큰이 낭비됩니다
- AI가 중요한 부분을 놓칠 수 있습니다

Brf.it은 이 문제를 해결합니다.

### 기능

- **지능형 코드 압축**: 함수 시그니처와 주석만 추출 (Tree-sitter 기반)
- **다양한 출력 포맷**: XML, Markdown 선택 가능
- **구조화된 메타데이터**: 디렉토리 트리와 심볼 리스트 포함
- **토큰 카운트**: 출력물의 토큰 수 자동 계산
- **간편한 사용**: `brfit <path>` 한 줄로 실행

### 지원 언어

- Go (`.go`)
- TypeScript (`.ts`, `.tsx`)
- JavaScript (`.js`, `.jsx`)

### 설치

#### npm (권장)

```bash
# npx로 바로 실행
npx brfit .

# 또는 전역 설치
npm install -g brfit
brfit .
```

#### 소스에서 빌드

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o bin/brfit ./cmd/brfit
```

### 사용법

```bash
brfit [path] [options]
```

#### 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `-m, --mode` | 출력 모드 | `sig` |
| `-f, --format` | 출력 포맷 (`xml`, `md`, `markdown`) | `xml` |
| `-o, --output` | 출력 파일 경로 | stdout |
| `-i, --ignore` | 무시 파일 (gitignore 패턴) | `.gitignore` |
| `--include-hidden` | 숨김 파일 포함 | `false` |
| `--no-tree` | 디렉토리 트리 생략 | `false` |
| `--no-tokens` | 토큰 수 계산 비활성화 | `false` |
| `--max-size` | 최대 파일 크기 (bytes) | `512000` (500KB) |
| `-v, --version` | 버전 정보 표시 | |

### 사용 예시

```bash
# 기본 사용법 (XML, stdout)
npx brfit .

# Markdown 포맷으로 출력
npx brfit . -f md

# 결과를 파일로 저장
npx brfit . -o output.xml

# 하위 디렉토리에 저장 (자동 생성)
npx brfit . -o build/output/result.xml

# 토큰 카운트 없이
npx brfit . --no-tokens

# 트리 없이 시그니처만
npx brfit . --no-tree

# 숨김 파일 포함
npx brfit . --include-hidden

# 커스텀 ignore 파일
npx brfit . -i .brfitignore

# 최대 파일 크기 설정
npx brfit . --max-size 1000000

# 버전 확인
npx brfit --version

# 도움말
npx brfit --help
```

### 출력 예시

#### XML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<brfit>
  <metadata>
    <tree>pkg/
└── scanner/
    └── scanner.go</tree>
    <symbols>
      - func ScanDirectory(root string, opts *ScanOptions) ([]FileEntry, error)
    </symbols>
  </metadata>
  <files>
    <file path="pkg/scanner/scanner.go" language="go">
      <signature>func ScanDirectory(root string, opts *ScanOptions) ([]FileEntry, error)</signature>
      <doc>ScanDirectory recursively scans the directory.</doc>
    </file>
  </files>
</brfit>
```

#### Markdown

```markdown
# Brf.it Output

## Directory Tree

```
pkg/
└── scanner/
    └── scanner.go
```

## Symbols

- `func ScanDirectory(root string, opts *ScanOptions) ([]FileEntry, error)`

---

## Files

### pkg/scanner/scanner.go

```go
func ScanDirectory(root string, opts *ScanOptions) ([]FileEntry, error)
```

> ScanDirectory recursively scans the directory.

---
```

---

## (English)

### What is Brf.it?

Brf.it is a CLI tool designed to deliver code to AI coding assistants more efficiently.

It extracts only function signatures and documentation from your codebase, transforming them into a format that AI can quickly understand. By removing unnecessary implementation details, it significantly reduces token usage.

### Why do you need it?

When sharing a large project's code with AI, copying all files directly leads to:

- Overwhelmingly long context
- Wasted tokens
- AI potentially missing important parts

Brf.it solves this problem.

### Features

- **Smart Code Compression**: Extract only function signatures and comments (Tree-sitter based)
- **Multiple Output Formats**: Choose between XML and Markdown
- **Structured Metadata**: Includes directory tree and symbol list
- **Token Counting**: Automatic token count calculation
- **Simple Usage**: Run with just `brfit <path>`

### Supported Languages

- Go (`.go`)
- TypeScript (`.ts`, `.tsx`)
- JavaScript (`.js`, `.jsx`)

### Installation

#### npm (Recommended)

```bash
# Run directly with npx
npx brfit .

# Or install globally
npm install -g brfit
brfit .
```

#### Build from Source

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o bin/brfit ./cmd/brfit
```

### Usage

```bash
brfit [path] [options]
```

#### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-m, --mode` | Output mode | `sig` |
| `-f, --format` | Output format (`xml`, `md`, `markdown`) | `xml` |
| `-o, --output` | Output file path | stdout |
| `-i, --ignore` | Ignore file (gitignore patterns) | `.gitignore` |
| `--include-hidden` | Include hidden files | `false` |
| `--no-tree` | Skip directory tree | `false` |
| `--no-tokens` | Disable token count | `false` |
| `--max-size` | Maximum file size (bytes) | `512000` (500KB) |
| `-v, --version` | Show version information | |

### Examples

```bash
# Basic usage (XML, stdout)
npx brfit .

# Output in Markdown format
npx brfit . -f md

# Save output to file
npx brfit . -o output.xml

# Save to subdirectory (auto-created)
npx brfit . -o build/output/result.xml

# Without token count
npx brfit . --no-tokens

# Skip directory tree
npx brfit . --no-tree

# Include hidden files
npx brfit . --include-hidden

# Custom ignore file
npx brfit . -i .brfitignore

# Set max file size
npx brfit . --max-size 1000000

# Show version
npx brfit --version

# Show help
npx brfit --help
```

### Output Examples

#### XML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<brfit>
  <metadata>
    <tree>pkg/
└── scanner/
    └── scanner.go</tree>
    <symbols>
      - func ScanDirectory(root string, opts *ScanOptions) ([]FileEntry, error)
    </symbols>
  </metadata>
  <files>
    <file path="pkg/scanner/scanner.go" language="go">
      <signature>func ScanDirectory(root string, opts *ScanOptions) ([]FileEntry, error)</signature>
      <doc>ScanDirectory recursively scans the directory.</doc>
    </file>
  </files>
</brfit>
```

#### Markdown

```markdown
# Brf.it Output

## Directory Tree

```
pkg/
└── scanner/
    └── scanner.go
```

## Symbols

- `func ScanDirectory(root string, opts *ScanOptions) ([]FileEntry, error)`

---

## Files

### pkg/scanner/scanner.go

```go
func ScanDirectory(root string, opts *ScanOptions) ([]FileEntry, error)
```

> ScanDirectory recursively scans the directory.

---
```

## License

MIT
