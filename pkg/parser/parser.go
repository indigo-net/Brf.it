// Package parser provides code parsing capabilities for brfit.
package parser

import (
	"path/filepath"
	"strings"
	"sync"
)

// Signature represents an extracted code signature (function, class, method, etc.).
type Signature struct {
	// Name is the identifier name (e.g., "Scan", "FileScanner").
	Name string

	// Kind is the type of signature (e.g., "function", "method", "class", "interface").
	Kind string

	// Text is the full signature text including parameters and return type.
	Text string

	// Doc is the documentation comment (if any).
	Doc string

	// Line is the starting line number (1-indexed).
	Line int

	// EndLine is the ending line number (1-indexed).
	EndLine int

	// Language is the source language (e.g., "go", "typescript").
	Language string

	// Exported indicates whether the signature is exported/public.
	Exported bool
}

// Node represents a node in the parsed AST.
type Node struct {
	// Type is the node type (e.g., "function_declaration", "class_definition").
	Type string

	// StartRow is the starting row (0-indexed).
	StartRow int

	// EndRow is the ending row (0-indexed).
	EndRow int

	// StartColumn is the starting column.
	StartColumn int

	// EndColumn is the ending column.
	EndColumn int

	// Text is the source text of the node.
	Text string

	// Children are child nodes.
	Children []Node
}

// FunctionCall represents a function/method call reference within a file.
type FunctionCall struct {
	// Caller is the name of the enclosing function (empty if top-level).
	Caller string

	// Callee is the called function/method name.
	Callee string

	// Line is the line number where the call occurs (1-indexed).
	Line int
}

// ParseResult contains the result of parsing a single file.
type ParseResult struct {
	// FilePath is the path to the parsed file.
	FilePath string

	// Language is the detected language.
	Language string

	// Signatures is the list of extracted signatures.
	Signatures []Signature

	// RawImports is the list of raw import/export statement text.
	RawImports []string

	// Calls is the list of function call references.
	Calls []FunctionCall

	// AST is the root node of the parsed AST (optional).
	AST *Node

	// Error is any error that occurred during parsing.
	Error error
}

// Options configures the parsing behavior.
type Options struct {
	// Language forces a specific language (auto-detected if empty).
	Language string

	// IncludeAST whether to include the full AST in the result.
	IncludeAST bool

	// IncludePrivate whether to include non-exported/private signatures.
	IncludePrivate bool

	// IncludeBody whether to include function/method bodies in the signature text.
	// When false (default), only the signature line is extracted.
	// When true, the full declaration including the body is extracted.
	IncludeBody bool

	// IncludeImports whether to include import/export statements in the result.
	IncludeImports bool

	// IncludeCalls whether to include function call references in the result.
	IncludeCalls bool
}

// Parser defines the interface for code parsers.
type Parser interface {
	// Parse parses the given content and returns extracted signatures.
	// Content is passed as []byte to avoid unnecessary string conversion
	// from os.ReadFile output.
	Parse(content []byte, opts *Options) (*ParseResult, error)

	// Languages returns the list of supported languages.
	Languages() []string
}

// Registry manages available parsers.
type Registry struct {
	mu      sync.RWMutex
	parsers map[string]Parser
}

// NewRegistry creates a new empty parser registry.
func NewRegistry() *Registry {
	return &Registry{
		parsers: make(map[string]Parser),
	}
}

// defaultRegistry is the global parser registry.
var defaultRegistry = NewRegistry()

// DefaultRegistry returns the global parser registry.
func DefaultRegistry() *Registry {
	return defaultRegistry
}

// Register adds a parser for the given language.
func (r *Registry) Register(lang string, parser Parser) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.parsers[lang] = parser
}

// Get returns the parser for the given language.
func (r *Registry) Get(lang string) (Parser, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	parser, ok := r.parsers[lang]
	return parser, ok
}

// Languages returns all registered languages.
func (r *Registry) Languages() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	langs := make([]string, 0, len(r.parsers))
	for lang := range r.parsers {
		langs = append(langs, lang)
	}
	return langs
}

// RegisterParser registers a parser in the default registry.
func RegisterParser(lang string, parser Parser) {
	defaultRegistry.Register(lang, parser)
}

// GetParser returns a parser from the default registry.
func GetParser(lang string) (Parser, bool) {
	return defaultRegistry.Get(lang)
}

// languageMapping maps file extensions to language names.
// This is the canonical source of truth for extension-to-language mapping.
// Immutable after package initialization; safe for concurrent reads.
var languageMapping = map[string]string{
	".go":    "go",
	".ts":    "typescript",
	".tsx":   "typescript",
	".js":    "javascript",
	".jsx":   "javascript",
	".py":    "python",
	".java":  "java",
	".rs":    "rust",
	".rb":    "ruby",
	".php":   "php",
	".c":     "c",
	".cpp":   "cpp",
	".h":     "cpp",
	".hpp":   "cpp",
	".cs":    "csharp",
	".swift": "swift",
	".kt":    "kotlin",
	".kts":   "kotlin",
	".lua":   "lua",
	".sh":    "shell",
	".bash":  "shell",
	".zsh":   "shell",
	".scala": "scala",
	".sc":    "scala",
	".ex":    "elixir",
	".exs":   "elixir",
	".sql":   "sql",
	".yaml":  "yaml",
	".yml":   "yaml",
	".toml":  "toml",
}

// LanguageMapping returns a copy of the canonical extension-to-language mapping.
func LanguageMapping() map[string]string {
	m := make(map[string]string, len(languageMapping))
	for k, v := range languageMapping {
		m[k] = v
	}
	return m
}

// DetectLanguage returns the language for a given file path.
func DetectLanguage(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if lang, ok := languageMapping[ext]; ok {
		return lang
	}
	return ""
}
