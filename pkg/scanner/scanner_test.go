package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewFileScanner(t *testing.T) {
	opts := DefaultScanOptions()
	opts.RootPath = "."

	scanner, err := NewFileScanner(opts)
	if err != nil {
		t.Fatalf("NewFileScanner returned error: %v", err)
	}

	if scanner == nil {
		t.Fatal("expected scanner to be non-nil")
	}

	if scanner.opts != opts {
		t.Error("expected opts to be set")
	}
}

func TestNewFileScannerNilOptions(t *testing.T) {
	scanner, err := NewFileScanner(nil)
	if err != nil {
		t.Fatalf("NewFileScanner with nil options returned error: %v", err)
	}

	if scanner.opts == nil {
		t.Fatal("expected default options to be set")
	}
}

func TestFileEntryDefaults(t *testing.T) {
	entry := FileEntry{}

	if entry.Path != "" {
		t.Errorf("expected empty Path, got '%s'", entry.Path)
	}
	if entry.Language != "" {
		t.Errorf("expected empty Language, got '%s'", entry.Language)
	}
	if entry.Size != 0 {
		t.Errorf("expected zero Size, got %d", entry.Size)
	}
}

func TestScanOptionsDefaults(t *testing.T) {
	opts := DefaultScanOptions()

	expectedExts := map[string]string{
		".go":  "go",
		".ts":  "typescript",
		".tsx": "typescript",
		".js":  "javascript",
		".jsx": "javascript",
	}

	for ext, lang := range expectedExts {
		if gotLang, ok := opts.SupportedExtensions[ext]; !ok {
			t.Errorf("expected extension '%s' to be supported", ext)
		} else if gotLang != lang {
			t.Errorf("expected extension '%s' to map to '%s', got '%s'", ext, lang, gotLang)
		}
	}

	if opts.IncludeHidden {
		t.Error("expected IncludeHidden to be false by default")
	}

	const expectedMaxSize = 512000
	if opts.MaxFileSize != expectedMaxSize {
		t.Errorf("expected MaxFileSize %d, got %d", expectedMaxSize, opts.MaxFileSize)
	}
}

func TestScanOptionsWithExtensions(t *testing.T) {
	opts := DefaultScanOptions()

	// Test extension detection
	tests := []struct {
		path      string
		expected  string
		supported bool
	}{
		{"main.go", "go", true},
		{"app.ts", "typescript", true},
		{"component.tsx", "typescript", true},
		{"index.js", "javascript", true},
		{"App.jsx", "javascript", true},
		{"README.md", "", false},
		{"config.json", "", false},
		{"style.css", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			lang, ok := opts.GetLanguage(tt.path)
			if ok != tt.supported {
				t.Errorf("expected supported %v, got %v", tt.supported, ok)
			}
			if lang != tt.expected {
				t.Errorf("expected language '%s', got '%s'", tt.expected, lang)
			}
		})
	}
}

func TestScannerInterface(t *testing.T) {
	// Verify FileScanner implements Scanner interface
	var _ Scanner = (*FileScanner)(nil)
}

func TestScanEmptyDirectory(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir

	scanner, _ := NewFileScanner(opts)
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if len(result.Files) != 0 {
		t.Errorf("expected 0 entries in empty directory, got %d", len(result.Files))
	}
}

func TestScanSingleFile(t *testing.T) {
	// Create temp directory with a Go file
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test Go file
	testFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(testFile, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir

	scanner, _ := NewFileScanner(opts)
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	if len(result.Files) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result.Files))
		return
	}

	if result.Files[0].Language != "go" {
		t.Errorf("expected language 'go', got '%s'", result.Files[0].Language)
	}
}

func TestScanFilterByExtension(t *testing.T) {
	// Create temp directory with multiple file types
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := []struct {
		name     string
		content  string
		expected bool
	}{
		{"main.go", "package main\n", true},
		{"app.ts", "const x = 1;\n", true},
		{"README.md", "# Test\n", false},
		{"config.json", "{}\n", false},
	}

	for _, f := range files {
		path := filepath.Join(tmpDir, f.name)
		if err := os.WriteFile(path, []byte(f.content), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir

	scanner, _ := NewFileScanner(opts)
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	// Should only get .go and .ts files
	if len(result.Files) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result.Files))
		for _, e := range result.Files {
			t.Logf("  - %s (%s)", e.Path, e.Language)
		}
	}
}

func TestScanExcludeHidden(t *testing.T) {
	// Create temp directory with hidden files
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create visible and hidden files
	visibleFile := filepath.Join(tmpDir, "visible.go")
	hiddenFile := filepath.Join(tmpDir, ".hidden.go")
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	hiddenDirFile := filepath.Join(hiddenDir, "inside.go")

	if err := os.WriteFile(visibleFile, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("failed to create visible file: %v", err)
	}
	if err := os.WriteFile(hiddenFile, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("failed to create hidden file: %v", err)
	}
	if err := os.Mkdir(hiddenDir, 0755); err != nil {
		t.Fatalf("failed to create hidden dir: %v", err)
	}
	if err := os.WriteFile(hiddenDirFile, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("failed to create file in hidden dir: %v", err)
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir
	opts.IncludeHidden = false

	scanner, _ := NewFileScanner(opts)
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	// Should only get visible.go
	if len(result.Files) != 1 {
		t.Errorf("expected 1 entry (hidden files excluded), got %d", len(result.Files))
		for _, e := range result.Files {
			t.Logf("  - %s", e.Path)
		}
	}
}

func TestScanIncludeHidden(t *testing.T) {
	// Create temp directory with hidden files
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create visible and hidden files
	visibleFile := filepath.Join(tmpDir, "visible.go")
	hiddenFile := filepath.Join(tmpDir, ".hidden.go")

	if err := os.WriteFile(visibleFile, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("failed to create visible file: %v", err)
	}
	if err := os.WriteFile(hiddenFile, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("failed to create hidden file: %v", err)
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir
	opts.IncludeHidden = true

	scanner, _ := NewFileScanner(opts)
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	// Should get both files
	if len(result.Files) != 2 {
		t.Errorf("expected 2 entries (hidden files included), got %d", len(result.Files))
	}
}

func TestScanMaxFileSize(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create small and large files
	smallFile := filepath.Join(tmpDir, "small.go")
	largeFile := filepath.Join(tmpDir, "large.go")

	smallContent := []byte("package main\n")
	largeContent := make([]byte, 600000) // 600KB, exceeds default 500KB limit

	if err := os.WriteFile(smallFile, smallContent, 0644); err != nil {
		t.Fatalf("failed to create small file: %v", err)
	}
	if err := os.WriteFile(largeFile, largeContent, 0644); err != nil {
		t.Fatalf("failed to create large file: %v", err)
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir
	opts.MaxFileSize = 512000 // 500KB

	scanner, _ := NewFileScanner(opts)
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	// Should only get small file
	if len(result.Files) != 1 {
		t.Errorf("expected 1 entry (large file excluded), got %d", len(result.Files))
	}

	// Large file should be counted as skipped
	if result.SkippedCount < 1 {
		t.Errorf("expected at least 1 skipped file, got %d", result.SkippedCount)
	}
}

func TestScanGitignore(t *testing.T) {
	// Create temp directory with gitignore
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create .gitignore
	gitignore := filepath.Join(tmpDir, ".gitignore")
	if err := os.WriteFile(gitignore, []byte("ignored.go\nnode_modules/\n"), 0644); err != nil {
		t.Fatalf("failed to create gitignore: %v", err)
	}

	// Create files
	normalFile := filepath.Join(tmpDir, "normal.go")
	ignoredFile := filepath.Join(tmpDir, "ignored.go")
	nodeModulesDir := filepath.Join(tmpDir, "node_modules")
	nodeModulesFile := filepath.Join(nodeModulesDir, "index.js")

	if err := os.WriteFile(normalFile, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("failed to create normal file: %v", err)
	}
	if err := os.WriteFile(ignoredFile, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("failed to create ignored file: %v", err)
	}
	if err := os.Mkdir(nodeModulesDir, 0755); err != nil {
		t.Fatalf("failed to create node_modules: %v", err)
	}
	if err := os.WriteFile(nodeModulesFile, []byte("module.exports = {};\n"), 0644); err != nil {
		t.Fatalf("failed to create node_modules file: %v", err)
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir
	opts.IgnoreFile = gitignore

	scanner, _ := NewFileScanner(opts)
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	// Should only get normal.go
	if len(result.Files) != 1 {
		t.Errorf("expected 1 entry (gitignore applied), got %d", len(result.Files))
		for _, e := range result.Files {
			t.Logf("  - %s", e.Path)
		}
	}
}

func TestScanNestedDirectories(t *testing.T) {
	// Create nested directory structure
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create nested directories with files
	dirs := []string{
		"src",
		"src/components",
		"src/utils",
		"pkg",
		"pkg/scanner",
	}

	for _, dir := range dirs {
		path := filepath.Join(tmpDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
	}

	// Create files in each directory
	files := []string{
		"main.go",
		"src/app.ts",
		"src/components/Button.tsx",
		"src/utils/helper.ts",
		"pkg/scanner/scanner.go",
	}

	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		if err := os.WriteFile(path, []byte("// test\n"), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", file, err)
		}
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir

	scanner, _ := NewFileScanner(opts)
	result, err := scanner.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	// Should get all 5 files
	if len(result.Files) != 5 {
		t.Errorf("expected 5 entries, got %d", len(result.Files))
		for _, e := range result.Files {
			t.Logf("  - %s (%s)", e.Path, e.Language)
		}
	}

	// Verify total size is calculated
	if result.TotalSize == 0 {
		t.Error("expected TotalSize to be non-zero")
	}
}
