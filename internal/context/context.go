// Package context provides global configuration and packaging logic for brfit.
package context

import (
	"github.com/indigo-net/Brf.it/pkg/extractor"
	"github.com/indigo-net/Brf.it/pkg/formatter"
	"github.com/indigo-net/Brf.it/pkg/scanner"
	"github.com/indigo-net/Brf.it/pkg/tokenizer"
)

// Options contains CLI options for packaging.
type Options struct {
	// Path is the target path to scan.
	Path string

	// Format is the output format ("xml" or "md").
	Format string

	// Output is the output file path (empty = stdout).
	Output string

	// IgnoreFile is the custom ignore file path.
	IgnoreFile string

	// IncludeHidden determines whether to include hidden files.
	IncludeHidden bool

	// IncludeTree determines whether to include directory tree.
	IncludeTree bool

	// IncludePrivate determines whether to include private symbols.
	IncludePrivate bool

	// MaxFileSize is the maximum file size in bytes.
	MaxFileSize int64
}

// DefaultOptions returns Options with sensible defaults.
func DefaultOptions() *Options {
	return &Options{
		Format:        "xml",
		IgnoreFile:    ".gitignore",
		IncludeTree:   true,
		IncludeHidden: false,
		MaxFileSize:   512000, // 500KB
	}
}

// Result contains the final packaged output.
type Result struct {
	// Content is the formatted output bytes.
	Content []byte

	// TotalSignatures is the total number of signatures.
	TotalSignatures int

	// TotalFiles is the number of processed files.
	TotalFiles int

	// TotalSize is the total size of processed files.
	TotalSize int64

	// TokenCount is the number of tokens in the output.
	// Returns 0 if token counting is disabled or tokenizer is not set.
	TokenCount int
}

// Packager orchestrates scanning, extraction, and formatting.
type Packager struct {
	scanner    scanner.Scanner
	extractor  extractor.Extractor
	formatters map[string]formatter.Formatter
	tokenizer  tokenizer.Tokenizer
}

// NewPackager creates a new Packager with the given dependencies.
// Tokenizer is set to NoOpTokenizer by default; use SetTokenizer to change.
func NewPackager(
	s scanner.Scanner,
	e extractor.Extractor,
	f map[string]formatter.Formatter,
) *Packager {
	return &Packager{
		scanner:    s,
		extractor:  e,
		formatters: f,
		tokenizer:  tokenizer.NewNoOpTokenizer(),
	}
}

// SetTokenizer sets the tokenizer for the packager.
// Pass nil to disable token counting (uses NoOpTokenizer).
func (p *Packager) SetTokenizer(t tokenizer.Tokenizer) {
	if t == nil {
		p.tokenizer = tokenizer.NewNoOpTokenizer()
	} else {
		p.tokenizer = t
	}
}

// Package processes files and returns formatted output.
func (p *Packager) Package(opts *Options) (*Result, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	// 1. Scan files
	scanResult, err := p.scanner.Scan()
	if err != nil {
		return nil, err
	}

	// 2. Extract signatures
	extractOpts := &extractor.ExtractOptions{
		IncludePrivate: opts.IncludePrivate,
	}
	extractResult, err := p.extractor.Extract(scanResult, extractOpts)
	if err != nil {
		return nil, err
	}

	// 3. Build directory tree
	var treeStr string
	if opts.IncludeTree && len(scanResult.Files) > 0 {
		paths := make([]string, len(scanResult.Files))
		for i, f := range scanResult.Files {
			paths[i] = f.Path
		}
		treeStr = BuildTree(opts.Path, paths)
	}

	// 4. Convert ExtractedFile to FileData
	files := make([]formatter.FileData, len(extractResult.Files))
	for i, ef := range extractResult.Files {
		files[i] = formatter.FileData{
			Path:       ef.Path,
			Language:   ef.Language,
			Signatures: ef.Signatures,
			Error:      ef.Error,
		}
	}

	// 5. Create PackageData
	packageData := &formatter.PackageData{
		Tree:            treeStr,
		Files:           files,
		TotalSignatures: extractResult.TotalSignatures,
		TotalSize:       extractResult.TotalSize,
	}

	// 6. Get formatter (normalize format and fallback to xml)
	format := normalizeFormat(opts.Format)
	f, ok := p.formatters[format]
	if !ok {
		f = p.formatters["xml"]
	}

	// 7. Format output
	content, err := f.Format(packageData)
	if err != nil {
		return nil, err
	}

	// 8. Calculate token count
	tokenCount, _ := p.tokenizer.Count(string(content))

	return &Result{
		Content:         content,
		TotalSignatures: extractResult.TotalSignatures,
		TotalFiles:      len(extractResult.Files),
		TotalSize:       extractResult.TotalSize,
		TokenCount:      tokenCount,
	}, nil
}

// NewDefaultPackager creates a Packager with default dependencies.
// Tokenizer is set to TiktokenTokenizer if available, otherwise NoOpTokenizer.
func NewDefaultPackager(scanOpts *scanner.ScanOptions) (*Packager, error) {
	s, err := scanner.NewFileScanner(scanOpts)
	if err != nil {
		return nil, err
	}

	e := extractor.NewDefaultFileExtractor()

	formatters := map[string]formatter.Formatter{
		"xml":      formatter.NewXMLFormatter(),
		"markdown": formatter.NewMarkdownFormatter(),
	}

	p := NewPackager(s, e, formatters)

	// Try to set up tiktoken tokenizer (graceful fallback to NoOp)
	if tt, err := tokenizer.NewTiktokenTokenizer(); err == nil {
		p.SetTokenizer(tt)
	}

	return p, nil
}

// normalizeFormat converts CLI format flag to formatter key.
func normalizeFormat(format string) string {
	switch format {
	case "md":
		return "markdown"
	default:
		return format
	}
}
