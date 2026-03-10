<p align="center">
  <img src="assets/logo.png" alt="Brf.it logo" width="128">
</p>

# Brf.it

🌐 [English](README.md) | [한국어](docs/ko/README.md) | [日本語](docs/ja/README.md) | [हिन्दी](docs/hi/README.md) | [Deutsch](docs/de/README.md)

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![Coverage](https://codecov.io/gh/indigo-net/Brf.it/branch/main/graph/badge.svg)](https://codecov.io/gh/indigo-net/Brf.it)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Package your codebase for AI comprehension.**
>
> `50 tokens → 8 tokens` — Same information, fewer tokens.

[Installation](#installation) · [Quick Start](#quick-start) · [Supported Languages](#supported-languages)

---

<br/>

## How It Works

Instead of feeding raw code to AI assistants:

<table>
<tr>
<th>Before (50+ tokens)</th>
<th>After with brfit (8 tokens)</th>
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

## Quick Start

### Installation

**macOS (Homebrew)**

```bash
brew install indigo-net/tap/brfit
```

**Linux / macOS (Script)**

```bash
curl -fsSL https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.sh | sh
```

**Windows (PowerShell)**

```powershell
irm https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.ps1 | iex
```

**From Source**

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

### First Run

```bash
brfit .                    # Analyze current directory
brfit . -f md              # Output in Markdown
brfit . -o briefing.xml    # Save to file
```

---

<br/>

## See It In Action

**[SAMPLE.md](SAMPLE.md)** | **[SAMPLE.xml](SAMPLE.xml)**

This project, packaged by brfit itself. Auto-generated on every commit.

---

<br/>

## Features

| Feature | Description |
|---------|-------------|
| Tree-sitter Based | Accurate AST parsing for language structure analysis |
| Multiple Formats | XML and Markdown output support |
| Token Counting | Automatic output token calculation |
| Gitignore Aware | Automatically excludes unnecessary files |
| Cross-Platform | Linux, macOS, and Windows support |

---

<br/>

## Supported Languages

| Language | Extensions | Documentation |
|----------|------------|---------------|
| Go | `.go` | [Go Guide](docs/languages/go.md) |
| TypeScript | `.ts`, `.tsx` | [TypeScript Guide](docs/languages/typescript.md) |
| JavaScript | `.js`, `.jsx` | [TypeScript Guide](docs/languages/typescript.md) |
| Python | `.py` | [Python Guide](docs/languages/python.md) |
| C | `.c` | [C Guide](docs/languages/c.md) |
| C++ | `.cpp`, `.hpp`, `.h` | [C++ Guide](docs/languages/cpp.md) |
| Java | `.java` | [Java Guide](docs/languages/java.md) |
| Rust | `.rs` | [Rust Guide](docs/languages/rust.md) |
| Swift | `.swift` | [Swift Guide](docs/languages/swift.md) |
| Kotlin | `.kt`, `.kts` | [Kotlin Guide](docs/languages/kotlin.md) |
| C# | `.cs` | [C# Guide](docs/languages/csharp.md) |
| Lua | `.lua` | [Lua Guide](docs/languages/lua.md) |
| PHP | `.php` | [PHP Guide](docs/languages/php.md) |
| Ruby | `.rb` | [Ruby Guide](docs/languages/ruby.md) |
| Scala | `.scala`, `.sc` | [Scala Guide](docs/languages/scala.md) |
| Elixir | `.ex`, `.exs` | [Elixir Guide](docs/languages/elixir.md) |

---

<br/>

## CLI Reference

```bash
brfit [path] [options]
```

### Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--format` | `-f` | Output format (`xml`, `md`) | `xml` |
| `--output` | `-o` | Output file path | stdout |
| `--include-body` | | Include function bodies | `false` |
| `--include-imports` | | Include import statements | `false` |
| `--ignore` | `-i` | Ignore file path | `.gitignore` |
| `--include-hidden` | | Include hidden files | `false` |
| `--no-tree` | | Skip directory tree | `false` |
| `--no-tokens` | | Disable token counting | `false` |
| `--max-size` | | Max file size (bytes) | `512000` |
| `--version` | `-v` | Show version | |

### Examples

```bash
# Copy to clipboard for AI assistants
brfit . | pbcopy              # macOS
brfit . | xclip               # Linux
brfit . | clip                # Windows

# Analyze project and save
brfit ./my-project -o briefing.xml

# Include function bodies (full code)
brfit . --include-body

# Skip directory tree output
brfit . --no-tree

# Include imports (verbatim)
brfit . --include-imports
```

---

<br/>

## License

MIT License — Use freely in personal and commercial projects.

See [LICENSE](LICENSE) for details.
