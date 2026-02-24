package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

// TypeScriptQuery implements LanguageQuery for TypeScript language.
type TypeScriptQuery struct {
	language *sitter.Language
	query    []byte
}

// NewTypeScriptQuery creates a new TypeScript language query.
func NewTypeScriptQuery() *TypeScriptQuery {
	return &TypeScriptQuery{
		language: sitter.NewLanguage(tree_sitter_typescript.LanguageTypescript()),
		query:    []byte(typeScriptQueryPattern),
	}
}

// Language returns the TypeScript Tree-sitter language.
func (q *TypeScriptQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the TypeScript query pattern.
func (q *TypeScriptQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for TypeScript queries.
func (q *TypeScriptQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *TypeScriptQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_declaration":   "function",
		"method_definition":      "method",
		"class_declaration":      "class",
		"interface_declaration":  "interface",
		"type_alias_declaration": "type",
		"arrow_function":         "function",
		"variable_declaration":   "variable",
		"variable_declarator":    "arrow",
		"lexical_declaration":    "arrow",
		"export_statement":       "export",
	}
}

// typeScriptQueryPattern is the Tree-sitter query for extracting TypeScript signatures.
const typeScriptQueryPattern = `
; Function declarations
(function_declaration
  name: (identifier) @name
) @signature @kind

; Exported function declarations
(export_statement
  (function_declaration
    name: (identifier) @name
  )
) @signature @kind

; Arrow functions in variable declarations (capture full declaration with const/let/var)
(lexical_declaration
  (variable_declarator
    name: (identifier) @name
    value: (arrow_function)
  )
) @signature @kind

; Method definitions
(method_definition
  name: (property_identifier) @name
) @signature @kind

; Class declarations
(class_declaration
  name: (type_identifier) @name
) @signature @kind

; Interface declarations
(interface_declaration
  name: (type_identifier) @name
) @signature @kind

; Type alias declarations
(type_alias_declaration
  name: (type_identifier) @name
) @signature @kind

; Comments (documentation)
(comment) @doc
`
