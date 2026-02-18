// Package parser provides code parsing capabilities for brfit.
package parser

import "sync"

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

// ParseResult contains the result of parsing a single file.
type ParseResult struct {
	// FilePath is the path to the parsed file.
	FilePath string

	// Language is the detected language.
	Language string

	// Signatures is the list of extracted signatures.
	Signatures []Signature

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
}

// Parser defines the interface for code parsers.
type Parser interface {
	// Parse parses the given content and returns extracted signatures.
	Parse(content string, opts *Options) (*ParseResult, error)

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
