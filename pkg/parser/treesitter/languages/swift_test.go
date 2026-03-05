package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_swift "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/swift"
)

func TestSwiftQueryLanguage(t *testing.T) {
	query := NewSwiftQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestSwiftQueryPattern(t *testing.T) {
	query := NewSwiftQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_swift.Language())
	_, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
}

func TestSwiftQueryImportPattern(t *testing.T) {
	query := NewSwiftQuery()
	pattern := query.ImportQuery()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty import query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_swift.Language())
	_, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid import query pattern: %v", err)
	}
}

func TestSwiftQueryExtractFunction(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_swift.Language())
	parser.SetLanguage(lang)

	code := []byte(`
/// Greets someone.
public func greet(name: String) -> String {
    return "Hello, \(name)"
}

func privateFunc() -> Bool {
    return true
}

public func asyncThrows() async throws -> Data {
    return Data()
}

@discardableResult
public func annotated() -> Bool {
    return true
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewSwiftQuery()
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

	expectedNames := []string{"greet", "privateFunc", "asyncThrows", "annotated"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find function '%s'", expected)
		}
	}
}

func TestSwiftQueryExtractTypes(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_swift.Language())
	parser.SetLanguage(lang)

	code := []byte(`
/// A point structure.
public struct Point {
    public let x: Double
    public let y: Double
}

public class Vehicle {
    var speed: Int = 0
}

public enum Direction {
    case north
    case south
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewSwiftQuery()
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

	expectedNames := []string{"Point", "Vehicle", "Direction"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find type '%s'", expected)
		}
	}
}

func TestSwiftQueryExtractProtocol(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_swift.Language())
	parser.SetLanguage(lang)

	code := []byte(`
public protocol Drawable {
    func draw()
    var color: String { get }
}

protocol Codable {
    func encode() -> Data
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewSwiftQuery()
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

	expectedNames := []string{"Drawable", "Codable"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find protocol '%s'", expected)
		}
	}
}

func TestSwiftQueryExtractExtension(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_swift.Language())
	parser.SetLanguage(lang)

	code := []byte(`
extension Point: Drawable {
    func draw() {}
}

extension String {
    func trimmed() -> String {
        return self
    }
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewSwiftQuery()
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

	expectedNames := []string{"Point", "String"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find extension '%s'", expected)
		}
	}
}

func TestSwiftQueryExtractProperties(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_swift.Language())
	parser.SetLanguage(lang)

	code := []byte(`
public let PI = 3.14159
public var count = 0
let maxSize: Int = 1024
var name: String = "test"
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewSwiftQuery()
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

	expectedNames := []string{"PI", "count", "maxSize", "name"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find property '%s'", expected)
		}
	}
}

func TestSwiftQueryExtractInitDeinit(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_swift.Language())
	parser.SetLanguage(lang)

	code := []byte(`
class MyClass {
    init(value: Int) {
        // init
    }

    convenience init() {
        self.init(value: 0)
    }

    deinit {
        // cleanup
    }
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewSwiftQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	foundKinds := make(map[string]bool)
	matchCount := 0

	for {
		match := matches.Next()
		if match == nil {
			break
		}
		matchCount++

		for _, c := range match.Captures {
			capName := captureNames[c.Index]
			if capName == "kind" {
				foundKinds[c.Node.Kind()] = true
			}
		}
	}

	// Should find at least 2 matches: init_declaration(s) and deinit_declaration
	if matchCount < 2 {
		t.Errorf("expected at least 2 matches for init/deinit, got %d", matchCount)
	}

	if !foundKinds["init_declaration"] {
		t.Error("expected to find init_declaration")
	}
	if !foundKinds["deinit_declaration"] {
		t.Error("expected to find deinit_declaration")
	}
}

func TestSwiftQueryExtractSubscript(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_swift.Language())
	parser.SetLanguage(lang)

	code := []byte(`
class Matrix {
    subscript(row: Int, col: Int) -> Double {
        return 0.0
    }
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewSwiftQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	foundSubscript := false

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		for _, c := range match.Captures {
			capName := captureNames[c.Index]
			if capName == "kind" && c.Node.Kind() == "subscript_declaration" {
				foundSubscript = true
			}
		}
	}

	if !foundSubscript {
		t.Error("expected to find subscript_declaration")
	}
}

func TestSwiftQueryExtractImport(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_swift.Language())
	parser.SetLanguage(lang)

	code := []byte(`
import Foundation
import UIKit
import SwiftUI
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewSwiftQuery()
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

	// Should find 3 import declarations
	if count < 3 {
		t.Errorf("expected at least 3 import declarations, got %d", count)
	}
}

func TestSwiftQueryExtractGenerics(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_swift.Language())
	parser.SetLanguage(lang)

	code := []byte(`
public func compare<T: Comparable>(a: T, b: T) -> T {
    return a > b ? a : b
}

public struct Container<Element> {
    var items: [Element] = []
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewSwiftQuery()
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

	expectedNames := []string{"compare", "Container"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestSwiftQueryKindMapping(t *testing.T) {
	query := NewSwiftQuery()
	mapping := query.KindMapping()

	expectedMappings := map[string]string{
		"function_declaration":          "function",
		"class_declaration":             "class",
		"protocol_declaration":          "interface",
		"typealias_declaration":         "type",
		"property_declaration":          "variable",
		"init_declaration":              "constructor",
		"deinit_declaration":            "destructor",
		"subscript_declaration":         "method",
		"operator_declaration":          "function",
		"protocol_function_declaration": "function",
	}

	for nodeType, expected := range expectedMappings {
		if actual, ok := mapping[nodeType]; !ok {
			t.Errorf("expected mapping for '%s'", nodeType)
		} else if actual != expected {
			t.Errorf("expected '%s' -> '%s', got '%s'", nodeType, expected, actual)
		}
	}
}

func TestSwiftQueryCaptures(t *testing.T) {
	query := NewSwiftQuery()
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
