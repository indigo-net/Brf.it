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
| Shell | `.sh`, `.bash`, `.zsh` | [Shell Guide](shell) |
| PHP | `.php` | [PHP Guide](php) |

## Extraction Capabilities

Each language extracts:

- **Functions/Methods** — Signatures with parameters and return types
- **Types** — Classes, structs, interfaces, enums
- **Variables** — Constants and variable declarations
- **Imports** — Module dependencies (with `--include-imports`)
- **Documentation** — Doc comments and annotations

## Language-Specific Features

Each language may have unique features:

- **Go** — Structs, interfaces, type aliases
- **TypeScript/JavaScript** — Classes, interfaces, type aliases, arrow functions
- **Python** — Classes, decorators, type hints
- **C/C++** — Structs, unions, typedefs, preprocessor directives
- **Java** — Classes, interfaces, annotations, enums
- **Rust** — Structs, enums, traits, impls, macros
- **Swift** — Structs, classes, protocols, extensions
- **Kotlin** — Classes, interfaces, objects, companion objects
- **C#** — Classes, interfaces, structs, delegates, events
- **Lua** — Functions, local functions, module functions, methods
- **Shell** — Functions, variables
- **PHP** — Classes, interfaces, traits, enums, constants
