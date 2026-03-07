package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/indigo-net/Brf.it/internal/config"

	// Import treesitter parser to register it
	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"
)

func TestExecuteHelp(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set args for help
	os.Args = []string{"brfit", "--help"}

	// Execute and capture panic (since help exits)
	defer func() {
		os.Stdout = old
		recover()
		_ = w.Close()
	}()

	Execute()

	_ = r.Close()
	os.Stdout = old
}

func TestExecuteVersion(t *testing.T) {
	// Capture stdout to get the version output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a fresh command to avoid global state issues
	testCfg := config.DefaultConfig()
	cmd := newRootCommandWithConfig(testCfg)
	cmd.SetArgs([]string{"--version"})

	// Execute in goroutine since it may call os.Exit
	done := make(chan bool)
	go func() {
		cmd.Execute()
		done <- true
	}()

	// Wait for completion with timeout
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for command execution")
	}
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check for new format: "brfit VERSION (commit: COMMIT, built: DATE)"
	if !strings.Contains(output, "brfit") || !strings.Contains(output, "commit:") {
		t.Errorf("expected version output to contain 'brfit' and 'commit:', got '%s'", output)
	}
}

func TestNewRootCommand(t *testing.T) {
	testCfg := config.DefaultConfig()
	cmd := newRootCommandWithConfig(testCfg)

	if cmd.Use != "brfit [path] [options]" {
		t.Errorf("expected Use 'brfit [path] [options]', got '%s'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be set")
	}

	if cmd.Long == "" {
		t.Error("expected Long description to be set")
	}

	// Check flags exist
	flags := []string{"mode", "format", "output", "ignore", "include-hidden", "no-tree", "no-tokens", "max-size"}
	for _, flag := range flags {
		f := cmd.Flags().Lookup(flag)
		if f == nil {
			t.Errorf("expected flag '%s' to exist", flag)
		}
	}

	// Check version flag
	if cmd.Flags().Lookup("version") == nil {
		t.Error("expected version flag to exist")
	}
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectedPath string
		expectedMode string
		expectedFmt  string
		expectedOut  string
		expectedHide bool
		expectedTree bool
		expectedTok  bool
		expectedSize int64
	}{
		{
			name:         "default values",
			args:         []string{"brfit"},
			expectedPath: ".",
			expectedMode: "sig",
			expectedFmt:  "xml",
			expectedOut:  "",
			expectedHide: false,
			expectedTree: false,
			expectedTok:  false,
			expectedSize: 512000,
		},
		{
			name:         "with path argument",
			args:         []string{"brfit", "./src"},
			expectedPath: "./src",
			expectedMode: "sig",
			expectedFmt:  "xml",
		},
		{
			name:         "with all flags",
			args:         []string{"brfit", "/project", "--mode", "sig", "--format", "md", "--output", "out.xml", "--include-hidden", "--no-tree", "--no-tokens", "--max-size", "1000000"},
			expectedPath: "/project",
			expectedMode: "sig",
			expectedFmt:  "md",
			expectedOut:  "out.xml",
			expectedHide: true,
			expectedTree: true,
			expectedTok:  true,
			expectedSize: 1000000,
		},
		{
			name:         "short flags",
			args:         []string{"brfit", "-m", "sig", "-f", "md", "-o", "output.md"},
			expectedPath: ".",
			expectedMode: "sig",
			expectedFmt:  "md",
			expectedOut:  "output.md",
		},
		{
			name:         "custom ignore file",
			args:         []string{"brfit", "-i", ".brfitignore"},
			expectedPath: ".",
			expectedMode: "sig",
			expectedFmt:  "xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 각 테스트마다 독립적인 config와 command 생성
			testCfg := config.DefaultConfig()
			cmd := newRootCommandWithConfig(testCfg)

			// Set args (without program name)
			cmd.SetArgs(tt.args[1:])

			// Execute ParseFlags
			if err := cmd.ParseFlags(tt.args[1:]); err != nil {
				t.Fatalf("ParseFlags returned error: %v", err)
			}

			// Set path from positional args
			testCfg.Path = "."
			remaining := cmd.Flags().Args()
			if len(remaining) > 0 {
				testCfg.Path = remaining[0]
			}

			if testCfg.Path != tt.expectedPath {
				t.Errorf("expected Path '%s', got '%s'", tt.expectedPath, testCfg.Path)
			}
			if testCfg.Mode != tt.expectedMode {
				t.Errorf("expected Mode '%s', got '%s'", tt.expectedMode, testCfg.Mode)
			}
			if testCfg.Format != tt.expectedFmt {
				t.Errorf("expected Format '%s', got '%s'", tt.expectedFmt, testCfg.Format)
			}
			if testCfg.Output != tt.expectedOut {
				t.Errorf("expected Output '%s', got '%s'", tt.expectedOut, testCfg.Output)
			}
			if testCfg.IncludeHidden != tt.expectedHide {
				t.Errorf("expected IncludeHidden %v, got %v", tt.expectedHide, testCfg.IncludeHidden)
			}
			if testCfg.NoTree != tt.expectedTree {
				t.Errorf("expected NoTree %v, got %v", tt.expectedTree, testCfg.NoTree)
			}
			if testCfg.NoTokens != tt.expectedTok {
				t.Errorf("expected NoTokens %v, got %v", tt.expectedTok, testCfg.NoTokens)
			}
			if tt.expectedSize != 0 && testCfg.MaxFileSize != tt.expectedSize {
				t.Errorf("expected MaxFileSize %d, got %d", tt.expectedSize, testCfg.MaxFileSize)
			}
		})
	}
}

func TestRootCommandIntegration(t *testing.T) {
	// Create temp directory with sample Go file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(testFile, []byte("package main\n\nfunc Add(a, b int) int { return a + b }\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run command with XML format
	cmd := newRootCommandWithConfig(config.DefaultConfig())
	cmd.SetArgs([]string{tmpDir, "-f", "xml"})

	err := cmd.Execute()

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Verify output contains expected XML
	if !strings.Contains(output, "<?xml") {
		t.Error("expected XML output")
	}
	if !strings.Contains(output, "<brfit>") {
		t.Error("expected <brfit> root element")
	}
	if !strings.Contains(output, "func Add(a, b int) int") {
		t.Error("expected function signature in output")
	}
}

func TestRootCommandIntegrationMarkdown(t *testing.T) {
	// Create temp directory with sample Go file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(testFile, []byte("package main\n\nfunc Subtract(a, b int) int { return a - b }\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run command with Markdown format
	cmd := newRootCommandWithConfig(config.DefaultConfig())
	cmd.SetArgs([]string{tmpDir, "-f", "md"})

	err := cmd.Execute()

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Fatalf("command failed: %v", err)
	}

	// Verify output contains expected Markdown
	if !strings.Contains(output, "# Code Summary") {
		t.Error("expected Markdown header")
	}
	if !strings.Contains(output, "func Subtract(a, b int) int") {
		t.Error("expected function signature in output")
	}
}

func TestRootCommandPathNotFound(t *testing.T) {
	cmd := newRootCommandWithConfig(config.DefaultConfig())
	cmd.SetArgs([]string{"/nonexistent/path/12345"})

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for non-existent path")
	}
	if !strings.Contains(err.Error(), "path not found") {
		t.Errorf("expected 'path not found' error, got: %v", err)
	}
}

func TestWriteToFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "subdir", "output.xml")

	content := []byte("<?xml version=\"1.0\"?>\n<test>content</test>\n")

	if err := writeToFile(outputPath, content); err != nil {
		t.Fatalf("writeToFile failed: %v", err)
	}

	// Verify file was created
	readContent, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if string(readContent) != string(content) {
		t.Errorf("file content mismatch: got %q, want %q", string(readContent), string(content))
	}
}
