// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	tree_sitter_lua "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/lua"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// LuaQuery implements LanguageQuery for Lua language.
type LuaQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}

// NewLuaQuery creates a new Lua language query.
func NewLuaQuery() *LuaQuery {
	return &LuaQuery{
		language: sitter.NewLanguage(tree_sitter_lua.Language()),
		query:    []byte(luaQueryPattern),
	}
}

// Language returns the Lua Tree-sitter language.
func (q *LuaQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Lua query pattern.
func (q *LuaQuery) Query() []byte {
	return q.query
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *LuaQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_declaration": "function",
		"variable_declaration": "variable",
		"assignment_statement": "variable",
	}
}

// ImportQuery returns the Lua import query pattern.
func (q *LuaQuery) ImportQuery() []byte {
	return []byte(luaImportQueryPattern)
}

// luaImportQueryPattern is the Tree-sitter query for extracting Lua require() calls.
// Captures the full variable_declaration containing require() as @import_path.
// The (#eq? @_fn "require") predicate filters to only match require() calls.
// Go-side filtering in extractImports() additionally checks @_fn == "require" as a fallback
// in case the tree-sitter predicate is not evaluated by the Go binding.
const luaImportQueryPattern = `
; local json = require("json")
(variable_declaration
  (assignment_statement
    (expression_list
      value: (function_call
        name: (identifier) @_fn
        arguments: (arguments (string))
      )
    )
  )
) @import_path
(#eq? @_fn "require")
`

// luaQueryPattern is the Tree-sitter query for extracting Lua signatures.
//
// Tree-sitter query syntax:
// - (@name) captures the identifier name
// - (@signature) captures the full declaration text
// - (@doc) captures the comment/documentation
// - (@kind) captures the node type for kind mapping
//
// Lua function declarations use a unified function_declaration node type
// for global, local, module (M.func), and method (M:func) patterns.
// The refineLuaFunctionKind() function in parser.go distinguishes between
// function, local_function, module_function, and method.
const luaQueryPattern = `
; Function declarations (global, local, module, method)
; Covers: function foo(), local function foo(), function M.foo(), function M:foo()
(function_declaration
  name: [
    (identifier) @name
    (dot_index_expression field: (identifier) @name)
    (method_index_expression method: (identifier) @name)
  ]
) @signature @kind

; Variable declarations with function assignment: local foo = function() end
(variable_declaration
  (assignment_statement
    (variable_list
      name: (identifier) @name)
    (expression_list
      value: (function_definition))
  )
) @signature @kind

; Variable declarations with table constructor: local M = {}
(variable_declaration
  (assignment_statement
    (variable_list
      name: (identifier) @name)
    (expression_list
      value: (table_constructor))
  )
) @signature @kind

; Comments (LuaDoc --- and regular --)
(comment) @doc
`
