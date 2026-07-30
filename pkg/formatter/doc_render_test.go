package formatter

import (
	"strings"
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser"
)

// TestMarkdownFormatterMultiLineDocQuoted verifies that every line of a
// multi-line doc comment is prefixed with the blockquote marker, so the
// rendered Markdown does not break out of the quote after the first line.
func TestMarkdownFormatterMultiLineDocQuoted(t *testing.T) {
	f := NewMarkdownFormatter()
	data := &PackageData{
		Files: []FileData{
			{
				Path:     "test.go",
				Language: "go",
				Signatures: []parser.Signature{
					{
						Name:     "Add",
						Kind:     "function",
						Text:     "func Add(a, b int) int",
						Doc:      "Add sums two ints.\nSecond line of docs.",
						Line:     1,
						Language: "go",
						Exported: true,
					},
				},
			},
		},
		TotalSignatures: 1,
	}

	out, err := f.Format(data)
	if err != nil {
		t.Fatalf("Format failed: %v", err)
	}

	want := "> Add sums two ints.\n> Second line of docs."
	if !strings.Contains(string(out), want) {
		t.Errorf("multi-line doc not fully quoted;\nwant to contain: %q\ngot:\n%s", want, out)
	}
}
