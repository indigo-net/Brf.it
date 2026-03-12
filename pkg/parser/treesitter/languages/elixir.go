// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	tree_sitter_elixir "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/elixir"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// ElixirQuery implements LanguageQuery for Elixir language.
type ElixirQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}

// NewElixirQuery creates a new Elixir language query.
func NewElixirQuery() *ElixirQuery {
	return &ElixirQuery{
		language: sitter.NewLanguage(tree_sitter_elixir.Language()),
		query:    []byte(elixirQueryPattern),
	}
}

// Language returns the Elixir Tree-sitter language.
func (q *ElixirQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Elixir query pattern.
func (q *ElixirQuery) Query() []byte {
	return q.query
}

// KindMapping returns the mapping from node types to Signature kinds.
// Elixir uses generic node types (call, unary_operator) which are refined
// on the Go side in parser.go based on the signature text.
func (q *ElixirQuery) KindMapping() map[string]string {
	return map[string]string{
		"call":           "function",
		"unary_operator": "type",
	}
}

// ImportQuery returns the Elixir import query pattern.
func (q *ElixirQuery) ImportQuery() []byte {
	return []byte(elixirImportQueryPattern)
}

// elixirImportQueryPattern is the Tree-sitter query for extracting Elixir
// import/alias/use/require statements. These are all plain `call` nodes
// in tree-sitter-elixir.
const elixirImportQueryPattern = `
; import statements: import Module
(call
  target: (identifier)
  (arguments
    (alias))) @import_path

; import with options: import Module, only: [...]
(call
  target: (identifier)
  (arguments
    (alias)
    (keywords))) @import_path
`

// elixirQueryPattern is the Tree-sitter query for extracting Elixir signatures.
//
// Tree-sitter query syntax:
// - (@name) captures the identifier name
// - (@signature) captures the full declaration text
// - (@doc) captures the comment/documentation
// - (@kind) captures the node type for kind mapping
//
// In tree-sitter-elixir, all definition macros (def, defmodule, etc.) are
// represented as generic `call` nodes. Module attributes (@spec, @type, etc.)
// are `unary_operator` nodes. Go-side refinement in parser.go distinguishes
// between definition types by inspecting the signature text prefix.
//
// Elixir declaration types captured:
// - defmodule, defprotocol, defimpl (module-level definitions)
// - def, defp (public/private functions)
// - defmacro, defmacrop (public/private macros)
// - defguard, defguardp (public/private guards)
// - defdelegate (delegated functions)
// - defstruct (struct field definitions)
// - @spec, @type, @typep, @opaque, @callback (type specifications)
const elixirQueryPattern = `
; Module/protocol/impl definitions: defmodule MyModule do...end
(call
  target: (identifier)
  (arguments
    (alias) @name)
  (do_block)
) @signature @kind

; defimpl with keyword options: defimpl Protocol, for: Module do...end
(call
  target: (identifier)
  (arguments
    (alias) @name
    (keywords))
  (do_block)
) @signature @kind

; Function/macro definitions with arguments: def foo(args) do...end
(call
  target: (identifier)
  (arguments
    (call
      target: (identifier) @name))
  (do_block)
) @signature @kind

; Function/macro definitions with guard clause: def foo(args) when guard do...end
(call
  target: (identifier)
  (arguments
    (binary_operator
      left: (call
        target: (identifier) @name)
      operator: "when"))
  (do_block)
) @signature @kind

; Zero-arity function definitions: def foo do...end
(call
  target: (identifier)
  (arguments
    (identifier) @name)
  (do_block)
) @signature @kind

; Zero-arity function with guard: def foo when guard do...end
(call
  target: (identifier)
  (arguments
    (binary_operator
      left: (identifier) @name
      operator: "when"))
  (do_block)
) @signature @kind

; Guard definitions without do_block: defguard is_positive(x) when ...
(call
  target: (identifier)
  (arguments
    (binary_operator
      left: (call
        target: (identifier) @name)
      operator: "when"))
) @signature @kind

; defdelegate: defdelegate foo(args), to: Bar
(call
  target: (identifier)
  (arguments
    (call
      target: (identifier) @name)
    (keywords))
) @signature @kind

; defstruct with list: defstruct [:field1, :field2]
(call
  target: (identifier) @name
  (arguments
    (list))
) @signature @kind

; defstruct with keywords: defstruct field: default_value
(call
  target: (identifier) @name
  (arguments
    (keywords))
) @signature @kind

; Module attributes: @spec, @type, @typep, @opaque, @callback
(unary_operator
  operator: "@"
  operand: (call
    target: (identifier) @name)
) @signature @kind

; Line comments
(comment) @doc
`
