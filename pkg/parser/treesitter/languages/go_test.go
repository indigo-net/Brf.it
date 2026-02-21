package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

func TestGoQueryLanguage(t *testing.T) {
	query := NewGoQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestGoQueryPattern(t *testing.T) {
	query := NewGoQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_go.Language())
	_, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
}

func TestGoQueryExtractFunction(t *testing.T) {
	// Integration test: parse Go code and extract function
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_go.Language())
	parser.SetLanguage(lang)

	code := []byte(`package main

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewGoQuery()
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
	var funcKindNode *sitter.Node
	for {
		match := matches.Next()
		if match == nil {
			break
		}

		captures := make(map[string]string)
		var kindNode *sitter.Node
		for i, c := range match.Captures {
			name := captureNames[c.Index]
			captures[name] = string(code[c.Node.StartByte():c.Node.EndByte()])
			if name == "kind" {
				kindNode = &match.Captures[i].Node
			}
		}

		// Check if this is a function declaration
		if captures["name"] == "Add" {
			funcCaptures = captures
			funcKindNode = kindNode
			break
		}
	}

	if funcCaptures == nil {
		t.Fatal("expected to find function 'Add'")
	}

	if funcCaptures["name"] != "Add" {
		t.Errorf("expected name 'Add', got '%s'", funcCaptures["name"])
	}

	// Get the actual node type for kind
	if funcKindNode == nil {
		t.Fatal("expected kind node")
	}
	if funcKindNode.Kind() != "function_declaration" {
		t.Errorf("expected kind 'function_declaration', got '%s'", funcKindNode.Kind())
	}
}
