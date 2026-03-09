// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	tree_sitter_kotlin "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/kotlin"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// KotlinQuery implements LanguageQuery for Kotlin language.
type KotlinQuery struct {
	language *sitter.Language
	query    []byte
}

// NewKotlinQuery creates a new Kotlin language query.
func NewKotlinQuery() *KotlinQuery {
	return &KotlinQuery{
		language: sitter.NewLanguage(tree_sitter_kotlin.Language()),
		query:    []byte(kotlinQueryPattern),
	}
}

// Language returns the Kotlin Tree-sitter language.
func (q *KotlinQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Kotlin query pattern.
func (q *KotlinQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for Kotlin queries.
func (q *KotlinQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *KotlinQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_declaration":  "function",
		"class_declaration":     "class",
		"object_declaration":    "class",
		"companion_object":      "class",
		"property_declaration":  "variable",
		"type_alias":            "type",
		"enum_entry":            "variable",
		"secondary_constructor": "constructor",
	}
}

// ImportQuery returns the Kotlin import query pattern.
func (q *KotlinQuery) ImportQuery() []byte {
	return []byte(kotlinImportQueryPattern)
}

// kotlinImportQueryPattern is the Tree-sitter query for extracting Kotlin import statements.
const kotlinImportQueryPattern = `
; Import statements
(import_header) @import_path
`

// kotlinQueryPattern is the Tree-sitter query for extracting Kotlin signatures.
//
// Tree-sitter query syntax:
// - (@name) captures the identifier name
// - (@signature) captures the full declaration text
// - (@doc) captures the comment/documentation
// - (@kind) captures the node type for kind mapping
//
// Note: In tree-sitter-kotlin, interface declarations are represented as
// class_declaration nodes with an "interface" keyword. The refineKotlinClassKind()
// function in parser.go distinguishes between class, interface, and enum class.
const kotlinQueryPattern = `
; Function declarations (regular, suspend, inline, extension, operator, infix, tailrec)
(function_declaration
  (simple_identifier) @name
) @signature @kind

; Class declarations (class, data class, sealed class, enum class, interface, annotation class, value class)
(class_declaration
  (type_identifier) @name
) @signature @kind

; Object declarations (singleton)
(object_declaration
  (type_identifier) @name
) @signature @kind

; Companion object with explicit name (e.g., companion object Factory)
(companion_object
  (type_identifier) @name
) @signature @kind

; Property declarations (val/var, const val, lateinit, delegated)
(property_declaration
  (variable_declaration
    (simple_identifier) @name
  )
) @signature @kind

; Type alias
(type_alias
  (type_identifier) @name
) @signature @kind

; Enum entries
(enum_entry
  (simple_identifier) @name
) @signature @kind

; Secondary constructors
(secondary_constructor) @signature @kind

; Line comments
(line_comment) @doc

; Block/multiline comments
(multiline_comment) @doc
`
