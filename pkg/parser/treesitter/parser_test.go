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

	expected := []string{"go", "typescript", "tsx", "java", "cpp", "rust", "swift", "kotlin", "csharp", "lua", "php", "ruby", "scala", "elixir"}
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

	result, err := p.Parse([]byte(code), &parser.Options{Language: "go"})
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

	result, err := p.Parse([]byte(code), &parser.Options{Language: "typescript"})
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

	code := `print("hello")`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "haskell"})
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

	for _, lang := range []string{"go", "typescript", "tsx", "java", "cpp", "rust", "swift", "kotlin", "csharp", "lua", "ruby", "scala", "elixir"} {
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
	result, err := p.Parse([]byte(code), &parser.Options{Language: "go", IncludeBody: false})
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
	result, err := p.Parse([]byte(code), &parser.Options{Language: "go", IncludeBody: true})
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
	result, err := p.Parse([]byte(code), &parser.Options{Language: "typescript", IncludeBody: false})
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
	result, err := p.Parse([]byte(code), &parser.Options{Language: "typescript", IncludeBody: false})
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
        // Private method - now included (no private filtering)
    }
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "java"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	// Should have: User class, User constructor, getName method, internalMethod
	// Private methods are now included (no private filtering)
	if len(result.Signatures) < 4 {
		t.Errorf("expected at least 4 signatures, got %d", len(result.Signatures))
	}

	var foundClass, foundConstructor, foundPublicMethod, foundPrivateMethod bool
	for _, sig := range result.Signatures {
		switch sig.Name {
		case "User":
			if sig.Kind == "class" {
				foundClass = true
			} else if sig.Kind == "constructor" {
				foundConstructor = true
			}
		case "getName":
			foundPublicMethod = true
			if sig.Kind != "method" {
				t.Errorf("expected kind 'method', got '%s'", sig.Kind)
			}
		case "internalMethod":
			foundPrivateMethod = true
			if sig.Kind != "method" {
				t.Errorf("expected kind 'method', got '%s'", sig.Kind)
			}
		}
	}

	if !foundClass {
		t.Error("expected to find 'User' class")
	}
	if !foundConstructor {
		t.Error("expected to find 'User' constructor")
	}
	if !foundPublicMethod {
		t.Error("expected to find 'getName' method")
	}
	if !foundPrivateMethod {
		t.Error("expected to find 'internalMethod' (private methods now included)")
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

	result, err := p.Parse([]byte(code), &parser.Options{Language: "java", IncludeBody: false})
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

	result, err := p.Parse([]byte(code), &parser.Options{Language: "java", IncludeBody: false})
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

func TestTreeSitterParserParseCpp(t *testing.T) {
	p := NewTreeSitterParser()

	code := `#include <iostream>

// Simple class without constructor (to avoid name collision)
class SimpleClass {
public:
    void doSomething();
};

namespace utils {
    int helper(int x);
}

template<typename T>
T getMax(T a, T b) {
    return (a > b) ? a : b;
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "cpp"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Signatures) < 3 {
		t.Errorf("expected at least 3 signatures, got %d", len(result.Signatures))
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Check for class
	if kind, ok := foundNames["SimpleClass"]; !ok {
		t.Error("expected to find class 'SimpleClass'")
	} else if kind != "class" {
		t.Errorf("expected kind 'class' for SimpleClass, got '%s'", kind)
	}

	// Check for namespace
	if kind, ok := foundNames["utils"]; !ok {
		t.Error("expected to find namespace 'utils'")
	} else if kind != "namespace" {
		t.Errorf("expected kind 'namespace' for utils, got '%s'", kind)
	}
}

func TestCppSignatureOnlyExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `class Calculator {
public:
    int add(int a, int b) {
        return a + b;
    }
};

int multiply(int a, int b) {
    return a * b;
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "cpp", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	for _, sig := range result.Signatures {
		switch sig.Name {
		case "Calculator":
			// Class should not contain body
			if contains(sig.Text, "int add") {
				t.Errorf("class signature should not contain methods, got '%s'", sig.Text)
			}
		case "multiply":
			// Function should not contain body
			if contains(sig.Text, "return") {
				t.Errorf("function signature should not contain body, got '%s'", sig.Text)
			}
		}
	}
}

func TestCppTemplateExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `template<typename T>
class Box {
    T value;
public:
    T getValue() const;
};

template<typename T, typename U>
T convert(U input) {
    return static_cast<T>(input);
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "cpp", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Check for template class
	if _, ok := foundNames["Box"]; !ok {
		t.Error("expected to find template class 'Box'")
	}

	// Check for template function
	if _, ok := foundNames["convert"]; !ok {
		t.Error("expected to find template function 'convert'")
	}
}

func TestCppAutoRegistration(t *testing.T) {
	registry := parser.DefaultRegistry()

	p, ok := registry.Get("cpp")
	if !ok {
		t.Error("expected parser for 'cpp' to be registered")
	}
	if p == nil {
		t.Error("expected non-nil parser for 'cpp'")
	}
}

func TestCppImportExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `#include <iostream>
#include <vector>
#include "myheader.h"

int main() {
    return 0;
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "cpp", IncludeImports: true})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.RawImports) != 3 {
		t.Errorf("expected 3 imports, got %d", len(result.RawImports))
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

	result, err := p.Parse([]byte(code), &parser.Options{Language: "go"})
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

	result, err := p.Parse([]byte(code), &parser.Options{Language: "typescript"})
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

	result, err := p.Parse([]byte(code), &parser.Options{Language: "python"})
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

	result, err := p.Parse([]byte(code), &parser.Options{Language: "java"})
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

	result, err := p.Parse([]byte(code), &parser.Options{Language: "c"})
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

	result, err := p.Parse([]byte(code), &parser.Options{Language: "go", IncludeBody: false})
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

// === Rust tests ===

func TestTreeSitterParserParseRust(t *testing.T) {
	p := NewTreeSitterParser()

	code := `/// Adds two numbers together.
pub fn add(a: i32, b: i32) -> i32 {
    a + b
}

pub struct Point {
    pub x: f64,
    pub y: f64,
}

impl Point {
    pub fn new(x: f64, y: f64) -> Self {
        Point { x, y }
    }
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "rust"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Signatures) < 3 {
		t.Errorf("expected at least 3 signatures, got %d", len(result.Signatures))
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Check for function
	if kind, ok := foundNames["add"]; !ok {
		t.Error("expected to find function 'add'")
	} else if kind != "function" {
		t.Errorf("expected kind 'function' for 'add', got '%s'", kind)
	}

	// Check for struct
	if kind, ok := foundNames["Point"]; !ok {
		t.Error("expected to find struct 'Point'")
	} else if kind != "struct" && kind != "impl" {
		t.Errorf("expected kind 'struct' or 'impl' for 'Point', got '%s'", kind)
	}

	// Check for method in impl block
	if kind, ok := foundNames["new"]; !ok {
		t.Error("expected to find method 'new'")
	} else if kind != "function" {
		t.Errorf("expected kind 'function' for 'new', got '%s'", kind)
	}
}

func TestRustSignatureOnlyExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `pub fn add(a: i32, b: i32) -> i32 {
    a + b
}

pub struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    pub fn area(&self) -> u32 {
        self.width * self.height
    }
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "rust", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	for _, sig := range result.Signatures {
		switch sig.Name {
		case "add":
			// Function should not contain body
			if contains(sig.Text, "a + b") {
				t.Errorf("function signature should not contain body, got '%s'", sig.Text)
			}
		case "Rectangle":
			// Struct should not contain body
			if contains(sig.Text, "width: u32") {
				t.Errorf("struct signature should not contain body, got '%s'", sig.Text)
			}
		case "area":
			// Method should not contain body
			if contains(sig.Text, "self.width * self.height") {
				t.Errorf("method signature should not contain body, got '%s'", sig.Text)
			}
		}
	}
}

func TestRustImportExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `use std::collections::HashMap;
use std::io::{self, Read, Write};
use crate::utils::*;

pub fn main() {}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "rust", IncludeImports: true})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.RawImports) < 3 {
		t.Errorf("expected at least 3 imports, got %d", len(result.RawImports))
	}
}

func TestRustAutoRegistration(t *testing.T) {
	registry := parser.DefaultRegistry()

	p, ok := registry.Get("rust")
	if !ok {
		t.Error("expected parser for 'rust' to be registered")
	}
	if p == nil {
		t.Error("expected non-nil parser for 'rust'")
	}
}

func TestRustConstAndStaticExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `pub const MAX_SIZE: usize = 1024;
const PRIVATE_CONST: i32 = 42;
pub static GLOBAL_COUNTER: AtomicUsize = AtomicUsize::new(0);
static mut MUTABLE_STATIC: i32 = 0;

pub fn main() {}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "rust"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Constants and statics should be found as "variable"
	expectedVars := []string{"MAX_SIZE", "PRIVATE_CONST", "GLOBAL_COUNTER", "MUTABLE_STATIC"}
	for _, name := range expectedVars {
		if kind, ok := foundNames[name]; !ok {
			t.Errorf("expected to find '%s'", name)
		} else if kind != "variable" {
			t.Errorf("expected kind 'variable' for '%s', got '%s'", name, kind)
		}
	}
}

func TestRustMacroExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `macro_rules! say_hello {
    () => {
        println!("Hello!");
    };
}

pub fn main() {}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "rust"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Macro should be found
	if kind, ok := foundNames["say_hello"]; !ok {
		t.Error("expected to find macro 'say_hello'")
	} else if kind != "macro" {
		t.Errorf("expected kind 'macro' for 'say_hello', got '%s'", kind)
	}
}

func TestRustGenericsAndLifetimes(t *testing.T) {
	p := NewTreeSitterParser()

	code := `pub fn with_lifetime<'a>(s: &'a str) -> &'a str {
    s
}

pub struct Container<T> {
    value: T,
}

pub fn generic_fn<T: Clone>(item: T) -> T {
    item.clone()
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "rust", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	foundNames := make(map[string]bool)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = true
		// Verify generics and lifetimes are preserved in signature
		if sig.Name == "with_lifetime" && !contains(sig.Text, "'a") {
			t.Errorf("expected lifetime to be preserved, got '%s'", sig.Text)
		}
		if sig.Name == "Container" && !contains(sig.Text, "<T>") {
			t.Errorf("expected generic to be preserved, got '%s'", sig.Text)
		}
		if sig.Name == "generic_fn" && !contains(sig.Text, "T: Clone") {
			t.Errorf("expected trait bound to be preserved, got '%s'", sig.Text)
		}
	}

	expectedNames := []string{"with_lifetime", "Container", "generic_fn"}
	for _, name := range expectedNames {
		if !foundNames[name] {
			t.Errorf("expected to find '%s'", name)
		}
	}
}

func TestTreeSitterParserParseSwift(t *testing.T) {
	p := NewTreeSitterParser()

	code := `/// A 2D point structure.
public struct Point {
    public let x: Double
    public let y: Double

    /// Calculate distance from origin.
    public func distance() -> Double {
        return sqrt(x * x + y * y)
    }
}

public protocol Drawable {
    func draw()
}

extension Point: Drawable {
    func draw() {}
}

public func greet(name: String) -> String {
    return "Hello"
}

public let PI = 3.14159
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "swift"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Signatures) < 4 {
		t.Errorf("expected at least 4 signatures, got %d", len(result.Signatures))
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Check for struct
	if kind, ok := foundNames["Point"]; !ok {
		t.Error("expected to find struct 'Point'")
	} else if kind != "struct" && kind != "type" {
		t.Errorf("expected kind 'struct' or 'type' for 'Point', got '%s'", kind)
	}

	// Check for function
	if kind, ok := foundNames["greet"]; !ok {
		t.Error("expected to find function 'greet'")
	} else if kind != "function" {
		t.Errorf("expected kind 'function' for 'greet', got '%s'", kind)
	}

	// Check for protocol
	if kind, ok := foundNames["Drawable"]; !ok {
		t.Error("expected to find protocol 'Drawable'")
	} else if kind != "interface" {
		t.Errorf("expected kind 'interface' for 'Drawable', got '%s'", kind)
	}

	// Check for property
	if kind, ok := foundNames["PI"]; !ok {
		t.Error("expected to find property 'PI'")
	} else if kind != "variable" {
		t.Errorf("expected kind 'variable' for 'PI', got '%s'", kind)
	}

	// Verify body stripping (signature only mode)
	for _, sig := range result.Signatures {
		if sig.Name == "greet" {
			expected := "public func greet(name: String) -> String"
			if sig.Text != expected {
				t.Errorf("expected signature '%s', got '%s'", expected, sig.Text)
			}
		}
		if sig.Name == "distance" {
			expected := "public func distance() -> Double"
			if sig.Text != expected {
				t.Errorf("expected signature '%s', got '%s'", expected, sig.Text)
			}
		}
	}
}

func TestTreeSitterParserParseKotlin(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package com.example

/** User data class. */
data class User(val id: Long, val name: String) {
    fun isValid(): Boolean = name.isNotEmpty()
}

interface UserRepository {
    fun getUser(id: Long): User?
    fun save(user: User): Boolean
}

object AppConfig {
    const val VERSION = "1.0.0"
}

fun <T : Any> requireNotNull(value: T?): T {
    return value ?: throw IllegalArgumentException()
}

typealias UserCallback = (User) -> Unit

enum class Role {
    ADMIN,
    USER,
    GUEST;
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "kotlin"})
	if err != nil {
		t.Fatalf("failed to parse Kotlin: %v", err)
	}
	if len(result.Signatures) == 0 {
		t.Fatal("expected signatures but got none")
	}

	expectedNames := map[string]bool{
		"User": true, "isValid": true,
		"UserRepository": true, "getUser": true, "save": true,
		"AppConfig": true, "VERSION": true,
		"requireNotNull": true,
		"UserCallback":   true,
		"Role":           true, "ADMIN": true, "USER": true, "GUEST": true,
	}

	foundNames := make(map[string]bool)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = true
	}

	for name := range expectedNames {
		if !foundNames[name] {
			t.Errorf("expected to find '%s' in signatures", name)
		}
	}
}

func TestKotlinBodyStripping(t *testing.T) {
	tests := []struct {
		input    string
		kind     string
		expected string
	}{
		{
			input:    "fun greet(name: String): String { return \"Hello\" }",
			kind:     "function",
			expected: "fun greet(name: String): String",
		},
		{
			input:    "fun double(x: Int) = x * 2",
			kind:     "function",
			expected: "fun double(x: Int) = x * 2",
		},
		{
			input:    "data class User(val name: String) { fun greet() = name }",
			kind:     "class",
			expected: "data class User(val name: String)",
		},
		{
			input:    "sealed class Result<out T> { data class Success<T>(val d: T) : Result<T>() }",
			kind:     "class",
			expected: "sealed class Result<out T>",
		},
		{
			input:    "interface Repository<T> { fun getAll(): List<T> }",
			kind:     "interface",
			expected: "interface Repository<T>",
		},
	}

	for _, tt := range tests {
		result := stripKotlinBody(tt.input, tt.kind)
		if result != tt.expected {
			t.Errorf("stripKotlinBody(%q, %q) = %q, want %q", tt.input, tt.kind, result, tt.expected)
		}
	}
}

func TestRefineKotlinClassKind(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"class User", "class"},
		{"data class Point", "class"},
		{"sealed class Result", "class"},
		{"abstract class Base", "class"},
		{"interface Repository", "interface"},
		{"enum class Color", "enum"},
		{"annotation class Api", "class"},
		{"open class Vehicle", "class"},
		{"fun interface Predicate", "interface"},
		{"@Serializable data class User", "class"},
		{"@Retention(AnnotationRetention.RUNTIME) interface Marker", "interface"},
		{"@JvmField(name = \"x\") enum class Color", "enum"},
		{"public sealed class State", "class"},
		// Edge cases: malformed annotations (defensive)
		{"@Foo(bar class X", "class"}, // unmatched parenthesis - stop stripping
		{"@Foo) class Y", "class"},    // depth negative - stop stripping
		{"@ class Z", "class"},        // empty annotation - class keyword recognized
		// Normal annotation processing
		{"@A(@B(val=1)) class T", "class"}, // nested annotation
	}

	for _, tt := range tests {
		result := refineKotlinClassKind(tt.input)
		if result != tt.expected {
			t.Errorf("refineKotlinClassKind(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestKotlinAutoRegistration(t *testing.T) {
	registry := parser.DefaultRegistry()

	p, ok := registry.Get("kotlin")
	if !ok {
		t.Error("expected parser for 'kotlin' to be registered")
	}
	if p == nil {
		t.Error("expected non-nil parser for 'kotlin'")
	}
}

func TestKotlinSignatureOnlyExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `package com.example

class Calculator {
    fun add(a: Int, b: Int): Int {
        return a + b
    }
}

interface Repository<T> {
    fun findById(id: String): T?
    fun save(entity: T)
}

enum class Status {
    PENDING, ACTIVE, COMPLETED
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "kotlin", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	for _, sig := range result.Signatures {
		switch sig.Name {
		case "Calculator":
			if contains(sig.Text, "fun add") {
				t.Errorf("class signature should not contain methods, got '%s'", sig.Text)
			}
		case "add":
			if contains(sig.Text, "return") {
				t.Errorf("method signature should not contain body, got '%s'", sig.Text)
			}
		case "Repository":
			if contains(sig.Text, "fun findById") {
				t.Errorf("interface signature should not contain methods, got '%s'", sig.Text)
			}
		}
	}
}

func TestKotlinImportExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `import kotlin.collections.List
import kotlin.io.println
import kotlinx.coroutines.launch

fun main() {}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "kotlin", IncludeImports: true})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.RawImports) < 3 {
		t.Errorf("expected at least 3 imports, got %d", len(result.RawImports))
	}
}

func TestParsePanicRecovery(t *testing.T) {
	p := NewTreeSitterParser()

	tests := []struct {
		name     string
		content  string
		language string
	}{
		{
			name:     "empty content",
			content:  "",
			language: "go",
		},
		{
			name:     "binary-like content",
			content:  string([]byte{0x00, 0x01, 0x02, 0xff, 0xfe}),
			language: "go",
		},
		{
			name:     "extremely nested braces",
			content:  strings.Repeat("{", 1000) + strings.Repeat("}", 1000),
			language: "go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic regardless of input
			_, _ = p.Parse([]byte(tt.content), &parser.Options{Language: tt.language})
		})
	}
}

func TestTreeSitterParserParseCSharp(t *testing.T) {
	p := NewTreeSitterParser()

	code := `using System;

namespace MyApp {
    /// <summary>Calculator class.</summary>
    public class Calculator {
        public const int MaxValue = 100;
        private int _count;

        public Calculator() { _count = 0; }
        ~Calculator() { }

        public int Add(int a, int b) { return a + b; }
        public string Title { get; set; }
    }

    public struct Point { }

    public interface IDrawable {
        void Draw();
    }

    public enum Color {
        Red,
        Green,
        Blue
    }

    public record Person(string Name, int Age);
    public delegate void Handler(string msg);
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "csharp"})
	if err != nil {
		t.Fatalf("failed to parse C#: %v", err)
	}
	if len(result.Signatures) == 0 {
		t.Fatal("expected signatures but got none")
	}

	expectedNames := map[string]bool{
		"MyApp": true, "Calculator": true, "MaxValue": true,
		"Add": true, "Title": true,
		"Point": true, "IDrawable": true, "Draw": true,
		"Color": true, "Red": true, "Green": true, "Blue": true,
		"Person": true, "Handler": true,
	}

	foundNames := make(map[string]bool)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = true
	}

	for name := range expectedNames {
		if !foundNames[name] {
			t.Errorf("expected to find '%s' in signatures", name)
		}
	}

	// Verify non-static field is excluded
	if foundNames["_count"] {
		t.Error("should not find non-static instance field '_count'")
	}
}

func TestCSharpSignatureOnlyExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `using System;

public class Calculator {
    public int Add(int a, int b) {
        return a + b;
    }

    public async Task<int> AddAsync(int a, int b) => a + b;
}

public interface IService {
    void Execute();
}

public enum Status {
    Active, Inactive
}

public record Person(string Name, int Age);
public record struct Measurement(double Value);
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "csharp", IncludeBody: false})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	for _, sig := range result.Signatures {
		switch sig.Name {
		case "Calculator":
			if contains(sig.Text, "public int Add") {
				t.Errorf("class signature should not contain methods, got '%s'", sig.Text)
			}
		case "Add":
			if contains(sig.Text, "return") {
				t.Errorf("method signature should not contain body, got '%s'", sig.Text)
			}
		case "AddAsync":
			if contains(sig.Text, "a + b") {
				t.Errorf("expression-bodied method should have body stripped, got '%s'", sig.Text)
			}
		case "IService":
			if contains(sig.Text, "void Execute") {
				t.Errorf("interface signature should not contain methods, got '%s'", sig.Text)
			}
		case "Status":
			if sig.Kind != "enum" {
				t.Errorf("expected kind 'enum', got '%s'", sig.Kind)
			}
		case "Person":
			if sig.Kind != "record" {
				t.Errorf("expected kind 'record', got '%s'", sig.Kind)
			}
		case "Measurement":
			if sig.Kind != "struct" {
				t.Errorf("expected kind 'struct' for record struct, got '%s'", sig.Kind)
			}
		}
	}
}

func TestCSharpOperatorNameSynthesis(t *testing.T) {
	p := NewTreeSitterParser()

	code := `public class Vec {
    public static Vec operator +(Vec a, Vec b) => a;
    public static implicit operator int(Vec v) => 0;
    public static explicit operator string(Vec v) => "";
    public int this[int idx] => idx;
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "csharp"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	expectedNames := map[string]bool{
		"operator+":                true,
		"implicit operator int":    true,
		"explicit operator string": true,
		"this":                     true,
	}

	foundNames := make(map[string]bool)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = true
	}

	for name := range expectedNames {
		if !foundNames[name] {
			t.Errorf("expected to find synthesized name '%s', found: %v", name, foundNames)
		}
	}
}

func TestCSharpStaticFieldExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `public class Config {
    public const int MAX = 100;
    public static readonly string NAME = "cfg";
    public static int Counter = 0;
    private int _instanceField;
    private readonly string _name;
}
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "csharp"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	foundNames := make(map[string]string)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = sig.Kind
	}

	// Static/const fields should be found as "variable"
	for _, name := range []string{"MAX", "NAME", "Counter"} {
		if kind, ok := foundNames[name]; !ok {
			t.Errorf("expected to find static/const field '%s'", name)
		} else if kind != "variable" {
			t.Errorf("expected kind 'variable' for '%s', got '%s'", name, kind)
		}
	}

	// Instance fields should NOT be found
	if _, ok := foundNames["_instanceField"]; ok {
		t.Error("should not find non-static instance field '_instanceField'")
	}
	if _, ok := foundNames["_name"]; ok {
		t.Error("should not find non-static instance field '_name'")
	}
}

func TestCSharpImportExtraction(t *testing.T) {
	p := NewTreeSitterParser()

	code := `using System;
using System.Collections.Generic;
using static System.Math;
global using System.Linq;

public class Foo { }
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "csharp", IncludeImports: true})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.RawImports) < 4 {
		t.Errorf("expected at least 4 imports, got %d", len(result.RawImports))
	}
}

func TestCSharpAutoRegistration(t *testing.T) {
	registry := parser.DefaultRegistry()

	p, ok := registry.Get("csharp")
	if !ok {
		t.Error("expected parser for 'csharp' to be registered")
	}
	if p == nil {
		t.Error("expected non-nil parser for 'csharp'")
	}
}

func TestCSharpBodyStripping(t *testing.T) {
	tests := []struct {
		input    string
		kind     string
		expected string
	}{
		{
			input:    "public int Add(int a, int b) { return a + b; }",
			kind:     "method",
			expected: "public int Add(int a, int b)",
		},
		{
			input:    "public int Double(int x) => x * 2;",
			kind:     "method",
			expected: "public int Double(int x)",
		},
		{
			input:    "public Calculator() { _count = 0; }",
			kind:     "constructor",
			expected: "public Calculator()",
		},
		{
			input:    "~Calculator() { }",
			kind:     "destructor",
			expected: "~Calculator()",
		},
		{
			input:    "public class Foo<T> where T : class { public void Bar() { } }",
			kind:     "class",
			expected: "public class Foo<T> where T : class",
		},
		{
			input:    "public struct Point { public int X; }",
			kind:     "struct",
			expected: "public struct Point",
		},
		{
			input:    "public interface IFoo { void Bar(); }",
			kind:     "interface",
			expected: "public interface IFoo",
		},
		{
			input:    "public enum Color { Red, Green }",
			kind:     "enum",
			expected: "public enum Color",
		},
		{
			input:    "public record Person(string Name, int Age);",
			kind:     "record",
			expected: "public record Person(string Name, int Age);",
		},
		{
			input:    "namespace MyApp { public class Foo { } }",
			kind:     "namespace",
			expected: "namespace MyApp",
		},
		{
			input:    "public string Title { get; set; }",
			kind:     "variable",
			expected: "public string Title { get; set; }",
		},
		{
			input:    "public int ReadOnly => 42;",
			kind:     "variable",
			expected: "public int ReadOnly",
		},
		{
			input:    "public delegate void Action<T>(T item);",
			kind:     "type",
			expected: "public delegate void Action<T>(T item);",
		},
		{
			input:    "void Execute();",
			kind:     "method",
			expected: "void Execute();",
		},
	}

	for _, tt := range tests {
		result := stripCSharpBody(tt.input, tt.kind)
		if result != tt.expected {
			t.Errorf("stripCSharpBody(%q, %q) = %q, want %q", tt.input, tt.kind, result, tt.expected)
		}
	}
}

func TestFindCSharpBodyStart(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"public class Foo { }", 17},
		{"class Foo<T> where T : class { }", 29},
		{"List<Dictionary<int, string>> { }", 30},
		{"no brace here", -1},
	}

	for _, tt := range tests {
		result := findCSharpBodyStart(tt.input)
		if result != tt.expected {
			t.Errorf("findCSharpBodyStart(%q) = %d, want %d", tt.input, result, tt.expected)
		}
	}
}

func TestIsExpressionBodied(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"public int X => 42;", true},
		{"public int Add(int a, int b) => a + b;", true},
		{"public int Add(int a, int b) { return a + b; }", false},
		{"public string Name { get; set; }", false},
	}

	for _, tt := range tests {
		result := isExpressionBodied(tt.input)
		if result != tt.expected {
			t.Errorf("isExpressionBodied(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestExtractCSharpOperatorName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"public static Vec operator +(Vec a, Vec b) => a;", "operator+"},
		{"public static bool operator ==(Vec a, Vec b) => true;", "operator=="},
		{"public static Vec operator -(Vec a) => a;", "operator-"},
	}

	for _, tt := range tests {
		result := extractCSharpOperatorName(tt.input)
		if result != tt.expected {
			t.Errorf("extractCSharpOperatorName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestExtractCSharpConversionOperatorName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"public static implicit operator int(Vec v) => 0;", "implicit operator int"},
		{"public static explicit operator string(Vec v) => \"\";", "explicit operator string"},
	}

	for _, tt := range tests {
		result := extractCSharpConversionOperatorName(tt.input)
		if result != tt.expected {
			t.Errorf("extractCSharpConversionOperatorName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFindFunctionBodyStart(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "normal function",
			input:    "func foo(a int) {",
			expected: 16,
		},
		{
			name:     "generic function",
			input:    "func foo<T>(a T) {",
			expected: 17,
		},
		{
			name:     "nested parentheses",
			input:    "func foo(a func(b int)) {",
			expected: 24,
		},
		{
			name:     "unbalanced leading close paren",
			input:    ")(x) {",
			expected: 5,
		},
		{
			name:     "pure unbalanced close paren",
			input:    ") {",
			expected: -1,
		},
		{
			name:     "no body",
			input:    "func foo(a int)",
			expected: -1,
		},
		{
			name:     "empty string",
			input:    "",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findFunctionBodyStart(tt.input)
			if result != tt.expected {
				t.Errorf("findFunctionBodyStart(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParsePanicRecoveryMechanism(t *testing.T) {
	// Verify that if a panic occurs inside Parse, it returns an error
	// rather than crashing. A nil LanguageQuery in the queries map causes
	// a nil pointer dereference panic when query.Language() is called.
	p := &TreeSitterParser{
		queries: map[string]LanguageQuery{
			"go": nil, // nil interface triggers panic on method call
		},
	}

	result, err := p.Parse([]byte("package main"), &parser.Options{Language: "go"})
	if err == nil {
		t.Fatal("expected error from recovered panic, got nil")
	}
	if result != nil {
		t.Fatal("expected nil result from recovered panic")
	}
	if !strings.Contains(err.Error(), "panic recovered") {
		t.Errorf("expected panic recovery error message, got: %v", err)
	}
}

func TestTreeSitterParserParseLua(t *testing.T) {
	p := NewTreeSitterParser()

	code := `local M = {}

--- Greets a person by name.
-- @param name string The person's name
function M.greet(name)
    print("Hello, " .. name)
end

function M:init(config)
    self.config = config
end

local function helper()
    return 42
end

function globalFunc(a, b)
    return a + b
end

local callback = function(err, result)
    if err then return nil end
    return result
end
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "lua"})
	if err != nil {
		t.Fatalf("failed to parse Lua: %v", err)
	}
	if len(result.Signatures) == 0 {
		t.Fatal("expected signatures but got none")
	}

	expectedNames := map[string]bool{
		"M":          true,
		"greet":      true,
		"init":       true,
		"helper":     true,
		"globalFunc": true,
		"callback":   true,
	}

	foundNames := make(map[string]bool)
	for _, sig := range result.Signatures {
		foundNames[sig.Name] = true
	}

	for name := range expectedNames {
		if !foundNames[name] {
			t.Errorf("expected to find '%s' in signatures", name)
		}
	}

	// Verify kind refinement
	kindMap := make(map[string]string)
	for _, sig := range result.Signatures {
		kindMap[sig.Name] = sig.Kind
	}

	expectedKinds := map[string]string{
		"greet":      "module_function",
		"init":       "method",
		"helper":     "local_function",
		"globalFunc": "function",
		"M":          "variable",
	}

	for name, expected := range expectedKinds {
		actual, ok := kindMap[name]
		if !ok {
			t.Errorf("expected '%s' to be present in signatures", name)
			continue
		}
		if actual != expected {
			t.Errorf("expected '%s' kind = '%s', got '%s'", name, expected, actual)
		}
	}

	// Verify body stripping
	for _, sig := range result.Signatures {
		if sig.Kind == "module_function" || sig.Kind == "method" || sig.Kind == "function" || sig.Kind == "local_function" {
			if strings.Contains(sig.Text, "\n") {
				t.Errorf("expected body stripped for '%s', got multiline: %q", sig.Name, sig.Text)
			}
		}
	}

	// Note: doc comments are captured separately in Tree-sitter and may not
	// be associated with the following declaration (same as other languages)
}

func TestTreeSitterParserParseLuaImports(t *testing.T) {
	p := NewTreeSitterParser()

	code := `local json = require("json")
local utils = require("app.utils")
local lfs = require("lfs")

local M = {}
function M.run() end
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "lua", IncludeImports: true})
	if err != nil {
		t.Fatalf("failed to parse Lua imports: %v", err)
	}
	if len(result.RawImports) < 3 {
		t.Errorf("expected at least 3 imports, got %d", len(result.RawImports))
	}

	// Verify raw import text contains require()
	for _, imp := range result.RawImports {
		if !strings.Contains(imp, "require(") {
			t.Errorf("expected raw import text to contain 'require(', got: %q", imp)
		}
	}
}

func TestLuaBodyStripping(t *testing.T) {
	tests := []struct {
		input    string
		kind     string
		expected string
	}{
		{
			input:    "function greet(name)\n    print(name)\nend",
			kind:     "function",
			expected: "function greet(name)",
		},
		{
			input:    "function M.create(name)\n    return {}\nend",
			kind:     "module_function",
			expected: "function M.create(name)",
		},
		{
			input:    "function M:init(config)\n    self.config = config\nend",
			kind:     "method",
			expected: "function M:init(config)",
		},
		{
			input:    "local function helper()\n    return 42\nend",
			kind:     "local_function",
			expected: "local function helper()",
		},
		{
			input:    "local M = {}",
			kind:     "variable",
			expected: "local M = {}",
		},
	}

	for _, tt := range tests {
		result := stripLuaBody(tt.input, tt.kind)
		if result != tt.expected {
			t.Errorf("stripLuaBody(%q, %q) = %q, want %q", tt.input, tt.kind, result, tt.expected)
		}
	}
}

func TestRefineLuaFunctionKind(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"function globalFunc(a, b)\n  return a + b\nend", "function"},
		{"local function helper()\n  return 42\nend", "local_function"},
		{"function M.create(name)\n  return {}\nend", "module_function"},
		{"function M:init(config)\n  self.config = config\nend", "method"},
	}

	for _, tt := range tests {
		result := refineLuaFunctionKind(tt.input)
		if result != tt.expected {
			t.Errorf("refineLuaFunctionKind(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestIsPythonMethod(t *testing.T) {
	tests := []struct {
		name      string
		signature string
		expected  bool
	}{
		{"simple method with self", "def greet(self) -> str:", true},
		{"method with self and params", "def greet(self, name: str) -> str:", true},
		{"classmethod with cls", "def create(cls, data: dict) -> 'MyClass':", true},
		{"plain function", "def greet(name: str) -> str:", false},
		{"no params", "def greet():", false},
		{"no parens", "def greet:", false},
		{"nested parens in type hint", "def process(self, data: Dict[str, Tuple[int, ...]]) -> None:", true},
		{"default value with function call", "def setup(self, config: Config = Config()) -> None:", true},
		{"nested parens not method", "def process(data: Dict[str, Tuple[int, ...]]) -> None:", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPythonMethod(tt.signature)
			if result != tt.expected {
				t.Errorf("isPythonMethod(%q) = %v, want %v", tt.signature, result, tt.expected)
			}
		})
	}
}

func TestTreeSitterParserParsePHP(t *testing.T) {
	p := NewTreeSitterParser()

	code := `<?php

namespace App\Services;

use App\Models\User;

/**
 * UserService handles user-related operations.
 */
class UserService {
    private $repository;

    public function __construct($repo) {
        $this->repository = $repo;
    }

    public function findUser($id) {
        return $this->repository->find($id);
    }
}

interface RepositoryInterface {
    public function find($id);
    public function save($entity);
}

trait Loggable {
    public function log($message) {
        echo $message;
    }
}

function helper($data) {
    return $data;
}

const MAX_ITEMS = 100;
const APP_NAME = "Brf.it";
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "php"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	// Expected signatures: UserService class, __construct, findUser, RepositoryInterface,
	// find, save, Loggable trait, log, helper function, MAX_ITEMS, APP_NAME
	if len(result.Signatures) < 8 {
		t.Errorf("expected at least 8 signatures, got %d", len(result.Signatures))
	}

	// Find UserService class
	var foundClass, foundMethod, foundInterface, foundTrait, foundFunction, foundConst bool
	for _, sig := range result.Signatures {
		switch sig.Name {
		case "UserService":
			foundClass = true
			if sig.Kind != "class" {
				t.Errorf("expected UserService kind 'class', got '%s'", sig.Kind)
			}
		case "findUser":
			foundMethod = true
			if sig.Kind != "method" {
				t.Errorf("expected findUser kind 'method', got '%s'", sig.Kind)
			}
		case "RepositoryInterface":
			foundInterface = true
			if sig.Kind != "interface" {
				t.Errorf("expected RepositoryInterface kind 'interface', got '%s'", sig.Kind)
			}
		case "Loggable":
			foundTrait = true
			if sig.Kind != "type" {
				t.Errorf("expected Loggable kind 'type', got '%s'", sig.Kind)
			}
		case "helper":
			foundFunction = true
			if sig.Kind != "function" {
				t.Errorf("expected helper kind 'function', got '%s'", sig.Kind)
			}
		case "MAX_ITEMS":
			foundConst = true
			if sig.Kind != "variable" {
				t.Errorf("expected MAX_ITEMS kind 'variable', got '%s'", sig.Kind)
			}
		}
	}

	if !foundClass {
		t.Error("expected to find 'UserService' class")
	}
	if !foundMethod {
		t.Error("expected to find 'findUser' method")
	}
	if !foundInterface {
		t.Error("expected to find 'RepositoryInterface' interface")
	}
	if !foundTrait {
		t.Error("expected to find 'Loggable' trait")
	}
	if !foundFunction {
		t.Error("expected to find 'helper' function")
	}
	if !foundConst {
		t.Error("expected to find 'MAX_ITEMS' const")
	}
}

func TestTreeSitterParserParsePHPImports(t *testing.T) {
	p := NewTreeSitterParser()

	code := `<?php
use App\Services\UserService;
use App\Models\Product;
require 'vendor/autoload.php';
include 'config.php';
`

	result, err := p.Parse([]byte(code), &parser.Options{Language: "php", IncludeImports: true})
	if err != nil {
		t.Fatalf("failed to parse PHP imports: %v", err)
	}
	if len(result.RawImports) < 3 {
		t.Errorf("expected at least 3 imports, got %d", len(result.RawImports))
	}

	// Verify raw import text contains use/require/include
	for _, imp := range result.RawImports {
		if !strings.Contains(imp, "use ") && !strings.Contains(imp, "require") && !strings.Contains(imp, "include") {
			t.Errorf("expected raw import text to contain use/require/include, got: %q", imp)
		}
	}
}

func TestPHPBodyStripping(t *testing.T) {
	tests := []struct {
		input    string
		kind     string
		expected string
	}{
		{
			input:    "function greet($name) {\n    return 'Hello, ' . $name;\n}",
			kind:     "function",
			expected: "function greet($name)",
		},
		{
			input:    "public function findUser($id) {\n    return $this->repo->find($id);\n}",
			kind:     "method",
			expected: "public function findUser($id)",
		},
		{
			input:    "class UserService {\n    private $repo;\n}",
			kind:     "class",
			expected: "class UserService",
		},
		{
			input:    "interface Repository {\n    public function find($id);\n}",
			kind:     "interface",
			expected: "interface Repository",
		},
		{
			input:    "enum Status {\n    case Active;\n    case Inactive;\n}",
			kind:     "enum",
			expected: "enum Status",
		},
		{
			input:    "trait Loggable {\n    public function log($msg) {}\n}",
			kind:     "type",
			expected: "trait Loggable",
		},
		{
			input:    "const MAX_SIZE = 100;",
			kind:     "variable",
			expected: "const MAX_SIZE = 100;",
		},
	}

	for _, tt := range tests {
		result := stripPHPBody(tt.input, tt.kind)
		if result != tt.expected {
			t.Errorf("stripPHPBody(%q, %q) = %q, want %q", tt.input, tt.kind, result, tt.expected)
		}
	}
}

func TestTreeSitterParserParseRuby(t *testing.T) {
	p := NewTreeSitterParser()

	code := []byte(`# User class for managing accounts
class User
  attr_reader :name, :email

  def initialize(name, email)
    @name = name
    @email = email
  end

  def self.create(attrs)
    new(attrs[:name], attrs[:email])
  end

  def display
    "#{name} <#{email}>"
  end
end

module Helpers
  def self.format(str)
    str.strip
  end
end

MAX_RETRIES = 3

def greet(name)
  puts "Hello, #{name}!"
end
`)

	opts := &parser.Options{
		Language: "ruby",
	}

	result, err := p.Parse(code, opts)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Language != "ruby" {
		t.Errorf("expected language 'ruby', got '%s'", result.Language)
	}

	expectedNames := map[string]bool{
		"User": false, "initialize": false, "create": false,
		"display": false, "Helpers": false, "format": false,
		"MAX_RETRIES": false, "greet": false,
	}

	for _, sig := range result.Signatures {
		if _, ok := expectedNames[sig.Name]; ok {
			expectedNames[sig.Name] = true
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("expected to find signature '%s'", name)
		}
	}
}

func TestTreeSitterParserParseRubyImports(t *testing.T) {
	p := NewTreeSitterParser()

	code := []byte(`require "json"
require "net/http"
require_relative "helper"

class App
  def run
    data = JSON.parse("{}")
  end
end
`)

	opts := &parser.Options{
		Language:       "ruby",
		IncludeImports: true,
	}

	result, err := p.Parse(code, opts)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.RawImports) < 3 {
		t.Errorf("expected at least 3 imports, got %d: %v", len(result.RawImports), result.RawImports)
	}
}

func TestTreeSitterParserParseScala(t *testing.T) {
	p := NewTreeSitterParser()

	code := []byte(`// User management
trait Greeter {
  def greet(name: String): String
}

sealed trait Shape

abstract class Vehicle(val wheels: Int) {
  def description: String
}

class Person(val name: String, var age: Int) extends Greeter {
  def greet(name: String): String = s"Hello, $name"
  private def helper(): Unit = ()
}

case class Point(x: Double, y: Double) {
  def distance: Double = math.sqrt(x * x + y * y)
}

object MathUtils {
  val PI: Double = 3.14159
  var counter: Int = 0
  def add(a: Int, b: Int): Int = a + b
}

type StringList = List[String]

enum Color {
  case Red, Green, Blue
}

lazy val expensive: String = "computed"
`)

	opts := &parser.Options{
		Language: "scala",
	}

	result, err := p.Parse(code, opts)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.Language != "scala" {
		t.Errorf("expected language 'scala', got '%s'", result.Language)
	}

	expectedNames := map[string]bool{
		"Greeter": false, "greet": false, "Shape": false,
		"Vehicle": false, "description": false,
		"Person": false, "helper": false,
		"Point": false, "distance": false,
		"MathUtils": false, "PI": false, "counter": false, "add": false,
		"StringList": false, "Color": false, "expensive": false,
	}

	for _, sig := range result.Signatures {
		if _, ok := expectedNames[sig.Name]; ok {
			expectedNames[sig.Name] = true
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("expected to find signature '%s'", name)
		}
	}

	// Verify stripBody works for methods
	for _, sig := range result.Signatures {
		if sig.Name == "add" {
			if strings.Contains(sig.Text, "a + b") {
				t.Errorf("expected body to be stripped from 'add', got: %s", sig.Text)
			}
			if !strings.Contains(sig.Text, "def add(a: Int, b: Int): Int") {
				t.Errorf("expected signature to contain 'def add(a: Int, b: Int): Int', got: %s", sig.Text)
			}
		}
		if sig.Name == "Person" {
			if strings.Contains(sig.Text, "def greet") {
				t.Errorf("expected class body to be stripped from 'Person', got: %s", sig.Text)
			}
		}
	}
}

func TestTreeSitterParserParseScalaImports(t *testing.T) {
	p := NewTreeSitterParser()

	code := []byte(`import scala.collection.mutable
import scala.collection.mutable.{ListBuffer, ArrayBuffer}
import java.util._

class App {
  def run(): Unit = ()
}
`)

	opts := &parser.Options{
		Language:       "scala",
		IncludeImports: true,
	}

	result, err := p.Parse(code, opts)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.RawImports) < 3 {
		t.Errorf("expected at least 3 imports, got %d: %v", len(result.RawImports), result.RawImports)
	}
}

func TestTreeSitterParserParseElixir(t *testing.T) {
	p := NewTreeSitterParser()

	code := []byte(`defmodule Calculator do
  @moduledoc "A simple calculator module"

  @type number :: integer() | float()
  @spec add(number(), number()) :: number()

  def add(a, b) do
    a + b
  end

  defp validate(x) when is_number(x) do
    :ok
  end

  defstruct [:value, :operation]
end
`)

	result, err := p.Parse([]byte(code), &parser.Options{Language: "elixir"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Signatures) < 3 {
		t.Errorf("expected at least 3 signatures, got %d", len(result.Signatures))
		for _, sig := range result.Signatures {
			t.Logf("  sig: name=%s kind=%s text=%q", sig.Name, sig.Kind, sig.Text)
		}
	}

	// Verify specific signatures
	foundModule := false
	foundFunc := false
	foundType := false
	for _, sig := range result.Signatures {
		if sig.Name == "Calculator" && sig.Kind == "class" {
			foundModule = true
		}
		if sig.Name == "add" && sig.Kind == "function" {
			foundFunc = true
			// Body should be stripped (no "do...end")
			if strings.Contains(sig.Text, "\n") {
				t.Errorf("expected stripped body for add, got: %q", sig.Text)
			}
		}
		if sig.Name == "number" && sig.Kind == "type" {
			foundType = true
		}
	}

	if !foundModule {
		t.Error("expected to find module 'Calculator'")
	}
	if !foundFunc {
		t.Error("expected to find function 'add'")
	}
	if !foundType {
		t.Error("expected to find type 'number'")
	}
}

func TestTreeSitterParserParseElixirImports(t *testing.T) {
	p := NewTreeSitterParser()

	code := []byte(`defmodule MyApp do
  import Enum
  import String, only: [trim: 1]
  alias MyApp.Accounts.User
  use GenServer
  require Logger

  def start_link(opts) do
    GenServer.start_link(__MODULE__, opts)
  end
end
`)

	opts := &parser.Options{
		Language:       "elixir",
		IncludeImports: true,
	}

	result, err := p.Parse(code, opts)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(result.RawImports) < 2 {
		t.Errorf("expected at least 2 imports, got %d: %v", len(result.RawImports), result.RawImports)
	}

	// Verify defmodule is NOT included in imports
	for _, imp := range result.RawImports {
		if strings.HasPrefix(imp, "defmodule") {
			t.Errorf("defmodule should not be in imports: %q", imp)
		}
	}

	// Verify actual imports are present
	foundImport := false
	for _, imp := range result.RawImports {
		if strings.HasPrefix(imp, "import ") || strings.HasPrefix(imp, "alias ") ||
			strings.HasPrefix(imp, "use ") || strings.HasPrefix(imp, "require ") {
			foundImport = true
			break
		}
	}
	if !foundImport {
		t.Errorf("expected at least one import/alias/use/require statement, got: %v", result.RawImports)
	}
}

func TestRefineElixirCallKind(t *testing.T) {
	tests := []struct {
		text string
		want string
	}{
		{"defmodule MyApp do", "class"},
		{"defprotocol Printable do", "interface"},
		{"defimpl Printable, for: Integer do", "impl"},
		{"def hello(name) do", "function"},
		{"defp validate(x) when is_number(x) do", "function"},
		{"defmacro unless(condition, do: block) do", "macro"},
		{"defmacrop private_macro(x) do", "macro"},
		{"defguard is_positive(x) when is_integer(x) and x > 0", "function"},
		{"defguardp is_even(x) when rem(x, 2) == 0", "function"},
		{"defdelegate keys(map), to: Map", "function"},
		{"defstruct [:name, :email]", "struct"},
		{"if condition do", ""},
		{"case value do", ""},
		{"Enum.map(list, fn x -> x end)", ""},
		{"IO.puts(\"hello\")", ""},
	}

	for _, tt := range tests {
		got := refineElixirCallKind(tt.text)
		if got != tt.want {
			t.Errorf("refineElixirCallKind(%q) = %q, want %q", tt.text, got, tt.want)
		}
	}
}

func TestRefineElixirAttrKind(t *testing.T) {
	tests := []struct {
		text         string
		capturedName string
		wantKind     string
		wantName     string
	}{
		{"@spec add(integer(), integer()) :: integer()", "spec", "type", "add"},
		{"@spec foo :: bar", "spec", "type", "foo"},
		{"@type color :: :red | :green | :blue", "type", "type", "color"},
		{"@typep internal_state :: map()", "typep", "type", "internal_state"},
		{"@opaque hidden :: %__MODULE__{}", "opaque", "type", "hidden"},
		{"@callback handle_event(term()) :: {:ok, term()}", "callback", "type", "handle_event"},
		{"@doc \"Some documentation\"", "doc", "", ""},
		{"@moduledoc \"Module docs\"", "moduledoc", "", ""},
		{"@behaviour GenServer", "behaviour", "", ""},
	}

	for _, tt := range tests {
		gotKind, gotName := refineElixirAttrKind(tt.text, tt.capturedName)
		if gotKind != tt.wantKind || gotName != tt.wantName {
			t.Errorf("refineElixirAttrKind(%q, %q) = (%q, %q), want (%q, %q)",
				tt.text, tt.capturedName, gotKind, gotName, tt.wantKind, tt.wantName)
		}
	}
}

func TestStripElixirBody(t *testing.T) {
	tests := []struct {
		text string
		kind string
		want string
	}{
		{"def add(a, b) do\n  a + b\nend", "function", "def add(a, b)"},
		{"defp validate(x) do\n  :ok\nend", "function", "defp validate(x)"},
		{"defmodule MyApp do\n  use GenServer\nend", "class", "defmodule MyApp"},
		{"def add(a, b), do: a + b", "function", "def add(a, b)"},
		{"@spec add(integer(), integer()) :: integer()", "type", "@spec add(integer(), integer()) :: integer()"},
		{"@type color :: :red | :green | :blue", "type", "@type color :: :red | :green | :blue"},
		{"defstruct [:name, :email]", "struct", "defstruct [:name, :email]"},
		{"def zero_arity do\n  :ok\nend", "function", "def zero_arity"},
	}

	for _, tt := range tests {
		got := stripElixirBody(tt.text, tt.kind)
		if got != tt.want {
			t.Errorf("stripElixirBody(%q, %q) = %q, want %q", tt.text, tt.kind, got, tt.want)
		}
	}
}
