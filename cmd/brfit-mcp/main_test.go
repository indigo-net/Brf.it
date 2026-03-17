package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	_ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"
)

func TestSummarizeProject(t *testing.T) {
	// Create a temp directory with a Go file
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(goFile, []byte("package main\n\nfunc Hello() string { return \"hi\" }\n"), 0644); err != nil {
		t.Fatal(err)
	}

	handler := makeSummarizeProject(tmpDir)
	_, output, err := handler(context.Background(), &mcp.CallToolRequest{}, SummarizeProjectInput{
		Format: "xml",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.TotalFiles == 0 {
		t.Error("expected at least 1 file")
	}
	if output.TotalSignatures == 0 {
		t.Error("expected at least 1 signature")
	}
	if output.Content == "" {
		t.Error("expected non-empty content")
	}
}

func TestSummarizeFile(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmpDir, "pkg"), 0755); err != nil {
		t.Fatal(err)
	}
	goFile := filepath.Join(tmpDir, "pkg", "lib.go")
	if err := os.WriteFile(goFile, []byte("package pkg\n\nfunc Greet(name string) string { return name }\n"), 0644); err != nil {
		t.Fatal(err)
	}
	// Also create a file outside the include pattern
	otherFile := filepath.Join(tmpDir, "other.go")
	if err := os.WriteFile(otherFile, []byte("package main\n\nfunc Other() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	handler := makeSummarizeFile(tmpDir)
	_, output, err := handler(context.Background(), &mcp.CallToolRequest{}, SummarizeFileInput{
		Include: "pkg/**/*.go",
		Format:  "json",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.TotalFiles != 1 {
		t.Errorf("expected 1 file, got %d", output.TotalFiles)
	}
	if output.TotalSignatures == 0 {
		t.Error("expected at least 1 signature")
	}
}

func TestListLanguages(t *testing.T) {
	_, output, err := handleListLanguages(context.Background(), &mcp.CallToolRequest{}, ListLanguagesInput{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.Count == 0 {
		t.Error("expected at least 1 supported language")
	}
	if len(output.Languages) != output.Count {
		t.Errorf("languages length %d does not match count %d", len(output.Languages), output.Count)
	}

	// Verify languages are sorted
	for i := 1; i < len(output.Languages); i++ {
		if output.Languages[i] < output.Languages[i-1] {
			t.Errorf("languages not sorted: %q before %q", output.Languages[i-1], output.Languages[i])
			break
		}
	}

	// Verify Go is in the list (since we import treesitter parsers)
	found := false
	for _, lang := range output.Languages {
		if lang == "go" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected 'go' in languages list, got: %v", output.Languages)
	}
}

func TestSummarizeProjectInvalidPath(t *testing.T) {
	handler := makeSummarizeProject("/nonexistent/path")
	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, SummarizeProjectInput{})
	if err == nil {
		t.Error("expected error for nonexistent path")
	}
}

func TestPathTraversal(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		errPart string
	}{
		{"relative traversal", "../../etc", "resolves outside the project root"},
		{"absolute escape", "/etc", "absolute path"},
		{"dot-dot in middle", "subdir/../../etc", "resolves outside the project root"},
	}

	for _, tt := range tests {
		t.Run("project/"+tt.name, func(t *testing.T) {
			handler := makeSummarizeProject(tmpDir)
			_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, SummarizeProjectInput{
				Path: tt.path,
			})
			if err == nil {
				t.Error("expected error for path traversal attempt")
			}
			if err != nil && !strings.Contains(err.Error(), tt.errPart) {
				t.Errorf("expected %q error, got: %v", tt.errPart, err)
			}
		})

		t.Run("file/"+tt.name, func(t *testing.T) {
			handler := makeSummarizeFile(tmpDir)
			_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, SummarizeFileInput{
				Path: tt.path,
			})
			if err == nil {
				t.Error("expected error for path traversal attempt")
			}
			if err != nil && !strings.Contains(err.Error(), tt.errPart) {
				t.Errorf("expected %q error, got: %v", tt.errPart, err)
			}
		})
	}
}

func TestInvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(goFile, []byte("package main\n\nfunc Hello() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	handler := makeSummarizeProject(tmpDir)
	_, _, err := handler(context.Background(), &mcp.CallToolRequest{}, SummarizeProjectInput{
		Format: "yaml",
	})
	if err == nil {
		t.Error("expected error for invalid format")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid format") {
		t.Errorf("expected 'invalid format' error, got: %v", err)
	}
}

func TestValidSubdirectoryPath(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	goFile := filepath.Join(subDir, "main.go")
	if err := os.WriteFile(goFile, []byte("package main\n\nfunc Hello() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	handler := makeSummarizeProject(tmpDir)
	// Use relative path (absolute paths are now rejected)
	_, output, err := handler(context.Background(), &mcp.CallToolRequest{}, SummarizeProjectInput{
		Path: "sub",
	})
	if err != nil {
		t.Fatalf("unexpected error for valid subdirectory: %v", err)
	}
	if output.TotalFiles == 0 {
		t.Error("expected at least 1 file")
	}
}

func TestResolvePathAbsoluteRejected(t *testing.T) {
	tmpDir := t.TempDir()

	// Absolute paths should be rejected even if they point within the root
	absInRoot := filepath.Join(tmpDir, "sub")
	_, err := resolvePath(tmpDir, absInRoot)
	if err == nil {
		t.Error("expected error for absolute path input")
	}
	if err != nil && !strings.Contains(err.Error(), "absolute path") {
		t.Errorf("expected 'absolute path' error, got: %v", err)
	}
}

func TestResolvePathValidRelative(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	resolved, err := resolvePath(tmpDir, "sub")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Clean(filepath.Join(tmpDir, "sub"))
	if resolved != expected {
		t.Errorf("expected %q, got %q", expected, resolved)
	}
}

func TestResolvePathSymlinkEscape(t *testing.T) {
	tmpDir := t.TempDir()
	outsideDir := t.TempDir()

	// Create a file outside the root
	outsideFile := filepath.Join(outsideDir, "secret.txt")
	if err := os.WriteFile(outsideFile, []byte("secret"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a symlink inside root that points outside
	symlinkPath := filepath.Join(tmpDir, "escape")
	if err := os.Symlink(outsideDir, symlinkPath); err != nil {
		t.Skipf("cannot create symlinks: %v", err)
	}

	_, err := resolvePath(tmpDir, "escape")
	if err == nil {
		t.Error("expected error for symlink escaping project root")
	}
	if err != nil && !strings.Contains(err.Error(), "resolves outside the project root via symlink") {
		t.Errorf("expected symlink escape error, got: %v", err)
	}
}

func TestResolvePathSymlinkWithinRoot(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "real")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a symlink within root pointing to another location within root
	symlinkPath := filepath.Join(tmpDir, "link")
	if err := os.Symlink(subDir, symlinkPath); err != nil {
		t.Skipf("cannot create symlinks: %v", err)
	}

	resolved, err := resolvePath(tmpDir, "link")
	if err != nil {
		t.Fatalf("unexpected error for symlink within root: %v", err)
	}
	expected := filepath.Clean(filepath.Join(tmpDir, "link"))
	if resolved != expected {
		t.Errorf("expected %q, got %q", expected, resolved)
	}
}

func TestResolvePathEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	resolved, err := resolvePath(tmpDir, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != tmpDir {
		t.Errorf("expected %q, got %q", tmpDir, resolved)
	}
}
