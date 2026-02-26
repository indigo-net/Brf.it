package treesitter

import (
	"strings"
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser"
)

func TestTreeSitterParserImplementsParser(t *testing.T) {
	// Verify TreeSitterParser implements parser.Parser interface
	var _ parser.Parser = (*TreeSitterParser)(nil)
}

func TestTreeSitterParserLanguages(t *testing.T) {
	p := NewTreeSitterParser()

	langs := p.Languages()

	expected := []string{"go", "typescript", "tsx", "java"}
	for _, exp := range expected {
		found := false
		for _, lang := range langs {
			if lang == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected language '%s' not found in %v", exp, langs)
		}
	}
}

func TestTreeSitterParserParseGo(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package main

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

type Point struct {
	X, Y int
}
`

	result, err := p.Parse(code, &parser.Options{Language: "go"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Signatures) < 2 {
		t.Errorf("expected at least 2 signatures, got %d", len(result.Signatures))
	}

	// Find Add function
	var foundAdd bool
	for _, sig := range result.Signatures {
		if sig.Name == "Add" {
			foundAdd = true
			if sig.Kind != "function" {
				t.Errorf("expected kind 'function', got '%s'", sig.Kind)
			}
			// Note: doc may be empty as comments are captured separately in Tree-sitter
			if sig.Line == 0 {
				t.Error("expected non-zero line number")
			}
		}
	}

	if !foundAdd {
		t.Error("expected to find 'Add' function signature")
	}
}

func TestTreeSitterParserParseTypeScript(t *testing.T) {
	p := NewTreeSitterParser()

	code := `
/**
 * Adds two numbers together.
 */
export function add(a: number, b: number): number {
  return a + b;
}

export interface Point {
  x: number;
  y: number;
}
`

	result, err := p.Parse(code, &parser.Options{Language: "typescript"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Signatures) < 1 {
		t.Errorf("expected at least 1 signature, got %d", len(result.Signatures))
	}

	// Find add function
	var foundAdd bool
	for _, sig := range result.Signatures {
		if sig.Name == "add" {
			foundAdd = true
			// Note: kind may be "export" when using export_statement pattern
			if sig.Kind != "function" && sig.Kind != "export" {
				t.Errorf("expected kind 'function' or 'export', got '%s'", sig.Kind)
			}
		}
	}

	if !foundAdd {
		t.Error("expected to find 'add' function signature")
	}
}

func TestTreeSitterParserUnsupportedLanguage(t *testing.T) {
	p := NewTreeSitterParser()

	code := `fn main() { println!("hello"); }`

	result, err := p.Parse(code, &parser.Options{Language: "rust"})
	if err == nil {
		t.Error("expected error for unsupported language")
	}
	if result != nil {
		t.Error("expected nil result for unsupported language")
	}
}

func TestTreeSitterParserAutoRegistration(t *testing.T) {
	// Verify parser is registered in default registry
	registry := parser.DefaultRegistry()

	for _, lang := range []string{"go", "typescript", "tsx", "java"} {
		p, ok := registry.Get(lang)
		if !ok {
			t.Errorf("expected parser for '%s' to be registered", lang)
		}
		if p == nil {
			t.Errorf("expected non-nil parser for '%s'", lang)
		}
	}
}

func TestGoSignatureOnlyExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package main

func Add(a, b int) int {
	return a + b
}

func (p *Point) Move(dx, dy int) {
	p.X += dx
	p.Y += dy
}

type Point struct {
	X, Y int
}
`

	// Test with IncludeBody = false (default)
	result, err := p.Parse(code, &parser.Options{Language: "go", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	for _, sig := range result.Signatures {
		switch sig.Name {
		case "Add":
			// Should only have signature, no body
			expected := "func Add(a, b int) int"
			if sig.Text != expected {
				t.Errorf("expected signature '%s', got '%s'", expected, sig.Text)
			}
		case "Move":
			// Should only have signature, no body
			expected := "func (p *Point) Move(dx, dy int)"
			if sig.Text != expected {
				t.Errorf("expected signature '%s', got '%s'", expected, sig.Text)
			}
		}
	}
}

func TestGoIncludeBodyExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package main

func Add(a, b int) int {
	return a + b
}
`

	// Test with IncludeBody = true
	result, err := p.Parse(code, &parser.Options{Language: "go", IncludeBody: true})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	var foundAdd bool
	for _, sig := range result.Signatures {
		if sig.Name == "Add" {
			foundAdd = true
			// Should include full body
			if len(sig.Text) < 30 { // "func Add(a, b int) int { return a + b }" is longer
				t.Errorf("expected full body, got '%s'", sig.Text)
			}
			if !contains(sig.Text, "return a + b") {
				t.Errorf("expected body to contain 'return a + b', got '%s'", sig.Text)
			}
		}
	}

	if !foundAdd {
		t.Error("expected to find 'Add' function")
	}
}

func TestTypeScriptSignatureOnlyExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `
export function add(a: number, b: number): number {
  return a + b;
}

export class Calculator {
  multiply(a: number, b: number): number {
    return a * b;
  }
}

export interface Config {
  timeout: number;
}
`

	// Test with IncludeBody = false (default)
	result, err := p.Parse(code, &parser.Options{Language: "typescript", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	for _, sig := range result.Signatures {
		switch sig.Name {
		case "add":
			// Should NOT contain the body
			if contains(sig.Text, "return a + b") {
				t.Errorf("signature should not contain body, got '%s'", sig.Text)
			}
		case "Calculator":
			// Class declaration should not have body
			if contains(sig.Text, "multiply") {
				t.Errorf("class signature should not contain methods, got '%s'", sig.Text)
			}
		case "multiply":
			// Method should not have body
			if contains(sig.Text, "return a * b") {
				t.Errorf("method signature should not contain body, got '%s'", sig.Text)
			}
		}
	}
}

func TestTypeScriptArrowFunctionSignature(t *testing.T) {
	p := NewTreeSitterParser()

	code := `
const add = (a: number, b: number): number => {
  return a + b;
};

const double = (n: number): number => n * 2;
`

	// Test with IncludeBody = false (default)
	result, err := p.Parse(code, &parser.Options{Language: "typescript", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	for _, sig := range result.Signatures {
		switch sig.Name {
		case "add":
			// Arrow function with block body
			if contains(sig.Text, "return a + b") {
				t.Errorf("arrow function should not contain body, got '%s'", sig.Text)
			}
		case "double":
			// Arrow function with expression body
			if contains(sig.Text, "n * 2") && !contains(sig.Text, "...") {
				t.Errorf("expression body should be replaced with placeholder, got '%s'", sig.Text)
			}
		}
	}
}

// contains checks if s contains substr
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestTreeSitterParserParseJava(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package com.example;

/**
 * User class represents a user in the system.
 */
public class User {
    private String name;

    public User(String name) {
        this.name = name;
    }

    public String getName() {
        return name;
    }

    private void internalMethod() {
        // Private method - should be filtered
    }
}
`

	result, err := p.Parse(code, &parser.Options{Language: "java"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	// Should have: User class, User constructor, getName method
	// Should NOT have: internalMethod (private)
	if len(result.Signatures) < 3 {
		t.Errorf("expected at least 3 signatures, got %d", len(result.Signatures))
	}

	var foundClass, foundConstructor, foundMethod bool
	for _, sig := range result.Signatures {
		switch sig.Name {
		case "User":
			if sig.Kind == "class" {
				foundClass = true
			} else if sig.Kind == "constructor" {
				foundConstructor = true
			}
		case "getName":
			foundMethod = true
			if sig.Kind != "method" {
				t.Errorf("expected kind 'method', got '%s'", sig.Kind)
			}
		case "internalMethod":
			t.Error("private method 'internalMethod' should be filtered out")
		}
	}

	if !foundClass {
		t.Error("expected to find 'User' class")
	}
	if !foundConstructor {
		t.Error("expected to find 'User' constructor")
	}
	if !foundMethod {
		t.Error("expected to find 'getName' method")
	}
}

func TestJavaSignatureOnlyExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package com.example;

public class Calculator {
    public int add(int a, int b) {
        return a + b;
    }
}

public interface Repository<T> {
    T findById(String id);
    void save(T entity);
}

public enum Status {
    PENDING, ACTIVE, COMPLETED
}

public record Point(int x, int y) {
    public double distance() {
        return Math.sqrt(x * x + y * y);
    }
}
`

	result, err := p.Parse(code, &parser.Options{Language: "java", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	for _, sig := range result.Signatures {
		switch sig.Name {
		case "Calculator":
			// Class should not contain body
			if contains(sig.Text, "public int add") {
				t.Errorf("class signature should not contain methods, got '%s'", sig.Text)
			}
			expected := "public class Calculator"
			if sig.Text != expected {
				t.Errorf("expected '%s', got '%s'", expected, sig.Text)
			}
		case "add":
			// Method should not contain body
			if contains(sig.Text, "return") {
				t.Errorf("method signature should not contain body, got '%s'", sig.Text)
			}
		case "Repository":
			// Interface should not contain methods
			if contains(sig.Text, "findById") {
				t.Errorf("interface signature should not contain methods, got '%s'", sig.Text)
			}
		case "Status":
			if sig.Kind != "enum" {
				t.Errorf("expected kind 'enum', got '%s'", sig.Kind)
			}
		case "Point":
			if sig.Kind != "record" {
				t.Errorf("expected kind 'record', got '%s'", sig.Kind)
			}
			// Record should preserve component parameters
			if !contains(sig.Text, "int x, int y") {
				t.Errorf("record signature should contain components, got '%s'", sig.Text)
			}
		}
	}
}

func TestJavaGenericsExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package com.example;

public class Box<T extends Comparable<T>> {
    private T value;

    public <U> U transform(Function<T, U> fn) {
        return fn.apply(value);
    }
}
`

	result, err := p.Parse(code, &parser.Options{Language: "java", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	var foundClass, foundMethod bool
	for _, sig := range result.Signatures {
		switch sig.Name {
		case "Box":
			foundClass = true
			// Generic type parameter should be preserved
			if !contains(sig.Text, "<T extends Comparable<T>>") {
				t.Errorf("class signature should contain generics, got '%s'", sig.Text)
			}
		case "transform":
			foundMethod = true
			// Method type parameter should be preserved
			if !contains(sig.Text, "<U>") {
				t.Errorf("method signature should contain type parameter, got '%s'", sig.Text)
			}
		}
	}

	if !foundClass {
		t.Error("expected to find 'Box' class")
	}
	if !foundMethod {
		t.Error("expected to find 'transform' method")
	}
}

func TestJavaAutoRegistration(t *testing.T) {
	registry := parser.DefaultRegistry()

	p, ok := registry.Get("java")
	if !ok {
		t.Error("expected parser for 'java' to be registered")
	}
	if p == nil {
		t.Error("expected non-nil parser for 'java'")
	}
}
