// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

// JavaQuery implements LanguageQuery for Java language.
type JavaQuery struct {
	language *sitter.Language
	query    []byte
}

// NewJavaQuery creates a new Java language query.
func NewJavaQuery() *JavaQuery {
	return &JavaQuery{
		language: sitter.NewLanguage(tree_sitter_java.Language()),
		query:    []byte(javaQueryPattern),
	}
}

// Language returns the Java Tree-sitter language.
func (q *JavaQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Java query pattern.
func (q *JavaQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for Java queries.
func (q *JavaQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *JavaQuery) KindMapping() map[string]string {
	return map[string]string{
		"class_declaration":           "class",
		"interface_declaration":       "interface",
		"method_declaration":          "method",
		"constructor_declaration":     "constructor",
		"enum_declaration":            "enum",
		"annotation_type_declaration": "annotation",
		"record_declaration":          "record",
		"field_declaration":           "field",
	}
}

// ImportQuery returns the Java import query pattern.
func (q *JavaQuery) ImportQuery() []byte {
	return []byte(javaImportQueryPattern)
}

// javaImportQueryPattern is the Tree-sitter query for extracting Java imports.
const javaImportQueryPattern = `
; import statements (capture full declaration)
(import_declaration) @import_path
`

// javaQueryPattern is the Tree-sitter query for extracting Java signatures.
//
// Tree-sitter query syntax:
// - (@name) captures the identifier name
// - (@signature) captures the full declaration text
// - (@doc) captures the comment/documentation
// - (@kind) captures the node type for kind mapping
const javaQueryPattern = `
; Class declarations (includes inner classes)
(class_declaration
  name: (identifier) @name
) @signature @kind

; Interface declarations
(interface_declaration
  name: (identifier) @name
) @signature @kind

; Method declarations
(method_declaration
  name: (identifier) @name
) @signature @kind

; Constructor declarations
(constructor_declaration
  name: (identifier) @name
) @signature @kind

; Enum declarations
(enum_declaration
  name: (identifier) @name
) @signature @kind

; Annotation type declarations (@interface)
(annotation_type_declaration
  name: (identifier) @name
) @signature @kind

; Record declarations (Java 14+)
(record_declaration
  name: (identifier) @name
) @signature @kind

; Field declarations (static fields filtered in parser.go)
(field_declaration
  (variable_declarator
    name: (identifier) @name
  )
) @signature @kind

; Comments (Javadoc and regular)
(line_comment) @doc
(block_comment) @doc
`
