// Package extractor provides signature extraction from scanned files.
package extractor

import (
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/indigo-net/Brf.it/pkg/parser"
	"github.com/indigo-net/Brf.it/pkg/scanner"
)

// ExtractedFile represents a file with its extracted signatures.
type ExtractedFile struct {
	// Path is the file path.
	Path string

	// Language is the detected language.
	Language string

	// Signatures is the list of extracted signatures.
	Signatures []parser.Signature

	// RawImports is the list of raw import/export statement text.
	RawImports []string

	// Size is the file size in bytes.
	Size int64

	// Error is any error that occurred during extraction.
	Error error
}

// ExtractResult contains the results of an extraction operation.
type ExtractResult struct {
	// Files is the list of extracted files.
	Files []ExtractedFile

	// TotalSignatures is the total number of signatures extracted.
	TotalSignatures int

	// TotalSize is the total size of processed files.
	TotalSize int64

	// ErrorCount is the number of files that had errors.
	ErrorCount int
}

// ExtractOptions configures the extraction behavior.
type ExtractOptions struct {
	// IncludePrivate whether to include non-exported/private signatures.
	IncludePrivate bool

	// IncludeBody whether to include function/method bodies.
	IncludeBody bool

	// IncludeImports whether to include import/export statements.
	IncludeImports bool

	// Concurrency is the number of concurrent workers.
	// 0 = auto (runtime.NumCPU()), 1 = sequential.
	Concurrency int

	// MaxFileSize is the maximum file size in bytes for TOCTOU re-check.
	// If positive, file content size is verified after reading.
	MaxFileSize int64
}

// Extractor defines the interface for signature extraction.
type Extractor interface {
	// Extract extracts signatures from the given scan result.
	Extract(scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error)
}

// FileExtractor implements Extractor using Scanner and Parser Registry.
type FileExtractor struct {
	registry *parser.Registry
}

// NewFileExtractor creates a new FileExtractor with the given registry.
func NewFileExtractor(registry *parser.Registry) *FileExtractor {
	return &FileExtractor{
		registry: registry,
	}
}

// NewDefaultFileExtractor creates a FileExtractor with the default registry.
func NewDefaultFileExtractor() *FileExtractor {
	return NewFileExtractor(parser.DefaultRegistry())
}

// DefaultExtractOptions returns ExtractOptions with sensible defaults.
// Concurrency defaults to 0 (auto = runtime.NumCPU()).
func DefaultExtractOptions() *ExtractOptions {
	return &ExtractOptions{
		Concurrency: 0,
	}
}

// Extract implements Extractor interface.
func (e *FileExtractor) Extract(scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error) {
	if opts == nil {
		opts = &ExtractOptions{}
	}

	// Validate concurrency before accessing scanResult to catch invalid input early.
	if opts.Concurrency < 0 {
		return nil, fmt.Errorf("concurrency must be >= 0, got %d", opts.Concurrency)
	}

	if scanResult == nil {
		return &ExtractResult{}, nil
	}

	files := scanResult.Files
	n := len(files)
	if n == 0 {
		return &ExtractResult{}, nil
	}

	// Resolve concurrency
	concurrency := opts.Concurrency
	if concurrency == 0 {
		concurrency = runtime.NumCPU()
	}
	if concurrency > n {
		concurrency = n
	}

	// Sequential fast path
	if concurrency == 1 {
		return e.extractSequential(files, opts), nil
	}

	// Concurrent path: semaphore + indexed goroutines
	extracted := make([]ExtractedFile, n)
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, fileEntry := range files {
		wg.Add(1)
		sem <- struct{}{} // Acquire
		go func(idx int, entry scanner.FileEntry) {
			defer wg.Done()
			defer func() { <-sem }() // Release
			extracted[idx] = e.extractFile(entry, opts)
		}(i, fileEntry)
	}
	wg.Wait()

	// Aggregate (sequential)
	result := &ExtractResult{Files: extracted}
	for _, ef := range extracted {
		result.TotalSignatures += len(ef.Signatures)
		result.TotalSize += ef.Size
		if ef.Error != nil {
			result.ErrorCount++
		}
	}
	return result, nil
}

// extractSequential processes files sequentially.
func (e *FileExtractor) extractSequential(files []scanner.FileEntry, opts *ExtractOptions) *ExtractResult {
	result := &ExtractResult{}
	for _, fileEntry := range files {
		extracted := e.extractFile(fileEntry, opts)
		result.Files = append(result.Files, extracted)
		result.TotalSignatures += len(extracted.Signatures)
		result.TotalSize += extracted.Size
		if extracted.Error != nil {
			result.ErrorCount++
		}
	}
	return result
}

// extractFile extracts signatures from a single file.
func (e *FileExtractor) extractFile(entry scanner.FileEntry, opts *ExtractOptions) ExtractedFile {
	extracted := ExtractedFile{
		Path:     entry.Path,
		Language: entry.Language,
		Size:     entry.Size,
	}

	// Get parser for language
	p, ok := e.registry.Get(entry.Language)
	if !ok {
		extracted.Error = fmt.Errorf("no parser for language: %s", entry.Language)
		return extracted
	}

	// Read file content
	content, err := os.ReadFile(entry.Path)
	if err != nil {
		extracted.Error = fmt.Errorf("failed to read file: %w", err)
		return extracted
	}

	// TOCTOU guard: re-check file size after reading
	if opts.MaxFileSize > 0 && int64(len(content)) > opts.MaxFileSize {
		extracted.Error = fmt.Errorf("file size changed since scan: %s (%d > %d)",
			entry.Path, len(content), opts.MaxFileSize)
		return extracted
	}

	// Parse content
	parseResult, err := p.Parse(string(content), &parser.Options{
		Language:       entry.Language,
		IncludePrivate: opts.IncludePrivate,
		IncludeBody:    opts.IncludeBody,
		IncludeImports: opts.IncludeImports,
	})
	if err != nil {
		extracted.Error = fmt.Errorf("failed to parse: %w", err)
		return extracted
	}

	extracted.Signatures = parseResult.Signatures
	extracted.RawImports = parseResult.RawImports
	return extracted
}
