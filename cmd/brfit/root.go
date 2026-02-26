package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/indigo-net/Brf.it/internal/config"
	"github.com/indigo-net/Brf.it/internal/context"
	"github.com/indigo-net/Brf.it/pkg/scanner"
	"github.com/spf13/cobra"

	// Import treesitter parser to register Go/TypeScript parsers
	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"
)

// Build information (set by main.go from ldflags)
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// SetBuildInfo sets build information from ldflags (called by main.go)
func SetBuildInfo(v, c, d string) {
	Version = v
	Commit = c
	Date = d
}

// cfg holds the global configuration for the CLI.
var cfg *config.Config

// rootCmd represents the base command when called without any subcommands.
var rootCmd *cobra.Command

func init() {
	cfg = config.DefaultConfig()
	rootCmd = newRootCommandWithConfig(cfg)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// NewRootCommand creates and returns a new root command instance with default config.
// This is useful for testing command structure without executing.
func NewRootCommand() *cobra.Command {
	return newRootCommandWithConfig(config.DefaultConfig())
}

// newRootCommandWithConfig creates a root command bound to the given config.
// This allows tests to use isolated config instances.
func newRootCommandWithConfig(c *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "brfit [path] [options]",
		Short: "Brief your code for AI assistants",
		Long: `Brf.it extracts function signatures and documentation from your codebase,
transforming them into a format optimized for AI coding assistants.

By removing implementation details, it significantly reduces token usage while
preserving the essential information AI needs to understand your project.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRoot(cmd, args, c)
		},
	}

	// Add flags bound to the provided config
	addFlags(cmd, c)

	return cmd
}

// addFlags adds all CLI flags to the given command, bound to the provided config.
func addFlags(cmd *cobra.Command, c *config.Config) {
	// Mode flag
	cmd.Flags().StringVarP(&c.Mode, "mode", "m", c.Mode,
		"output mode: \"sig\" (signature only)")

	// Format flag
	cmd.Flags().StringVarP(&c.Format, "format", "f", c.Format,
		"output format: \"xml\" | \"md\"")

	// Output flag
	cmd.Flags().StringVarP(&c.Output, "output", "o", c.Output,
		"output file path (default: stdout)")

	// Ignore file flag
	cmd.Flags().StringVarP(&c.IgnoreFile, "ignore", "i", c.IgnoreFile,
		"custom ignore file (default: .gitignore)")

	// Boolean flags
	cmd.Flags().BoolVar(&c.IncludeHidden, "include-hidden", c.IncludeHidden,
		"include hidden files (dotfiles)")

	cmd.Flags().BoolVar(&c.IncludeBody, "include-body", c.IncludeBody,
		"include function/method bodies (default: signatures only)")

	cmd.Flags().BoolVar(&c.IncludeImports, "include-imports", c.IncludeImports,
		"include import/export statements in output")

	cmd.Flags().BoolVar(&c.NoTree, "no-tree", c.NoTree,
		"skip directory tree in output")

	cmd.Flags().BoolVar(&c.NoTokens, "no-tokens", c.NoTokens,
		"disable token count calculation")

	// Max file size
	cmd.Flags().Int64Var(&c.MaxFileSize, "max-size", c.MaxFileSize,
		"maximum file size in bytes (default: 512000 = 500KB)")

	// Version flag
	cmd.Flags().BoolP("version", "v", false, "print version information")
}

// runRoot is the main execution function for the root command.
func runRoot(cmd *cobra.Command, args []string, c *config.Config) error {
	// Check if version flag is set
	versionFlag, _ := cmd.Flags().GetBool("version")
	if versionFlag {
		fmt.Printf("brfit %s (commit: %s, built: %s)\n", Version, Commit, Date)
		return nil
	}

	// Parse path argument
	c.Path = "."
	if len(args) > 0 {
		c.Path = args[0]
	}

	// Validate path exists
	if _, err := os.Stat(c.Path); os.IsNotExist(err) {
		return fmt.Errorf("path not found: %s", c.Path)
	}

	// Validate configuration
	if err := c.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Create scan options from config
	scanOpts := &scanner.ScanOptions{
		RootPath:            c.Path,
		SupportedExtensions: c.SupportedExtensions(),
		IgnoreFile:          c.IgnoreFile,
		IncludeHidden:       c.IncludeHidden,
		MaxFileSize:         c.MaxFileSize,
	}

	// Create packager with default dependencies
	packager, err := context.NewDefaultPackager(scanOpts)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	// Disable tokenizer if --no-tokens flag is set
	if c.NoTokens {
		packager.SetTokenizer(nil)
	}

	// Convert config to options
	opts := c.ToOptions()

	// Execute packaging
	result, err := packager.Package(opts)
	if err != nil {
		return fmt.Errorf("processing failed: %w", err)
	}

	// Write output
	if err := writeOutput(result, c); err != nil {
		return fmt.Errorf("output failed: %w", err)
	}

	// Print summary to stderr (doesn't pollute stdout)
	fmt.Fprintf(os.Stderr, "Files: %d, Signatures: %d",
		result.TotalFiles, result.TotalSignatures)
	if result.TokenCount > 0 {
		fmt.Fprintf(os.Stderr, ", Tokens: %d", result.TokenCount)
	}
	fmt.Fprintln(os.Stderr)

	return nil
}

// writeOutput writes the result to stdout or a file.
func writeOutput(result *context.Result, c *config.Config) error {
	if c.Output == "" {
		// Write to stdout (direct []byte output for efficiency)
		_, err := os.Stdout.Write(result.Content)
		return err
	}

	// Write to file
	return writeToFile(c.Output, result.Content)
}

// writeToFile writes content to a file, creating parent directories if needed.
func writeToFile(path string, content []byte) error {
	// Create parent directory if needed
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	return os.WriteFile(path, content, 0644)
}
