# Code Summary: /home/runner/work/Brf.it/Brf.it

*brf.it 0.13.0*

---

## Files

### /home/runner/work/Brf.it/Brf.it/cmd/brfit/main.go

```go
version = "dev" // variable
commit  = "none" // variable
date    = "unknown" // variable
func main() // function
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
Version = "dev" // variable
Commit  = "none" // variable
Date    = "unknown" // variable
func SetBuildInfo(v, c, d string) // function
cfg *config.Config // variable
rootCmd *cobra.Command // variable
func init() // function
func Execute() // function
func NewRootCommand() *cobra.Command // function
func newRootCommandWithConfig(c *config.Config) *cobra.Command // function
func addFlags(cmd *cobra.Command, c *config.Config) // function
func runRoot(cmd *cobra.Command, args []string, c *config.Config) error // function
func writeOutput(result *context.Result, c *config.Config) error // function
func writeToFile(path string, content []byte) error // function
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
func TestExecuteHelp(t *testing.T) // function
func TestExecuteVersion(t *testing.T) // function
buf bytes.Buffer // variable
func TestNewRootCommand(t *testing.T) // function
func TestParseFlags(t *testing.T) // function
func TestRootCommandIntegration(t *testing.T) // function
buf bytes.Buffer // variable
func TestRootCommandIntegrationMarkdown(t *testing.T) // function
buf bytes.Buffer // variable
func TestRootCommandPathNotFound(t *testing.T) // function
func TestWriteToFile(t *testing.T) // function
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
} // type
func DefaultConfig() *Config // function
func (c *Config) Validate() error // method
func (c *Config) SupportedExtensions() map[string]string // method
func (c *Config) ToOptions() *pkgcontext.Options // method
```

---

### /home/runner/work/Brf.it/Brf.it/internal/config/config_test.go

**Imports:**
- `import "testing"`

```go
func TestDefaultConfig(t *testing.T) // function
expectedMaxSize = 512000 // variable
func TestConfigValidate(t *testing.T) // function
func TestConfigSupportedLanguages(t *testing.T) // function
func containsString(s, substr string) bool // function
func containsSubstring(s, substr string) bool // function
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
} // type
func DefaultOptions() *Options // function
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
} // type
type Packager struct {
	scanner    scanner.Scanner
	extractor  extractor.Extractor
	formatters map[string]formatter.Formatter
	tokenizer  tokenizer.Tokenizer
} // type
func NewPackager(
	s scanner.Scanner,
	e extractor.Extractor,
	f map[string]formatter.Formatter,
) *Packager // function
func (p *Packager) SetTokenizer(t tokenizer.Tokenizer) // method
func (p *Packager) Package(opts *Options) (*Result, error) // method
treeStr string // variable
func NewDefaultPackager(scanOpts *scanner.ScanOptions) (*Packager, error) // function
func normalizeFormat(format string) string // function
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
} // type
func (m *mockScanner) Scan() (*scanner.ScanResult, error) // method
type mockExtractor struct {
	result *extractor.ExtractResult
	err    error
} // type
func (m *mockExtractor) Extract(_ *scanner.ScanResult, _ *extractor.ExtractOptions) (*extractor.ExtractResult, error) // method
func TestPackagerPackage(t *testing.T) // function
func TestPackagerPackageMarkdown(t *testing.T) // function
func TestPackagerPackageMarkdownFull(t *testing.T) // function
func TestPackagerUnknownFormat(t *testing.T) // function
func TestPackagerSetTokenizer(t *testing.T) // function
func TestPackagerWithTiktokenTokenizer(t *testing.T) // function
func TestPackagerTokenizerConsistency(t *testing.T) // function
func TestBuildTree(t *testing.T) // function
func TestBuildTreeStructure(t *testing.T) // function
func TestDefaultOptions(t *testing.T) // function
func TestNormalizeFormat(t *testing.T) // function
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
} // type
func BuildTree(root string, paths []string) string // function
buf strings.Builder // variable
func renderNode(buf *strings.Builder, n *treeNode, prefix string, isRoot bool) // function
newPrefix string // variable
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
} // type
type ExtractResult struct {
	// Files is the list of extracted files.
	Files []ExtractedFile

	// TotalSignatures is the total number of signatures extracted.
	TotalSignatures int

	// TotalSize is the total size of processed files.
	TotalSize int64

	// ErrorCount is the number of files that had errors.
	ErrorCount int
} // type
type ExtractOptions struct {
	// IncludePrivate whether to include non-exported/private signatures.
	IncludePrivate bool

	// IncludeBody whether to include function/method bodies.
	IncludeBody bool

	// IncludeImports whether to include import/export statements.
	IncludeImports bool

	// Concurrency is the number of concurrent workers (0 = sequential).
	Concurrency int
} // type
type Extractor interface {
	// Extract extracts signatures from the given scan result.
	Extract(scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error)
} // type
type FileExtractor struct {
	registry *parser.Registry
} // type
func NewFileExtractor(registry *parser.Registry) *FileExtractor // function
func NewDefaultFileExtractor() *FileExtractor // function
func (e *FileExtractor) Extract(scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error) // method
func (e *FileExtractor) extractFile(entry scanner.FileEntry, opts *ExtractOptions) ExtractedFile // method
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
func TestFileExtractorImplementsExtractor(t *testing.T) // function
_ Extractor = (*FileExtractor)(nil) // variable
func TestFileExtractorExtract(t *testing.T) // function
foundAdd bool // variable
func TestFileExtractorUnsupportedLanguage(t *testing.T) // function
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
} // type
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
} // type
type Formatter interface {
	// Format formats the package data and returns the output bytes.
	Format(data *PackageData) ([]byte, error)

	// Name returns the formatter name (e.g., "xml", "markdown").
	Name() string
} // type
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/formatter_test.go

**Imports:**
- `import "fmt"`
- `import "strings"`
- `import "testing"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`

```go
func TestXMLFormatterImplementsFormatter(t *testing.T) // function
_ Formatter = (*XMLFormatter)(nil) // variable
func TestMarkdownFormatterImplementsFormatter(t *testing.T) // function
_ Formatter = (*MarkdownFormatter)(nil) // variable
func TestXMLFormatterFormat(t *testing.T) // function
func TestXMLFormatterFormatWithError(t *testing.T) // function
func TestMarkdownFormatterFormat(t *testing.T) // function
func TestFormatterNames(t *testing.T) // function
func TestXMLFormatterEscapeXML(t *testing.T) // function
func TestMarkdownFormatterEscapeMarkdown(t *testing.T) // function
func TestXMLFormatterEmptyData(t *testing.T) // function
func TestMarkdownFormatterEmptyData(t *testing.T) // function
func TestMarkdownFormatterEmptyFile(t *testing.T) // function
func TestMarkdownFormatterEmptyFileWithImports(t *testing.T) // function
func TestXMLFormatterEmptyFile(t *testing.T) // function
func TestKindToTag(t *testing.T) // function
func TestXMLFormatterKindTags(t *testing.T) // function
func TestMarkdownFormatterKindComment(t *testing.T) // function
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/helpers.go

```go
func getEmptyComment(lang string) string // function
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/helpers_test.go

**Imports:**
- `import "testing"`

```go
func TestGetEmptyComment(t *testing.T) // function
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/markdown.go

**Imports:**
- `import "bytes"`
- `import "fmt"`
- `import "strings"`

```go
type MarkdownFormatter struct{} // type
func NewMarkdownFormatter() *MarkdownFormatter // function
func (f *MarkdownFormatter) Name() string // method
func (f *MarkdownFormatter) Format(data *PackageData) ([]byte, error) // method
buf bytes.Buffer // variable
imports, exports []string // variable
func escapeMarkdown(s string) string // function
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/xml.go

**Imports:**
- `import "bytes"`
- `import "fmt"`
- `import "strings"`

```go
type XMLFormatter struct{} // type
func NewXMLFormatter() *XMLFormatter // function
func (f *XMLFormatter) Name() string // method
func (f *XMLFormatter) Format(data *PackageData) ([]byte, error) // method
buf bytes.Buffer // variable
func escapeXML(s string) string // function
func kindToTag(kind string) string // function
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
} // type
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
} // type
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
} // type
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
} // type
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
} // type
type Parser interface {
	// Parse parses the given content and returns extracted signatures.
	Parse(content string, opts *Options) (*ParseResult, error)

	// Languages returns the list of supported languages.
	Languages() []string
} // type
type Registry struct {
	mu      sync.RWMutex
	parsers map[string]Parser
} // type
func NewRegistry() *Registry // function
defaultRegistry = NewRegistry() // variable
func DefaultRegistry() *Registry // function
func (r *Registry) Register(lang string, parser Parser) // method
func (r *Registry) Get(lang string) (Parser, bool) // method
func (r *Registry) Languages() []string // method
func RegisterParser(lang string, parser Parser) // function
func GetParser(lang string) (Parser, bool) // function
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
} // variable
func DetectLanguage(path string) string // function
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/parser_test.go

**Imports:**
- `import "testing"`

```go
func TestSignatureDefaults(t *testing.T) // function
func TestParseResultDefaults(t *testing.T) // function
func TestNodeKind(t *testing.T) // function
func TestParserInterface(t *testing.T) // function
_ Parser = (*MockParser)(nil) // variable
type MockParser struct {
	signatures []Signature
	err        error
} // type
func (m *MockParser) Parse(content string, opts *Options) (*ParseResult, error) // method
func (m *MockParser) Languages() []string // method
func TestMockParser(t *testing.T) // function
func TestRegistry(t *testing.T) // function
func TestDefaultRegistry(t *testing.T) // function
func TestDetectLanguage(t *testing.T) // function
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
} // type
func NewCQuery() *CQuery // function
func (q *CQuery) Language() *sitter.Language // method
func (q *CQuery) Query() []byte // method
func (q *CQuery) Captures() []string // method
func (q *CQuery) KindMapping() map[string]string // method
func (q *CQuery) ImportQuery() []byte // method
cImportQueryPattern = `
; #include directives (capture full statement)
(preproc_include) @import_path
` // variable
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
` // variable
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/c_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"`

```go
func TestCQueryLanguage(t *testing.T) // function
func TestCQueryPattern(t *testing.T) // function
func TestCQueryExtractFunction(t *testing.T) // function
funcCaptures map[string]string // variable
func TestCQueryExtractStruct(t *testing.T) // function
func TestCQueryExtractMacro(t *testing.T) // function
func TestCQueryExtractEnum(t *testing.T) // function
func TestCQueryExtractTypedef(t *testing.T) // function
func TestCQueryExtractGlobalVariables(t *testing.T) // function
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
} // type
func NewCppQuery() *CppQuery // function
func (q *CppQuery) Language() *sitter.Language // method
func (q *CppQuery) Query() []byte // method
func (q *CppQuery) Captures() []string // method
func (q *CppQuery) KindMapping() map[string]string // method
func (q *CppQuery) ImportQuery() []byte // method
cppImportQueryPattern = `
; #include directives (capture full statement)
(preproc_include) @import_path
` // variable
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
` // variable
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/cpp_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"`

```go
func TestCppQueryLanguage(t *testing.T) // function
func TestCppQueryPattern(t *testing.T) // function
func TestCppQueryExtractFunction(t *testing.T) // function
funcCaptures map[string]string // variable
func TestCppQueryExtractClass(t *testing.T) // function
func TestCppQueryExtractMethod(t *testing.T) // function
func TestCppQueryExtractConstructorDestructor(t *testing.T) // function
func TestCppQueryExtractNamespace(t *testing.T) // function
func TestCppQueryExtractTemplate(t *testing.T) // function
func TestCppQueryExtractStruct(t *testing.T) // function
func TestCppQueryExtractEnum(t *testing.T) // function
func TestCppQueryExtractMacro(t *testing.T) // function
func TestCppQueryExtractTypedef(t *testing.T) // function
func TestCppQueryExtractIncludes(t *testing.T) // function
imports []string // variable
func TestCppQueryNestedNamespaces(t *testing.T) // function
func TestCppQueryMultipleInheritance(t *testing.T) // function
func TestCppQueryEmptyFile(t *testing.T) // function
func TestCppQueryOnlyComments(t *testing.T) // function
nameCount int // variable
docCount int // variable
func TestCppQueryKindMapping(t *testing.T) // function
func TestCppQueryCaptures(t *testing.T) // function
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/go.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"`

```go
captureName      = "name" // variable
captureSignature = "signature" // variable
captureDoc       = "doc" // variable
captureKind      = "kind" // variable
type GoQuery struct {
	language *sitter.Language
	query    []byte
} // type
func NewGoQuery() *GoQuery // function
func (q *GoQuery) Language() *sitter.Language // method
func (q *GoQuery) Query() []byte // method
func (q *GoQuery) Captures() []string // method
func (q *GoQuery) KindMapping() map[string]string // method
func (q *GoQuery) ImportQuery() []byte // method
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
` // variable
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
` // variable
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/go_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"`

```go
func TestGoQueryLanguage(t *testing.T) // function
func TestGoQueryPattern(t *testing.T) // function
func TestGoQueryExtractFunction(t *testing.T) // function
funcCaptures map[string]string // variable
funcKindNode *sitter.Node // variable
kindNode *sitter.Node // variable
func TestGoQueryExtractConstAndVar(t *testing.T) // function
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
} // type
func NewJavaQuery() *JavaQuery // function
func (q *JavaQuery) Language() *sitter.Language // method
func (q *JavaQuery) Query() []byte // method
func (q *JavaQuery) Captures() []string // method
func (q *JavaQuery) KindMapping() map[string]string // method
func (q *JavaQuery) ImportQuery() []byte // method
javaImportQueryPattern = `
; import statements (capture full declaration)
(import_declaration) @import_path
` // variable
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
` // variable
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/java_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"`

```go
func TestJavaQueryLanguage(t *testing.T) // function
func TestJavaQueryPattern(t *testing.T) // function
func TestJavaQueryKindMapping(t *testing.T) // function
func TestJavaQueryExtractClass(t *testing.T) // function
foundClass, foundMethod bool // variable
func TestJavaQueryExtractInterface(t *testing.T) // function
func TestJavaQueryExtractEnum(t *testing.T) // function
foundEnum bool // variable
func TestJavaQueryExtractAnnotationType(t *testing.T) // function
foundAnnotation bool // variable
func TestJavaQueryExtractRecord(t *testing.T) // function
func TestJavaQueryExtractGenerics(t *testing.T) // function
func TestJavaQueryExtractFieldDeclarations(t *testing.T) // function
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
} // type
func NewPythonQuery() *PythonQuery // function
func (q *PythonQuery) Language() *sitter.Language // method
func (q *PythonQuery) Query() []byte // method
func (q *PythonQuery) Captures() []string // method
func (q *PythonQuery) KindMapping() map[string]string // method
func (q *PythonQuery) ImportQuery() []byte // method
pythonImportQueryPattern = `
; import module (capture full statement)
(import_statement) @import_path

; from module import ... (capture full statement)
(import_from_statement) @import_path
` // variable
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
` // variable
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/python_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"`

```go
func TestPythonQueryLanguage(t *testing.T) // function
func TestPythonQueryPattern(t *testing.T) // function
func TestPythonQueryExtractFunction(t *testing.T) // function
funcCaptures map[string]string // variable
func TestPythonQueryExtractClass(t *testing.T) // function
func TestPythonQueryExtractAsyncFunction(t *testing.T) // function
funcCaptures map[string]string // variable
func TestPythonQueryExtractModuleLevelVariables(t *testing.T) // function
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/rust.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"`

```go
type RustQuery struct {
	language *sitter.Language
	query    []byte
} // type
func NewRustQuery() *RustQuery // function
func (q *RustQuery) Language() *sitter.Language // method
func (q *RustQuery) Query() []byte // method
func (q *RustQuery) Captures() []string // method
func (q *RustQuery) KindMapping() map[string]string // method
func (q *RustQuery) ImportQuery() []byte // method
rustImportQueryPattern = `
; Use declarations (capture full statement)
(use_declaration) @import_path

; Extern crate declarations
(extern_crate_declaration) @import_path
` // variable
rustQueryPattern = `
; Functions (including async, unsafe, const, extern)
(function_item
  name: (identifier) @name
) @signature @kind

; Struct declarations
(struct_item
  name: (type_identifier) @name
) @signature @kind

; Enum declarations
(enum_item
  name: (type_identifier) @name
) @signature @kind

; Trait declarations
(trait_item
  name: (type_identifier) @name
) @signature @kind

; Type aliases
(type_item
  name: (type_identifier) @name
) @signature @kind

; Impl blocks (capture the whole impl signature)
(impl_item
  type: (type_identifier) @name
) @signature @kind

; Impl blocks for generic types
(impl_item
  type: (generic_type
    type: (type_identifier) @name
  )
) @signature @kind

; Trait impl blocks (impl Trait for Type)
(impl_item
  trait: (type_identifier)
  type: (type_identifier) @name
) @signature @kind

; Constants
(const_item
  name: (identifier) @name
) @signature @kind

; Statics
(static_item
  name: (identifier) @name
) @signature @kind

; Modules
(mod_item
  name: (identifier) @name
) @signature @kind

; Macro definitions (macro_rules!)
(macro_definition
  name: (identifier) @name
) @signature @kind

; Union declarations
(union_item
  name: (type_identifier) @name
) @signature @kind

; Foreign mod (extern "C" blocks)
(foreign_mod_item) @signature @kind

; Associated types in traits
(associated_type
  name: (type_identifier) @name
) @signature @kind

; Function signatures in traits (without body)
(function_signature_item
  name: (identifier) @name
) @signature @kind

; Doc comments (/// and //!)
(line_comment) @doc

; Block doc comments (/** and /*!)
(block_comment) @doc
` // variable
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/rust_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"`

```go
func TestRustQueryLanguage(t *testing.T) // function
func TestRustQueryPattern(t *testing.T) // function
func TestRustQueryImportPattern(t *testing.T) // function
func TestRustQueryExtractFunction(t *testing.T) // function
func TestRustQueryExtractTypes(t *testing.T) // function
func TestRustQueryExtractImplAndMethods(t *testing.T) // function
func TestRustQueryExtractConstAndStatic(t *testing.T) // function
func TestRustQueryExtractMacro(t *testing.T) // function
func TestRustQueryExtractModule(t *testing.T) // function
func TestRustQueryExtractUse(t *testing.T) // function
func TestRustQueryExtractGenericsAndLifetimes(t *testing.T) // function
func TestRustQueryEmptyFile(t *testing.T) // function
func TestRustQueryKindMapping(t *testing.T) // function
func TestRustQueryCaptures(t *testing.T) // function
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
} // type
func NewTypeScriptQuery() *TypeScriptQuery // function
func (q *TypeScriptQuery) Language() *sitter.Language // method
func (q *TypeScriptQuery) Query() []byte // method
func (q *TypeScriptQuery) Captures() []string // method
func (q *TypeScriptQuery) KindMapping() map[string]string // method
func (q *TypeScriptQuery) ImportQuery() []byte // method
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
` // variable
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
` // variable
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/typescript_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"`

```go
func TestTypeScriptQueryLanguage(t *testing.T) // function
func TestTypeScriptQueryPattern(t *testing.T) // function
func TestTypeScriptQueryExtractFunction(t *testing.T) // function
funcCaptures map[string]string // variable
func TestTypeScriptQueryExtractModuleLevelVariables(t *testing.T) // function
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
func init() // function
type TreeSitterParser struct {
	queries map[string]LanguageQuery
} // type
func NewTreeSitterParser() *TreeSitterParser // function
func (p *TreeSitterParser) Parse(content string, opts *parser.Options) (*parser.ParseResult, error) // method
imports []parser.ImportExport // variable
func (p *TreeSitterParser) Languages() []string // method
func (p *TreeSitterParser) extractSignatures(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
) []parser.Signature // method
signatures []parser.Signature // variable
kindNode *sitter.Node // variable
func cleanComment(text string) string // function
func isExported(name, language string) bool // function
func stripBody(text, kind, language string) string // function
func stripGoBody(text, kind string) string // function
tsFunctionBodyRe = regexp.MustCompile(`\s*\{[\s\S]*\}\s*$`) // variable
tsArrowBodyRe = regexp.MustCompile(`\s*=>\s*[\s\S]+$`) // variable
tsClassBodyRe = regexp.MustCompile(`\s*\{[\s\S]*\}\s*$`) // variable
func stripTypeScriptBody(text, kind string) string // function
func stripTSFunctionBody(text string) string // function
func findFunctionBodyStart(text string) int // function
func findTSClassBodyStart(text string) int // function
func stripPythonBody(text, kind string) string // function
func findPythonBodyStart(text string) int // function
func stripCBody(text, kind string) string // function
func stripCppBody(text, kind string) string // function
func findCppBodyStart(text string) int // function
func isPythonMethod(signature string) bool // function
func stripJavaBody(text, kind string) string // function
func findJavaBodyStart(text string) int // function
func stripRustBody(text, kind string) string // function
func findRustBodyStart(text string) int // function
func (p *TreeSitterParser) extractImports(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
) []parser.ImportExport // method
imports []parser.ImportExport // variable
imp parser.ImportExport // variable
hasExportType bool // variable
func cleanImportPath(path string) string // function
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/parser_test.go

**Imports:**
- `import "strings"`
- `import "testing"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`

```go
func TestTreeSitterParserImplementsParser(t *testing.T) // function
_ parser.Parser = (*TreeSitterParser)(nil) // variable
func TestTreeSitterParserLanguages(t *testing.T) // function
func TestTreeSitterParserParseGo(t *testing.T) // function
foundAdd bool // variable
func TestTreeSitterParserParseTypeScript(t *testing.T) // function
foundAdd bool // variable
func TestTreeSitterParserUnsupportedLanguage(t *testing.T) // function
func TestTreeSitterParserAutoRegistration(t *testing.T) // function
func TestGoSignatureOnlyExtraction(t *testing.T) // function
func TestGoIncludeBodyExtraction(t *testing.T) // function
foundAdd bool // variable
func TestTypeScriptSignatureOnlyExtraction(t *testing.T) // function
func TestTypeScriptArrowFunctionSignature(t *testing.T) // function
func contains(s, substr string) bool // function
func TestTreeSitterParserParseJava(t *testing.T) // function
foundClass, foundConstructor, foundPublicMethod, foundPrivateMethod bool // variable
func TestJavaSignatureOnlyExtraction(t *testing.T) // function
func TestJavaGenericsExtraction(t *testing.T) // function
foundClass, foundMethod bool // variable
func TestJavaAutoRegistration(t *testing.T) // function
func TestTreeSitterParserParseCpp(t *testing.T) // function
func TestCppSignatureOnlyExtraction(t *testing.T) // function
func TestCppTemplateExtraction(t *testing.T) // function
func TestCppAutoRegistration(t *testing.T) // function
func TestCppImportExtraction(t *testing.T) // function
func TestGoVariableExtraction(t *testing.T) // function
func TestTypeScriptVariableExtraction(t *testing.T) // function
func TestPythonVariableExtraction(t *testing.T) // function
func TestJavaStaticFieldExtraction(t *testing.T) // function
func TestCGlobalVariableExtraction(t *testing.T) // function
func TestVariableSignaturePreservesValue(t *testing.T) // function
func TestTreeSitterParserParseRust(t *testing.T) // function
func TestRustSignatureOnlyExtraction(t *testing.T) // function
func TestRustImportExtraction(t *testing.T) // function
func TestRustAutoRegistration(t *testing.T) // function
func TestRustConstAndStaticExtraction(t *testing.T) // function
func TestRustMacroExtraction(t *testing.T) // function
func TestRustGenericsAndLifetimes(t *testing.T) // function
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
} // type
CaptureName      = "name" // variable
CaptureSignature = "signature" // variable
CaptureDoc       = "doc" // variable
CaptureKind      = "kind" // variable
CaptureImportPath = "import_path" // variable
CaptureExportName = "export_name" // variable
CaptureImportType = "import_type" // variable
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
} // variable
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/query_test.go

**Imports:**
- `import "testing"`

```go
func TestCaptureDefinitions(t *testing.T) // function
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
} // type
type ScanResult struct {
	// Files is the list of matched files.
	Files []FileEntry

	// TotalSize is the sum of all matched file sizes.
	TotalSize int64

	// SkippedCount is the number of files skipped (too large, unsupported, etc.).
	SkippedCount int
} // type
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
} // type
func DefaultScanOptions() *ScanOptions // function
func (o *ScanOptions) GetLanguage(path string) (string, bool) // method
func IsHidden(name string) bool // function
type Scanner interface {
	// Scan performs the scan and returns scan results.
	Scan() (*ScanResult, error)
} // type
type FileScanner struct {
	opts       *ScanOptions
	ignorer    *ignore.GitIgnore
	ignorerErr error
	logger     *log.Logger
} // type
func NewFileScanner(opts *ScanOptions) (*FileScanner, error) // function
func (s *FileScanner) Scan() (*ScanResult, error) // method
func (s *FileScanner) checkFile(path string, info os.FileInfo) (FileEntry, bool) // method
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/scanner/scanner_test.go

**Imports:**
- `import "os"`
- `import "path/filepath"`
- `import "testing"`

```go
func TestNewFileScanner(t *testing.T) // function
func TestNewFileScannerNilOptions(t *testing.T) // function
func TestFileEntryDefaults(t *testing.T) // function
func TestScanOptionsDefaults(t *testing.T) // function
expectedMaxSize = 512000 // variable
func TestScanOptionsWithExtensions(t *testing.T) // function
func TestScannerInterface(t *testing.T) // function
_ Scanner = (*FileScanner)(nil) // variable
func TestScanEmptyDirectory(t *testing.T) // function
func TestScanSingleFile(t *testing.T) // function
func TestScanFilterByExtension(t *testing.T) // function
func TestScanExcludeHidden(t *testing.T) // function
func TestScanIncludeHidden(t *testing.T) // function
func TestScanMaxFileSize(t *testing.T) // function
func TestScanGitignore(t *testing.T) // function
func TestScanNestedDirectories(t *testing.T) // function
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/tokenizer/tiktoken.go

**Imports:**
- `import "github.com/pkoukk/tiktoken-go"`

```go
type TiktokenTokenizer struct {
	encoding string
	tke      *tiktoken.Tiktoken
} // type
_ Tokenizer = (*TiktokenTokenizer)(nil) // variable
func NewTiktokenTokenizer() (*TiktokenTokenizer, error) // function
func (t *TiktokenTokenizer) Count(text string) (int, error) // method
func (t *TiktokenTokenizer) Name() string // method
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
} // type
type NoOpTokenizer struct{} // type
_ Tokenizer = (*NoOpTokenizer)(nil) // variable
func NewNoOpTokenizer() *NoOpTokenizer // function
func (t *NoOpTokenizer) Count(_ string) (int, error) // method
func (t *NoOpTokenizer) Name() string // method
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/tokenizer/tokenizer_test.go

**Imports:**
- `import "strings"`
- `import "testing"`

```go
func TestNoOpTokenizerImplementsTokenizer(t *testing.T) // function
_ Tokenizer = (*NoOpTokenizer)(nil) // variable
func TestTiktokenTokenizerImplementsTokenizer(t *testing.T) // function
_ Tokenizer = (*TiktokenTokenizer)(nil) // variable
func TestNoOpTokenizerCount(t *testing.T) // function
func TestNoOpTokenizerName(t *testing.T) // function
func TestTiktokenTokenizerCount(t *testing.T) // function
func TestTiktokenTokenizerName(t *testing.T) // function
func TestTiktokenTokenizerConsistency(t *testing.T) // function
func TestTiktokenTokenizerSpecialCharacters(t *testing.T) // function
```

---

