# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.9.2] - 2026-02-27

### Added
- Empty file indicator: files with no signatures show `// (empty)` comment
  - Markdown: language-specific comment (Python: `# (empty)`, HTML: `<!-- (empty) -->`)
  - XML: `<!-- empty -->`
- Helper function `getEmptyComment()` for language-aware empty comments

## [0.9.1] - 2026-02-26

### Changed
- Imports now rendered inside each file block (not separate section)
- Full import statement output instead of path-only
  - TypeScript: `import { useState } from 'react';`
  - Go: `import "fmt"` / `import alias "path"`
  - Python: `import os` / `from pathlib import Path`
  - Java: `import java.util.List;`
  - C: `#include <stdio.h>`

## [0.9.0] - 2026-02-26

### Added
- `--include-imports` flag to extract import/export statements
  - Go: `import` declarations
  - TypeScript/JavaScript: `import` statements and `export` declarations
  - Python: `import` and `from ... import` statements
  - Java: `import` declarations (including static imports)
  - C: `#include` directives
- Imports rendered in separate `<imports>` section (XML) or "Imports & Exports" section (Markdown)
- File-grouped output for better readability

## [0.8.0] - 2026-02-26

### Added
- Module-level constant/variable extraction for all supported languages
  - Go: `const` and `var` declarations (including grouped declarations)
  - TypeScript/JavaScript: module-level `const`, `let`, `var` (local variables excluded)
  - Python: module-level assignments (class/function-level excluded)
  - Java: `static` fields only (instance fields excluded)
  - C: global variable declarations (local variables excluded)
- Variables are output with their values preserved (not stripped like function bodies)
- New "variable" kind for extracted constants/variables

### Changed
- Improved duplicate detection for TypeScript arrow functions

## [0.7.0] - 2026-02-26

### Added
- Java language support (`.java`)
  - Class, interface, method, constructor extraction
  - Enum, annotation (`@interface`), record (Java 14+) extraction
  - Private member filtering (excluded by default)
  - Generic type parameter preservation
  - Inner/nested class support
- Java language documentation in 5 languages (EN, KO, JA, HI, DE)
- Updated Supported Languages table in all README versions

## [0.6.0] - 2026-02-26

### Added
- Cross-platform install scripts (Bash and PowerShell)
- MIT License

### Changed
- Improved CLAUDE.md with new language addition checklist and debugging tips

## [0.5.0] - 2026-02-26

### Added
- C language support (`.c`, `.h`)
  - Function definitions and declarations extraction
  - Struct and enum extraction
  - Typedef extraction
  - Object-like and function-like macro extraction
  - Pointer return type support for functions
  - Multi-line macro support
- C language documentation in 5 languages (EN, KO, JA, HI, DE)
- Updated Supported Languages table in all README versions

## [0.4.0] - 2026-02-26

### Added
- Python (`.py`) language support
  - Function and async function extraction
  - Class and method extraction
  - Type hints preserved in signatures
  - Automatic method detection (self/cls first parameter)
- Language documentation: `docs/languages/go.md`, `typescript.md`, `python.md`
- Documentation links in README Supported Languages table

## [0.3.0] - 2026-02-24

### Added
- `--include-body` flag to include function/method bodies in output (default: signatures only)
- JavaScript (`.js`) and JSX (`.jsx`) language support

### Changed
- Arrow function query improved to include `const`/`let`/`var` keywords in signature

### Fixed
- Body stripping logic for Go, TypeScript, and JavaScript functions

## [0.2.7] - 2026-02-24

### Changed
- Removed duplicate `Symbols` section from XML and Markdown output
- Signatures are already listed per-file, so the summary was redundant

## [0.2.6] - 2026-02-24

### Fixed
- GoReleaser v2.6 compatibility (format singular form)

## [0.2.5] - 2026-02-24

### Fixed
- GoReleaser v2 syntax and goamd64 variant specification

## [0.2.4] - 2026-02-24

### Fixed
- Downgraded to GoReleaser v1 syntax for compatibility

[0.9.0]: https://github.com/indigo-net/Brf.it/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/indigo-net/Brf.it/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/indigo-net/Brf.it/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/indigo-net/Brf.it/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/indigo-net/Brf.it/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/indigo-net/Brf.it/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/indigo-net/Brf.it/compare/v0.2.7...v0.3.0
[0.2.7]: https://github.com/indigo-net/Brf.it/compare/v0.2.6...v0.2.7
[0.2.6]: https://github.com/indigo-net/Brf.it/compare/v0.2.5...v0.2.6
[0.2.5]: https://github.com/indigo-net/Brf.it/compare/v0.2.4...v0.2.5
[0.2.4]: https://github.com/indigo-net/Brf.it/releases/tag/v0.2.4
