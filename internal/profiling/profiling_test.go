package profiling

import (
	"os"
	"testing"
)

func TestGetMemoryStats(t *testing.T) {
	stats := GetMemoryStats()

	if stats.Alloc == 0 {
		t.Error("expected non-zero Alloc")
	}
	if stats.Sys == 0 {
		t.Error("expected non-zero Sys")
	}
	if stats.GoroutineCount < 1 {
		t.Error("expected at least 1 goroutine")
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    uint64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		result := FormatBytes(tt.bytes)
		if result != tt.expected {
			t.Errorf("FormatBytes(%d) = %s, want %s", tt.bytes, result, tt.expected)
		}
	}
}

func TestWriteHeapProfile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "heap_profile_*.pprof")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	if err := WriteHeapProfile(tmpfile.Name()); err != nil {
		t.Errorf("WriteHeapProfile failed: %v", err)
	}

	// Verify file was written
	info, err := os.Stat(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("heap profile file is empty")
	}
}

func TestStartCPUProfile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "cpu_profile_*.pprof")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	stop, err := StartCPUProfile(tmpfile.Name())
	if err != nil {
		t.Fatalf("StartCPUProfile failed: %v", err)
	}

	// Do some work
	for i := 0; i < 1000; i++ {
		_ = GetMemoryStats()
	}

	stop()

	// Verify file was written
	info, err := os.Stat(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() == 0 {
		t.Error("CPU profile file is empty")
	}
}
