package extractor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	scanResult, err := sc.Scan()
	if err != nil {
		t.Fatal(err)
	}

	// Extract signatures
	extractor := NewDefaultFileExtractor()
	result, err := extractor.Extract(scanResult, nil)
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

	result, err := extractor.Extract(scanResult, opts)
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

	result, err := extractor.Extract(scanResult, opts)
	if err != nil {
		t.Fatal(err)
	}

	if result.ErrorCount != 0 {
		t.Errorf("expected 0 errors when TOCTOU disabled, got %d", result.ErrorCount)
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

	result, err := extractor.Extract(scanResult, nil)
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
