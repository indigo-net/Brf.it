package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_c_sharp "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/csharp"
)

// extractCSharpNames is a test helper that parses C# code and returns
// all captured @name values from the query matches.
func extractCSharpNames(t *testing.T, code []byte) map[string]bool {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_c_sharp.Language())
	parser.SetLanguage(lang)
	tree := parser.Parse(code, nil)
	defer tree.Close()
	query := NewCSharpQuery()
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

func TestCSharpQueryLanguage(t *testing.T) {
	query := NewCSharpQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestCSharpQueryPattern(t *testing.T) {
	query := NewCSharpQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_c_sharp.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
	defer q.Close()
}

func TestCSharpQueryImportPattern(t *testing.T) {
	query := NewCSharpQuery()
	pattern := query.ImportQuery()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty import query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_c_sharp.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid import query pattern: %v", err)
	}
	defer q.Close()
}

func TestCSharpQueryExtractFunction(t *testing.T) {
	code := []byte(`
public class Calculator {
    public int Add(int a, int b) {
        return a + b;
    }

    public async Task<int> AddAsync(int a, int b) => a + b;

    public static void DoStuff() { }

    public static T Parse<T>(string input) where T : struct {
        return default;
    }

    public static int Square(this int x) => x * x;

    public int Double(int x) => x * 2;
}
`)

	foundNames := extractCSharpNames(t, code)

	expectedNames := []string{"Add", "AddAsync", "DoStuff", "Parse", "Square", "Double"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find method '%s'", expected)
		}
	}
}

func TestCSharpQueryExtractTypes(t *testing.T) {
	code := []byte(`
public class User { }
public struct Point { }
public interface IDrawable { void Draw(); }
public enum Color { Red, Green, Blue }
public record Person(string Name, int Age);
public record class PersonClass(string Name);
public record struct Measurement(double Value);
public delegate void Action<T>(T item);
`)

	foundNames := extractCSharpNames(t, code)

	expectedNames := []string{
		"User", "Point", "IDrawable", "Color",
		"Person", "PersonClass", "Measurement", "Action",
	}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find type '%s'", expected)
		}
	}
}

func TestCSharpQueryExtractConstructorDestructor(t *testing.T) {
	code := []byte(`
public class MyClass {
    public MyClass() { }
    public MyClass(int value) { }
    static MyClass() { }
    ~MyClass() { }
}
`)

	foundNames := extractCSharpNames(t, code)

	if !foundNames["MyClass"] {
		t.Error("expected to find constructor/destructor 'MyClass'")
	}
}

func TestCSharpQueryExtractProperties(t *testing.T) {
	code := []byte(`
public class Config {
    public string Name { get; set; }
    public int ReadOnly => 42;
    public int Value { get; init; }
    public string Title { get; private set; }
}
`)

	foundNames := extractCSharpNames(t, code)

	expectedNames := []string{"Name", "ReadOnly", "Value", "Title"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find property '%s'", expected)
		}
	}
}

func TestCSharpQueryExtractFields(t *testing.T) {
	code := []byte(`
public class Config {
    public const int MaxValue = 100;
    public static readonly string Name = "Calc";
    public static int Counter = 0;
    private int _count;
}
`)

	foundNames := extractCSharpNames(t, code)

	expectedNames := []string{"MaxValue", "Name", "Counter", "_count"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find field '%s'", expected)
		}
	}
}

func TestCSharpQueryExtractEvents(t *testing.T) {
	code := []byte(`
public class EventSource {
    public event EventHandler Changed;
    public event EventHandler<int> ValueChanged { add {} remove {} }
}
`)

	foundNames := extractCSharpNames(t, code)

	expectedNames := []string{"Changed", "ValueChanged"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find event '%s'", expected)
		}
	}
}

func TestCSharpQueryExtractOperators(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_c_sharp.Language())
	parser.SetLanguage(lang)

	code := []byte(`
public class Vector {
    public static Vector operator +(Vector a, Vector b) => a;
    public static implicit operator int(Vector v) => 0;
    public static explicit operator string(Vector v) => "";
    public int this[int index] => index;
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCSharpQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), code)
	captureNames := q.CaptureNames()

	var foundOperator, foundConversion, foundIndexer bool
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, c := range match.Captures {
			if captureNames[c.Index] == "kind" {
				switch c.Node.Kind() {
				case "operator_declaration":
					foundOperator = true
				case "conversion_operator_declaration":
					foundConversion = true
				case "indexer_declaration":
					foundIndexer = true
				}
			}
		}
	}

	if !foundOperator {
		t.Error("expected to find operator declaration")
	}
	if !foundConversion {
		t.Error("expected to find conversion operator declaration")
	}
	if !foundIndexer {
		t.Error("expected to find indexer declaration")
	}
}

func TestCSharpQueryExtractImport(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_c_sharp.Language())
	parser.SetLanguage(lang)

	code := []byte(`
using System;
using System.Collections.Generic;
using static System.Math;
global using System.Linq;
using MyAlias = System.Text.StringBuilder;
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewCSharpQuery()
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

	if count < 5 {
		t.Errorf("expected at least 5 using directives, got %d", count)
	}
}

func TestCSharpQueryExtractGenerics(t *testing.T) {
	code := []byte(`
public class Container<T> where T : class {
    public T GetItem() { return default; }
}

public interface IRepository<T> where T : new() {
    T FindById(int id);
}

public class Utils {
    public static T Parse<T, U>(U input) where T : struct where U : class {
        return default;
    }
}
`)

	foundNames := extractCSharpNames(t, code)

	expectedNames := []string{"Container", "GetItem", "IRepository", "FindById", "Parse"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestCSharpQueryKindMapping(t *testing.T) {
	query := NewCSharpQuery()
	mapping := query.KindMapping()

	expectedMappings := map[string]string{
		"class_declaration":                 "class",
		"struct_declaration":                "struct",
		"interface_declaration":             "interface",
		"enum_declaration":                  "enum",
		"record_declaration":                "record",
		"delegate_declaration":              "type",
		"method_declaration":                "method",
		"constructor_declaration":           "constructor",
		"destructor_declaration":            "destructor",
		"property_declaration":              "variable",
		"field_declaration":                 "field",
		"event_declaration":                 "variable",
		"event_field_declaration":           "variable",
		"indexer_declaration":               "method",
		"operator_declaration":              "function",
		"conversion_operator_declaration":   "function",
		"namespace_declaration":             "namespace",
		"file_scoped_namespace_declaration": "namespace",
		"enum_member_declaration":           "variable",
	}

	for nodeType, expected := range expectedMappings {
		if actual, ok := mapping[nodeType]; !ok {
			t.Errorf("expected mapping for '%s'", nodeType)
		} else if actual != expected {
			t.Errorf("expected '%s' -> '%s', got '%s'", nodeType, expected, actual)
		}
	}
}

func TestCSharpQueryCaptures(t *testing.T) {
	query := NewCSharpQuery()
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

func TestCSharpQueryExtractNamespace(t *testing.T) {
	code := []byte(`
namespace MyApp {
    public class Foo { }
}

namespace Nested.Namespace {
    public class Bar { }
}
`)

	foundNames := extractCSharpNames(t, code)

	expectedNames := []string{"MyApp", "Foo", "Bar"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}

	// Test file-scoped namespace
	code2 := []byte(`namespace FileScoped;

public class MyClass { }
`)

	foundNames2 := extractCSharpNames(t, code2)

	if !foundNames2["FileScoped"] {
		t.Error("expected to find file-scoped namespace 'FileScoped'")
	}
	if !foundNames2["MyClass"] {
		t.Error("expected to find 'MyClass' in file-scoped namespace")
	}
}

func TestCSharpQueryExtractRecords(t *testing.T) {
	code := []byte(`
public record Person(string Name, int Age);
public record class Employee(string Name, string Department);
public record struct Measurement(double Value);
public record Wrapper<T>(T Value);
`)

	foundNames := extractCSharpNames(t, code)

	expectedNames := []string{"Person", "Employee", "Measurement", "Wrapper"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find record '%s'", expected)
		}
	}
}
