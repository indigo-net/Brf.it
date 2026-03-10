// Package config provides CLI configuration management for brfit.
package config

import (
	"errors"
	"fmt"
	"os"

	pkgcontext "github.com/indigo-net/Brf.it/internal/context"
)

// MaxFileSizeUpperBound is the maximum allowed value for MaxFileSize (10MB).
// Values above this threshold trigger a warning (not an error).
const MaxFileSizeUpperBound = 10 * 1024 * 1024

// Config holds all configuration options for the brfit CLI.
type Config struct {
	// Path is the root directory or file to process.
	Path string

	// Version is the brf.it version string.
	Version string

	// Mode determines what to extract. Currently only "sig" (signature) is supported.
	Mode string

	// Format specifies the output format: "xml" or "md".
	Format string

	// Output is the file path to write output. Empty means stdout.
	Output string

	// IgnoreFile is the path to the ignore file (default: .gitignore).
	IgnoreFile string

	// IncludeHidden determines whether to include hidden files (dotfiles).
	IncludeHidden bool

	// IncludeBody determines whether to include function/method bodies.
	// When false (default), only signatures are extracted.
	IncludeBody bool

	// IncludeImports determines whether to include import/export statements.
	IncludeImports bool

	// DedupeImports deduplicates imports across files and shows them globally.
	// Requires IncludeImports to be true.
	DedupeImports bool

	// NoTree skips directory tree generation in output.
	NoTree bool

	// NoTokens disables token count calculation.
	NoTokens bool

	// NoSchema skips the schema section in XML output.
	NoSchema bool

	// MaxFileSize is the maximum file size in bytes to process.
	MaxFileSize int64

	// MaxDocLength is the maximum length of documentation comments in characters.
	// 0 means no limit (default).
	MaxDocLength int
}

// DefaultConfig returns a Config with all default values set.
func DefaultConfig() *Config {
	return &Config{
		Mode:           "sig",
		Format:         "xml",
		Output:         "",
		IgnoreFile:     ".gitignore",
		IncludeHidden:  false,
		IncludeBody:    false,
		IncludeImports: false,
		NoTree:         false,
		NoTokens:       false,
		MaxFileSize:    512000, // 500KB
		MaxDocLength:   0,      // no limit
	}
}

// Validate checks if the configuration values are valid.
// Returns an error describing the first invalid field, or nil if valid.
func (c *Config) Validate() error {
	// Validate mode
	if c.Mode != "sig" {
		return fmt.Errorf("invalid mode '%s': only 'sig' mode is supported", c.Mode)
	}

	// Validate format (accept "xml", "md", "markdown", "json")
	validFormats := map[string]bool{"xml": true, "md": true, "markdown": true, "json": true}
	if !validFormats[c.Format] {
		return fmt.Errorf("invalid format '%s': must be 'xml', 'md', 'markdown', or 'json'", c.Format)
	}

	// Validate max file size
	if c.MaxFileSize <= 0 {
		return errors.New("max file size must be positive")
	}

	// Warn if max file size exceeds upper bound (not an error)
	if c.MaxFileSize > MaxFileSizeUpperBound {
		fmt.Fprintf(os.Stderr, "[brfit] WARN: max-size %d bytes exceeds recommended upper bound of %d bytes (10MB)\n",
			c.MaxFileSize, MaxFileSizeUpperBound)
	}

	return nil
}

// SupportedExtensions returns a map of file extensions to language names.
func (c *Config) SupportedExtensions() map[string]string {
	return map[string]string{
		".go":    "go",
		".ts":    "typescript",
		".tsx":   "typescript",
		".js":    "javascript",
		".jsx":   "javascript",
		".py":    "python",
		".c":     "c",
		".cpp":   "cpp",
		".hpp":   "cpp",
		".h":     "cpp",
		".java":  "java",
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
		".scala": "scala",
		".sc":    "scala",
		".ex":    "elixir",
		".exs":   "elixir",
		".sql":   "sql",
	}
}

// ToOptions converts Config to packager Options.
func (c *Config) ToOptions() *pkgcontext.Options {
	return &pkgcontext.Options{
		Path:           c.Path,
		Version:        c.Version,
		Format:         c.Format,
		Output:         c.Output,
		IgnoreFile:     c.IgnoreFile,
		IncludeHidden:  c.IncludeHidden,
		IncludeBody:    c.IncludeBody,
		IncludeImports: c.IncludeImports,
		DedupeImports:  c.DedupeImports,
		IncludeTree:    !c.NoTree,
		IncludePrivate: false, // Future: add --include-private flag
		MaxFileSize:    c.MaxFileSize,
		MaxDocLength:   c.MaxDocLength,
		NoSchema:       c.NoSchema,
	}
}
