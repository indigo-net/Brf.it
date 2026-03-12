package context

import (
	"strings"
	"testing"
)

func TestBuildTokenTree(t *testing.T) {
	tests := []struct {
		name     string
		root     string
		files    []FileTokenCount
		contains []string
	}{
		{
			name:  "empty files",
			root:  "/project",
			files: nil,
			contains: nil,
		},
		{
			name: "single file",
			root: "/project",
			files: []FileTokenCount{
				{Path: "/project/main.go", Tokens: 150},
			},
			contains: []string{"main.go (150 tokens)", "Total: 150 tokens"},
		},
		{
			name: "multiple files in directory",
			root: "/project",
			files: []FileTokenCount{
				{Path: "/project/pkg/scanner/scanner.go", Tokens: 1234},
				{Path: "/project/pkg/scanner/scanner_test.go", Tokens: 567},
				{Path: "/project/main.go", Tokens: 42},
			},
			contains: []string{
				"main.go (42 tokens)",
				"pkg (1,801 tokens)",
				"scanner.go (1,234 tokens)",
				"scanner_test.go (567 tokens)",
				"Total: 1,843 tokens",
			},
		},
		{
			name: "directory token sum",
			root: "/project",
			files: []FileTokenCount{
				{Path: "/project/a/x.go", Tokens: 100},
				{Path: "/project/a/y.go", Tokens: 200},
			},
			contains: []string{
				"a (300 tokens)",
				"x.go (100 tokens)",
				"y.go (200 tokens)",
				"Total: 300 tokens",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildTokenTree(tt.root, tt.files)

			if tt.files == nil {
				if result != "" {
					t.Errorf("expected empty string for nil files, got %q", result)
				}
				return
			}

			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("expected output to contain %q\ngot:\n%s", s, result)
				}
			}
		})
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{12, "12"},
		{123, "123"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
	}

	for _, tt := range tests {
		result := formatNumber(tt.input)
		if result != tt.expected {
			t.Errorf("formatNumber(%d) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
