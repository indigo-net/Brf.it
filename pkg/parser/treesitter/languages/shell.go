// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_bash "github.com/tree-sitter/tree-sitter-bash/bindings/go"
)

// ShellQuery implements LanguageQuery for Shell/Bash language.
type ShellQuery struct {
	language *sitter.Language
	query    []byte
}

// NewShellQuery creates a new Shell/Bash language query.
func NewShellQuery() *ShellQuery {
	return &ShellQuery{
		language: sitter.NewLanguage(tree_sitter_bash.Language()),
		query:    []byte(shellQueryPattern),
	}
}

// Language returns the Shell/Bash Tree-sitter language.
func (q *ShellQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the Shell/Bash query pattern.
func (q *ShellQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for Shell/Bash queries.
func (q *ShellQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from node types to Signature kinds.
func (q *ShellQuery) KindMapping() map[string]string {
	return map[string]string{
		"function_definition": "function",
		"variable_assignment": "variable",
	}
}

// ImportQuery returns the Shell/Bash import query pattern.
func (q *ShellQuery) ImportQuery() []byte {
	return []byte(shellImportQueryPattern)
}

// shellImportQueryPattern is the Tree-sitter query for extracting Shell source/include statements.
// Captures source commands: source /path/to/file or . /path/to/file
const shellImportQueryPattern = `
; source /path/to/file and . /path/to/file
; Capture command nodes. Go-side filtering will check if command name is "source" or "."
(command
  name: (command_name) @name
) @import_path
`

// shellQueryPattern is the Tree-sitter query for extracting Shell/Bash signatures.
//
// Tree-sitter query syntax:
// - (@name) captures the identifier name
// - (@signature) captures the full declaration text
// - (@doc) captures the comment/documentation
// - (@kind) captures the node type for kind mapping
//
// Bash function definitions can be:
// - function foo { ... }
// - function foo() { ... }
// - foo() { ... }
const shellQueryPattern = `
; Function definitions: function foo { } or function foo() { } or foo() { }
(function_definition
  name: (word) @name
) @signature @kind

; Variable assignments: FOO=bar, FOO="bar", FOO=$(cmd)
(variable_assignment
  name: (variable_name) @name
) @signature @kind

; Comments (# ...)
(comment) @doc
`
