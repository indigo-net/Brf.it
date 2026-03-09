// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_scala "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/scala"
)

// ScalaQuery implements LanguageQuery for Scala language.
type ScalaQuery struct {
	language *sitter.Language
	query    []byte
}

// NewScalaQuery creates a new Scala language query.
func NewScalaQuery() *ScalaQuery {
	return &ScalaQuery{
		language: sitter.NewLanguage(tree_sitter_scala.Language()),
		query:    []byte(scalaQueryPattern),
	}
}

// Language returns the Scala Tree-sitter language.
func (q *ScalaQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Scala query pattern.
func (q *ScalaQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for Scala queries.
func (q *ScalaQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *ScalaQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_definition":  "method",
		"function_declaration": "method",
		"class_definition":     "class",
		"trait_definition":     "trait",
		"object_definition":    "class",
		"val_definition":       "variable",
		"val_declaration":      "variable",
		"var_definition":       "variable",
		"var_declaration":      "variable",
		"type_definition":      "type",
		"enum_definition":      "enum",
		"given_definition":     "variable",
		"extension_definition": "method",
	}
}

// ImportQuery returns the Scala import query pattern.
func (q *ScalaQuery) ImportQuery() []byte {
	return []byte(scalaImportQueryPattern)
}

// scalaImportQueryPattern is the Tree-sitter query for extracting Scala import statements.
const scalaImportQueryPattern = `
; Import statements
(import_declaration) @import_path
`

// scalaQueryPattern is the Tree-sitter query for extracting Scala signatures.
//
// Tree-sitter query syntax:
// - (@name) captures the identifier name
// - (@signature) captures the full declaration text
// - (@doc) captures the comment/documentation
// - (@kind) captures the node type for kind mapping
//
// Scala declaration types captured:
// - def (function definitions and declarations)
// - class (regular, abstract, case, sealed, implicit)
// - trait (regular, sealed)
// - object (standalone, companion)
// - val/var (including lazy val, implicit val)
// - type aliases
// - enum (Scala 3)
// - given (Scala 3)
// - extension (Scala 3)
const scalaQueryPattern = `
; Function definitions (def with body)
(function_definition
  name: (identifier) @name
) @signature @kind

; Function declarations (abstract methods in traits/classes, no body)
(function_declaration
  name: (identifier) @name
) @signature @kind

; Class definitions (class, abstract class, case class, sealed class, implicit class)
(class_definition
  name: (identifier) @name
) @signature @kind

; Trait definitions (trait, sealed trait)
(trait_definition
  name: (identifier) @name
) @signature @kind

; Object definitions (singleton, companion)
(object_definition
  name: (identifier) @name
) @signature @kind

; Val definitions (val, lazy val, implicit val)
(val_definition
  pattern: (identifier) @name
) @signature @kind

; Val declarations (abstract val in traits)
(val_declaration
  name: (identifier) @name
) @signature @kind

; Var definitions
(var_definition
  pattern: (identifier) @name
) @signature @kind

; Var declarations (abstract var in traits)
(var_declaration
  name: (identifier) @name
) @signature @kind

; Type aliases
(type_definition
  name: (type_identifier) @name
) @signature @kind

; Enum definitions (Scala 3)
(enum_definition
  name: (identifier) @name
) @signature @kind

; Given definitions (Scala 3)
(given_definition
  name: (identifier) @name
) @signature @kind

; Extension definitions (Scala 3)
(extension_definition) @signature @kind

; Line comments
(comment) @doc
`
