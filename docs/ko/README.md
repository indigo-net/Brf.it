<p align="center">
  <img src="../../assets/logo.png" alt="Brf.it logo" width="128">
</p>

# Brf.it

🌐 [English](../../README.md) | [한국어](README.md) | [日本語](../ja/README.md) | [हिन्दी](../hi/README.md) | [Deutsch](../de/README.md)

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **코드베이스를 AI가 이해하기 쉬운 형태로 패키징**
>
> `50 토큰 → 8 토큰` — 같은 정보, 더 적은 토큰.

[설치](#설치) · [빠른 시작](#빠른-시작) · [지원 언어](#지원-언어)

---

<br/>

## 동작 방식

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
<function>
  export async function fetchUser(
    id: string
  ): Promise<User>
</function>
```

</td>
</tr>
</table>

---

<br/>

## 빠른 시작

### 설치

**macOS (Homebrew)**

```bash
brew install indigo-net/tap/brfit
```

**Linux / macOS (스크립트)**

```bash
curl -fsSL https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.sh | sh
```

**Windows (PowerShell)**

```powershell
irm https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.ps1 | iex
```

**소스에서 빌드**

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

### 첫 실행

```bash
brfit .                    # 현재 디렉토리 분석
brfit . -f md              # Markdown 출력
brfit . -o briefing.xml    # 파일로 저장
```

---

<br/>

## 실제 사용 예시

**[SAMPLE.md](SAMPLE.md)** | **[SAMPLE.xml](SAMPLE.xml)**

brfit으로 패키징한 이 프로젝트 자체입니다. 커밋마다 자동 생성됩니다.

---

<br/>

## 주요 기능

| 기능 | 설명 |
|------|------|
| Tree-sitter 기반 | 정확한 AST 파싱으로 언어 구조 분석 |
| 다중 포맷 | XML, Markdown 출력 지원 |
| 토큰 카운팅 | 출력 토큰 수 자동 계산 |
| gitignore 인식 | 불필요한 파일 자동 제외 |
| 크로스 플랫폼 | Linux, macOS, Windows 지원 |

---

<br/>

## 지원 언어

| 언어 | 확장자 | 문서 |
|------|--------|------|
| Go | `.go` | [Go 가이드](languages/go.md) |
| TypeScript | `.ts`, `.tsx` | [TypeScript 가이드](languages/typescript.md) |
| JavaScript | `.js`, `.jsx` | [TypeScript 가이드](languages/typescript.md) |
| Python | `.py` | [Python 가이드](languages/python.md) |
| C | `.c`, `.h` | [C 가이드](languages/c.md) |
| Java | `.java` | [Java 가이드](languages/java.md) |
| Rust | `.rs` | [Rust 가이드](languages/rust.md) |
| Swift | `.swift` | [Swift 가이드](languages/swift.md) |
| Kotlin | `.kt`, `.kts` | [Kotlin 가이드](languages/kotlin.md) |
| C# | `.cs` | [C# 가이드](languages/csharp.md) |
| Lua | `.lua` | [Lua 가이드](languages/lua.md) |
| PHP | `.php` | [PHP 가이드](languages/php.md) |
| Ruby | `.rb` | [Ruby 가이드](languages/ruby.md) |
| Scala | `.scala`, `.sc` | [Scala 가이드](languages/scala.md) |
| Elixir | `.ex`, `.exs` | [Elixir 가이드](languages/elixir.md) |
| SQL | `.sql` | [SQL 가이드](languages/sql.md) |
| YAML | `.yaml`, `.yml` | [YAML 가이드](languages/yaml.md) |
| TOML | `.toml` | [TOML 가이드](languages/toml.md) |

---

<br/>

## CLI 레퍼런스

```bash
brfit [경로] [옵션]
```

### 옵션

| 옵션 | 단축 | 설명 | 기본값 |
|------|------|------|--------|
| `--format` | `-f` | 출력 형식 (`xml`, `md`) | `xml` |
| `--output` | `-o` | 출력 파일 경로 | stdout |
| `--include-body` | | 함수 본문 포함 | `false` |
| `--include-imports` | | import 문 포함 | `false` |
| `--include-private` | | 비공개/unexported 심볼 포함 | `false` |
| `--ignore` | `-i` | ignore 파일 경로 (여러 번 지정 가능) | `.gitignore` |
| `--include` | | 포함할 glob 패턴 (여러 번 지정 가능) | |
| `--exclude` | | 제외할 glob 패턴 (여러 번 지정 가능) | |
| `--include-hidden` | | 숨김 파일 포함 | `false` |
| `--no-tree` | | 디렉토리 트리 생략 | `false` |
| `--no-tokens` | | 토큰 수 계산 비활성화 | `false` |
| `--max-size` | | 최대 파일 크기 (바이트) | `512000` |
| `--changed` | | Git 변경 파일만 스캔 | `false` |
| `--since` | | 특정 커밋/태그 이후 변경된 파일만 스캔 | |
| `--token-tree` | | 파일별 토큰 수 트리 출력 | `false` |
| `--version` | `-v` | 버전 표시 | |

### 예제

```bash
# AI 어시스턴트에 전달 (클립보드 복사)
brfit . | pbcopy              # macOS
brfit . | xclip               # Linux
brfit . | clip                # Windows

# 프로젝트 분석 후 파일 저장
brfit ./my-project -o briefing.xml

# 함수 본문 포함 (전체 코드)
brfit . --include-body

# 디렉토리 트리 출력 생략
brfit . --no-tree

# import 포함 (원본 그대로)
brfit . --include-imports
```

---

<br/>

## 라이선스

MIT 라이선스 — 개인 및 상업 프로젝트에서 자유롭게 사용 가능합니다.

자세한 내용은 [LICENSE](LICENSE)를 참조하세요.
