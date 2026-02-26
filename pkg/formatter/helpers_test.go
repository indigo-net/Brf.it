package formatter

import "testing"

func TestGetEmptyComment(t *testing.T) {
	tests := []struct {
		lang     string
		expected string
	}{
		{"go", "// (empty)"},
		{"typescript", "// (empty)"},
		{"javascript", "// (empty)"},
		{"java", "// (empty)"},
		{"c", "// (empty)"},
		{"python", "# (empty)"},
		{"ruby", "# (empty)"},
		{"html", "<!-- (empty) -->"},
		{"xml", "<!-- (empty) -->"},
		{"unknown", "// (empty)"},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			result := getEmptyComment(tt.lang)
			if result != tt.expected {
				t.Errorf("getEmptyComment(%q) = %q, want %q", tt.lang, result, tt.expected)
			}
		})
	}
}
