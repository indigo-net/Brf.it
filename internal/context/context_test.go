package context

import (
	"fmt"
	"strings"
	"testing"

	"github.com/indigo-net/Brf.it/pkg/extractor"
	"github.com/indigo-net/Brf.it/pkg/formatter"
	"github.com/indigo-net/Brf.it/pkg/parser"
	"github.com/indigo-net/Brf.it/pkg/scanner"
	"github.com/indigo-net/Brf.it/pkg/tokenizer"
)

// mockScanner implements scanner.Scanner for testing
type mockScanner struct {
	result *scanner.ScanResult
	err    error
}

func (m *mockScanner) Scan() (*scanner.ScanResult, error) {
	return m.result, m.err
}

// mockExtractor implements extractor.Extractor for testing
type mockExtractor struct {
	result *extractor.ExtractResult
	err    error
}

func (m *mockExtractor) Extract(_ *scanner.ScanResult, _ *extractor.ExtractOptions) (*extractor.ExtractResult, error) {
	return m.result, m.err
}

func TestPackagerPackage(t *testing.T) {
	mockScan := &mockScanner{
		result: &scanner.ScanResult{
			Files: []scanner.FileEntry{
				{Path: "pkg/test.go", Language: "go", Size: 100},
			},
			TotalSize: 100,
		},
	}

	mockExt := &mockExtractor{
		result: &extractor.ExtractResult{
			Files: []extractor.ExtractedFile{
				{
					Path:     "pkg/test.go",
					Language: "go",
					Signatures: []parser.Signature{
						{
							Name:     "Add",
							Kind:     "function",
							Text:     "func Add(a, b int) int",
							Line:     5,
							Language: "go",
							Exported: true,
						},
					},
					Size: 100,
				},
			},
			TotalSignatures: 1,
			TotalSize:       100,
		},
	}

	formatters := map[string]formatter.Formatter{
		"xml":      formatter.NewXMLFormatter(),
		"markdown": formatter.NewMarkdownFormatter(),
	}

	p := NewPackager(mockScan, mockExt, formatters)

	result, err := p.Package(&Options{
		Path:        ".",
		Format:      "xml",
		IncludeTree: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.TotalSignatures != 1 {
		t.Errorf("expected 1 signature, got %d", result.TotalSignatures)
	}

	if result.TotalFiles != 1 {
		t.Errorf("expected 1 file, got %d", result.TotalFiles)
	}

	if !strings.Contains(string(result.Content), "<brfit>") {
		t.Error("expected XML output to contain <brfit>")
	}

	// TokenCount should be 0 (NoOpTokenizer by default)
	if result.TokenCount != 0 {
		t.Errorf("expected TokenCount 0 (NoOpTokenizer), got %d", result.TokenCount)
	}
}

func TestPackagerPackageMarkdown(t *testing.T) {
	mockScan := &mockScanner{
		result: &scanner.ScanResult{
			Files: []scanner.FileEntry{
				{Path: "test.go"},
			},
		},
	}

	mockExt := &mockExtractor{
		result: &extractor.ExtractResult{
			Files: []extractor.ExtractedFile{
				{
					Path:     "test.go",
					Language: "go",
					Signatures: []parser.Signature{
						{Text: "func Test()", Language: "go"},
					},
				},
			},
			TotalSignatures: 1,
		},
	}

	formatters := map[string]formatter.Formatter{
		"xml":      formatter.NewXMLFormatter(),
		"markdown": formatter.NewMarkdownFormatter(),
	}

	p := NewPackager(mockScan, mockExt, formatters)

	// Test with "md" - should be normalized to "markdown"
	result, err := p.Package(&Options{
		Path:        ".",
		Format:      "md",
		IncludeTree: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(result.Content), "# Code Summary") {
		t.Error("expected markdown header")
	}
}

func TestPackagerPackageMarkdownFull(t *testing.T) {
	mockScan := &mockScanner{
		result: &scanner.ScanResult{
			Files: []scanner.FileEntry{
				{Path: "test.go"},
			},
		},
	}

	mockExt := &mockExtractor{
		result: &extractor.ExtractResult{
			Files: []extractor.ExtractedFile{
				{
					Path:       "test.go",
					Language:   "go",
					Signatures: []parser.Signature{{Text: "func Test()", Language: "go"}},
				},
			},
			TotalSignatures: 1,
		},
	}

	formatters := map[string]formatter.Formatter{
		"xml":      formatter.NewXMLFormatter(),
		"markdown": formatter.NewMarkdownFormatter(),
	}

	p := NewPackager(mockScan, mockExt, formatters)

	// Test with "markdown" directly
	result, err := p.Package(&Options{
		Path:        ".",
		Format:      "markdown",
		IncludeTree: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(result.Content), "# Code Summary") {
		t.Error("expected markdown header")
	}
}

func TestPackagerUnknownFormat(t *testing.T) {
	mockScan := &mockScanner{
		result: &scanner.ScanResult{},
	}

	mockExt := &mockExtractor{
		result: &extractor.ExtractResult{},
	}

	formatters := map[string]formatter.Formatter{
		"xml": formatter.NewXMLFormatter(),
	}

	p := NewPackager(mockScan, mockExt, formatters)

	result, err := p.Package(&Options{
		Format: "unknown",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Should fallback to xml
	if !strings.Contains(string(result.Content), "<?xml") {
		t.Error("expected fallback to XML format")
	}
}

func TestPackagerSetTokenizer(t *testing.T) {
	mockScan := &mockScanner{
		result: &scanner.ScanResult{
			Files: []scanner.FileEntry{
				{Path: "test.go"},
			},
		},
	}

	mockExt := &mockExtractor{
		result: &extractor.ExtractResult{
			Files: []extractor.ExtractedFile{
				{
					Path:       "test.go",
					Language:   "go",
					Signatures: []parser.Signature{{Text: "func Test()", Language: "go"}},
				},
			},
			TotalSignatures: 1,
		},
	}

	formatters := map[string]formatter.Formatter{
		"xml": formatter.NewXMLFormatter(),
	}

	p := NewPackager(mockScan, mockExt, formatters)

	// Default: NoOpTokenizer
	result, err := p.Package(&Options{Format: "xml"})
	if err != nil {
		t.Fatal(err)
	}
	if result.TokenCount != 0 {
		t.Errorf("expected TokenCount 0 with NoOpTokenizer, got %d", result.TokenCount)
	}

	// Set nil tokenizer (should use NoOpTokenizer)
	p.SetTokenizer(nil)
	result, err = p.Package(&Options{Format: "xml"})
	if err != nil {
		t.Fatal(err)
	}
	if result.TokenCount != 0 {
		t.Errorf("expected TokenCount 0 with nil tokenizer, got %d", result.TokenCount)
	}
}

func TestPackagerWithTiktokenTokenizer(t *testing.T) {
	mockScan := &mockScanner{
		result: &scanner.ScanResult{
			Files: []scanner.FileEntry{
				{Path: "test.go"},
			},
		},
	}

	mockExt := &mockExtractor{
		result: &extractor.ExtractResult{
			Files: []extractor.ExtractedFile{
				{
					Path:       "test.go",
					Language:   "go",
					Signatures: []parser.Signature{{Text: "func Test()", Language: "go"}},
				},
			},
			TotalSignatures: 1,
		},
	}

	formatters := map[string]formatter.Formatter{
		"xml": formatter.NewXMLFormatter(),
	}

	p := NewPackager(mockScan, mockExt, formatters)

	// Set TiktokenTokenizer
	tt, err := tokenizer.NewTiktokenTokenizer()
	if err != nil {
		t.Fatalf("failed to create tiktoken tokenizer: %v", err)
	}
	p.SetTokenizer(tt)

	result, err := p.Package(&Options{Format: "xml"})
	if err != nil {
		t.Fatal(err)
	}

	// TokenCount should be > 0 with TiktokenTokenizer
	if result.TokenCount <= 0 {
		t.Errorf("expected TokenCount > 0 with TiktokenTokenizer, got %d", result.TokenCount)
	}
}

func TestPackagerTokenizerConsistency(t *testing.T) {
	mockScan := &mockScanner{
		result: &scanner.ScanResult{
			Files: []scanner.FileEntry{
				{Path: "test.go"},
			},
		},
	}

	mockExt := &mockExtractor{
		result: &extractor.ExtractResult{
			Files: []extractor.ExtractedFile{
				{
					Path:       "test.go",
					Language:   "go",
					Signatures: []parser.Signature{{Text: "func Test()", Language: "go"}},
				},
			},
			TotalSignatures: 1,
		},
	}

	formatters := map[string]formatter.Formatter{
		"xml": formatter.NewXMLFormatter(),
	}

	p := NewPackager(mockScan, mockExt, formatters)

	tt, err := tokenizer.NewTiktokenTokenizer()
	if err != nil {
		t.Fatalf("failed to create tiktoken tokenizer: %v", err)
	}
	p.SetTokenizer(tt)

	// Multiple calls should return consistent token counts
	result1, _ := p.Package(&Options{Format: "xml"})
	result2, _ := p.Package(&Options{Format: "xml"})

	if result1.TokenCount != result2.TokenCount {
		t.Errorf("inconsistent token counts: %d vs %d", result1.TokenCount, result2.TokenCount)
	}
}

func TestBuildTree(t *testing.T) {
	tests := []struct {
		name     string
		root     string
		paths    []string
		contains string
	}{
		{
			name:     "empty paths",
			root:     ".",
			paths:    []string{},
			contains: "",
		},
		{
			name:     "single file",
			root:     ".",
			paths:    []string{"main.go"},
			contains: "main.go",
		},
		{
			name:     "nested structure",
			root:     ".",
			paths:    []string{"pkg/scanner/scanner.go", "pkg/parser/parser.go"},
			contains: "pkg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildTree(tt.root, tt.paths)
			if tt.contains == "" {
				if result != "" {
					t.Errorf("expected empty string, got %q", result)
				}
			} else {
				if !strings.Contains(result, tt.contains) {
					t.Errorf("expected tree to contain %q, got %q", tt.contains, result)
				}
			}
		})
	}
}

func TestBuildTreeStructure(t *testing.T) {
	paths := []string{
		"pkg/scanner/scanner.go",
		"pkg/parser/parser.go",
	}

	result := BuildTree(".", paths)

	// Should contain pkg directory
	if !strings.Contains(result, "pkg") {
		t.Error("expected 'pkg' in tree")
	}

	// Should contain scanner.go
	if !strings.Contains(result, "scanner.go") {
		t.Error("expected 'scanner.go' in tree")
	}

	// Should contain parser.go
	if !strings.Contains(result, "parser.go") {
		t.Error("expected 'parser.go' in tree")
	}

	// Should use tree connectors
	if !strings.Contains(result, "├──") && !strings.Contains(result, "└──") {
		t.Error("expected tree connectors in output")
	}
}

func TestPackagerNoStdImportsPassthrough(t *testing.T) {
	mockScan := &mockScanner{
		result: &scanner.ScanResult{
			Files: []scanner.FileEntry{
				{Path: "main.go", Language: "go", Size: 100},
			},
			TotalSize: 100,
		},
	}

	mockExt := &mockExtractor{
		result: &extractor.ExtractResult{
			Files: []extractor.ExtractedFile{
				{
					Path:     "main.go",
					Language: "go",
					Signatures: []parser.Signature{
						{Name: "Main", Kind: "function", Text: "func Main()"},
					},
					RawImports: []string{
						`import (
	"fmt"
	"github.com/spf13/cobra"
)`,
					},
					Size: 100,
				},
			},
			TotalSignatures: 1,
			TotalSize:       100,
		},
	}

	formatters := map[string]formatter.Formatter{
		"xml":      formatter.NewXMLFormatter(),
		"markdown": formatter.NewMarkdownFormatter(),
	}

	p := NewPackager(mockScan, mockExt, formatters)

	// With IncludeImports=true, imports should be included verbatim
	result, err := p.Package(&Options{
		Path:           ".",
		Format:         "xml",
		IncludeTree:    false,
		IncludeImports: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	output := string(result.Content)
	if !strings.Contains(output, "fmt") {
		t.Error("expected import 'fmt' to be present")
	}
	if !strings.Contains(output, "cobra") {
		t.Error("expected import 'cobra' to be present")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	if opts.Format != "xml" {
		t.Errorf("expected default format 'xml', got %q", opts.Format)
	}

	if opts.MaxFileSize != 512000 {
		t.Errorf("expected MaxFileSize 512000, got %d", opts.MaxFileSize)
	}

	if !opts.IncludeTree {
		t.Error("expected IncludeTree to be true by default")
	}
}

func TestNormalizeFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"xml", "xml"},
		{"md", "markdown"},
		{"markdown", "markdown"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeFormat(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeFormat(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuildGlobalImports(t *testing.T) {
	tests := []struct {
		name           string
		files          []formatter.FileData
		expectedCount  int
		expectedImport string
		expectedFiles  int
	}{
		{
			name: "single file with imports",
			files: []formatter.FileData{
				{
					Path:       "a.go",
					RawImports: []string{`import "fmt"`, `import "os"`},
				},
			},
			expectedCount:  2,
			expectedImport: `import "fmt"`,
			expectedFiles:  1,
		},
		{
			name: "multiple files with overlapping imports",
			files: []formatter.FileData{
				{
					Path:       "a.go",
					RawImports: []string{`import "fmt"`, `import "os"`},
				},
				{
					Path:       "b.go",
					RawImports: []string{`import "fmt"`, `import "strings"`},
				},
			},
			expectedCount:  3,
			expectedImport: `import "fmt"`,
			expectedFiles:  2, // fmt appears in 2 files
		},
		{
			name: "file with error should be skipped",
			files: []formatter.FileData{
				{
					Path:       "a.go",
					RawImports: []string{`import "fmt"`},
				},
				{
					Path:       "b.go",
					Error:      fmt.Errorf("parse error"),
					RawImports: []string{`import "os"`},
				},
			},
			expectedCount:  1,
			expectedImport: `import "fmt"`,
			expectedFiles:  1,
		},
		{
			name:           "empty files",
			files:          []formatter.FileData{},
			expectedCount:  0,
			expectedImport: "",
			expectedFiles:  0,
		},
		{
			name: "duplicate imports in same file counted once",
			files: []formatter.FileData{
				{
					Path:       "a.go",
					RawImports: []string{`import "fmt"`, `import "fmt"`},
				},
			},
			expectedCount:  1,
			expectedImport: `import "fmt"`,
			expectedFiles:  1, // counted once per file
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildGlobalImports(tt.files)

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d imports, got %d", tt.expectedCount, len(result))
			}

			if tt.expectedImport != "" {
				found := false
				for _, ic := range result {
					if ic.Import == tt.expectedImport {
						found = true
						if ic.Count != tt.expectedFiles {
							t.Errorf("expected import %q count %d, got %d", tt.expectedImport, tt.expectedFiles, ic.Count)
						}
						break
					}
				}
				if !found {
					t.Errorf("expected import %q not found", tt.expectedImport)
				}
			}
		})
	}
}

func TestBuildGlobalImportsSorting(t *testing.T) {
	files := []formatter.FileData{
		{
			Path:       "a.go",
			RawImports: []string{`import "fmt"`},
		},
		{
			Path:       "b.go",
			RawImports: []string{`import "fmt"`},
		},
		{
			Path:       "c.go",
			RawImports: []string{`import "fmt"`, `import "os"`},
		},
	}

	result := buildGlobalImports(files)

	// Should be sorted by count descending
	if len(result) != 2 {
		t.Fatalf("expected 2 imports, got %d", len(result))
	}

	// fmt appears in 3 files, os in 1
	if result[0].Import != `import "fmt"` || result[0].Count != 3 {
		t.Errorf("expected fmt with count 3 first, got %s with count %d", result[0].Import, result[0].Count)
	}

	if result[1].Import != `import "os"` || result[1].Count != 1 {
		t.Errorf("expected os with count 1 second, got %s with count %d", result[1].Import, result[1].Count)
	}
}
