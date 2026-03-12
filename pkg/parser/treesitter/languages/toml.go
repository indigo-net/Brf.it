// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	tree_sitter_toml "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/toml"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// TOMLQuery implements LanguageQuery for TOML language.
type TOMLQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}

// NewTOMLQuery creates a new TOML language query.
func NewTOMLQuery() *TOMLQuery {
	return &TOMLQuery{
		language: sitter.NewLanguage(tree_sitter_toml.Language()),
		query:    []byte(tomlQueryPattern),
	}
}

// Language returns the TOML Tree-sitter language.
func (q *TOMLQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the TOML query pattern.
func (q *TOMLQuery) Query() []byte {
	return q.query
}

var tomlKindMapping = map[string]string{
	"table":               "namespace",
	"table_array_element": "namespace",
	"pair":                "variable",
}

// KindMapping returns the mapping from TOML node types to Signature kinds.
func (q *TOMLQuery) KindMapping() map[string]string {
	return tomlKindMapping
}

// ImportQuery returns nil since TOML has no import system.
func (q *TOMLQuery) ImportQuery() []byte {
	return nil
}

// tomlQueryPattern is the Tree-sitter query for extracting TOML signatures.
//
// TOML has three main structures:
// - [table] sections (standard tables)
// - [[table_array]] sections (array of tables)
// - key = value pairs
//
// Tables use bare_key or quoted_key for names.
// Pairs use bare_key, quoted_key, or dotted_key.
const tomlQueryPattern = `
; Standard table sections [section]
(table
  (bare_key) @name) @signature @kind

; Standard table sections with quoted key ["section"]
(table
  (quoted_key) @name) @signature @kind

; Standard table sections with dotted key [section.subsection]
(table
  (dotted_key) @name) @signature @kind

; Array of tables [[section]]
(table_array_element
  (bare_key) @name) @signature @kind

; Array of tables with quoted key [["section"]]
(table_array_element
  (quoted_key) @name) @signature @kind

; Array of tables with dotted key [[section.subsection]]
(table_array_element
  (dotted_key) @name) @signature @kind

; Key-value pairs (bare key)
(pair
  (bare_key) @name) @signature @kind

; Key-value pairs (quoted key)
(pair
  (quoted_key) @name) @signature @kind

; Key-value pairs (dotted key)
(pair
  (dotted_key) @name) @signature @kind

; TOML comments
(comment) @doc
`
