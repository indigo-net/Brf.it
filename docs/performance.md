---
layout: default
title: Performance
parent: Documentation
nav_order: 5
---

# Performance

Benchmark results for Brf.it v0.20.0 on Intel Core i7-9750H (6C/12T, 2.60 GHz), macOS, Go 1.24.

## Real-World Performance

| Project | Files | Time | Format |
|---------|-------|------|--------|
| Brf.it itself (~95 Go files, 150+ signatures) | 95 | ~0.7s | XML |
| Brf.it itself (~95 Go files, 150+ signatures) | 95 | ~0.7s | Markdown |

## Scanner Benchmarks

File system scanning with language detection, gitignore filtering, and glob matching.

| Scenario | Time/op | Allocs/op | Memory/op |
|----------|---------|-----------|-----------|
| 10 files | 137 us | 82 | 10 KB |
| 50 files | 460 us | 326 | 40 KB |
| 100 files | 870 us | 628 | 79 KB |
| 500 files | 1.1 ms | 808 | 95 KB |
| With .gitignore patterns | 227 us | 76 | 9.8 KB |

Scanner performance scales linearly with file count. Gitignore filtering adds negligible overhead.

## Parser Benchmarks (Tree-sitter)

Parsing and signature extraction via Tree-sitter CGO bindings.

### Go

| Size | Time/op | Allocs/op | Memory/op |
|------|---------|-----------|-----------|
| 10 functions | 367 us | 93 | 7.1 KB |
| 50 functions + docs | 2.4 ms | 1,418 | 74 KB |
| 100 functions + docs | 6.3 ms | 4,821 | 246 KB |
| 500 functions + docs | 46 ms | 44,035 | 2.2 MB |

### TypeScript

| Size | Time/op | Allocs/op | Memory/op |
|------|---------|-----------|-----------|
| 10 functions | 618 us | 95 | 7.6 KB |
| 50 functions | 3.4 ms | 420 | 30 KB |
| 100 functions | 5.8 ms | 823 | 63 KB |

### Python

| Size | Time/op | Allocs/op | Memory/op |
|------|---------|-----------|-----------|
| 10 functions | 442 us | 105 | 7.9 KB |
| 50 functions | 2.2 ms | 470 | 31 KB |
| 100 functions | 4.4 ms | 923 | 66 KB |

Parse time is dominated by Tree-sitter's CGO boundary crossing and query execution. Go files with doc comments use more memory due to comment node extraction.

## Formatter Benchmarks

Output generation from extracted signatures.

### XML

| Size | Time/op | Memory/op |
|------|---------|-----------|
| 5 files, 10 signatures | 19 us | 33 KB |
| 20 files, 50 signatures | 337 us | 524 KB |
| 50 files, 100 signatures | 1.5 ms | 2.1 MB |
| 100 files, 200 signatures | 5.9 ms | 8.4 MB |

### Markdown

| Size | Time/op | Memory/op |
|------|---------|-----------|
| 5 files, 10 signatures | 7.7 us | 16 KB |
| 20 files, 50 signatures | 111 us | 262 KB |
| 50 files, 100 signatures | 971 us | 2.1 MB |

### JSON

| Size | Time/op | Memory/op |
|------|---------|-----------|
| 5 files, 10 signatures | 31 us | 14 KB |
| 20 files, 50 signatures | 596 us | 270 KB |
| 50 files, 100 signatures | 3.1 ms | 1.6 MB |

Markdown is the fastest output format. XML is ~3x slower due to XML escaping. JSON is slowest due to marshaling overhead.

## Running Benchmarks

```bash
# All benchmarks
go test ./... -bench=. -benchmem

# Scanner only
go test ./pkg/scanner/ -bench=. -benchmem

# Parser only (specific language)
go test ./pkg/parser/treesitter/ -bench=BenchmarkParseGo -benchmem

# Formatter comparison
go test ./pkg/formatter/ -bench=BenchmarkFormatterComparison -benchmem

# With CPU profiling
go test ./pkg/parser/treesitter/ -bench=BenchmarkParseGo -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

## Bottleneck Analysis

| Stage | % of Total Time | Bottleneck |
|-------|----------------|------------|
| Scan | ~5% | File system I/O |
| Parse | ~85% | Tree-sitter CGO calls |
| Format | ~10% | String building, XML escaping |

Parsing dominates total runtime. The CGO boundary between Go and Tree-sitter's C library is the primary overhead source. For large codebases, enabling `PreloadContent` (default since v0.20.0) eliminates redundant file reads between scan and parse stages.
