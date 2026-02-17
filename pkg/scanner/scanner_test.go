package scanner

import (
	"testing"
)

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
