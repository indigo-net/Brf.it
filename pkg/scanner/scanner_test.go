package scanner

import (
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetBaseName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{"normal file", "path/to/file.go", "file.go"},
		{"hidden file", "path/to/.hidden", ".hidden"},
		{"root path", ".", ""},
		{"empty path", "", ""},
		{"normal directory", "src/pkg", "pkg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getBaseName(tt.path)
			// Note: filepath.Base(".") returns "." which getBaseName converts to ""
			if result != tt.expected {
				t.Errorf("getBaseName(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

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
	result, err := scanner.Scan(context.Background())
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
	result, err := scanner.Scan(context.Background())
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
	result, err := scanner.Scan(context.Background())
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
	result, err := scanner.Scan(context.Background())
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
	result, err := scanner.Scan(context.Background())
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
	result, err := scanner.Scan(context.Background())
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
	opts.IgnoreFiles = []string{gitignore}

	scanner, _ := NewFileScanner(opts)
	result, err := scanner.Scan(context.Background())
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

func TestScanGitignoreLoadFailureWarning(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(testFile, []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("user-specified ignore file missing warns", func(t *testing.T) {
		opts := DefaultScanOptions()
		opts.RootPath = tmpDir
		opts.IgnoreFiles = []string{filepath.Join(tmpDir, "nonexistent-gitignore")}

		sc, err := NewFileScanner(opts)
		if err != nil {
			t.Fatal(err)
		}

		var buf bytes.Buffer
		sc.logger = log.New(&buf, "[brfit] ", 0)

		_, err = sc.Scan(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(buf.String(), "WARN") || !strings.Contains(buf.String(), "ignore file") {
			t.Errorf("expected warning about ignore file failure, got: %q", buf.String())
		}
	})

	t.Run("default gitignore missing does not warn", func(t *testing.T) {
		opts := DefaultScanOptions()
		opts.RootPath = tmpDir
		// IgnoreFiles defaults to [".gitignore"] which doesn't exist in tmpDir

		sc, err := NewFileScanner(opts)
		if err != nil {
			t.Fatal(err)
		}

		var buf bytes.Buffer
		sc.logger = log.New(&buf, "[brfit] ", 0)

		_, err = sc.Scan(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		if buf.Len() > 0 {
			t.Errorf("expected no warning for default missing .gitignore, got: %q", buf.String())
		}
	})

	t.Run("warning emitted only once on repeated Scan calls", func(t *testing.T) {
		opts := DefaultScanOptions()
		opts.RootPath = tmpDir
		opts.IgnoreFiles = []string{filepath.Join(tmpDir, "nonexistent-gitignore")}

		sc, err := NewFileScanner(opts)
		if err != nil {
			t.Fatal(err)
		}

		var buf bytes.Buffer
		sc.logger = log.New(&buf, "[brfit] ", 0)

		_, _ = sc.Scan(context.Background())
		_, _ = sc.Scan(context.Background())

		warnCount := strings.Count(buf.String(), "WARN")
		if warnCount != 1 {
			t.Errorf("expected exactly 1 warning, got %d: %q", warnCount, buf.String())
		}
	})
}

func TestScanWalkDirPermissionDenied(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping permission test: running as root")
	}

	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	normalFile := filepath.Join(tmpDir, "normal.go")
	if err := os.WriteFile(normalFile, []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}

	noAccessDir := filepath.Join(tmpDir, "noaccess")
	if err := os.Mkdir(noAccessDir, 0755); err != nil {
		t.Fatal(err)
	}
	innerFile := filepath.Join(noAccessDir, "inner.go")
	if err := os.WriteFile(innerFile, []byte("package inner\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(noAccessDir, 0000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(noAccessDir, 0755)

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir

	sc, err := NewFileScanner(opts)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	sc.logger = log.New(&buf, "[brfit] ", 0)

	result, err := sc.Scan(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Files) < 1 {
		t.Error("expected at least 1 file")
	}

	hasWarning := false
	for _, w := range result.Warnings {
		if strings.Contains(w, "permission denied") {
			hasWarning = true
			break
		}
	}
	if !hasWarning {
		t.Errorf("expected permission denied warning in result.Warnings, got: %v", result.Warnings)
	}
}

func TestScanSymlinkSkip(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	normalFile := filepath.Join(tmpDir, "normal.go")
	if err := os.WriteFile(normalFile, []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}

	symlinkFile := filepath.Join(tmpDir, "link.go")
	if err := os.Symlink(normalFile, symlinkFile); err != nil {
		t.Skip("symlinks not supported on this platform")
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir

	sc, err := NewFileScanner(opts)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	sc.logger = log.New(&buf, "[brfit] ", 0)

	result, err := sc.Scan(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Files) != 1 {
		t.Errorf("expected 1 file (symlink excluded), got %d", len(result.Files))
		for _, f := range result.Files {
			t.Logf("  - %s", f.Path)
		}
	}

	hasSymlinkWarning := false
	for _, w := range result.Warnings {
		if strings.Contains(w, "symlink") {
			hasSymlinkWarning = true
			break
		}
	}
	if !hasSymlinkWarning {
		t.Errorf("expected symlink warning in result.Warnings, got: %v", result.Warnings)
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
	result, err := scanner.Scan(context.Background())
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

func TestLogOutputNoDoubleNewline(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file that exceeds MaxFileSize to trigger "skipping large file" warning.
	// Use valid Go source content repeated to reach 1KB+ size.
	largeFile := filepath.Join(tmpDir, "large.go")
	data := bytes.Repeat([]byte("package main\n"), 100)
	if err := os.WriteFile(largeFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir
	opts.MaxFileSize = 512 // Set low limit so 1KB file triggers warning

	sc, err := NewFileScanner(opts)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	sc.logger = log.New(&buf, "", 0)

	_, err = sc.Scan(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	output := buf.String()
	if !strings.Contains(output, "skipping large file") {
		t.Fatalf("expected 'skipping large file' warning, got: %q", output)
	}

	if strings.Contains(output, "\n\n") {
		t.Errorf("log output contains double newline (Printf format string should not end with \\n): %q", output)
	}
}

func TestScanMultipleIgnoreFiles(t *testing.T) {
	// Verify that two separate ignore files both take effect
	tmpDir, err := os.MkdirTemp("", "brfit-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source files and files that should be ignored
	files := map[string]string{
		"main.go":    "package main\n",
		"app.ts":     "const x = 1;\n",
		"debug.log":  "some log output\n",
		"cache.tmp":  "temporary data\n",
		"server.log": "another log\n",
		"data.tmp":   "more temp data\n",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", name, err)
		}
	}

	// Create two separate ignore files
	ignoreLog := filepath.Join(tmpDir, ".ignore-logs")
	if err := os.WriteFile(ignoreLog, []byte("*.log\n"), 0644); err != nil {
		t.Fatalf("failed to create ignore-logs: %v", err)
	}

	ignoreTmp := filepath.Join(tmpDir, ".ignore-tmp")
	if err := os.WriteFile(ignoreTmp, []byte("*.tmp\n"), 0644); err != nil {
		t.Fatalf("failed to create ignore-tmp: %v", err)
	}

	opts := DefaultScanOptions()
	opts.RootPath = tmpDir
	opts.IgnoreFiles = []string{ignoreLog, ignoreTmp}

	scanner, err := NewFileScanner(opts)
	if err != nil {
		t.Fatalf("NewFileScanner returned error: %v", err)
	}

	result, err := scanner.Scan(context.Background())
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}

	// Should only get main.go and app.ts (*.log and *.tmp excluded)
	if len(result.Files) != 2 {
		t.Errorf("expected 2 entries (*.log and *.tmp excluded), got %d", len(result.Files))
		for _, e := range result.Files {
			t.Logf("  - %s (%s)", e.Path, e.Language)
		}
	}

	// Verify none of the ignored extensions are present
	for _, f := range result.Files {
		base := filepath.Base(f.Path)
		if strings.HasSuffix(base, ".log") {
			t.Errorf("expected .log files to be excluded, found: %s", f.Path)
		}
		if strings.HasSuffix(base, ".tmp") {
			t.Errorf("expected .tmp files to be excluded, found: %s", f.Path)
		}
	}
}

func TestFilepathBaseEdgeCases(t *testing.T) {
	// Document expected behavior of filepath.Base for edge cases
	tests := []struct {
		path     string
		expected string
	}{
		{".", "."},
		{"", "."},
		{"dir", "dir"},
		{"dir/subdir", "subdir"},
		{"/path/to/file", "file"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := filepath.Base(tt.path); got != tt.expected {
				t.Errorf("filepath.Base(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

func TestScanIncludePatterns(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure:
	// tmpDir/
	//   pkg/scanner/scanner.go
	//   pkg/scanner/scanner_test.go
	//   cmd/main.go
	//   README.md (unsupported)
	dirs := []string{
		filepath.Join(tmpDir, "pkg", "scanner"),
		filepath.Join(tmpDir, "cmd"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
	}
	files := map[string]string{
		"pkg/scanner/scanner.go":      "package scanner",
		"pkg/scanner/scanner_test.go": "package scanner",
		"cmd/main.go":                 "package main",
	}
	for rel, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, rel), []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	defaultOpts := DefaultScanOptions()

	tests := []struct {
		name            string
		includePatterns []string
		excludePatterns []string
		wantFiles       int
		wantPaths       []string
	}{
		{
			name:      "no patterns - all files",
			wantFiles: 3,
		},
		{
			name:            "include only pkg/**/*.go",
			includePatterns: []string{"pkg/**/*.go"},
			wantFiles:       2,
			wantPaths:       []string{"pkg/scanner/scanner.go", "pkg/scanner/scanner_test.go"},
		},
		{
			name:            "exclude test files",
			excludePatterns: []string{"**/*_test.go"},
			wantFiles:       2,
			wantPaths:       []string{"pkg/scanner/scanner.go", "cmd/main.go"},
		},
		{
			name:            "include pkg + exclude tests",
			includePatterns: []string{"pkg/**/*.go"},
			excludePatterns: []string{"**/*_test.go"},
			wantFiles:       1,
			wantPaths:       []string{"pkg/scanner/scanner.go"},
		},
		{
			name:            "include cmd only",
			includePatterns: []string{"cmd/**/*.go"},
			wantFiles:       1,
			wantPaths:       []string{"cmd/main.go"},
		},
		{
			name:            "no match",
			includePatterns: []string{"nonexistent/**"},
			wantFiles:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &ScanOptions{
				RootPath:            tmpDir,
				SupportedExtensions: defaultOpts.SupportedExtensions,
				MaxFileSize:         defaultOpts.MaxFileSize,
				IncludePatterns:     tt.includePatterns,
				ExcludePatterns:     tt.excludePatterns,
			}

			s, err := NewFileScanner(opts)
			if err != nil {
				t.Fatalf("NewFileScanner: %v", err)
			}

			result, err := s.Scan(context.Background())
			if err != nil {
				t.Fatalf("Scan: %v", err)
			}

			if len(result.Files) != tt.wantFiles {
				paths := make([]string, len(result.Files))
				for i, f := range result.Files {
					paths[i] = f.Path
				}
				t.Errorf("got %d files %v, want %d", len(result.Files), paths, tt.wantFiles)
			}

			if tt.wantPaths != nil {
				for _, wantRel := range tt.wantPaths {
					wantAbs := filepath.Join(tmpDir, wantRel)
					found := false
					for _, f := range result.Files {
						if f.Path == wantAbs {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected file %s not found in results", wantRel)
					}
				}
			}
		})
	}
}

func TestScanExcludeDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	dirs := []string{
		filepath.Join(tmpDir, "src"),
		filepath.Join(tmpDir, "vendor", "lib"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatalf("MkdirAll: %v", err)
		}
	}
	files := map[string]string{
		"src/main.go":        "package main",
		"vendor/lib/util.go": "package lib",
	}
	for rel, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, rel), []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	defaultOpts := DefaultScanOptions()
	opts := &ScanOptions{
		RootPath:            tmpDir,
		SupportedExtensions: defaultOpts.SupportedExtensions,
		MaxFileSize:         defaultOpts.MaxFileSize,
		ExcludePatterns:     []string{"vendor/**"},
	}

	s, err := NewFileScanner(opts)
	if err != nil {
		t.Fatalf("NewFileScanner: %v", err)
	}

	result, err := s.Scan(context.Background())
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}

	if len(result.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(result.Files))
	}
	if len(result.Files) > 0 && !strings.HasSuffix(result.Files[0].Path, "main.go") {
		t.Errorf("expected main.go, got %s", result.Files[0].Path)
	}
}

func TestScanSingleFileWithIncludePattern(t *testing.T) {
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(goFile, []byte("package main"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	defaultOpts := DefaultScanOptions()

	tests := []struct {
		name            string
		includePatterns []string
		wantFiles       int
	}{
		{
			name:      "no include pattern",
			wantFiles: 1,
		},
		{
			name:            "matching include pattern",
			includePatterns: []string{"**/*.go"},
			wantFiles:       1,
		},
		{
			name:            "matching simple pattern",
			includePatterns: []string{"*.go"},
			wantFiles:       1,
		},
		{
			name:            "non-matching include pattern",
			includePatterns: []string{"**/*.ts"},
			wantFiles:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &ScanOptions{
				RootPath:            goFile,
				SupportedExtensions: defaultOpts.SupportedExtensions,
				MaxFileSize:         defaultOpts.MaxFileSize,
				IncludePatterns:     tt.includePatterns,
			}

			s, err := NewFileScanner(opts)
			if err != nil {
				t.Fatalf("NewFileScanner: %v", err)
			}

			result, err := s.Scan(context.Background())
			if err != nil {
				t.Fatalf("Scan: %v", err)
			}

			if len(result.Files) != tt.wantFiles {
				t.Errorf("got %d files, want %d", len(result.Files), tt.wantFiles)
			}
		})
	}
}

func TestScanChangedFilesWhitelist(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	for _, f := range []struct {
		path    string
		content string
	}{
		{"main.go", "package main\nfunc main() {}"},
		{"util.go", "package main\nfunc util() {}"},
		{"lib.go", "package main\nfunc lib() {}"},
		{"test.py", "def test(): pass"},
	} {
		if err := os.WriteFile(filepath.Join(tmpDir, f.path), []byte(f.content), 0644); err != nil {
			t.Fatalf("write %s: %v", f.path, err)
		}
	}

	defaultOpts := DefaultScanOptions()

	tests := []struct {
		name         string
		changedFiles map[string]bool
		wantFiles    int
		wantNames    []string
	}{
		{
			name:         "nil whitelist includes all",
			changedFiles: nil,
			wantFiles:    4,
		},
		{
			name:         "whitelist with two Go files",
			changedFiles: map[string]bool{"main.go": true, "util.go": true},
			wantFiles:    2,
			wantNames:    []string{"main.go", "util.go"},
		},
		{
			name:         "whitelist with one file",
			changedFiles: map[string]bool{"test.py": true},
			wantFiles:    1,
			wantNames:    []string{"test.py"},
		},
		{
			name:         "empty whitelist excludes all",
			changedFiles: map[string]bool{},
			wantFiles:    0,
		},
		{
			name:         "whitelist with nonexistent file",
			changedFiles: map[string]bool{"nonexistent.go": true},
			wantFiles:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &ScanOptions{
				RootPath:            tmpDir,
				SupportedExtensions: defaultOpts.SupportedExtensions,
				MaxFileSize:         defaultOpts.MaxFileSize,
				ChangedFiles:        tt.changedFiles,
			}

			s, err := NewFileScanner(opts)
			if err != nil {
				t.Fatalf("NewFileScanner: %v", err)
			}

			result, err := s.Scan(context.Background())
			if err != nil {
				t.Fatalf("Scan: %v", err)
			}

			if len(result.Files) != tt.wantFiles {
				names := make([]string, len(result.Files))
				for i, f := range result.Files {
					names[i] = filepath.Base(f.Path)
				}
				t.Errorf("got %d files %v, want %d", len(result.Files), names, tt.wantFiles)
			}

			if tt.wantNames != nil {
				gotNames := make(map[string]bool)
				for _, f := range result.Files {
					gotNames[filepath.Base(f.Path)] = true
				}
				for _, name := range tt.wantNames {
					if !gotNames[name] {
						t.Errorf("expected file %q not found in results", name)
					}
				}
			}
		})
	}
}
