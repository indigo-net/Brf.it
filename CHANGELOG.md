# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- `--changed` flag to scan only git-modified files (#179)
- `--since <ref>` flag to scan files changed since a specific commit/tag (#179)
- YAML 언어 지원 (`.yaml`, `.yml`)
- TOML 언어 지원 (`.toml`)
- `--include-private` CLI flag to include non-exported/private symbols (#177)
- `--include` / `--exclude` glob pattern filtering with doublestar support (#176)

## [0.19.0] - 2026-03-09

### Added
- Scala language support (#25)
  - Classes, traits, objects, methods, val/var, type aliases
  - Enum, given definitions (Scala 3)
  - Extension methods extracted individually
  - Import statement extraction
- Ruby language support (#21)
  - Classes, modules, methods, singleton methods, constants
  - attr_accessor/reader/writer
  - require/require_relative import extraction
- PHP language support (#22)
  - Classes, interfaces, traits, enums, functions, methods
  - Properties, constants, namespaces
  - use/require/include import extraction
- Shell/Bash language support (#118)
  - Functions, variables, aliases
- Lua language support (#26)
  - Functions, local functions, module functions
  - require import extraction
- JSON output format (`--format json`) (#73)
- Docusaurus documentation website (#112)
- Automated release GitHub Workflow (#114)
- verbatim import/export output and `--no-std-imports` flag removal (#109)

### Changed
- `normalizeKind()`: added `trait` and `impl` Kind mappings
- XML/JSON formatter Kind mapping logic unified (#138)
- Lua import filtering variable name improved (#120)
- Unused TypeScript regex variable removed (#123)

### Fixed
- Binary file detection and parse skip logic (#134)
- Tree-sitter node byte range explicit bound check (#137)
- Negative depth guard unified across all body-start functions (#133)
- PHP parser test type mismatch and config validation test (#129)
- Windows UNC path `filepath.Base` handling (#116)
- refineKotlinClassKind annotation parsing defense (#113)
- isPythonMethod nested parenthesis handling (#106)
- scanner.go `log.Printf` unnecessary `\n` removal (#101)

### Performance
- sync.Pool for Tree-sitter Parser pooling (#131)
- sync.Pool for QueryCursor pooling (#135)
- Tree-sitter Query compile caching (#61)
- `removeBlankLines` single-pass strings.Builder (#136)
- `Parser.Parse` interface changed to accept `[]byte` (#126)
- `escapeXML()` single-pass optimization (#125)
- Empty metadata field output optimization (#124)
- `stripTSFunctionBody` duplicate string scan removal (#122)
- Empty string `cleanComment` call prevention (#121)
- Capture loop unnecessary string conversion removal (#132)

## [0.18.0] - 2026-03-06

### Added
- C# language support (#81)
  - Classes, structs, interfaces, records, enums
  - Methods, constructors, properties, fields, events
  - Delegates, operators, indexers, finalizers
  - Namespaces, using directives, generic type parameters
- Parallel file extraction with worker pool (#60)
  - Configurable concurrency via `DefaultExtractOptions()`
  - scanResult nil guard and negative concurrency validation
- `Warnings` field in scan results for non-fatal errors (#77)

### Fixed
- `findFunctionBodyStart` parenDepth negative depth guard (#97)
- Markdown output: removed unnecessary kind comments (#84)
- TOCTOU race condition defense in file scanner (#79)
- Tree-sitter query error propagation (#78)
- WalkDir error classification and warning separation (#77)
- CGO panic recovery (#76)
- `.gitignore` load failure warning output (#75)
- `MaxFileSize` upper limit exceeded warning (#74)

## [0.17.0] - 2026-03-05

### Added
- Kotlin language support (#23)
  - Functions, classes, interfaces, objects, properties, typealias, enum entries
  - suspend, inline, data, sealed, annotation, value class keywords
  - Import statement extraction, companion object, inner class
- GitHub Pages documentation site (#14)
- GitHub Issue-based development workflow (#12)
- PR review automation skill and agent
- CONTRIBUTING.md contribution guide (#37)
- SECURITY.md security policy (#39)

### Changed
- Swift tree-sitter switched to vendor approach (#58)
- Test code sitter.Query resource cleanup across all languages (#57)

### Fixed
- refineKotlinClassKind classification error with parenthesized annotations (#23)
- companion_object duplicate query pattern removal (#23)
- actions/labeler@v5 syntax fix (#14)

## [0.16.0] - 2026-03-03

### Added
- Swift language support (`.swift` files)
  - Functions, classes, structs, enums, protocols, extensions
  - init/deinit, subscript, operator declarations
  - Properties (let/var), generic type parameters
  - Import statement extraction
- `add-language-support` project skill (`.claude/skills/add-language-support.md`)

## [0.15.0] - 2026-03-03

### Added
- `--no-std-imports` flag to exclude standard library imports from output
- `isStdLibImport()` function supporting 6 languages (Go, Python, JS/TS, C/C++, Java, Rust)
- Python stdlib module list (~150 modules) and Node.js builtin module list
- 35+ test cases for stdlib detection across all supported languages

## [0.13.0] - 2026-03-02

### Changed
- **Breaking Change**: XML `<signature>` tag replaced with semantic tags
  - `<function>`: function, method, constructor, destructor, arrow
  - `<type>`: class, interface, type, struct, enum, record, annotation, typedef, namespace, template
  - `<variable>`: variable, field, macro, export
  - `<signature>`: fallback for unknown kinds
- Markdown output now includes Kind as comment suffix (e.g., `func Add() // function`)
- Schema section updated with new tag descriptions

### Added
- `kindToTag()` function for Kind-to-tag mapping
- Comprehensive test coverage for all Kind values (19 kinds)
- Language-specific edge case tests (TypeScript arrow/export, C++ constructor/destructor, etc.)

## [0.12.0] - 2026-02-28

### Changed
- XML metadata restructured: `version`, `path`, `schema` moved inside `<metadata>` section
- XML schema now uses `<tag name="..." description="..." />` format inside `<schema>` element
- Removed XML header comments (version/path info now in structured elements)

## [0.11.0] - 2026-02-28

### Added
- Output header with path-based title: `# Code Summary: {path}` (Markdown), `<!-- ... | Code Summary: {path} -->` (XML)
- Version display: `*brf.it {version}*` (Markdown), included in header comment (XML)
- XML Schema description comment with tag explanations (file, signature, imports, import, export, doc, error)
- `RootPath` and `Version` fields in `PackageData` struct

## [0.10.0] - 2026-02-27

### Added
- C++ language support with full extraction capabilities
  - Classes, structs, methods, constructors, destructors
  - Namespaces (including nested)
  - Template classes and functions
  - Enums, typedefs, macros
  - `#include` directives (with `--include-imports`)
  - Extensions: `.cpp`, `.hpp`, `.h`
- Documentation: `docs/languages/cpp.md`

### Changed
- `.h` files now parsed as C++ (previously C)
- Java: removed private method filtering - all access levels now extracted
- Removed unused `isJavaPrivate()` function

## [0.9.3] - 2026-02-27

### Fixed
- Go module-level symbol extraction bug: `isExported()` was filtering by uppercase rule
  - Before: only exported symbols (uppercase first letter) were included
  - After: all module-level symbols included (same as other languages)
  - Fixes `main.go` appearing empty in CODE_PACKAGE output

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

[0.19.0]: https://github.com/indigo-net/Brf.it/compare/v0.18.0...v0.19.0
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
