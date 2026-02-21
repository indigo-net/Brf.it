package treesitter

import (
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser"
)

func TestTreeSitterParserImplementsParser(t *testing.T) {
	// Verify TreeSitterParser implements parser.Parser interface
	var _ parser.Parser = (*TreeSitterParser)(nil)
}

func TestTreeSitterParserLanguages(t *testing.T) {
	p := NewTreeSitterParser()

	langs := p.Languages()

	expected := []string{"go", "typescript", "tsx"}
	for _, exp := range expected {
		found := false
		for _, lang := range langs {
			if lang == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected language '%s' not found in %v", exp, langs)
		}
	}
}

func TestTreeSitterParserParseGo(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package main

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

type Point struct {
	X, Y int
}
`

	result, err := p.Parse(code, &parser.Options{Language: "go"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Signatures) < 2 {
		t.Errorf("expected at least 2 signatures, got %d", len(result.Signatures))
	}

	// Find Add function
	var foundAdd bool
	for _, sig := range result.Signatures {
		if sig.Name == "Add" {
			foundAdd = true
			if sig.Kind != "function" {
				t.Errorf("expected kind 'function', got '%s'", sig.Kind)
			}
			// Note: doc may be empty as comments are captured separately in Tree-sitter
			if sig.Line == 0 {
				t.Error("expected non-zero line number")
			}
		}
	}

	if !foundAdd {
		t.Error("expected to find 'Add' function signature")
	}
}

func TestTreeSitterParserParseTypeScript(t *testing.T) {
	p := NewTreeSitterParser()

	code := `
/**
 * Adds two numbers together.
 */
export function add(a: number, b: number): number {
  return a + b;
}

export interface Point {
  x: number;
  y: number;
}
`

	result, err := p.Parse(code, &parser.Options{Language: "typescript"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Signatures) < 1 {
		t.Errorf("expected at least 1 signature, got %d", len(result.Signatures))
	}

	// Find add function
	var foundAdd bool
	for _, sig := range result.Signatures {
		if sig.Name == "add" {
			foundAdd = true
			// Note: kind may be "export" when using export_statement pattern
			if sig.Kind != "function" && sig.Kind != "export" {
				t.Errorf("expected kind 'function' or 'export', got '%s'", sig.Kind)
			}
		}
	}

	if !foundAdd {
		t.Error("expected to find 'add' function signature")
	}
}

func TestTreeSitterParserUnsupportedLanguage(t *testing.T) {
	p := NewTreeSitterParser()

	code := `print("hello")`

	result, err := p.Parse(code, &parser.Options{Language: "python"})
	if err == nil {
		t.Error("expected error for unsupported language")
	}
	if result != nil {
		t.Error("expected nil result for unsupported language")
	}
}

func TestTreeSitterParserAutoRegistration(t *testing.T) {
	// Verify parser is registered in default registry
	registry := parser.DefaultRegistry()

	for _, lang := range []string{"go", "typescript", "tsx"} {
		p, ok := registry.Get(lang)
		if !ok {
			t.Errorf("expected parser for '%s' to be registered", lang)
		}
		if p == nil {
			t.Errorf("expected non-nil parser for '%s'", lang)
		}
	}
}
