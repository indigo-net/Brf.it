// Package extractor provides signature extraction from scanned files.
package extractor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

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

	// Calls is the list of function call references.
	Calls []parser.FunctionCall

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

	// IncludeCalls whether to include function call references.
	IncludeCalls bool

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
	// The context controls cancellation and timeout for the extraction.
	Extract(ctx context.Context, scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error)
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
func (e *FileExtractor) Extract(ctx context.Context, scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error) {
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

	// Check context before starting work
	if err := ctx.Err(); err != nil {
		return nil, err
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
		return e.extractSequential(ctx, files, opts)
	}

	// Concurrent path: semaphore + indexed goroutines
	extracted := make([]ExtractedFile, n)
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var cancelErr error
	var cancelOnce sync.Once

	for i, fileEntry := range files {
		// Check context before launching each goroutine
		select {
		case <-ctx.Done():
			cancelOnce.Do(func() { cancelErr = ctx.Err() })
		case sem <- struct{}{}: // Acquire
		}
		if cancelErr != nil {
			break
		}

		wg.Add(1)
		go func(idx int, entry scanner.FileEntry) {
			defer wg.Done()
			defer func() { <-sem }() // Release
			extracted[idx] = e.extractFile(ctx, entry, opts)
		}(i, fileEntry)
	}
	// Wait for goroutines with context awareness.
	// If context is cancelled, give in-flight goroutines a grace period
	// to finish before returning (e.g. CGO parser calls cannot be interrupted).
	waitDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitDone)
	}()
	select {
	case <-waitDone:
		// All goroutines completed normally
	case <-ctx.Done():
		cancelOnce.Do(func() { cancelErr = ctx.Err() })
		select {
		case <-waitDone:
			// Goroutines finished within grace period
		case <-time.After(10 * time.Second):
			// In-flight goroutines did not finish; return to avoid blocking forever.
			// Goroutines will eventually complete on their own.
		}
	}

	// On context cancellation, discard partial results. The caller
	// requested cancellation, so incomplete extraction is not useful.
	if cancelErr != nil {
		return nil, cancelErr
	}

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
func (e *FileExtractor) extractSequential(ctx context.Context, files []scanner.FileEntry, opts *ExtractOptions) (*ExtractResult, error) {
	result := &ExtractResult{}
	for _, fileEntry := range files {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		extracted := e.extractFile(ctx, fileEntry, opts)
		result.Files = append(result.Files, extracted)
		result.TotalSignatures += len(extracted.Signatures)
		result.TotalSize += extracted.Size
		if extracted.Error != nil {
			result.ErrorCount++
		}
	}
	return result, nil
}

// binarySniffSize is the number of bytes inspected for NUL to detect binary content.
const binarySniffSize = 512

// isBinaryContent reports whether content appears to be binary by checking
// for a NUL byte in the first 512 bytes.
func isBinaryContent(content []byte) bool {
	sniff := content
	if len(sniff) > binarySniffSize {
		sniff = sniff[:binarySniffSize]
	}
	return bytes.ContainsRune(sniff, 0)
}

// extractFile extracts signatures from a single file.
func (e *FileExtractor) extractFile(ctx context.Context, entry scanner.FileEntry, opts *ExtractOptions) ExtractedFile {
	extracted := ExtractedFile{
		Path:     entry.Path,
		Language: entry.Language,
		Size:     entry.Size,
	}

	// Check context before expensive I/O
	if err := ctx.Err(); err != nil {
		extracted.Error = err
		return extracted
	}

	// Get parser for language
	p, ok := e.registry.Get(entry.Language)
	if !ok {
		langs := e.registry.Languages()
		sort.Strings(langs)
		extracted.Error = fmt.Errorf("no parser for language %q in %q. Available parsers: %s", entry.Language, entry.Path, strings.Join(langs, ", "))
		return extracted
	}

	// Use preloaded content if available, otherwise read from disk
	content := entry.Content
	if content == nil {
		var err error
		content, err = os.ReadFile(entry.Path)
		if err != nil {
			extracted.Error = fmt.Errorf("failed to read file %q: %w", entry.Path, err)
			return extracted
		}
	}

	// Skip binary files
	if isBinaryContent(content) {
		extracted.Error = fmt.Errorf("skipping binary file %q (detected NUL byte in content)", entry.Path)
		return extracted
	}

	// TOCTOU guard: re-check file size after reading
	if opts.MaxFileSize > 0 && int64(len(content)) > opts.MaxFileSize {
		extracted.Error = fmt.Errorf("file size changed since scan: %q (%d > %d bytes limit)",
			entry.Path, len(content), opts.MaxFileSize)
		return extracted
	}

	// Parse content (no string conversion needed)
	parseResult, err := p.Parse(content, &parser.Options{
		Language:       entry.Language,
		IncludePrivate: opts.IncludePrivate,
		IncludeBody:    opts.IncludeBody,
		IncludeImports: opts.IncludeImports,
		IncludeCalls:   opts.IncludeCalls,
	})
	if err != nil {
		extracted.Error = fmt.Errorf("failed to parse %q: %w", entry.Path, err)
		return extracted
	}

	extracted.Signatures = parseResult.Signatures
	extracted.RawImports = parseResult.RawImports
	extracted.Calls = parseResult.Calls
	return extracted
}
