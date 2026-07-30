package treesitter

import (
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser"
)

// TestParseAttachesDocComment verifies that a documentation comment
// immediately preceding a declaration is attached to that signature's Doc.
func TestParseAttachesDocComment(t *testing.T) {
	code := []byte(`package main

// Greet returns a friendly greeting for the given name.
func Greet(name string) string {
	return "hi " + name
}
`)

	p := NewTreeSitterParser()
	result, err := p.Parse(code, &parser.Options{Language: "go"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	var greet *parser.Signature
	for i := range result.Signatures {
		if result.Signatures[i].Name == "Greet" {
			greet = &result.Signatures[i]
			break
		}
	}
	if greet == nil {
		t.Fatalf("Greet signature not found; got %d signatures", len(result.Signatures))
	}

	want := "Greet returns a friendly greeting for the given name."
	if greet.Doc != want {
		t.Errorf("Doc not attached: got %q, want %q", greet.Doc, want)
	}
}

// TestParseAttachesMultiLineDoc verifies that a contiguous run of single-line
// comments above a declaration is joined into that signature's Doc.
func TestParseAttachesMultiLineDoc(t *testing.T) {
	code := []byte(`package main

// Add returns the sum of two integers.
// It is used in the arithmetic examples.
func Add(a, b int) int {
	return a + b
}
`)

	p := NewTreeSitterParser()
	result, err := p.Parse(code, &parser.Options{Language: "go"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	var add *parser.Signature
	for i := range result.Signatures {
		if result.Signatures[i].Name == "Add" {
			add = &result.Signatures[i]
			break
		}
	}
	if add == nil {
		t.Fatalf("Add signature not found; got %d signatures", len(result.Signatures))
	}

	want := "Add returns the sum of two integers.\nIt is used in the arithmetic examples."
	if add.Doc != want {
		t.Errorf("multi-line Doc not joined: got %q, want %q", add.Doc, want)
	}
}

// TestParseDoesNotAttachSeparatedComment verifies that a comment separated from
// a declaration by a blank line is NOT attached (matches godoc semantics).
func TestParseDoesNotAttachSeparatedComment(t *testing.T) {
	code := []byte(`package main

// This is a stray comment, not documentation.

func Lonely() {}
`)

	p := NewTreeSitterParser()
	result, err := p.Parse(code, &parser.Options{Language: "go"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	var lonely *parser.Signature
	for i := range result.Signatures {
		if result.Signatures[i].Name == "Lonely" {
			lonely = &result.Signatures[i]
			break
		}
	}
	if lonely == nil {
		t.Fatalf("Lonely signature not found; got %d signatures", len(result.Signatures))
	}

	if lonely.Doc != "" {
		t.Errorf("separated comment should not attach: got Doc=%q", lonely.Doc)
	}
}

// TestParseDoesNotAttachTrailingComment verifies that a trailing/inline comment
// on a declaration line is not mistaken for the following declaration's doc.
func TestParseDoesNotAttachTrailingComment(t *testing.T) {
	code := []byte(`package main

func Foo() {} // trailing note about Foo
func Bar() {}
`)

	p := NewTreeSitterParser()
	result, err := p.Parse(code, &parser.Options{Language: "go"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	byName := map[string]string{}
	for i := range result.Signatures {
		byName[result.Signatures[i].Name] = result.Signatures[i].Doc
	}

	if got := byName["Bar"]; got != "" {
		t.Errorf("trailing comment leaked onto next declaration: Bar.Doc=%q", got)
	}
	if got := byName["Foo"]; got != "" {
		t.Errorf("trailing comment should not be Foo's doc either: Foo.Doc=%q", got)
	}
}

// TestParseAttachesBlockDoc verifies that a leading block comment (/* ... */) is
// attached as the following declaration's doc.
func TestParseAttachesBlockDoc(t *testing.T) {
	code := []byte(`package main

/* Widget builds a widget. */
func Widget() {}
`)

	p := NewTreeSitterParser()
	result, err := p.Parse(code, &parser.Options{Language: "go"})
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	var widget *parser.Signature
	for i := range result.Signatures {
		if result.Signatures[i].Name == "Widget" {
			widget = &result.Signatures[i]
			break
		}
	}
	if widget == nil {
		t.Fatalf("Widget signature not found; got %d signatures", len(result.Signatures))
	}

	want := "Widget builds a widget."
	if widget.Doc != want {
		t.Errorf("block doc not attached: got %q, want %q", widget.Doc, want)
	}
}
