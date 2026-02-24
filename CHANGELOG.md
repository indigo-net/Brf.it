# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.3.0]: https://github.com/indigo-net/Brf.it/compare/v0.2.7...v0.3.0
[0.2.7]: https://github.com/indigo-net/Brf.it/compare/v0.2.6...v0.2.7
[0.2.6]: https://github.com/indigo-net/Brf.it/compare/v0.2.5...v0.2.6
[0.2.5]: https://github.com/indigo-net/Brf.it/compare/v0.2.4...v0.2.5
[0.2.4]: https://github.com/indigo-net/Brf.it/releases/tag/v0.2.4
