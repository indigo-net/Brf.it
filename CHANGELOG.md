# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.7.0]: https://github.com/indigo-net/Brf.it/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/indigo-net/Brf.it/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/indigo-net/Brf.it/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/indigo-net/Brf.it/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/indigo-net/Brf.it/compare/v0.2.7...v0.3.0
[0.2.7]: https://github.com/indigo-net/Brf.it/compare/v0.2.6...v0.2.7
[0.2.6]: https://github.com/indigo-net/Brf.it/compare/v0.2.5...v0.2.6
[0.2.5]: https://github.com/indigo-net/Brf.it/compare/v0.2.4...v0.2.5
[0.2.4]: https://github.com/indigo-net/Brf.it/releases/tag/v0.2.4
