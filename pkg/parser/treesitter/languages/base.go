package languages

import "strings"

// hasVisibilityPrefix checks if the first line of sigText starts with the given modifier.
// This avoids false positives from matching modifiers inside function bodies.
func hasVisibilityPrefix(sigText, modifier string) bool {
	first := sigText
	if idx := strings.IndexByte(sigText, '\n'); idx >= 0 {
		first = sigText[:idx]
	}
	return strings.HasPrefix(strings.TrimSpace(first), modifier+" ") ||
		strings.Contains(first, " "+modifier+" ")
}

// BaseQuery provides default implementations for common LanguageQuery methods.
// Embed this struct to get the default Captures() implementation.
type BaseQuery struct{}

// Captures returns the standard capture names used across all language queries.
func (BaseQuery) Captures() []string {
	return []string{
		captureName,
		captureSignature,
		captureDoc,
		captureKind,
	}
}

// ImportQuery returns nil by default (no import extraction support).
func (BaseQuery) ImportQuery() []byte {
	return nil
}

// CallQuery returns nil by default (no call extraction support).
func (BaseQuery) CallQuery() []byte {
	return nil
}

// IsExported returns true for any non-empty name by default.
// Languages with explicit visibility rules should override this method.
func (BaseQuery) IsExported(name, _ string) bool {
	return len(name) > 0
}
