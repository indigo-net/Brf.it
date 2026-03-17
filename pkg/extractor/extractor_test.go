package extractor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/indigo-net/Brf.it/pkg/parser"
	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter" // Register Tree-sitter parsers
	"github.com/indigo-net/Brf.it/pkg/scanner"
)

func TestFileExtractorImplementsExtractor(t *testing.T) {
	// Verify FileExtractor implements Extractor interface
	var _ Extractor = (*FileExtractor)(nil)
}

func TestFileExtractorExtract(t *testing.T) {
	// Create temp directory with test files
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test Go file
	testFile := filepath.Join(tmpDir, "test.go")
	testCode := `package test

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

type Point struct {
	X, Y int
}
`
	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	// Scan files
	defaultOpts := scanner.DefaultScanOptions()
	scanOpts := &scanner.ScanOptions{
		RootPath:            tmpDir,
		SupportedExtensions: defaultOpts.SupportedExtensions,
		MaxFileSize:         defaultOpts.MaxFileSize,
	}
	sc, err := scanner.NewFileScanner(scanOpts)
	if err != nil {
		t.Fatal(err)
	}
	scanResult, err := sc.Scan(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Extract signatures
	extractor := NewDefaultFileExtractor()
	result, err := extractor.Extract(context.Background(), scanResult, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Verify results
	if len(result.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(result.Files))
	}

	if result.TotalSignatures < 1 {
		t.Errorf("expected at least 1 signature, got %d", result.TotalSignatures)
	}

	// Find Add function
	var foundAdd bool
	for _, sig := range result.Files[0].Signatures {
		if sig.Name == "Add" {
			foundAdd = true
			if sig.Kind != "function" {
				t.Errorf("expected kind 'function', got '%s'", sig.Kind)
			}
		}
	}

	if !foundAdd {
		t.Error("expected to find 'Add' function signature")
	}
}

func TestFileExtractorTOCTOUGuard(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "big.go")
	bigContent := "package big\n" + strings.Repeat("// padding\n", 100)
	if err := os.WriteFile(testFile, []byte(bigContent), 0644); err != nil {
		t.Fatal(err)
	}

	scanResult := &scanner.ScanResult{
		Files: []scanner.FileEntry{
			{Path: testFile, Language: "go", Size: 10},
		},
	}

	extractor := NewDefaultFileExtractor()
	opts := &ExtractOptions{
		MaxFileSize: 50,
	}

	result, err := extractor.Extract(context.Background(), scanResult, opts)
	if err != nil {
		t.Fatal(err)
	}

	if result.ErrorCount != 1 {
		t.Errorf("expected 1 error (TOCTOU size mismatch), got %d", result.ErrorCount)
	}

	if result.Files[0].Error == nil {
		t.Error("expected error for file size mismatch")
	} else if !strings.Contains(result.Files[0].Error.Error(), "file size changed") {
		t.Errorf("expected 'file size changed' error, got: %v", result.Files[0].Error)
	}
}

func TestFileExtractorTOCTOUGuardDisabled(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.go")
	testCode := "package test\n\nfunc Foo() {}\n"
	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	scanResult := &scanner.ScanResult{
		Files: []scanner.FileEntry{
			{Path: testFile, Language: "go", Size: int64(len(testCode))},
		},
	}

	extractor := NewDefaultFileExtractor()
	opts := &ExtractOptions{MaxFileSize: 0}

	result, err := extractor.Extract(context.Background(), scanResult, opts)
	if err != nil {
		t.Fatal(err)
	}

	if result.ErrorCount != 0 {
		t.Errorf("expected 0 errors when TOCTOU disabled, got %d", result.ErrorCount)
	}
}

func TestExtractConcurrencySequential(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.go")
	testCode := `package test

func Hello() string { return "hello" }
`
	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	scanResult := &scanner.ScanResult{
		Files: []scanner.FileEntry{
			{Path: testFile, Language: "go", Size: int64(len(testCode))},
		},
	}

	extractor := NewDefaultFileExtractor()
	result, err := extractor.Extract(context.Background(), scanResult, &ExtractOptions{Concurrency: 1})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(result.Files))
	}
	if result.TotalSignatures < 1 {
		t.Errorf("expected at least 1 signature, got %d", result.TotalSignatures)
	}
}

func TestExtractConcurrencyDeterministicOrder(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create multiple Go files
	fileNames := []string{"a.go", "b.go", "c.go", "d.go", "e.go"}
	var entries []scanner.FileEntry
	for i, name := range fileNames {
		path := filepath.Join(tmpDir, name)
		code := fmt.Sprintf("package test\n\nfunc Func%d() {}\n", i)
		if err := os.WriteFile(path, []byte(code), 0644); err != nil {
			t.Fatal(err)
		}
		entries = append(entries, scanner.FileEntry{
			Path:     path,
			Language: "go",
			Size:     int64(len(code)),
		})
	}

	scanResult := &scanner.ScanResult{Files: entries}
	extractor := NewDefaultFileExtractor()
	result, err := extractor.Extract(context.Background(), scanResult, &ExtractOptions{Concurrency: 2})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Files) != len(fileNames) {
		t.Fatalf("expected %d files, got %d", len(fileNames), len(result.Files))
	}

	// Verify order matches input
	for i, name := range fileNames {
		expectedPath := filepath.Join(tmpDir, name)
		if result.Files[i].Path != expectedPath {
			t.Errorf("file[%d]: expected path %s, got %s", i, expectedPath, result.Files[i].Path)
		}
	}
}

func TestExtractConcurrencyEmptyFiles(t *testing.T) {
	scanResult := &scanner.ScanResult{Files: []scanner.FileEntry{}}
	extractor := NewDefaultFileExtractor()

	result, err := extractor.Extract(context.Background(), scanResult, &ExtractOptions{Concurrency: 4})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(result.Files))
	}
	if result.TotalSignatures != 0 {
		t.Errorf("expected 0 signatures, got %d", result.TotalSignatures)
	}
}

func TestExtractNilScanResult(t *testing.T) {
	extractor := NewDefaultFileExtractor()
	result, err := extractor.Extract(context.Background(), nil, nil)
	if err != nil {
		t.Fatalf("expected no error for nil scanResult, got: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(result.Files))
	}
}

func TestExtractNegativeConcurrency(t *testing.T) {
	scanResult := &scanner.ScanResult{Files: []scanner.FileEntry{}}
	extractor := NewDefaultFileExtractor()
	_, err := extractor.Extract(context.Background(), scanResult, &ExtractOptions{Concurrency: -1})
	if err == nil {
		t.Fatal("expected error for negative concurrency, got nil")
	}
}

func TestDefaultExtractOptions(t *testing.T) {
	opts := DefaultExtractOptions()
	if opts == nil {
		t.Fatal("expected non-nil options")
	}
	if opts.Concurrency != 0 {
		t.Errorf("expected Concurrency=0 (auto), got %d", opts.Concurrency)
	}
}

func TestExtractConcurrencyWithErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create one valid file
	validFile := filepath.Join(tmpDir, "valid.go")
	validCode := "package test\n\nfunc Valid() {}\n"
	if err := os.WriteFile(validFile, []byte(validCode), 0644); err != nil {
		t.Fatal(err)
	}

	// One entry points to a non-existent file
	scanResult := &scanner.ScanResult{
		Files: []scanner.FileEntry{
			{Path: validFile, Language: "go", Size: int64(len(validCode))},
			{Path: filepath.Join(tmpDir, "missing.go"), Language: "go", Size: 100},
			{Path: validFile, Language: "go", Size: int64(len(validCode))},
		},
	}

	extractor := NewDefaultFileExtractor()
	result, err := extractor.Extract(context.Background(), scanResult, &ExtractOptions{Concurrency: 2})
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Files) != 3 {
		t.Fatalf("expected 3 files, got %d", len(result.Files))
	}
	if result.ErrorCount != 1 {
		t.Errorf("expected 1 error, got %d", result.ErrorCount)
	}

	// Error should be at index 1 (the missing file)
	if result.Files[0].Error != nil {
		t.Errorf("file[0] should succeed, got error: %v", result.Files[0].Error)
	}
	if result.Files[1].Error == nil {
		t.Error("file[1] should have error for missing file")
	}
	if result.Files[2].Error != nil {
		t.Errorf("file[2] should succeed, got error: %v", result.Files[2].Error)
	}
}

func TestFileExtractorUnsupportedLanguage(t *testing.T) {
	// Create registry without any parsers
	registry := parser.NewRegistry()
	extractor := NewFileExtractor(registry)

	scanResult := &scanner.ScanResult{
		Files: []scanner.FileEntry{
			{Path: "test.py", Language: "python", Size: 100},
		},
	}

	result, err := extractor.Extract(context.Background(), scanResult, nil)
	if err != nil {
		t.Fatal(err)
	}

	if result.ErrorCount != 1 {
		t.Errorf("expected 1 error, got %d", result.ErrorCount)
	}

	if result.Files[0].Error == nil {
		t.Error("expected error for unsupported language")
	}
}

func TestExtractCanceledContext(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.go")
	testCode := "package test\n\nfunc Foo() {}\n"
	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	scanResult := &scanner.ScanResult{
		Files: []scanner.FileEntry{
			{Path: testFile, Language: "go", Size: int64(len(testCode))},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	ext := NewDefaultFileExtractor()

	// Sequential path
	_, err = ext.Extract(ctx, scanResult, &ExtractOptions{Concurrency: 1})
	if err == nil {
		t.Fatal("expected error for canceled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got: %v", err)
	}

	// Concurrent path
	_, err = ext.Extract(ctx, scanResult, &ExtractOptions{Concurrency: 2})
	if err == nil {
		t.Fatal("expected error for canceled context (concurrent), got nil")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled (concurrent), got: %v", err)
	}
}

func TestExtractDeadlineExceededContext(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.go")
	testCode := "package test\n\nfunc Foo() {}\n"
	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	scanResult := &scanner.ScanResult{
		Files: []scanner.FileEntry{
			{Path: testFile, Language: "go", Size: int64(len(testCode))},
		},
	}

	// Create an already-expired deadline
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancel()

	ext := NewDefaultFileExtractor()

	_, err = ext.Extract(ctx, scanResult, &ExtractOptions{Concurrency: 1})
	if err == nil {
		t.Fatal("expected error for expired deadline, got nil")
	}
	if err != context.DeadlineExceeded {
		t.Errorf("expected context.DeadlineExceeded, got: %v", err)
	}
}

func TestExtractConcurrentCancelReturnsPromptly(t *testing.T) {
	// Verify that concurrent Extract returns promptly when context is
	// already cancelled. Using a pre-cancelled context guarantees the
	// cancellation code path is exercised.
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create multiple test files to ensure concurrent path is taken
	var files []scanner.FileEntry
	for i := 0; i < 10; i++ {
		p := filepath.Join(tmpDir, fmt.Sprintf("test%d.go", i))
		code := fmt.Sprintf("package test\n\nfunc Foo%d() {}\n", i)
		if err := os.WriteFile(p, []byte(code), 0644); err != nil {
			t.Fatal(err)
		}
		files = append(files, scanner.FileEntry{Path: p, Language: "go", Size: int64(len(code))})
	}

	scanResult := &scanner.ScanResult{Files: files}

	// Pre-cancel context so cancellation is guaranteed before goroutines start
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ext := NewDefaultFileExtractor()
	start := time.Now()
	_, err = ext.Extract(ctx, scanResult, &ExtractOptions{Concurrency: 4})
	elapsed := time.Since(start)

	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
	// Must return nearly instantly with a pre-cancelled context
	if elapsed > 5*time.Second {
		t.Errorf("Extract took too long after context cancellation: %v", elapsed)
	}
}

func TestIsBinaryContent(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{"empty", []byte{}, false},
		{"text", []byte("hello world"), false},
		{"binary with NUL", []byte("hel\x00lo"), true},
		{"UTF-16 LE BOM", []byte{0xFF, 0xFE, 'h', 0x00, 'i', 0x00}, false},
		{"UTF-16 BE BOM", []byte{0xFE, 0xFF, 0x00, 'h', 0x00, 'i'}, false},
		{"NUL without BOM", []byte{0x00, 0x01, 0x02}, true},
		{"UTF-8 BOM", []byte{0xEF, 0xBB, 0xBF, 'h', 'e', 'l', 'l', 'o'}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBinaryContent(tt.content)
			if got != tt.expected {
				t.Errorf("isBinaryContent(%q) = %v, want %v", tt.content, got, tt.expected)
			}
		})
	}
}

func TestHasUTF16BOM(t *testing.T) {
	tests := []struct {
		name     string
		content  []byte
		expected bool
	}{
		{"empty", []byte{}, false},
		{"single byte", []byte{0xFF}, false},
		{"UTF-16 LE", []byte{0xFF, 0xFE}, true},
		{"UTF-16 BE", []byte{0xFE, 0xFF}, true},
		{"not BOM", []byte{0xFE, 0xFE}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasUTF16BOM(tt.content)
			if got != tt.expected {
				t.Errorf("hasUTF16BOM(%v) = %v, want %v", tt.content, got, tt.expected)
			}
		})
	}
}
