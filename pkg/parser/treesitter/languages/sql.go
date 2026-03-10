// Package languages provides language-specific Tree-sitter queries.
package languages

import (
	tree_sitter_sql "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/sql"
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// SQLQuery implements LanguageQuery for SQL language.
type SQLQuery struct {
	language *sitter.Language
	query    []byte
}

// NewSQLQuery creates a new SQL language query.
func NewSQLQuery() *SQLQuery {
	return &SQLQuery{
		language: sitter.NewLanguage(tree_sitter_sql.Language()),
		query:    []byte(sqlQueryPattern),
	}
}

// Language returns the SQL Tree-sitter language.
func (q *SQLQuery) Language() *sitter.Language {
	return q.language
}

// Query returns the SQL query pattern.
func (q *SQLQuery) Query() []byte {
	return q.query
}

// Captures returns the capture names for SQL queries.
func (q *SQLQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// KindMapping returns the mapping from SQL DDL node types to Signature kinds.
func (q *SQLQuery) KindMapping() map[string]string {
	return map[string]string{
		"create_table":             "struct",
		"create_function":          "function",
		"create_view":              "type",
		"create_materialized_view": "type",
		"create_index":             "variable",
		"create_trigger":           "function",
		"create_type":              "type",
		"create_schema":            "namespace",
		"create_sequence":          "variable",
		"alter_table":              "type",
	}
}

// ImportQuery returns nil since SQL has no import system.
func (q *SQLQuery) ImportQuery() []byte {
	return nil
}

// sqlQueryPattern is the Tree-sitter query for extracting SQL DDL signatures.
//
// SQL DDL statements are captured as whole nodes. The object name is extracted
// from object_reference (most DDL) or via Go-side text parsing (CREATE INDEX,
// CREATE SCHEMA) where the name is a bare identifier.
//
// For CREATE TRIGGER, two object_references exist (trigger name + table name);
// the first match wins via seenLines deduplication in parser.go.
const sqlQueryPattern = `
; CREATE TABLE
(create_table
  (object_reference) @name) @signature @kind

; CREATE FUNCTION
(create_function
  (object_reference) @name) @signature @kind

; CREATE VIEW
(create_view
  (object_reference) @name) @signature @kind

; CREATE MATERIALIZED VIEW
(create_materialized_view
  (object_reference) @name) @signature @kind

; CREATE INDEX (name extracted Go-side)
(create_index) @signature @kind

; CREATE TYPE
(create_type
  (object_reference) @name) @signature @kind

; CREATE TRIGGER (first object_reference = trigger name)
(create_trigger
  (object_reference) @name) @signature @kind

; CREATE SCHEMA (bare identifier)
(create_schema
  (identifier) @name) @signature @kind

; CREATE SEQUENCE
(create_sequence
  (object_reference) @name) @signature @kind

; ALTER TABLE
(alter_table
  (object_reference) @name) @signature @kind

; SQL comments (-- single-line)
(comment) @doc

; SQL multi-line comments (/* ... */ are "marginalia" in tree-sitter-sql)
(marginalia) @doc
`
