// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"
)

// RustQuery implements LanguageQuery for Rust language.
type RustQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}

// NewRustQuery creates a new Rust language query.
func NewRustQuery() *RustQuery {
	return &RustQuery{
		language: sitter.NewLanguage(tree_sitter_rust.Language()),
		query:    []byte(rustQueryPattern),
	}
}

// Language returns the Rust Tree-sitter language.
func (q *RustQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Rust query pattern.
func (q *RustQuery) Query() []byte {
	return q.query
}

var rustKindMapping = map[string]string{
	"function_item":           "function",
	"struct_item":             "struct",
	"enum_item":               "enum",
	"trait_item":              "trait",
	"type_item":               "type",
	"impl_item":               "impl",
	"const_item":              "variable",
	"static_item":             "variable",
	"mod_item":                "namespace",
	"macro_definition":        "macro",
	"foreign_mod_item":        "namespace",
	"union_item":              "struct",
	"associated_type":         "type",
	"function_signature_item": "function",
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *RustQuery) KindMapping() map[string]string {
	return rustKindMapping
}

// ImportQuery returns the Rust import query pattern.
func (q *RustQuery) ImportQuery() []byte {
	return []byte(rustImportQueryPattern)
}

// rustImportQueryPattern is the Tree-sitter query for extracting Rust use statements.
const rustImportQueryPattern = `
; Use declarations (capture full statement)
(use_declaration) @import_path

; Extern crate declarations
(extern_crate_declaration) @import_path
`

// rustQueryPattern is the Tree-sitter query for extracting Rust signatures.
//
// Tree-sitter query syntax:
// - (@name) captures the identifier name
// - (@signature) captures the full declaration text
// - (@doc) captures the comment/documentation
// - (@kind) captures the node type for kind mapping
const rustQueryPattern = `
; Functions (including async, unsafe, const, extern)
(function_item
  name: (identifier) @name
) @signature @kind

; Struct declarations
(struct_item
  name: (type_identifier) @name
) @signature @kind

; Enum declarations
(enum_item
  name: (type_identifier) @name
) @signature @kind

; Trait declarations
(trait_item
  name: (type_identifier) @name
) @signature @kind

; Type aliases
(type_item
  name: (type_identifier) @name
) @signature @kind

; Impl blocks (capture the whole impl signature)
(impl_item
  type: (type_identifier) @name
) @signature @kind

; Impl blocks for generic types
(impl_item
  type: (generic_type
    type: (type_identifier) @name
  )
) @signature @kind

; Trait impl blocks (impl Trait for Type)
(impl_item
  trait: (type_identifier)
  type: (type_identifier) @name
) @signature @kind

; Constants
(const_item
  name: (identifier) @name
) @signature @kind

; Statics
(static_item
  name: (identifier) @name
) @signature @kind

; Modules
(mod_item
  name: (identifier) @name
) @signature @kind

; Macro definitions (macro_rules!)
(macro_definition
  name: (identifier) @name
) @signature @kind

; Union declarations
(union_item
  name: (type_identifier) @name
) @signature @kind

; Foreign mod (extern "C" blocks)
(foreign_mod_item) @signature @kind

; Associated types in traits
(associated_type
  name: (type_identifier) @name
) @signature @kind

; Function signatures in traits (without body)
(function_signature_item
  name: (identifier) @name
) @signature @kind

; Doc comments (/// and //!)
(line_comment) @doc

; Block doc comments (/** and /*!)
(block_comment) @doc
`
