// Package extractor provides signature extraction from scanned files.
package extractor

import (
	"fmt"
	"os"

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

	// Imports is the list of extracted import/export statements.
	Imports []parser.ImportExport

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

	// Concurrency is the number of concurrent workers (0 = sequential).
	Concurrency int
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

// Extract implements Extractor interface.
func (e *FileExtractor) Extract(scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error) {
	if opts == nil {
		opts = &ExtractOptions{}
	}

	result := &ExtractResult{}

	// Sequential processing (concurrency can be added later)
	for _, fileEntry := range scanResult.Files {
		extracted := e.extractFile(fileEntry, opts)
		result.Files = append(result.Files, extracted)
		result.TotalSignatures += len(extracted.Signatures)
		result.TotalSize += extracted.Size
		if extracted.Error != nil {
			result.ErrorCount++
		}
	}

	return result, nil
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
	extracted.Imports = parseResult.Imports
	return extracted
}
