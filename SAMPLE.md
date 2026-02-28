# Code Summary: /home/runner/work/Brf.it/Brf.it

*brf.it 0.11.0*

---

## Files

### /home/runner/work/Brf.it/Brf.it/cmd/brfit/main.go

```go
version = "dev"
commit  = "none"
date    = "unknown"
func main()
```

---

### /home/runner/work/Brf.it/Brf.it/cmd/brfit/root.go

**Imports:**
- `import "fmt"`
- `import "os"`
- `import "path/filepath"`
- `import "github.com/indigo-net/Brf.it/internal/config"`
- `import "github.com/indigo-net/Brf.it/internal/context"`
- `import "github.com/indigo-net/Brf.it/pkg/scanner"`
- `import "github.com/spf13/cobra"`
- `import _ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"`

```go
Version = "dev"
Commit  = "none"
Date    = "unknown"
func SetBuildInfo(v, c, d string)
cfg *config.Config
rootCmd *cobra.Command
func init()
func Execute()
func NewRootCommand() *cobra.Command
func newRootCommandWithConfig(c *config.Config) *cobra.Command
func addFlags(cmd *cobra.Command, c *config.Config)
func runRoot(cmd *cobra.Command, args []string, c *config.Config) error
func writeOutput(result *context.Result, c *config.Config) error
func writeToFile(path string, content []byte) error
```

---

### /home/runner/work/Brf.it/Brf.it/cmd/brfit/root_test.go

**Imports:**
- `import "bytes"`
- `import "os"`
- `import "path/filepath"`
- `import "strings"`
- `import "testing"`
- `import "time"`
- `import "github.com/indigo-net/Brf.it/internal/config"`
- `import _ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"`

```go
func TestExecuteHelp(t *testing.T)
func TestExecuteVersion(t *testing.T)
buf bytes.Buffer
func TestNewRootCommand(t *testing.T)
func TestParseFlags(t *testing.T)
func TestRootCommandIntegration(t *testing.T)
buf bytes.Buffer
func TestRootCommandIntegrationMarkdown(t *testing.T)
buf bytes.Buffer
func TestRootCommandPathNotFound(t *testing.T)
func TestWriteToFile(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/internal/config/config.go

**Imports:**
- `import "errors"`
- `import "fmt"`
- `import pkgcontext "github.com/indigo-net/Brf.it/internal/context"`

```go
type Config struct {
	// Path is the root directory or file to process.
	Path string

	// Version is the brf.it version string.
	Version string

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

	// IncludeImports determines whether to include import/export statements.
	IncludeImports bool

	// NoTree skips directory tree generation in output.
	NoTree bool

	// NoTokens disables token count calculation.
	NoTokens bool

	// MaxFileSize is the maximum file size in bytes to process.
	MaxFileSize int64
}
func DefaultConfig() *Config
func (c *Config) Validate() error
func (c *Config) SupportedExtensions() map[string]string
func (c *Config) ToOptions() *pkgcontext.Options
```

---

### /home/runner/work/Brf.it/Brf.it/internal/config/config_test.go

**Imports:**
- `import "testing"`

```go
func TestDefaultConfig(t *testing.T)
expectedMaxSize = 512000
func TestConfigValidate(t *testing.T)
func TestConfigSupportedLanguages(t *testing.T)
func containsString(s, substr string) bool
func containsSubstring(s, substr string) bool
```

---

### /home/runner/work/Brf.it/Brf.it/internal/context/context.go

**Imports:**
- `import "github.com/indigo-net/Brf.it/pkg/extractor"`
- `import "github.com/indigo-net/Brf.it/pkg/formatter"`
- `import "github.com/indigo-net/Brf.it/pkg/scanner"`
- `import "github.com/indigo-net/Brf.it/pkg/tokenizer"`

```go
type Options struct {
	// Path is the target path to scan.
	Path string

	// Version is the brf.it version string.
	Version string

	// Format is the output format ("xml" or "md").
	Format string

	// Output is the output file path (empty = stdout).
	Output string

	// IgnoreFile is the custom ignore file path.
	IgnoreFile string

	// IncludeHidden determines whether to include hidden files.
	IncludeHidden bool

	// IncludeBody determines whether to include function/method bodies.
	IncludeBody bool

	// IncludeImports determines whether to include import/export statements.
	IncludeImports bool

	// IncludeTree determines whether to include directory tree.
	IncludeTree bool

	// IncludePrivate determines whether to include private symbols.
	IncludePrivate bool

	// MaxFileSize is the maximum file size in bytes.
	MaxFileSize int64
}
func DefaultOptions() *Options
type Result struct {
	// Content is the formatted output bytes.
	Content []byte

	// TotalSignatures is the total number of signatures.
	TotalSignatures int

	// TotalFiles is the number of processed files.
	TotalFiles int

	// TotalSize is the total size of processed files.
	TotalSize int64

	// TokenCount is the number of tokens in the output.
	// Returns 0 if token counting is disabled or tokenizer is not set.
	TokenCount int
}
type Packager struct {
	scanner    scanner.Scanner
	extractor  extractor.Extractor
	formatters map[string]formatter.Formatter
	tokenizer  tokenizer.Tokenizer
}
func NewPackager(
	s scanner.Scanner,
	e extractor.Extractor,
	f map[string]formatter.Formatter,
) *Packager
func (p *Packager) SetTokenizer(t tokenizer.Tokenizer)
func (p *Packager) Package(opts *Options) (*Result, error)
treeStr string
func NewDefaultPackager(scanOpts *scanner.ScanOptions) (*Packager, error)
func normalizeFormat(format string) string
```

---

### /home/runner/work/Brf.it/Brf.it/internal/context/context_test.go

**Imports:**
- `import "strings"`
- `import "testing"`
- `import "github.com/indigo-net/Brf.it/pkg/extractor"`
- `import "github.com/indigo-net/Brf.it/pkg/formatter"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`
- `import "github.com/indigo-net/Brf.it/pkg/scanner"`
- `import "github.com/indigo-net/Brf.it/pkg/tokenizer"`

```go
type mockScanner struct {
	result *scanner.ScanResult
	err    error
}
func (m *mockScanner) Scan() (*scanner.ScanResult, error)
type mockExtractor struct {
	result *extractor.ExtractResult
	err    error
}
func (m *mockExtractor) Extract(_ *scanner.ScanResult, _ *extractor.ExtractOptions) (*extractor.ExtractResult, error)
func TestPackagerPackage(t *testing.T)
func TestPackagerPackageMarkdown(t *testing.T)
func TestPackagerPackageMarkdownFull(t *testing.T)
func TestPackagerUnknownFormat(t *testing.T)
func TestPackagerSetTokenizer(t *testing.T)
func TestPackagerWithTiktokenTokenizer(t *testing.T)
func TestPackagerTokenizerConsistency(t *testing.T)
func TestBuildTree(t *testing.T)
func TestBuildTreeStructure(t *testing.T)
func TestDefaultOptions(t *testing.T)
func TestNormalizeFormat(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/internal/context/tree.go

**Imports:**
- `import "path/filepath"`
- `import "sort"`
- `import "strings"`

```go
type treeNode struct {
	children map[string]*treeNode
}
func BuildTree(root string, paths []string) string
buf strings.Builder
func renderNode(buf *strings.Builder, n *treeNode, prefix string, isRoot bool)
newPrefix string
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/extractor/extractor.go

**Imports:**
- `import "fmt"`
- `import "os"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`
- `import "github.com/indigo-net/Brf.it/pkg/scanner"`

```go
type ExtractedFile struct {
	// Path is the file path.
	Path string

	// Language is the detected language.
	Language string

	// Signatures is the list of extracted signatures.
	Signatures []parser.Signature

	// Imports is the list of extracted import/export statements.
	Imports []parser.ImportExport

	// Size is the file size in bytes.
	Size int64

	// Error is any error that occurred during extraction.
	Error error
}
type ExtractResult struct {
	// Files is the list of extracted files.
	Files []ExtractedFile

	// TotalSignatures is the total number of signatures extracted.
	TotalSignatures int

	// TotalSize is the total size of processed files.
	TotalSize int64

	// ErrorCount is the number of files that had errors.
	ErrorCount int
}
type ExtractOptions struct {
	// IncludePrivate whether to include non-exported/private signatures.
	IncludePrivate bool

	// IncludeBody whether to include function/method bodies.
	IncludeBody bool

	// IncludeImports whether to include import/export statements.
	IncludeImports bool

	// Concurrency is the number of concurrent workers (0 = sequential).
	Concurrency int
}
type Extractor interface {
	// Extract extracts signatures from the given scan result.
	Extract(scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error)
}
type FileExtractor struct {
	registry *parser.Registry
}
func NewFileExtractor(registry *parser.Registry) *FileExtractor
func NewDefaultFileExtractor() *FileExtractor
func (e *FileExtractor) Extract(scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error)
func (e *FileExtractor) extractFile(entry scanner.FileEntry, opts *ExtractOptions) ExtractedFile
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/extractor/extractor_test.go

**Imports:**
- `import "os"`
- `import "path/filepath"`
- `import "testing"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`
- `import _ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"`
- `import "github.com/indigo-net/Brf.it/pkg/scanner"`

```go
func TestFileExtractorImplementsExtractor(t *testing.T)
_ Extractor = (*FileExtractor)(nil)
func TestFileExtractorExtract(t *testing.T)
foundAdd bool
func TestFileExtractorUnsupportedLanguage(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/formatter.go

**Imports:**
- `import "github.com/indigo-net/Brf.it/pkg/parser"`

```go
type FileData struct {
	// Path is the file path.
	Path string

	// Language is the detected language.
	Language string

	// Signatures is the list of extracted signatures.
	Signatures []parser.Signature

	// Imports is the list of extracted import/export statements.
	Imports []parser.ImportExport

	// Error is any error that occurred during extraction.
	Error error
}
type PackageData struct {
	// RootPath is the root path being packaged.
	RootPath string

	// Version is the brf.it version string.
	Version string

	// Tree is the directory tree string.
	Tree string

	// Files is the list of file data.
	Files []FileData

	// TotalSignatures is the total number of signatures.
	TotalSignatures int

	// TotalSize is the total size of processed files.
	TotalSize int64

	// IncludeImports indicates whether imports should be rendered.
	IncludeImports bool
}
type Formatter interface {
	// Format formats the package data and returns the output bytes.
	Format(data *PackageData) ([]byte, error)

	// Name returns the formatter name (e.g., "xml", "markdown").
	Name() string
}
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/formatter_test.go

**Imports:**
- `import "fmt"`
- `import "strings"`
- `import "testing"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`

```go
func TestXMLFormatterImplementsFormatter(t *testing.T)
_ Formatter = (*XMLFormatter)(nil)
func TestMarkdownFormatterImplementsFormatter(t *testing.T)
_ Formatter = (*MarkdownFormatter)(nil)
func TestXMLFormatterFormat(t *testing.T)
func TestXMLFormatterFormatWithError(t *testing.T)
func TestMarkdownFormatterFormat(t *testing.T)
func TestFormatterNames(t *testing.T)
func TestXMLFormatterEscapeXML(t *testing.T)
func TestMarkdownFormatterEscapeMarkdown(t *testing.T)
func TestXMLFormatterEmptyData(t *testing.T)
func TestMarkdownFormatterEmptyData(t *testing.T)
func TestMarkdownFormatterEmptyFile(t *testing.T)
func TestMarkdownFormatterEmptyFileWithImports(t *testing.T)
func TestXMLFormatterEmptyFile(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/helpers.go

```go
func getEmptyComment(lang string) string
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/helpers_test.go

**Imports:**
- `import "testing"`

```go
func TestGetEmptyComment(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/markdown.go

**Imports:**
- `import "bytes"`
- `import "fmt"`
- `import "strings"`

```go
type MarkdownFormatter struct{}
func NewMarkdownFormatter() *MarkdownFormatter
func (f *MarkdownFormatter) Name() string
func (f *MarkdownFormatter) Format(data *PackageData) ([]byte, error)
buf bytes.Buffer
imports, exports []string
func escapeMarkdown(s string) string
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/xml.go

**Imports:**
- `import "bytes"`
- `import "fmt"`
- `import "strings"`

```go
type XMLFormatter struct{}
func NewXMLFormatter() *XMLFormatter
func (f *XMLFormatter) Name() string
xmlSchemaComment = `<!--
Schema:
| Tag       | Description                              |
|-----------|------------------------------------------|
| file      | Source file (path, language attributes)  |
| signature | Function, type, or variable declaration  |
| imports   | Import statements container              |
| import    | Single import statement                  |
| export    | Single export statement                  |
| doc       | Documentation comment                    |
| error     | Parse error message                      |
-->
`
func (f *XMLFormatter) Format(data *PackageData) ([]byte, error)
buf bytes.Buffer
func escapeXML(s string) string
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/parser.go

**Imports:**
- `import "path/filepath"`
- `import "sync"`

```go
type ImportExport struct {
	// Type is "import" or "export".
	Type string

	// Path is the module path (e.g., "fmt", "react", "./utils").
	Path string

	// Name is the export name (for named exports).
	Name string

	// Line is the line number (1-indexed).
	Line int
}
type Signature struct {
	// Name is the identifier name (e.g., "Scan", "FileScanner").
	Name string

	// Kind is the type of signature (e.g., "function", "method", "class", "interface").
	Kind string

	// Text is the full signature text including parameters and return type.
	Text string

	// Doc is the documentation comment (if any).
	Doc string

	// Line is the starting line number (1-indexed).
	Line int

	// EndLine is the ending line number (1-indexed).
	EndLine int

	// Language is the source language (e.g., "go", "typescript").
	Language string

	// Exported indicates whether the signature is exported/public.
	Exported bool
}
type Node struct {
	// Type is the node type (e.g., "function_declaration", "class_definition").
	Type string

	// StartRow is the starting row (0-indexed).
	StartRow int

	// EndRow is the ending row (0-indexed).
	EndRow int

	// StartColumn is the starting column.
	StartColumn int

	// EndColumn is the ending column.
	EndColumn int

	// Text is the source text of the node.
	Text string

	// Children are child nodes.
	Children []Node
}
type ParseResult struct {
	// FilePath is the path to the parsed file.
	FilePath string

	// Language is the detected language.
	Language string

	// Signatures is the list of extracted signatures.
	Signatures []Signature

	// Imports is the list of extracted import/export statements.
	Imports []ImportExport

	// AST is the root node of the parsed AST (optional).
	AST *Node

	// Error is any error that occurred during parsing.
	Error error
}
type Options struct {
	// Language forces a specific language (auto-detected if empty).
	Language string

	// IncludeAST whether to include the full AST in the result.
	IncludeAST bool

	// IncludePrivate whether to include non-exported/private signatures.
	IncludePrivate bool

	// IncludeBody whether to include function/method bodies in the signature text.
	// When false (default), only the signature line is extracted.
	// When true, the full declaration including the body is extracted.
	IncludeBody bool

	// IncludeImports whether to include import/export statements in the result.
	IncludeImports bool
}
type Parser interface {
	// Parse parses the given content and returns extracted signatures.
	Parse(content string, opts *Options) (*ParseResult, error)

	// Languages returns the list of supported languages.
	Languages() []string
}
type Registry struct {
	mu      sync.RWMutex
	parsers map[string]Parser
}
func NewRegistry() *Registry
defaultRegistry = NewRegistry()
func DefaultRegistry() *Registry
func (r *Registry) Register(lang string, parser Parser)
func (r *Registry) Get(lang string) (Parser, bool)
func (r *Registry) Languages() []string
func RegisterParser(lang string, parser Parser)
func GetParser(lang string) (Parser, bool)
LanguageMapping = map[string]string{
	".go":    "go",
	".ts":    "typescript",
	".tsx":   "tsx",
	".js":    "javascript",
	".jsx":   "jsx",
	".py":    "python",
	".java":  "java",
	".rs":    "rust",
	".rb":    "ruby",
	".php":   "php",
	".c":     "c",
	".cpp":   "cpp",
	".h":     "c",
	".hpp":   "cpp",
	".cs":    "csharp",
	".swift": "swift",
	".kt":    "kotlin",
}
func DetectLanguage(path string) string
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/parser_test.go

**Imports:**
- `import "testing"`

```go
func TestSignatureDefaults(t *testing.T)
func TestParseResultDefaults(t *testing.T)
func TestNodeKind(t *testing.T)
func TestParserInterface(t *testing.T)
_ Parser = (*MockParser)(nil)
type MockParser struct {
	signatures []Signature
	err        error
}
func (m *MockParser) Parse(content string, opts *Options) (*ParseResult, error)
func (m *MockParser) Languages() []string
func TestMockParser(t *testing.T)
func TestRegistry(t *testing.T)
func TestDefaultRegistry(t *testing.T)
func TestDetectLanguage(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/c.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"`

```go
type CQuery struct {
	language *sitter.Language
	query    []byte
}
func NewCQuery() *CQuery
func (q *CQuery) Language() *sitter.Language
func (q *CQuery) Query() []byte
func (q *CQuery) Captures() []string
func (q *CQuery) KindMapping() map[string]string
func (q *CQuery) ImportQuery() []byte
cImportQueryPattern = `
; #include directives (capture full statement)
(preproc_include) @import_path
`
cQueryPattern = `
; Function definitions - direct declarator
(function_definition
  declarator: (function_declarator
    declarator: (identifier) @name
  )
) @signature @kind

; Function definitions - pointer return type
(function_definition
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Function declarations (prototypes) - direct declarator
(declaration
  declarator: (function_declarator
    declarator: (identifier) @name
  )
) @signature @kind

; Function declarations (prototypes) - pointer return type
(declaration
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Struct specifiers
(struct_specifier
  name: (type_identifier) @name
) @signature @kind

; Enum specifiers
(enum_specifier
  name: (type_identifier) @name
) @signature @kind

; Typedef
(type_definition
  declarator: (type_identifier) @name
) @signature @kind

; Function-like macros
(preproc_function_def
  name: (identifier) @name
) @signature @kind

; Object-like macros
(preproc_def
  name: (identifier) @name
) @signature @kind

; Global variable declarations (with initializer)
(translation_unit
  (declaration
    declarator: (init_declarator
      declarator: (identifier) @name
    )
  ) @signature @kind
)

; Global variable declarations (simple identifier, e.g., extern)
(translation_unit
  (declaration
    declarator: (identifier) @name
  ) @signature @kind
)

; Global pointer variable declarations
(translation_unit
  (declaration
    declarator: (pointer_declarator
      declarator: (identifier) @name
    )
  ) @signature @kind
)

; Global pointer variable declarations (with initializer)
(translation_unit
  (declaration
    declarator: (init_declarator
      declarator: (pointer_declarator
        declarator: (identifier) @name
      )
    )
  ) @signature @kind
)

; Comments
(comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/c_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"`

```go
func TestCQueryLanguage(t *testing.T)
func TestCQueryPattern(t *testing.T)
func TestCQueryExtractFunction(t *testing.T)
funcCaptures map[string]string
func TestCQueryExtractStruct(t *testing.T)
func TestCQueryExtractMacro(t *testing.T)
func TestCQueryExtractEnum(t *testing.T)
func TestCQueryExtractTypedef(t *testing.T)
func TestCQueryExtractGlobalVariables(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/cpp.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"`

```go
type CppQuery struct {
	language *sitter.Language
	query    []byte
}
func NewCppQuery() *CppQuery
func (q *CppQuery) Language() *sitter.Language
func (q *CppQuery) Query() []byte
func (q *CppQuery) Captures() []string
func (q *CppQuery) KindMapping() map[string]string
func (q *CppQuery) ImportQuery() []byte
cppImportQueryPattern = `
; #include directives (capture full statement)
(preproc_include) @import_path
`
cppQueryPattern = `
; Function definitions - direct declarator
(function_definition
  declarator: (function_declarator
    declarator: (identifier) @name
  )
) @signature @kind

; Function definitions - pointer return type
(function_definition
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Function definitions - reference return type
(function_definition
  declarator: (reference_declarator
    (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Function declarations (prototypes) - direct declarator
(declaration
  declarator: (function_declarator
    declarator: (identifier) @name
  )
) @signature @kind

; Function declarations (prototypes) - pointer return type
(declaration
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Class definitions
(class_specifier
  name: (type_identifier) @name
) @signature @kind

; Struct specifiers
(struct_specifier
  name: (type_identifier) @name
) @signature @kind

; Enum specifiers
(enum_specifier
  name: (type_identifier) @name
) @signature @kind

; Typedef
(type_definition
  declarator: (type_identifier) @name
) @signature @kind

; Function-like macros
(preproc_function_def
  name: (identifier) @name
) @signature @kind

; Object-like macros
(preproc_def
  name: (identifier) @name
) @signature @kind

; Method declarations in class (regular methods)
(field_declaration
  declarator: (function_declarator
    declarator: (field_identifier) @name
  )
) @signature @kind

; Method declarations with pointer return type
(field_declaration
  declarator: (pointer_declarator
    declarator: (function_declarator
      declarator: (field_identifier) @name
    )
  )
) @signature @kind

; Method declarations with reference return type
(field_declaration
  declarator: (reference_declarator
    (function_declarator
      declarator: (field_identifier) @name
    )
  )
) @signature @kind

; Constructor declarations (in class body)
(function_definition
  declarator: (function_declarator
    declarator: (qualified_identifier
      name: (identifier) @name
    )
  )
) @signature @kind

; Destructor definitions (outside class)
(function_definition
  declarator: (function_declarator
    declarator: (destructor_name
      (identifier) @name
    )
  )
) @signature @kind

; Destructor declarations in class (captured via declaration node)
(declaration
  declarator: (function_declarator
    declarator: (destructor_name
      (identifier) @name
    )
  )
) @signature @kind

; Namespace definitions
(namespace_definition
  name: (namespace_identifier) @name
) @signature @kind

; Template function definitions
(template_declaration
  (function_definition
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Template function definitions - pointer return type
(template_declaration
  (function_definition
    declarator: (pointer_declarator
      declarator: (function_declarator
        declarator: (identifier) @name
      )
    )
  )
) @signature @kind

; Template class definitions
(template_declaration
  (class_specifier
    name: (type_identifier) @name
  )
) @signature @kind

; Template struct definitions
(template_declaration
  (struct_specifier
    name: (type_identifier) @name
  )
) @signature @kind

; Template declarations (standalone)
(template_declaration
  (declaration
    declarator: (function_declarator
      declarator: (identifier) @name
    )
  )
) @signature @kind

; Comments
(comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/cpp_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"`

```go
func TestCppQueryLanguage(t *testing.T)
func TestCppQueryPattern(t *testing.T)
func TestCppQueryExtractFunction(t *testing.T)
funcCaptures map[string]string
func TestCppQueryExtractClass(t *testing.T)
func TestCppQueryExtractMethod(t *testing.T)
func TestCppQueryExtractConstructorDestructor(t *testing.T)
func TestCppQueryExtractNamespace(t *testing.T)
func TestCppQueryExtractTemplate(t *testing.T)
func TestCppQueryExtractStruct(t *testing.T)
func TestCppQueryExtractEnum(t *testing.T)
func TestCppQueryExtractMacro(t *testing.T)
func TestCppQueryExtractTypedef(t *testing.T)
func TestCppQueryExtractIncludes(t *testing.T)
imports []string
func TestCppQueryNestedNamespaces(t *testing.T)
func TestCppQueryMultipleInheritance(t *testing.T)
func TestCppQueryEmptyFile(t *testing.T)
func TestCppQueryOnlyComments(t *testing.T)
nameCount int
docCount int
func TestCppQueryKindMapping(t *testing.T)
func TestCppQueryCaptures(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/go.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"`

```go
captureName      = "name"
captureSignature = "signature"
captureDoc       = "doc"
captureKind      = "kind"
type GoQuery struct {
	language *sitter.Language
	query    []byte
}
func NewGoQuery() *GoQuery
func (q *GoQuery) Language() *sitter.Language
func (q *GoQuery) Query() []byte
func (q *GoQuery) Captures() []string
func (q *GoQuery) KindMapping() map[string]string
func (q *GoQuery) ImportQuery() []byte
goImportQueryPattern = `
; Single import (capture full spec including alias)
(import_declaration
  (import_spec) @import_path
)

; Multi-line imports (capture each spec)
(import_declaration
  (import_spec_list
    (import_spec) @import_path
  )
)
`
goQueryPattern = `
; Function declarations
(function_declaration
  name: (identifier) @name
) @signature @kind

; Method declarations
(method_declaration
  name: (field_identifier) @name
) @signature @kind

; Type declarations (struct, interface, etc.)
(type_declaration
  (type_spec
    name: (type_identifier) @name
  )
) @signature @kind

; Package-level const specs (captures each const individually)
(const_spec
  name: (identifier) @name
) @signature @kind

; Package-level var specs (captures each var individually)
(var_spec
  name: (identifier) @name
) @signature @kind

; Comments (documentation)
(comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/go_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"`

```go
func TestGoQueryLanguage(t *testing.T)
func TestGoQueryPattern(t *testing.T)
func TestGoQueryExtractFunction(t *testing.T)
funcCaptures map[string]string
funcKindNode *sitter.Node
kindNode *sitter.Node
func TestGoQueryExtractConstAndVar(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/java.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"`

```go
type JavaQuery struct {
	language *sitter.Language
	query    []byte
}
func NewJavaQuery() *JavaQuery
func (q *JavaQuery) Language() *sitter.Language
func (q *JavaQuery) Query() []byte
func (q *JavaQuery) Captures() []string
func (q *JavaQuery) KindMapping() map[string]string
func (q *JavaQuery) ImportQuery() []byte
javaImportQueryPattern = `
; import statements (capture full declaration)
(import_declaration) @import_path
`
javaQueryPattern = `
; Class declarations (includes inner classes)
(class_declaration
  name: (identifier) @name
) @signature @kind

; Interface declarations
(interface_declaration
  name: (identifier) @name
) @signature @kind

; Method declarations
(method_declaration
  name: (identifier) @name
) @signature @kind

; Constructor declarations
(constructor_declaration
  name: (identifier) @name
) @signature @kind

; Enum declarations
(enum_declaration
  name: (identifier) @name
) @signature @kind

; Annotation type declarations (@interface)
(annotation_type_declaration
  name: (identifier) @name
) @signature @kind

; Record declarations (Java 14+)
(record_declaration
  name: (identifier) @name
) @signature @kind

; Field declarations (static fields filtered in parser.go)
(field_declaration
  (variable_declarator
    name: (identifier) @name
  )
) @signature @kind

; Comments (Javadoc and regular)
(line_comment) @doc
(block_comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/java_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"`

```go
func TestJavaQueryLanguage(t *testing.T)
func TestJavaQueryPattern(t *testing.T)
func TestJavaQueryKindMapping(t *testing.T)
func TestJavaQueryExtractClass(t *testing.T)
foundClass, foundMethod bool
func TestJavaQueryExtractInterface(t *testing.T)
func TestJavaQueryExtractEnum(t *testing.T)
foundEnum bool
func TestJavaQueryExtractAnnotationType(t *testing.T)
foundAnnotation bool
func TestJavaQueryExtractRecord(t *testing.T)
func TestJavaQueryExtractGenerics(t *testing.T)
func TestJavaQueryExtractFieldDeclarations(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/python.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"`

```go
type PythonQuery struct {
	language *sitter.Language
	query    []byte
}
func NewPythonQuery() *PythonQuery
func (q *PythonQuery) Language() *sitter.Language
func (q *PythonQuery) Query() []byte
func (q *PythonQuery) Captures() []string
func (q *PythonQuery) KindMapping() map[string]string
func (q *PythonQuery) ImportQuery() []byte
pythonImportQueryPattern = `
; import module (capture full statement)
(import_statement) @import_path

; from module import ... (capture full statement)
(import_from_statement) @import_path
`
pythonQueryPattern = `
; Function definitions (includes async def, methods)
(function_definition
  name: (identifier) @name
) @signature @kind

; Class definitions
(class_definition
  name: (identifier) @name
) @signature @kind

; Module-level assignments (simple and with type annotations)
(module
  (expression_statement
    (assignment
      left: (identifier) @name
    )
  ) @signature @kind
)

; Comments
(comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/python_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"`

```go
func TestPythonQueryLanguage(t *testing.T)
func TestPythonQueryPattern(t *testing.T)
func TestPythonQueryExtractFunction(t *testing.T)
funcCaptures map[string]string
func TestPythonQueryExtractClass(t *testing.T)
func TestPythonQueryExtractAsyncFunction(t *testing.T)
funcCaptures map[string]string
func TestPythonQueryExtractModuleLevelVariables(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/typescript.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"`

```go
type TypeScriptQuery struct {
	language *sitter.Language
	query    []byte
}
func NewTypeScriptQuery() *TypeScriptQuery
func (q *TypeScriptQuery) Language() *sitter.Language
func (q *TypeScriptQuery) Query() []byte
func (q *TypeScriptQuery) Captures() []string
func (q *TypeScriptQuery) KindMapping() map[string]string
func (q *TypeScriptQuery) ImportQuery() []byte
typeScriptImportQueryPattern = `
; Import statements (capture full statement)
(import_statement) @import_path

; Export statements with source (re-exports)
(export_statement
  source: (string)
) @import_path @export_type

; Named exports without source (local exports)
(export_statement
  declaration: (_
    name: (identifier) @export_name
  )
)

; Export clause (export { foo, bar })
(export_statement
  (export_clause
    (export_specifier
      name: (identifier) @export_name
    )
  )
)
`
typeScriptQueryPattern = `
; Function declarations
(function_declaration
  name: (identifier) @name
) @signature @kind

; Exported function declarations
(export_statement
  (function_declaration
    name: (identifier) @name
  )
) @signature @kind

; Arrow functions in variable declarations (capture full declaration with const/let/var)
(lexical_declaration
  (variable_declarator
    name: (identifier) @name
    value: (arrow_function)
  )
) @signature @kind

; Module-level const/let declarations with values (captures all module-level)
; Deduplication for arrow functions is handled in parser.go
(program
  (lexical_declaration
    (variable_declarator
      name: (identifier) @name
      value: (_)
    )
  ) @signature @kind
)

; Module-level const/let without initial value (TypeScript declares)
(program
  (lexical_declaration
    (variable_declarator
      name: (identifier) @name
      !value
    )
  ) @signature @kind
)

; Exported module-level const/let declarations with values
(program
  (export_statement
    declaration: (lexical_declaration
      (variable_declarator
        name: (identifier) @name
        value: (_)
      )
    )
  ) @signature @kind
)

; Exported module-level const/let without initial value
(program
  (export_statement
    declaration: (lexical_declaration
      (variable_declarator
        name: (identifier) @name
        !value
      )
    )
  ) @signature @kind
)

; Method definitions
(method_definition
  name: (property_identifier) @name
) @signature @kind

; Class declarations
(class_declaration
  name: (type_identifier) @name
) @signature @kind

; Interface declarations
(interface_declaration
  name: (type_identifier) @name
) @signature @kind

; Type alias declarations
(type_alias_declaration
  name: (type_identifier) @name
) @signature @kind

; Comments (documentation)
(comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/typescript_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"`

```go
func TestTypeScriptQueryLanguage(t *testing.T)
func TestTypeScriptQueryPattern(t *testing.T)
func TestTypeScriptQueryExtractFunction(t *testing.T)
funcCaptures map[string]string
func TestTypeScriptQueryExtractModuleLevelVariables(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/parser.go

**Imports:**
- `import "fmt"`
- `import "regexp"`
- `import "strings"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`
- `import "github.com/indigo-net/Brf.it/pkg/parser/treesitter/languages"`

```go
func init()
type TreeSitterParser struct {
	queries map[string]LanguageQuery
}
func NewTreeSitterParser() *TreeSitterParser
func (p *TreeSitterParser) Parse(content string, opts *parser.Options) (*parser.ParseResult, error)
imports []parser.ImportExport
func (p *TreeSitterParser) Languages() []string
func (p *TreeSitterParser) extractSignatures(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
) []parser.Signature
signatures []parser.Signature
kindNode *sitter.Node
func cleanComment(text string) string
func isExported(name, language string) bool
func stripBody(text, kind, language string) string
func stripGoBody(text, kind string) string
tsFunctionBodyRe = regexp.MustCompile(`\s*\{[\s\S]*\}\s*$`)
tsArrowBodyRe = regexp.MustCompile(`\s*=>\s*[\s\S]+$`)
tsClassBodyRe = regexp.MustCompile(`\s*\{[\s\S]*\}\s*$`)
func stripTypeScriptBody(text, kind string) string
func stripTSFunctionBody(text string) string
func findFunctionBodyStart(text string) int
func findTSClassBodyStart(text string) int
func stripPythonBody(text, kind string) string
func findPythonBodyStart(text string) int
func stripCBody(text, kind string) string
func stripCppBody(text, kind string) string
func findCppBodyStart(text string) int
func isPythonMethod(signature string) bool
func stripJavaBody(text, kind string) string
func findJavaBodyStart(text string) int
func (p *TreeSitterParser) extractImports(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
) []parser.ImportExport
imports []parser.ImportExport
imp parser.ImportExport
hasExportType bool
func cleanImportPath(path string) string
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/parser_test.go

**Imports:**
- `import "strings"`
- `import "testing"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`

```go
func TestTreeSitterParserImplementsParser(t *testing.T)
_ parser.Parser = (*TreeSitterParser)(nil)
func TestTreeSitterParserLanguages(t *testing.T)
func TestTreeSitterParserParseGo(t *testing.T)
foundAdd bool
func TestTreeSitterParserParseTypeScript(t *testing.T)
foundAdd bool
func TestTreeSitterParserUnsupportedLanguage(t *testing.T)
func TestTreeSitterParserAutoRegistration(t *testing.T)
func TestGoSignatureOnlyExtraction(t *testing.T)
func TestGoIncludeBodyExtraction(t *testing.T)
foundAdd bool
func TestTypeScriptSignatureOnlyExtraction(t *testing.T)
func TestTypeScriptArrowFunctionSignature(t *testing.T)
func contains(s, substr string) bool
func TestTreeSitterParserParseJava(t *testing.T)
foundClass, foundConstructor, foundPublicMethod, foundPrivateMethod bool
func TestJavaSignatureOnlyExtraction(t *testing.T)
func TestJavaGenericsExtraction(t *testing.T)
foundClass, foundMethod bool
func TestJavaAutoRegistration(t *testing.T)
func TestTreeSitterParserParseCpp(t *testing.T)
func TestCppSignatureOnlyExtraction(t *testing.T)
func TestCppTemplateExtraction(t *testing.T)
func TestCppAutoRegistration(t *testing.T)
func TestCppImportExtraction(t *testing.T)
func TestGoVariableExtraction(t *testing.T)
func TestTypeScriptVariableExtraction(t *testing.T)
func TestPythonVariableExtraction(t *testing.T)
func TestJavaStaticFieldExtraction(t *testing.T)
func TestCGlobalVariableExtraction(t *testing.T)
func TestVariableSignaturePreservesValue(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/query.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`

```go
type LanguageQuery interface {
	// Language returns the Tree-sitter language for parsing.
	Language() *sitter.Language

	// Query returns the Tree-sitter query pattern for signature extraction.
	Query() []byte

	// ImportQuery returns the Tree-sitter query pattern for import/export extraction.
	// Returns nil if the language doesn't support import extraction.
	ImportQuery() []byte

	// Captures returns the list of capture names used in the query.
	Captures() []string

	// KindMapping maps Tree-sitter node types to Signature kinds.
	KindMapping() map[string]string
}
CaptureName      = "name"
CaptureSignature = "signature"
CaptureDoc       = "doc"
CaptureKind      = "kind"
CaptureImportPath = "import_path"
CaptureExportName = "export_name"
CaptureImportType = "import_type"
DefaultKindMapping = map[string]string{
	"function_declaration": "function",
	"method_declaration":   "method",
	"type_declaration":     "type",
	"struct_type":          "struct",
	"interface_type":       "interface",
	"class_declaration":    "class",
	"arrow_function":       "function",
	"function_expression":  "function",
	"method_definition":    "method",
}
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/query_test.go

**Imports:**
- `import "testing"`

```go
func TestCaptureDefinitions(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/scanner/scanner.go

**Imports:**
- `import "io/fs"`
- `import "log"`
- `import "os"`
- `import "path/filepath"`
- `import "strings"`
- `import ignore "github.com/sabhiram/go-gitignore"`

```go
type FileEntry struct {
	// Path is the absolute or relative path to the file.
	Path string

	// Language is the detected programming language (e.g., "go", "typescript").
	Language string

	// Size is the file size in bytes.
	Size int64
}
type ScanResult struct {
	// Files is the list of matched files.
	Files []FileEntry

	// TotalSize is the sum of all matched file sizes.
	TotalSize int64

	// SkippedCount is the number of files skipped (too large, unsupported, etc.).
	SkippedCount int
}
type ScanOptions struct {
	// RootPath is the directory or file to scan.
	RootPath string

	// SupportedExtensions maps file extensions to language names.
	SupportedExtensions map[string]string

	// IgnoreFile is the path to the gitignore file (default: .gitignore).
	IgnoreFile string

	// IncludeHidden determines whether to include hidden files (dotfiles).
	IncludeHidden bool

	// MaxFileSize is the maximum file size in bytes to include.
	MaxFileSize int64
}
func DefaultScanOptions() *ScanOptions
func (o *ScanOptions) GetLanguage(path string) (string, bool)
func IsHidden(name string) bool
type Scanner interface {
	// Scan performs the scan and returns scan results.
	Scan() (*ScanResult, error)
}
type FileScanner struct {
	opts       *ScanOptions
	ignorer    *ignore.GitIgnore
	ignorerErr error
	logger     *log.Logger
}
func NewFileScanner(opts *ScanOptions) (*FileScanner, error)
func (s *FileScanner) Scan() (*ScanResult, error)
func (s *FileScanner) checkFile(path string, info os.FileInfo) (FileEntry, bool)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/scanner/scanner_test.go

**Imports:**
- `import "os"`
- `import "path/filepath"`
- `import "testing"`

```go
func TestNewFileScanner(t *testing.T)
func TestNewFileScannerNilOptions(t *testing.T)
func TestFileEntryDefaults(t *testing.T)
func TestScanOptionsDefaults(t *testing.T)
expectedMaxSize = 512000
func TestScanOptionsWithExtensions(t *testing.T)
func TestScannerInterface(t *testing.T)
_ Scanner = (*FileScanner)(nil)
func TestScanEmptyDirectory(t *testing.T)
func TestScanSingleFile(t *testing.T)
func TestScanFilterByExtension(t *testing.T)
func TestScanExcludeHidden(t *testing.T)
func TestScanIncludeHidden(t *testing.T)
func TestScanMaxFileSize(t *testing.T)
func TestScanGitignore(t *testing.T)
func TestScanNestedDirectories(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/tokenizer/tiktoken.go

**Imports:**
- `import "github.com/pkoukk/tiktoken-go"`

```go
type TiktokenTokenizer struct {
	encoding string
	tke      *tiktoken.Tiktoken
}
_ Tokenizer = (*TiktokenTokenizer)(nil)
func NewTiktokenTokenizer() (*TiktokenTokenizer, error)
func (t *TiktokenTokenizer) Count(text string) (int, error)
func (t *TiktokenTokenizer) Name() string
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/tokenizer/tokenizer.go

```go
type Tokenizer interface {
	// Count returns the number of tokens in the given text.
	// Returns 0 and error if counting fails.
	Count(text string) (int, error)

	// Name returns the tokenizer name (e.g., "tiktoken-cl100k", "noop").
	Name() string
}
type NoOpTokenizer struct{}
_ Tokenizer = (*NoOpTokenizer)(nil)
func NewNoOpTokenizer() *NoOpTokenizer
func (t *NoOpTokenizer) Count(_ string) (int, error)
func (t *NoOpTokenizer) Name() string
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/tokenizer/tokenizer_test.go

**Imports:**
- `import "strings"`
- `import "testing"`

```go
func TestNoOpTokenizerImplementsTokenizer(t *testing.T)
_ Tokenizer = (*NoOpTokenizer)(nil)
func TestTiktokenTokenizerImplementsTokenizer(t *testing.T)
_ Tokenizer = (*TiktokenTokenizer)(nil)
func TestNoOpTokenizerCount(t *testing.T)
func TestNoOpTokenizerName(t *testing.T)
func TestTiktokenTokenizerCount(t *testing.T)
func TestTiktokenTokenizerName(t *testing.T)
func TestTiktokenTokenizerConsistency(t *testing.T)
func TestTiktokenTokenizerSpecialCharacters(t *testing.T)
```

---

