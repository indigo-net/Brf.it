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
		name string
		path string
	}{
		{"relative traversal", "../../etc"},
		{"absolute escape", "/etc"},
		{"dot-dot in middle", "subdir/../../etc"},
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
			if err != nil && !strings.Contains(err.Error(), "outside the project root") {
				t.Errorf("expected 'outside the project root' error, got: %v", err)
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
			if err != nil && !strings.Contains(err.Error(), "outside the project root") {
				t.Errorf("expected 'outside the project root' error, got: %v", err)
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
	_, output, err := handler(context.Background(), &mcp.CallToolRequest{}, SummarizeProjectInput{
		Path: subDir,
	})
	if err != nil {
		t.Fatalf("unexpected error for valid subdirectory: %v", err)
	}
	if output.TotalFiles == 0 {
		t.Error("expected at least 1 file")
	}
}
