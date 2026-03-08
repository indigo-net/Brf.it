---
layout: default
title: Language Guides
nav_order: 4
has_children: true
---

# Language Guides

Brf.it supports multiple programming languages with Tree-sitter based parsing. Each language has specific extraction targets and features.

## Supported Languages

| Language | Extensions | Guide |
|----------|------------|-------|
| Go | `.go` | [Go Guide](go) |
| TypeScript | `.ts`, `.tsx` | [TypeScript Guide](typescript) |
| JavaScript | `.js`, `.jsx` | [TypeScript Guide](typescript) |
| Python | `.py` | [Python Guide](python) |
| Java | `.java` | [Java Guide](java) |
| C | `.c` | [C Guide](c) |
| C++ | `.cpp`, `.hpp`, `.h` | [C++ Guide](cpp) |
| Rust | `.rs` | [Rust Guide](rust) |
| Swift | `.swift` | [Swift Guide](swift) |
| Kotlin | `.kt`, `.kts` | [Kotlin Guide](kotlin) |
| C# | `.cs` | [C# Guide](csharp) |
| Lua | `.lua` | [Lua Guide](lua) |
| PHP | `.php` | [PHP Guide](php) |

## Extraction Capabilities

Each language extracts:

- **Functions/Methods** — Signatures with parameters and return types
- **Types** — Classes, structs, interfaces, enums
- **Variables** — Constants and variable declarations
- **Documentation** — Comments (JSDoc, GoDoc, docstrings, etc.)

## Output Examples

### XML Format

```xml
<file path="example.go" language="go">
  <signature kind="function" line="5">
    <name>DoSomething</name>
    <text>func DoSomething(input string) error</text>
    <doc>DoSomething performs an action.</doc>
  </signature>
</file>
```

### Markdown Format

````markdown
### example.go

```go
func DoSomething(input string) error
```
DoSomething performs an action.
````
