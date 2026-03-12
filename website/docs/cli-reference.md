---
sidebar_position: 2
title: CLI Reference
---

# CLI Reference

Complete reference for all Brf.it CLI options.

## Synopsis

```bash
brfit [path] [options]
```

## Arguments

| Argument | Description | Default |
|----------|-------------|---------|
| `path` | Directory or file to process | `.` (current directory) |

## Options

### Output Mode

| Option | Description | Default |
|--------|-------------|---------|
| `-m, --mode <mode>` | Output mode: `sig` (signatures only) | `sig` |
| `-f, --format <format>` | Output format: `xml` or `md` | `xml` |
| `-o, --output <file>` | Write output to file | stdout |

### Content Control

| Option | Description | Default |
|--------|-------------|---------|
| `--include-body` | Include function implementations | false |
| `--include-imports` | Include import/export statements | false |
| `--no-tree` | Skip directory tree in output | false |
| `--no-tokens` | Disable token counting | false |

### File Selection

| Option | Description | Default |
|--------|-------------|---------|
| `-i, --ignore <file>` | Custom ignore file path (can be specified multiple times) | `.gitignore` |
| `--include-hidden` | Include hidden files (starting with `.`) | false |
| `--max-size <bytes>` | Maximum file size to process | `512000` (500KB) |

### Help

| Option | Description |
|--------|-------------|
| `-h, --help` | Show help message |
| `-v, --version` | Show version |

## Examples

### Basic Usage

```bash
# Scan current directory, output XML to stdout
brfit .

# Scan specific directory
brfit ./src

# Scan single file
brfit ./src/main.go
```

### Markdown Output

```bash
# Output as Markdown
brfit . -f md

# Save to file
brfit . -f md -o context.md
```

### Include Imports

```bash
# Include import/export statements for dependency context
brfit . --include-imports -f md
```

### Custom Ignore File

```bash
# Use custom ignore rules
brfit . -i .brfitignore
```

### Large File Handling

```bash
# Increase max file size limit (1MB)
brfit . --max-size 1048576

# Include hidden files
brfit . --include-hidden
```

### Minimal Output

```bash
# Skip tree and token count for minimal output
brfit . --no-tree --no-tokens
```

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | Error (file not found, parse error, etc.) |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `BRFIT_LOG_LEVEL` | Set log level: `debug`, `info`, `warn`, `error` |

## Configuration File

Brf.it respects `.gitignore` by default. Create a `.brfitignore` file for custom ignore rules:

```gitignore
# Example .brfitignore
node_modules/
dist/
*.min.js
*.test.ts
vendor/
```

## Integration Examples

### Makefile

```makefile
context:
	brfit . -f md --include-imports -o AI_CONTEXT.md
```

### package.json

```json
{
  "scripts": {
    "context": "brfit . -f md --include-imports -o AI_CONTEXT.md"
  }
}
```

### Git Hook (pre-commit)

```bash
#!/bin/bash
# Update AI context before commits
brfit . -f md --include-imports -o AI_CONTEXT.md
git add AI_CONTEXT.md
```
