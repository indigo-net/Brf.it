package languages

import (
	"testing"
	"unsafe"

	tree_sitter_yaml "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/yaml"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// Helper: extract names captured by @name from YAML query matches.
func extractYAMLNames(t *testing.T, code string) []string {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(unsafe.Pointer(tree_sitter_yaml.Language()))
	parser.SetLanguage(lang)

	src := []byte(code)
	tree := parser.Parse(src, nil)
	defer tree.Close()

	query := NewYAMLQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), src)

	captureNames := q.CaptureNames()
	var names []string
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, c := range match.Captures {
			if captureNames[c.Index] == "name" {
				names = append(names, string(src[c.Node.StartByte():c.Node.EndByte()]))
			}
		}
	}
	return names
}

func TestYAMLQueryLanguage(t *testing.T) {
	q := NewYAMLQuery()
	if q.Language() == nil {
		t.Fatal("Language() returned nil")
	}
}

func TestYAMLQueryPattern(t *testing.T) {
	q := NewYAMLQuery()
	lang := q.Language()
	_, err := sitter.NewQuery(lang, string(q.Query()))
	if err != nil {
		t.Fatalf("failed to compile YAML query pattern: %v", err)
	}
}

func TestYAMLQueryImportPattern(t *testing.T) {
	q := NewYAMLQuery()
	if q.ImportQuery() != nil {
		t.Fatal("YAML ImportQuery() should return nil")
	}
}

func TestYAMLQueryExtractSimpleKeys(t *testing.T) {
	code := `
name: myapp
version: "1.0"
description: A sample application
`
	names := extractYAMLNames(t, code)
	expected := map[string]bool{"name": false, "version": false, "description": false}
	for _, name := range names {
		if _, ok := expected[name]; ok {
			expected[name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("expected to find key '%s'", name)
		}
	}
}

func TestYAMLQueryExtractNestedKeys(t *testing.T) {
	code := `
database:
  host: localhost
  port: 5432
`
	names := extractYAMLNames(t, code)
	expected := map[string]bool{"database": false, "host": false, "port": false}
	for _, name := range names {
		if _, ok := expected[name]; ok {
			expected[name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("expected to find key '%s'", name)
		}
	}
}

func TestYAMLQueryExtractComments(t *testing.T) {
	code := `
# This is a comment
name: myapp
# Another comment
version: "1.0"
`
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(unsafe.Pointer(tree_sitter_yaml.Language()))
	parser.SetLanguage(lang)

	src := []byte(code)
	tree := parser.Parse(src, nil)
	defer tree.Close()

	query := NewYAMLQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), src)

	captureNames := q.CaptureNames()
	docCount := 0
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, c := range match.Captures {
			if captureNames[c.Index] == "doc" {
				docCount++
			}
		}
	}
	if docCount < 2 {
		t.Errorf("expected at least 2 comment captures, got %d", docCount)
	}
}

func TestYAMLQueryKindMapping(t *testing.T) {
	q := NewYAMLQuery()
	km := q.KindMapping()

	if got, ok := km["block_mapping_pair"]; !ok {
		t.Error("missing kind mapping for 'block_mapping_pair'")
	} else if got != "variable" {
		t.Errorf("kind mapping for 'block_mapping_pair': got '%s', want 'variable'", got)
	}
}

func TestYAMLQueryCaptures(t *testing.T) {
	q := NewYAMLQuery()
	captures := q.Captures()

	expected := map[string]bool{
		"name":      false,
		"signature": false,
		"doc":       false,
		"kind":      false,
	}
	for _, c := range captures {
		expected[c] = true
	}
	for name, found := range expected {
		if !found {
			t.Errorf("missing capture '%s'", name)
		}
	}
}
