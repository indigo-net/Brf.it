---
layout: default
title: Home
nav_exclude: true
---

# Brf.it

{: .fs-9 }

> Package your codebase for AI comprehension.
>
> `50 tokens → 8 tokens` — Same information, fewer tokens.
{: .fs-6 .fw-300 }

[Get Started](docs/getting-started){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 } [View on GitHub](https://github.com/indigo-net/Brf.it){: .btn .fs-5 .mb-4 .mb-md-0 }

---

## How It Works

Instead of feeding raw code to AI assistants:

### Before (50+ tokens)
{: .text-grey-dk-000 }

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

### After with brfit (8 tokens)
{: .text-green-100 }

```xml
<function>
  export async function fetchUser(
    id: string
  ): Promise<User>
</function>
```

---

## Features

| Feature | Description |
|---------|-------------|
| **Tree-sitter Based** | Accurate AST parsing for language structure analysis |
| **Multiple Formats** | XML and Markdown output support |
| **Token Counting** | Automatic output token calculation |
| **Gitignore Aware** | Automatically excludes unnecessary files |
| **Cross-Platform** | Linux, macOS, and Windows support |

---

## Quick Install

**macOS (Homebrew)**
{: .code-label }

```bash
brew install indigo-net/tap/brfit
```

**Linux / macOS (Script)**
{: .code-label }

```bash
curl -fsSL https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.sh | sh
```

**Windows (PowerShell)**
{: .code-label }

```powershell
irm https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.ps1 | iex
```

---

## Supported Languages

| Language | Extensions |
|----------|------------|
| Go | `.go` |
| TypeScript | `.ts`, `.tsx` |
| JavaScript | `.js`, `.jsx` |
| Python | `.py` |
| C | `.c` |
| C++ | `.cpp`, `.hpp`, `.h` |
| Java | `.java` |
| Rust | `.rs` |
| Swift | `.swift` |

See [Language Guides](docs/languages/) for details.

---

## See It In Action

- [SAMPLE.md](SAMPLE.md) — This project in Markdown format
- [SAMPLE.xml](SAMPLE.xml) — This project in XML format

Auto-generated on every commit.
