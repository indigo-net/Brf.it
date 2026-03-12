package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_ruby "github.com/tree-sitter/tree-sitter-ruby/bindings/go"
)

// RubyQuery implements LanguageQuery for Ruby language.
type RubyQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}

// NewRubyQuery creates a new Ruby language query.
func NewRubyQuery() *RubyQuery {
	return &RubyQuery{
		language: sitter.NewLanguage(tree_sitter_ruby.Language()),
		query:    []byte(rubyQueryPattern),
	}
}

// Language returns the Ruby Tree-sitter language.
func (q *RubyQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Ruby query pattern.
func (q *RubyQuery) Query() []byte {
	return q.query
}

var rubyKindMapping = map[string]string{
	"method":           "method",
	"singleton_method": "method",
	"class":            "class",
	"module":           "namespace",
	"assignment":       "variable",
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *RubyQuery) KindMapping() map[string]string {
	return rubyKindMapping
}

// ImportQuery returns the Ruby import query pattern.
func (q *RubyQuery) ImportQuery() []byte {
	return []byte(rubyImportQueryPattern)
}

// rubyImportQueryPattern is the Tree-sitter query for extracting Ruby require statements.
// Captures both require "lib" and require_relative "lib" calls.
const rubyImportQueryPattern = `
; require "library" / require_relative "library"
(call
  method: (identifier) @_fn
  arguments: (argument_list
    (string)
  )
) @import_path
(#match? @_fn "^require")
`

// rubyQueryPattern is the Tree-sitter query for extracting Ruby signatures.
//
// Captures:
// - (@name) identifier name
// - (@signature) full declaration text
// - (@doc) comment/documentation
// - (@kind) node type for kind mapping
const rubyQueryPattern = `
; Instance methods and top-level functions: def foo(args) ... end
(method
  name: (identifier) @name
) @signature @kind

; Class methods: def self.foo(args) ... end
(singleton_method
  name: (identifier) @name
) @signature @kind

; Class definitions: class Foo ... end
(class
  name: (constant) @name
) @signature @kind

; Module definitions: module Foo ... end
(module
  name: (constant) @name
) @signature @kind

; Top-level constant assignments: FOO = value
(program
  (assignment
    left: (constant) @name
  ) @signature @kind
)

; Comments
(comment) @doc
`
