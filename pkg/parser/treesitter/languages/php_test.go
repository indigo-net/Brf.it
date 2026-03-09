package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_php "github.com/tree-sitter/tree-sitter-php/bindings/go"
)

// extractPHPNames is a test helper that parses PHP code.
func extractPHPNames(t *testing.T, code []byte) map[string]bool {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_php.LanguagePHP())
	parser.SetLanguage(lang)
	tree := parser.Parse(code, nil)
	defer tree.Close()
	query := NewPHPQuery()
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
			if captureNames[c.Index] == "name" {
				foundNames[string(code[c.Node.StartByte():c.Node.EndByte()])] = true
			}
		}
	}
	return foundNames
}

func TestPHPQueryLanguage(t *testing.T) {
	query := NewPHPQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestPHPQueryPattern(t *testing.T) {
	query := NewPHPQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	lang := sitter.NewLanguage(tree_sitter_php.LanguagePHP())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
	defer q.Close()
}

func TestPHPQueryImportPattern(t *testing.T) {
	query := NewPHPQuery()
	pattern := query.ImportQuery()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty import query pattern")
	}

	lang := sitter.NewLanguage(tree_sitter_php.LanguagePHP())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid import query pattern: %v", err)
	}
	defer q.Close()
}

func TestPHPQueryExtractFunction(t *testing.T) {
	code := []byte(`<?php
function greet($name) {
    return "Hello, " . $name;
}

function add($a, $b) {
    return $a + $b;
}

function calculate($x, $y, $op) {
    return $op === '+' ? $x + $y : $x - $y;
}
`)

	foundNames := extractPHPNames(t, code)

	expectedNames := []string{"greet", "add", "calculate"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find function '%s'", expected)
		}
	}
}

func TestPHPQueryExtractClass(t *testing.T) {
	code := []byte(`<?php
class User {
    private $name;

    public function getName() {
        return $this->name;
    }
}

class Product {
    public $title;
    public $price;
}
`)

	foundNames := extractPHPNames(t, code)

	expectedNames := []string{"User", "getName", "Product"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestPHPQueryExtractInterface(t *testing.T) {
	code := []byte(`<?php
interface Logger {
    public function log($message);
}

interface Repository {
    public function find($id);
    public function save($entity);
}
`)

	foundNames := extractPHPNames(t, code)

	expectedNames := []string{"Logger", "log", "Repository", "find", "save"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestPHPQueryExtractTrait(t *testing.T) {
	code := []byte(`<?php
trait Loggable {
    public function logActivity() {
        // implementation
    }
}

trait HasName {
    abstract protected function getName();
}
`)

	foundNames := extractPHPNames(t, code)

	expectedNames := []string{"Loggable", "logActivity", "HasName", "getName"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestPHPQueryExtractEnum(t *testing.T) {
	code := []byte(`<?php
enum Status {
    case Active;
    case Inactive;
}

enum Priority: int {
    case Low = 1;
    case High = 2;
}
`)

	foundNames := extractPHPNames(t, code)

	expectedNames := []string{"Status", "Priority"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find enum '%s'", expected)
		}
	}
}

func TestPHPQueryExtractVariable(t *testing.T) {
	code := []byte(`<?php
const MAX_SIZE = 100;
const APP_NAME = "Brf.it";
`)

	foundNames := extractPHPNames(t, code)

	expectedNames := []string{"MAX_SIZE", "APP_NAME"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find variable '%s'", expected)
		}
	}
}

func TestPHPQueryExtractImport(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_php.LanguagePHP())
	parser.SetLanguage(lang)

	code := []byte(`<?php
use App\Services\UserService;
use App\Models\Product;
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewPHPQuery()
	q, err := sitter.NewQuery(lang, string(query.ImportQuery()))
	if err != nil {
		t.Fatalf("failed to create import query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	count := 0
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		count++
	}

	// Should find import statements
	if count < 1 {
		t.Errorf("expected at least 1 import declaration, got %d", count)
	}
}

func TestPHPQueryKindMapping(t *testing.T) {
	query := NewPHPQuery()
	mapping := query.KindMapping()

	expectedMappings := map[string]string{
		"function_definition":       "function",
		"method_declaration":        "method",
		"class_declaration":         "class",
		"interface_declaration":     "interface",
		"trait_declaration":         "type",
		"enum_declaration":          "enum",
		"const_declaration":         "variable",
		"property_declaration":      "variable",
		"namespace_use_declaration": "import",
	}

	for nodeType, expected := range expectedMappings {
		if actual, ok := mapping[nodeType]; !ok {
			t.Errorf("expected mapping for '%s'", nodeType)
		} else if actual != expected {
			t.Errorf("expected '%s' -> '%s', got '%s'", nodeType, expected, actual)
		}
	}
}

func TestPHPQueryCaptures(t *testing.T) {
	query := NewPHPQuery()
	captures := query.Captures()

	expected := []string{"name", "signature", "doc", "kind"}
	if len(captures) != len(expected) {
		t.Errorf("expected %d captures, got %d", len(expected), len(captures))
	}

	for i, exp := range expected {
		if captures[i] != exp {
			t.Errorf("expected capture[%d] = '%s', got '%s'", i, exp, captures[i])
		}
	}
}
