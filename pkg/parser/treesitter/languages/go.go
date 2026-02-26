// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

// Capture name constants (must match treesitter package constants).
const (
	captureName      = "name"
	captureSignature = "signature"
	captureDoc       = "doc"
	captureKind      = "kind"
)

// GoQuery implements LanguageQuery for Go language.
type GoQuery struct {
	language *sitter.Language
	query    []byte
}

// NewGoQuery creates a new Go language query.
func NewGoQuery() *GoQuery {
	return &GoQuery{
		language: sitter.NewLanguage(tree_sitter_go.Language()),
		query:    []byte(goQueryPattern),
	}
}

// Language returns the Go Tree-sitter language.
func (q *GoQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Go query pattern.
func (q *GoQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for Go queries.
func (q *GoQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *GoQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_declaration": "function",
		"method_declaration":   "method",
		"type_declaration":     "type",
		"const_declaration":    "variable",
		"var_declaration":      "variable",
		"const_spec":           "variable",
		"var_spec":             "variable",
	}
}

// goQueryPattern is the Tree-sitter query for extracting Go signatures.
//
// Tree-sitter query syntax:
// - (@name) captures the identifier name
// - (@signature) captures the full declaration text
// - (@doc) captures the comment/documentation
// - (@kind) captures the node type for kind mapping
const goQueryPattern = `
; Function declarations
(function_declaration
  name: (identifier) @name
) @signature @kind

; Method declarations
(method_declaration
  name: (field_identifier) @name
) @signature @kind

; Type declarations (struct, interface, etc.)
(type_declaration
  (type_spec
    name: (type_identifier) @name
  )
) @signature @kind

; Package-level const specs (captures each const individually)
(const_spec
  name: (identifier) @name
) @signature @kind

; Package-level var specs (captures each var individually)
(var_spec
  name: (identifier) @name
) @signature @kind

; Comments (documentation)
(comment) @doc
`
