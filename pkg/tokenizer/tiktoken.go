package tokenizer

import (
	"github.com/pkoukk/tiktoken-go"
)

// TiktokenTokenizer implements Tokenizer using tiktoken with cl100k_base encoding.
// This encoding is compatible with GPT-4 and GPT-3.5-turbo models.
type TiktokenTokenizer struct {
	encoding string
	tke      *tiktoken.Tiktoken
}

// Ensure TiktokenTokenizer implements Tokenizer interface.
var _ Tokenizer = (*TiktokenTokenizer)(nil)

// NewTiktokenTokenizer creates a new TiktokenTokenizer with cl100k_base encoding.
// Returns error if the encoding fails to load.
func NewTiktokenTokenizer() (*TiktokenTokenizer, error) {
	tke, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return nil, err
	}

	return &TiktokenTokenizer{
		encoding: "cl100k_base",
		tke:      tke,
	}, nil
}

// Count returns the number of tokens in the given text.
func (t *TiktokenTokenizer) Count(text string) (int, error) {
	if text == "" {
		return 0, nil
	}

	tokens := t.tke.Encode(text, nil, nil)
	return len(tokens), nil
}

// Name returns the tokenizer name with encoding.
func (t *TiktokenTokenizer) Name() string {
	return "tiktoken-" + t.encoding
}
