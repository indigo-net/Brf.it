// Package scanner provides file system scanning capabilities for brfit.
package scanner

import (
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// FileEntry represents a single file discovered during scanning.
type FileEntry struct {
	// Path is the absolute or relative path to the file.
	Path string

	// Language is the detected programming language (e.g., "go", "typescript").
	Language string

	// Size is the file size in bytes.
	Size int64
}

// ScanResult contains the results of a scan operation.
type ScanResult struct {
	// Files is the list of matched files.
	Files []FileEntry

	// TotalSize is the sum of all matched file sizes.
	TotalSize int64

	// SkippedCount is the number of files skipped (too large, unsupported, etc.).
	SkippedCount int
}

// ScanOptions configures the scanning behavior.
type ScanOptions struct {
	// RootPath is the directory or file to scan.
	RootPath string

	// SupportedExtensions maps file extensions to language names.
	SupportedExtensions map[string]string

	// IgnoreFile is the path to the gitignore file (default: .gitignore).
	IgnoreFile string

	// IncludeHidden determines whether to include hidden files (dotfiles).
	IncludeHidden bool

	// MaxFileSize is the maximum file size in bytes to include.
	MaxFileSize int64
}

// DefaultScanOptions returns a ScanOptions with sensible defaults.
func DefaultScanOptions() *ScanOptions {
	return &ScanOptions{
		SupportedExtensions: map[string]string{
			".go":  "go",
			".ts":  "typescript",
			".tsx": "typescript",
			".js":  "javascript",
			".jsx": "javascript",
		},
		IgnoreFile:    ".gitignore",
		IncludeHidden: false,
		MaxFileSize:   512000, // 500KB
	}
}

// GetLanguage returns the language for a given file path and whether it's supported.
func (o *ScanOptions) GetLanguage(path string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(path))
	lang, ok := o.SupportedExtensions[ext]
	return lang, ok
}

// IsHidden checks if a file or directory name is hidden (starts with dot).
func IsHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

// Scanner defines the interface for file system scanning.
type Scanner interface {
	// Scan performs the scan and returns scan results.
	Scan() (*ScanResult, error)
}

// FileScanner implements Scanner for file system traversal.
type FileScanner struct {
	opts       *ScanOptions
	ignorer    *ignore.GitIgnore
	ignorerErr error
}

// NewFileScanner creates a new FileScanner with the given options.
// If opts is nil, default options are used.
func NewFileScanner(opts *ScanOptions) (*FileScanner, error) {
	if opts == nil {
		opts = DefaultScanOptions()
	}

	scanner := &FileScanner{
		opts: opts,
	}

	// Try to load gitignore file
	if opts.IgnoreFile != "" {
		ignorer, err := ignore.CompileIgnoreFile(opts.IgnoreFile)
		if err != nil {
			// Store error but don't fail - gitignore is optional
			scanner.ignorerErr = err
		} else {
			scanner.ignorer = ignorer
		}
	}

	return scanner, nil
}
