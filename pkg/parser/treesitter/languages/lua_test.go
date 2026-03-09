package languages

import (
	"testing"

	tree_sitter_lua "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/lua"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// extractLuaNames is a test helper that parses Lua code and returns
// all captured @name values from the query matches.
func extractLuaNames(t *testing.T, code []byte) map[string]bool {
	t.Helper()
	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_lua.Language())
	parser.SetLanguage(lang)
	tree := parser.Parse(code, nil)
	defer tree.Close()
	query := NewLuaQuery()
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

func TestLuaQueryLanguage(t *testing.T) {
	query := NewLuaQuery()
	lang := query.Language()

	if lang == nil {
		t.Fatal("expected non-nil language")
	}
}

func TestLuaQueryPattern(t *testing.T) {
	query := NewLuaQuery()
	pattern := query.Query()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_lua.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid query pattern: %v", err)
	}
	defer q.Close()
}

func TestLuaQueryImportPattern(t *testing.T) {
	query := NewLuaQuery()
	pattern := query.ImportQuery()

	if len(pattern) == 0 {
		t.Fatal("expected non-empty import query pattern")
	}

	// Compile query to verify syntax
	lang := sitter.NewLanguage(tree_sitter_lua.Language())
	q, err := sitter.NewQuery(lang, string(pattern))
	if err != nil {
		t.Fatalf("invalid import query pattern: %v", err)
	}
	defer q.Close()
}

func TestLuaQueryExtractFunction(t *testing.T) {
	code := []byte(`
function greet(name)
    print("Hello, " .. name)
end

function add(a, b)
    return a + b
end

function calculate(x, y, op)
    if op == "add" then
        return x + y
    end
    return 0
end
`)

	foundNames := extractLuaNames(t, code)

	expectedNames := []string{"greet", "add", "calculate"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find function '%s'", expected)
		}
	}
}

func TestLuaQueryExtractLocalFunction(t *testing.T) {
	code := []byte(`
local function helper()
    return 42
end

local function validate(input)
    return input ~= nil
end

local function format(template, ...)
    return string.format(template, ...)
end
`)

	foundNames := extractLuaNames(t, code)

	expectedNames := []string{"helper", "validate", "format"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find local function '%s'", expected)
		}
	}
}

func TestLuaQueryExtractModuleFunction(t *testing.T) {
	code := []byte(`
local M = {}

function M.create(name)
    return { name = name }
end

function M.new(name, age)
    return setmetatable({name = name, age = age}, M)
end

function M.validate(obj)
    return obj.name ~= nil
end
`)

	foundNames := extractLuaNames(t, code)

	expectedNames := []string{"M", "create", "new", "validate"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find module function '%s'", expected)
		}
	}
}

func TestLuaQueryExtractMethod(t *testing.T) {
	code := []byte(`
local M = {}

function M:init(name)
    self.name = name
end

function M:destroy()
    self.name = nil
end

function M:getName()
    return self.name
end
`)

	foundNames := extractLuaNames(t, code)

	expectedNames := []string{"M", "init", "destroy", "getName"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find method '%s'", expected)
		}
	}
}

func TestLuaQueryExtractTableAssignment(t *testing.T) {
	code := []byte(`
local M = {}
local Config = {}
local Utils = {}
`)

	foundNames := extractLuaNames(t, code)

	expectedNames := []string{"M", "Config", "Utils"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find table assignment '%s'", expected)
		}
	}
}

func TestLuaQueryExtractFunctionAssignment(t *testing.T) {
	code := []byte(`
local foo = function()
    return "bar"
end

local transform = function(x)
    return x * 2
end
`)

	foundNames := extractLuaNames(t, code)

	expectedNames := []string{"foo", "transform"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find function assignment '%s'", expected)
		}
	}
}

func TestLuaQueryExtractImport(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_lua.Language())
	parser.SetLanguage(lang)

	code := []byte(`
local json = require("json")
local utils = require("app.utils")
local lfs = require("lfs")
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewLuaQuery()
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

	// Should find 3 require declarations
	if count < 3 {
		t.Errorf("expected at least 3 import declarations, got %d", count)
	}
}

// TestLuaQueryNonRequireFalsePositive verifies that non-require() function calls
// such as pcall("...") do not produce import false positives after Go-side filtering.
func TestLuaQueryNonRequireFalsePositive(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_lua.Language())
	parser.SetLanguage(lang)

	// pcall and xpcall are valid function calls with string arguments but must NOT
	// be treated as imports. Only require() should produce import matches.
	code := []byte(`
local json = require("json")
local ok = pcall("not_a_module")
local status = xpcall("handler", "msg")
local lfs = require("lfs")
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewLuaQuery()
	q, err := sitter.NewQuery(lang, string(query.ImportQuery()))
	if err != nil {
		t.Fatalf("failed to create import query: %v", err)
	}
	defer q.Close()

	captureNames := q.CaptureNames()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(q, tree.RootNode(), code)

	// Simulate the Go-side filtering that extractImports() performs:
	// skip matches where @_fn != "require".
	requireCount := 0
	for {
		match := matches.Next()
		if match == nil {
			break
		}
		fnName := ""
		for _, c := range match.Captures {
			if captureNames[c.Index] == "_fn" {
				fnName = string(code[c.Node.StartByte():c.Node.EndByte()])
			}
		}
		if fnName == "require" {
			requireCount++
		}
	}

	// Should match only the 2 require() calls, not pcall/xpcall.
	if requireCount != 2 {
		t.Errorf("expected exactly 2 require() matches after filtering, got %d", requireCount)
	}
}

func TestLuaQueryExtractDoc(t *testing.T) {
	parser := sitter.NewParser()
	defer parser.Close()

	lang := sitter.NewLanguage(tree_sitter_lua.Language())
	parser.SetLanguage(lang)

	code := []byte(`
--- This is a LuaDoc comment.
-- This is a regular comment.
--[[ This is a block comment ]]
function foo() end
`)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	query := NewLuaQuery()
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

	// Should find 3 comments (LuaDoc, regular, block)
	if docCount < 3 {
		t.Errorf("expected at least 3 doc captures, got %d", docCount)
	}
}

func TestLuaQueryKindMapping(t *testing.T) {
	query := NewLuaQuery()
	mapping := query.KindMapping()

	expectedMappings := map[string]string{
		"function_declaration": "function",
		"variable_declaration": "variable",
		"assignment_statement": "variable",
	}

	for nodeType, expected := range expectedMappings {
		if actual, ok := mapping[nodeType]; !ok {
			t.Errorf("expected mapping for '%s'", nodeType)
		} else if actual != expected {
			t.Errorf("expected '%s' -> '%s', got '%s'", nodeType, expected, actual)
		}
	}
}

func TestLuaQueryCaptures(t *testing.T) {
	query := NewLuaQuery()
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

func TestLuaQueryExtractMixed(t *testing.T) {
	code := []byte(`
local M = {}

--- Creates a new instance.
-- @param name string The name
function M.new(name)
    return setmetatable({name = name}, M)
end

function M:greet()
    return "Hello, " .. self.name
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
`)

	foundNames := extractLuaNames(t, code)

	expectedNames := []string{"M", "new", "greet", "helper", "globalFunc", "callback"}
	for _, expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("expected to find '%s'", expected)
		}
	}
}
