---
layout: default
title: TOML
parent: Language Guides
nav_order: 19
---

# TOML Support

[English](toml.md) | [한국어](../ko/languages/toml.md) | [日本語](../ja/languages/toml.md) | [हिन्दी](../hi/languages/toml.md) | [Deutsch](../de/languages/toml.md)

## Supported Extensions

- `.toml`

## Grammar

- [tree-sitter-toml](https://github.com/tree-sitter-grammars/tree-sitter-toml) v0.7.0 by tree-sitter-grammars

## Extraction Targets

| Element | Kind | XML Tag | Example |
|---------|------|---------|---------|
| Table | `namespace` | `<type>` | `[package]` |
| Array of Tables | `namespace` | `<type>` | `[[bin]]` |
| Key-Value Pair | `variable` | `<variable>` | `name = "myapp"` |

## Example

### Input

```toml
# Project configuration
name = "myapp"
version = "1.0.0"

[package]
authors = ["Alice"]
edition = "2024"

[[bin]]
name = "cli"
path = "src/main.rs"
```

### Output (XML)

```xml
<file path="config.toml" language="toml">
  <variable>name = "myapp"</variable>
  <variable>version = "1.0.0"</variable>
  <type>[package]</type>
  <type>[[bin]]</type>
</file>
```

## Notes

### Body Removal

When `--include-body` flag is not used:

- Table sections (`[table]`) are stripped of their body, showing only the header
- Array of tables (`[[table]]`) are stripped of their body, showing only the header
- Top-level key-value pairs preserve their values

Use `--include-private` to include non-exported/private symbols.

### Comments

- Single-line comments (`# comment`) are extracted as documentation

### Imports

- TOML has no import system; `--include-imports` has no effect

### Limitations

- Inline tables and inline arrays within values are not separately extracted
- Dotted keys (e.g., `physical.color = "orange"`) are captured as a single pair
