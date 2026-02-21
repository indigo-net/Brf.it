// Package tokenizer provides token counting for brfit output.
package tokenizer

// Tokenizer defines the interface for counting tokens in text.
type Tokenizer interface {
	// Count returns the number of tokens in the given text.
	// Returns 0 and error if counting fails.
	Count(text string) (int, error)

	// Name returns the tokenizer name (e.g., "tiktoken-cl100k", "noop").
	Name() string
}

// NoOpTokenizer is a tokenizer that always returns 0.
// Used as default when token counting is disabled or unavailable.
type NoOpTokenizer struct{}

// Ensure NoOpTokenizer implements Tokenizer interface.
var _ Tokenizer = (*NoOpTokenizer)(nil)

// NewNoOpTokenizer creates a new NoOpTokenizer.
func NewNoOpTokenizer() *NoOpTokenizer {
	return &NoOpTokenizer{}
}

// Count returns 0 and nil error (no-op).
func (t *NoOpTokenizer) Count(_ string) (int, error) {
	return 0, nil
}

// Name returns "noop".
func (t *NoOpTokenizer) Name() string {
	return "noop"
}
