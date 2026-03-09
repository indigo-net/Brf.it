// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_php "github.com/tree-sitter/tree-sitter-php/bindings/go"
)

// PHPQuery implements LanguageQuery for PHP language.
type PHPQuery struct {
	language *sitter.Language
	query    []byte
}

// NewPHPQuery creates a new PHP language query.
func NewPHPQuery() *PHPQuery {
	return &PHPQuery{
		language: sitter.NewLanguage(tree_sitter_php.LanguagePHP()),
		query:    []byte(phpQueryPattern),
	}
}

// Language returns the PHP Tree-sitter language.
func (q *PHPQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the PHP query pattern.
func (q *PHPQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for PHP queries.
func (q *PHPQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *PHPQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_definition":       "function",
		"method_declaration":        "method",
		"class_declaration":         "class",
		"interface_declaration":     "interface",
		"trait_declaration":         "type",
		"enum_declaration":          "enum",
		"const_declaration":         "variable",
		"property_declaration":      "variable",
		"namespace_use_declaration": "import",
	}
}

// ImportQuery returns the PHP import query pattern.
func (q *PHPQuery) ImportQuery() []byte {
	return []byte(phpImportQueryPattern)
}

// phpImportQueryPattern is the Tree-sitter query for extracting PHP use/include statements.
const phpImportQueryPattern = `
; use Namespace\\Class;
(namespace_use_declaration) @import_path

; include 'file.php';
(include_expression) @import_path

; require 'vendor/autoload.php';
(require_expression) @import_path

; include_once 'config.php';
(include_once_expression) @import_path

; require_once 'config.php';
(require_once_expression) @import_path
`

// phpQueryPattern is the Tree-sitter query for extracting PHP signatures.
const phpQueryPattern = `
; Function definitions: function name() {}
(function_definition name: (name) @name) @signature @kind

; Method declarations in classes
(method_declaration name: (name) @name) @signature @kind

; Class declarations
(class_declaration name: (name) @name) @signature @kind

; Interface declarations
(interface_declaration name: (name) @name) @signature @kind

; Trait declarations
(trait_declaration name: (name) @name) @signature @kind

; Enum declarations
(enum_declaration name: (name) @name) @signature @kind

; Const declarations: const NAME = value;
(const_declaration
  (const_element
    (name) @name
  )
) @signature @kind

; Property declarations: public $name;
(property_declaration
  (property_element
    (variable_name (name) @name)
  )
) @signature @kind

; Comments (PHPDoc and regular)
(comment) @doc
`
