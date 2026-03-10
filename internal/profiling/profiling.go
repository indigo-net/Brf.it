// Package profiling provides memory profiling utilities for brfit.
package profiling

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

// MemoryStats holds memory usage information.
type MemoryStats struct {
	// Alloc is bytes of allocated heap objects.
	Alloc uint64

	// TotalAlloc is cumulative bytes allocated for heap objects.
	TotalAlloc uint64

	// Sys is total bytes of memory obtained from the OS.
	Sys uint64

	// NumGC is the number of completed GC cycles.
	NumGC uint32

	// GoroutineCount is the number of goroutines.
	GoroutineCount int

	// HeapObjects is the number of allocated heap objects.
	HeapObjects uint64
}

// GetMemoryStats returns current memory statistics.
func GetMemoryStats() MemoryStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemoryStats{
		Alloc:          m.Alloc,
		TotalAlloc:     m.TotalAlloc,
		Sys:            m.Sys,
		NumGC:          m.NumGC,
		GoroutineCount: runtime.NumGoroutine(),
		HeapObjects:    m.HeapObjects,
	}
}

// FormatBytes formats bytes as a human-readable string.
func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// WriteHeapProfile writes a heap profile to the specified file.
func WriteHeapProfile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	runtime.GC() // Run GC before capturing heap profile
	return pprof.WriteHeapProfile(f)
}

// StartCPUProfile starts CPU profiling to the specified file.
// Returns a stop function that should be called to stop profiling.
func StartCPUProfile(filename string) (func(), error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		f.Close()
		return nil, err
	}

	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}, nil
}
