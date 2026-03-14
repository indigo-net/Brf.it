---
layout: default
title: Getting Started
nav_order: 2
---

# Getting Started

## Installation

### macOS (Homebrew)

```bash
brew install indigo-net/tap/brfit
```

### Linux / macOS (Script)

```bash
curl -fsSL https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.sh | sh
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/indigo-net/Brf.it/main/install.ps1 | iex
```

### From Source

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

## Verify Installation

```bash
brfit --version
```

---

## 5-Minute Quick Start

### Step 1: Analyze Your Project

Navigate to any project directory and run brfit:

```bash
cd ~/my-project
brfit .
```

This outputs an XML document containing all function signatures, types, and documentation from your codebase — stripping away implementation details that consume tokens but add little structural context.

### Step 2: Copy to Clipboard

**macOS**
{: .code-label }

```bash
brfit . | pbcopy
```

**Linux (with xclip)**
{: .code-label }

```bash
brfit . | xclip -selection clipboard
```

**Windows**
{: .code-label }

```bash
brfit . | clip
```

Now paste the result into your AI assistant (Claude, ChatGPT, etc.) along with your question.

### Step 3: Choose the Right Level of Detail

brfit offers multiple levels of information density:

```bash
# Signatures only (default) — smallest output, best for architecture questions
brfit .

# Signatures + imports — see dependencies between modules
brfit . --include-imports

# Full function bodies — when AI needs implementation details
brfit . --include-body
```

### Step 4: Focus on What Matters

```bash
# Only analyze Go files in the pkg directory
brfit . --include "pkg/**/*.go"

# Only files changed since your last release
brfit . --since v1.0.0

# Only uncommitted changes
brfit . --changed
```

### Step 5: Save and Reuse

```bash
# Save as XML for tools
brfit . -o context.xml

# Save as Markdown for documentation
brfit . -f md -o context.md
```

---

## Real-World Use Cases

### Using with Claude Code

When working on a large project with Claude Code, brfit helps you provide full project context within the token limit:

```bash
# Generate a project briefing and pipe it to your prompt
brfit . -f md --include-imports > /tmp/context.md
```

Paste the content into Claude with questions like:
- *"Based on this project structure, where should I add a new REST endpoint?"*
- *"What modules depend on the `scanner` package?"*

For Claude Desktop, use the MCP server for automatic integration:

```json
{
  "mcpServers": {
    "brfit": {
      "command": "brfit-mcp",
      "args": ["--root", "/path/to/project"]
    }
  }
}
```

### Using with ChatGPT / Other AI Assistants

```bash
# Copy project signatures to clipboard
brfit . | pbcopy

# Then paste into ChatGPT with your question:
# "Here is my project structure. How can I refactor the parser module?"
```

For larger projects, use `--include` to focus on the relevant area:

```bash
# Only the authentication module
brfit . --include "internal/auth/**" | pbcopy
```

### Code Review Preparation

Generate a focused summary of what changed for review context:

```bash
# What changed since the last tag
brfit . --changed -f md

# What changed since a specific commit
brfit . --since HEAD~10 -f md -o review-context.md
```

This gives reviewers a quick structural overview of the affected code without reading every line.

### Understanding Legacy Code

When onboarding to an unfamiliar codebase:

```bash
# Get a bird's-eye view of the entire project
brfit . -f md -o overview.md

# Then zoom into specific areas
brfit . --include "src/core/**" --include-imports -f md

# Check function call relationships
brfit . --call-graph --include "src/core/**"
```

Paste the output into an AI assistant and ask:
- *"Explain the architecture of this project"*
- *"What are the main entry points?"*
- *"How does data flow from input to output?"*

### Security Auditing

brfit automatically detects and redacts secrets when extracting code:

```bash
# Extract with bodies — secrets are automatically redacted
brfit . --include-body

# Disable redaction if needed (not recommended for sharing)
brfit . --include-body --no-security-check
```

Detected patterns include AWS keys, GitHub tokens, API keys, passwords, bearer tokens, private keys, Slack tokens, and more (12 types total).

### Token Budget Estimation

Use `--token-tree` to understand which parts of your codebase consume the most tokens:

```bash
brfit . --token-tree
```

This helps you decide which directories to include or exclude when working within AI context limits.

---

## FAQ

### How does brfit reduce token usage?

brfit uses [Tree-sitter](https://tree-sitter.github.io/) to parse your code into an AST and extract only the structural elements: function signatures, type definitions, interface declarations, and documentation comments. Implementation details (function bodies, inline comments, variable assignments) are stripped by default, typically reducing token count by 80-90%.

### Which languages are supported?

brfit supports 19 languages: Go, TypeScript, JavaScript, Python, C, C++, Java, Rust, Swift, Kotlin, C#, Lua, PHP, Ruby, Scala, Elixir, SQL, YAML, and TOML. See [Supported Languages](languages/) for per-language details.

### Does brfit lose important information?

By default, brfit extracts signatures and documentation — this is enough for AI to understand project architecture, dependencies, and API surfaces. For tasks that require implementation details (bug finding, detailed refactoring), use `--include-body` to include function bodies.

### Can I use brfit in CI/CD?

Yes. brfit is a single binary with no runtime dependencies. Common CI use cases:

```bash
# Generate context as a CI artifact
brfit . -o context.xml

# Check for accidentally committed secrets
brfit . --include-body 2>security-warnings.txt
```

### What's the difference between XML and Markdown output?

**XML** (`-f xml`, default) — structured, machine-readable. Best for tool integration and programmatic use.

**Markdown** (`-f md`) — human-readable. Best for documentation, code review context, and pasting into AI chat interfaces.

### How do I handle large projects?

For projects that exceed AI context limits even with brfit:

```bash
# Focus on specific directories
brfit . --include "src/api/**"

# Exclude test files and vendor code
brfit . --exclude "**/*_test.go" --exclude "vendor/**"

# Only recently changed files
brfit . --changed
```

### What is the MCP server?

`brfit-mcp` is a [Model Context Protocol](https://modelcontextprotocol.io/) server that lets AI agents call brfit directly without manual copy-paste. It provides `summarize_project` and `summarize_file` tools over stdio JSON-RPC. See [CLI Reference](cli-reference#mcp-server-brfit-mcp) for setup instructions.

---

## Next Steps

- [CLI Reference](cli-reference) — Full command-line options
- [Language Guides](languages/) — Language-specific documentation
