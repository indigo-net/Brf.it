package languages_test

import (
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser"

	// Blank import to register tree-sitter parsers
	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"
)

func TestGoCallExtraction(t *testing.T) {
	content := []byte(`package main

import "fmt"

func main() {
	fmt.Println("hello")
	helper()
}

func helper() {
	fmt.Sprintf("world")
}
`)

	p, ok := parser.GetParser("go")
	if !ok {
		t.Fatal("go parser not found")
	}

	result, err := p.Parse(content, &parser.Options{
		Language:     "go",
		IncludeCalls: true,
	})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(result.Calls) == 0 {
		t.Fatal("expected calls to be extracted")
	}

	// Verify we find expected calls
	found := make(map[string]bool)
	for _, call := range result.Calls {
		found[call.Callee] = true

		// Verify line numbers are positive
		if call.Line <= 0 {
			t.Errorf("expected positive line number, got %d for callee %q", call.Line, call.Callee)
		}
	}

	if !found["Println"] {
		t.Error("expected to find call to Println")
	}
	if !found["helper"] {
		t.Error("expected to find call to helper")
	}
	if !found["Sprintf"] {
		t.Error("expected to find call to Sprintf")
	}

	// Verify caller attribution
	for _, call := range result.Calls {
		switch call.Callee {
		case "Println", "helper":
			if call.Caller != "main" {
				t.Errorf("expected caller 'main' for %q, got %q", call.Callee, call.Caller)
			}
		case "Sprintf":
			if call.Caller != "helper" {
				t.Errorf("expected caller 'helper' for Sprintf, got %q", call.Caller)
			}
		}
	}
}

func TestTypeScriptCallExtraction(t *testing.T) {
	content := []byte(`function greet(name: string): void {
  console.log("Hello, " + name);
  format(name);
}

function format(s: string): string {
  return s.trim();
}
`)

	p, ok := parser.GetParser("typescript")
	if !ok {
		t.Fatal("typescript parser not found")
	}

	result, err := p.Parse(content, &parser.Options{
		Language:     "typescript",
		IncludeCalls: true,
	})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(result.Calls) == 0 {
		t.Fatal("expected calls to be extracted")
	}

	found := make(map[string]bool)
	for _, call := range result.Calls {
		found[call.Callee] = true
	}

	if !found["log"] {
		t.Error("expected to find call to log")
	}
	if !found["format"] {
		t.Error("expected to find call to format")
	}
	if !found["trim"] {
		t.Error("expected to find call to trim")
	}
}

func TestPythonCallExtraction(t *testing.T) {
	content := []byte(`def main():
    print("hello")
    result = process("data")

def process(data):
    return data.strip()
`)

	p, ok := parser.GetParser("python")
	if !ok {
		t.Fatal("python parser not found")
	}

	result, err := p.Parse(content, &parser.Options{
		Language:     "python",
		IncludeCalls: true,
	})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(result.Calls) == 0 {
		t.Fatal("expected calls to be extracted")
	}

	found := make(map[string]bool)
	for _, call := range result.Calls {
		found[call.Callee] = true
	}

	if !found["print"] {
		t.Error("expected to find call to print")
	}
	if !found["process"] {
		t.Error("expected to find call to process")
	}
	if !found["strip"] {
		t.Error("expected to find call to strip")
	}
}

func TestCallExtractionDisabledByDefault(t *testing.T) {
	content := []byte(`package main

func main() {
	helper()
}
`)

	p, ok := parser.GetParser("go")
	if !ok {
		t.Fatal("go parser not found")
	}

	result, err := p.Parse(content, &parser.Options{
		Language: "go",
		// IncludeCalls is false by default
	})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(result.Calls) != 0 {
		t.Errorf("expected no calls when IncludeCalls is false, got %d", len(result.Calls))
	}
}

func TestCallExtractionTopLevel(t *testing.T) {
	// C code with a top-level call (in main, but testing that calls outside
	// any function return empty caller)
	content := []byte(`#include <stdio.h>

void greet() {
    printf("hello");
}
`)

	p, ok := parser.GetParser("c")
	if !ok {
		t.Fatal("c parser not found")
	}

	result, err := p.Parse(content, &parser.Options{
		Language:       "c",
		IncludeCalls:   true,
		IncludePrivate: true,
	})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(result.Calls) == 0 {
		t.Fatal("expected calls to be extracted")
	}

	found := false
	for _, call := range result.Calls {
		if call.Callee == "printf" {
			found = true
			if call.Caller != "greet" {
				t.Errorf("expected caller 'greet' for printf, got %q", call.Caller)
			}
		}
	}
	if !found {
		t.Error("expected to find call to printf")
	}
}
