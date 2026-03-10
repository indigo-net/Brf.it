package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

// BenchmarkScanDirectory benchmarks scanning directories with varying file counts.
func BenchmarkScanDirectory(b *testing.B) {
	tests := []struct {
		name      string
		fileCount int
	}{
		{"10_files", 10},
		{"50_files", 50},
		{"100_files", 100},
		{"500_files", 500},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			// Create temp directory with test files
			tmpDir, err := os.MkdirTemp("", "benchmark-*")
			if err != nil {
				b.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			// Create test files
			for i := 0; i < tt.fileCount; i++ {
				filename := filepath.Join(tmpDir, "file_"+string(rune('a'+i%26))+string(rune('0'+i%10))+".go")
				content := []byte(`package test

func Func` + string(rune('A'+i%26)) + `() int {
	return ` + string(rune('0'+i%10)) + `
}
`)
				if err := os.WriteFile(filename, content, 0644); err != nil {
					b.Fatal(err)
				}
			}

			opts := DefaultScanOptions()
			opts.RootPath = tmpDir
			s, err := NewFileScanner(opts)
			if err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := s.Scan()
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkScanLargeFile benchmarks scanning a single large file.
func BenchmarkScanLargeFile(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark-large-*")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a large Go file
	largeContent := []byte(`package test

// LargeFile is a test file with many functions.
`)

	// Add many functions
	for i := 0; i < 1000; i++ {
		largeContent = append(largeContent, []byte(`
func Func`+string(rune('A'+i%26))+string(rune('a'+(i/26)%26))+`() int {
	// This is function number `+string(rune('0'+i%10))+`
	return `+string(rune('0'+i%10))+`
}
`)...)
	}

	filename := filepath.Join(tmpDir, "large.go")
	if err := os.WriteFile(filename, largeContent, 0644); err != nil {
		b.Fatal(err)
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir
	s, err := NewFileScanner(opts)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := s.Scan()
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkScanWithIgnore benchmarks scanning with gitignore patterns.
func BenchmarkScanWithIgnore(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "benchmark-ignore-*")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directory structure
	dirs := []string{"src", "vendor", "node_modules", ".git"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			b.Fatal(err)
		}
	}

	// Create files
	for i := 0; i < 50; i++ {
		for _, dir := range dirs {
			filename := filepath.Join(tmpDir, dir, "file.go")
			content := []byte(`package test; func F() {}`)
			if err := os.WriteFile(filename, content, 0644); err != nil {
				// Ignore if file exists
			}
		}
	}

	// Create .gitignore
	gitignore := []byte(`vendor/
node_modules/
.git/
`)
	if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), gitignore, 0644); err != nil {
		b.Fatal(err)
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir
	s, err := NewFileScanner(opts)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := s.Scan()
		if err != nil {
			b.Fatal(err)
		}
	}
}
