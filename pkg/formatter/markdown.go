package formatter

import (
	"bytes"
	"strconv"
	"strings"
)

// MarkdownFormatter implements Formatter for Markdown output.
type MarkdownFormatter struct{}

// NewMarkdownFormatter creates a new MarkdownFormatter.
func NewMarkdownFormatter() *MarkdownFormatter {
	return &MarkdownFormatter{}
}

// Name returns the formatter name.
func (f *MarkdownFormatter) Name() string {
	return "markdown"
}

// Format implements Formatter interface.
func (f *MarkdownFormatter) Format(data *PackageData) ([]byte, error) {
	var buf bytes.Buffer

	// Header with path
	if data.RootPath != "" {
		buf.WriteString("# Code Summary: ")
		buf.WriteString(data.RootPath)
		buf.WriteString("\n\n")
	} else {
		buf.WriteString("# Code Summary\n\n")
	}

	// Version info
	if data.Version != "" {
		buf.WriteString("*brf.it ")
		buf.WriteString(data.Version)
		buf.WriteString("*\n\n")
	}

	// Directory Tree
	if data.Tree != "" {
		buf.WriteString("## Directory Tree\n\n")
		buf.WriteString("```\n")
		buf.WriteString(data.Tree)
		buf.WriteString("\n```\n\n")
	}

	// Global imports (when dedupe mode is enabled)
	if data.DedupeImports && len(data.GlobalImports) > 0 {
		buf.WriteString("## Global Imports\n\n")
		buf.WriteString("| Import | Files |\n")
		buf.WriteString("|--------|-------|\n")
		for _, ic := range data.GlobalImports {
			buf.WriteString("| `")
			buf.WriteString(escapeMarkdown(ic.Import))
			buf.WriteString("` | ")
			buf.WriteString(strconv.Itoa(ic.Count))
			buf.WriteString(" |\n")
		}
		buf.WriteString("\n")
	}

	// Files
	buf.WriteString("## Files\n\n")
	for _, file := range data.Files {
		buf.WriteString("### ")
		buf.WriteString(file.Path)
		buf.WriteString("\n\n")

		// Imports (within file block) - only if not deduping
		hasRenderedImports := false
		if file.Error == nil && data.IncludeImports && len(file.RawImports) > 0 && !data.DedupeImports {
			hasRenderedImports = true
		}

		if file.Error != nil {
			buf.WriteString("> **Error:** ")
			buf.WriteString(escapeMarkdown(file.Error.Error()))
			buf.WriteString("\n\n")
		} else {
			// 빈 파일 확인
			isEmpty := len(file.Signatures) == 0 && !hasRenderedImports

			buf.WriteString("```")
			buf.WriteString(file.Language)
			buf.WriteByte('\n')
			if isEmpty {
				buf.WriteString(getEmptyComment(file.Language))
				buf.WriteString("\n")
			} else {
				// Include imports at the top of the code block
				if hasRenderedImports {
					for _, imp := range file.RawImports {
						buf.WriteString(imp)
						buf.WriteString("\n")
					}
				}
				// Then include signatures
				for _, sig := range file.Signatures {
					buf.WriteString(sig.Text)
					buf.WriteString("\n")
				}
			}
			buf.WriteString("```\n")

			// Add docs as quotes (빈 파일이면 건너뜀)
			if !isEmpty {
				for _, sig := range file.Signatures {
					if sig.Doc != "" {
						buf.WriteString("> ")
						buf.WriteString(escapeMarkdown(truncateDoc(sig.Doc, data.MaxDocLength)))
						buf.WriteString("\n")
					}
				}
			}
		}

		buf.WriteByte('\n')
	}

	return buf.Bytes(), nil
}

// escapeMarkdown escapes special characters for Markdown content.
func escapeMarkdown(s string) string {
	// Only escape backticks to avoid breaking code blocks
	s = strings.ReplaceAll(s, "`", "\\`")
	return s
}
