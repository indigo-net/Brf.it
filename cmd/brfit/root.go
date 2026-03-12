package main

import (
	gocontext "context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/indigo-net/Brf.it/internal/config"
	"github.com/indigo-net/Brf.it/internal/context"
	"github.com/indigo-net/Brf.it/pkg/scanner"
	"github.com/indigo-net/Brf.it/pkg/tokenizer"
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
		"output format: \"xml\" | \"md\" | \"json\"")

	// Output flag
	cmd.Flags().StringVarP(&c.Output, "output", "o", c.Output,
		"output file path (default: stdout)")

	// Ignore file flag (supports multiple: -i .gitignore -i .myignore)
	cmd.Flags().StringArrayVarP(&c.IgnoreFiles, "ignore", "i", c.IgnoreFiles,
		"custom ignore file(s), can be specified multiple times (default: .gitignore)")

	// Glob pattern filters
	cmd.Flags().StringArrayVar(&c.IncludePatterns, "include", c.IncludePatterns,
		"glob pattern(s) to include, can be specified multiple times (e.g., \"pkg/**/*.go\")")

	cmd.Flags().StringArrayVar(&c.ExcludePatterns, "exclude", c.ExcludePatterns,
		"glob pattern(s) to exclude, can be specified multiple times (e.g., \"**/*_test.go\")")

	// Boolean flags
	cmd.Flags().BoolVar(&c.IncludeHidden, "include-hidden", c.IncludeHidden,
		"include hidden files (dotfiles)")

	cmd.Flags().BoolVar(&c.IncludeBody, "include-body", c.IncludeBody,
		"include function/method bodies (default: signatures only)")

	cmd.Flags().BoolVar(&c.IncludeImports, "include-imports", c.IncludeImports,
		"include import/export statements in output")

	cmd.Flags().BoolVar(&c.IncludePrivate, "include-private", c.IncludePrivate,
		"include non-exported/private symbols in output")

	cmd.Flags().BoolVar(&c.DedupeImports, "dedupe-imports", c.DedupeImports,
		"deduplicate imports across files (requires --include-imports)")

	cmd.Flags().BoolVar(&c.NoTree, "no-tree", c.NoTree,
		"skip directory tree in output")

	cmd.Flags().BoolVar(&c.NoTokens, "no-tokens", c.NoTokens,
		"disable token count calculation")

	// Max file size
	cmd.Flags().Int64Var(&c.MaxFileSize, "max-size", c.MaxFileSize,
		"maximum file size in bytes (default: 512000 = 500KB)")

	// Max doc length
	cmd.Flags().IntVar(&c.MaxDocLength, "max-doc-length", c.MaxDocLength,
		"maximum documentation comment length in characters (0 = no limit)")

	// Token tree flag
	cmd.Flags().BoolVar(&c.TokenTree, "token-tree", c.TokenTree,
		"output directory tree with per-file token counts")

	// No schema flag
	cmd.Flags().BoolVar(&c.NoSchema, "no-schema", c.NoSchema,
		"skip XML schema section in output")

	// Git change detection flags
	cmd.Flags().BoolVar(&c.Changed, "changed", c.Changed,
		"only scan files changed in git working tree (git diff --name-only HEAD)")

	cmd.Flags().StringVar(&c.Since, "since", c.Since,
		"only scan files changed since the specified commit/tag (e.g., \"v1.0.0\", \"HEAD~5\")")

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

	// Convert to absolute path for display
	absPath, err := filepath.Abs(c.Path)
	if err == nil {
		c.Path = absPath
	}

	// Set version
	c.Version = Version

	// Validate configuration
	if err := c.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Resolve git changed files if --changed or --since is specified
	var changedFiles map[string]bool
	if c.Changed || c.Since != "" {
		files, err := resolveChangedFiles(c.Path, c.Changed, c.Since)
		if err != nil {
			return fmt.Errorf("failed to resolve changed files: %w", err)
		}
		changedFiles = files
	}

	// Create scan options from config
	scanOpts := &scanner.ScanOptions{
		RootPath:            c.Path,
		SupportedExtensions: c.SupportedExtensions(),
		IgnoreFiles:         c.IgnoreFiles,
		IncludePatterns:     c.IncludePatterns,
		ExcludePatterns:     c.ExcludePatterns,
		ChangedFiles:        changedFiles,
		IncludeHidden:       c.IncludeHidden,
		MaxFileSize:         c.MaxFileSize,
	}

	// Handle --token-tree mode: scan files, count tokens per file, output tree
	if c.TokenTree {
		return runTokenTree(cmd.Context(), scanOpts, c.Path)
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
	result, err := packager.Package(cmd.Context(), opts)
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

// runTokenTree scans files, counts tokens per file, and outputs a directory tree
// with per-file token counts. This is a standalone mode that exits after output.
func runTokenTree(ctx gocontext.Context, scanOpts *scanner.ScanOptions, rootPath string) error {
	// Create scanner
	s, err := scanner.NewFileScanner(scanOpts)
	if err != nil {
		return fmt.Errorf("failed to initialize scanner: %w", err)
	}

	// Scan files
	scanResult, err := s.Scan(ctx)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if len(scanResult.Files) == 0 {
		fmt.Fprintln(os.Stderr, "No files found.")
		return nil
	}

	// Initialize tokenizer
	tok, err := tokenizer.NewTiktokenTokenizer()
	if err != nil {
		return fmt.Errorf("failed to initialize tokenizer: %w", err)
	}

	// Count tokens for each file
	fileCounts := make([]context.FileTokenCount, 0, len(scanResult.Files))
	for _, f := range scanResult.Files {
		content, err := os.ReadFile(f.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[brfit] WARN: cannot read %s: %v\n", f.Path, err)
			continue
		}
		count, err := tok.Count(content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[brfit] WARN: token count failed for %s: %v\n", f.Path, err)
			continue
		}
		fileCounts = append(fileCounts, context.FileTokenCount{
			Path:   f.Path,
			Tokens: count,
		})
	}

	// Build and output the token tree
	tree := context.BuildTokenTree(rootPath, fileCounts)
	fmt.Println(tree)

	return nil
}

// resolveChangedFiles runs git diff to get changed file paths and returns them
// as a set of relative paths (forward-slash separated, relative to rootPath).
// The resulting map is suitable for ScanOptions.ChangedFiles.
//
// When changed is true (--changed flag), untracked files are also included via
// git ls-files --others --exclude-standard.
func resolveChangedFiles(rootPath string, changed bool, since string) (map[string]bool, error) {
	// Determine directory to run git in
	dir := rootPath
	if info, err := os.Stat(rootPath); err == nil && !info.IsDir() {
		dir = filepath.Dir(rootPath)
	}

	// Get the git repository root
	topCmd := exec.Command("git", "rev-parse", "--show-toplevel")
	topCmd.Dir = dir
	topOut, err := topCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git not found or not a git repository: %w", err)
	}
	repoRoot := strings.TrimSpace(string(topOut))
	// Resolve symlinks to avoid mismatches (e.g., macOS /var -> /private/var)
	if resolved, err := filepath.EvalSymlinks(repoRoot); err == nil {
		repoRoot = resolved
	}

	// Determine which git diff command to run
	var diffArgs []string
	if since != "" {
		// --since <ref>: files changed between ref and HEAD + working tree
		diffArgs = []string{"diff", "--name-only", since}
	} else {
		// --changed: uncommitted changes (staged + unstaged)
		diffArgs = []string{"diff", "--name-only", "HEAD"}
	}

	diffCmd := exec.Command("git", diffArgs...)
	diffCmd.Dir = dir
	out, err := diffCmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("git diff failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("git not found or not a git repository: %w", err)
	}

	// Collect git-relative paths from diff output
	gitRelPaths := splitNonEmpty(string(out))

	// For --changed mode, also include untracked files
	if changed {
		lsCmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
		lsCmd.Dir = dir
		lsOut, err := lsCmd.Output()
		if err == nil {
			gitRelPaths = append(gitRelPaths, splitNonEmpty(string(lsOut))...)
		}
	}

	// Early return for empty output
	if len(gitRelPaths) == 0 {
		return map[string]bool{}, nil
	}

	// Resolve absolute rootPath for proper relative path computation
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		absRoot = rootPath
	}
	// If rootPath is a file, use its parent directory
	if info, err := os.Stat(absRoot); err == nil && !info.IsDir() {
		absRoot = filepath.Dir(absRoot)
	}
	// Resolve symlinks to match repoRoot (e.g., macOS /var -> /private/var)
	if resolved, err := filepath.EvalSymlinks(absRoot); err == nil {
		absRoot = resolved
	}

	// Convert git-relative paths to rootPath-relative paths
	files := make(map[string]bool, len(gitRelPaths))
	for _, gitRel := range gitRelPaths {
		// Convert git-relative path to absolute path
		absPath := filepath.Join(repoRoot, filepath.FromSlash(gitRel))

		// Compute path relative to rootPath
		rel, err := filepath.Rel(absRoot, absPath)
		if err != nil {
			continue
		}

		// Skip entries outside the scan root (paths that start with "..")
		if strings.HasPrefix(rel, "..") {
			continue
		}

		files[filepath.ToSlash(rel)] = true
	}

	return files, nil
}

// splitNonEmpty splits output by newlines and returns non-empty trimmed lines.
func splitNonEmpty(s string) []string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}
	lines := strings.Split(trimmed, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result
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
