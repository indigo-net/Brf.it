---
sidebar_position: 1
title: Supported Languages
---

# Supported Languages

Brf.it uses Tree-sitter for accurate syntax parsing across 14 programming languages.

## Overview

| Language | Extensions | Status |
|----------|------------|--------|
| Go | `.go` | ✅ Full support |
| TypeScript | `.ts`, `.tsx` | ✅ Full support |
| JavaScript | `.js`, `.jsx`, `.mjs`, `.cjs` | ✅ Full support |
| Python | `.py` | ✅ Full support |
| Java | `.java` | ✅ Full support |
| Kotlin | `.kt`, `.kts` | ✅ Full support |
| Rust | `.rs` | ✅ Full support |
| Ruby | `.rb`, `.rake` | ✅ Full support |
| PHP | `.php` | ✅ Full support |
| Swift | `.swift` | ✅ Full support |
| Scala | `.scala`, `.sc` | ✅ Full support |
| C/C++ | `.c`, `.cpp`, `.cc`, `.cxx`, `.h`, `.hpp` | ✅ Full support |
| YAML | `.yaml`, `.yml` | ✅ Full support |
| TOML | `.toml` | ✅ Full support |

## Extracted Elements

### All Languages

- Function/method signatures
- Struct/class/interface definitions
- Documentation comments (docstrings, JSDoc, etc.)
- Import/export statements (with `--include-imports`)

### Language-Specific

| Language | Special Handling |
|----------|------------------|
| Go | Methods on structs, interfaces |
| TypeScript | Type annotations, generics |
| Python | Decorators, type hints |
| Java | Annotations, generics |
| Kotlin | Extension functions, data classes |
| Rust | Traits, impl blocks |
| Ruby | Modules, class methods |
| PHP | Namespaces, traits |
| Swift | Protocols, extensions |
| Scala | Case classes, traits |
| C/C++ | Templates, namespaces |

## Tree-sitter Quality

All languages use official or well-maintained Tree-sitter grammars for reliable parsing:

- Handles syntax errors gracefully
- Preserves formatting in signatures
- Supports modern language features

## Adding New Languages

Want to add support for a new language? See the [contributing guide](https://github.com/indigo-net/Brf.it/blob/main/CONTRIBUTING.md).
