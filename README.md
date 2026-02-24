# Brf.it

[![Release](https://img.shields.io/github/v/release/indigo-net/Brf.it)](https://github.com/indigo-net/Brf.it/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/indigo-net/Brf.it)](https://goreportcard.com/report/github.com/indigo-net/Brf.it)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Brief your code for AI assistants.**

Brf.it extracts function signatures from your codebase, removing implementation details to dramatically reduce token usage while preserving the essential information AI needs.

---

## What It Does

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

## Installation

### macOS (Homebrew)

```bash
brew install indigo-net/tap/brfit
```

### From Release

Download the latest binary from [Releases](https://github.com/indigo-net/Brf.it/releases).

### From Source

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

---

## Usage

```bash
brfit [path] [options]
```

### Quick Examples

```bash
# Extract signatures from current directory
brfit .

# Output in Markdown format
brfit . -f md

# Save to file
brfit . -o output.xml

# Include function bodies (full code)
brfit . --include-body

# Skip directory tree
brfit . --no-tree
```

### CLI Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--format` | `-f` | Output format (`xml`, `md`) | `xml` |
| `--output` | `-o` | Output file path | stdout |
| `--include-body` | | Include function bodies | `false` |
| `--ignore` | `-i` | Ignore file path | `.gitignore` |
| `--include-hidden` | | Include hidden files | `false` |
| `--no-tree` | | Skip directory tree | `false` |
| `--no-tokens` | | Disable token counting | `false` |
| `--max-size` | | Max file size (bytes) | `512000` |
| `--version` | `-v` | Show version | |

---

## Supported Languages

| Language | Extensions |
|----------|------------|
| Go | `.go` |
| TypeScript | `.ts`, `.tsx` |
| JavaScript | `.js`, `.jsx` |

---

## Output Examples

### XML (default)

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

## License

MIT
