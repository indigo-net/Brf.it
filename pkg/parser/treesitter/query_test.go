package treesitter

import (
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser/treesitter/languages"
)

func TestLanguageQueryInterface(t *testing.T) {
	// Verify GoQuery implements LanguageQuery
	var _ LanguageQuery = (*languages.GoQuery)(nil)
	// Verify TypeScriptQuery implements LanguageQuery
	var _ LanguageQuery = (*languages.TypeScriptQuery)(nil)
}

func TestCaptureDefinitions(t *testing.T) {
	// Test that capture names are correctly defined
	expected := []string{"name", "signature", "doc", "kind"}

	for _, exp := range expected {
		switch exp {
		case CaptureName:
		case CaptureSignature:
		case CaptureDoc:
		case CaptureKind:
			// OK - constant exists
		default:
			t.Errorf("unexpected capture name: %s", exp)
		}
	}
}
