package formatter

import (
	"bytes"
	"fmt"
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

	buf.WriteString("# Brf.it Output\n\n")

	// Directory Tree
	if data.Tree != "" {
		buf.WriteString("## Directory Tree\n\n")
		buf.WriteString("```\n")
		buf.WriteString(data.Tree)
		buf.WriteString("\n```\n\n")
	}

	// Symbols
	if data.TotalSignatures > 0 {
		buf.WriteString("## Symbols\n\n")
		for _, file := range data.Files {
			for _, sig := range file.Signatures {
				buf.WriteString("- `")
				buf.WriteString(escapeMarkdown(sig.Text))
				buf.WriteString("`\n")
			}
		}
		buf.WriteString("\n---\n\n")
	}

	// Files
	buf.WriteString("## Files\n\n")
	for _, file := range data.Files {
		buf.WriteString(fmt.Sprintf("### %s\n\n", file.Path))

		if file.Error != nil {
			buf.WriteString("> **Error:** ")
			buf.WriteString(escapeMarkdown(file.Error.Error()))
			buf.WriteString("\n\n")
		} else {
			buf.WriteString(fmt.Sprintf("```%s\n", file.Language))
			for _, sig := range file.Signatures {
				buf.WriteString(sig.Text)
				buf.WriteString("\n")
			}
			buf.WriteString("```\n\n")

			// Add docs as quotes
			for _, sig := range file.Signatures {
				if sig.Doc != "" {
					buf.WriteString("> ")
					buf.WriteString(escapeMarkdown(sig.Doc))
					buf.WriteString("\n\n")
				}
			}
		}

		buf.WriteString("---\n\n")
	}

	return buf.Bytes(), nil
}

// escapeMarkdown escapes special characters for Markdown content.
func escapeMarkdown(s string) string {
	// Only escape backticks to avoid breaking code blocks
	s = strings.ReplaceAll(s, "`", "\\`")
	return s
}
