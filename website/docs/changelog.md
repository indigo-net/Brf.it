---
sidebar_position: 4
title: Changelog
description: All notable changes to Brf.it, following Keep a Changelog format and Semantic Versioning.
---

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

For the full changelog, see [CHANGELOG.md on GitHub](https://github.com/indigo-net/Brf.it/blob/main/CHANGELOG.md).

---

## [0.20.0] - 2026-03-13

### Added
- `brfit-mcp` MCP (Model Context Protocol) server binary for AI agents
- Remote Git repository analysis (`--remote` flag)
- Function call graph output (`--call-graph`) — Tree-sitter based (Go, TypeScript, Python, Java, Rust, C)
- Security secret detection and masking (`--security-check`) — 12 pattern types
- `--token-tree` for per-file token count tree
- `--changed` / `--since <ref>` for Git change detection
- `--include` / `--exclude` glob pattern filtering
- `--include-private` for unexported symbols
- `--dedupe-imports` for import deduplication
- `--no-schema` to skip XML schema section
- `--max-doc-length` for doc comment length limit
- Multiple `--ignore` flag support
- SQL, Elixir, YAML, TOML language support
- `context.Context` support across Scanner, Packager, and Extract
- Memory profiling, benchmarks, codecov integration, fuzzing tests

### Fixed
- Hardcoded language list replaced with dynamic Registry lookup
- LanguageMapping synchronization and single source of truth
- `captureNames` index bounds checking
- `extractFile` context propagation
- `body-strip` index 0 edge case

### Changed
- `LanguageQuery.Captures()` deduplication via common embedding
- `isExported()` simplification
- `sync.Map` + `RWMutex` dual-lock simplification
- Removed duplicate TreeSitterParser instantiation in `init()`

---

## [0.9.0] - 2026-03-03

### Added
- C# language support
- Lua language support
- PHP language support
- Ruby language support
- Scala language support

---

## [0.8.0] - 2026-03-03

### Added
- Swift language support
- Kotlin language support

---

## [0.7.0] - 2026-03-03

### Added
- C++ language support

---

## [0.6.0] - 2026-03-03

### Added
- C language support

---

## [0.5.0] - 2026-03-02

### Added
- Rust language support

---

## [0.4.0] - 2026-03-02

### Added
- Java language support

---

## [0.3.0] - 2026-03-02

### Added
- Python language support

---

## [0.2.0] - 2026-03-02

### Added
- TypeScript/JavaScript language support

---

## [0.1.0] - 2026-03-01

### Added
- Initial release
- Go language support
- XML and Markdown output formats
- Tree-sitter based parsing
- Token counting
- `.gitignore` integration
