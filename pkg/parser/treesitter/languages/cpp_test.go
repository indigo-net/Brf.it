package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
)

func TestCppQueryLanguage(t *testing.T) {
	query := NewCppQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestCppQueryPattern(t *testing.T) {
	query := NewCppQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	_, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
}

func TestCppQueryExtractFunction(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`// Add two integers
int add(int a, int b) {
    return a + b;
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
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

func TestCppQueryExtractClass(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`// User class
class User {
public:
    int id;
    std::string name;

    void setName(const std::string& n);
};
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
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
		t.Fatal("expected to find class 'User'")
	}
}

func TestCppQueryExtractMethod(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`class Calculator {
public:
    int add(int a, int b);
    double multiply(double a, double b);
};
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
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

	expectedNames := []string{"Calculator", "add", "multiply"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestCppQueryExtractConstructorDestructor(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`class Resource {
public:
    Resource();
    ~Resource();
};
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
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

	// Should find class name and destructor name (both are "Resource")
	if !foundNames["Resource"] {
		t.Error("expected to find 'Resource' (class and/or destructor)")
	}
}

func TestCppQueryExtractNamespace(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`namespace utils {
    int helper();
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	foundUtils := false
	foundHelper := false
	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, c := range match.Captures {
			name := captureNames[c.Index]
			if name == "name" {
				text := string(code[c.Node.StartByte():c.Node.EndByte()])
				if text == "utils" {
					foundUtils = true
				}
				if text == "helper" {
					foundHelper = true
				}
			}
		}
	}

	if !foundUtils {
		t.Error("expected to find namespace 'utils'")
	}
	if !foundHelper {
		t.Error("expected to find function 'helper'")
	}
}

func TestCppQueryExtractTemplate(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`template<typename T>
class Box {
    T value;
};

template<typename T>
T getMax(T a, T b) {
    return (a > b) ? a : b;
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
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

	expectedNames := []string{"Box", "getMax"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find template '%s'", expected)
		}
	}
}

func TestCppQueryExtractStruct(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`struct Point {
    int x;
    int y;
};
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
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
		t.Fatal("expected to find struct 'Point'")
	}
}

func TestCppQueryExtractEnum(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`enum Color {
    RED,
    GREEN,
    BLUE
};

enum class Status {
    Pending,
    Active,
    Done
};
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
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

	expectedNames := []string{"Color", "Status"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find enum '%s'", expected)
		}
	}
}

func TestCppQueryExtractMacro(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`#define MAX_SIZE 100
#define MIN(a, b) ((a) < (b) ? (a) : (b))
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
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

func TestCppQueryExtractTypedef(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`typedef unsigned int uint;
typedef std::vector<int> IntVec;
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
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

	expectedNames := []string{"uint", "IntVec"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find typedef '%s'", expected)
		}
	}
}

func TestCppQueryExtractIncludes(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`#include <iostream>
#include <vector>
#include "myheader.h"
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
	q, err := sitter.NewQuery(lang, string(query.ImportQuery()))
	if err != nil {
		t.Fatalf("failed to create import query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	var imports []string
	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, c := range match.Captures {
			name := captureNames[c.Index]
			if name == "import_path" {
				imports = append(imports, string(code[c.Node.StartByte():c.Node.EndByte()]))
			}
		}
	}

	if len(imports) != 3 {
		t.Errorf("expected 3 includes, got %d", len(imports))
	}
}

func TestCppQueryNestedNamespaces(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`namespace outer {
    namespace inner {
        void helper();
    }
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
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

	expectedNames := []string{"outer", "inner", "helper"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestCppQueryMultipleInheritance(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`class Base1 {};
class Base2 {};
class Derived : public Base1, public Base2 {
    void method();
};
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
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

	expectedNames := []string{"Base1", "Base2", "Derived", "method"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestCppQueryEmptyFile(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(``)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	count := 0
	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, c := range match.Captures {
			name := captureNames[c.Index]
			if name == "name" {
				count++
			}
		}
	}

	if count != 0 {
		t.Errorf("expected 0 captures for empty file, got %d", count)
	}
}

func TestCppQueryOnlyComments(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_cpp.Language())
	parser.SetLanguage(lang)

	code := []byte(`// This is a comment
/* This is a
   multi-line comment */
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCppQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	var nameCount int
	var docCount int
	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, c := range match.Captures {
			name := captureNames[c.Index]
			if name == "name" {
				nameCount++
			}
			if name == "doc" {
				docCount++
			}
		}
	}

	if nameCount != 0 {
		t.Errorf("expected 0 name captures, got %d", nameCount)
	}
	if docCount != 2 {
		t.Errorf("expected 2 doc captures, got %d", docCount)
	}
}

func TestCppQueryKindMapping(t *testing.T) {
	query := NewCppQuery()
	mapping := query.KindMapping()

	expectedMappings := map[string]string{
		"function_definition":  "function",
		"class_specifier":      "class",
		"struct_specifier":     "struct",
		"enum_specifier":       "enum",
		"namespace_definition": "namespace",
		"template_declaration": "template",
		"field_declaration":    "method",
	}

	for nodeType, expectedKind := range expectedMappings {
		if kind, ok := mapping[nodeType]; !ok {
			t.Errorf("expected mapping for '%s'", nodeType)
		} else if kind != expectedKind {
			t.Errorf("expected kind '%s' for '%s', got '%s'", expectedKind, nodeType, kind)
		}
	}
}

func TestCppQueryCaptures(t *testing.T) {
	query := NewCppQuery()
	captures := query.Captures()

	expected := []string{"name", "signature", "doc", "kind"}
	if len(captures) != len(expected) {
		t.Errorf("expected %d captures, got %d", len(expected), len(captures))
	}

	for _, exp := range expected {
		found := false
		for _, cap := range captures {
			if cap == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected capture '%s' not found", exp)
		}
	}
}
