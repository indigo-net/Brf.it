// Package context provides global configuration and packaging logic for brfit.
package context

import (
	"context"
	"sort"

	"github.com/indigo-net/Brf.it/pkg/extractor"
	"github.com/indigo-net/Brf.it/pkg/formatter"
	"github.com/indigo-net/Brf.it/pkg/scanner"
	"github.com/indigo-net/Brf.it/pkg/tokenizer"
)

// Options contains CLI options for packaging.
type Options struct {
	// Path is the target path to scan.
	Path string

	// Version is the brf.it version string.
	Version string

	// Format is the output format ("xml" or "md").
	Format string

	// Output is the output file path (empty = stdout).
	Output string

	// IgnoreFile is the custom ignore file path.
	IgnoreFile string

	// IncludeHidden determines whether to include hidden files.
	IncludeHidden bool

	// IncludeBody determines whether to include function/method bodies.
	IncludeBody bool

	// IncludeImports determines whether to include import/export statements.
	IncludeImports bool

	// DedupeImports deduplicates imports across files and shows them globally.
	// Requires IncludeImports to be true.
	DedupeImports bool

	// IncludeTree determines whether to include directory tree.
	IncludeTree bool

	// IncludePrivate determines whether to include private symbols.
	IncludePrivate bool

	// MaxFileSize is the maximum file size in bytes.
	MaxFileSize int64

	// MaxDocLength is the maximum length of documentation comments.
	// 0 means no limit (default).
	MaxDocLength int

	// NoSchema skips the schema section in XML output.
	NoSchema bool
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
		IncludeBody:    opts.IncludeBody,
		IncludeImports: opts.IncludeImports,
		MaxFileSize:    opts.MaxFileSize,
	}
	// TODO: propagate context from Package() caller once Package accepts context.Context
	extractResult, err := p.extractor.Extract(context.TODO(), scanResult, extractOpts)
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
			RawImports: ef.RawImports,
			Error:      ef.Error,
		}
	}

	// 4.5 Build global imports if DedupeImports is enabled
	var globalImports []formatter.ImportCount
	if opts.IncludeImports && opts.DedupeImports {
		globalImports = buildGlobalImports(files)
	}

	// 5. Create PackageData
	packageData := &formatter.PackageData{
		RootPath:        opts.Path,
		Version:         opts.Version,
		Tree:            treeStr,
		Files:           files,
		TotalSignatures: extractResult.TotalSignatures,
		TotalSize:       extractResult.TotalSize,
		IncludeImports:  opts.IncludeImports,
		DedupeImports:   opts.DedupeImports,
		GlobalImports:   globalImports,
		MaxDocLength:    opts.MaxDocLength,
		NoSchema:        opts.NoSchema,
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
		"json":     formatter.NewJSONFormatter(),
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

// buildGlobalImports collects and deduplicates imports from all files.
// Returns a list of unique imports with their usage counts, sorted by count (descending).
func buildGlobalImports(files []formatter.FileData) []formatter.ImportCount {
	importCounts := make(map[string]int)

	seen := make(map[string]bool)
	for _, file := range files {
		if file.Error != nil {
			continue
		}
		// Use a set to count each import only once per file
		clear(seen)
		for _, imp := range file.RawImports {
			if !seen[imp] {
				seen[imp] = true
				importCounts[imp]++
			}
		}
	}

	// Convert to slice and sort by count (descending)
	result := make([]formatter.ImportCount, 0, len(importCounts))
	for imp, count := range importCounts {
		result = append(result, formatter.ImportCount{
			Import: imp,
			Count:  count,
		})
	}

	// Sort by count descending, then by import string for stability
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Import < result[j].Import
	})

	return result
}
