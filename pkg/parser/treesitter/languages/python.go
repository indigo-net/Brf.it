package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

// PythonQuery implements LanguageQuery for Python language.
type PythonQuery struct {
	language *sitter.Language
	query    []byte
}

// NewPythonQuery creates a new Python language query.
func NewPythonQuery() *PythonQuery {
	return &PythonQuery{
		language: sitter.NewLanguage(tree_sitter_python.Language()),
		query:    []byte(pythonQueryPattern),
	}
}

// Language returns the Python Tree-sitter language.
func (q *PythonQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Python query pattern.
func (q *PythonQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for Python queries.
func (q *PythonQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *PythonQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_definition":  "function",
		"class_definition":     "class",
		"expression_statement": "variable",
		"assignment":           "variable",
	}
}

// ImportQuery returns the Python import query pattern.
func (q *PythonQuery) ImportQuery() []byte {
	return []byte(pythonImportQueryPattern)
}

// pythonImportQueryPattern is the Tree-sitter query for extracting Python imports.
const pythonImportQueryPattern = `
; import module
(import_statement
  name: (dotted_name) @import_path
)

; from module import ...
(import_from_statement
  module_name: (dotted_name) @import_path
)

; from . import ... (relative imports)
(import_from_statement
  module_name: (relative_import) @import_path
)
`

// pythonQueryPattern is the Tree-sitter query for extracting Python signatures.
const pythonQueryPattern = `
; Function definitions (includes async def, methods)
(function_definition
  name: (identifier) @name
) @signature @kind

; Class definitions
(class_definition
  name: (identifier) @name
) @signature @kind

; Module-level assignments (simple and with type annotations)
(module
  (expression_statement
    (assignment
      left: (identifier) @name
    )
  ) @signature @kind
)

; Comments
(comment) @doc
`
