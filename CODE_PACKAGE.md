# Brf.it Output

---

## Files

### cmd/brfit/main.go

```go
```

---

### cmd/brfit/root.go

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
func Execute()
func NewRootCommand() *cobra.Command
```

---

### cmd/brfit/root_test.go

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
func TestNewRootCommand(t *testing.T)
func TestParseFlags(t *testing.T)
func TestRootCommandIntegration(t *testing.T)
func TestRootCommandIntegrationMarkdown(t *testing.T)
func TestRootCommandPathNotFound(t *testing.T)
func TestWriteToFile(t *testing.T)
```

---

### internal/config/config.go

**Imports:**
- `import "errors"`
- `import "fmt"`
- `import pkgcontext "github.com/indigo-net/Brf.it/internal/context"`

```go
type Config struct {
	// Path is the root directory or file to process.
	Path string

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

### internal/config/config_test.go

**Imports:**
- `import "testing"`

```go
func TestDefaultConfig(t *testing.T)
func TestConfigValidate(t *testing.T)
func TestConfigSupportedLanguages(t *testing.T)
```

---

### internal/context/context.go

**Imports:**
- `import "github.com/indigo-net/Brf.it/pkg/extractor"`
- `import "github.com/indigo-net/Brf.it/pkg/formatter"`
- `import "github.com/indigo-net/Brf.it/pkg/scanner"`
- `import "github.com/indigo-net/Brf.it/pkg/tokenizer"`

```go
type Options struct {
	// Path is the target path to scan.
	Path string

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
func NewDefaultPackager(scanOpts *scanner.ScanOptions) (*Packager, error)
```

---

### internal/context/context_test.go

**Imports:**
- `import "strings"`
- `import "testing"`
- `import "github.com/indigo-net/Brf.it/pkg/extractor"`
- `import "github.com/indigo-net/Brf.it/pkg/formatter"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`
- `import "github.com/indigo-net/Brf.it/pkg/scanner"`
- `import "github.com/indigo-net/Brf.it/pkg/tokenizer"`

```go
func (m *mockScanner) Scan() (*scanner.ScanResult, error)
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

### internal/context/tree.go

**Imports:**
- `import "path/filepath"`
- `import "sort"`
- `import "strings"`

```go
func BuildTree(root string, paths []string) string
```

---

### pkg/extractor/extractor.go

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
```

---

### pkg/extractor/extractor_test.go

**Imports:**
- `import "os"`
- `import "path/filepath"`
- `import "testing"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`
- `import _ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"`
- `import "github.com/indigo-net/Brf.it/pkg/scanner"`

```go
func TestFileExtractorImplementsExtractor(t *testing.T)
func TestFileExtractorExtract(t *testing.T)
func TestFileExtractorUnsupportedLanguage(t *testing.T)
```

---

### pkg/formatter/formatter.go

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

### pkg/formatter/formatter_test.go

**Imports:**
- `import "fmt"`
- `import "strings"`
- `import "testing"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`

```go
func TestXMLFormatterImplementsFormatter(t *testing.T)
func TestMarkdownFormatterImplementsFormatter(t *testing.T)
func TestXMLFormatterFormat(t *testing.T)
func TestXMLFormatterFormatWithError(t *testing.T)
func TestMarkdownFormatterFormat(t *testing.T)
func TestFormatterNames(t *testing.T)
func TestXMLFormatterEscapeXML(t *testing.T)
func TestMarkdownFormatterEscapeMarkdown(t *testing.T)
func TestXMLFormatterEmptyData(t *testing.T)
func TestMarkdownFormatterEmptyData(t *testing.T)
```

---

### pkg/formatter/markdown.go

**Imports:**
- `import "bytes"`
- `import "fmt"`
- `import "strings"`

```go
type MarkdownFormatter struct{}
func NewMarkdownFormatter() *MarkdownFormatter
func (f *MarkdownFormatter) Name() string
func (f *MarkdownFormatter) Format(data *PackageData) ([]byte, error)
```

---

### pkg/formatter/xml.go

**Imports:**
- `import "bytes"`
- `import "fmt"`
- `import "strings"`

```go
type XMLFormatter struct{}
func NewXMLFormatter() *XMLFormatter
func (f *XMLFormatter) Name() string
func (f *XMLFormatter) Format(data *PackageData) ([]byte, error)
```

---

### pkg/parser/parser.go

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

### pkg/parser/parser_test.go

**Imports:**
- `import "testing"`

```go
func TestSignatureDefaults(t *testing.T)
func TestParseResultDefaults(t *testing.T)
func TestNodeKind(t *testing.T)
func TestParserInterface(t *testing.T)
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

### pkg/parser/treesitter/languages/c.go

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
```

---

### pkg/parser/treesitter/languages/c_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"`

```go
func TestCQueryLanguage(t *testing.T)
func TestCQueryPattern(t *testing.T)
func TestCQueryExtractFunction(t *testing.T)
func TestCQueryExtractStruct(t *testing.T)
func TestCQueryExtractMacro(t *testing.T)
func TestCQueryExtractEnum(t *testing.T)
func TestCQueryExtractTypedef(t *testing.T)
func TestCQueryExtractGlobalVariables(t *testing.T)
```

---

### pkg/parser/treesitter/languages/go.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"`

```go
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
```

---

### pkg/parser/treesitter/languages/go_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"`

```go
func TestGoQueryLanguage(t *testing.T)
func TestGoQueryPattern(t *testing.T)
func TestGoQueryExtractFunction(t *testing.T)
func TestGoQueryExtractConstAndVar(t *testing.T)
```

---

### pkg/parser/treesitter/languages/java.go

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
```

---

### pkg/parser/treesitter/languages/java_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"`

```go
func TestJavaQueryLanguage(t *testing.T)
func TestJavaQueryPattern(t *testing.T)
func TestJavaQueryKindMapping(t *testing.T)
func TestJavaQueryExtractClass(t *testing.T)
func TestJavaQueryExtractInterface(t *testing.T)
func TestJavaQueryExtractEnum(t *testing.T)
func TestJavaQueryExtractAnnotationType(t *testing.T)
func TestJavaQueryExtractRecord(t *testing.T)
func TestJavaQueryExtractGenerics(t *testing.T)
func TestJavaQueryExtractFieldDeclarations(t *testing.T)
```

---

### pkg/parser/treesitter/languages/python.go

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
```

---

### pkg/parser/treesitter/languages/python_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"`

```go
func TestPythonQueryLanguage(t *testing.T)
func TestPythonQueryPattern(t *testing.T)
func TestPythonQueryExtractFunction(t *testing.T)
func TestPythonQueryExtractClass(t *testing.T)
func TestPythonQueryExtractAsyncFunction(t *testing.T)
func TestPythonQueryExtractModuleLevelVariables(t *testing.T)
```

---

### pkg/parser/treesitter/languages/typescript.go

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
```

---

### pkg/parser/treesitter/languages/typescript_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"`

```go
func TestTypeScriptQueryLanguage(t *testing.T)
func TestTypeScriptQueryPattern(t *testing.T)
func TestTypeScriptQueryExtractFunction(t *testing.T)
func TestTypeScriptQueryExtractModuleLevelVariables(t *testing.T)
```

---

### pkg/parser/treesitter/parser.go

**Imports:**
- `import "fmt"`
- `import "regexp"`
- `import "strings"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`
- `import "github.com/indigo-net/Brf.it/pkg/parser/treesitter/languages"`

```go
type TreeSitterParser struct {
	queries map[string]LanguageQuery
}
func NewTreeSitterParser() *TreeSitterParser
func (p *TreeSitterParser) Parse(content string, opts *parser.Options) (*parser.ParseResult, error)
func (p *TreeSitterParser) Languages() []string
```

---

### pkg/parser/treesitter/parser_test.go

**Imports:**
- `import "strings"`
- `import "testing"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`

```go
func TestTreeSitterParserImplementsParser(t *testing.T)
func TestTreeSitterParserLanguages(t *testing.T)
func TestTreeSitterParserParseGo(t *testing.T)
func TestTreeSitterParserParseTypeScript(t *testing.T)
func TestTreeSitterParserUnsupportedLanguage(t *testing.T)
func TestTreeSitterParserAutoRegistration(t *testing.T)
func TestGoSignatureOnlyExtraction(t *testing.T)
func TestGoIncludeBodyExtraction(t *testing.T)
func TestTypeScriptSignatureOnlyExtraction(t *testing.T)
func TestTypeScriptArrowFunctionSignature(t *testing.T)
func TestTreeSitterParserParseJava(t *testing.T)
func TestJavaSignatureOnlyExtraction(t *testing.T)
func TestJavaGenericsExtraction(t *testing.T)
func TestJavaAutoRegistration(t *testing.T)
func TestGoVariableExtraction(t *testing.T)
func TestTypeScriptVariableExtraction(t *testing.T)
func TestPythonVariableExtraction(t *testing.T)
func TestJavaStaticFieldExtraction(t *testing.T)
func TestCGlobalVariableExtraction(t *testing.T)
func TestVariableSignaturePreservesValue(t *testing.T)
```

---

### pkg/parser/treesitter/query.go

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

### pkg/parser/treesitter/query_test.go

**Imports:**
- `import "testing"`

```go
func TestCaptureDefinitions(t *testing.T)
```

---

### pkg/scanner/scanner.go

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
```

---

### pkg/scanner/scanner_test.go

**Imports:**
- `import "os"`
- `import "path/filepath"`
- `import "testing"`

```go
func TestNewFileScanner(t *testing.T)
func TestNewFileScannerNilOptions(t *testing.T)
func TestFileEntryDefaults(t *testing.T)
func TestScanOptionsDefaults(t *testing.T)
func TestScanOptionsWithExtensions(t *testing.T)
func TestScannerInterface(t *testing.T)
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

### pkg/tokenizer/tiktoken.go

**Imports:**
- `import "github.com/pkoukk/tiktoken-go"`

```go
type TiktokenTokenizer struct {
	encoding string
	tke      *tiktoken.Tiktoken
}
func NewTiktokenTokenizer() (*TiktokenTokenizer, error)
func (t *TiktokenTokenizer) Count(text string) (int, error)
func (t *TiktokenTokenizer) Name() string
```

---

### pkg/tokenizer/tokenizer.go

```go
type Tokenizer interface {
	// Count returns the number of tokens in the given text.
	// Returns 0 and error if counting fails.
	Count(text string) (int, error)

	// Name returns the tokenizer name (e.g., "tiktoken-cl100k", "noop").
	Name() string
}
type NoOpTokenizer struct{}
func NewNoOpTokenizer() *NoOpTokenizer
func (t *NoOpTokenizer) Count(_ string) (int, error)
func (t *NoOpTokenizer) Name() string
```

---

### pkg/tokenizer/tokenizer_test.go

**Imports:**
- `import "strings"`
- `import "testing"`

```go
func TestNoOpTokenizerImplementsTokenizer(t *testing.T)
func TestTiktokenTokenizerImplementsTokenizer(t *testing.T)
func TestNoOpTokenizerCount(t *testing.T)
func TestNoOpTokenizerName(t *testing.T)
func TestTiktokenTokenizerCount(t *testing.T)
func TestTiktokenTokenizerName(t *testing.T)
func TestTiktokenTokenizerConsistency(t *testing.T)
func TestTiktokenTokenizerSpecialCharacters(t *testing.T)
```

---

