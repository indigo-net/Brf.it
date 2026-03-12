package main

import (
	"bytes"
	"os"
	"os/exec"
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
	flags := []string{"mode", "format", "output", "ignore", "include", "exclude", "include-hidden", "include-private", "no-tree", "no-tokens", "max-size", "changed", "since", "token-tree"}
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
		expectedHide    bool
		expectedPrivate bool
		expectedTree    bool
		expectedTok     bool
		expectedSize    int64
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
			args:            []string{"brfit", "/project", "--mode", "sig", "--format", "md", "--output", "out.xml", "--include-hidden", "--include-private", "--no-tree", "--no-tokens", "--max-size", "1000000"},
			expectedPath:    "/project",
			expectedMode:    "sig",
			expectedFmt:     "md",
			expectedOut:     "out.xml",
			expectedHide:    true,
			expectedPrivate: true,
			expectedTree:    true,
			expectedTok:     true,
			expectedSize:    1000000,
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
			if testCfg.IncludePrivate != tt.expectedPrivate {
				t.Errorf("expected IncludePrivate %v, got %v", tt.expectedPrivate, testCfg.IncludePrivate)
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

func TestResolveChangedFilesPathAnchoring(t *testing.T) {
	// Create a temporary git repo with a subdirectory, make a change,
	// and verify resolveChangedFiles returns paths relative to the subdirectory.
	tmpDir := t.TempDir()

	// Initialize a git repo
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = tmpDir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command %v failed: %v\n%s", args, err, out)
		}
	}

	run("git", "init")
	run("git", "config", "user.email", "test@test.com")
	run("git", "config", "user.name", "test")

	// Create subdirectory structure
	subDir := filepath.Join(tmpDir, "pkg", "scanner")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	otherDir := filepath.Join(tmpDir, "cmd")
	if err := os.MkdirAll(otherDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create initial files and commit
	for _, f := range []struct {
		path, content string
	}{
		{filepath.Join(subDir, "scanner.go"), "package scanner\n"},
		{filepath.Join(otherDir, "main.go"), "package main\n"},
	} {
		if err := os.WriteFile(f.path, []byte(f.content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	run("git", "add", "-A")
	run("git", "commit", "-m", "initial")

	// Modify a file in the subdirectory
	if err := os.WriteFile(filepath.Join(subDir, "scanner.go"), []byte("package scanner\n// changed\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test: resolveChangedFiles with rootPath = subDir should return "scanner.go"
	files, err := resolveChangedFiles(subDir, true, "")
	if err != nil {
		t.Fatalf("resolveChangedFiles failed: %v", err)
	}

	if !files["scanner.go"] {
		t.Errorf("expected 'scanner.go' in changed files, got: %v", files)
	}

	// cmd/main.go should NOT appear (it's outside the scan root)
	if files["../cmd/main.go"] || files["cmd/main.go"] || files["main.go"] {
		t.Errorf("expected files outside scan root to be excluded, got: %v", files)
	}
}

func TestResolveChangedFilesIncludesUntracked(t *testing.T) {
	tmpDir := t.TempDir()

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = tmpDir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command %v failed: %v\n%s", args, err, out)
		}
	}

	run("git", "init")
	run("git", "config", "user.email", "test@test.com")
	run("git", "config", "user.name", "test")

	// Create initial file and commit
	if err := os.WriteFile(filepath.Join(tmpDir, "existing.go"), []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run("git", "add", "-A")
	run("git", "commit", "-m", "initial")

	// Create a new untracked file
	if err := os.WriteFile(filepath.Join(tmpDir, "newfile.go"), []byte("package main\n// new\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// --changed (changed=true) should include untracked files
	files, err := resolveChangedFiles(tmpDir, true, "")
	if err != nil {
		t.Fatalf("resolveChangedFiles failed: %v", err)
	}

	if !files["newfile.go"] {
		t.Errorf("expected 'newfile.go' in changed files (untracked), got: %v", files)
	}

	// Modify existing.go and create a second commit so HEAD~1 is valid
	if err := os.WriteFile(filepath.Join(tmpDir, "existing.go"), []byte("package main\n// modified\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run("git", "add", "existing.go")
	run("git", "commit", "-m", "second")

	// --since (changed=false) should NOT include untracked files
	filesSince, err := resolveChangedFiles(tmpDir, false, "HEAD~1")
	if err != nil {
		t.Fatalf("resolveChangedFiles --since failed: %v", err)
	}

	if filesSince["newfile.go"] {
		t.Errorf("expected 'newfile.go' to NOT appear in --since mode, got: %v", filesSince)
	}
}

func TestResolveChangedFilesEmptyOutput(t *testing.T) {
	tmpDir := t.TempDir()

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = tmpDir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command %v failed: %v\n%s", args, err, out)
		}
	}

	run("git", "init")
	run("git", "config", "user.email", "test@test.com")
	run("git", "config", "user.name", "test")

	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run("git", "add", "-A")
	run("git", "commit", "-m", "initial")

	// No changes - should return empty map, not nil
	files, err := resolveChangedFiles(tmpDir, true, "")
	if err != nil {
		t.Fatalf("resolveChangedFiles failed: %v", err)
	}

	if files == nil {
		t.Error("expected non-nil empty map, got nil")
	}
	if len(files) != 0 {
		t.Errorf("expected empty map, got: %v", files)
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
