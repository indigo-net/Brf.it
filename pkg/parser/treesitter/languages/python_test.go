package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

func TestPythonQueryLanguage(t *testing.T) {
	query := NewPythonQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestPythonQueryPattern(t *testing.T) {
	query := NewPythonQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_python.Language())
	_, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
}

func TestPythonQueryExtractFunction(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_python.Language())
	parser.SetLanguage(lang)

	code := []byte(`# Greet a user by name
def greet(name: str) -> str:
    """Return a greeting message."""
    return f"Hello, {name}!"
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewPythonQuery()
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

		if captures["name"] == "greet" {
			funcCaptures = captures
			break
		}
	}

	if funcCaptures == nil {
		t.Fatal("expected to find function 'greet'")
	}

	if funcCaptures["name"] != "greet" {
		t.Errorf("expected name 'greet', got '%s'", funcCaptures["name"])
	}
}

func TestPythonQueryExtractClass(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_python.Language())
	parser.SetLanguage(lang)

	code := []byte(`class User:
    """User model."""

    def __init__(self, name: str):
        self.name = name

    def greet(self) -> str:
        return f"Hello, {self.name}!"
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewPythonQuery()
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

	expectedNames := []string{"User", "__init__", "greet"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestPythonQueryExtractAsyncFunction(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_python.Language())
	parser.SetLanguage(lang)

	code := []byte(`async def fetch_data(url: str) -> dict:
    """Fetch data from URL."""
    pass
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewPythonQuery()
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

		if captures["name"] == "fetch_data" {
			funcCaptures = captures
			break
		}
	}

	if funcCaptures == nil {
		t.Fatal("expected to find async function 'fetch_data'")
	}

	if funcCaptures["name"] != "fetch_data" {
		t.Errorf("expected name 'fetch_data', got '%s'", funcCaptures["name"])
	}
}

func TestPythonQueryExtractModuleLevelVariables(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_python.Language())
	parser.SetLanguage(lang)

	code := []byte(`API_URL = "https://api.example.com"
MAX_RETRIES: int = 3

def test():
    local_var = 1

class Config:
    CLASS_VAR = "value"
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewPythonQuery()
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
	expectedNames := []string{"API_URL", "MAX_RETRIES", "test", "Config"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}

	// Local and class-level variables should NOT be found
	unexpectedNames := []string{"local_var", "CLASS_VAR"}
	for _, unexpected := range unexpectedNames {
		if foundNames[unexpected] {
			t.Errorf("should not find '%s'", unexpected)
		}
	}
}
