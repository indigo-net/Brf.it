// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	tree_sitter_c_sharp "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/csharp"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// CSharpQuery implements LanguageQuery for C# language.
type CSharpQuery struct {
	language *sitter.Language
	query    []byte
}

// NewCSharpQuery creates a new C# language query.
func NewCSharpQuery() *CSharpQuery {
	return &CSharpQuery{
		language: sitter.NewLanguage(tree_sitter_c_sharp.Language()),
		query:    []byte(csharpQueryPattern),
	}
}

// Language returns the C# Tree-sitter language.
func (q *CSharpQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the C# query pattern.
func (q *CSharpQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for C# queries.
func (q *CSharpQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *CSharpQuery) KindMapping() map[string]string {
	return map[string]string{
		"class_declaration":                 "class",
		"struct_declaration":                "struct",
		"interface_declaration":             "interface",
		"enum_declaration":                  "enum",
		"record_declaration":                "record",
		"delegate_declaration":              "type",
		"method_declaration":                "method",
		"constructor_declaration":           "constructor",
		"destructor_declaration":            "destructor",
		"property_declaration":              "variable",
		"field_declaration":                 "field",
		"event_declaration":                 "variable",
		"event_field_declaration":           "variable",
		"indexer_declaration":               "method",
		"operator_declaration":              "function",
		"conversion_operator_declaration":   "function",
		"namespace_declaration":             "namespace",
		"file_scoped_namespace_declaration": "namespace",
		"enum_member_declaration":           "variable",
	}
}

// ImportQuery returns the C# import query pattern.
func (q *CSharpQuery) ImportQuery() []byte {
	return []byte(csharpImportQueryPattern)
}

// csharpImportQueryPattern is the Tree-sitter query for extracting C# using directives.
const csharpImportQueryPattern = `
; using directives (capture full declaration)
(using_directive) @import_path
`

// csharpQueryPattern is the Tree-sitter query for extracting C# signatures.
//
// Tree-sitter query syntax:
// - (@name) captures the identifier name
// - (@signature) captures the full declaration text
// - (@doc) captures the comment/documentation
// - (@kind) captures the node type for kind mapping
const csharpQueryPattern = `
; Class declarations
(class_declaration
  name: (identifier) @name
) @signature @kind

; Struct declarations
(struct_declaration
  name: (identifier) @name
) @signature @kind

; Interface declarations
(interface_declaration
  name: (identifier) @name
) @signature @kind

; Enum declarations
(enum_declaration
  name: (identifier) @name
) @signature @kind

; Record declarations (record, record class, record struct)
(record_declaration
  name: (identifier) @name
) @signature @kind

; Delegate declarations
(delegate_declaration
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

; Destructor declarations
(destructor_declaration
  name: (identifier) @name
) @signature @kind

; Property declarations
(property_declaration
  name: (identifier) @name
) @signature @kind

; Field declarations (static/const filtered in parser.go)
(field_declaration
  (variable_declaration
    (variable_declarator
      name: (identifier) @name
    )
  )
) @signature @kind

; Event field declarations (e.g., public event EventHandler Changed;)
(event_field_declaration
  (variable_declaration
    (variable_declarator
      name: (identifier) @name
    )
  )
) @signature @kind

; Event declarations with accessor body
(event_declaration
  name: (identifier) @name
) @signature @kind

; Indexer declarations (no name capture — synthesized in parser.go)
(indexer_declaration) @signature @kind

; Operator declarations (no name capture — synthesized in parser.go)
(operator_declaration) @signature @kind

; Conversion operator declarations (no name capture — synthesized in parser.go)
(conversion_operator_declaration) @signature @kind

; Namespace declarations
(namespace_declaration
  name: (_) @name
) @signature @kind

; File-scoped namespace declarations (C# 10+)
(file_scoped_namespace_declaration
  name: (_) @name
) @signature @kind

; Enum member declarations
(enum_member_declaration
  name: (identifier) @name
) @signature @kind

; Comments (XML doc comments and regular)
(comment) @doc
`
