package treesitter

import (
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser"
)

// FuzzParseGo tests that the Go parser does not panic on arbitrary input.
func FuzzParseGo(f *testing.F) {
	// Add seed corpus with valid and invalid Go code
	f.Add([]byte("package main\n\nfunc main() {}"))
	f.Add([]byte("func Add(a, b int) int { return a + b }"))
	f.Add([]byte("type Foo struct { Name string }"))
	f.Add([]byte(""))  // empty input
	f.Add([]byte("\x00")) // null byte
	f.Add([]byte("package main\n\n/* " + string(make([]byte, 10000)))) // large unclosed comment

	tsParser := NewTreeSitterParser()

	f.Fuzz(func(t *testing.T, data []byte) {
		// The parser should not panic on any input
		result, err := tsParser.Parse(data, nil)
		if err != nil {
			// Error is acceptable, panic is not
			return
		}
		// Result should be valid
		if result == nil {
			t.Error("expected non-nil result on successful parse")
		}
	})
}

// FuzzParseTypeScript tests that the TypeScript parser does not panic on arbitrary input.
func FuzzParseTypeScript(f *testing.F) {
	f.Add([]byte("function main() {}"))
	f.Add([]byte("const x: number = 1;"))
	f.Add([]byte("interface Foo { name: string }"))
	f.Add([]byte(""))
	f.Add([]byte("\x00"))
	f.Add([]byte("function test() { return `" + string(make([]byte, 5000)) + "` }"))

	tsParser := NewTreeSitterParser()

	f.Fuzz(func(t *testing.T, data []byte) {
		result, err := tsParser.Parse(data, &parser.Options{Language: "typescript"})
		if err != nil {
			return
		}
		if result == nil {
			t.Error("expected non-nil result on successful parse")
		}
	})
}

// FuzzParsePython tests that the Python parser does not panic on arbitrary input.
func FuzzParsePython(f *testing.F) {
	f.Add([]byte("def main():\n    pass"))
	f.Add([]byte("class Foo:\n    def __init__(self):\n        pass"))
	f.Add([]byte("import os\nfrom typing import List"))
	f.Add([]byte(""))
	f.Add([]byte("\x00"))
	f.Add([]byte("def f():\n    " + string(make([]byte, 5000))))

	tsParser := NewTreeSitterParser()

	f.Fuzz(func(t *testing.T, data []byte) {
		result, err := tsParser.Parse(data, &parser.Options{Language: "python"})
		if err != nil {
			return
		}
		if result == nil {
			t.Error("expected non-nil result on successful parse")
		}
	})
}

// FuzzParseJava tests that the Java parser does not panic on arbitrary input.
func FuzzParseJava(f *testing.F) {
	f.Add([]byte("public class Main { public static void main(String[] args) {} }"))
	f.Add([]byte("interface Foo { void bar(); }"))
	f.Add([]byte("package com.example;"))
	f.Add([]byte(""))
	f.Add([]byte("\x00"))

	tsParser := NewTreeSitterParser()

	f.Fuzz(func(t *testing.T, data []byte) {
		result, err := tsParser.Parse(data, &parser.Options{Language: "java"})
		if err != nil {
			return
		}
		if result == nil {
			t.Error("expected non-nil result on successful parse")
		}
	})
}

// FuzzParseRust tests that the Rust parser does not panic on arbitrary input.
func FuzzParseRust(f *testing.F) {
	f.Add([]byte("fn main() {}"))
	f.Add([]byte("struct Foo { name: String }"))
	f.Add([]byte("use std::collections::HashMap;"))
	f.Add([]byte(""))
	f.Add([]byte("\x00"))

	tsParser := NewTreeSitterParser()

	f.Fuzz(func(t *testing.T, data []byte) {
		result, err := tsParser.Parse(data, &parser.Options{Language: "rust"})
		if err != nil {
			return
		}
		if result == nil {
			t.Error("expected non-nil result on successful parse")
		}
	})
}

// FuzzParseC tests that the C parser does not panic on arbitrary input.
func FuzzParseC(f *testing.F) {
	f.Add([]byte("int main() { return 0; }"))
	f.Add([]byte("struct Foo { int x; };"))
	f.Add([]byte("#include <stdio.h>"))
	f.Add([]byte(""))
	f.Add([]byte("\x00"))

	tsParser := NewTreeSitterParser()

	f.Fuzz(func(t *testing.T, data []byte) {
		result, err := tsParser.Parse(data, &parser.Options{Language: "c"})
		if err != nil {
			return
		}
		if result == nil {
			t.Error("expected non-nil result on successful parse")
		}
	})
}

// FuzzParseJSON tests JSON parsing (via JavaScript/TypeScript parser).
func FuzzParseJSON(f *testing.F) {
	f.Add([]byte(`{"key": "value"}`))
	f.Add([]byte(`[1, 2, 3]`))
	f.Add([]byte(`{"nested": {"a": 1}}`))
	f.Add([]byte(""))
	f.Add([]byte("\x00"))
	f.Add([]byte(`{"a": ` + string(make([]byte, 10000)) + `}`))

	tsParser := NewTreeSitterParser()

	f.Fuzz(func(t *testing.T, data []byte) {
		result, err := tsParser.Parse(data, &parser.Options{Language: "javascript"})
		if err != nil {
			return
		}
		if result == nil {
			t.Error("expected non-nil result on successful parse")
		}
	})
}
