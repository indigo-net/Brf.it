package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"
)

func TestCQueryLanguage(t *testing.T) {
	query := NewCQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestCQueryPattern(t *testing.T) {
	query := NewCQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_c.Language())
	_, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
}

func TestCQueryExtractFunction(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_c.Language())
	parser.SetLanguage(lang)

	code := []byte(`// Add two integers
int add(int a, int b) {
    return a + b;
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

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

func TestCQueryExtractStruct(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_c.Language())
	parser.SetLanguage(lang)

	code := []byte(`// User structure
struct User {
    int id;
    char name[64];
};
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	foundUser := false
	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, c := range match.Captures {
			name := captureNames[c.Index]
			if name == "name" && string(code[c.Node.StartByte():c.Node.EndByte()]) == "User" {
				foundUser = true
			}
		}
	}

	if !foundUser {
		t.Fatal("expected to find struct 'User'")
	}
}

func TestCQueryExtractMacro(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_c.Language())
	parser.SetLanguage(lang)

	code := []byte(`#define MAX_SIZE 100
#define MIN(a, b) ((a) < (b) ? (a) : (b))
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCQuery()
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

	expectedNames := []string{"MAX_SIZE", "MIN"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find macro '%s'", expected)
		}
	}
}

func TestCQueryExtractEnum(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_c.Language())
	parser.SetLanguage(lang)

	code := []byte(`enum Color {
    RED,
    GREEN,
    BLUE
};
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	foundColor := false
	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, c := range match.Captures {
			name := captureNames[c.Index]
			if name == "name" && string(code[c.Node.StartByte():c.Node.EndByte()]) == "Color" {
				foundColor = true
			}
		}
	}

	if !foundColor {
		t.Fatal("expected to find enum 'Color'")
	}
}

func TestCQueryExtractTypedef(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_c.Language())
	parser.SetLanguage(lang)

	code := []byte(`typedef struct {
    int x;
    int y;
} Point;
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	foundPoint := false
	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, c := range match.Captures {
			name := captureNames[c.Index]
			if name == "name" && string(code[c.Node.StartByte():c.Node.EndByte()]) == "Point" {
				foundPoint = true
			}
		}
	}

	if !foundPoint {
		t.Fatal("expected to find typedef 'Point'")
	}
}
