// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_swift "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/swift"
)

// SwiftQuery implements LanguageQuery for Swift language.
type SwiftQuery struct {
	language *sitter.Language
	query    []byte
}

// NewSwiftQuery creates a new Swift language query.
func NewSwiftQuery() *SwiftQuery {
	return &SwiftQuery{
		language: sitter.NewLanguage(tree_sitter_swift.Language()),
		query:    []byte(swiftQueryPattern),
	}
}

// Language returns the Swift Tree-sitter language.
func (q *SwiftQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Swift query pattern.
func (q *SwiftQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for Swift queries.
func (q *SwiftQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *SwiftQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_declaration":          "function",
		"class_declaration":             "class",
		"protocol_declaration":          "interface",
		"typealias_declaration":         "type",
		"property_declaration":          "variable",
		"init_declaration":              "constructor",
		"deinit_declaration":            "destructor",
		"subscript_declaration":         "method",
		"operator_declaration":          "function",
		"protocol_function_declaration": "function",
	}
}

// ImportQuery returns the Swift import query pattern.
func (q *SwiftQuery) ImportQuery() []byte {
	return []byte(swiftImportQueryPattern)
}

// swiftImportQueryPattern is the Tree-sitter query for extracting Swift import statements.
const swiftImportQueryPattern = `
; Import declarations (capture full statement)
(import_declaration) @import_path
`

// swiftQueryPattern is the Tree-sitter query for extracting Swift signatures.
//
// Tree-sitter query syntax:
// - (@name) captures the identifier name
// - (@signature) captures the full declaration text
// - (@doc) captures the comment/documentation
// - (@kind) captures the node type for kind mapping
const swiftQueryPattern = `
; Functions
(function_declaration
  name: (simple_identifier) @name
) @signature @kind

; Classes, Structs, Enums (all use class_declaration node type)
(class_declaration
  name: (type_identifier) @name
) @signature @kind

; Extensions (name is in user_type child)
(class_declaration
  name: (user_type
    (type_identifier) @name
  )
) @signature @kind

; Protocol declarations
(protocol_declaration
  name: (type_identifier) @name
) @signature @kind

; Type aliases
(typealias_declaration
  name: (type_identifier) @name
) @signature @kind

; Properties (let/var)
(property_declaration
  name: (pattern
    (simple_identifier) @name
  )
) @signature @kind

; Initializers
(init_declaration) @signature @kind

; Deinitializers
(deinit_declaration) @signature @kind

; Subscript declarations
(subscript_declaration) @signature @kind

; Operator declarations
(operator_declaration
  (custom_operator) @name
) @signature @kind

; Protocol function declarations (methods in protocol body)
(protocol_function_declaration
  name: (simple_identifier) @name
) @signature @kind

; Doc comments (/// style)
(comment) @doc

; Multiline comments (/** style)
(multiline_comment) @doc
`
