---
sidebar_position: 1
title: Getting Started
slug: /
---

# Getting Started

Brf.it extracts function signatures from your codebase, giving AI assistants the context they need without the bloat.

## What is Brf.it?

Brf.it is a CLI tool that:
- Scans your codebase for source files
- Extracts function signatures and documentation using Tree-sitter
- Outputs compact, AI-friendly formats (XML or Markdown)
- Counts tokens to help manage AI costs

## Installation

### macOS (Homebrew)

```bash
brew tap indigo-net/tap
brew install brfit
```

### Linux

```bash
# Download latest release
curl -sSL https://github.com/indigo-net/Brf.it/releases/latest/download/brfit-linux-amd64 -o brfit
chmod +x brfit
sudo mv brfit /usr/local/bin/
```

### Windows (Scoop)

```bash
scoop bucket add indigo-net https://github.com/indigo-net/scoop-bucket
scoop install brfit
```

### From Source

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build -o brfit ./cmd/brfit
```

## Quick Start

Navigate to your project and run:

```bash
brfit .
```

This will:
1. Scan all supported files in the current directory
2. Extract function signatures (not implementations)
3. Output XML to stdout

### Common Options

```bash
# Output as Markdown
brfit . -f md

# Include import/export statements
brfit . --include-imports

# Save to file
brfit . -o context.md

# Skip token counting
brfit . --no-tokens
```

## Use Cases

### With Claude Code

```bash
# Generate context and pipe to clipboard
brfit . -f md --include-imports | pbcopy
```

Then paste into Claude Code with your prompt.

### With Cursor

```bash
# Generate context file for @file references
brfit ./src -f md -o .cursor/context.md
```

Reference with `@context.md` in your prompts.

### With GitHub Copilot

Add a context file to your project:

```bash
brfit . -f md -o AI_CONTEXT.md
```

Copilot will use this for better suggestions.

## Output Formats

### XML (default)

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

## Next Steps

- [CLI Reference](/docs/cli-reference) - All available options
- [Supported Languages](/docs/languages/) - Language-specific details
