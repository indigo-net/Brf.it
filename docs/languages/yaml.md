---
layout: default
title: YAML
parent: Language Guides
nav_order: 20
---

# YAML Support

[English](yaml.md) | [한국어](../ko/languages/yaml.md) | [日本語](../ja/languages/yaml.md) | [हिन्दी](../hi/languages/yaml.md) | [Deutsch](../de/languages/yaml.md)

## Supported Extensions

- `.yaml`
- `.yml`

## Grammar

- [tree-sitter-yaml](https://github.com/tree-sitter-grammars/tree-sitter-yaml) v0.7.2 by tree-sitter-grammars

## Extraction Targets

| Element | Kind | XML Tag | Example |
|---------|------|---------|---------|
| Key-Value Pair | `variable` | `<variable>` | `name: value` |

## Example

### Input

```yaml
# Application configuration
name: myapp
version: 1.0.0

database:
  host: localhost
  port: 5432

features:
  - logging
  - metrics
```

### Output (XML)

```xml
<file path="config.yaml" language="yaml">
  <variable>name: myapp</variable>
  <variable>version: 1.0.0</variable>
  <variable>database:</variable>
  <variable>features:</variable>
</file>
```

## Notes

### Body Removal

When `--include-body` flag is not used:

- Container keys (mappings with nested values) are stripped of their nested content, showing only the key
- Scalar key-value pairs preserve their values

Use `--include-private` to include non-exported/private symbols.

### Comments

- Single-line comments (`# comment`) are extracted as documentation

### Imports

- YAML has no import system; `--include-imports` has no effect

### Limitations

- Only top-level keys are captured as signatures to avoid excessive noise from deeply nested values
- Anchors and aliases are not specially handled
