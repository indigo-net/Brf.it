// Package formatter provides output formatting for brfit.
package formatter

import (
	"github.com/indigo-net/Brf.it/pkg/parser"
)

// FileData represents a file with its extracted data for formatting.
type FileData struct {
	// Path is the file path.
	Path string

	// Language is the detected language.
	Language string

	// Signatures is the list of extracted signatures.
	Signatures []parser.Signature

	// RawImports is the list of raw import/export statement text.
	RawImports []string

	// Error is any error that occurred during extraction.
	Error error
}

// PackageData contains all data needed for formatting output.
type PackageData struct {
	// RootPath is the root path being packaged.
	RootPath string

	// Version is the brf.it version string.
	Version string

	// Tree is the directory tree string.
	Tree string

	// Files is the list of file data.
	Files []FileData

	// TotalSignatures is the total number of signatures.
	TotalSignatures int

	// TotalSize is the total size of processed files.
	TotalSize int64

	// IncludeImports indicates whether imports should be rendered.
	IncludeImports bool

	// DedupeImports indicates whether imports should be deduplicated across files.
	// When true, imports are collected globally and shown in a separate section.
	DedupeImports bool

	// GlobalImports holds deduplicated imports with their usage counts.
	// Only populated when DedupeImports is true.
	GlobalImports []ImportCount
}

// ImportCount represents an import with its usage count across files.
type ImportCount struct {
	// Import is the raw import statement text.
	Import string
	// Count is the number of files that use this import.
	Count int
}

// Formatter defines the interface for output formatting.
type Formatter interface {
	// Format formats the package data and returns the output bytes.
	Format(data *PackageData) ([]byte, error)

	// Name returns the formatter name (e.g., "xml", "markdown").
	Name() string
}
