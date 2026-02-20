package main

import (
	"fmt"
	"os"

	"github.com/indigo-net/Brf.it/internal/config"
	"github.com/indigo-net/Brf.it/pkg/scanner"
	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "0.1.0"

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
		fmt.Printf("brfit version %s\n", Version)
		return nil
	}

	// Parse path argument
	c.Path = "."
	if len(args) > 0 {
		c.Path = args[0]
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

	// Create scanner
	fileScanner, err := scanner.NewFileScanner(scanOpts)
	if err != nil {
		return fmt.Errorf("failed to create scanner: %w", err)
	}

	// Perform scan
	result, err := fileScanner.Scan()
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	// Output results (placeholder - Phase 7 will add XML/MD formatting)
	fmt.Printf("Scanned: %s\n", c.Path)
	fmt.Printf("Found %d file(s) (%d bytes total)", len(result.Files), result.TotalSize)
	if result.SkippedCount > 0 {
		fmt.Printf(", %d skipped", result.SkippedCount)
	}
	fmt.Println()

	for _, entry := range result.Files {
		fmt.Printf("  - %s [%s] (%d bytes)\n", entry.Path, entry.Language, entry.Size)
	}

	return nil
}
