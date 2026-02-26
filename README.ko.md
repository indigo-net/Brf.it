# Brf.it

🌐 [English](README.md) | [한국어](README.ko.md) | [日本語](README.ja.md) | [हिन्दी](README.hi.md) | [Deutsch](README.de.md)

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**AI 어시스턴트를 위한 코드 브리핑 도구**

Brf.it은 코드베이스에서 함수 시그니처를 추출하여, 구현 세부사항을 제거하고 AI가 필요로 하는 핵심 정보만 남깁니다. 토큰 사용량을 획기적으로 줄일 수 있습니다.

---

## 주요 기능

AI 어시스턴트에 원본 코드를 전달하는 대신:

<table>
<tr>
<th>Before (50+ 토큰)</th>
<th>After with brfit (8 토큰)</th>
</tr>
<tr>
<td>

```typescript
export async function fetchUser(
  id: string
): Promise<User> {
  const response = await fetch(
    `${API_URL}/users/${id}`
  );
  if (!response.ok) {
    throw new Error('User not found');
  }
  const data = await response.json();
  return {
    id: data.id,
    name: data.name,
    email: data.email,
    createdAt: new Date(data.created_at)
  };
}
```

</td>
<td>

```xml
<signature>
  export async function fetchUser(
    id: string
  ): Promise<User>
</signature>
```

</td>
</tr>
</table>

---

## 설치

### macOS (Homebrew)

```bash
brew install indigo-net/tap/brfit
```

### 릴리즈에서 다운로드

[Releases](https://github.com/indigo-net/Brf.it/releases)에서 최신 바이너리를 다운로드하세요.

### 소스에서 빌드

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

---

## 사용법

```bash
brfit [경로] [옵션]
```

### 빠른 예제

```bash
# 현재 디렉토리에서 시그니처 추출
brfit .

# Markdown 형식으로 출력
brfit . -f md

# 파일로 저장
brfit . -o output.xml

# 함수 본문 포함 (전체 코드)
brfit . --include-body

# 디렉토리 트리 생략
brfit . --no-tree
```

### CLI 옵션

| 옵션 | 단축 | 설명 | 기본값 |
|------|------|------|--------|
| `--format` | `-f` | 출력 형식 (`xml`, `md`) | `xml` |
| `--output` | `-o` | 출력 파일 경로 | stdout |
| `--include-body` | | 함수 본문 포함 | `false` |
| `--ignore` | `-i` | ignore 파일 경로 | `.gitignore` |
| `--include-hidden` | | 숨김 파일 포함 | `false` |
| `--no-tree` | | 디렉토리 트리 생략 | `false` |
| `--no-tokens` | | 토큰 수 계산 비활성화 | `false` |
| `--max-size` | | 최대 파일 크기 (바이트) | `512000` |
| `--version` | `-v` | 버전 표시 | |

---

## 지원 언어

| 언어 | 확장자 | 문서 |
|------|--------|------|
| Go | `.go` | [Go 가이드](docs/languages/go.ko.md) |
| TypeScript | `.ts`, `.tsx` | [TypeScript 가이드](docs/languages/typescript.ko.md) |
| JavaScript | `.js`, `.jsx` | [TypeScript 가이드](docs/languages/typescript.ko.md) |
| Python | `.py` | [Python 가이드](docs/languages/python.ko.md) |
| C | `.c`, `.h` | [C 가이드](docs/languages/c.ko.md) |

---

## 출력 예제

### XML (기본)

```xml
<?xml version="1.0" encoding="UTF-8"?>
<brfit>
  <metadata>
    <tree>pkg/
└── scanner/
    └── scanner.go</tree>
  </metadata>
  <files>
    <file path="pkg/scanner/scanner.go" language="go">
      <signature>func Scan(root string) (*Result, error)</signature>
      <doc>Scan recursively scans the directory.</doc>
    </file>
  </files>
</brfit>
```

### Markdown

```markdown
# Brf.it Output

## Directory Tree

pkg/
└── scanner/
    └── scanner.go

## Files

### pkg/scanner/scanner.go

\`\`\`go
func Scan(root string) (*Result, error)
\`\`\`

> Scan recursively scans the directory.
```

---

## 라이선스

MIT
