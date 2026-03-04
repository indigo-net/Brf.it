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

	expected := []string{"go", "typescript", "tsx", "java", "cpp", "rust", "swift", "kotlin"}
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

	code := `puts "hello"`

	result, err := p.Parse(code, &parser.Options{Language: "ruby"})
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

	for _, lang := range []string{"go", "typescript", "tsx", "java", "cpp", "rust", "swift", "kotlin"} {
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
        // Private method - now included (no private filtering)
    }
}
`

	result, err := p.Parse(code, &parser.Options{Language: "java"})
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

	result, err := p.Parse(code, &parser.Options{Language: "cpp"})
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

	result, err := p.Parse(code, &parser.Options{Language: "cpp", IncludeBody: false})
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

	result, err := p.Parse(code, &parser.Options{Language: "cpp", IncludeBody: false})
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

	result, err := p.Parse(code, &parser.Options{Language: "cpp", IncludeImports: true})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Imports) != 3 {
		t.Errorf("expected 3 imports, got %d", len(result.Imports))
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

	result, err := p.Parse(code, &parser.Options{Language: "rust"})
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

	result, err := p.Parse(code, &parser.Options{Language: "rust", IncludeBody: false})
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

	result, err := p.Parse(code, &parser.Options{Language: "rust", IncludeImports: true})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(result.Imports) < 3 {
		t.Errorf("expected at least 3 imports, got %d", len(result.Imports))
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

	result, err := p.Parse(code, &parser.Options{Language: "rust"})
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

	result, err := p.Parse(code, &parser.Options{Language: "rust"})
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

	result, err := p.Parse(code, &parser.Options{Language: "rust", IncludeBody: false})
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

	result, err := p.Parse(code, &parser.Options{Language: "swift"})
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

	result, err := p.Parse(code, &parser.Options{Language: "kotlin"})
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
		"Role": true, "ADMIN": true, "USER": true, "GUEST": true,
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
		{"public sealed class State", "class"},
	}

	for _, tt := range tests {
		result := refineKotlinClassKind(tt.input)
		if result != tt.expected {
			t.Errorf("refineKotlinClassKind(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
