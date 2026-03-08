package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_bash "github.com/tree-sitter/tree-sitter-bash/bindings/go"
)

// BashQuery implements LanguageQuery for Bash/Shell language.
type BashQuery struct {
	language *sitter.Language
	query    []byte
}

// NewBashQuery creates a new Bash language query.
func NewBashQuery() *BashQuery {
	return &BashQuery{
		language: sitter.NewLanguage(tree_sitter_bash.Language()),
		query:    []byte(bashQueryPattern),
	}
}

// Language returns the Bash Tree-sitter language.
func (q *BashQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Bash query pattern.
func (q *BashQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for Bash queries.
func (q *BashQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *BashQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_definition":  "function",
		"variable_assignment":  "variable",
		"declaration_command":  "variable",
	}
}

// ImportQuery returns the Bash import query pattern.
func (q *BashQuery) ImportQuery() []byte {
	return []byte(bashImportQueryPattern)
}

// bashImportQueryPattern is the Tree-sitter query for extracting Bash source/include statements.
// Bash doesn't have a native import system, but `source` or `.` commands are used to include files.
const bashImportQueryPattern = `
; source /path/to/file.sh or . /path/to/file.sh
(command
  name: [
    (command_name (word) @_cmd)
    (string)
  ]
  argument: [
    (word) @import_path
    (string) @import_path
    (concatenation (word) @import_path)
  ]
  (#match? @_cmd "^(source|\\.)$")
)

; source "/path/to/file.sh"
(command
  name: (command_name (word) @_cmd)
  argument: (string) @import_path
  (#eq? @_cmd "source")
)
`

// bashQueryPattern is the Tree-sitter query for extracting Bash signatures.
//
// Tree-sitter query syntax:
// - (@name) captures the identifier name
// - (@signature) captures the full declaration text
// - (@doc) captures the comment/documentation
// - (@kind) captures the node type for kind mapping
//
// Bash function declarations:
// - function foo { ... }
// - function foo() { ... }
// - foo() { ... }
//
// Bash variables:
// - VAR=value
// - declare VAR
// - local VAR
// - readonly VAR
const bashQueryPattern = `
; Function declarations: function foo { ... } or foo() { ... }
; Both use (word) as the name child
(function_definition
  (word) @name
) @signature @kind

; Variable assignments: VAR=value
(variable_assignment
  name: (variable_name) @name
) @signature @kind

; Declaration commands: declare, local, readonly, typeset
(declaration_command
  (variable_name) @name
) @signature @kind

; Comments (including #)
(comment) @doc
`
