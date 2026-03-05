package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"
)

func TestRustQueryLanguage(t *testing.T) {
	query := NewRustQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestRustQueryPattern(t *testing.T) {
	query := NewRustQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_rust.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
	defer q.Close()
}

func TestRustQueryImportPattern(t *testing.T) {
	query := NewRustQuery()
	pattern := query.ImportQuery()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty import query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_rust.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid import query pattern: %v", err)
	}
	defer q.Close()
}

func TestRustQueryExtractFunction(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_rust.Language())
	parser.SetLanguage(lang)

	code := []byte(`
/// Adds two numbers together.
pub fn add(a: i32, b: i32) -> i32 {
    a + b
}

fn private_fn() -> bool {
    true
}

pub async fn async_fetch() -> Result<String, Error> {
    Ok("data".to_string())
}

pub const fn const_double(x: i32) -> i32 {
    x * 2
}

pub unsafe fn dangerous() {
    // unsafe code
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewRustQuery()
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

	expectedNames := []string{"add", "private_fn", "async_fetch", "const_double", "dangerous"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find function '%s'", expected)
		}
	}
}

func TestRustQueryExtractTypes(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_rust.Language())
	parser.SetLanguage(lang)

	code := []byte(`
/// A simple point structure.
pub struct Point {
    pub x: f64,
    pub y: f64,
}

enum Color {
    Red,
    Green,
    Blue,
}

pub trait Drawable {
    fn draw(&self);
}

pub type Coordinate = (f64, f64);

pub union IntOrFloat {
    i: i32,
    f: f32,
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewRustQuery()
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

	expectedNames := []string{"Point", "Color", "Drawable", "Coordinate", "IntOrFloat"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find type '%s'", expected)
		}
	}
}

func TestRustQueryExtractImplAndMethods(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_rust.Language())
	parser.SetLanguage(lang)

	code := []byte(`
pub struct Rectangle {
    width: u32,
    height: u32,
}

impl Rectangle {
    pub fn new(width: u32, height: u32) -> Self {
        Rectangle { width, height }
    }

    pub fn area(&self) -> u32 {
        self.width * self.height
    }

    fn private_method(&self) {}
}

impl Default for Rectangle {
    fn default() -> Self {
        Rectangle { width: 0, height: 0 }
    }
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewRustQuery()
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

	// Should find struct, impl block type, and methods
	expectedNames := []string{"Rectangle", "new", "area", "private_method", "default"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestRustQueryExtractConstAndStatic(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_rust.Language())
	parser.SetLanguage(lang)

	code := []byte(`
pub const MAX_SIZE: usize = 1024;

const PRIVATE_CONST: i32 = 42;

pub static GLOBAL_COUNTER: AtomicUsize = AtomicUsize::new(0);

static mut MUTABLE_STATIC: i32 = 0;
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewRustQuery()
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

	expectedNames := []string{"MAX_SIZE", "PRIVATE_CONST", "GLOBAL_COUNTER", "MUTABLE_STATIC"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestRustQueryExtractMacro(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_rust.Language())
	parser.SetLanguage(lang)

	code := []byte(`
macro_rules! say_hello {
    () => {
        println!("Hello!");
    };
    ($name:expr) => {
        println!("Hello, {}!", $name);
    };
}

macro_rules! vec_of_strings {
    ($($x:expr),*) => {
        vec![$($x.to_string()),*]
    };
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewRustQuery()
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

	expectedNames := []string{"say_hello", "vec_of_strings"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find macro '%s'", expected)
		}
	}
}

func TestRustQueryExtractModule(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_rust.Language())
	parser.SetLanguage(lang)

	code := []byte(`
pub mod utils {
    pub fn helper() {}
}

mod private_mod {
    fn internal() {}
}

mod external_mod;
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewRustQuery()
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

	expectedNames := []string{"utils", "private_mod", "external_mod", "helper", "internal"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestRustQueryExtractUse(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_rust.Language())
	parser.SetLanguage(lang)

	code := []byte(`
use std::collections::HashMap;
use std::io::{self, Read, Write};
use crate::utils::*;
extern crate serde;
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewRustQuery()
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

	// Should find 3 use declarations + 1 extern crate = 4
	if count < 4 {
		t.Errorf("expected at least 4 import declarations, got %d", count)
	}
}

func TestRustQueryExtractGenericsAndLifetimes(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_rust.Language())
	parser.SetLanguage(lang)

	code := []byte(`
pub struct Container<T> {
    value: T,
}

pub fn generic_fn<T: Clone>(item: T) -> T {
    item.clone()
}

pub fn with_lifetime<'a>(s: &'a str) -> &'a str {
    s
}

pub struct WithLifetime<'a, T> {
    reference: &'a T,
}

impl<'a, T: Clone> WithLifetime<'a, T> {
    pub fn get(&self) -> T {
        self.reference.clone()
    }
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewRustQuery()
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

	expectedNames := []string{"Container", "generic_fn", "with_lifetime", "WithLifetime", "get"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}

func TestRustQueryEmptyFile(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_rust.Language())
	parser.SetLanguage(lang)

	code := []byte(``)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewRustQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
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

	if count != 0 {
		t.Errorf("expected 0 matches for empty file, got %d", count)
	}
}

func TestRustQueryKindMapping(t *testing.T) {
	query := NewRustQuery()
	mapping := query.KindMapping()

	expectedMappings := map[string]string{
		"function_item":    "function",
		"struct_item":      "struct",
		"enum_item":        "enum",
		"trait_item":       "trait",
		"type_item":        "type",
		"impl_item":        "impl",
		"const_item":       "variable",
		"static_item":      "variable",
		"mod_item":         "namespace",
		"macro_definition": "macro",
	}

	for nodeType, expected := range expectedMappings {
		if actual, ok := mapping[nodeType]; !ok {
			t.Errorf("expected mapping for '%s'", nodeType)
		} else if actual != expected {
			t.Errorf("expected '%s' -> '%s', got '%s'", nodeType, expected, actual)
		}
	}
}

func TestRustQueryCaptures(t *testing.T) {
	query := NewRustQuery()
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
