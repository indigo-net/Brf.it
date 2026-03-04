# Contributing to Brf.it

Thank you for your interest in contributing to **Brf.it** — a CLI tool that extracts function signatures and documentation from codebases for AI coding assistants. We welcome contributions of all kinds: bug fixes, new language support, documentation improvements, and more.

---

## Table of Contents

- [Development Setup](#development-setup)
- [Project Architecture](#project-architecture)
- [Coding Standards](#coding-standards)
- [Commit Convention](#commit-convention)
- [Branch and PR Workflow](#branch-and-pr-workflow)
- [Running Tests](#running-tests)
- [Adding a New Language](#adding-a-new-language)
- [Reporting Issues](#reporting-issues)
- [License](#license)

---

## Development Setup

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| **Go** | 1.25+ | Core language |
| **C compiler** | gcc / clang | Required for CGO (Tree-sitter bindings) |
| **Git** | any recent | Version control |
| **goimports** | latest | Code formatting |

**Platform-specific setup:**

- **macOS**: Install Xcode Command Line Tools for the C compiler:
  ```bash
  xcode-select --install
  ```
- **Linux (Debian/Ubuntu)**: Install build essentials:
  ```bash
  sudo apt-get install build-essential
  ```

### Clone and Build

```bash
git clone https://github.com/indigo-net/Brf.it.git
cd Brf.it
go build ./cmd/brfit
```

### Verify Installation

```bash
# Run all tests
go test ./...

# Check the binary works
./brfit --version
```

If `go test` fails with linker errors, ensure your C compiler is properly installed (CGO is required for Tree-sitter).

---

## Project Architecture

### Directory Layout

```
Brf.it/
├── cmd/
│   └── brfit/
│       └── main.go              # CLI entry point (Cobra)
├── pkg/
│   ├── scanner/                 # File system scanning & .gitignore filtering
│   ├── parser/                  # Parser interface & registry
│   │   └── treesitter/          # Tree-sitter implementation
│   │       ├── languages/       # Per-language query definitions
│   │       └── grammars/        # Vendored C grammar sources
│   ├── extractor/               # Signature extraction orchestration
│   ├── formatter/               # Output formatters (XML, Markdown)
│   └── tokenizer/               # Token counting
├── internal/
│   ├── config/                  # CLI configuration & defaults
│   ├── context/                 # Build context (directory tree)
│   └── logger/                  # Logging utilities
├── docs/
│   ├── languages/               # Per-language guides (English)
│   ├── ko/                      # Korean translations
│   ├── ja/                      # Japanese translations
│   ├── de/                      # German translations
│   └── hi/                      # Hindi translations
└── assets/
    └── wasm/                    # Tree-sitter WASM files (if needed)
```

### Core Interfaces

The project is built around two key interfaces:

- **`Parser`** (`pkg/parser/parser.go`) — Defines how code files are parsed to extract signatures. The Tree-sitter implementation is the primary parser.
- **`LanguageQuery`** (`pkg/parser/treesitter/query.go`) — Defines per-language Tree-sitter query patterns, capture names, and kind mappings. Each supported language implements this interface.

These interfaces allow new languages and parser backends to be added without modifying existing code.

---

## Coding Standards

### Formatting

- **Always** run `gofmt` and `goimports` before committing. Code that doesn't pass `gofmt` will not be accepted.

### Naming

| Scope | Convention | Example |
|-------|-----------|---------|
| Exported (public) | `PascalCase` | `ScanDirectory`, `FileEntry` |
| Unexported (internal) | `camelCase` | `parseNode`, `kindMapping` |
| Packages | short, lowercase, single word | `scanner`, `parser`, `formatter` |

### Error Handling

- Always handle errors explicitly. Never discard errors with `_`.
  ```go
  // Good
  result, err := doSomething()
  if err != nil {
      return fmt.Errorf("failed to do something: %w", err)
  }

  // Bad — never do this
  result, _ := doSomething()
  ```

### Documentation

- All exported functions, types, and constants **must** have GoDoc comments.
  ```go
  // ScanDirectory recursively scans the given root directory
  // and returns a list of supported source files.
  func ScanDirectory(root string, opts *ScanOptions) ([]FileEntry, error) {
  ```

### Design Principles

- **Interface-first**: Define interfaces for extensibility (`Parser`, `LanguageQuery`, `Formatter`).
- **Composition over inheritance**: Use struct embedding to compose behavior.
- **Concurrency**: Use goroutines and channels for parallel file scanning, but apply worker pool patterns to avoid excessive goroutine creation.

---

## Commit Convention

All commit messages follow this format:

```
type: summary description (#issue)
```

- The summary is written in **Korean** by convention, but **English is also accepted** for external contributors.
- Use **present tense** (describe what the commit does, not what it did).

### Types

| Type | Description | Example |
|------|-------------|---------|
| `feat` | New feature | `feat: Add Rust language support (#42)` |
| `fix` | Bug fix | `fix: Fix memory leak on large files (#38)` |
| `docs` | Documentation | `docs: Update installation guide (#37)` |
| `style` | Formatting (no logic change) | `style: Apply gofmt` |
| `refactor` | Code restructuring (no behavior change) | `refactor: Separate parser interface` |
| `test` | Add or update tests | `test: Add scanner unit tests` |
| `chore` | Build, CI, or project config | `chore: Update goreleaser config` |

---

## Branch and PR Workflow

We follow an **issue-driven workflow**. Every contribution should be linked to a GitHub issue.

### Step-by-step

1. **Find or create an issue** describing the work.
2. **Fork the repository** and create a branch from `main`:
   ```bash
   git checkout -b feat/your-feature-name
   ```
   Branch name format: `{type}/{feature-name}` (e.g., `feat/swift-support`, `fix/memory-leak`, `docs/contributing-guide`).
3. **Make your changes** following the coding standards above.
4. **Commit** with a proper message referencing the issue:
   ```bash
   git commit -m "feat: Add Swift language support (#42)"
   ```
5. **Push** and open a Pull Request:
   ```bash
   git push origin feat/your-feature-name
   ```
6. In the PR description, include `Closes #XX` to auto-close the linked issue on merge.

### PR Checklist

Your PR should address the following (from our [PR template](.github/PULL_REQUEST_TEMPLATE.md)):

- [ ] Code follows the project conventions (Go style, naming, error handling)
- [ ] Commit messages follow the `type: description` format
- [ ] No unnecessary logs or comments left behind
- [ ] Unit tests pass: `go test ./pkg/...`
- [ ] Build succeeds: `go build ./cmd/brfit`
- [ ] Manual testing performed (if applicable)

---

## Running Tests

### Full Test Suite

```bash
go test ./...
```

> **Note**: CGO must be enabled (it is by default). If you see linker errors, ensure a C compiler is installed.

### Running Specific Tests

```bash
# Test a specific package
go test ./pkg/parser/treesitter/languages/ -v

# Test a specific function
go test ./pkg/scanner/ -run TestScanDirectory -v
```

### Parser Blank Import Pattern

When writing tests that use Tree-sitter parsers, you must trigger parser auto-registration with a blank import:

```go
import (
    _ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"
)
```

Without this import, parsers won't be registered and tests will fail silently.

### Table-Driven Tests

We use Go's table-driven test pattern:

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"basic case", "input1", "expected1"},
        {"edge case", "input2", "expected2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := doSomething(tt.input)
            if result != tt.expected {
                t.Errorf("got %q, want %q", result, tt.expected)
            }
        })
    }
}
```

---

## Adding a New Language

Adding a new programming language to Brf.it involves implementing a `LanguageQuery`, registering it with the parser, updating scanner configuration, writing tests, and adding documentation.

### High-Level Checklist

1. **Grammar**: Vendor the Tree-sitter grammar C sources into `pkg/parser/treesitter/grammars/<lang>/`
2. **LanguageQuery**: Implement `LanguageQuery` interface in `pkg/parser/treesitter/languages/<lang>.go`
3. **Registration**: Register the parser in `pkg/parser/treesitter/parser.go` (`init()`, `queries` map, `isExported()`, `stripBody()`)
4. **Scanner**: Add file extension mapping in `pkg/scanner/scanner.go` and `internal/config/config.go`
5. **Unit tests**: Write ~14 unit tests in `pkg/parser/treesitter/languages/<lang>_test.go`
6. **Integration tests**: Add 3 integration test cases in `pkg/parser/treesitter/parser_test.go`
7. **Documentation**: Create language guide in `docs/languages/<lang>.md` plus 4 translations (`docs/{ko,ja,hi,de}/languages/<lang>.md`), and add a row to the Supported Languages table in all 5 README files

### Key Principle

> Abbreviation is acceptable, omission is not.

When writing query patterns, ensure **all declaration types** of the target language are captured. A function body being stripped is fine — a function being entirely missing from output is a bug.

### Reference Implementation

- **Vendored grammar**: `pkg/parser/treesitter/grammars/kotlin/binding.go`
- **LanguageQuery**: `pkg/parser/treesitter/languages/kotlin.go`
- **Detailed guide**: See the internal skill document (`.claude/skills/add-language-support.md`) for a complete step-by-step walkthrough with templates and test patterns.

### Verification

```bash
# Build
go build ./cmd/brfit

# Unit tests for your language
go test ./pkg/parser/treesitter/languages/ -run <Lang> -v

# Integration tests
go test ./pkg/parser/treesitter/ -run <Lang> -v

# Full regression
go test ./...

# Manual test
echo '<sample code>' > /tmp/test.<ext>
./brfit /tmp -f xml
./brfit /tmp -f md
rm /tmp/test.<ext>
```

---

## Reporting Issues

### Bug Reports

When filing a bug report, please include:

- **Brf.it version**: output of `brfit --version`
- **OS and Go version**: output of `go version`
- **Steps to reproduce**: minimal commands or code to trigger the bug
- **Expected vs. actual behavior**
- **Error output**: full terminal output if applicable

### Feature Requests

Describe:
- **What** you'd like to see
- **Why** it would be useful
- **How** it might work (optional, but helpful)

### Security Issues

If you discover a security vulnerability, please **do not** open a public issue. Instead, contact the maintainers directly via the repository's security advisories.

---

## License

Brf.it is licensed under the [MIT License](LICENSE). By contributing, you agree that your contributions will be licensed under the same terms.
