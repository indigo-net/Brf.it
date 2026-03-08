package treesitter

import (
	"fmt"
	"regexp"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"

	"github.com/indigo-net/Brf.it/pkg/parser"
	"github.com/indigo-net/Brf.it/pkg/parser/treesitter/languages"
)

// init registers the TreeSitterParser with the default registry.
func init() {
	parser.RegisterParser("go", NewTreeSitterParser())
	parser.RegisterParser("typescript", NewTreeSitterParser())
	parser.RegisterParser("tsx", NewTreeSitterParser())
	parser.RegisterParser("javascript", NewTreeSitterParser())
	parser.RegisterParser("jsx", NewTreeSitterParser())
	parser.RegisterParser("python", NewTreeSitterParser())
	parser.RegisterParser("c", NewTreeSitterParser())
	parser.RegisterParser("java", NewTreeSitterParser())
	parser.RegisterParser("cpp", NewTreeSitterParser())
	parser.RegisterParser("rust", NewTreeSitterParser())
	parser.RegisterParser("swift", NewTreeSitterParser())
	parser.RegisterParser("kotlin", NewTreeSitterParser())
	parser.RegisterParser("csharp", NewTreeSitterParser())
	parser.RegisterParser("lua", NewTreeSitterParser())
	parser.RegisterParser("shell", NewTreeSitterParser())
	parser.RegisterParser("php", NewTreeSitterParser())
}

// TreeSitterParser implements parser.Parser using Tree-sitter.
type TreeSitterParser struct {
	queries map[string]LanguageQuery
}

// NewTreeSitterParser creates a new Tree-sitter based parser.
func NewTreeSitterParser() *TreeSitterParser {
	return &TreeSitterParser{
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
		},
	}
}

// Parse parses the given content and returns extracted signatures.
func (p *TreeSitterParser) Parse(content string, opts *parser.Options) (result *parser.ParseResult, err error) {
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
		return nil, fmt.Errorf("language must be specified")
	}

	// Get language query
	query, ok := p.queries[lang]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}

	// Create parser
	sitterParser := sitter.NewParser()
	defer sitterParser.Close()

	// Set language
	tsLang := query.Language()
	if err := sitterParser.SetLanguage(tsLang); err != nil {
		return nil, fmt.Errorf("failed to set language: %w", err)
	}

	// Parse content
	tree := sitterParser.Parse([]byte(content), nil)
	defer tree.Close()

	if tree == nil {
		return nil, fmt.Errorf("failed to parse content")
	}

	// Extract signatures
	signatures, err := p.extractSignatures(tree.RootNode(), []byte(content), query, opts)
	if err != nil {
		return nil, fmt.Errorf("signature extraction failed: %w", err)
	}

	// Extract imports if requested
	var rawImports []string
	if opts.IncludeImports {
		rawImports, err = p.extractImports(tree.RootNode(), []byte(content), query, opts)
		if err != nil {
			return nil, fmt.Errorf("import extraction failed: %w", err)
		}
	}

	return &parser.ParseResult{
		Language:   lang,
		Signatures: signatures,
		RawImports: rawImports,
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
	var signatures []parser.Signature

	// Create query
	query, err := sitter.NewQuery(langQuery.Language(), string(langQuery.Query()))
	if err != nil {
		return nil, fmt.Errorf("failed to create signature query for %s: %w", opts.Language, err)
	}
	defer query.Close()

	// Execute query
	qc := sitter.NewQueryCursor()
	defer qc.Close()

	matches := qc.Matches(query, root, content)

	// Process matches
	kindMapping := langQuery.KindMapping()
	captureNames := query.CaptureNames()

	// Track seen signatures by line number to avoid duplicates
	// (e.g., TypeScript arrow functions can be captured by multiple patterns)
	seenLines := make(map[int]bool)

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		sig := parser.Signature{}
		var kindNode *sitter.Node

		for _, capture := range match.Captures {
			name := captureNames[capture.Index]
			node := capture.Node
			text := string(content[node.StartByte():node.EndByte()])

			switch name {
			case CaptureName:
				sig.Name = text
			case CaptureSignature:
				sig.Text = strings.TrimSpace(text)
				sig.Line = int(node.StartPosition().Row) + 1
				sig.EndLine = int(node.EndPosition().Row) + 1
			case CaptureDoc:
				if text != "" {
					sig.Doc = cleanComment(text)
				}
			case CaptureKind:
				kindNode = &node
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
		}

		// Only add if we have a name and signature
		if sig.Name != "" && sig.Text != "" {
			// Skip duplicates (same line already captured by another pattern)
			if seenLines[sig.Line] {
				continue
			}
			seenLines[sig.Line] = true

			// Filter private if needed
			if !opts.IncludePrivate && !isExported(sig.Name, opts.Language) {
				continue
			}

			// Strip body if IncludeBody is false (default)
			if !opts.IncludeBody {
				sig.Text = stripBody(sig.Text, sig.Kind, opts.Language)
			}

			sig.Language = opts.Language
			sig.Exported = isExported(sig.Name, opts.Language)
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

// isExported checks if a name is exported/public.
func isExported(name, language string) bool {
	if len(name) == 0 {
		return false
	}

	switch language {
	case "go":
		// Go: module-level symbols are always included
		// (Tree-sitter query already captures only package-level declarations)
		return true
	case "typescript", "tsx", "javascript", "jsx":
		// TypeScript/JavaScript: assume all found signatures are exported
		// (since we query for export_statement patterns)
		return true
	case "python":
		// Python: all elements are considered public (no private filtering)
		return true
	case "c":
		// C: all functions are considered exported (static functions handled separately)
		return true
	case "cpp":
		// C++: all elements are considered exported (access control in class context is complex)
		return true
	case "java":
		// Java: all elements are considered exported (visibility modifiers not filtered)
		return true
	case "rust":
		// Rust: all elements are considered public (user requested private extraction too)
		return true
	case "swift":
		// Swift: all elements are considered public (visibility modifiers preserved in signature text)
		return true
	case "kotlin":
		// Kotlin: default visibility is public, all elements are considered exported
		return true
	case "csharp":
		// C#: all elements are considered exported (visibility modifiers preserved in signature text)
		return true
	case "lua":
		// Lua: all elements are considered public (no access modifiers)
		return true
	case "shell":
		// Shell/Bash: all functions and variables are public
		return true
	case "php":
		// PHP: all elements are considered exported (visibility modifiers preserved in signature text)
		return true
	default:
		return false
	}
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
	default:
		return text
	}
}

// stripGoBody removes the body from Go function/method/type declarations.
func stripGoBody(text, kind string) string {
	switch kind {
	case "function", "method":
		// Find the opening brace and remove everything from there
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

// Regex patterns for TypeScript body stripping
var (
	// Matches function body: starts with { and ends with matching }
	tsFunctionBodyRe = regexp.MustCompile(`\s*\{[\s\S]*\}\s*$`)
	// Matches arrow function body: => { ... } or => expression
	tsArrowBodyRe = regexp.MustCompile(`\s*=>\s*[\s\S]+$`)
	// Matches class body
	tsClassBodyRe = regexp.MustCompile(`\s*\{[\s\S]*\}\s*$`)
)

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
			parenDepth--
		case '[':
			bracketDepth++
		case ']':
			bracketDepth--
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
			parenDepth--
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
			parenDepth--
			if parenDepth == 0 {
				parenEnd = i
				goto found
			}
		case '[':
			bracketDepth++
		case ']':
			bracketDepth--
		}
	}
found:
	if parenEnd < 0 || parenEnd <= parenStart+1 {
		return false
	}

	params := signature[parenStart+1 : parenEnd]
	firstParam := strings.TrimSpace(strings.Split(params, ",")[0])

	// Remove type annotation if present (e.g., "self: Self" -> "self")
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
			parenDepth--
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
			parenDepth--
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
			parenDepth--
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
			parenDepth--
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
			parenDepth--
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
			parenDepth--
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
			parenDepth--
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
			parenDepth--
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
	var imports []string

	importQueryBytes := langQuery.ImportQuery()
	if importQueryBytes == nil || len(importQueryBytes) == 0 {
		return imports, nil
	}

	// Create query
	query, err := sitter.NewQuery(langQuery.Language(), string(importQueryBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create import query for %s: %w", opts.Language, err)
	}
	defer query.Close()

	// Execute query
	qc := sitter.NewQueryCursor()
	defer qc.Close()

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
			name := captureNames[capture.Index]
			node := capture.Node
			text := string(content[node.StartByte():node.EndByte()])

			switch name {
			case CaptureImportPath:
				importNode = &node
			case CaptureLuaRequireFn:
				luaRequireFn = text
			}
		}

		// Go-side filtering: if @_fn was captured, it must be "require".
		// The tree-sitter (#eq? @_fn "require") predicate is not evaluated
		// by the go-tree-sitter binding at runtime, so we enforce it here.
		if luaRequireFn != "" && luaRequireFn != "require" {
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

			rawText := string(content[startByte:(*importNode).EndByte()])
			// Remove blank lines (Go module group separators, etc.)
			rawText = removeBlankLines(rawText)
			if rawText != "" {
				imports = append(imports, rawText)
			}
		}
	}

	return imports, nil
}

// removeBlankLines removes empty lines from the import text.
// This is used to clean up Go import blocks that may have blank lines
// between import groups.
func removeBlankLines(text string) string {
	lines := strings.Split(text, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}
