// Package treesitter provides Tree-sitter based parser implementations.
package treesitter

import (
	sitter "github.com/tree-sitter/go-tree-sitter"
)

// LanguageQuery defines the interface for language-specific Tree-sitter queries.
type LanguageQuery interface {
	// Language returns the Tree-sitter language for parsing.
	Language() *sitter.Language

	// Query returns the Tree-sitter query pattern for signature extraction.
	Query() []byte

	// ImportQuery returns the Tree-sitter query pattern for import/export extraction.
	// Returns nil if the language doesn't support import extraction.
	ImportQuery() []byte

	// Captures returns the list of capture names used in the query.
	Captures() []string

	// KindMapping maps Tree-sitter node types to Signature kinds.
	KindMapping() map[string]string
}

// Capture names used across all language queries.
const (
	CaptureName      = "name"
	CaptureSignature = "signature"
	CaptureDoc       = "doc"
	CaptureKind      = "kind"
)

// Capture names for import queries.
const (
	CaptureImportPath = "import_path"
	CaptureExportName = "export_name"
	CaptureImportType = "import_type"
)

// DefaultKindMapping provides default kind mappings (can be overridden per language).
var DefaultKindMapping = map[string]string{
	"function_declaration": "function",
	"method_declaration":   "method",
	"type_declaration":     "type",
	"struct_type":          "struct",
	"interface_type":       "interface",
	"class_declaration":    "class",
	"arrow_function":       "function",
	"function_expression":  "function",
	"method_definition":    "method",
}
