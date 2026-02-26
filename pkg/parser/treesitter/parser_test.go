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

// === Module-level variable extraction tests ===

func TestGoVariableExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package main

const MaxSize = 100

var DefaultConfig = Config{}

const (
	MinValue = 0
	MaxValue = 1000
)

func Add(a, b int) int {
	return a + b
}
`

	result, err := p.Parse(code, &parser.Options{Language: "go"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Variables should be found
	expectedVars := []string{"MaxSize", "DefaultConfig", "MinValue", "MaxValue"}
	for _, name := range expectedVars {
		if kind, ok := foundNames[name]; !ok {
			t.Errorf("expected to find variable '%s'", name)
		} else if kind != "variable" {
			t.Errorf("expected kind 'variable' for '%s', got '%s'", name, kind)
		}
	}

	// Function should still be found
	if kind, ok := foundNames["Add"]; !ok {
		t.Error("expected to find function 'Add'")
	} else if kind != "function" {
		t.Errorf("expected kind 'function' for 'Add', got '%s'", kind)
	}
}

func TestTypeScriptVariableExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `const API_URL = "https://api.example.com";
export const MAX_RETRIES = 3;
let counter = 0;
const arrowFn = (x: number) => x * 2;

export function add(a: number, b: number): number {
  const localVar = 1;
  return a + b;
}
`

	result, err := p.Parse(code, &parser.Options{Language: "typescript"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Module-level variables should be found
	moduleVars := []string{"API_URL", "MAX_RETRIES", "counter"}
	for _, name := range moduleVars {
		if _, ok := foundNames[name]; !ok {
			t.Errorf("expected to find module-level variable '%s'", name)
		}
	}

	// Arrow function should be found (either as arrow or variable)
	if _, ok := foundNames["arrowFn"]; !ok {
		t.Error("expected to find arrow function 'arrowFn'")
	}

	// Function should still be found
	if _, ok := foundNames["add"]; !ok {
		t.Error("expected to find function 'add'")
	}

	// Local variable should NOT be found
	if _, ok := foundNames["localVar"]; ok {
		t.Error("should not find local variable 'localVar'")
	}
}

func TestPythonVariableExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `API_URL = "https://api.example.com"
MAX_RETRIES: int = 3

def test():
    local_var = 1
    return local_var

class Config:
    CLASS_VAR = "value"
`

	result, err := p.Parse(code, &parser.Options{Language: "python"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Module-level variables should be found
	moduleVars := []string{"API_URL", "MAX_RETRIES"}
	for _, name := range moduleVars {
		if kind, ok := foundNames[name]; !ok {
			t.Errorf("expected to find module-level variable '%s'", name)
		} else if kind != "variable" {
			t.Errorf("expected kind 'variable' for '%s', got '%s'", name, kind)
		}
	}

	// Function and class should still be found
	if _, ok := foundNames["test"]; !ok {
		t.Error("expected to find function 'test'")
	}
	if _, ok := foundNames["Config"]; !ok {
		t.Error("expected to find class 'Config'")
	}

	// Local and class variables should NOT be found
	if _, ok := foundNames["local_var"]; ok {
		t.Error("should not find local variable 'local_var'")
	}
	if _, ok := foundNames["CLASS_VAR"]; ok {
		t.Error("should not find class variable 'CLASS_VAR'")
	}
}

func TestJavaStaticFieldExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package com.example;

public class Config {
    public static final String API_URL = "https://api.example.com";
    public static int MAX_RETRIES = 3;
    private String instanceField = "value";

    public void method() {
        int localVar = 1;
    }
}
`

	result, err := p.Parse(code, &parser.Options{Language: "java"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Static fields should be found as "variable"
	staticFields := []string{"API_URL", "MAX_RETRIES"}
	for _, name := range staticFields {
		if kind, ok := foundNames[name]; !ok {
			t.Errorf("expected to find static field '%s'", name)
		} else if kind != "variable" {
			t.Errorf("expected kind 'variable' for '%s', got '%s'", name, kind)
		}
	}

	// Class and method should still be found
	if _, ok := foundNames["Config"]; !ok {
		t.Error("expected to find class 'Config'")
	}
	if _, ok := foundNames["method"]; !ok {
		t.Error("expected to find method 'method'")
	}

	// Non-static instance field should NOT be found (private is filtered, non-static is filtered)
	if _, ok := foundNames["instanceField"]; ok {
		t.Error("should not find non-static instance field 'instanceField'")
	}
}

func TestCGlobalVariableExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `int global_count = 0;
static char* buffer;
extern int shared_value;
const int MAX_SIZE = 100;

int add(int a, int b);

void test() {
    int local_var = 1;
}
`

	result, err := p.Parse(code, &parser.Options{Language: "c"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Global variables should be found as "variable"
	globalVars := []string{"global_count", "buffer", "shared_value", "MAX_SIZE"}
	for _, name := range globalVars {
		if kind, ok := foundNames[name]; !ok {
			t.Errorf("expected to find global variable '%s'", name)
		} else if kind != "variable" {
			t.Errorf("expected kind 'variable' for '%s', got '%s'", name, kind)
		}
	}

	// Function prototype and definition should still be found as "function"
	if kind, ok := foundNames["add"]; !ok {
		t.Error("expected to find function prototype 'add'")
	} else if kind != "function" {
		t.Errorf("expected kind 'function' for 'add', got '%s'", kind)
	}

	if kind, ok := foundNames["test"]; !ok {
		t.Error("expected to find function 'test'")
	} else if kind != "function" {
		t.Errorf("expected kind 'function' for 'test', got '%s'", kind)
	}

	// Local variable should NOT be found
	if _, ok := foundNames["local_var"]; ok {
		t.Error("should not find local variable 'local_var'")
	}
}

func TestVariableSignaturePreservesValue(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package main

const MaxSize = 100
var DefaultTimeout = 30 * time.Second
`

	result, err := p.Parse(code, &parser.Options{Language: "go", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	for _, sig := range result.Signatures {
		switch sig.Name {
		case "MaxSize":
			// Variable signature should preserve the value
			if !contains(sig.Text, "100") {
				t.Errorf("expected variable signature to contain value '100', got '%s'", sig.Text)
			}
		case "DefaultTimeout":
			// Variable signature should preserve the value expression
			if !contains(sig.Text, "30 * time.Second") {
				t.Errorf("expected variable signature to contain value, got '%s'", sig.Text)
			}
		}
	}
}
