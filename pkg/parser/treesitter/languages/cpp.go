// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
)

// CppQuery implements LanguageQuery for C++ language.
type CppQuery struct {
	language *sitter.Language
	query    []byte
}

// NewCppQuery creates a new C++ language query.
func NewCppQuery() *CppQuery {
	return &CppQuery{
		language: sitter.NewLanguage(tree_sitter_cpp.Language()),
		query:    []byte(cppQueryPattern),
	}
}

// Language returns the C++ Tree-sitter language.
func (q *CppQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the C++ query pattern.
func (q *CppQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for C++ queries.
func (q *CppQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *CppQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_definition":      "function",
		"declaration":              "function",
		"struct_specifier":         "struct",
		"enum_specifier":           "enum",
		"type_definition":          "typedef",
		"preproc_function_def":     "macro",
		"preproc_def":              "macro",
		"class_specifier":          "class",
		"field_declaration":        "method",
		"template_declaration":     "template",
		"namespace_definition":     "namespace",
	}
}

// ImportQuery returns the C++ import query pattern.
func (q *CppQuery) ImportQuery() []byte {
	return []byte(cppImportQueryPattern)
}

// cppImportQueryPattern is the Tree-sitter query for extracting C++ #include directives.
const cppImportQueryPattern = `
; #include directives (capture full statement)
(preproc_include) @import_path
`

// cppQueryPattern is the Tree-sitter query for extracting C++ signatures.
const cppQueryPattern = `
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

; Function definitions - reference return type
(function_definition
  declarator: (reference_declarator
    (function_declarator
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

; Class definitions
(class_specifier
  name: (type_identifier) @name
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

; Method declarations in class (regular methods)
(field_declaration
  declarator: (function_declarator
    declarator: (field_identifier) @name
  )
) @signature @kind

; Method declarations with pointer return type
(field_declaration
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (field_identifier) @name
    )
  )
) @signature @kind

; Method declarations with reference return type
(field_declaration
  declarator: (reference_declarator
    (function_declarator
      declarator: (field_identifier) @name
    )
  )
) @signature @kind

; Constructor declarations (in class body)
(function_definition
  declarator: (function_declarator
    declarator: (qualified_identifier
      name: (identifier) @name
    )
  )
) @signature @kind

; Destructor definitions (outside class)
(function_definition
  declarator: (function_declarator
    declarator: (destructor_name
      (identifier) @name
    )
  )
) @signature @kind

; Destructor declarations in class (captured via declaration node)
(declaration
  declarator: (function_declarator
    declarator: (destructor_name
      (identifier) @name
    )
  )
) @signature @kind

; Namespace definitions
(namespace_definition
  name: (namespace_identifier) @name
) @signature @kind

; Template function definitions
(template_declaration
  (function_definition
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Template function definitions - pointer return type
(template_declaration
  (function_definition
    declarator: (pointer_declarator
      declarator: (function_declarator
        declarator: (identifier) @name
      )
    )
  )
) @signature @kind

; Template class definitions
(template_declaration
  (class_specifier
    name: (type_identifier) @name
  )
) @signature @kind

; Template struct definitions
(template_declaration
  (struct_specifier
    name: (type_identifier) @name
  )
) @signature @kind

; Template declarations (standalone)
(template_declaration
  (declaration
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Comments
(comment) @doc
`
