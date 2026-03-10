# Memory Profiling Guide

This guide explains how to profile and optimize memory usage in Brf.it.

## Quick Start

### Using the profiling package

```go
import "github.com/indigo-net/Brf.it/internal/profiling"

// Get current memory stats
stats := profiling.GetMemoryStats()
fmt.Printf("Memory: %s\n", profiling.FormatBytes(stats.Alloc))

// Capture heap profile
if err := profiling.WriteHeapProfile("heap.prof"); err != nil {
    log.Fatal(err)
}

// Capture CPU profile
stop, err := profiling.StartCPUProfile("cpu.prof")
if err != nil {
    log.Fatal(err)
}
defer stop()
```

## Memory Profiling with pprof

### 1. Capture a heap profile

```bash
# While running brfit with large project
go tool pprof http://localhost:6060/debug/pprof/heap

# Or capture to file
curl -o heap.prof http://localhost:6060/debug/pprof/heap
```

### 2. Analyze heap profile

```bash
# Interactive mode
go tool pprof heap.prof

# Top memory consumers
go tool pprof -top heap.prof

# Web interface
go tool pprof -http=:8080 heap.prof
```

### 3. Common pprof commands

```
(pprof) top10        # Top 10 memory allocators
(pprof) list Scan    # Show memory allocations in Scan functions
(pprof) web          # Open in browser
(pprof) png > out.png # Generate PNG graph
```

## CPU Profiling

### Capture CPU profile

```bash
# Profile for 30 seconds
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

## Memory Optimization Tips

### 1. Reduce allocations

```go
// Bad: Creates new slice each time
func processFiles(files []string) []string {
    var results []string
    for _, f := range files {
        results = append(results, process(f))
    }
    return results
}

// Good: Pre-allocate slice
func processFiles(files []string) []string {
    results := make([]string, 0, len(files))
    for _, f := range files {
        results = append(results, process(f))
    }
    return results
}
```

### 2. Use sync.Pool for reusable buffers

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func process() {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buf)
    buf.Reset()
    // Use buffer...
}
```

### 3. Avoid string concatenation in loops

```go
// Bad
var s string
for _, item := range items {
    s += item  // Creates new string each iteration
}

// Good
var sb strings.Builder
for _, item := range items {
    sb.WriteString(item)
}
s := sb.String()
```

## Benchmarks

Run benchmarks with memory profiling:

```bash
# Basic benchmark
go test -bench=. ./pkg/scanner/

# With memory allocation stats
go test -bench=. -benchmem ./pkg/scanner/

# With allocation profile
go test -bench=. -memprofile=mem.prof ./pkg/scanner/
go tool pprof mem.prof
```

## Testing with Large Projects

For testing memory usage with large projects:

```bash
# Profile brfit running on itself
brfit . -f xml --no-tokens > /dev/null &

# Capture memory
curl -o heap.prof http://localhost:6060/debug/pprof/heap
go tool pprof -top heap.prof
```

## Common Memory Issues

### Issue: High memory with many files

**Symptom**: OOM with 10,000+ files

**Solution**:
- Use streaming instead of loading all files into memory
- Process files in batches
- Use `sync.Pool` for reusable buffers

### Issue: Memory not released after processing

**Symptom**: Memory stays high after processing completes

**Solution**:
- Check for goroutine leaks
- Ensure channels are closed
- Use `runtime.GC()` to force garbage collection (debug only)

### Issue: Large allocations during parsing

**Symptom**: Memory spikes during Tree-sitter parsing

**Solution**:
- Process files sequentially instead of in parallel
- Limit concurrent file processing
- Consider chunking large files
