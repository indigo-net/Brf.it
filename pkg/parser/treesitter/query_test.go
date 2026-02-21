package treesitter

import (
	"testing"
)

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
