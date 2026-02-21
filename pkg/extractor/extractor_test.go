package extractor

import (
	"os"
	"path/filepath"
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
