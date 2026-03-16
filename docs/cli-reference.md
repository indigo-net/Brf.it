---
layout: default
title: CLI Reference
nav_order: 3
---

# CLI Reference

## Usage

```bash
brfit [path] [options]
```

## Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--format` | `-f` | Output format (`xml`, `md`) | `xml` |
| `--output` | `-o` | Output file path | stdout |
| `--include-body` | | Include function bodies | `false` |
| `--include-imports` | | Include import statements | `false` |
| `--include-private` | | Include non-exported/private symbols | `false` |
| `--no-std-imports` | | Exclude stdlib imports | `false` |
| `--ignore` | `-i` | Ignore file path (can be specified multiple times) | `.gitignore` |
| `--include` | | Glob pattern(s) to include (can be specified multiple times) | |
| `--exclude` | | Glob pattern(s) to exclude (can be specified multiple times) | |
| `--include-hidden` | | Include hidden files | `false` |
| `--no-tree` | | Skip directory tree | `false` |
| `--no-tokens` | | Disable token counting | `false` |
| `--max-size` | | Max file size (bytes) | `512000` |
| `--changed` | | Only scan git-modified files (tracked + untracked) | `false` |
| `--since` | | Only scan files changed since commit/tag (e.g., `v1.0.0`, `HEAD~5`) | |
| `--token-tree` | | Show per-file token count tree with directory totals | `false` |
| `--security-check` / `--no-security-check` | | Detect and redact secrets in extracted code | `true` |
| `--call-graph` | | Extract function/method call relationships per file | `false` |
| `--strict` | | Exit with code 1 if any file has parsing errors (CI quality gate) | `false` |
| `--version` | `-v` | Show version | |
| `--help` | `-h` | Show help | |

## Output Formats

### XML (`-f xml`)

Default format. Structured and machine-readable.

```xml
<brfit>
  <metadata>
    <tree>...</tree>
    <tokens>1234</tokens>
  </metadata>
  <files>
    <file path="src/main.go" language="go">
      <function>func main()</function>
      <doc>Main entry point</doc>
    </file>
  </files>
</brfit>
```

### Markdown (`-f md`)

Human-readable format, great for documentation.

````markdown
# Project: myproject

## Directory Structure
```
src/
├── main.go
└── utils.go
```

## Files

### src/main.go

```go
func main()
```
Main entry point
````

## Examples

### Basic Usage

```bash
# Analyze current directory
brfit .

# Analyze specific directory
brfit ./src

# Output in Markdown
brfit . -f md
```

### Saving to File

```bash
# Save XML output
brfit . -o briefing.xml

# Save Markdown output
brfit . -f md -o briefing.md
```

### Including More Content

```bash
# Include function bodies
brfit . --include-body

# Include imports
brfit . --include-imports

# Include imports, exclude stdlib
brfit . --include-imports --no-std-imports
```

### Performance Options

```bash
# Skip directory tree for faster output
brfit . --no-tree

# Disable token counting
brfit . --no-tokens

# Increase max file size limit
brfit . --max-size 1048576  # 1MB
```

### Token Tree

```bash
# Show per-file token counts in tree format
brfit . --token-tree

# Combine with exclude patterns
brfit . --token-tree --exclude "vendor/**"
```

### Git Change Detection

```bash
# Only scan files changed in git working tree
brfit . --changed

# Only scan files changed since a tag
brfit . --since v1.0.0

# Combine with format options
brfit . --changed -f md -o changes.md

# Only changes since 5 commits ago
brfit . --since HEAD~5
```

### Security Check

```bash
# Security check is enabled by default
# Detected secrets are replaced with [REDACTED] and warnings are printed to stderr
brfit . --include-body

# Disable security check
brfit . --include-body --no-security-check
```

**Detected patterns (12 types):**
- AWS Access Key ID (`AKIA...`)
- AWS Secret Access Key
- GitHub Token (`ghp_`, `gho_`, `ghs_`, `ghr_`, `github_pat_`)
- Generic API Key patterns (`api_key`, `apikey`, `api-key`)
- Password patterns (`password`, `passwd`, `pwd`)
- Bearer tokens (`Bearer ...`)
- Private keys (`-----BEGIN ... PRIVATE KEY-----`)
- Slack tokens (`xoxb-`, `xoxp-`, `xoxo-`, `xapp-`)
- Google API Key (`AIza...`)
- Heroku API Key

### Call Graph

```bash
# Extract function call relationships
brfit . --call-graph

# Combine with format options
brfit . --call-graph -f md

# Call graph for specific files only
brfit . --call-graph --include "pkg/**/*.go"
```

Output shows caller-callee references per file:

```
## src/main.go

| Caller | Callee | Line |
|--------|--------|------|
| main | parseArgs | 15 |
| main | runServer | 23 |
| runServer | handleRequest | 45 |
```

**Supported languages:** Go, TypeScript/JavaScript, Python, Java, Rust, C

### Custom Ignore File

```bash
# Use custom ignore file
brfit . -i .brfitignore

# Include hidden files (normally excluded)
brfit . --include-hidden
```

## Integration with AI Tools

### Claude / ChatGPT

```bash
# Copy to clipboard and paste into AI chat
brfit . | pbcopy    # macOS
brfit . | xclip     # Linux
brfit . | clip      # Windows
```

### Piping to Files

```bash
# Generate XML for tools that accept file input
brfit . -o context.xml

# Use in scripts
brfit . --no-tokens --no-tree > signatures.xml
```

---

## MCP Server (`brfit-mcp`)

`brfit-mcp` is a standalone [Model Context Protocol](https://modelcontextprotocol.io/) server binary that exposes brfit's code analysis as tools for AI agents. It communicates over stdio using JSON-RPC.

### Usage

```bash
brfit-mcp --root /path/to/project
```

### Tools

#### `summarize_project`

Extract signatures from a project directory.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `path` | Project directory path | Server `--root` value |
| `format` | Output format (`xml`, `md`) | `xml` |
| `include_body` | Include function bodies | `false` |
| `include_imports` | Include import statements | `false` |
| `call_graph` | Extract function call relationships | `false` |

#### `summarize_file`

Extract signatures from files matching a glob pattern.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `path` | Base directory path | Server `--root` value |
| `include` | Glob pattern(s) to match files | |
| `format` | Output format (`xml`, `md`) | `xml` |

### Claude Desktop Configuration

Add to your Claude Desktop `claude_desktop_config.json`:

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
