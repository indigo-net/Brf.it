package treesitter

import (
	"fmt"
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
		}

		// Only add if we have a name and signature
		if sig.Name != "" && sig.Text != "" {
			// Filter private if needed
			if !opts.IncludePrivate && !isExported(sig.Name, opts.Language) {
				continue
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
	default:
		return false
	}
}
