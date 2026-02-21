package tokenizer

import (
	"strings"
	"testing"
)

func TestNoOpTokenizerImplementsTokenizer(t *testing.T) {
	var _ Tokenizer = (*NoOpTokenizer)(nil)
}

func TestTiktokenTokenizerImplementsTokenizer(t *testing.T) {
	var _ Tokenizer = (*TiktokenTokenizer)(nil)
}

func TestNoOpTokenizerCount(t *testing.T) {
	tokenizer := NewNoOpTokenizer()

	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"simple", "hello world"},
		{"long", "this is a longer text with multiple words"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := tokenizer.Count(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if count != 0 {
				t.Errorf("expected 0, got %d", count)
			}
		})
	}
}

func TestNoOpTokenizerName(t *testing.T) {
	tokenizer := NewNoOpTokenizer()

	if tokenizer.Name() != "noop" {
		t.Errorf("expected name 'noop', got '%s'", tokenizer.Name())
	}
}

func TestTiktokenTokenizerCount(t *testing.T) {
	tokenizer, err := NewTiktokenTokenizer()
	if err != nil {
		t.Fatalf("failed to create tokenizer: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		minCount int
		maxCount int
	}{
		{
			name:     "empty string",
			input:    "",
			minCount: 0,
			maxCount: 0,
		},
		{
			name:     "simple text",
			input:    "hello world",
			minCount: 2,
			maxCount: 3,
		},
		{
			name:     "code snippet",
			input:    "func main() { fmt.Println(\"hello\") }",
			minCount: 8,
			maxCount: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := tokenizer.Count(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if count < tt.minCount {
				t.Errorf("expected at least %d tokens, got %d", tt.minCount, count)
			}
			if count > tt.maxCount {
				t.Errorf("expected at most %d tokens, got %d", tt.maxCount, count)
			}
		})
	}
}

func TestTiktokenTokenizerName(t *testing.T) {
	tokenizer, err := NewTiktokenTokenizer()
	if err != nil {
		t.Fatalf("failed to create tokenizer: %v", err)
	}

	if !strings.Contains(tokenizer.Name(), "tiktoken") {
		t.Errorf("expected name to contain 'tiktoken', got '%s'", tokenizer.Name())
	}

	if !strings.Contains(tokenizer.Name(), "cl100k") {
		t.Errorf("expected name to contain 'cl100k', got '%s'", tokenizer.Name())
	}
}

func TestTiktokenTokenizerConsistency(t *testing.T) {
	tokenizer, err := NewTiktokenTokenizer()
	if err != nil {
		t.Fatalf("failed to create tokenizer: %v", err)
	}

	text := "package main\n\nfunc main() {}\n"

	count1, err := tokenizer.Count(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	count2, err := tokenizer.Count(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count1 != count2 {
		t.Errorf("inconsistent token count: %d vs %d", count1, count2)
	}
}

func TestTiktokenTokenizerSpecialCharacters(t *testing.T) {
	tokenizer, err := NewTiktokenTokenizer()
	if err != nil {
		t.Fatalf("failed to create tokenizer: %v", err)
	}

	tests := []struct {
		name  string
		input string
	}{
		{"unicode", "hello 안녕하세요"},
		{"emoji", "code 👍 test"},
		{"xml tags", "<xml><element>content</element></xml>"},
		{"markdown", "# Header\n\n**bold** and *italic*"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := tokenizer.Count(tt.input)
			if err != nil {
				t.Errorf("unexpected error for %s: %v", tt.name, err)
			}
			if count < 0 {
				t.Errorf("negative token count: %d", count)
			}
		})
	}
}
