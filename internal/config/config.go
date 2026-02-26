// Package config provides CLI configuration management for brfit.
package config

import (
	"errors"
	"fmt"

	pkgcontext "github.com/indigo-net/Brf.it/internal/context"
)

// Config holds all configuration options for the brfit CLI.
type Config struct {
	// Path is the root directory or file to process.
	Path string

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

	// NoTree skips directory tree generation in output.
	NoTree bool

	// NoTokens disables token count calculation.
	NoTokens bool

	// MaxFileSize is the maximum file size in bytes to process.
	MaxFileSize int64
}

// DefaultConfig returns a Config with all default values set.
func DefaultConfig() *Config {
	return &Config{
		Mode:          "sig",
		Format:        "xml",
		Output:        "",
		IgnoreFile:    ".gitignore",
		IncludeHidden: false,
		IncludeBody:   false,
		NoTree:        false,
		NoTokens:      false,
		MaxFileSize:   512000, // 500KB
	}
}

// Validate checks if the configuration values are valid.
// Returns an error describing the first invalid field, or nil if valid.
func (c *Config) Validate() error {
	// Validate mode
	if c.Mode != "sig" {
		return fmt.Errorf("invalid mode '%s': only 'sig' mode is supported", c.Mode)
	}

	// Validate format (accept "xml", "md", "markdown")
	if c.Format != "xml" && c.Format != "md" && c.Format != "markdown" {
		return fmt.Errorf("invalid format '%s': must be 'xml', 'md', or 'markdown'", c.Format)
	}

	// Validate max file size
	if c.MaxFileSize <= 0 {
		return errors.New("max file size must be positive")
	}

	return nil
}

// SupportedExtensions returns a map of file extensions to language names.
func (c *Config) SupportedExtensions() map[string]string {
	return map[string]string{
		".go":  "go",
		".ts":  "typescript",
		".tsx": "typescript",
		".js":  "javascript",
		".jsx": "javascript",
		".py":  "python",
		".c":   "c",
		".h":   "c",
	}
}

// ToOptions converts Config to packager Options.
func (c *Config) ToOptions() *pkgcontext.Options {
	return &pkgcontext.Options{
		Path:           c.Path,
		Format:         c.Format,
		Output:         c.Output,
		IgnoreFile:     c.IgnoreFile,
		IncludeHidden:  c.IncludeHidden,
		IncludeBody:    c.IncludeBody,
		IncludeTree:    !c.NoTree,
		IncludePrivate: false, // Future: add --include-private flag
		MaxFileSize:    c.MaxFileSize,
	}
}
