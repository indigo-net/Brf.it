package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_bash "github.com/tree-sitter/tree-sitter-bash/bindings/go"
)

// extractShellNames is a test helper that parses Shell/Bash code and returns
// all captured @name values from the query matches.
func extractShellNames(t *testing.T, code []byte) map[string]bool {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_bash.Language())
	parser.SetLanguage(lang)
	tree := parser.Parse(code, nil)
	defer tree.Close()
	query := NewShellQuery()
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

func TestShellQueryLanguage(t *testing.T) {
	query := NewShellQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestShellQueryPattern(t *testing.T) {
	query := NewShellQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_bash.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
	defer q.Close()
}

func TestShellQueryImportPattern(t *testing.T) {
	query := NewShellQuery()
	pattern := query.ImportQuery()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty import query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_bash.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid import query pattern: %v", err)
	}
	defer q.Close()
}

func TestShellQueryExtractFunction(t *testing.T) {
	code := []byte(`
#!/bin/bash

function greet() {
    echo "Hello, $1"
}

function add() {
    echo $(($1 + $2))
}
`)

	foundNames := extractShellNames(t, code)

	expectedNames := []string{"greet", "add"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find function '%s'", expected)
		}
	}
}

func TestShellQueryExtractFunctionWithoutKeyword(t *testing.T) {
	code := []byte(`
#!/bin/bash

deploy() {
    echo "Deploying..."
}

build() {
    npm run build
}

test_all() {
    npm test
}
`)

	foundNames := extractShellNames(t, code)

	expectedNames := []string{"deploy", "build", "test_all"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find function '%s'", expected)
		}
	}
}

func TestShellQueryExtractVariable(t *testing.T) {
	code := []byte(`
#!/bin/bash

NAME="John"
AGE=30
PI=3.14
CURRENT_DIR=$(pwd)
FILES=$(ls -la)
`)

	foundNames := extractShellNames(t, code)

	expectedNames := []string{"NAME", "AGE", "PI", "CURRENT_DIR", "FILES"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find variable '%s'", expected)
		}
	 }
}

func TestShellQueryExtractImport(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_bash.Language())
	parser.SetLanguage(lang)

	code := []byte(`
#!/bin/bash

source /path/to/utils.sh
. /path/to/config.sh
source ./helpers.sh
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewShellQuery()
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

	// Should find import statements (all commands including source and .)
	if count < 1 {
		t.Errorf("expected at least 1 import declaration, got %d", count)
	}
}

func TestShellQueryExtractDoc(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_bash.Language())
	parser.SetLanguage(lang)

	code := []byte(`
#!/bin/bash

# This is a comment
# Multi-line comment

function foo() {
    echo "bar"
}
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewShellQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)
	captureNames := q.CaptureNames()

	docCount := 0
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, c := range match.Captures {
			if captureNames[c.Index] == "doc" {
				docCount++
			}
		}
	}

	// Should find comments
	if docCount < 2 {
		t.Errorf("expected at least 2 doc captures, got %d", docCount)
	}
}

func TestShellQueryKindMapping(t *testing.T) {
	query := NewShellQuery()
	mapping := query.KindMapping()

	expectedMappings := map[string]string{
		"function_definition": "function",
		"variable_assignment": "variable",
	}

	for nodeType, expected := range expectedMappings {
		if actual, ok := mapping[nodeType]; !ok {
			t.Errorf("expected mapping for '%s'", nodeType)
		} else if actual != expected {
			t.Errorf("expected '%s' -> '%s', got '%s'", nodeType, expected, actual)
		}
	}
}

func TestShellQueryCaptures(t *testing.T) {
	query := NewShellQuery()
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

func TestShellQueryExtractMixed(t *testing.T) {
	code := []byte(`
#!/bin/bash

# Configuration
APP_NAME="myapp"
VERSION="1.0.0"

# Deploy function
function deploy() {
    echo "Deploying $APP_NAME v$VERSION"
}

# Build function
build() {
    npm run build
}

# Test function
function run_tests() {
    npm test
}

ENVIRONMENT="production"
`)

	foundNames := extractShellNames(t, code)

	expectedNames := []string{"APP_NAME", "VERSION", "deploy", "build", "run_tests", "ENVIRONMENT"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}
