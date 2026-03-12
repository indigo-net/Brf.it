// Package main provides the brfit MCP (Model Context Protocol) server.
// This server exposes brfit's code analysis capabilities as MCP tools,
// allowing AI agents to directly request code summaries and analysis.
//
// Usage:
//
//	brfit-mcp [--root <path>]
//
// Claude Desktop configuration:
//
//	{
//	  "mcpServers": {
//	    "brfit": {
//	      "command": "brfit-mcp",
//	      "args": ["--root", "/path/to/project"]
//	    }
//	  }
//	}
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/indigo-net/Brf.it/internal/config"
	pkgcontext "github.com/indigo-net/Brf.it/internal/context"
	"github.com/indigo-net/Brf.it/pkg/scanner"
	"github.com/modelcontextprotocol/go-sdk/mcp"

	// Register Tree-sitter parsers.
	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"
)

// Build information (set by ldflags).
var (
	version = "dev"
)

func main() {
	rootPath := flag.String("root", ".", "default project root path")
	flag.Parse()

	// Resolve absolute path
	absRoot, err := filepath.Abs(*rootPath)
	if err != nil {
		log.Fatalf("failed to resolve root path: %v", err)
	}

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "brfit",
			Version: version,
		},
		nil,
	)

	// Tool: summarize_project
	mcp.AddTool(server, &mcp.Tool{
		Name:        "summarize_project",
		Description: "Extract function signatures and documentation from a project directory, optimized for AI consumption. Returns XML/Markdown/JSON formatted output.",
	}, makeSummarizeProject(absRoot))

	// Tool: summarize_file
	mcp.AddTool(server, &mcp.Tool{
		Name:        "summarize_file",
		Description: "Extract function signatures and documentation from specific files matching a glob pattern.",
	}, makeSummarizeFile(absRoot))

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

// SummarizeProjectInput defines the input for the summarize_project tool.
type SummarizeProjectInput struct {
	Path          string `json:"path,omitempty" jsonschema:"project directory path (defaults to server root)"`
	Format        string `json:"format,omitempty" jsonschema:"output format: xml, md, or json (default: xml)"`
	IncludeBody   bool   `json:"include_body,omitempty" jsonschema:"include function bodies (default: false)"`
	IncludeImport bool   `json:"include_imports,omitempty" jsonschema:"include import statements (default: false)"`
	CallGraph     bool   `json:"call_graph,omitempty" jsonschema:"include function call graph (default: false)"`
}

// SummarizeProjectOutput defines the output for the summarize_project tool.
type SummarizeProjectOutput struct {
	Content         string `json:"content" jsonschema:"the formatted project summary"`
	TotalFiles      int    `json:"total_files" jsonschema:"number of files processed"`
	TotalSignatures int    `json:"total_signatures" jsonschema:"number of signatures extracted"`
	TokenCount      int    `json:"token_count,omitempty" jsonschema:"estimated token count of output"`
}

func makeSummarizeProject(defaultRoot string) func(context.Context, *mcp.CallToolRequest, SummarizeProjectInput) (*mcp.CallToolResult, SummarizeProjectOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SummarizeProjectInput) (*mcp.CallToolResult, SummarizeProjectOutput, error) {
		path, err := resolvePath(defaultRoot, input.Path)
		if err != nil {
			return nil, SummarizeProjectOutput{}, err
		}

		// Validate path exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, SummarizeProjectOutput{}, fmt.Errorf("path not found: %s", path)
		}

		format, err := validateFormat(input.Format)
		if err != nil {
			return nil, SummarizeProjectOutput{}, err
		}

		cfg := config.DefaultConfig()
		cfg.Path = path
		cfg.Format = format
		cfg.IncludeBody = input.IncludeBody
		cfg.IncludeImports = input.IncludeImport
		cfg.CallGraph = input.CallGraph

		result, err := runPackager(ctx, cfg)
		if err != nil {
			return nil, SummarizeProjectOutput{}, err
		}

		return nil, SummarizeProjectOutput{
			Content:         string(result.Content),
			TotalFiles:      result.TotalFiles,
			TotalSignatures: result.TotalSignatures,
			TokenCount:      result.TokenCount,
		}, nil
	}
}

// SummarizeFileInput defines the input for the summarize_file tool.
type SummarizeFileInput struct {
	Path    string `json:"path" jsonschema:"project directory path"`
	Include string `json:"include" jsonschema:"glob pattern to include (e.g. 'pkg/**/*.go')"`
	Format  string `json:"format,omitempty" jsonschema:"output format: xml, md, or json (default: xml)"`
}

// SummarizeFileOutput defines the output.
type SummarizeFileOutput struct {
	Content         string `json:"content" jsonschema:"the formatted file summary"`
	TotalFiles      int    `json:"total_files" jsonschema:"number of files processed"`
	TotalSignatures int    `json:"total_signatures" jsonschema:"number of signatures extracted"`
}

func makeSummarizeFile(defaultRoot string) func(context.Context, *mcp.CallToolRequest, SummarizeFileInput) (*mcp.CallToolResult, SummarizeFileOutput, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, input SummarizeFileInput) (*mcp.CallToolResult, SummarizeFileOutput, error) {
		path, err := resolvePath(defaultRoot, input.Path)
		if err != nil {
			return nil, SummarizeFileOutput{}, err
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, SummarizeFileOutput{}, fmt.Errorf("path not found: %s", path)
		}

		format, err := validateFormat(input.Format)
		if err != nil {
			return nil, SummarizeFileOutput{}, err
		}

		cfg := config.DefaultConfig()
		cfg.Path = path
		cfg.Format = format
		if input.Include != "" {
			cfg.IncludePatterns = []string{input.Include}
		}

		result, err := runPackager(ctx, cfg)
		if err != nil {
			return nil, SummarizeFileOutput{}, err
		}

		return nil, SummarizeFileOutput{
			Content:         string(result.Content),
			TotalFiles:      result.TotalFiles,
			TotalSignatures: result.TotalSignatures,
		}, nil
	}
}

// validFormats is the set of accepted output format values.
var validFormats = map[string]bool{"xml": true, "md": true, "markdown": true, "json": true}

// resolvePath resolves an input path relative to the default root, ensuring
// the result stays within the root directory to prevent path traversal.
// If inputPath is empty, defaultRoot is returned unchanged.
func resolvePath(defaultRoot, inputPath string) (string, error) {
	if inputPath == "" {
		return defaultRoot, nil
	}

	absPath, err := filepath.Abs(inputPath)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}

	// Ensure the resolved path is within the project root.
	rel, err := filepath.Rel(defaultRoot, absPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("path %q is outside the project root %q", inputPath, defaultRoot)
	}

	return absPath, nil
}

// validateFormat checks that format is a supported value, returning the
// validated format or "xml" as default when format is empty.
func validateFormat(format string) (string, error) {
	if format == "" {
		return "xml", nil
	}
	if !validFormats[format] {
		return "", fmt.Errorf("invalid format %q: must be xml, md, markdown, or json", format)
	}
	return format, nil
}

// runPackager creates and runs a Packager with the given config.
func runPackager(ctx context.Context, cfg *config.Config) (*pkgcontext.Result, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration error: %w", err)
	}

	scanOpts := &scanner.ScanOptions{
		RootPath:            cfg.Path,
		SupportedExtensions: cfg.SupportedExtensions(),
		IgnoreFiles:         cfg.IgnoreFiles,
		IncludePatterns:     cfg.IncludePatterns,
		ExcludePatterns:     cfg.ExcludePatterns,
		IncludeHidden:       cfg.IncludeHidden,
		MaxFileSize:         cfg.MaxFileSize,
	}

	packager, err := pkgcontext.NewDefaultPackager(scanOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize: %w", err)
	}

	opts := cfg.ToOptions()
	result, err := packager.Package(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("processing failed: %w", err)
	}

	return result, nil
}
