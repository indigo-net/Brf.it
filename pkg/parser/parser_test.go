package parser

import (
	"testing"
)

func TestSignatureDefaults(t *testing.T) {
	sig := Signature{}

	if sig.Name != "" {
		t.Errorf("expected empty Name, got '%s'", sig.Name)
	}
	if sig.Kind != "" {
		t.Errorf("expected empty Kind, got '%s'", sig.Kind)
	}
	if sig.Text != "" {
		t.Errorf("expected empty Text, got '%s'", sig.Text)
	}
	if sig.Doc != "" {
		t.Errorf("expected empty Doc, got '%s'", sig.Doc)
	}
	if sig.Line != 0 {
		t.Errorf("expected zero Line, got %d", sig.Line)
	}
}

func TestParseResultDefaults(t *testing.T) {
	result := ParseResult{}

	if result.Signatures != nil {
		t.Errorf("expected nil Signatures, got %v", result.Signatures)
	}
	if result.Error != nil {
		t.Errorf("expected nil Error, got %v", result.Error)
	}
}

func TestNodeKind(t *testing.T) {
	node := Node{
		Type:     "function_declaration",
		StartRow: 10,
		EndRow:   20,
	}

	if node.Type != "function_declaration" {
		t.Errorf("expected Type 'function_declaration', got '%s'", node.Type)
	}
}
