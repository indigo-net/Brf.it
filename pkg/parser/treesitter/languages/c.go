// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"
)

// CQuery implements LanguageQuery for C language.
type CQuery struct {
	language *sitter.Language
	query    []byte
}

// NewCQuery creates a new C language query.
func NewCQuery() *CQuery {
	return &CQuery{
		language: sitter.NewLanguage(tree_sitter_c.Language()),
		query:    []byte(cQueryPattern),
	}
}

// Language returns the C Tree-sitter language.
func (q *CQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the C query pattern.
func (q *CQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for C queries.
func (q *CQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *CQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_definition":  "function",
		"declaration":          "function", // function prototypes
		"struct_specifier":     "struct",
		"enum_specifier":       "enum",
		"type_definition":      "typedef",
		"preproc_function_def": "macro",
		"preproc_def":          "macro",
		"global_variable":      "variable", // mapped from declaration patterns
	}
}

// ImportQuery returns the C import query pattern.
func (q *CQuery) ImportQuery() []byte {
	return []byte(cImportQueryPattern)
}

// cImportQueryPattern is the Tree-sitter query for extracting C #include directives.
const cImportQueryPattern = `
; #include "header.h"
(preproc_include
  path: (string_literal) @import_path
)

; #include <header.h>
(preproc_include
  path: (system_lib_string) @import_path
)
`

// cQueryPattern is the Tree-sitter query for extracting C signatures.
const cQueryPattern = `
; Function definitions - direct declarator
(function_definition
  declarator: (function_declarator
    declarator: (identifier) @name
  )
) @signature @kind

; Function definitions - pointer return type
(function_definition
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Function declarations (prototypes) - direct declarator
(declaration
  declarator: (function_declarator
    declarator: (identifier) @name
  )
) @signature @kind

; Function declarations (prototypes) - pointer return type
(declaration
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Struct specifiers
(struct_specifier
  name: (type_identifier) @name
) @signature @kind

; Enum specifiers
(enum_specifier
  name: (type_identifier) @name
) @signature @kind

; Typedef
(type_definition
  declarator: (type_identifier) @name
) @signature @kind

; Function-like macros
(preproc_function_def
  name: (identifier) @name
) @signature @kind

; Object-like macros
(preproc_def
  name: (identifier) @name
) @signature @kind

; Global variable declarations (with initializer)
(translation_unit
  (declaration
    declarator: (init_declarator
      declarator: (identifier) @name
    )
  ) @signature @kind
)

; Global variable declarations (simple identifier, e.g., extern)
(translation_unit
  (declaration
    declarator: (identifier) @name
  ) @signature @kind
)

; Global pointer variable declarations
(translation_unit
  (declaration
    declarator: (pointer_declarator
      declarator: (identifier) @name
    )
  ) @signature @kind
)

; Global pointer variable declarations (with initializer)
(translation_unit
  (declaration
    declarator: (init_declarator
      declarator: (pointer_declarator
        declarator: (identifier) @name
      )
    )
  ) @signature @kind
)

; Comments
(comment) @doc
`
