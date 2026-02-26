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
		},
	}
}

// Parse parses the given content and returns extracted signatures.
func (p *TreeSitterParser) Parse(content string, opts *parser.Options) (*parser.ParseResult, error) {
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
	signatures := p.extractSignatures(tree.RootNode(), []byte(content), query, opts)

	return &parser.ParseResult{
		Language:   lang,
		Signatures: signatures,
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
) []parser.Signature {
	var signatures []parser.Signature

	// Create query
	query, err := sitter.NewQuery(langQuery.Language(), string(langQuery.Query()))
	if err != nil {
		return signatures
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
				sig.Doc = cleanComment(text)
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
		}

		// Only add if we have a name and signature
		if sig.Name != "" && sig.Text != "" {
			// Skip duplicates (same line already captured by another pattern)
			if seenLines[sig.Line] {
				continue
			}
			seenLines[sig.Line] = true

			// Filter private if needed
			if !opts.IncludePrivate {
				// Java: check for private modifier in signature text
				if opts.Language == "java" && isJavaPrivate(sig.Text) {
					continue
				} else if !isExported(sig.Name, opts.Language) {
					continue
				}
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

	return signatures
}

// cleanComment removes comment markers from the text.
func cleanComment(text string) string {
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
		// Go: first letter uppercase
		return name[0] >= 'A' && name[0] <= 'Z'
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
	case "java":
		// Java: visibility is determined by modifiers, not name
		// Private filtering is done in extractSignatures using isJavaPrivate()
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
	case "java":
		return stripJavaBody(text, kind)
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
	if strings.Contains(text, "=>") {
		arrowIdx := strings.Index(text, "=>")
		if arrowIdx > 0 {
			// Remove everything from => onwards, keep only the signature
			return strings.TrimSpace(text[:arrowIdx])
		}
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
			parenDepth--
			if parenDepth == 0 {
				lastParenClose = i
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

// isPythonMethod checks if a Python function is actually a method
// by looking for self or cls as the first parameter.
func isPythonMethod(signature string) bool {
	parenStart := strings.Index(signature, "(")
	parenEnd := strings.Index(signature, ")")
	if parenStart < 0 || parenEnd < 0 || parenEnd <= parenStart+1 {
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

// isJavaPrivate checks if a Java signature has private visibility.
func isJavaPrivate(signature string) bool {
	words := strings.Fields(signature)
	for _, word := range words {
		if word == "private" {
			return true
		}
		// Stop after modifiers (hit declaration keyword)
		if word == "class" || word == "interface" || word == "enum" ||
			word == "void" || word == "@interface" || word == "record" {
			break
		}
	}
	return false
}
