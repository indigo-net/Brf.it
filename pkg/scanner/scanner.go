// Package scanner provides file system scanning capabilities for brfit.
package scanner

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
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

	// Warnings contains non-fatal issues encountered during scanning.
	Warnings []string
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
			".go":   "go",
			".ts":   "typescript",
			".tsx":  "typescript",
			".js":   "javascript",
			".jsx":  "javascript",
			".py":   "python",
			".c":    "c",
			".cpp":  "cpp",
			".hpp":  "cpp",
			".h":    "cpp",
			".java": "java",
			".rs":    "rust",
			".swift": "swift",
			".kt":    "kotlin",
			".kts":   "kotlin",
			".cs":    "csharp",
			".lua":   "lua",
			".sh":    "shell",
			".bash":  "shell",
			".zsh":   "shell",
			".php":   "php",
			".rb":    "ruby",
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

// getBaseName extracts the base name from path, handling UNC edge cases.
// Returns empty string for paths where filepath.Base returns "." (empty path),
// which can occur with certain special paths like Windows UNC roots.
func getBaseName(path string) string {
	name := filepath.Base(path)
	if name == "." {
		return ""
	}
	return name
}

// Scanner defines the interface for file system scanning.
type Scanner interface {
	// Scan performs the scan and returns scan results.
	Scan() (*ScanResult, error)
}

// FileScanner implements Scanner for file system traversal.
type FileScanner struct {
	opts             *ScanOptions
	ignorer          *ignore.GitIgnore
	ignorerErr       error
	ignorerErrWarned bool
	logger           *log.Logger
}

// NewFileScanner creates a new FileScanner with the given options.
// If opts is nil, default options are used.
func NewFileScanner(opts *ScanOptions) (*FileScanner, error) {
	if opts == nil {
		opts = DefaultScanOptions()
	}

	scanner := &FileScanner{
		opts:   opts,
		logger: log.New(os.Stderr, "[brfit] ", 0),
	}

	// Try to load gitignore file
	if opts.IgnoreFile != "" {
		ignorer, err := ignore.CompileIgnoreFile(opts.IgnoreFile)
		if err != nil {
			// Default .gitignore not found is normal; only store error for
			// user-specified ignore files or non-file-not-found errors.
			if !(errors.Is(err, os.ErrNotExist) && opts.IgnoreFile == ".gitignore") {
				scanner.ignorerErr = err
			}
		} else {
			scanner.ignorer = ignorer
		}
	}

	return scanner, nil
}

// Scan implements the Scanner interface.
// It recursively traverses the directory tree and returns matching files.
func (s *FileScanner) Scan() (*ScanResult, error) {
	result := &ScanResult{}

	// Warn once if gitignore loading failed
	if s.ignorerErr != nil && !s.ignorerErrWarned {
		s.logger.Printf("WARN: failed to load ignore file: %v", s.ignorerErr)
		s.ignorerErrWarned = true
	}

	// Check if root path is a file
	info, err := os.Stat(s.opts.RootPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		// Single file - check if it matches criteria
		if entry, ok := s.checkFile(s.opts.RootPath, info); ok {
			result.Files = append(result.Files, entry)
			result.TotalSize = entry.Size
		} else {
			result.SkippedCount = 1
		}
		return result, nil
	}

	// Directory - walk recursively using WalkDir (more efficient than Walk)
	err = filepath.WalkDir(s.opts.RootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			var warning string
			switch {
			case os.IsPermission(err):
				warning = fmt.Sprintf("permission denied: %s", path)
			case os.IsNotExist(err):
				warning = fmt.Sprintf("file not found: %s", path)
			default:
				warning = fmt.Sprintf("skipping: %s: %v", path, err)
			}
			s.logger.Printf("WARN: %s", warning)
			result.Warnings = append(result.Warnings, warning)
			return nil
		}

		// Skip symlinks
		if d.Type()&os.ModeSymlink != 0 {
			warning := fmt.Sprintf("skipping symlink: %s", path)
			s.logger.Printf("WARN: %s", warning)
			result.Warnings = append(result.Warnings, warning)
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Handle directories
		if d.IsDir() {
			// Skip hidden directories (e.g., .git, .idea), but not the root directory
			if path != s.opts.RootPath {
				name := getBaseName(path)
				// Edge case: empty name means filepath.Base returned "." (empty/special path)
				// Continue traversal without applying hidden check
				if name == "" {
					return nil
				}
				if !s.opts.IncludeHidden && IsHidden(name) {
					return filepath.SkipDir
				}
				// Check gitignore for directory
				if s.ignorer != nil && s.ignorer.MatchesPath(path) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Get file info for size check
		info, err := d.Info()
		if err != nil {
			warning := fmt.Sprintf("skipping: %s: %v", path, err)
			s.logger.Printf("WARN: %s", warning)
			result.Warnings = append(result.Warnings, warning)
			return nil
		}

		// Check file
		if entry, ok := s.checkFile(path, info); ok {
			result.Files = append(result.Files, entry)
			result.TotalSize += entry.Size
		} else {
			result.SkippedCount++
		}

		return nil
	})

	return result, err
}

// checkFile checks if a file should be included in the scan results.
func (s *FileScanner) checkFile(path string, info os.FileInfo) (FileEntry, bool) {
	// Check hidden
	name := getBaseName(path)
	// UNC root paths return empty - skip hidden check and include the file
	if name == "" {
		// Fall through to other checks
	} else if !s.opts.IncludeHidden && IsHidden(name) {
		return FileEntry{}, false
	}

	// Check gitignore
	if s.ignorer != nil && s.ignorer.MatchesPath(path) {
		return FileEntry{}, false
	}

	// Check extension
	language, ok := s.opts.GetLanguage(path)
	if !ok {
		return FileEntry{}, false
	}

	// Check file size - log warning for large files
	if info.Size() > s.opts.MaxFileSize {
		s.logger.Printf("WARN: skipping large file %s (%d bytes > %d limit)",
			path, info.Size(), s.opts.MaxFileSize)
		return FileEntry{}, false
	}

	return FileEntry{
		Path:     path,
		Language: language,
		Size:     info.Size(),
	}, true
}
