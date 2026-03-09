package formatter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/indigo-net/Brf.it/pkg/parser"
)

func TestXMLFormatterImplementsFormatter(t *testing.T) {
	var _ Formatter = (*XMLFormatter)(nil)
}

func TestMarkdownFormatterImplementsFormatter(t *testing.T) {
	var _ Formatter = (*MarkdownFormatter)(nil)
}

func TestJSONFormatterImplementsFormatter(t *testing.T) {
	var _ Formatter = (*JSONFormatter)(nil)
}

func TestXMLFormatterFormat(t *testing.T) {
	formatter := NewXMLFormatter()

	data := &PackageData{
		Version:  "v0.12.0",
		RootPath: "/path/to/project",
		Tree:     "pkg/\n└── test.go",
		Files: []FileData{
			{
				Path:     "pkg/test.go",
				Language: "go",
				Signatures: []parser.Signature{
					{
						Name:     "Add",
						Kind:     "function",
						Text:     "func Add(a, b int) int",
						Doc:      "Add returns the sum of two integers.",
						Line:     5,
						Language: "go",
						Exported: true,
					},
				},
			},
		},
		TotalSignatures: 1,
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	// Verify XML structure
	outputStr := string(output)

	if !strings.Contains(outputStr, `<?xml version="1.0"`) {
		t.Error("expected XML declaration")
	}

	if !strings.Contains(outputStr, "<brfit>") {
		t.Error("expected <brfit> root element")
	}

	if !strings.Contains(outputStr, "<metadata>") {
		t.Error("expected <metadata> element")
	}

	// Verify metadata contains version, path, and schema
	if !strings.Contains(outputStr, "<version>v0.12.0</version>") {
		t.Error("expected <version> element in metadata")
	}

	if !strings.Contains(outputStr, "<path>/path/to/project</path>") {
		t.Error("expected <path> element in metadata")
	}

	if !strings.Contains(outputStr, "<schema>") {
		t.Error("expected <schema> element in metadata")
	}

	if !strings.Contains(outputStr, `<tag name="metadata" description="Project metadata container" />`) {
		t.Error("expected metadata tag in schema")
	}

	if !strings.Contains(outputStr, `<tag name="file" description="Source file (path, language attributes)" />`) {
		t.Error("expected file tag in schema")
	}

	if !strings.Contains(outputStr, `<file path="pkg/test.go" language="go"`) {
		t.Error("expected file element with path and language attributes")
	}

	if !strings.Contains(outputStr, "<function>func Add(a, b int) int</function>") {
		t.Error("expected function element for Kind='function'")
	}

	if !strings.Contains(outputStr, "<doc>Add returns the sum of two integers.</doc>") {
		t.Error("expected doc element")
	}
}

func TestXMLFormatterFormatWithError(t *testing.T) {
	formatter := NewXMLFormatter()

	data := &PackageData{
		Files: []FileData{
			{
				Path:     "test.py",
				Language: "python",
				Error:    fmt.Errorf("no parser for language: python"),
			},
		},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "<error>no parser for language: python</error>") {
		t.Error("expected error element")
	}
}

func TestMarkdownFormatterFormat(t *testing.T) {
	formatter := NewMarkdownFormatter()

	data := &PackageData{
		Tree: "pkg/\n└── test.go",
		Files: []FileData{
			{
				Path:     "pkg/test.go",
				Language: "go",
				Signatures: []parser.Signature{
					{
						Name:     "Add",
						Kind:     "function",
						Text:     "func Add(a, b int) int",
						Doc:      "Add returns the sum.",
						Line:     5,
						Language: "go",
						Exported: true,
					},
				},
			},
		},
		TotalSignatures: 1,
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "# Code Summary") {
		t.Error("expected header")
	}

	if !strings.Contains(outputStr, "## Directory Tree") {
		t.Error("expected Directory Tree section")
	}

	if !strings.Contains(outputStr, "## Files") {
		t.Error("expected Files section")
	}

	if !strings.Contains(outputStr, "### pkg/test.go") {
		t.Error("expected file path as heading")
	}

	if !strings.Contains(outputStr, "```go") {
		t.Error("expected code block with language")
	}

	if !strings.Contains(outputStr, "func Add(a, b int) int") {
		t.Error("expected signature in code block")
	}

	if strings.Contains(outputStr, "// function") {
		t.Error("kind comment should not appear in markdown output")
	}
}

func TestFormatterNames(t *testing.T) {
	xmlFormatter := NewXMLFormatter()
	if xmlFormatter.Name() != "xml" {
		t.Errorf("expected name 'xml', got '%s'", xmlFormatter.Name())
	}

	mdFormatter := NewMarkdownFormatter()
	if mdFormatter.Name() != "markdown" {
		t.Errorf("expected name 'markdown', got '%s'", mdFormatter.Name())
	}
}

func TestXMLFormatterEscapeXML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"a & b", "a &amp; b"},
		{"<tag>", "&lt;tag&gt;"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"it's", "it&apos;s"},
		{"all & < > \" '", "all &amp; &lt; &gt; &quot; &apos;"},
	}

	for _, tt := range tests {
		result := escapeXML(tt.input)
		if result != tt.expected {
			t.Errorf("escapeXML(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestMarkdownFormatterEscapeMarkdown(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"code `here`", "code \\`here\\`"},
		{"no escape needed", "no escape needed"},
	}

	for _, tt := range tests {
		result := escapeMarkdown(tt.input)
		if result != tt.expected {
			t.Errorf("escapeMarkdown(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestXMLFormatterEmptyData(t *testing.T) {
	formatter := NewXMLFormatter()

	data := &PackageData{
		Files: []FileData{},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "<brfit>") {
		t.Error("expected <brfit> root element")
	}

	if !strings.Contains(outputStr, "<files>\n  </files>") {
		t.Error("expected empty files section")
	}
}

func TestMarkdownFormatterEmptyData(t *testing.T) {
	formatter := NewMarkdownFormatter()

	data := &PackageData{
		Files: []FileData{},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, "# Code Summary") {
		t.Error("expected header")
	}

	if !strings.Contains(outputStr, "## Files") {
		t.Error("expected Files section")
	}
}

func TestMarkdownFormatterEmptyFile(t *testing.T) {
	formatter := NewMarkdownFormatter()
	data := &PackageData{
		Files: []FileData{
			{
				Path:       "cmd/main.go",
				Language:   "go",
				Signatures: []parser.Signature{}, // 빈 시그니처
				RawImports: []string{},
			},
		},
		IncludeImports: false,
	}

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	output := string(result)
	if !strings.Contains(output, "// (empty)") {
		t.Errorf("Expected empty comment, got:\n%s", output)
	}
}

func TestMarkdownFormatterEmptyFileWithImports(t *testing.T) {
	formatter := NewMarkdownFormatter()
	data := &PackageData{
		Files: []FileData{
			{
				Path:       "cmd/main.go",
				Language:   "go",
				Signatures: []parser.Signature{},
				RawImports: []string{`import "fmt"`},
			},
		},
		IncludeImports: true,
	}

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	output := string(result)
	// imports가 있으면 빈 파일이 아님
	if strings.Contains(output, "// (empty)") {
		t.Errorf("Should not show empty comment when imports exist, got:\n%s", output)
	}
}

func TestXMLFormatterEmptyFile(t *testing.T) {
	formatter := NewXMLFormatter()
	data := &PackageData{
		Files: []FileData{
			{
				Path:       "cmd/main.go",
				Language:   "go",
				Signatures: []parser.Signature{},
				RawImports: []string{},
			},
		},
		IncludeImports: false,
	}

	result, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	output := string(result)
	if !strings.Contains(output, "<!-- empty -->") {
		t.Errorf("Expected empty comment, got:\n%s", output)
	}
}

func TestKindToTag(t *testing.T) {
	tests := []struct {
		kind     string
		expected string
	}{
		// function 그룹
		{"function", "function"},
		{"method", "function"},
		{"constructor", "function"},
		{"destructor", "function"},
		{"arrow", "function"},
		{"local_function", "function"},
		{"module_function", "function"},

		// type 그룹
		{"class", "type"},
		{"interface", "type"},
		{"type", "type"},
		{"struct", "type"},
		{"enum", "type"},
		{"record", "type"},
		{"annotation", "type"},
		{"typedef", "type"},
		{"namespace", "type"},
		{"template", "type"},

		// variable 그룹
		{"variable", "variable"},
		{"field", "variable"},
		{"macro", "variable"},
		{"export", "variable"},

		// fallback
		{"", "signature"},
		{"unknown", "signature"},
		{"custom_kind", "signature"},
	}

	for _, tt := range tests {
		got := kindToTag(tt.kind)
		if got != tt.expected {
			t.Errorf("kindToTag(%q) = %q, want %q", tt.kind, got, tt.expected)
		}
	}
}

func TestXMLFormatterKindTags(t *testing.T) {
	formatter := NewXMLFormatter()

	tests := []struct {
		name     string
		kind     string
		text     string
		wantTag  string
	}{
		// TypeScript arrow function
		{"ts_arrow", "arrow", "const add = (a, b) => a + b", "function"},
		// TypeScript export
		{"ts_export", "export", "export { foo }", "variable"},
		// C++ constructor
		{"cpp_constructor", "constructor", "MyClass::MyClass()", "function"},
		// C++ destructor
		{"cpp_destructor", "destructor", "MyClass::~MyClass()", "function"},
		// Python method (런타임 판단 후)
		{"py_method", "method", "def foo(self):", "function"},
		// Java static field
		{"java_variable", "variable", "public static final int X = 1", "variable"},
		// Go type
		{"go_type", "type", "type Config struct { Path string }", "type"},
		// Go struct (direct)
		{"go_struct", "struct", "struct { X int }", "type"},
		// Java class
		{"java_class", "class", "public class MyClass", "type"},
		// Java interface
		{"java_interface", "interface", "public interface MyInterface", "type"},
		// C macro
		{"c_macro", "macro", "#define MAX 100", "variable"},
		// Empty kind fallback
		{"empty_kind", "", "func foo()", "signature"},
		// Unknown kind fallback
		{"unknown_kind", "something_new", "unknown declaration", "signature"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &PackageData{
				Files: []FileData{
					{
						Path:     "test.go",
						Language: "go",
						Signatures: []parser.Signature{
							{
								Name: "test",
								Kind: tt.kind,
								Text: tt.text,
							},
						},
					},
				},
			}

			output, err := formatter.Format(data)
			if err != nil {
				t.Fatal(err)
			}

			outputStr := string(output)
			expectedOpen := fmt.Sprintf("<%s>", tt.wantTag)
			expectedClose := fmt.Sprintf("</%s>", tt.wantTag)

			if !strings.Contains(outputStr, expectedOpen) {
				t.Errorf("expected opening tag %s in output:\n%s", expectedOpen, outputStr)
			}
			if !strings.Contains(outputStr, expectedClose) {
				t.Errorf("expected closing tag %s in output:\n%s", expectedClose, outputStr)
			}
		})
	}
}

func TestXMLFormatterVerbatimImports(t *testing.T) {
	formatter := NewXMLFormatter()

	data := &PackageData{
		Files: []FileData{
			{
				Path:     "main.go",
				Language: "go",
				Signatures: []parser.Signature{
					{Name: "Main", Kind: "function", Text: "func Main()"},
				},
				RawImports: []string{
					`import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)`,
				},
			},
		},
		IncludeImports: true,
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)

	// imports section should be present
	if !strings.Contains(outputStr, "<imports>") {
		t.Error("expected <imports> section")
	}

	// raw import block should be included verbatim (including multi-line format)
	if !strings.Contains(outputStr, "fmt") {
		t.Error("expected import 'fmt' to be included")
	}
	if !strings.Contains(outputStr, "cobra") {
		t.Error("expected import 'cobra' to be included")
	}
}

func TestMarkdownFormatterVerbatimImports(t *testing.T) {
	formatter := NewMarkdownFormatter()

	data := &PackageData{
		Files: []FileData{
			{
				Path:     "main.go",
				Language: "go",
				Signatures: []parser.Signature{
					{Name: "Main", Kind: "function", Text: "func Main()"},
				},
				RawImports: []string{
					`import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)`,
				},
			},
		},
		IncludeImports: true,
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)

	// imports should be included in the code block
	if !strings.Contains(outputStr, "fmt") {
		t.Error("expected import 'fmt' to be included")
	}
	if !strings.Contains(outputStr, "cobra") {
		t.Error("expected import 'cobra' to be included")
	}
}

func TestJSONFormatterFormat(t *testing.T) {
	formatter := NewJSONFormatter()

	data := &PackageData{
		Version:  "v0.17.0",
		RootPath: "/path/to/project",
		Tree:     "pkg/\n└── test.go",
		Files: []FileData{
			{
				Path:     "pkg/test.go",
				Language: "go",
				Signatures: []parser.Signature{
					{
						Name:     "Add",
						Kind:     "function",
						Text:     "func Add(a, b int) int",
						Doc:      "Add returns the sum of two integers.",
						Line:     5,
						Language: "go",
						Exported: true,
					},
				},
			},
		},
		TotalSignatures: 1,
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)

	// Verify JSON structure
	if !strings.Contains(outputStr, `"version":"v0.17.0"`) {
		t.Error("expected version in JSON output")
	}

	if !strings.Contains(outputStr, `"path":"/path/to/project"`) {
		t.Error("expected path in JSON output")
	}

	if !strings.Contains(outputStr, `"tree":"`) {
		t.Error("expected tree in JSON output")
	}

	if !strings.Contains(outputStr, `"files":[`) {
		t.Error("expected files array in JSON output")
	}

	if !strings.Contains(outputStr, `"path":"pkg/test.go"`) {
		t.Error("expected file path in JSON output")
	}

	if !strings.Contains(outputStr, `"language":"go"`) {
		t.Error("expected language in JSON output")
	}

	if !strings.Contains(outputStr, `"kind":"function"`) {
		t.Error("expected kind in JSON output")
	}

	if !strings.Contains(outputStr, `"text":"func Add(a, b int) int"`) {
		t.Error("expected signature text in JSON output")
	}

	if !strings.Contains(outputStr, `"doc":"Add returns the sum of two integers."`) {
		t.Error("expected doc in JSON output")
	}
}

func TestJSONFormatterKindNormalization(t *testing.T) {
	formatter := NewJSONFormatter()

	tests := []struct {
		name         string
		kind         string
		expectedKind string
	}{
		{"method", "method", "function"},
		{"constructor", "constructor", "function"},
		{"class", "class", "type"},
		{"interface", "interface", "type"},
		{"struct", "struct", "type"},
		{"variable", "variable", "variable"},
		{"field", "field", "variable"},
		{"unknown", "unknown_kind", "unknown_kind"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &PackageData{
				Files: []FileData{
					{
						Path:     "test.go",
						Language: "go",
						Signatures: []parser.Signature{
							{Kind: tt.kind, Text: "test"},
						},
					},
				},
			}

			output, err := formatter.Format(data)
			if err != nil {
				t.Fatal(err)
			}

			outputStr := string(output)
			expected := fmt.Sprintf(`"kind":"%s"`, tt.expectedKind)
			if !strings.Contains(outputStr, expected) {
				t.Errorf("expected kind %s, got: %s", tt.expectedKind, outputStr)
			}
		})
	}
}

func TestJSONFormatterWithImports(t *testing.T) {
	formatter := NewJSONFormatter()

	data := &PackageData{
		Files: []FileData{
			{
				Path:     "main.go",
				Language: "go",
				Signatures: []parser.Signature{
					{Name: "Main", Kind: "function", Text: "func Main()"},
				},
				RawImports: []string{
					`import "fmt"`,
					`import "os"`,
				},
			},
		},
		IncludeImports: true,
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)

	if !strings.Contains(outputStr, `"imports":[`) {
		t.Error("expected imports array in JSON output")
	}

	// JSON escapes quotes, so we check for the escaped version
	if !strings.Contains(outputStr, `import \"fmt\"`) {
		t.Error("expected 'fmt' import in JSON output")
	}

	if !strings.Contains(outputStr, `import \"os\"`) {
		t.Error("expected 'os' import in JSON output")
	}
}

func TestJSONFormatterWithError(t *testing.T) {
	formatter := NewJSONFormatter()

	data := &PackageData{
		Files: []FileData{
			{
				Path:     "test.py",
				Language: "python",
				Error:    fmt.Errorf("no parser for language: python"),
			},
		},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, `"error":"no parser for language: python"`) {
		t.Error("expected error field in JSON output")
	}
}

func TestJSONFormatterName(t *testing.T) {
	formatter := NewJSONFormatter()
	if formatter.Name() != "json" {
		t.Errorf("expected name 'json', got '%s'", formatter.Name())
	}
}

func TestXMLFormatterWithNoSchema(t *testing.T) {
	formatter := NewXMLFormatter()
	data := &PackageData{
		NoSchema: true,
		Files: []FileData{
			{
				Path:     "test.go",
				Language: "go",
			},
		},
	}

	output, err := formatter.Format(data)
	if err != nil {
		t.Fatal(err)
	}

	outputStr := string(output)
	if strings.Contains(outputStr, "<schema>") {
		t.Error("expected no <schema> section with --no-schema flag")
	}
}
