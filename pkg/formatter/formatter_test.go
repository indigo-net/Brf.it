package formatter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser"
)

func TestXMLFormatterImplementsFormatter(t *testing.T) {
	var _ Formatter = (*XMLFormatter)(nil)
}

func TestMarkdownFormatterImplementsFormatter(t *testing.T) {
	var _ Formatter = (*MarkdownFormatter)(nil)
}

func TestXMLFormatterFormat(t *testing.T) {
	formatter := NewXMLFormatter()

	data := &PackageData{
		Tree: "pkg/\n└── test.go",
		Files: []FileData{
			{
				Path:     "pkg/test.go",
				Language: "go",
				Signatures: []parser.Signature{
					{
						Name:     "Add",
						Kind:     "function",
						Text:     "func Add(a, b int) int",
						Doc:      "Add returns the sum of two integers.",
						Line:     5,
						Language: "go",
						Exported: true,
					},
				},
			},
		},
		TotalSignatures: 1,
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	// Verify XML structure
	outputStr := string(output)

	if !strings.Contains(outputStr, `<?xml version="1.0"`) {
		t.Error("expected XML declaration")
	}

	if !strings.Contains(outputStr, "<brfit>") {
		t.Error("expected <brfit> root element")
	}

	if !strings.Contains(outputStr, "<metadata>") {
		t.Error("expected <metadata> element")
	}

	if !strings.Contains(outputStr, `<file path="pkg/test.go" language="go"`) {
		t.Error("expected file element with path and language attributes")
	}

	if !strings.Contains(outputStr, "<signature>func Add(a, b int) int</signature>") {
		t.Error("expected signature element")
	}

	if !strings.Contains(outputStr, "<doc>Add returns the sum of two integers.</doc>") {
		t.Error("expected doc element")
	}
}

func TestXMLFormatterFormatWithError(t *testing.T) {
	formatter := NewXMLFormatter()

	data := &PackageData{
		Files: []FileData{
			{
				Path:     "test.py",
				Language: "python",
				Error:    fmt.Errorf("no parser for language: python"),
			},
		},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "<error>no parser for language: python</error>") {
		t.Error("expected error element")
	}
}

func TestMarkdownFormatterFormat(t *testing.T) {
	formatter := NewMarkdownFormatter()

	data := &PackageData{
		Tree: "pkg/\n└── test.go",
		Files: []FileData{
			{
				Path:     "pkg/test.go",
				Language: "go",
				Signatures: []parser.Signature{
					{
						Name:     "Add",
						Kind:     "function",
						Text:     "func Add(a, b int) int",
						Doc:      "Add returns the sum.",
						Line:     5,
						Language: "go",
						Exported: true,
					},
				},
			},
		},
		TotalSignatures: 1,
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "# Brf.it Output") {
		t.Error("expected header")
	}

	if !strings.Contains(outputStr, "## Directory Tree") {
		t.Error("expected Directory Tree section")
	}

	if !strings.Contains(outputStr, "## Symbols") {
		t.Error("expected Symbols section")
	}

	if !strings.Contains(outputStr, "## Files") {
		t.Error("expected Files section")
	}

	if !strings.Contains(outputStr, "### pkg/test.go") {
		t.Error("expected file path as heading")
	}

	if !strings.Contains(outputStr, "```go") {
		t.Error("expected code block with language")
	}

	if !strings.Contains(outputStr, "func Add(a, b int) int") {
		t.Error("expected signature in code block")
	}
}

func TestFormatterNames(t *testing.T) {
	xmlFormatter := NewXMLFormatter()
	if xmlFormatter.Name() != "xml" {
		t.Errorf("expected name 'xml', got '%s'", xmlFormatter.Name())
	}

	mdFormatter := NewMarkdownFormatter()
	if mdFormatter.Name() != "markdown" {
		t.Errorf("expected name 'markdown', got '%s'", mdFormatter.Name())
	}
}

func TestXMLFormatterEscapeXML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"a & b", "a &amp; b"},
		{"<tag>", "&lt;tag&gt;"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"it's", "it&apos;s"},
		{"all & < > \" '", "all &amp; &lt; &gt; &quot; &apos;"},
	}

	for _, tt := range tests {
		result := escapeXML(tt.input)
		if result != tt.expected {
			t.Errorf("escapeXML(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestMarkdownFormatterEscapeMarkdown(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"code `here`", "code \\`here\\`"},
		{"no escape needed", "no escape needed"},
	}

	for _, tt := range tests {
		result := escapeMarkdown(tt.input)
		if result != tt.expected {
			t.Errorf("escapeMarkdown(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestXMLFormatterEmptyData(t *testing.T) {
	formatter := NewXMLFormatter()

	data := &PackageData{
		Files: []FileData{},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "<brfit>") {
		t.Error("expected <brfit> root element")
	}

	if !strings.Contains(outputStr, "<files>\n  </files>") {
		t.Error("expected empty files section")
	}
}

func TestMarkdownFormatterEmptyData(t *testing.T) {
	formatter := NewMarkdownFormatter()

	data := &PackageData{
		Files: []FileData{},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "# Brf.it Output") {
		t.Error("expected header")
	}

	if !strings.Contains(outputStr, "## Files") {
		t.Error("expected Files section")
	}
}
