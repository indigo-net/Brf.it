package languages

import (
	"testing"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_bash "github.com/tree-sitter/tree-sitter-bash/bindings/go"
)

func TestBashQueryLanguage(t *testing.T) {
	query := NewBashQuery()
	if query.Language() == nil {
		t.Error("Language() returned nil")
	}
}

func TestBashQueryPattern(t *testing.T) {
	query := NewBashQuery()
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_bash.Language())
	parser.SetLanguage(lang)

	// Create query to verify pattern compiles
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to compile query: %v", err)
	}
	defer q.Close()
}

func TestBashQueryImportPattern(t *testing.T) {
	query := NewBashQuery()
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_bash.Language())
	parser.SetLanguage(lang)

	// Create query to verify pattern compiles
	q, err := sitter.NewQuery(lang, string(query.ImportQuery()))
	if err != nil {
		t.Fatalf("failed to compile import query: %v", err)
	}
	defer q.Close()
}

func TestBashQueryExtractFunction(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_bash.Language())
	parser.SetLanguage(lang)

	code := []byte(`#!/bin/bash

# Simple function
function greet {
    echo "Hello, World!"
}

# Function with parentheses
foo() {
    echo "foo"
}

# Function with arguments
function process_args() {
    local input="$1"
    echo "$input"
}
`)
	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewBashQuery()
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

	expected := []string{"greet", "foo", "process_args"}
	for _, exp := range expected {
		if !foundNames[exp] {
			t.Errorf("expected to find function '%s'", exp)
		}
	}
}

func TestBashQueryExtractVariables(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_bash.Language())
	parser.SetLanguage(lang)

	code := []byte(`#!/bin/bash

# Variable assignments
NAME="value"
COUNT=42
PATH_VAR="/usr/local/bin"

# Declaration commands
declare VERBOSE
local DEBUG
readonly VERSION="1.0"
`)
	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewBashQuery()
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

	expected := []string{"NAME", "COUNT", "PATH_VAR", "VERBOSE", "DEBUG", "VERSION"}
	for _, exp := range expected {
		if !foundNames[exp] {
			t.Errorf("expected to find variable '%s'", exp)
		}
	}
}

func TestBashQueryExtractComments(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_bash.Language())
	parser.SetLanguage(lang)

	code := []byte(`#!/bin/bash
# This is a comment
# Another comment

function test {
    echo "test"
}
`)
	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewBashQuery()
	q, err := sitter.NewQuery(lang, string(query.Query()))
	if err != nil {
		t.Fatalf("failed to create query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	commentCount := 0
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, c := range match.Captures {
			if captureNames[c.Index] == "doc" {
				commentCount++
			}
		}
	}

	if commentCount < 2 {
		t.Errorf("expected at least 2 comments, got %d", commentCount)
	}
}

func TestBashQueryExtractSourceStatements(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_bash.Language())
	parser.SetLanguage(lang)

	code := []byte(`#!/bin/bash

source /path/to/lib.sh
. ./config.sh
source "./utils.sh"
`)
	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewBashQuery()
	q, err := sitter.NewQuery(lang, string(query.ImportQuery()))
	if err != nil {
		t.Fatalf("failed to create import query: %v", err)
	}
	defer q.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()
	matches := qc.Matches(q, tree.RootNode(), code)

	captureNames := q.CaptureNames()
	imports := []string{}
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		for _, c := range match.Captures {
			if captureNames[c.Index] == "import_path" {
				imports = append(imports, string(code[c.Node.StartByte():c.Node.EndByte()]))
			}
		}
	}

	if len(imports) < 1 {
		t.Errorf("expected at least 1 source statement, got %d", len(imports))
	}
}

func TestBashQueryKindMapping(t *testing.T) {
	query := NewBashQuery()
	mapping := query.KindMapping()

	expectedKinds := []string{"function_definition", "variable_assignment", "declaration_command"}
	for _, kind := range expectedKinds {
		if _, ok := mapping[kind]; !ok {
			t.Errorf("expected kind mapping for '%s'", kind)
		}
	}
}

func TestBashQueryCaptures(t *testing.T) {
	query := NewBashQuery()
	captures := query.Captures()

	expectedCaptures := []string{"name", "signature", "doc", "kind"}
	for _, exp := range expectedCaptures {
		found := false
		for _, c := range captures {
			if c == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected capture '%s'", exp)
		}
	}
}
