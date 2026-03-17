package treesitter

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/indigo-net/Brf.it/pkg/parser"
	"github.com/indigo-net/Brf.it/pkg/parser/treesitter/languages"
)

// init registers a single shared TreeSitterParser with the default registry.
// One instance handles all languages via its internal query map, sharing
// sync.Pool resources (parsers, cursors) across languages for better reuse.
func init() {
	p := NewTreeSitterParser()
	for lang := range p.queries {
		parser.RegisterParser(lang, p)
	}
}

// queryType distinguishes between signature and import queries for caching.
type queryType int

const (
	queryTypeSignature queryType = iota
	queryTypeImport
	queryTypeCall
)

// elixirDefPrefixes are byte-slice prefixes used to filter out non-import
// Elixir definitions (defmodule/defprotocol/defimpl) from import queries.
var elixirDefPrefixes = [][]byte{
	[]byte("defmodule "),
	[]byte("defprotocol "),
	[]byte("defimpl "),
}

// supportedLangs is the list of languages with Tree-sitter parsers available.
const supportedLangs = "go, typescript, tsx, javascript, jsx, python, c, java, cpp, rust, swift, kotlin, csharp, lua, shell, php, ruby, scala, elixir, sql, yaml, toml"

// queryCacheKey combines language and query type for cache lookup.
type queryCacheKey struct {
	lang string
	typ  queryType
}

// TreeSitterParser implements parser.Parser using Tree-sitter.
type TreeSitterParser struct {
	queries         map[string]LanguageQuery
	compiledQueries sync.Map // map[queryCacheKey]*sitter.Query
	parserPool      sync.Pool
	cursorPool      sync.Pool
	mu              sync.RWMutex // guards query lifetime around Close
	closed          bool
}

// pooledParser wraps a sitter.Parser with the last language set,
// allowing SetLanguage to be skipped when the same language is reused.
type pooledParser struct {
	parser   *sitter.Parser
	lastLang *sitter.Language
}

// NewTreeSitterParser creates a new Tree-sitter based parser.
func NewTreeSitterParser() *TreeSitterParser {
	p := &TreeSitterParser{
		queries: map[string]LanguageQuery{
			"go":         languages.NewGoQuery(),
			"typescript": languages.NewTypeScriptQuery(),
			"tsx":        languages.NewTypeScriptQuery(), // TSX uses TypeScript grammar
			"javascript": languages.NewTypeScriptQuery(), // JS uses TypeScript grammar (subset)
			"jsx":        languages.NewTypeScriptQuery(), // JSX uses TypeScript grammar
			"python":     languages.NewPythonQuery(),
			"c":          languages.NewCQuery(),
			"java":       languages.NewJavaQuery(),
			"cpp":        languages.NewCppQuery(),
			"rust":       languages.NewRustQuery(),
			"swift":      languages.NewSwiftQuery(),
			"kotlin":     languages.NewKotlinQuery(),
			"csharp":     languages.NewCSharpQuery(),
			"lua":        languages.NewLuaQuery(),
			"shell":      languages.NewShellQuery(),
			"php":        languages.NewPHPQuery(),
			"ruby":       languages.NewRubyQuery(),
			"scala":      languages.NewScalaQuery(),
			"elixir":     languages.NewElixirQuery(),
			"sql":        languages.NewSQLQuery(),
			"yaml":       languages.NewYAMLQuery(),
			"toml":       languages.NewTOMLQuery(),
		},
	}
	p.parserPool = sync.Pool{
		New: func() any {
			return &pooledParser{parser: sitter.NewParser()}
		},
	}
	p.cursorPool = sync.Pool{
		New: func() any {
			return sitter.NewQueryCursor()
		},
	}
	return p
}

// Close releases all cached Tree-sitter Query objects.
// This should be called when the parser is no longer needed,
// especially in long-running processes like brfit-mcp.
// After Close, Parse returns an error instead of accessing freed memory.
// Close acquires a write lock to ensure no concurrent Parse is using queries.
func (p *TreeSitterParser) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	p.compiledQueries.Range(func(key, value any) bool {
		if q, ok := value.(*sitter.Query); ok {
			q.Close()
		}
		p.compiledQueries.Delete(key)
		return true
	})
}

// getOrCreateQuery returns a cached query or creates and caches a new one.
// The returned query should NOT be closed by the caller - it's managed by the cache.
// Callers must hold p.mu.RLock() to ensure queries are not freed by Close().
func (p *TreeSitterParser) getOrCreateQuery(lang string, langQuery LanguageQuery, typ queryType) (*sitter.Query, error) {
	key := queryCacheKey{lang: lang, typ: typ}

	// Fast path: check cache (sync.Map is inherently thread-safe)
	if cached, ok := p.compiledQueries.Load(key); ok {
		return cached.(*sitter.Query), nil
	}

	// Slow path: create query and attempt to store it
	var queryBytes []byte
	switch typ {
	case queryTypeSignature:
		queryBytes = langQuery.Query()
	case queryTypeImport:
		queryBytes = langQuery.ImportQuery()
	case queryTypeCall:
		queryBytes = langQuery.CallQuery()
	}
	queryStr := string(queryBytes)

	query, err := sitter.NewQuery(langQuery.Language(), queryStr)
	if err != nil {
		return nil, err
	}

	// LoadOrStore ensures only one query per key is retained.
	// If another goroutine won the race, close our duplicate and use theirs.
	actual, loaded := p.compiledQueries.LoadOrStore(key, query)
	if loaded {
		query.Close()
		return actual.(*sitter.Query), nil
	}
	return query, nil
}

// Parse parses the given content and returns extracted signatures.
// Parse holds a read lock to ensure Close() cannot free queries mid-parse.
func (p *TreeSitterParser) Parse(content []byte, opts *parser.Options) (result *parser.ParseResult, err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return nil, fmt.Errorf("parser is closed")
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("tree-sitter panic recovered: %v", r)
			result = nil
		}
	}()

	if opts == nil {
		opts = &parser.Options{}
	}

	// Determine language
	lang := opts.Language
	if lang == "" {
		return nil, fmt.Errorf("language must be specified in Options.Language (received empty string). Supported: %s", supportedLangs)
	}

	// Get language query
	query, ok := p.queries[lang]
	if !ok {
		return nil, fmt.Errorf("unsupported language %q. Available parsers: %s", lang, supportedLangs)
	}

	// Get parser from pool
	pp := p.parserPool.Get().(*pooledParser)
	defer p.parserPool.Put(pp)

	// Set language (skip if already set to the same language)
	tsLang := query.Language()
	if pp.lastLang != tsLang {
		if err := pp.parser.SetLanguage(tsLang); err != nil {
			return nil, fmt.Errorf("failed to initialize %q parser: %w (grammar may be corrupted or incompatible)", lang, err)
		}
		pp.lastLang = tsLang
	}

	// Parse content (no conversion needed - already []byte)
	tree := pp.parser.Parse(content, nil)
	if tree == nil {
		return nil, fmt.Errorf("failed to parse content for language %q (content may be malformed or parser error occurred)", lang)
	}
	defer tree.Close()

	// Extract signatures
	signatures, err := p.extractSignatures(tree.RootNode(), content, query, opts)
	if err != nil {
		return nil, fmt.Errorf("signature extraction failed: %w", err)
	}

	// Extract imports if requested
	var rawImports []string
	if opts.IncludeImports {
		rawImports, err = p.extractImports(tree.RootNode(), content, query, opts)
		if err != nil {
			return nil, fmt.Errorf("import extraction failed: %w", err)
		}
	}

	// Extract calls if requested (uses full signatures for caller attribution)
	var calls []parser.FunctionCall
	if opts.IncludeCalls && query.CallQuery() != nil {
		calls, err = p.extractCalls(tree.RootNode(), content, query, opts, signatures)
		if err != nil {
			return nil, fmt.Errorf("call extraction failed: %w", err)
		}
	}

	// Filter private symbols after call extraction (calls need all signatures for caller attribution)
	if !opts.IncludePrivate {
		filtered := signatures[:0]
		for _, sig := range signatures {
			if sig.Exported {
				filtered = append(filtered, sig)
			}
		}
		signatures = filtered
	}

	return &parser.ParseResult{
		Language:   lang,
		Signatures: signatures,
		RawImports: rawImports,
		Calls:      calls,
	}, nil
}

// Languages returns the list of supported languages.
func (p *TreeSitterParser) Languages() []string {
	langs := make([]string, 0, len(p.queries))
	for lang := range p.queries {
		langs = append(langs, lang)
	}
	return langs
}

// extractSignatures extracts signatures from the AST using the language query.
func (p *TreeSitterParser) extractSignatures(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
) ([]parser.Signature, error) {
	signatures := make([]parser.Signature, 0, 32)

	// Get cached query (or create if first time)
	query, err := p.getOrCreateQuery(opts.Language, langQuery, queryTypeSignature)
	if err != nil {
		return nil, fmt.Errorf("failed to create signature query for %s: %w", opts.Language, err)
	}
	// Note: query.Close() is NOT called here because the query is cached for reuse

	// Execute query
	qc := p.cursorPool.Get().(*sitter.QueryCursor)
	defer p.cursorPool.Put(qc)

	matches := qc.Matches(query, root, content)

	// Process matches
	kindMapping := langQuery.KindMapping()
	captureNames := query.CaptureNames()

	// Track seen signatures by (line, column, name) to avoid duplicates from
	// overlapping query patterns (e.g., TypeScript arrow functions).
	// Using a composite key prevents false deduplication when two
	// distinct symbols start on the same line or even at the same position.
	type dedupKey struct {
		line   int
		column int
		name   string
	}
	seen := make(map[dedupKey]bool)

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		sig := parser.Signature{}
		sigColumn := 0
		var kindNode *sitter.Node

		for _, capture := range match.Captures {
			if int(capture.Index) >= len(captureNames) {
				continue
			}
			name := captureNames[capture.Index]
			node := capture.Node

			if name == CaptureKind {
				kindNode = &node
				continue
			}

			start, end := node.StartByte(), node.EndByte()
			if end > uint(len(content)) || start > end {
				continue
			}
			raw := content[start:end]

			switch name {
			case CaptureName:
				sig.Name = string(raw)
			case CaptureSignature:
				sig.Text = string(bytes.TrimSpace(raw))
				sig.Line = int(node.StartPosition().Row) + 1
				sigColumn = int(node.StartPosition().Column)
				sig.EndLine = int(node.EndPosition().Row) + 1
			case CaptureDoc:
				if len(raw) > 0 {
					sig.Doc = cleanComment(string(raw))
				}
			}
		}

		// Map kind if present
		if kindNode != nil {
			kind := kindNode.Kind()
			if mapped, ok := kindMapping[kind]; ok {
				sig.Kind = mapped
			} else {
				sig.Kind = kind
			}

			// Python: distinguish methods from functions by checking first parameter
			if opts.Language == "python" && sig.Kind == "function" {
				if isPythonMethod(sig.Text) {
					sig.Kind = "method"
				}
			}

			// Java: filter out non-static fields (only keep static fields as "variable")
			if opts.Language == "java" && sig.Kind == "field" {
				if !strings.Contains(sig.Text, "static") {
					continue // skip non-static instance fields
				}
				sig.Kind = "variable" // remap to variable for consistency
			}

			// Swift: class_declaration is used for struct, class, enum, and extension
			// Refine kind based on the declaration keyword
			if opts.Language == "swift" && kind == "class_declaration" {
				sig.Kind = refineSwiftClassKind(sig.Text)
			}

			// Swift: init/deinit/subscript have no name capture, synthesize from kind
			if opts.Language == "swift" && sig.Name == "" {
				switch kind {
				case "init_declaration":
					sig.Name = "init"
				case "deinit_declaration":
					sig.Name = "deinit"
				case "subscript_declaration":
					sig.Name = "subscript"
				}
			}

			// Kotlin: class_declaration is used for class, interface, enum class, etc.
			// Refine kind based on the declaration keyword
			if opts.Language == "kotlin" && kind == "class_declaration" {
				sig.Kind = refineKotlinClassKind(sig.Text)
			}

			// Kotlin: companion_object may not have a name, synthesize "Companion"
			if opts.Language == "kotlin" && sig.Name == "" {
				switch kind {
				case "companion_object":
					sig.Name = "Companion"
				case "secondary_constructor":
					sig.Name = "constructor"
				}
			}

			// C#: filter out non-static/non-const fields (like Java)
			if opts.Language == "csharp" && sig.Kind == "field" {
				if !strings.Contains(sig.Text, "static") && !strings.Contains(sig.Text, "const") {
					continue // skip instance fields
				}
				sig.Kind = "variable" // remap to variable for consistency
			}

			// C#: synthesize names for indexer, operator, and conversion operator
			if opts.Language == "csharp" && sig.Name == "" {
				switch kind {
				case "indexer_declaration":
					sig.Name = "this"
				case "operator_declaration":
					sig.Name = extractCSharpOperatorName(sig.Text)
				case "conversion_operator_declaration":
					sig.Name = extractCSharpConversionOperatorName(sig.Text)
				}
			}

			// C#: record struct → "struct" kind; record/record class → "record" kind (의도적)
			if opts.Language == "csharp" && kind == "record_declaration" {
				if strings.Contains(sig.Text, "record struct") {
					sig.Kind = "struct"
				}
			}

			// C: distinguish between function prototypes and variable declarations
			// Both are "declaration" node type, but function prototypes have ()
			if opts.Language == "c" && kind == "declaration" {
				if strings.Contains(sig.Text, "(") && strings.Contains(sig.Text, ")") {
					// It's a function prototype - keep as "function"
					sig.Kind = "function"
				} else {
					// It's a variable declaration
					sig.Kind = "variable"
				}
			}

			// Lua: refine function_declaration kind based on text pattern
			// (global function, local function, module function M.foo, method M:foo)
			if opts.Language == "lua" && kind == "function_declaration" {
				sig.Kind = refineLuaFunctionKind(sig.Text)
			}

			// Elixir: call and unary_operator are generic node types.
			// Refine kind based on signature text prefix (def, defmodule, @spec, etc.)
			// and filter out non-definition calls (if, case, for, etc.)
			if opts.Language == "elixir" {
				if kind == "call" {
					refinedKind := refineElixirCallKind(sig.Text)
					if refinedKind == "" {
						sig.Name = "" // filter out non-definition calls
					} else {
						sig.Kind = refinedKind
					}
				} else if kind == "unary_operator" {
					refinedKind, realName := refineElixirAttrKind(sig.Text, sig.Name)
					if refinedKind == "" {
						sig.Name = "" // filter out @doc, @moduledoc, etc.
					} else {
						sig.Kind = refinedKind
						sig.Name = realName
					}
				}
			}

			// SQL: extract name from DDL text when not captured by query
			// (CREATE INDEX and CREATE SCHEMA have bare identifiers)
			if opts.Language == "sql" && sig.Name == "" && sig.Text != "" {
				sig.Name = extractSQLDDLName(sig.Text)
				// Fallback: use DDL keyword + line number if name extraction fails
				if sig.Name == "" {
					sig.Name = sqlDDLFallbackName(sig.Text, sig.Line)
				}
			}
		}

		// Only add if we have a name and signature
		if sig.Name != "" && sig.Text != "" {
			// Skip duplicates (same line+column+name already captured by another pattern)
			dk := dedupKey{sig.Line, sigColumn, sig.Name}
			if seen[dk] {
				continue
			}
			seen[dk] = true

			sig.Exported = langQuery.IsExported(sig.Name, sig.Text)

			// Strip body if IncludeBody is false (default)
			if !opts.IncludeBody {
				sig.Text = stripBody(sig.Text, sig.Kind, opts.Language)
			}

			sig.Language = opts.Language
			signatures = append(signatures, sig)
		}
	}

	return signatures, nil
}

// cleanComment removes comment markers from the text.
func cleanComment(text string) string {
	// LuaDoc (--- prefix) — check before -- to avoid partial match
	if strings.HasPrefix(text, "---") {
		return strings.TrimSpace(strings.TrimPrefix(text, "---"))
	}
	// Lua block comment --[[ ... ]]
	if strings.HasPrefix(text, "--[[") {
		inner := strings.TrimPrefix(text, "--[[")
		inner = strings.TrimSuffix(inner, "]]")
		return strings.TrimSpace(inner)
	}
	// Lua single-line (-- prefix)
	if strings.HasPrefix(text, "--") {
		return strings.TrimSpace(strings.TrimPrefix(text, "--"))
	}

	// Elixir single-line (# prefix)
	if strings.HasPrefix(text, "#") {
		return strings.TrimSpace(strings.TrimPrefix(text, "#"))
	}

	// Remove // prefix for single-line comments
	if strings.HasPrefix(text, "//") {
		return strings.TrimSpace(strings.TrimPrefix(text, "//"))
	}

	// Remove /* */ for multi-line comments
	if strings.HasPrefix(text, "/*") && strings.HasSuffix(text, "*/") {
		inner := strings.TrimPrefix(text, "/*")
		inner = strings.TrimSuffix(inner, "*/")
		return strings.TrimSpace(inner)
	}

	return strings.TrimSpace(text)
}

// stripBody removes the function/method body from the signature text.
// It preserves only the signature line (declaration without implementation).
func stripBody(text, kind, language string) string {
	switch language {
	case "go":
		return stripGoBody(text, kind)
	case "typescript", "tsx", "javascript", "jsx":
		return stripTypeScriptBody(text, kind)
	case "python":
		return stripPythonBody(text, kind)
	case "c":
		return stripCBody(text, kind)
	case "cpp":
		return stripCppBody(text, kind)
	case "java":
		return stripJavaBody(text, kind)
	case "rust":
		return stripRustBody(text, kind)
	case "swift":
		return stripSwiftBody(text, kind)
	case "kotlin":
		return stripKotlinBody(text, kind)
	case "csharp":
		return stripCSharpBody(text, kind)
	case "lua":
		return stripLuaBody(text, kind)
	case "shell":
		return stripShellBody(text, kind)
	case "php":
		return stripPHPBody(text, kind)
	case "ruby":
		return stripRubyBody(text, kind)
	case "scala":
		return stripScalaBody(text, kind)
	case "elixir":
		return stripElixirBody(text, kind)
	case "sql":
		return stripSQLBody(text, kind)
	case "yaml":
		return stripYAMLBody(text, kind)
	case "toml":
		return stripTOMLBody(text, kind)
	default:
		return text
	}
}

// stripGoBody removes the body from Go function/method/type declarations.
func stripGoBody(text, kind string) string {
	switch kind {
	case "function", "method":
		// Find the opening brace and remove everything from there.
		// Use > 0 (not >= 0): if '{' is at index 0, there's no signature
		// text to keep, so we return the original text unchanged.
		// This pattern is intentionally used in all stripXxxBody functions.
		braceIdx := strings.Index(text, "{")
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "type":
		// For type declarations, keep the entire type spec
		// e.g., "type Foo struct { ... }" -> "type Foo struct { ... }"
		// or "type Foo interface { ... }" -> keep full interface
		// Actually for signatures, we might want to keep the structure
		// but for v0.3.0, let's keep type declarations as-is for now
		return text
	case "variable":
		// Variables: keep full text with value
		return text
	}
	return text
}

// stripTypeScriptBody removes the body from TypeScript/JavaScript declarations.
func stripTypeScriptBody(text, kind string) string {
	switch kind {
	case "function", "method", "export":
		// Find the opening brace for function body
		// Need to be careful with type annotations that contain { }
		// "export" kind is used for exported function declarations
		result := stripTSFunctionBody(text)
		return result
	case "class":
		// For classes, remove the class body but keep the declaration
		// e.g., "class Foo extends Bar { ... }" -> "class Foo extends Bar"
		braceIdx := findTSClassBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "interface", "type":
		// Keep interface/type declarations as-is (they define structure)
		return text
	case "arrow":
		// Arrow functions in variable declarations
		if strings.Contains(text, "=>") {
			return stripTSFunctionBody(text)
		}
		return text
	case "variable":
		// Module-level variables: keep full text with value
		// But if it's an arrow function, strip the body
		if strings.Contains(text, "=>") {
			return stripTSFunctionBody(text)
		}
		return text
	}
	// Default: try to strip body if it looks like a function
	if strings.Contains(text, "{") || strings.Contains(text, "=>") {
		return stripTSFunctionBody(text)
	}
	return text
}

// stripTSFunctionBody removes the function body from a TypeScript function.
// It handles regular functions, methods, and arrow functions.
func stripTSFunctionBody(text string) string {
	// Handle arrow functions: const foo = (args): type => { body } or const foo = (args): type => expr
	// > 0: '=>' at index 0 means no signature before the arrow
	if arrowIdx := strings.Index(text, "=>"); arrowIdx > 0 {
		// Remove everything from => onwards, keep only the signature
		return strings.TrimSpace(text[:arrowIdx])
	}

	// Handle regular functions: function foo(args): type { body }
	// Find the last ) before { that's part of the signature (not in type annotation)
	braceIdx := findFunctionBodyStart(text)
	if braceIdx > 0 {
		return strings.TrimSpace(text[:braceIdx])
	}

	return text
}

// findFunctionBodyStart finds the index where the function body starts.
// It handles nested braces in type annotations.
func findFunctionBodyStart(text string) int {
	// Find the opening brace that starts the function body
	// This is tricky because type annotations can contain { }
	// We look for { that follows ) or a type annotation

	parenDepth := 0
	angleDepth := 0
	lastParenClose := -1

	for i, ch := range text {
		switch ch {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
				if parenDepth == 0 {
					lastParenClose = i
				}
			}
		case '<':
			angleDepth++
		case '>':
			if angleDepth > 0 {
				angleDepth--
			}
		case '{':
			// Only consider { as body start if we're not inside angle brackets
			// and we've seen the closing paren of the parameter list
			if angleDepth == 0 && parenDepth == 0 && lastParenClose >= 0 {
				return i
			}
		}
	}

	return -1
}

// findTSClassBodyStart finds where the class body starts.
func findTSClassBodyStart(text string) int {
	// Class body starts after class declaration, implements, extends clauses
	// Look for the first { at depth 0
	angleDepth := 0
	for i, ch := range text {
		switch ch {
		case '<':
			angleDepth++
		case '>':
			if angleDepth > 0 {
				angleDepth--
			}
		case '{':
			if angleDepth == 0 {
				return i
			}
		}
	}
	return -1
}

// stripPythonBody removes the body from Python function/class declarations.
func stripPythonBody(text, kind string) string {
	switch kind {
	case "function", "method":
		// Find the colon that ends the signature and remove everything after
		colonIdx := findPythonBodyStart(text)
		if colonIdx > 0 {
			return strings.TrimSpace(text[:colonIdx])
		}
	case "class":
		// For classes, also strip at the colon
		colonIdx := findPythonBodyStart(text)
		if colonIdx > 0 {
			return strings.TrimSpace(text[:colonIdx])
		}
	case "variable":
		// Variables: keep full text with value
		return text
	}
	return text
}

// findPythonBodyStart finds the index of the colon that starts the body.
// It handles colons inside type annotations like Dict[str, int].
func findPythonBodyStart(text string) int {
	parenDepth := 0
	bracketDepth := 0

	for i, ch := range text {
		switch ch {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '[':
			bracketDepth++
		case ']':
			if bracketDepth > 0 {
				bracketDepth--
			}
		case ':':
			// Only consider : as body start if we're not inside any brackets
			if parenDepth == 0 && bracketDepth == 0 {
				return i
			}
		}
	}
	return -1
}

// stripCBody removes the body from C function declarations.
func stripCBody(text, kind string) string {
	switch kind {
	case "function":
		// Function: remove everything after {
		braceIdx := strings.Index(text, "{")
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "struct", "enum", "typedef", "macro":
		// Type definitions and macros: keep full text
		return text
	case "variable":
		// Variables: keep full text with value
		return text
	}
	return text
}

// stripCppBody removes the body from C++ declarations.
func stripCppBody(text, kind string) string {
	switch kind {
	case "function", "method", "constructor", "destructor":
		// Remove everything after {
		braceIdx := findCppBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "class", "struct", "namespace":
		// For classes/structs/namespaces, remove the body
		braceIdx := findCppBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "template":
		// Templates: find the underlying declaration and strip its body
		braceIdx := findCppBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "enum", "typedef", "macro":
		// Keep full text
		return text
	}
	return text
}

// findCppBodyStart finds the index where the C++ body starts.
// Handles nested angle brackets for templates.
func findCppBodyStart(text string) int {
	parenDepth := 0
	angleDepth := 0

	for i, ch := range text {
		switch ch {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '<':
			angleDepth++
		case '>':
			if angleDepth > 0 {
				angleDepth--
			}
		case '{':
			if angleDepth == 0 && parenDepth == 0 {
				return i
			}
		}
	}
	return -1
}

// isPythonMethod checks if a Python function is actually a method
// by looking for self or cls as the first parameter.
func isPythonMethod(signature string) bool {
	parenStart := strings.Index(signature, "(")
	if parenStart < 0 {
		return false
	}

	// Find matching closing parenthesis (handles nested parens/brackets)
	parenDepth := 0
	bracketDepth := 0
	parenEnd := -1
	for i := parenStart; i < len(signature); i++ {
		switch signature[i] {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
				if parenDepth == 0 {
					parenEnd = i
					goto found
				}
			}
		case '[':
			bracketDepth++
		case ']':
			if bracketDepth > 0 {
				bracketDepth--
			}
		}
	}
found:
	if parenEnd < 0 || parenEnd <= parenStart+1 {
		return false
	}

	params := signature[parenStart+1 : parenEnd]
	firstParam := strings.TrimSpace(strings.Split(params, ",")[0])

	// Remove type annotation if present (e.g., "self: Self" -> "self")
	// > 0: ':' at index 0 means no parameter name before the colon
	if colonIdx := strings.Index(firstParam, ":"); colonIdx > 0 {
		firstParam = strings.TrimSpace(firstParam[:colonIdx])
	}

	return firstParam == "self" || firstParam == "cls"
}

// stripJavaBody removes the body from Java declarations.
func stripJavaBody(text, kind string) string {
	switch kind {
	case "method", "constructor":
		// Abstract methods end with ; (no body)
		if strings.HasSuffix(strings.TrimSpace(text), ";") {
			return text
		}
		braceIdx := findJavaBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "class", "interface", "enum", "annotation", "record":
		braceIdx := findJavaBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "field":
		// Fields: keep full text with value
		return text
	}
	return text
}

// findJavaBodyStart finds the index where the Java body starts.
// Handles nested angle brackets for generics.
func findJavaBodyStart(text string) int {
	parenDepth := 0
	angleDepth := 0

	for i, ch := range text {
		switch ch {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '<':
			angleDepth++
		case '>':
			if angleDepth > 0 {
				angleDepth--
			}
		case '{':
			if angleDepth == 0 && parenDepth == 0 {
				return i
			}
		}
	}
	return -1
}

// stripRustBody removes the body from Rust declarations.
func stripRustBody(text, kind string) string {
	switch kind {
	case "function", "method":
		// Remove everything after {
		braceIdx := findRustBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "struct", "enum", "trait", "impl", "namespace":
		// For types, remove the body
		braceIdx := findRustBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "type", "variable", "macro":
		// Type aliases, constants, statics, macros: keep full text
		return text
	}
	return text
}

// findRustBodyStart finds the index where the Rust body starts.
// Handles nested angle brackets for generics and lifetime annotations.
func findRustBodyStart(text string) int {
	parenDepth := 0
	angleDepth := 0

	for i, ch := range text {
		switch ch {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '<':
			angleDepth++
		case '>':
			if angleDepth > 0 {
				angleDepth--
			}
		case '{':
			if angleDepth == 0 && parenDepth == 0 {
				return i
			}
		}
	}
	return -1
}

// refineSwiftClassKind determines the specific kind of a Swift class_declaration
// based on the declaration keyword (struct, class, enum, extension, actor).
func refineSwiftClassKind(text string) string {
	trimmed := strings.TrimSpace(text)
	// Strip leading modifiers (public, private, internal, open, final, etc.)
	for _, prefix := range []string{"public ", "private ", "internal ", "open ", "fileprivate ", "final ", "@objc "} {
		for strings.HasPrefix(trimmed, prefix) {
			trimmed = strings.TrimPrefix(trimmed, prefix)
			trimmed = strings.TrimSpace(trimmed)
		}
	}
	// Strip attributes like @available(...) etc.
	for strings.HasPrefix(trimmed, "@") {
		// Skip past the attribute
		if idx := strings.Index(trimmed, " "); idx >= 0 {
			trimmed = strings.TrimSpace(trimmed[idx:])
		} else {
			break
		}
	}

	switch {
	case strings.HasPrefix(trimmed, "struct "):
		return "struct"
	case strings.HasPrefix(trimmed, "enum "):
		return "enum"
	case strings.HasPrefix(trimmed, "extension "):
		return "type"
	case strings.HasPrefix(trimmed, "actor "):
		return "class"
	default:
		return "class"
	}
}

// stripSwiftBody removes the body from Swift declarations.
func stripSwiftBody(text, kind string) string {
	switch kind {
	case "function", "method", "constructor", "destructor":
		braceIdx := findSwiftBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "class", "struct", "enum", "interface", "type":
		braceIdx := findSwiftBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "variable":
		// Properties: keep full text (includes type annotation and default value)
		return text
	}
	return text
}

// findSwiftBodyStart finds the index where the Swift body starts.
// Handles nested angle brackets for generics.
func findSwiftBodyStart(text string) int {
	parenDepth := 0
	angleDepth := 0

	for i, ch := range text {
		switch ch {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '<':
			angleDepth++
		case '>':
			if angleDepth > 0 {
				angleDepth--
			}
		case '{':
			if angleDepth == 0 && parenDepth == 0 {
				return i
			}
		}
	}
	return -1
}

// stripKotlinBody removes the body from Kotlin declarations.
func stripKotlinBody(text, kind string) string {
	switch kind {
	case "function":
		// Single-expression functions (fun foo() = expr) have no braces
		if !strings.Contains(text, "{") && strings.Contains(text, "=") {
			return text
		}
		braceIdx := findKotlinBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "class", "interface", "enum", "constructor":
		braceIdx := findKotlinBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "variable", "type":
		// Properties and type aliases: keep full text
		return text
	}
	return text
}

// findKotlinBodyStart finds the index where the Kotlin body starts.
// Handles nested angle brackets for generics and parentheses for parameters.
func findKotlinBodyStart(text string) int {
	parenDepth := 0
	angleDepth := 0

	for i, ch := range text {
		switch ch {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '<':
			angleDepth++
		case '>':
			if angleDepth > 0 {
				angleDepth--
			}
		case '{':
			if angleDepth == 0 && parenDepth == 0 {
				return i
			}
		}
	}
	return -1
}

// refineKotlinClassKind determines the specific kind of a Kotlin class_declaration
// based on the declaration keyword (class, interface, enum class, etc.).
// In tree-sitter-kotlin, interface and enum class share the class_declaration node type.
func refineKotlinClassKind(text string) string {
	trimmed := strings.TrimSpace(text)
	// Strip leading modifiers
	for _, prefix := range []string{
		"public ", "private ", "internal ", "protected ",
		"open ", "final ", "abstract ", "sealed ",
		"data ", "inner ", "annotation ", "value ",
		"external ", "actual ", "expect ", "override ",
	} {
		for strings.HasPrefix(trimmed, prefix) {
			trimmed = strings.TrimPrefix(trimmed, prefix)
			trimmed = strings.TrimSpace(trimmed)
		}
	}
	// Strip annotations like @Serializable, @Retention(AnnotationRetention.RUNTIME)
	for strings.HasPrefix(trimmed, "@") {
		end := 1
		depth := 0
		for end < len(trimmed) {
			c := trimmed[end]
			if c == '(' {
				depth++
			} else if c == ')' {
				depth--
				if depth < 0 {
					break // defensive: unmatched closing parenthesis
				}
			} else if c == ' ' && depth == 0 {
				break
			}
			end++
		}
		if depth != 0 {
			break // defensive: unmatched parentheses, stop stripping annotations
		}
		trimmed = strings.TrimSpace(trimmed[end:])
	}
	// Strip fun keyword for functional interfaces (fun interface)
	trimmed = strings.TrimPrefix(trimmed, "fun ")
	trimmed = strings.TrimSpace(trimmed)

	switch {
	case strings.HasPrefix(trimmed, "interface "):
		return "interface"
	case strings.HasPrefix(trimmed, "enum class "):
		return "enum"
	default:
		return "class"
	}
}

// refineLuaFunctionKind determines the specific kind of a Lua function_declaration
// based on the declaration text (method M:foo, module function M.foo, local function, global function).
func refineLuaFunctionKind(text string) string {
	trimmed := strings.TrimSpace(text)
	funcIdx := strings.Index(trimmed, "function")
	parenIdx := strings.Index(trimmed, "(")
	if funcIdx >= 0 && parenIdx > funcIdx {
		between := trimmed[funcIdx:parenIdx]
		if strings.Contains(between, ":") {
			return "method"
		}
		if strings.Contains(between, ".") {
			return "module_function"
		}
	}
	if strings.HasPrefix(trimmed, "local ") {
		return "local_function"
	}
	return "function"
}

// stripLuaBody removes the body from Lua declarations.
// Lua functions use function...end blocks; the body starts after the parameter list ")".
func stripLuaBody(text, kind string) string {
	switch kind {
	case "function", "method", "module_function", "local_function":
		parenIdx := strings.Index(text, ")")
		if parenIdx >= 0 {
			return strings.TrimSpace(text[:parenIdx+1])
		}
		// Fallback: take first line
		if nlIdx := strings.Index(text, "\n"); nlIdx > 0 {
			return strings.TrimSpace(text[:nlIdx])
		}
	}
	return text
}

// stripPHPBody removes the body from PHP declarations.
// PHP uses { } blocks; the body starts after the opening brace.
func stripPHPBody(text, kind string) string {
	switch kind {
	case "function", "method":
		// Find the opening brace and remove everything after
		braceIdx := findPHPBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "class", "interface", "trait", "enum", "type":
		// For classes/interfaces/traits, remove the body
		braceIdx := findPHPBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "variable":
		// Properties/Constants: keep full text with value
		return text
	}
	return text
}

// stripRubyBody removes the body from Ruby declarations.
// Ruby methods use def...end blocks; we keep only the first line (signature).
func stripRubyBody(text, kind string) string {
	switch kind {
	case "method":
		// For methods: keep from "def" up to and including the closing ")"
		// or end of first line if no params
		parenIdx := strings.Index(text, ")")
		if parenIdx >= 0 {
			return strings.TrimSpace(text[:parenIdx+1])
		}
		// No params - take first line
		if nlIdx := strings.Index(text, "\n"); nlIdx > 0 {
			return strings.TrimSpace(text[:nlIdx])
		}
	case "class", "namespace":
		// For class/module: keep first line (class Foo < Bar / module Foo)
		if nlIdx := strings.Index(text, "\n"); nlIdx > 0 {
			return strings.TrimSpace(text[:nlIdx])
		}
	case "variable":
		// Constants: keep full text with value
		return text
	}
	return text
}

// stripScalaBody removes the body from Scala declarations.
// Scala uses { } blocks for class/trait/object bodies, and = for def/val/var bodies.
func stripScalaBody(text, kind string) string {
	switch kind {
	case "method":
		// For methods: strip body after = (but keep return type)
		// e.g., "def add(a: Int, b: Int): Int = a + b" → "def add(a: Int, b: Int): Int"
		// e.g., "def greet(name: String): String = { ... }" → "def greet(name: String): String"
		bodyIdx := findScalaBodyStart(text)
		if bodyIdx >= 0 {
			return strings.TrimSpace(text[:bodyIdx])
		}
		// Abstract method or no body: keep first line
		if nlIdx := strings.Index(text, "\n"); nlIdx > 0 {
			return strings.TrimSpace(text[:nlIdx])
		}
	case "class", "trait", "enum":
		// For class/trait/object/enum: strip { ... } body, keep declaration line
		braceIdx := findScalaBodyStart(text)
		if braceIdx >= 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
		// No body (e.g., sealed trait Shape): keep first line
		if nlIdx := strings.Index(text, "\n"); nlIdx > 0 {
			return strings.TrimSpace(text[:nlIdx])
		}
	case "type":
		// Type aliases: keep full text
		return text
	case "variable":
		// val/var: keep full text with value
		return text
	}
	return text
}

// findScalaBodyStart finds the index of the body start in a Scala declaration.
// For methods: finds = that is not part of => or >=/<= operator.
// For classes/traits: finds opening {.
// Returns -1 if not found.
func findScalaBodyStart(text string) int {
	parenDepth := 0
	bracketDepth := 0
	inString := false
	prevCh := rune(0)
	for i, ch := range text {
		if inString {
			if ch == '"' && prevCh != '\\' {
				inString = false
			}
			prevCh = ch
			continue
		}
		switch ch {
		case '"':
			inString = true
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '[':
			bracketDepth++
		case ']':
			if bracketDepth > 0 {
				bracketDepth--
			}
		case '{':
			if parenDepth == 0 && bracketDepth == 0 {
				return i
			}
		case '=':
			if parenDepth == 0 && bracketDepth == 0 {
				// Skip => (lambda arrow) and operators like >=, <=
				if i+1 < len(text) && text[i+1] == '>' {
					// => arrow, skip
				} else if i > 0 && (text[i-1] == '>' || text[i-1] == '<' || text[i-1] == '!' || text[i-1] == '=') {
					// >=, <=, !=, == operators, skip
				} else {
					return i
				}
			}
		}
		prevCh = ch
	}
	return -1
}

// stripShellBody removes the body from Shell/Bash declarations.
// Shell functions use { } blocks; the body starts after the opening brace.
func stripShellBody(text, kind string) string {
	switch kind {
	case "function":
		// Find the opening brace and remove everything after
		braceIdx := strings.Index(text, "{")
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
		// Function without braces (rare): take first line
		if nlIdx := strings.Index(text, "\n"); nlIdx > 0 {
			return strings.TrimSpace(text[:nlIdx])
		}
	case "variable":
		// Variables: keep full text with value
		return text
	}
	return text
}

// findPHPBodyStart finds the index where the PHP body starts.
// It handles nested parentheses and brackets.
func findPHPBodyStart(text string) int {
	parenDepth := 0

	for i, ch := range text {
		switch ch {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '{':
			if parenDepth == 0 {
				return i
			}
		}
	}
	return -1
}

// stripCSharpBody removes the body from C# declarations.
func stripCSharpBody(text, kind string) string {
	switch kind {
	case "method", "constructor", "destructor", "function":
		// Expression-bodied members: remove => expr
		if isExpressionBodied(text) {
			arrowIdx := findCSharpArrowIndex(text)
			if arrowIdx > 0 {
				return strings.TrimSpace(text[:arrowIdx])
			}
		}
		// Abstract/interface methods end with ;
		if strings.HasSuffix(strings.TrimSpace(text), ";") {
			return text
		}
		braceIdx := findCSharpBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "class", "struct", "interface", "enum", "record", "namespace":
		braceIdx := findCSharpBodyStart(text)
		if braceIdx > 0 {
			return strings.TrimSpace(text[:braceIdx])
		}
	case "variable":
		// Properties: auto-properties ({ get; set; }) are kept as-is
		// Expression-bodied properties: remove => expr
		if isExpressionBodied(text) {
			arrowIdx := findCSharpArrowIndex(text)
			if arrowIdx > 0 {
				return strings.TrimSpace(text[:arrowIdx])
			}
		}
		// Auto-properties with accessor list — keep full text
		return text
	case "type":
		// Delegate declarations: no body, return as-is
		return text
	case "field":
		// Fields: keep full text
		return text
	}
	return text
}

// findCSharpBodyStart finds the index where the C# body starts.
// It handles nested angle brackets for generics, parentheses, and square brackets.
func findCSharpBodyStart(text string) int {
	parenDepth := 0
	angleDepth := 0
	bracketDepth := 0

	for i, ch := range text {
		switch ch {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '<':
			angleDepth++
		case '>':
			if angleDepth > 0 {
				angleDepth--
			}
		case '[':
			bracketDepth++
		case ']':
			if bracketDepth > 0 {
				bracketDepth--
			}
		case '{':
			if angleDepth == 0 && parenDepth == 0 && bracketDepth == 0 {
				return i
			}
		}
	}
	return -1
}

// isExpressionBodied checks if a C# declaration uses expression body syntax (=>).
// Distinguishes => from => inside lambda expressions by checking context.
func isExpressionBodied(text string) bool {
	// Look for => that is part of the member declaration (not inside a body)
	parenDepth := 0
	angleDepth := 0
	bracketDepth := 0

	for i := 0; i < len(text)-1; i++ {
		ch := text[i]
		switch ch {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '<':
			angleDepth++
		case '>':
			if angleDepth > 0 {
				angleDepth--
			}
		case '[':
			bracketDepth++
		case ']':
			if bracketDepth > 0 {
				bracketDepth--
			}
		case '{':
			// If we hit { before =>, it's not expression-bodied
			if parenDepth == 0 && angleDepth == 0 && bracketDepth == 0 {
				return false
			}
		case '=':
			if i+1 < len(text) && text[i+1] == '>' && parenDepth == 0 && angleDepth == 0 && bracketDepth == 0 {
				return true
			}
		}
	}
	return false
}

// findCSharpArrowIndex finds the index of the => in an expression-bodied member.
func findCSharpArrowIndex(text string) int {
	parenDepth := 0
	angleDepth := 0
	bracketDepth := 0

	for i := 0; i < len(text)-1; i++ {
		ch := text[i]
		switch ch {
		case '(':
			parenDepth++
		case ')':
			if parenDepth > 0 {
				parenDepth--
			}
		case '<':
			angleDepth++
		case '>':
			if angleDepth > 0 {
				angleDepth--
			}
		case '[':
			bracketDepth++
		case ']':
			if bracketDepth > 0 {
				bracketDepth--
			}
		case '{':
			if parenDepth == 0 && angleDepth == 0 && bracketDepth == 0 {
				return -1
			}
		case '=':
			if i+1 < len(text) && text[i+1] == '>' && parenDepth == 0 && angleDepth == 0 && bracketDepth == 0 {
				return i
			}
		}
	}
	return -1
}

// extractCSharpOperatorName extracts the operator symbol from an operator declaration.
// e.g., "public static int operator +(int a, int b) => ..." -> "operator+"
func extractCSharpOperatorName(text string) string {
	idx := strings.Index(text, "operator")
	if idx < 0 {
		return "operator"
	}
	rest := strings.TrimSpace(text[idx+len("operator"):])
	// The operator symbol follows: +, -, *, /, ==, !=, etc.
	// Find the first ( and take everything before it as the operator
	parenIdx := strings.Index(rest, "(")
	if parenIdx > 0 {
		op := strings.TrimSpace(rest[:parenIdx])
		return "operator" + op
	}
	return "operator"
}

// extractCSharpConversionOperatorName extracts the conversion operator name.
// e.g., "public static implicit operator int(...)" -> "implicit operator int"
// e.g., "public static explicit operator string(...)" -> "explicit operator string"
func extractCSharpConversionOperatorName(text string) string {
	// Find "implicit" or "explicit" keyword
	keyword := ""
	if strings.Contains(text, "implicit") {
		keyword = "implicit"
	} else if strings.Contains(text, "explicit") {
		keyword = "explicit"
	}

	// Find "operator" keyword and extract the type after it
	idx := strings.Index(text, "operator")
	if idx < 0 {
		return keyword + " operator"
	}
	rest := strings.TrimSpace(text[idx+len("operator"):])
	// The target type is between "operator" and "("
	parenIdx := strings.Index(rest, "(")
	if parenIdx > 0 {
		targetType := strings.TrimSpace(rest[:parenIdx])
		return keyword + " operator " + targetType
	}
	return keyword + " operator"
}

// extractImports extracts import/export statements from the AST as raw text.
func (p *TreeSitterParser) extractImports(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
) ([]string, error) {
	imports := make([]string, 0, 8)

	importQueryBytes := langQuery.ImportQuery()
	if importQueryBytes == nil || len(importQueryBytes) == 0 {
		return imports, nil
	}

	// Get cached query (or create if first time)
	query, err := p.getOrCreateQuery(opts.Language, langQuery, queryTypeImport)
	if err != nil {
		return nil, fmt.Errorf("failed to create import query for %s: %w", opts.Language, err)
	}
	// Note: query.Close() is NOT called here because the query is cached for reuse

	// Execute query
	qc := p.cursorPool.Get().(*sitter.QueryCursor)
	defer p.cursorPool.Put(qc)

	matches := qc.Matches(query, root, content)
	captureNames := query.CaptureNames()

	// Track seen imports by start position to avoid duplicates
	seenPositions := make(map[uint]bool)

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		// luaRequireFn holds the value of @_fn capture for Lua function-call import patterns
		// (e.g., Lua require()). Used for Go-side predicate filtering.
		luaRequireFn := ""
		var importNode *sitter.Node

		for _, capture := range match.Captures {
			if int(capture.Index) >= len(captureNames) {
				continue
			}
			name := captureNames[capture.Index]
			node := capture.Node

			switch name {
			case CaptureImportPath:
				importNode = &node
			case CaptureLuaRequireFn:
				start, end := node.StartByte(), node.EndByte()
				if end <= uint(len(content)) && start <= end {
					luaRequireFn = string(content[start:end])
				}
			}
		}

		// Go-side filtering: if @_fn was captured, it must start with "require".
		// This covers Lua's require() and Ruby's require/require_relative.
		// The tree-sitter predicate is not evaluated by the go-tree-sitter binding.
		if luaRequireFn != "" && !strings.HasPrefix(luaRequireFn, "require") {
			continue
		}

		// Extract raw text from the import node
		if importNode != nil {
			startByte := (*importNode).StartByte()
			// Deduplicate by start position
			if seenPositions[startByte] {
				continue
			}
			seenPositions[startByte] = true

			endByte := (*importNode).EndByte()
			if endByte > uint(len(content)) || startByte > endByte {
				continue
			}
			rawBytes := content[startByte:endByte]

			// Elixir: the broad import query pattern also matches defmodule/defprotocol/defimpl
			// (since they also take alias arguments). Filter them out on byte slice
			// to avoid string conversion for skipped entries.
			if opts.Language == "elixir" {
				skip := false
				for _, defKw := range elixirDefPrefixes {
					if bytes.HasPrefix(rawBytes, defKw) {
						skip = true
						break
					}
				}
				if skip {
					continue
				}
			}

			// Remove blank lines (Go module group separators, etc.)
			rawText := removeBlankLines(string(rawBytes))
			if rawText != "" {
				imports = append(imports, rawText)
			}
		}
	}

	return imports, nil
}

// extractCalls extracts function call references from the AST using the language call query.
// It matches each call to its enclosing function based on line ranges from signatures.
func (p *TreeSitterParser) extractCalls(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
	signatures []parser.Signature,
) ([]parser.FunctionCall, error) {
	calls := make([]parser.FunctionCall, 0, 16)

	callQueryBytes := langQuery.CallQuery()
	if len(callQueryBytes) == 0 {
		return calls, nil
	}

	// Get cached query (or create if first time)
	query, err := p.getOrCreateQuery(opts.Language, langQuery, queryTypeCall)
	if err != nil {
		return nil, fmt.Errorf("failed to create call query for %s: %w", opts.Language, err)
	}

	// Execute query
	qc := p.cursorPool.Get().(*sitter.QueryCursor)
	defer p.cursorPool.Put(qc)

	matches := qc.Matches(query, root, content)
	captureNames := query.CaptureNames()

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		var callee string
		var callLine int

		for _, capture := range match.Captures {
			if int(capture.Index) >= len(captureNames) {
				continue
			}
			name := captureNames[capture.Index]
			node := capture.Node

			if name == CaptureCallee {
				start, end := node.StartByte(), node.EndByte()
				if end > uint(len(content)) || start > end {
					continue
				}
				callee = string(content[start:end])
				callLine = int(node.StartPosition().Row) + 1
			}
		}

		if callee == "" {
			continue
		}

		// Find enclosing function
		caller := findEnclosingFunction(signatures, callLine)

		calls = append(calls, parser.FunctionCall{
			Caller: caller,
			Callee: callee,
			Line:   callLine,
		})
	}

	return calls, nil
}

// findEnclosingFunction returns the name of the innermost function/method
// whose Line..EndLine range contains the given line. When multiple signatures
// overlap (e.g., a class containing a method), the narrowest range wins.
// Returns empty string if the call is at top-level (not inside any function).
func findEnclosingFunction(signatures []parser.Signature, line int) string {
	bestName := ""
	bestSpan := int(^uint(0) >> 1) // max int

	for _, sig := range signatures {
		// Skip signatures with invalid EndLine (e.g., parse errors where
		// EndLine is 0); they cannot reliably enclose any line.
		if sig.EndLine == 0 {
			continue
		}
		if sig.Line <= line && line <= sig.EndLine {
			span := sig.EndLine - sig.Line
			if span < bestSpan {
				bestSpan = span
				bestName = sig.Name
			}
		}
	}
	return bestName
}

// removeBlankLines removes empty lines from the import text.
// This is used to clean up Go import blocks that may have blank lines
// between import groups.
func removeBlankLines(text string) string {
	var buf strings.Builder
	buf.Grow(len(text))
	first := true
	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) != "" {
			if !first {
				buf.WriteByte('\n')
			}
			buf.WriteString(line)
			first = false
		}
	}
	return buf.String()
}

// elixirDefKeywords is the set of Elixir definition macro names.
var elixirDefKeywords = map[string]string{
	"defmodule":   "class",
	"defprotocol": "interface",
	"defimpl":     "impl",
	"def":         "function",
	"defp":        "function",
	"defmacro":    "macro",
	"defmacrop":   "macro",
	"defguard":    "function",
	"defguardp":   "function",
	"defdelegate": "function",
	"defstruct":   "struct",
}

// refineElixirCallKind determines the actual kind for an Elixir call node
// by checking the first word of the signature text. Returns "" for
// non-definition calls (if, case, for, etc.) which should be filtered out.
func refineElixirCallKind(text string) string {
	// Extract the first word (the macro name)
	firstWord := text
	if idx := strings.IndexAny(text, " \t\n("); idx > 0 {
		firstWord = text[:idx]
	}
	if kind, ok := elixirDefKeywords[firstWord]; ok {
		return kind
	}
	return ""
}

// elixirAttrKeywords is the set of Elixir module attribute names that are declarations.
var elixirAttrKeywords = map[string]bool{
	"spec":     true,
	"type":     true,
	"typep":    true,
	"opaque":   true,
	"callback": true,
}

// refineElixirAttrKind determines the kind and real name for an Elixir
// module attribute (@spec, @type, etc.). Returns "" kind for non-declaration
// attributes (@doc, @moduledoc, etc.) which should be filtered out.
func refineElixirAttrKind(text, capturedName string) (string, string) {
	if !elixirAttrKeywords[capturedName] {
		return "", ""
	}

	// Extract the real name from the text after "@attr_name "
	// e.g., "@spec hello(integer()) :: atom()" → "hello"
	// e.g., "@type my_type :: integer()" → "my_type"
	prefix := "@" + capturedName + " "
	if !strings.HasPrefix(text, prefix) {
		return "type", capturedName
	}
	rest := text[len(prefix):]

	// Extract the identifier (name before '(' or ' ' or '::')
	realName := rest
	for i, ch := range rest {
		if ch == '(' || ch == ' ' || ch == ':' {
			realName = rest[:i]
			break
		}
	}
	if realName == "" {
		return "type", capturedName
	}
	return "type", realName
}

// stripElixirBody removes the do...end block from Elixir declarations.
func stripElixirBody(text, kind string) string {
	// @spec, @type, etc. — no body to strip
	if strings.HasPrefix(text, "@") {
		return text
	}

	// defstruct — keep full text (no body block)
	if kind == "struct" {
		return text
	}

	// Find " do\n" pattern (block-style body)
	if idx := strings.Index(text, " do\n"); idx >= 0 {
		return strings.TrimRight(text[:idx], " \t")
	}
	if idx := strings.Index(text, " do\r\n"); idx >= 0 {
		return strings.TrimRight(text[:idx], " \t")
	}

	// Handle ", do:" (inline body)
	if idx := strings.Index(text, ", do:"); idx >= 0 {
		return strings.TrimRight(text[:idx], " \t")
	}

	// Multiline — return first line
	if idx := strings.IndexByte(text, '\n'); idx > 0 {
		return strings.TrimRight(text[:idx], " \t")
	}

	return text
}

// sqlDDLFallbackName generates a fallback name for SQL DDL statements when
// the name cannot be extracted. It uses the DDL keyword (e.g., "INDEX",
// "SCHEMA") plus the line number.
func sqlDDLFallbackName(text string, line int) string {
	upper := strings.ToUpper(text)
	keywords := []string{"INDEX", "SCHEMA", "TABLE", "VIEW", "PROCEDURE", "FUNCTION", "TRIGGER", "SEQUENCE", "TYPE"}
	for _, kw := range keywords {
		if strings.Contains(upper, kw) {
			return fmt.Sprintf("<%s:L%d>", strings.ToLower(kw), line)
		}
	}
	return fmt.Sprintf("<ddl:L%d>", line)
}

// extractSQLDDLName extracts the object name from a SQL DDL statement text.
// Used for CREATE INDEX and CREATE SCHEMA where the name is not captured
// by the tree-sitter query pattern.
func extractSQLDDLName(text string) string {
	upper := strings.ToUpper(text)

	// CREATE [UNIQUE] INDEX [IF NOT EXISTS] name ON ...
	if idx := strings.Index(upper, "INDEX"); idx >= 0 {
		rest := strings.TrimSpace(text[idx+5:])
		// Skip optional "IF NOT EXISTS"
		upperRest := strings.ToUpper(rest)
		if strings.HasPrefix(upperRest, "IF NOT EXISTS") {
			rest = strings.TrimSpace(rest[13:])
		} else if strings.HasPrefix(upperRest, "IF") {
			rest = strings.TrimSpace(rest[2:])
			if upperRestAfterIf := strings.ToUpper(rest); strings.HasPrefix(upperRestAfterIf, "NOT") {
				rest = strings.TrimSpace(rest[3:])
				if upperRestAfterNot := strings.ToUpper(rest); strings.HasPrefix(upperRestAfterNot, "EXISTS") {
					rest = strings.TrimSpace(rest[6:])
				}
			}
		}
		// Take the next word (possibly schema.name)
		return extractNextSQLIdentifier(rest)
	}

	// CREATE SCHEMA [IF NOT EXISTS] name
	if idx := strings.Index(upper, "SCHEMA"); idx >= 0 {
		rest := strings.TrimSpace(text[idx+6:])
		upperRest := strings.ToUpper(rest)
		if strings.HasPrefix(upperRest, "IF NOT EXISTS") {
			rest = strings.TrimSpace(rest[13:])
		}
		return extractNextSQLIdentifier(rest)
	}

	return ""
}

// extractNextSQLIdentifier extracts the next SQL identifier from text.
// Handles plain identifiers, schema-qualified (schema.name), and quoted identifiers.
func extractNextSQLIdentifier(text string) string {
	if len(text) == 0 {
		return ""
	}

	// Handle quoted identifier
	if text[0] == '"' || text[0] == '`' {
		quote := text[0]
		end := strings.IndexByte(text[1:], quote)
		if end < 0 {
			return ""
		}
		name := text[1 : end+1]
		// Check for schema.name after closing quote
		rest := text[end+2:]
		if len(rest) > 0 && rest[0] == '.' {
			return name + "." + extractNextSQLIdentifier(rest[1:])
		}
		return name
	}

	// Plain identifier: letters, digits, underscores
	end := 0
	for end < len(text) && (text[end] == '_' || text[end] == '.' ||
		(text[end] >= 'a' && text[end] <= 'z') ||
		(text[end] >= 'A' && text[end] <= 'Z') ||
		(text[end] >= '0' && text[end] <= '9')) {
		end++
	}
	if end == 0 {
		return ""
	}

	name := text[:end]
	// Don't include trailing dot
	name = strings.TrimRight(name, ".")
	// Skip SQL keywords that might appear as "next word"
	upper := strings.ToUpper(name)
	if upper == "ON" || upper == "CONCURRENTLY" {
		rest := strings.TrimSpace(text[end:])
		return extractNextSQLIdentifier(rest)
	}
	return name
}

// stripSQLBody removes the body from SQL DDL statements.
// - Functions/Procedures: strips AS $$ ... $$ or BEGIN...END body
// - Views: strips AS SELECT... query
// - Tables: keeps column definitions (they ARE the schema)
// - Others: keeps full text
func stripSQLBody(text, kind string) string {
	upper := strings.ToUpper(text)

	// Functions/Procedures: strip body
	if kind == "function" && !strings.Contains(upper, "TRIGGER") {
		return stripSQLFunctionBody(text)
	}

	// Views: strip AS SELECT...
	if kind == "type" && strings.Contains(upper, " VIEW ") {
		return stripSQLViewBody(text)
	}

	return text
}

// stripSQLFunctionBody strips the function/procedure body.
// Keeps: CREATE [OR REPLACE] FUNCTION name(args) RETURNS type [LANGUAGE lang]
// Strips: AS $$ ... $$ or AS $tag$ ... $tag$ or function_body content
func stripSQLFunctionBody(text string) string {
	upper := strings.ToUpper(text)

	// Find "AS" followed by body delimiter ($$, $tag$, ', or BEGIN)
	// Try newline variants first to handle multi-line formatting
	asIdx := -1
	searchFrom := 0
	for {
		// Try " AS\n", " AS\r\n", " AS " in order
		bestIdx := -1
		bestLen := 0
		for _, sep := range []string{" AS\n", " AS\r\n", " AS "} {
			idx := strings.Index(upper[searchFrom:], sep)
			if idx >= 0 && (bestIdx < 0 || idx < bestIdx) {
				bestIdx = idx
				bestLen = len(sep)
			}
		}
		if bestIdx < 0 {
			break
		}
		pos := searchFrom + bestIdx
		after := strings.TrimSpace(text[pos+bestLen:])
		// Check for dollar-quote ($$, $tag$), single-quote, or BEGIN
		if strings.HasPrefix(after, "$") || strings.HasPrefix(after, "'") ||
			strings.HasPrefix(strings.ToUpper(after), "BEGIN") {
			asIdx = pos
			break
		}
		searchFrom = pos + bestLen
	}

	if asIdx >= 0 {
		result := strings.TrimSpace(text[:asIdx])
		// Check if LANGUAGE clause is after the body — append it
		langIdx := strings.LastIndex(upper, "LANGUAGE ")
		if langIdx > asIdx {
			// Extract "LANGUAGE plpgsql" part
			langPart := text[langIdx:]
			// Take until semicolon or end
			if semi := strings.IndexByte(langPart, ';'); semi > 0 {
				langPart = langPart[:semi]
			}
			if nl := strings.IndexByte(langPart, '\n'); nl > 0 {
				langPart = langPart[:nl]
			}
			result += " " + strings.TrimSpace(langPart)
		}
		return result
	}

	// Fallback: first line only
	if idx := strings.IndexByte(text, '\n'); idx > 0 {
		return strings.TrimSpace(text[:idx])
	}

	return text
}

// stripSQLViewBody strips the AS SELECT... part from CREATE VIEW statements.
func stripSQLViewBody(text string) string {
	upper := strings.ToUpper(text)
	// Search newline variants first to avoid matching column alias AS inside SELECT
	for _, sep := range []string{" AS\n", " AS\r\n", " AS "} {
		idx := strings.Index(upper, sep)
		if idx >= 0 {
			return strings.TrimSpace(text[:idx])
		}
	}
	return text
}

// stripYAMLBody strips nested block values from YAML mapping pairs.
// For "key: value" pairs, keeps the full line. For container keys with
// nested block mappings, keeps only "key:" (first line).
// The kind parameter is unused because all YAML captures are
// block_mapping_pair nodes and share the same stripping logic.
func stripYAMLBody(text, _ string) string {
	idx := strings.IndexByte(text, '\n')
	if idx < 0 {
		return text // single-line pair, keep as-is
	}
	return strings.TrimSpace(text[:idx])
}

// stripTOMLBody strips the body of TOML table/table_array sections.
// For [table] and [[array]] headers, keeps only the header line.
// For key=value pairs, keeps the full text.
func stripTOMLBody(text, kind string) string {
	switch kind {
	case "namespace":
		// Table or table_array_element: keep only the header line ([section] or [[section]])
		idx := strings.IndexByte(text, '\n')
		if idx < 0 {
			return text
		}
		return strings.TrimSpace(text[:idx])
	default:
		return text
	}
}
