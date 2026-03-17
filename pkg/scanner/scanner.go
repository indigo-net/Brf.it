// Package scanner provides file system scanning capabilities for brfit.
package scanner

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	ignore "github.com/sabhiram/go-gitignore"
	"github.com/indigo-net/Brf.it/pkg/parser"
)

// FileEntry represents a single file discovered during scanning.
type FileEntry struct {
	// Path is the absolute or relative path to the file.
	Path string

	// Language is the detected programming language (e.g., "go", "typescript").
	Language string

	// Size is the file size in bytes.
	Size int64

	// Content holds the file bytes when PreloadContent is enabled.
	// nil when content was not preloaded.
	Content []byte
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

	// IgnoreFiles is the list of ignore file paths (default: [".gitignore"]).
	IgnoreFiles []string

	// IncludePatterns is a list of glob patterns to include.
	// If non-empty, only files matching at least one pattern are included.
	// Supports doublestar (**) patterns.
	IncludePatterns []string

	// ExcludePatterns is a list of glob patterns to exclude.
	// Files matching any pattern are excluded.
	// Supports doublestar (**) patterns.
	ExcludePatterns []string

	// ChangedFiles is an optional whitelist of file paths (relative to RootPath).
	// When non-nil, only files in this list are included in scan results.
	// Used by --changed and --since flags to restrict scanning to git-changed files.
	ChangedFiles map[string]bool

	// IncludeHidden determines whether to include hidden files (dotfiles).
	IncludeHidden bool

	// MaxFileSize is the maximum file size in bytes to include.
	MaxFileSize int64

	// PreloadContent reads file content during scan so downstream consumers
	// (e.g., the extractor) can skip a redundant os.ReadFile call.
	PreloadContent bool

	// MaxTotalPreloadSize limits the total bytes preloaded into memory when
	// PreloadContent is true. Once this budget is exceeded, remaining files
	// are included in the scan results but with Content set to nil (the
	// extractor will fall back to on-demand os.ReadFile). A value of 0 means
	// no limit. Default: 0.
	MaxTotalPreloadSize int64
}

// DefaultScanOptions returns a ScanOptions with sensible defaults.
// SupportedExtensions is derived from parser.LanguageMapping (single source of truth).
func DefaultScanOptions() *ScanOptions {
	return &ScanOptions{
		SupportedExtensions: copyLanguageMapping(),
		IgnoreFiles:   []string{".gitignore"},
		IncludeHidden: false,
		MaxFileSize:   512000, // 500KB
	}
}

// copyLanguageMapping returns a copy of parser.LanguageMapping().
func copyLanguageMapping() map[string]string {
	return parser.LanguageMapping()
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
	// The context allows cancellation of long-running scans.
	Scan(ctx context.Context) (*ScanResult, error)
}

// FileScanner implements Scanner for file system traversal.
type FileScanner struct {
	opts              *ScanOptions
	ignorers          []*ignore.GitIgnore
	ignorerErrs       []error
	ignorerErrsWarned bool
	logger            *log.Logger
	rootIsFile        bool
	preloadedSize     int64 // tracks total bytes preloaded so far
}

// NewFileScanner creates a new FileScanner with the given options.
// If opts is nil, default options are used.
func NewFileScanner(opts *ScanOptions) (*FileScanner, error) {
	if opts == nil {
		opts = DefaultScanOptions()
	}

	// Validate glob patterns at construction time (fail-fast)
	for _, p := range opts.IncludePatterns {
		if !doublestar.ValidatePattern(p) {
			return nil, fmt.Errorf("invalid include pattern %q", p)
		}
	}
	for _, p := range opts.ExcludePatterns {
		if !doublestar.ValidatePattern(p) {
			return nil, fmt.Errorf("invalid exclude pattern %q", p)
		}
	}

	// Cache whether RootPath is a file (not a directory) to avoid
	// repeated os.Stat calls in relPath().
	var rootIsFile bool
	if info, err := os.Stat(opts.RootPath); err == nil && !info.IsDir() {
		rootIsFile = true
	}

	s := &FileScanner{
		opts:       opts,
		logger:     log.New(os.Stderr, "[brfit] ", 0),
		rootIsFile: rootIsFile,
	}

	// Try to load each ignore file
	for _, ignoreFile := range opts.IgnoreFiles {
		if ignoreFile == "" {
			continue
		}
		ignorer, err := ignore.CompileIgnoreFile(ignoreFile)
		if err != nil {
			// Default .gitignore not found is normal; only store error for
			// user-specified ignore files or non-file-not-found errors.
			if !(errors.Is(err, os.ErrNotExist) && ignoreFile == ".gitignore") {
				s.ignorerErrs = append(s.ignorerErrs, fmt.Errorf("%s: %w", ignoreFile, err))
			}
		} else {
			s.ignorers = append(s.ignorers, ignorer)
		}
	}

	return s, nil
}

// Scan implements the Scanner interface.
// It recursively traverses the directory tree and returns matching files.
func (s *FileScanner) Scan(ctx context.Context) (*ScanResult, error) {
	result := &ScanResult{}

	// Warn once if any ignore file loading failed
	if len(s.ignorerErrs) > 0 && !s.ignorerErrsWarned {
		for _, err := range s.ignorerErrs {
			s.logger.Printf("WARN: failed to load ignore file: %v", err)
		}
		s.ignorerErrsWarned = true
	}

	// Check if root path is a file
	info, err := os.Stat(s.opts.RootPath)
	if err != nil {
		return nil, fmt.Errorf("cannot access path %q: %w", s.opts.RootPath, err)
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
		// Check context cancellation
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}

		if err != nil {
			var warning string
			switch {
			case os.IsPermission(err):
				warning = fmt.Sprintf("permission denied: %s (check file permissions or run with appropriate privileges)", path)
			case os.IsNotExist(err):
				warning = fmt.Sprintf("file not found: %s (may have been deleted during scan)", path)
			default:
				warning = fmt.Sprintf("skipping: %s: %v", path, err)
			}
			s.logger.Printf("WARN: %s", warning)
			result.Warnings = append(result.Warnings, warning)
			return nil
		}

		// Skip symlinks
		if d.Type()&os.ModeSymlink != 0 {
			warning := fmt.Sprintf("skipping symlink: %s (symlinks are not followed for security)", path)
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
				// Check gitignore for directory (skip if ANY ignorer matches)
				if s.matchesIgnore(path) {
					return filepath.SkipDir
				}
				// Check exclude patterns for directory
				if s.matchesExcludeDir(s.relPath(path)) {
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

// relPath returns the path relative to the root, using forward slashes.
// When RootPath is a file (not a directory), uses the parent directory as base
// so that the file's name is preserved for glob matching.
func (s *FileScanner) relPath(path string) string {
	base := s.opts.RootPath
	// If RootPath is a file, use its parent as the base for relative paths.
	// This ensures glob matching works (e.g., "**/*.go" matches "main.go").
	if s.rootIsFile {
		base = filepath.Dir(base)
	}
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}

// matchesInclude returns true if the path matches at least one include pattern.
// Returns true if no include patterns are configured (no filtering).
// rel is the pre-computed relative path from relPath().
func (s *FileScanner) matchesInclude(rel string) bool {
	if len(s.opts.IncludePatterns) == 0 {
		return true
	}
	for _, pattern := range s.opts.IncludePatterns {
		// Patterns are validated in NewFileScanner; error is unreachable
		if matched, _ := doublestar.Match(pattern, rel); matched {
			return true
		}
	}
	return false
}

// matchesExclude returns true if the path matches any exclude pattern.
// rel is the pre-computed relative path from relPath().
func (s *FileScanner) matchesExclude(rel string) bool {
	if len(s.opts.ExcludePatterns) == 0 {
		return false
	}
	for _, pattern := range s.opts.ExcludePatterns {
		// Patterns are validated in NewFileScanner; error is unreachable
		if matched, _ := doublestar.Match(pattern, rel); matched {
			return true
		}
	}
	return false
}

// matchesExcludeDir returns true if a directory should be pruned.
// Handles patterns like "vendor/**" by also matching the directory name itself.
// rel is the pre-computed relative path from relPath().
func (s *FileScanner) matchesExcludeDir(rel string) bool {
	if len(s.opts.ExcludePatterns) == 0 {
		return false
	}
	for _, pattern := range s.opts.ExcludePatterns {
		// Patterns are validated in NewFileScanner; error is unreachable
		if matched, _ := doublestar.Match(pattern, rel); matched {
			return true
		}
		// "vendor/**" should also prune the "vendor" directory node
		if stripped := strings.TrimSuffix(pattern, "/**"); stripped != pattern {
			if matched, _ := doublestar.Match(stripped, rel); matched {
				return true
			}
		}
	}
	return false
}

// matchesIgnore returns true if the path matches any of the loaded ignore patterns.
func (s *FileScanner) matchesIgnore(path string) bool {
	for _, ig := range s.ignorers {
		if ig.MatchesPath(path) {
			return true
		}
	}
	return false
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

	// Check gitignore (skip if ANY ignorer matches)
	if s.matchesIgnore(path) {
		return FileEntry{}, false
	}

	// Compute relative path once for all pattern matching
	rel := s.relPath(path)

	// Check exclude patterns
	if s.matchesExclude(rel) {
		return FileEntry{}, false
	}

	// Check include patterns
	if !s.matchesInclude(rel) {
		return FileEntry{}, false
	}

	// Check changed files whitelist
	if s.opts.ChangedFiles != nil {
		if !s.opts.ChangedFiles[rel] {
			return FileEntry{}, false
		}
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

	entry := FileEntry{
		Path:     path,
		Language: language,
		Size:     info.Size(),
	}

	if s.opts.PreloadContent {
		// Skip preloading if total preloaded size would exceed budget.
		budget := s.opts.MaxTotalPreloadSize
		if budget > 0 && s.preloadedSize+info.Size() > budget {
			// Still include the file but without preloaded content;
			// the extractor will fall back to on-demand os.ReadFile.
		} else {
			content, err := os.ReadFile(path)
			if err != nil {
				s.logger.Printf("WARN: failed to preload %s: %v (file will be read on-demand)", path, err)
				// Include the file without content; the extractor will
				// fall back to on-demand os.ReadFile.
			} else {
				entry.Content = content
				s.preloadedSize += int64(len(content))
			}
		}
	}

	return entry, true
}
