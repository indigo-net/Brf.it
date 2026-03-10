# Examples

Sample source files and their brfit output for each supported language.

## Structure

Each language directory contains:

- **Source file** — a small, self-contained code sample
- **output.xml** — brfit XML output (`brfit <dir> -f xml --no-tree --no-tokens`)
- **output.md** — brfit Markdown output (`brfit <dir> -f md --no-tree --no-tokens`)

```
examples/
├── go/          # Go — structs, interfaces, methods
├── typescript/  # TypeScript — interfaces, generics, classes
├── python/      # Python — dataclasses, enums, type hints
├── java/        # Java — interfaces, classes, generics
├── rust/        # Rust — structs, enums, traits, impl blocks
├── sql/         # SQL — tables, functions, views, indexes
└── README.md
```

## Quick Start

Generate output for a single language:

```bash
brfit examples/go -f xml --no-tree --no-tokens
```

Generate output for all examples:

```bash
for lang in go typescript python java rust sql; do
  brfit "examples/$lang" -f xml --no-tree --no-tokens -o "examples/$lang/output.xml"
  brfit "examples/$lang" -f md  --no-tree --no-tokens -o "examples/$lang/output.md"
done
```

## What brfit Extracts

| Element | XML Tag | Examples |
|---------|---------|----------|
| Functions, methods | `<function>` | `func Distance() float64`, `def create(self, title: str) -> Task` |
| Types, classes, interfaces | `<type>` | `type Shape interface`, `class TaskRepository`, `CREATE TABLE` |
| Variables, constants | `<variable>` | `export const formatUser`, `CREATE INDEX` |

Body content (function implementations, view queries) is stripped — only signatures and type definitions are preserved.

## Languages

| Language | Source | Signatures | Highlights |
|----------|--------|------------|------------|
| Go | `main.go` | 9 | struct fields, interface methods, receiver methods |
| TypeScript | `app.ts` | 10 | generic class, async methods, type guards |
| Python | `api.py` | 9 | dataclass, enum, type hints |
| Java | `ShapeService.java` | 18 | interface hierarchy, generics, streams |
| Rust | `lib.rs` | 13 | generic struct, enum, trait, impl block |
| SQL | `schema.sql` | 6 | DDL: tables, functions, views, indexes, types |
