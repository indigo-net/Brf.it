// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	tree_sitter_yaml "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/yaml"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// YAMLQuery implements LanguageQuery for YAML language.
type YAMLQuery struct {
	BaseQuery
	language *sitter.Language
	query    []byte
}

// NewYAMLQuery creates a new YAML language query.
func NewYAMLQuery() *YAMLQuery {
	return &YAMLQuery{
		language: sitter.NewLanguage(tree_sitter_yaml.Language()),
		query:    []byte(yamlQueryPattern),
	}
}

// Language returns the YAML Tree-sitter language.
func (q *YAMLQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the YAML query pattern.
func (q *YAMLQuery) Query() []byte {
	return q.query
}

var yamlKindMapping = map[string]string{
	"block_mapping_pair": "variable",
}

// KindMapping returns the mapping from YAML node types to Signature kinds.
func (q *YAMLQuery) KindMapping() map[string]string {
	return yamlKindMapping
}

// ImportQuery returns nil since YAML has no import system.
func (q *YAMLQuery) ImportQuery() []byte {
	return nil
}

// yamlQueryPattern is the Tree-sitter query for extracting YAML key-value signatures.
//
// YAML structures are captured as block_mapping_pair nodes. The key field
// provides the name, and the entire pair is the signature.
// Only top-level keys are captured to avoid excessive noise from nested values.
const yamlQueryPattern = `
; Top-level key-value pairs
(block_mapping_pair
  key: (_) @name) @signature @kind

; YAML comments
(comment) @doc
`
