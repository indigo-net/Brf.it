package languages

import (
	"testing"
	"unsafe"

	tree_sitter_toml "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/toml"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// Helper: extract (name, kind) tuples from TOML query matches.
type tomlCapture struct {
	Name string
	Kind string
}

func extractTOMLCaptures(t *testing.T, code string) []tomlCapture {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(unsafe.Pointer(tree_sitter_toml.Language()))
	parser.SetLanguage(lang)

	src := []byte(code)
	tree := parser.Parse(src, nil)
	defer tree.Close()

	query := NewTOMLQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), src)

	captureNames := q.CaptureNames()
	var captures []tomlCapture
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		var tc tomlCapture
		for _, c := range match.Captures {
			text := string(src[c.Node.StartByte():c.Node.EndByte()])
			switch captureNames[c.Index] {
			case "name":
				tc.Name = text
			case "kind":
				tc.Kind = c.Node.Kind()
			}
		}
		if tc.Name != "" {
			captures = append(captures, tc)
		}
	}
	return captures
}

// Helper: extract just names from TOML query matches.
func extractTOMLNames(t *testing.T, code string) []string {
	t.Helper()
	captures := extractTOMLCaptures(t, code)
	var names []string
	for _, c := range captures {
		names = append(names, c.Name)
	}
	return names
}

func TestTOMLQueryLanguage(t *testing.T) {
	q := NewTOMLQuery()
	if q.Language() == nil {
		t.Fatal("Language() returned nil")
	}
}

func TestTOMLQueryPattern(t *testing.T) {
	q := NewTOMLQuery()
	lang := q.Language()
	_, err := sitter.NewQuery(lang, string(q.Query()))
	if err != nil {
		t.Fatalf("failed to compile TOML query pattern: %v", err)
	}
}

func TestTOMLQueryImportPattern(t *testing.T) {
	q := NewTOMLQuery()
	if q.ImportQuery() != nil {
		t.Fatal("TOML ImportQuery() should return nil")
	}
}

func TestTOMLQueryExtractTables(t *testing.T) {
	code := `
[package]
name = "myapp"

[dependencies]
serde = "1.0"
`
	captures := extractTOMLCaptures(t, code)
	tableCount := 0
	for _, c := range captures {
		if c.Kind == "table" {
			tableCount++
		}
	}
	if tableCount < 2 {
		t.Errorf("expected at least 2 table captures, got %d", tableCount)
	}
}

func TestTOMLQueryExtractTableArrays(t *testing.T) {
	code := `
[[bin]]
name = "cli"

[[bin]]
name = "server"
`
	captures := extractTOMLCaptures(t, code)
	arrayCount := 0
	for _, c := range captures {
		if c.Kind == "table_array_element" {
			arrayCount++
		}
	}
	if arrayCount < 2 {
		t.Errorf("expected at least 2 table_array_element captures, got %d", arrayCount)
	}
}

func TestTOMLQueryExtractPairs(t *testing.T) {
	code := `
name = "myapp"
version = "1.0.0"
edition = "2021"
`
	names := extractTOMLNames(t, code)
	expected := map[string]bool{"name": false, "version": false, "edition": false}
	for _, name := range names {
		if _, ok := expected[name]; ok {
			expected[name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("expected to find pair '%s'", name)
		}
	}
}

func TestTOMLQueryExtractDottedKey(t *testing.T) {
	code := `
[server.database]
host = "localhost"
`
	names := extractTOMLNames(t, code)
	found := false
	for _, name := range names {
		if name == "server.database" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected dotted key 'server.database', got %v", names)
	}
}

func TestTOMLQueryExtractComments(t *testing.T) {
	code := `
# This is a comment
[package]
# Another comment
name = "myapp"
`
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(unsafe.Pointer(tree_sitter_toml.Language()))
	parser.SetLanguage(lang)

	src := []byte(code)
	tree := parser.Parse(src, nil)
	defer tree.Close()

	query := NewTOMLQuery()
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

func TestTOMLQueryKindMapping(t *testing.T) {
	q := NewTOMLQuery()
	km := q.KindMapping()

	expectedKinds := map[string]string{
		"table":               "namespace",
		"table_array_element": "namespace",
		"pair":                "variable",
	}

	for nodeType, expectedKind := range expectedKinds {
		if got, ok := km[nodeType]; !ok {
			t.Errorf("missing kind mapping for '%s'", nodeType)
		} else if got != expectedKind {
			t.Errorf("kind mapping for '%s': got '%s', want '%s'", nodeType, got, expectedKind)
		}
	}
}

func TestTOMLQueryCaptures(t *testing.T) {
	q := NewTOMLQuery()
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
