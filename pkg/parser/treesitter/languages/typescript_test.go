package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

func TestTypeScriptQueryLanguage(t *testing.T) {
	query := NewTypeScriptQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestTypeScriptQueryPattern(t *testing.T) {
	query := NewTypeScriptQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_typescript.LanguageTypescript())
	_, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
}

func TestTypeScriptQueryExtractFunction(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_typescript.LanguageTypescript())
	parser.SetLanguage(lang)

	code := []byte(`
/**
 * Adds two numbers together.
 */
export function add(a: number, b: number): number {
  return a + b;
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewTypeScriptQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	// Iterate through all matches to find the function declaration
	captureNames := q.CaptureNames()
	var funcCaptures map[string]string
	for {
		match := matches.Next()
		if match == nil {
			break
		}

		captures := make(map[string]string)
		for _, c := range match.Captures {
			name := captureNames[c.Index]
			captures[name] = string(code[c.Node.StartByte():c.Node.EndByte()])
		}

		// Check if this is the add function
		if captures["name"] == "add" {
			funcCaptures = captures
			break
		}
	}

	if funcCaptures == nil {
		t.Fatal("expected to find function 'add'")
	}

	if funcCaptures["name"] != "add" {
		t.Errorf("expected name 'add', got '%s'", funcCaptures["name"])
	}
}

func TestTypeScriptQueryExtractModuleLevelVariables(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_typescript.LanguageTypescript())
	parser.SetLanguage(lang)

	code := []byte(`const API_URL = "https://api.example.com";
export const MAX_RETRIES = 3;
let counter = 0;
const arrowFn = () => {};

function test() {
    const localVar = 1;
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewTypeScriptQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	foundNames := make(map[string]bool)

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, c := range match.Captures {
			name := captureNames[c.Index]
			if name == "name" {
				foundNames[string(code[c.Node.StartByte():c.Node.EndByte()])] = true
			}
		}
	}

	// Module-level variables should be found
	expectedNames := []string{"API_URL", "MAX_RETRIES", "counter", "arrowFn", "test"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}

	// Local variable should NOT be found (filtered by program wrapper)
	if foundNames["localVar"] {
		t.Error("should not find local variable 'localVar'")
	}
}
