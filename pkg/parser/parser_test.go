package parser

import (
	"testing"
)

func TestSignatureDefaults(t *testing.T) {
	sig := Signature{}

	if sig.Name != "" {
		t.Errorf("expected empty Name, got '%s'", sig.Name)
	}
	if sig.Kind != "" {
		t.Errorf("expected empty Kind, got '%s'", sig.Kind)
	}
	if sig.Text != "" {
		t.Errorf("expected empty Text, got '%s'", sig.Text)
	}
	if sig.Doc != "" {
		t.Errorf("expected empty Doc, got '%s'", sig.Doc)
	}
	if sig.Line != 0 {
		t.Errorf("expected zero Line, got %d", sig.Line)
	}
}

func TestParseResultDefaults(t *testing.T) {
	result := ParseResult{}

	if result.Signatures != nil {
		t.Errorf("expected nil Signatures, got %v", result.Signatures)
	}
	if result.Error != nil {
		t.Errorf("expected nil Error, got %v", result.Error)
	}
}

func TestNodeKind(t *testing.T) {
	node := Node{
		Type:     "function_declaration",
		StartRow: 10,
		EndRow:   20,
	}

	if node.Type != "function_declaration" {
		t.Errorf("expected Type 'function_declaration', got '%s'", node.Type)
	}
}

func TestParserInterface(t *testing.T) {
	// Verify MockParser implements Parser interface
	var _ Parser = (*MockParser)(nil)
}

// MockParser is a mock implementation for testing.
type MockParser struct {
	signatures []Signature
	err        error
}

func (m *MockParser) Parse(content string, opts *Options) (*ParseResult, error) {
	if m.err != nil {
		return &ParseResult{Error: m.err}, m.err
	}
	return &ParseResult{
		Signatures: m.signatures,
	}, nil
}

func (m *MockParser) Languages() []string {
	return []string{"go", "typescript", "javascript"}
}

func TestMockParser(t *testing.T) {
	mock := &MockParser{
		signatures: []Signature{
			{Name: "Test", Kind: "function", Text: "func Test()", Line: 1},
		},
	}

	result, err := mock.Parse("package main", nil)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Signatures) != 1 {
		t.Errorf("expected 1 signature, got %d", len(result.Signatures))
	}

	if result.Signatures[0].Name != "Test" {
		t.Errorf("expected signature name 'Test', got '%s'", result.Signatures[0].Name)
	}
}

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	// Register mock parser
	mock := &MockParser{}
	registry.Register("go", mock)

	// Test Get
	parser, ok := registry.Get("go")
	if !ok {
		t.Fatal("expected to find parser for 'go'")
	}
	if parser != mock {
		t.Error("expected same parser instance")
	}

	// Test Get non-existent
	_, ok = registry.Get("python")
	if ok {
		t.Error("expected not to find parser for 'python'")
	}

	// Test Languages
	langs := registry.Languages()
	if len(langs) != 1 || langs[0] != "go" {
		t.Errorf("expected ['go'], got %v", langs)
	}
}

func TestDefaultRegistry(t *testing.T) {
	registry := DefaultRegistry()
	if registry == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"main.go", "go"},
		{"app.ts", "typescript"},
		{"component.tsx", "tsx"},
		{"index.js", "javascript"},
		{"App.jsx", "jsx"},
		{"README.md", ""},
		{"config.json", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			lang := DetectLanguage(tt.path)
			if lang != tt.expected {
				t.Errorf("expected language '%s', got '%s'", tt.expected, lang)
			}
		})
	}
}
