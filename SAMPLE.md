# Code Summary: /home/runner/work/Brf.it/Brf.it

*brf.it 0.18.0*

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
- `import "os"`
- `import pkgcontext "github.com/indigo-net/Brf.it/internal/context"`

```go
MaxFileSizeUpperBound = 10 * 1024 * 1024
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
- `import "bytes"`
- `import "os"`
- `import "strings"`
- `import "testing"`

```go
func TestDefaultConfig(t *testing.T)
expectedMaxSize = 512000
func TestConfigValidate(t *testing.T)
func TestConfigSupportedLanguages(t *testing.T)
func TestValidateMaxFileSizeUpperBound(t *testing.T)
buf bytes.Buffer
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
func TestPackagerNoStdImportsPassthrough(t *testing.T)
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
- `import "bytes"`
- `import "fmt"`
- `import "os"`
- `import "runtime"`
- `import "sync"`
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

	// RawImports is the list of raw import/export statement text.
	RawImports []string

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

	// Concurrency is the number of concurrent workers.
	// 0 = auto (runtime.NumCPU()), 1 = sequential.
	Concurrency int

	// MaxFileSize is the maximum file size in bytes for TOCTOU re-check.
	// If positive, file content size is verified after reading.
	MaxFileSize int64
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
func DefaultExtractOptions() *ExtractOptions
func (e *FileExtractor) Extract(scanResult *scanner.ScanResult, opts *ExtractOptions) (*ExtractResult, error)
wg sync.WaitGroup
func (e *FileExtractor) extractSequential(files []scanner.FileEntry, opts *ExtractOptions) *ExtractResult
binarySniffSize = 512
func isBinaryContent(content []byte) bool
func (e *FileExtractor) extractFile(entry scanner.FileEntry, opts *ExtractOptions) ExtractedFile
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/extractor/extractor_test.go

**Imports:**
- `import "fmt"`
- `import "os"`
- `import "path/filepath"`
- `import "strings"`
- `import "testing"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`
- `import _ "github.com/indigo-net/Brf.it/pkg/parser/treesitter"`
- `import "github.com/indigo-net/Brf.it/pkg/scanner"`

```go
func TestFileExtractorImplementsExtractor(t *testing.T)
_ Extractor = (*FileExtractor)(nil)
func TestFileExtractorExtract(t *testing.T)
foundAdd bool
func TestFileExtractorTOCTOUGuard(t *testing.T)
func TestFileExtractorTOCTOUGuardDisabled(t *testing.T)
func TestExtractConcurrencySequential(t *testing.T)
func TestExtractConcurrencyDeterministicOrder(t *testing.T)
entries []scanner.FileEntry
func TestExtractConcurrencyEmptyFiles(t *testing.T)
func TestExtractNilScanResult(t *testing.T)
func TestExtractNegativeConcurrency(t *testing.T)
func TestDefaultExtractOptions(t *testing.T)
func TestExtractConcurrencyWithErrors(t *testing.T)
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

	// RawImports is the list of raw import/export statement text.
	RawImports []string

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
func TestJSONFormatterImplementsFormatter(t *testing.T)
_ Formatter = (*JSONFormatter)(nil)
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
func TestKindToTag(t *testing.T)
func TestXMLFormatterKindTags(t *testing.T)
func TestXMLFormatterVerbatimImports(t *testing.T)
func TestMarkdownFormatterVerbatimImports(t *testing.T)
func TestJSONFormatterFormat(t *testing.T)
func TestJSONFormatterKindNormalization(t *testing.T)
func TestJSONFormatterWithImports(t *testing.T)
func TestJSONFormatterWithError(t *testing.T)
func TestJSONFormatterName(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/helpers.go

```go
func normalizeKind(kind string) string
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

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/json.go

**Imports:**
- `import "encoding/json"`

```go
type JSONFormatter struct{}
func NewJSONFormatter() *JSONFormatter
func (f *JSONFormatter) Name() string
type jsonOutput struct {
	Version string     `json:"version,omitempty"`
	Path    string     `json:"path,omitempty"`
	Tree    string     `json:"tree,omitempty"`
	Files   []jsonFile `json:"files"`
}
type jsonFile struct {
	Path       string        `json:"path"`
	Language   string        `json:"language"`
	Signatures []jsonSig     `json:"signatures,omitempty"`
	Imports    []string      `json:"imports,omitempty"`
	Error      string        `json:"error,omitempty"`
}
type jsonSig struct {
	Kind string `json:"kind"`
	Text string `json:"text"`
	Doc  string `json:"doc,omitempty"`
}
func (f *JSONFormatter) Format(data *PackageData) ([]byte, error)
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
func (f *XMLFormatter) Format(data *PackageData) ([]byte, error)
buf bytes.Buffer
func escapeXML(s string) string
needsEscape bool
buf strings.Builder
func kindToTag(kind string) string
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/parser.go

**Imports:**
- `import "path/filepath"`
- `import "sync"`

```go
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

	// RawImports is the list of raw import/export statement text.
	RawImports []string

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
	// Content is passed as []byte to avoid unnecessary string conversion
	// from os.ReadFile output.
	Parse(content []byte, opts *Options) (*ParseResult, error)

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
	".kts":   "kotlin",
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
func (m *MockParser) Parse(content []byte, opts *Options) (*ParseResult, error)
func (m *MockParser) Languages() []string
func TestMockParser(t *testing.T)
func TestRegistry(t *testing.T)
func TestDefaultRegistry(t *testing.T)
func TestDetectLanguage(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/csharp/binding.go

**Imports:**
- `import "C"`
- `import "unsafe"`

```go
func Language() unsafe.Pointer
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/csharp/scanner.c

**Imports:**
- `#include "tree_sitter/alloc.h"`
- `#include "tree_sitter/array.h"`
- `#include "tree_sitter/parser.h"`
- `#include <wctype.h>`

```c
enum TokenType {
    OPT_SEMI,
    INTERPOLATION_REGULAR_START,
    INTERPOLATION_VERBATIM_START,
    INTERPOLATION_RAW_START,
    INTERPOLATION_START_QUOTE,
    INTERPOLATION_END_QUOTE,
    INTERPOLATION_OPEN_BRACE,
    INTERPOLATION_CLOSE_BRACE,
    INTERPOLATION_STRING_CONTENT,
    RAW_STRING_START,
    RAW_STRING_END,
    RAW_STRING_CONTENT,
}
typedef enum {
    REGULAR = 1 << 0,
    VERBATIM = 1 << 1,
    RAW = 1 << 2,
} StringType;
typedef struct {
    uint8_t dollar_count;
    uint8_t open_brace_count;
    uint8_t quote_count;
    StringType string_type;
} Interpolation;
static inline bool is_regular(Interpolation *interpolation)
static inline bool is_verbatim(Interpolation *interpolation)
static inline bool is_raw(Interpolation *interpolation)
typedef struct {
    uint8_t quote_count;
    Array(Interpolation) interpolation_stack;
} Scanner;
static inline void advance(TSLexer *lexer)
static inline void skip(TSLexer *lexer)
void *tree_sitter_c_sharp_external_scanner_create()
void tree_sitter_c_sharp_external_scanner_destroy(void *payload)
unsigned tree_sitter_c_sharp_external_scanner_serialize(void *payload, char *buffer)
void tree_sitter_c_sharp_external_scanner_deserialize(void *payload, const char *buffer, unsigned length)
bool tree_sitter_c_sharp_external_scanner_scan(void *payload, TSLexer *lexer, const bool *valid_symbols)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/csharp/tree_sitter/alloc.h

**Imports:**
- `#include <stdbool.h>`
- `#include <stdio.h>`
- `#include <stdlib.h>`

```cpp
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/csharp/tree_sitter/array.h

**Imports:**
- `#include "./alloc.h"`
- `#include <assert.h>`
- `#include <stdbool.h>`
- `#include <stdint.h>`
- `#include <stdlib.h>`
- `#include <string.h>`

```cpp
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity)        \
  ((self)->contents = _array__reserve(           \
    (void *)(self)->contents, &(self)->capacity, \
    array_elem_size(self), new_capacity)         \
  )
#define array_delete(self) _array__delete((self), (void *)(self)->contents, sizeof(*self))
#define array_push(self, element)                                 \
  do {                                                            \
    (self)->contents = _array__grow(                              \
      (void *)(self)->contents, (self)->size, &(self)->capacity,  \
      1, array_elem_size(self)                                    \
    );                                                            \
   (self)->contents[(self)->size++] = (element);                  \
  } while(0)
#define array_grow_by(self, count)                                               \
  do {                                                                           \
    if ((count) == 0) break;                                                     \
    (self)->contents = _array__grow(                                             \
      (self)->contents, (self)->size, &(self)->capacity,                         \
      count, array_elem_size(self)                                               \
    );                                                                           \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count);                                                     \
  } while (0)
#define array_push_all(self, other) \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, other_contents)                 \
  (self)->contents = _array__splice(                              \
    (void*)(self)->contents, &(self)->size, &(self)->capacity,    \
    array_elem_size(self), (self)->size, 0, count, other_contents \
  )
#define array_splice(self, _index, old_count, new_count, new_contents) \
  (self)->contents = _array__splice(                                   \
    (void *)(self)->contents, &(self)->size, &(self)->capacity,        \
    array_elem_size(self), _index, old_count, new_count, new_contents  \
  )
#define array_insert(self, _index, element)                     \
  (self)->contents = _array__splice(                            \
    (void *)(self)->contents, &(self)->size, &(self)->capacity, \
    array_elem_size(self), _index, 0, 1, &(element)             \
  )
#define array_erase(self, _index) \
  _array__erase((void *)(self)->contents, &(self)->size, array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other)                                   \
  (self)->contents = _array__assign(                                \
    (void *)(self)->contents, &(self)->size, &(self)->capacity,     \
    (const void *)(other)->contents, (other)->size, array_elem_size(self) \
  )
#define array_swap(self, other)                                     \
  do {                                                              \
    struct Swap swapped_contents = _array__swap(                    \
      (void *)(self)->contents, &(self)->size, &(self)->capacity,   \
      (void *)(other)->contents, &(other)->size, &(other)->capacity \
    );                                                              \
    (self)->contents = swapped_contents.self_contents;              \
    (other)->contents = swapped_contents.other_contents;            \
  } while (0)
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__delete(void *self, void *contents, size_t self_size)
static inline void _array__erase(void* self_contents, uint32_t *size,
                                size_t element_size, uint32_t index)
static inline void *_array__reserve(void *contents, uint32_t *capacity,
                                  size_t element_size, uint32_t new_capacity)
static inline void *_array__assign(void* self_contents, uint32_t *self_size, uint32_t *self_capacity,
                                 const void *other_contents, uint32_t other_size, size_t element_size)
struct Swap
struct Swap
struct Swap
static inline void *_array__grow(void *contents, uint32_t size, uint32_t *capacity,
                               uint32_t count, size_t element_size)
static inline void *_array__splice(void *self_contents, uint32_t *size, uint32_t *capacity,
                                 size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/csharp/tree_sitter/parser.h

**Imports:**
- `#include <stdbool.h>`
- `#include <stdint.h>`
- `#include <stdlib.h>`

```cpp
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
struct TSLanguageMetadata
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
  uint16_t reserved_word_set_id;
} TSLexerMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(const TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/kotlin/binding.go

**Imports:**
- `import "C"`
- `import "unsafe"`

```go
func Language() unsafe.Pointer
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/kotlin/scanner.c

**Imports:**
- `#include "tree_sitter/array.h"`
- `#include "tree_sitter/parser.h"`
- `#include <string.h>`
- `#include <wctype.h>`

```c
enum TokenType {
  AUTOMATIC_SEMICOLON,
  IMPORT_LIST_DELIMITER,
  SAFE_NAV,
  MULTILINE_COMMENT,
  STRING_START,
  STRING_END,
  STRING_CONTENT,
  PRIMARY_CONSTRUCTOR_KEYWORD,
  IMPORT_DOT,
}
#define DELIMITER_LENGTH 3
typedef char Delimiter;
typedef Array(Delimiter) Stack;
static inline void stack_push(Stack *stack, char chr, bool triple)
static inline Delimiter stack_pop(Stack *stack)
static inline void skip(TSLexer *lexer)
static inline void advance(TSLexer *lexer)
static bool scan_string_start(TSLexer *lexer, Stack *stack)
static bool scan_string_content(TSLexer *lexer, Stack *stack)
static bool scan_multiline_comment(TSLexer *lexer)
static bool scan_whitespace_and_comments(TSLexer *lexer)
static bool is_word_char(int32_t c)
static bool scan_for_word(TSLexer *lexer, const char* word, unsigned len)
static bool check_word(TSLexer *lexer, const char *word, unsigned len)
static bool check_modifier_then_constructor(TSLexer *lexer)
static bool scan_automatic_semicolon(TSLexer *lexer, const bool *valid_symbols)
static bool scan_safe_nav(TSLexer *lexer)
static bool scan_line_sep(TSLexer *lexer)
static bool scan_import_list_delimiter(TSLexer *lexer)
static bool scan_import_dot(TSLexer *lexer)
bool tree_sitter_kotlin_external_scanner_scan(void *payload, TSLexer *lexer, const bool *valid_symbols)
void *tree_sitter_kotlin_external_scanner_create()
void tree_sitter_kotlin_external_scanner_destroy(void *payload)
unsigned tree_sitter_kotlin_external_scanner_serialize(void *payload, char *buffer)
void tree_sitter_kotlin_external_scanner_deserialize(void *payload, const char *buffer, unsigned length)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/kotlin/tree_sitter/alloc.h

**Imports:**
- `#include <stdbool.h>`
- `#include <stdio.h>`
- `#include <stdlib.h>`

```cpp
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/kotlin/tree_sitter/array.h

**Imports:**
- `#include "./alloc.h"`
- `#include <assert.h>`
- `#include <stdbool.h>`
- `#include <stdint.h>`
- `#include <stdlib.h>`
- `#include <string.h>`

```cpp
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity) \
  _array__reserve((Array *)(self), array_elem_size(self), new_capacity)
#define array_delete(self) _array__delete((Array *)(self))
#define array_push(self, element)                            \
  (_array__grow((Array *)(self), 1, array_elem_size(self)), \
   (self)->contents[(self)->size++] = (element))
#define array_grow_by(self, count) \
  do { \
    if ((count) == 0) break; \
    _array__grow((Array *)(self), count, array_elem_size(self)); \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count); \
  } while (0)
#define array_push_all(self, other)                                       \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, contents)                    \
  _array__splice(                                               \
    (Array *)(self), array_elem_size(self), (self)->size, \
    0, count,  contents                                        \
  )
#define array_splice(self, _index, old_count, new_count, new_contents)  \
  _array__splice(                                                       \
    (Array *)(self), array_elem_size(self), _index,                \
    old_count, new_count, new_contents                                 \
  )
#define array_insert(self, _index, element) \
  _array__splice((Array *)(self), array_elem_size(self), _index, 0, 1, &(element))
#define array_erase(self, _index) \
  _array__erase((Array *)(self), array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other) \
  _array__assign((Array *)(self), (const Array *)(other), array_elem_size(self))
#define array_swap(self, other) \
  _array__swap((Array *)(self), (Array *)(other))
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__delete(Array *self)
static inline void _array__erase(Array *self, size_t element_size,
                                uint32_t index)
static inline void _array__reserve(Array *self, size_t element_size, uint32_t new_capacity)
static inline void _array__assign(Array *self, const Array *other, size_t element_size)
static inline void _array__swap(Array *self, Array *other)
static inline void _array__grow(Array *self, uint32_t count, size_t element_size)
static inline void _array__splice(Array *self, size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/kotlin/tree_sitter/parser.h

**Imports:**
- `#include <stdbool.h>`
- `#include <stdint.h>`
- `#include <stdlib.h>`

```cpp
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSFieldMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/binding.go

**Imports:**
- `import "C"`
- `import "unsafe"`

```go
func Language() unsafe.Pointer
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/parser.c

**Imports:**
- `#include "tree_sitter/parser.h"`

```c
#define LANGUAGE_VERSION 15
#define STATE_COUNT 282
#define LARGE_STATE_COUNT 2
#define SYMBOL_COUNT 143
#define ALIAS_COUNT 0
#define TOKEN_COUNT 73
#define EXTERNAL_TOKEN_COUNT 6
#define FIELD_COUNT 23
#define MAX_ALIAS_SEQUENCE_LENGTH 7
#define MAX_RESERVED_WORD_SET_SIZE 0
#define PRODUCTION_ID_COUNT 70
#define SUPERTYPE_COUNT 4
enum ts_symbol_identifiers {
  sym_identifier = 1,
  sym_hash_bang_line = 2,
  anon_sym_return = 3,
  anon_sym_SEMI = 4,
  anon_sym_EQ = 5,
  anon_sym_COMMA = 6,
  anon_sym_COLON_COLON = 7,
  sym_break_statement = 8,
  anon_sym_goto = 9,
  anon_sym_do = 10,
  anon_sym_end = 11,
  anon_sym_while = 12,
  anon_sym_repeat = 13,
  anon_sym_until = 14,
  anon_sym_if = 15,
  anon_sym_then = 16,
  anon_sym_elseif = 17,
  anon_sym_else = 18,
  anon_sym_for = 19,
  anon_sym_in = 20,
  anon_sym_function = 21,
  anon_sym_local = 22,
  anon_sym_global = 23,
  anon_sym_DOT = 24,
  anon_sym_COLON = 25,
  anon_sym_STAR = 26,
  anon_sym_LT = 27,
  anon_sym_GT = 28,
  sym_nil = 29,
  sym_false = 30,
  sym_true = 31,
  sym_number = 32,
  anon_sym_DQUOTE = 33,
  anon_sym_SQUOTE = 34,
  aux_sym__doublequote_string_content_token1 = 35,
  aux_sym__singlequote_string_content_token1 = 36,
  sym_escape_sequence = 37,
  sym_vararg_expression = 38,
  anon_sym_LPAREN = 39,
  anon_sym_RPAREN = 40,
  anon_sym_LBRACK = 41,
  anon_sym_RBRACK = 42,
  anon_sym_LBRACE = 43,
  anon_sym_RBRACE = 44,
  anon_sym_or = 45,
  anon_sym_and = 46,
  anon_sym_LT_EQ = 47,
  anon_sym_EQ_EQ = 48,
  anon_sym_TILDE_EQ = 49,
  anon_sym_GT_EQ = 50,
  anon_sym_PIPE = 51,
  anon_sym_TILDE = 52,
  anon_sym_AMP = 53,
  anon_sym_LT_LT = 54,
  anon_sym_GT_GT = 55,
  anon_sym_PLUS = 56,
  anon_sym_DASH = 57,
  anon_sym_SLASH = 58,
  anon_sym_SLASH_SLASH = 59,
  anon_sym_PERCENT = 60,
  anon_sym_DOT_DOT = 61,
  anon_sym_CARET = 62,
  anon_sym_not = 63,
  anon_sym_POUND = 64,
  anon_sym_DASH_DASH = 65,
  aux_sym_comment_token1 = 66,
  sym__block_comment_start = 67,
  sym__block_comment_content = 68,
  sym__block_comment_end = 69,
  sym__block_string_start = 70,
  sym__block_string_content = 71,
  sym__block_string_end = 72,
  sym_chunk = 73,
  sym__block = 74,
  sym_statement = 75,
  sym_return_statement = 76,
  sym_empty_statement = 77,
  sym_assignment_statement = 78,
  sym__variable_assignment_varlist = 79,
  sym__variable_assignment_explist = 80,
  sym_label_statement = 81,
  sym_goto_statement = 82,
  sym_do_statement = 83,
  sym_while_statement = 84,
  sym_repeat_statement = 85,
  sym_if_statement = 86,
  sym_elseif_statement = 87,
  sym_else_statement = 88,
  sym_for_statement = 89,
  sym_for_generic_clause = 90,
  sym_for_numeric_clause = 91,
  sym__name_list = 92,
  sym_declaration = 93,
  sym_function_declaration = 94,
  sym__local_function_declaration = 95,
  sym__global_function_declaration = 96,
  sym__function_name = 97,
  sym__function_name_prefix_expression = 98,
  sym__function_name_dot_index_expression = 99,
  sym__function_name_method_index_expression = 100,
  sym_variable_declaration = 101,
  sym__global_variable_declaration = 102,
  sym__variable_assignment = 103,
  sym__att_name_list = 104,
  sym__global_implicit_variable_declaration = 105,
  sym__attrib = 106,
  sym__expression_list = 107,
  sym_expression = 108,
  sym_string = 109,
  sym__quote_string = 110,
  aux_sym__doublequote_string_content = 111,
  aux_sym__singlequote_string_content = 112,
  sym__block_string = 113,
  sym_function_definition = 114,
  sym__function_body = 115,
  sym_parameters = 116,
  sym__parameter_list = 117,
  sym__vararg_parameter = 118,
  sym__prefix_expression = 119,
  sym_variable = 120,
  sym_bracket_index_expression = 121,
  sym_dot_index_expression = 122,
  sym_function_call = 123,
  sym_method_index_expression = 124,
  sym_arguments = 125,
  sym_parenthesized_expression = 126,
  sym_table_constructor = 127,
  sym__field_list = 128,
  sym__field_sep = 129,
  sym_field = 130,
  sym_binary_expression = 131,
  sym_unary_expression = 132,
  sym_comment = 133,
  sym__contextual_keyword = 134,
  aux_sym_chunk_repeat1 = 135,
  aux_sym__variable_assignment_varlist_repeat1 = 136,
  aux_sym__variable_assignment_explist_repeat1 = 137,
  aux_sym_if_statement_repeat1 = 138,
  aux_sym__name_list_repeat1 = 139,
  aux_sym__att_name_list_repeat1 = 140,
  aux_sym__expression_list_repeat1 = 141,
  aux_sym__field_list_repeat1 = 142,
}
enum ts_field_identifiers {
  field_alternative = 1,
  field_arguments = 2,
  field_attribute = 3,
  field_body = 4,
  field_clause = 5,
  field_condition = 6,
  field_consequence = 7,
  field_content = 8,
  field_end = 9,
  field_field = 10,
  field_global_declaration = 11,
  field_left = 12,
  field_local_declaration = 13,
  field_method = 14,
  field_name = 15,
  field_operand = 16,
  field_operator = 17,
  field_parameters = 18,
  field_right = 19,
  field_start = 20,
  field_step = 21,
  field_table = 22,
  field_value = 23,
}
static bool ts_lex(TSLexer *lexer, TSStateId state)
static bool ts_lex_keywords(TSLexer *lexer, TSStateId state)
enum ts_external_scanner_symbol_identifiers {
  ts_external_token__block_comment_start = 0,
  ts_external_token__block_comment_content = 1,
  ts_external_token__block_comment_end = 2,
  ts_external_token__block_string_start = 3,
  ts_external_token__block_string_content = 4,
  ts_external_token__block_string_end = 5,
}
void *tree_sitter_lua_external_scanner_create(void);
void tree_sitter_lua_external_scanner_destroy(void *);
bool tree_sitter_lua_external_scanner_scan(void *, TSLexer *, const bool *);
unsigned tree_sitter_lua_external_scanner_serialize(void *, char *);
void tree_sitter_lua_external_scanner_deserialize(void *, const char *, unsigned);
#define TS_PUBLIC
#define TS_PUBLIC __declspec(dllexport)
#define TS_PUBLIC __attribute__((visibility("default")))
TS_PUBLIC const TSLanguage *tree_sitter_lua(void)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/scanner.c

**Imports:**
- `#include <stdio.h>`
- `#include "tree_sitter/alloc.h"`
- `#include "tree_sitter/parser.h"`
- `#include <wctype.h>`

```c
enum TokenType {
  BLOCK_COMMENT_START,
  BLOCK_COMMENT_CONTENT,
  BLOCK_COMMENT_END,

  BLOCK_STRING_START,
  BLOCK_STRING_CONTENT,
  BLOCK_STRING_END,
}
static inline void consume(TSLexer *lexer)
static inline void skip(TSLexer *lexer)
static inline bool consume_char(char c, TSLexer *lexer)
static inline uint8_t consume_and_count_char(char c, TSLexer *lexer)
static inline void skip_whitespaces(TSLexer *lexer)
typedef struct {
  char ending_char;
  uint8_t level_count;
} Scanner;
static inline void reset_state(Scanner *scanner)
void *tree_sitter_lua_external_scanner_create()
void tree_sitter_lua_external_scanner_destroy(void *payload)
unsigned tree_sitter_lua_external_scanner_serialize(void *payload, char *buffer)
void tree_sitter_lua_external_scanner_deserialize(void *payload, const char *buffer, unsigned length)
static bool scan_block_start(Scanner *scanner, TSLexer *lexer)
static bool scan_block_end(Scanner *scanner, TSLexer *lexer)
static bool scan_block_content(Scanner *scanner, TSLexer *lexer)
static bool scan_comment_start(Scanner *scanner, TSLexer *lexer)
static bool scan_comment_content(Scanner *scanner, TSLexer *lexer)
bool tree_sitter_lua_external_scanner_scan(void *payload, TSLexer *lexer, const bool *valid_symbols)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/tree_sitter/alloc.h

**Imports:**
- `#include <stdbool.h>`
- `#include <stdio.h>`
- `#include <stdlib.h>`

```cpp
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/tree_sitter/array.h

**Imports:**
- `#include "./alloc.h"`
- `#include <assert.h>`
- `#include <stdbool.h>`
- `#include <stdint.h>`
- `#include <stdlib.h>`
- `#include <string.h>`

```cpp
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity)        \
  ((self)->contents = _array__reserve(           \
    (void *)(self)->contents, &(self)->capacity, \
    array_elem_size(self), new_capacity)         \
  )
#define array_delete(self)                           \
  do {                                               \
    if ((self)->contents) ts_free((self)->contents); \
    (self)->contents = NULL;                         \
    (self)->size = 0;                                \
    (self)->capacity = 0;                            \
  } while (0)
#define array_push(self, element)                                 \
  do {                                                            \
    (self)->contents = _array__grow(                              \
      (void *)(self)->contents, (self)->size, &(self)->capacity,  \
      1, array_elem_size(self)                                    \
    );                                                            \
   (self)->contents[(self)->size++] = (element);                  \
  } while(0)
#define array_grow_by(self, count)                                               \
  do {                                                                           \
    if ((count) == 0) break;                                                     \
    (self)->contents = _array__grow(                                             \
      (self)->contents, (self)->size, &(self)->capacity,                         \
      count, array_elem_size(self)                                               \
    );                                                                           \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count);                                                     \
  } while (0)
#define array_push_all(self, other) \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, other_contents)                 \
  (self)->contents = _array__splice(                              \
    (void*)(self)->contents, &(self)->size, &(self)->capacity,    \
    array_elem_size(self), (self)->size, 0, count, other_contents \
  )
#define array_splice(self, _index, old_count, new_count, new_contents) \
  (self)->contents = _array__splice(                                   \
    (void *)(self)->contents, &(self)->size, &(self)->capacity,        \
    array_elem_size(self), _index, old_count, new_count, new_contents  \
  )
#define array_insert(self, _index, element)                     \
  (self)->contents = _array__splice(                            \
    (void *)(self)->contents, &(self)->size, &(self)->capacity, \
    array_elem_size(self), _index, 0, 1, &(element)             \
  )
#define array_erase(self, _index) \
  _array__erase((void *)(self)->contents, &(self)->size, array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other)                                   \
  (self)->contents = _array__assign(                                \
    (void *)(self)->contents, &(self)->size, &(self)->capacity,     \
    (const void *)(other)->contents, (other)->size, array_elem_size(self) \
  )
#define array_swap(self, other)                                     \
  do {                                                              \
    void *_array_swap_tmp = (void *)(self)->contents;               \
    (self)->contents = (other)->contents;                           \
    (other)->contents = _array_swap_tmp;                            \
    _array__swap(&(self)->size, &(self)->capacity,                  \
                 &(other)->size, &(other)->capacity);               \
  } while (0)
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__erase(void* self_contents, uint32_t *size,
                                size_t element_size, uint32_t index)
static inline void *_array__reserve(void *contents, uint32_t *capacity,
                                  size_t element_size, uint32_t new_capacity)
static inline void *_array__assign(void* self_contents, uint32_t *self_size, uint32_t *self_capacity,
                                 const void *other_contents, uint32_t other_size, size_t element_size)
static inline void _array__swap(uint32_t *self_size, uint32_t *self_capacity,
                               uint32_t *other_size, uint32_t *other_capacity)
static inline void *_array__grow(void *contents, uint32_t size, uint32_t *capacity,
                               uint32_t count, size_t element_size)
static inline void *_array__splice(void *self_contents, uint32_t *size, uint32_t *capacity,
                                 size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/lua/tree_sitter/parser.h

**Imports:**
- `#include <stdbool.h>`
- `#include <stdint.h>`
- `#include <stdlib.h>`

```cpp
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
struct TSLanguageMetadata
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
  uint16_t reserved_word_set_id;
} TSLexerMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(const TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/swift/binding.go

**Imports:**
- `import "C"`
- `import "unsafe"`

```go
func Language() unsafe.Pointer
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/swift/scanner.c

**Imports:**
- `#include "tree_sitter/parser.h"`
- `#include <string.h>`
- `#include <wctype.h>`

```c
#define TOKEN_COUNT 33
enum TokenType {
    BLOCK_COMMENT,
    RAW_STR_PART,
    RAW_STR_CONTINUING_INDICATOR,
    RAW_STR_END_PART,
    IMPLICIT_SEMI,
    EXPLICIT_SEMI,
    ARROW_OPERATOR,
    DOT_OPERATOR,
    CONJUNCTION_OPERATOR,
    DISJUNCTION_OPERATOR,
    NIL_COALESCING_OPERATOR,
    EQUAL_SIGN,
    EQ_EQ,
    PLUS_THEN_WS,
    MINUS_THEN_WS,
    BANG,
    THROWS_KEYWORD,
    RETHROWS_KEYWORD,
    DEFAULT_KEYWORD,
    WHERE_KEYWORD,
    ELSE_KEYWORD,
    CATCH_KEYWORD,
    AS_KEYWORD,
    AS_QUEST,
    AS_BANG,
    ASYNC_KEYWORD,
    CUSTOM_OPERATOR,
    HASH_SYMBOL,
    DIRECTIVE_IF,
    DIRECTIVE_ELSEIF,
    DIRECTIVE_ELSE,
    DIRECTIVE_ENDIF,
    FAKE_TRY_BANG
}
#define OPERATOR_COUNT 20
enum IllegalTerminatorGroup {
    ALPHANUMERIC,
    OPERATOR_SYMBOLS,
    OPERATOR_OR_DOT,
    NON_WHITESPACE
}
enum IllegalTerminatorGroup
enum TokenType
#define RESERVED_OP_COUNT 31
static bool is_cross_semi_token(enum TokenType op)
#define NON_CONSUMING_CROSS_SEMI_CHAR_COUNT 3
enum ParseDirective {
    CONTINUE_PARSING_NOTHING_FOUND,
    CONTINUE_PARSING_TOKEN_FOUND,
    CONTINUE_PARSING_SLASH_CONSUMED,
    STOP_PARSING_NOTHING_FOUND,
    STOP_PARSING_TOKEN_FOUND,
    STOP_PARSING_END_OF_FILE
}
struct ScannerState {
    uint32_t ongoing_raw_str_hash_count;
}
void *tree_sitter_swift_external_scanner_create()
struct ScannerState
void tree_sitter_swift_external_scanner_destroy(void *payload)
void tree_sitter_swift_external_scanner_reset(void *payload)
struct ScannerState
unsigned tree_sitter_swift_external_scanner_serialize(void *payload, char *buffer)
struct ScannerState
void tree_sitter_swift_external_scanner_deserialize(
    void *payload,
    const char *buffer,
    unsigned length
)
struct ScannerState
static void advance(TSLexer *lexer)
static bool should_treat_as_wspace(int32_t character)
static int32_t encountered_op_count(bool *encountered_operator)
static bool any_reserved_ops(uint8_t *encountered_reserved_ops)
static bool is_legal_custom_operator(
    int32_t char_idx,
    int32_t first_char,
    int32_t cur_char
)
static bool eat_operators(
    TSLexer *lexer,
    const bool *valid_symbols,
    bool mark_end,
    const int32_t prior_char,
    enum TokenType *symbol_result
)
enum TokenType
enum IllegalTerminatorGroup
enum ParseDirective
enum TokenType
enum ParseDirective
enum TokenType
enum ParseDirective
enum ParseDirective
enum TokenType
enum TokenType
#define DIRECTIVE_COUNT 4
enum TokenType
enum TokenType
static bool eat_raw_str_part(
    struct ScannerState *state,
    TSLexer *lexer,
    const bool *valid_symbols,
    enum TokenType *symbol_result
)
struct ScannerState
enum TokenType
bool tree_sitter_swift_external_scanner_scan(
    void *payload,
    TSLexer *lexer,
    const bool *valid_symbols
)
struct ScannerState
enum TokenType
enum ParseDirective
enum TokenType
enum ParseDirective
enum TokenType
enum TokenType
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/swift/tree_sitter/alloc.h

**Imports:**
- `#include <stdbool.h>`
- `#include <stdio.h>`
- `#include <stdlib.h>`

```cpp
#define TREE_SITTER_ALLOC_H_
#define ts_malloc  ts_current_malloc
#define ts_calloc  ts_current_calloc
#define ts_realloc ts_current_realloc
#define ts_free    ts_current_free
#define ts_malloc  malloc
#define ts_calloc  calloc
#define ts_realloc realloc
#define ts_free    free
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/swift/tree_sitter/array.h

**Imports:**
- `#include "./alloc.h"`
- `#include <assert.h>`
- `#include <stdbool.h>`
- `#include <stdint.h>`
- `#include <stdlib.h>`
- `#include <string.h>`

```cpp
#define TREE_SITTER_ARRAY_H_
#define Array(T)       \
  struct {             \
    T *contents;       \
    uint32_t size;     \
    uint32_t capacity; \
  }
#define array_init(self) \
  ((self)->size = 0, (self)->capacity = 0, (self)->contents = NULL)
#define array_new() \
  { NULL, 0, 0 }
#define array_get(self, _index) \
  (assert((uint32_t)(_index) < (self)->size), &(self)->contents[_index])
#define array_front(self) array_get(self, 0)
#define array_back(self) array_get(self, (self)->size - 1)
#define array_clear(self) ((self)->size = 0)
#define array_reserve(self, new_capacity) \
  _array__reserve((Array *)(self), array_elem_size(self), new_capacity)
#define array_delete(self) _array__delete((Array *)(self))
#define array_push(self, element)                            \
  (_array__grow((Array *)(self), 1, array_elem_size(self)), \
   (self)->contents[(self)->size++] = (element))
#define array_grow_by(self, count) \
  do { \
    if ((count) == 0) break; \
    _array__grow((Array *)(self), count, array_elem_size(self)); \
    memset((self)->contents + (self)->size, 0, (count) * array_elem_size(self)); \
    (self)->size += (count); \
  } while (0)
#define array_push_all(self, other)                                       \
  array_extend((self), (other)->size, (other)->contents)
#define array_extend(self, count, contents)                    \
  _array__splice(                                               \
    (Array *)(self), array_elem_size(self), (self)->size, \
    0, count,  contents                                        \
  )
#define array_splice(self, _index, old_count, new_count, new_contents)  \
  _array__splice(                                                       \
    (Array *)(self), array_elem_size(self), _index,                \
    old_count, new_count, new_contents                                 \
  )
#define array_insert(self, _index, element) \
  _array__splice((Array *)(self), array_elem_size(self), _index, 0, 1, &(element))
#define array_erase(self, _index) \
  _array__erase((Array *)(self), array_elem_size(self), _index)
#define array_pop(self) ((self)->contents[--(self)->size])
#define array_assign(self, other) \
  _array__assign((Array *)(self), (const Array *)(other), array_elem_size(self))
#define array_swap(self, other) \
  _array__swap((Array *)(self), (Array *)(other))
#define array_elem_size(self) (sizeof *(self)->contents)
#define array_search_sorted_with(self, compare, needle, _index, _exists) \
  _array__search_sorted(self, 0, compare, , needle, _index, _exists)
#define array_search_sorted_by(self, field, needle, _index, _exists) \
  _array__search_sorted(self, 0, _compare_int, field, needle, _index, _exists)
#define array_insert_sorted_with(self, compare, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_with(self, compare, &(value), &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
#define array_insert_sorted_by(self, field, value) \
  do { \
    unsigned _index, _exists; \
    array_search_sorted_by(self, field, (value) field, &_index, &_exists); \
    if (!_exists) array_insert(self, _index, value); \
  } while (0)
static inline void _array__delete(Array *self)
static inline void _array__erase(Array *self, size_t element_size,
                                uint32_t index)
static inline void _array__reserve(Array *self, size_t element_size, uint32_t new_capacity)
static inline void _array__assign(Array *self, const Array *other, size_t element_size)
static inline void _array__swap(Array *self, Array *other)
static inline void _array__grow(Array *self, uint32_t count, size_t element_size)
static inline void _array__splice(Array *self, size_t element_size,
                                 uint32_t index, uint32_t old_count,
                                 uint32_t new_count, const void *elements)
#define _array__search_sorted(self, start, compare, suffix, needle, _index, _exists) \
  do { \
    *(_index) = start; \
    *(_exists) = false; \
    uint32_t size = (self)->size - *(_index); \
    if (size == 0) break; \
    int comparison; \
    while (size > 1) { \
      uint32_t half_size = size / 2; \
      uint32_t mid_index = *(_index) + half_size; \
      comparison = compare(&((self)->contents[mid_index] suffix), (needle)); \
      if (comparison <= 0) *(_index) = mid_index; \
      size -= half_size; \
    } \
    comparison = compare(&((self)->contents[*(_index)] suffix), (needle)); \
    if (comparison == 0) *(_exists) = true; \
    else if (comparison < 0) *(_index) += 1; \
  } while (0)
#define _compare_int(a, b) ((int)*(a) - (int)(b))
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/grammars/swift/tree_sitter/parser.h

**Imports:**
- `#include <stdbool.h>`
- `#include <stdint.h>`
- `#include <stdlib.h>`

```cpp
#define TREE_SITTER_PARSER_H_
#define ts_builtin_sym_error ((TSSymbol)-1)
#define ts_builtin_sym_end 0
#define TREE_SITTER_SERIALIZATION_BUFFER_SIZE 1024
typedef uint16_t TSStateId;
typedef uint16_t TSSymbol;
typedef uint16_t TSFieldId;
struct TSLanguage
typedef struct {
  TSFieldId field_id;
  uint8_t child_index;
  bool inherited;
} TSFieldMapEntry;
typedef struct {
  uint16_t index;
  uint16_t length;
} TSFieldMapSlice;
typedef struct {
  bool visible;
  bool named;
  bool supertype;
} TSSymbolMetadata;
struct TSLexer
struct TSLexer
typedef enum {
  TSParseActionTypeShift,
  TSParseActionTypeReduce,
  TSParseActionTypeAccept,
  TSParseActionTypeRecover,
} TSParseActionType;
typedef union {
  struct {
    uint8_t type;
    TSStateId state;
    bool extra;
    bool repetition;
  } shift;
  struct {
    uint8_t type;
    uint8_t child_count;
    TSSymbol symbol;
    int16_t dynamic_precedence;
    uint16_t production_id;
  } reduce;
  uint8_t type;
} TSParseAction;
typedef struct {
  uint16_t lex_state;
  uint16_t external_lex_state;
} TSLexMode;
typedef union {
  TSParseAction action;
  struct {
    uint8_t count;
    bool reusable;
  } entry;
} TSParseActionEntry;
typedef struct {
  int32_t start;
  int32_t end;
} TSCharacterRange;
struct TSLanguage
static inline bool set_contains(TSCharacterRange *ranges, uint32_t len, int32_t lookahead)
#define UNUSED __pragma(warning(suppress : 4101))
#define UNUSED __attribute__((unused))
#define START_LEXER()           \
  bool result = false;          \
  bool skip = false;            \
  UNUSED                        \
  bool eof = false;             \
  int32_t lookahead;            \
  goto start;                   \
  next_state:                   \
  lexer->advance(lexer, skip);  \
  start:                        \
  skip = false;                 \
  lookahead = lexer->lookahead;
#define ADVANCE(state_value) \
  {                          \
    state = state_value;     \
    goto next_state;         \
  }
#define ADVANCE_MAP(...)                                              \
  {                                                                   \
    static const uint16_t map[] = { __VA_ARGS__ };                    \
    for (uint32_t i = 0; i < sizeof(map) / sizeof(map[0]); i += 2) {  \
      if (map[i] == lookahead) {                                      \
        state = map[i + 1];                                           \
        goto next_state;                                              \
      }                                                               \
    }                                                                 \
  }
#define SKIP(state_value) \
  {                       \
    skip = true;          \
    state = state_value;  \
    goto next_state;      \
  }
#define ACCEPT_TOKEN(symbol_value)     \
  result = true;                       \
  lexer->result_symbol = symbol_value; \
  lexer->mark_end(lexer);
#define END_STATE() return result;
#define SMALL_STATE(id) ((id) - LARGE_STATE_COUNT)
#define STATE(id) id
#define ACTIONS(id) id
#define SHIFT(state_value)            \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value)          \
    }                                 \
  }}
#define SHIFT_REPEAT(state_value)     \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .state = (state_value),         \
      .repetition = true              \
    }                                 \
  }}
#define SHIFT_EXTRA()                 \
  {{                                  \
    .shift = {                        \
      .type = TSParseActionTypeShift, \
      .extra = true                   \
    }                                 \
  }}
#define REDUCE(symbol_name, children, precedence, prod_id) \
  {{                                                       \
    .reduce = {                                            \
      .type = TSParseActionTypeReduce,                     \
      .symbol = symbol_name,                               \
      .child_count = children,                             \
      .dynamic_precedence = precedence,                    \
      .production_id = prod_id                             \
    },                                                     \
  }}
#define RECOVER()                    \
  {{                                 \
    .type = TSParseActionTypeRecover \
  }}
#define ACCEPT_INPUT()              \
  {{                                \
    .type = TSParseActionTypeAccept \
  }}
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

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/csharp.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_c_sharp "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/csharp"`

```go
type CSharpQuery struct {
	language *sitter.Language
	query    []byte
}
func NewCSharpQuery() *CSharpQuery
func (q *CSharpQuery) Language() *sitter.Language
func (q *CSharpQuery) Query() []byte
func (q *CSharpQuery) Captures() []string
func (q *CSharpQuery) KindMapping() map[string]string
func (q *CSharpQuery) ImportQuery() []byte
csharpImportQueryPattern = `
; using directives (capture full declaration)
(using_directive) @import_path
`
csharpQueryPattern = `
; Class declarations
(class_declaration
  name: (identifier) @name
) @signature @kind

; Struct declarations
(struct_declaration
  name: (identifier) @name
) @signature @kind

; Interface declarations
(interface_declaration
  name: (identifier) @name
) @signature @kind

; Enum declarations
(enum_declaration
  name: (identifier) @name
) @signature @kind

; Record declarations (record, record class, record struct)
(record_declaration
  name: (identifier) @name
) @signature @kind

; Delegate declarations
(delegate_declaration
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

; Destructor declarations
(destructor_declaration
  name: (identifier) @name
) @signature @kind

; Property declarations
(property_declaration
  name: (identifier) @name
) @signature @kind

; Field declarations (static/const filtered in parser.go)
(field_declaration
  (variable_declaration
    (variable_declarator
      name: (identifier) @name
    )
  )
) @signature @kind

; Event field declarations (e.g., public event EventHandler Changed;)
(event_field_declaration
  (variable_declaration
    (variable_declarator
      name: (identifier) @name
    )
  )
) @signature @kind

; Event declarations with accessor body
(event_declaration
  name: (identifier) @name
) @signature @kind

; Indexer declarations (no name capture — synthesized in parser.go)
(indexer_declaration) @signature @kind

; Operator declarations (no name capture — synthesized in parser.go)
(operator_declaration) @signature @kind

; Conversion operator declarations (no name capture — synthesized in parser.go)
(conversion_operator_declaration) @signature @kind

; Namespace declarations
(namespace_declaration
  name: (_) @name
) @signature @kind

; File-scoped namespace declarations (C# 10+)
(file_scoped_namespace_declaration
  name: (_) @name
) @signature @kind

; Enum member declarations
(enum_member_declaration
  name: (identifier) @name
) @signature @kind

; Comments (XML doc comments and regular)
(comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/csharp_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_c_sharp "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/csharp"`

```go
func extractCSharpNames(t *testing.T, code []byte) map[string]bool
func TestCSharpQueryLanguage(t *testing.T)
func TestCSharpQueryPattern(t *testing.T)
func TestCSharpQueryImportPattern(t *testing.T)
func TestCSharpQueryExtractFunction(t *testing.T)
func TestCSharpQueryExtractTypes(t *testing.T)
func TestCSharpQueryExtractConstructorDestructor(t *testing.T)
func TestCSharpQueryExtractProperties(t *testing.T)
func TestCSharpQueryExtractFields(t *testing.T)
func TestCSharpQueryExtractEvents(t *testing.T)
func TestCSharpQueryExtractOperators(t *testing.T)
foundOperator, foundConversion, foundIndexer bool
func TestCSharpQueryExtractImport(t *testing.T)
func TestCSharpQueryExtractGenerics(t *testing.T)
func TestCSharpQueryKindMapping(t *testing.T)
func TestCSharpQueryCaptures(t *testing.T)
func TestCSharpQueryExtractNamespace(t *testing.T)
func TestCSharpQueryExtractRecords(t *testing.T)
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
; Import declarations (capture entire declaration)
(import_declaration) @import_path
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

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/kotlin.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_kotlin "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/kotlin"`

```go
type KotlinQuery struct {
	language *sitter.Language
	query    []byte
}
func NewKotlinQuery() *KotlinQuery
func (q *KotlinQuery) Language() *sitter.Language
func (q *KotlinQuery) Query() []byte
func (q *KotlinQuery) Captures() []string
func (q *KotlinQuery) KindMapping() map[string]string
func (q *KotlinQuery) ImportQuery() []byte
kotlinImportQueryPattern = `
; Import statements
(import_header) @import_path
`
kotlinQueryPattern = `
; Function declarations (regular, suspend, inline, extension, operator, infix, tailrec)
(function_declaration
  (simple_identifier) @name
) @signature @kind

; Class declarations (class, data class, sealed class, enum class, interface, annotation class, value class)
(class_declaration
  (type_identifier) @name
) @signature @kind

; Object declarations (singleton)
(object_declaration
  (type_identifier) @name
) @signature @kind

; Companion object with explicit name (e.g., companion object Factory)
(companion_object
  (type_identifier) @name
) @signature @kind

; Property declarations (val/var, const val, lateinit, delegated)
(property_declaration
  (variable_declaration
    (simple_identifier) @name
  )
) @signature @kind

; Type alias
(type_alias
  (type_identifier) @name
) @signature @kind

; Enum entries
(enum_entry
  (simple_identifier) @name
) @signature @kind

; Secondary constructors
(secondary_constructor) @signature @kind

; Line comments
(line_comment) @doc

; Block/multiline comments
(multiline_comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/kotlin_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_kotlin "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/kotlin"`

```go
func extractKotlinNames(t *testing.T, code []byte) map[string]bool
func TestKotlinQueryLanguage(t *testing.T)
func TestKotlinQueryPattern(t *testing.T)
func TestKotlinQueryImportPattern(t *testing.T)
func TestKotlinQueryExtractFunction(t *testing.T)
func TestKotlinQueryExtractTypes(t *testing.T)
func TestKotlinQueryExtractInterface(t *testing.T)
func TestKotlinQueryExtractObject(t *testing.T)
func TestKotlinQueryExtractProperties(t *testing.T)
func TestKotlinQueryExtractTypeAlias(t *testing.T)
func TestKotlinQueryExtractEnumEntry(t *testing.T)
func TestKotlinQueryExtractImport(t *testing.T)
func TestKotlinQueryExtractGenerics(t *testing.T)
func TestKotlinQueryKindMapping(t *testing.T)
func TestKotlinQueryCaptures(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/lua.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_lua "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/lua"`

```go
type LuaQuery struct {
	language *sitter.Language
	query    []byte
}
func NewLuaQuery() *LuaQuery
func (q *LuaQuery) Language() *sitter.Language
func (q *LuaQuery) Query() []byte
func (q *LuaQuery) Captures() []string
func (q *LuaQuery) KindMapping() map[string]string
func (q *LuaQuery) ImportQuery() []byte
luaImportQueryPattern = `
; local json = require("json")
(variable_declaration
  (assignment_statement
    (expression_list
      value: (function_call
        name: (identifier) @_fn
        arguments: (arguments (string))
      )
    )
  )
) @import_path
(#eq? @_fn "require")
`
luaQueryPattern = `
; Function declarations (global, local, module, method)
; Covers: function foo(), local function foo(), function M.foo(), function M:foo()
(function_declaration
  name: [
    (identifier) @name
    (dot_index_expression field: (identifier) @name)
    (method_index_expression method: (identifier) @name)
  ]
) @signature @kind

; Variable declarations with function assignment: local foo = function() end
(variable_declaration
  (assignment_statement
    (variable_list
      name: (identifier) @name)
    (expression_list
      value: (function_definition))
  )
) @signature @kind

; Variable declarations with table constructor: local M = {}
(variable_declaration
  (assignment_statement
    (variable_list
      name: (identifier) @name)
    (expression_list
      value: (table_constructor))
  )
) @signature @kind

; Comments (LuaDoc --- and regular --)
(comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/lua_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_lua "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/lua"`

```go
func extractLuaNames(t *testing.T, code []byte) map[string]bool
func TestLuaQueryLanguage(t *testing.T)
func TestLuaQueryPattern(t *testing.T)
func TestLuaQueryImportPattern(t *testing.T)
func TestLuaQueryExtractFunction(t *testing.T)
func TestLuaQueryExtractLocalFunction(t *testing.T)
func TestLuaQueryExtractModuleFunction(t *testing.T)
func TestLuaQueryExtractMethod(t *testing.T)
func TestLuaQueryExtractTableAssignment(t *testing.T)
func TestLuaQueryExtractFunctionAssignment(t *testing.T)
func TestLuaQueryExtractImport(t *testing.T)
func TestLuaQueryNonRequireFalsePositive(t *testing.T)
func TestLuaQueryExtractDoc(t *testing.T)
func TestLuaQueryKindMapping(t *testing.T)
func TestLuaQueryCaptures(t *testing.T)
func TestLuaQueryExtractMixed(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/php.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_php "github.com/tree-sitter/tree-sitter-php/bindings/go"`

```go
type PHPQuery struct {
	language *sitter.Language
	query    []byte
}
func NewPHPQuery() *PHPQuery
func (q *PHPQuery) Language() *sitter.Language
func (q *PHPQuery) Query() []byte
func (q *PHPQuery) Captures() []string
func (q *PHPQuery) KindMapping() map[string]string
func (q *PHPQuery) ImportQuery() []byte
phpImportQueryPattern = `
; use Namespace\\Class;
(namespace_use_declaration) @import_path

; include 'file.php';
(include_expression) @import_path

; require 'vendor/autoload.php';
(require_expression) @import_path

; include_once 'config.php';
(include_once_expression) @import_path

; require_once 'config.php';
(require_once_expression) @import_path
`
phpQueryPattern = `
; Function definitions: function name() {}
(function_definition name: (name) @name) @signature @kind

; Method declarations in classes
(method_declaration name: (name) @name) @signature @kind

; Class declarations
(class_declaration name: (name) @name) @signature @kind

; Interface declarations
(interface_declaration name: (name) @name) @signature @kind

; Trait declarations
(trait_declaration name: (name) @name) @signature @kind

; Enum declarations
(enum_declaration name: (name) @name) @signature @kind

; Const declarations: const NAME = value;
(const_declaration
  (const_element
    (name) @name
  )
) @signature @kind

; Property declarations: public $name;
(property_declaration
  (property_element
    (variable_name (name) @name)
  )
) @signature @kind

; Comments (PHPDoc and regular)
(comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/php_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_php "github.com/tree-sitter/tree-sitter-php/bindings/go"`

```go
func extractPHPNames(t *testing.T, code []byte) map[string]bool
func TestPHPQueryLanguage(t *testing.T)
func TestPHPQueryPattern(t *testing.T)
func TestPHPQueryImportPattern(t *testing.T)
func TestPHPQueryExtractFunction(t *testing.T)
func TestPHPQueryExtractClass(t *testing.T)
func TestPHPQueryExtractInterface(t *testing.T)
func TestPHPQueryExtractTrait(t *testing.T)
func TestPHPQueryExtractEnum(t *testing.T)
func TestPHPQueryExtractVariable(t *testing.T)
func TestPHPQueryExtractImport(t *testing.T)
func TestPHPQueryKindMapping(t *testing.T)
func TestPHPQueryCaptures(t *testing.T)
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

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/ruby.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_ruby "github.com/tree-sitter/tree-sitter-ruby/bindings/go"`

```go
type RubyQuery struct {
	language *sitter.Language
	query    []byte
}
func NewRubyQuery() *RubyQuery
func (q *RubyQuery) Language() *sitter.Language
func (q *RubyQuery) Query() []byte
func (q *RubyQuery) Captures() []string
func (q *RubyQuery) KindMapping() map[string]string
func (q *RubyQuery) ImportQuery() []byte
rubyImportQueryPattern = `
; require "library" / require_relative "library"
(call
  method: (identifier) @_fn
  arguments: (argument_list
    (string)
  )
) @import_path
(#match? @_fn "^require")
`
rubyQueryPattern = `
; Instance methods and top-level functions: def foo(args) ... end
(method
  name: (identifier) @name
) @signature @kind

; Class methods: def self.foo(args) ... end
(singleton_method
  name: (identifier) @name
) @signature @kind

; Class definitions: class Foo ... end
(class
  name: (constant) @name
) @signature @kind

; Module definitions: module Foo ... end
(module
  name: (constant) @name
) @signature @kind

; Top-level constant assignments: FOO = value
(program
  (assignment
    left: (constant) @name
  ) @signature @kind
)

; Comments
(comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/ruby_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_ruby "github.com/tree-sitter/tree-sitter-ruby/bindings/go"`

```go
func extractRubyNames(t *testing.T, code []byte) map[string]bool
func extractRubyImports(t *testing.T, code []byte) []string
imports []string
func TestRubyQueryLanguage(t *testing.T)
func TestRubyQueryPattern(t *testing.T)
func TestRubyQueryImportPattern(t *testing.T)
func TestRubyQueryExtractFunction(t *testing.T)
func TestRubyQueryExtractTypes(t *testing.T)
func TestRubyQueryExtractClassMethods(t *testing.T)
func TestRubyQueryExtractModuleMethods(t *testing.T)
func TestRubyQueryExtractConstants(t *testing.T)
func TestRubyQueryExtractAttrAccessors(t *testing.T)
func TestRubyQueryExtractImport(t *testing.T)
func TestRubyQueryExtractNestedClasses(t *testing.T)
func TestRubyQueryKindMapping(t *testing.T)
func TestRubyQueryCaptures(t *testing.T)
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
}
func NewRustQuery() *RustQuery
func (q *RustQuery) Language() *sitter.Language
func (q *RustQuery) Query() []byte
func (q *RustQuery) Captures() []string
func (q *RustQuery) KindMapping() map[string]string
func (q *RustQuery) ImportQuery() []byte
rustImportQueryPattern = `
; Use declarations (capture full statement)
(use_declaration) @import_path

; Extern crate declarations
(extern_crate_declaration) @import_path
`
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
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/rust_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"`

```go
func TestRustQueryLanguage(t *testing.T)
func TestRustQueryPattern(t *testing.T)
func TestRustQueryImportPattern(t *testing.T)
func TestRustQueryExtractFunction(t *testing.T)
func TestRustQueryExtractTypes(t *testing.T)
func TestRustQueryExtractImplAndMethods(t *testing.T)
func TestRustQueryExtractConstAndStatic(t *testing.T)
func TestRustQueryExtractMacro(t *testing.T)
func TestRustQueryExtractModule(t *testing.T)
func TestRustQueryExtractUse(t *testing.T)
func TestRustQueryExtractGenericsAndLifetimes(t *testing.T)
func TestRustQueryEmptyFile(t *testing.T)
func TestRustQueryKindMapping(t *testing.T)
func TestRustQueryCaptures(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/shell.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_bash "github.com/tree-sitter/tree-sitter-bash/bindings/go"`

```go
type ShellQuery struct {
	language *sitter.Language
	query    []byte
}
func NewShellQuery() *ShellQuery
func (q *ShellQuery) Language() *sitter.Language
func (q *ShellQuery) Query() []byte
func (q *ShellQuery) Captures() []string
func (q *ShellQuery) KindMapping() map[string]string
func (q *ShellQuery) ImportQuery() []byte
shellImportQueryPattern = `
; source /path/to/file and . /path/to/file
; Capture command nodes. Go-side filtering will check if command name is "source" or "."
(command
  name: (command_name) @name
) @import_path
`
shellQueryPattern = `
; Function definitions: function foo { } or function foo() { } or foo() { }
(function_definition
  name: (word) @name
) @signature @kind

; Variable assignments: FOO=bar, FOO="bar", FOO=$(cmd)
(variable_assignment
  name: (variable_name) @name
) @signature @kind

; Comments (# ...)
(comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/shell_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_bash "github.com/tree-sitter/tree-sitter-bash/bindings/go"`

```go
func extractShellNames(t *testing.T, code []byte) map[string]bool
func TestShellQueryLanguage(t *testing.T)
func TestShellQueryPattern(t *testing.T)
func TestShellQueryImportPattern(t *testing.T)
func TestShellQueryExtractFunction(t *testing.T)
func TestShellQueryExtractFunctionWithoutKeyword(t *testing.T)
func TestShellQueryExtractVariable(t *testing.T)
func TestShellQueryExtractImport(t *testing.T)
func TestShellQueryExtractDoc(t *testing.T)
func TestShellQueryKindMapping(t *testing.T)
func TestShellQueryCaptures(t *testing.T)
func TestShellQueryExtractMixed(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/swift.go

**Imports:**
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_swift "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/swift"`

```go
type SwiftQuery struct {
	language *sitter.Language
	query    []byte
}
func NewSwiftQuery() *SwiftQuery
func (q *SwiftQuery) Language() *sitter.Language
func (q *SwiftQuery) Query() []byte
func (q *SwiftQuery) Captures() []string
func (q *SwiftQuery) KindMapping() map[string]string
func (q *SwiftQuery) ImportQuery() []byte
swiftImportQueryPattern = `
; Import declarations (capture full statement)
(import_declaration) @import_path
`
swiftQueryPattern = `
; Functions
(function_declaration
  name: (simple_identifier) @name
) @signature @kind

; Classes, Structs, Enums (all use class_declaration node type)
(class_declaration
  name: (type_identifier) @name
) @signature @kind

; Extensions (name is in user_type child)
(class_declaration
  name: (user_type
    (type_identifier) @name
  )
) @signature @kind

; Protocol declarations
(protocol_declaration
  name: (type_identifier) @name
) @signature @kind

; Type aliases
(typealias_declaration
  name: (type_identifier) @name
) @signature @kind

; Properties (let/var)
(property_declaration
  name: (pattern
    (simple_identifier) @name
  )
) @signature @kind

; Initializers
(init_declaration) @signature @kind

; Deinitializers
(deinit_declaration) @signature @kind

; Subscript declarations
(subscript_declaration) @signature @kind

; Operator declarations
(operator_declaration
  (custom_operator) @name
) @signature @kind

; Protocol function declarations (methods in protocol body)
(protocol_function_declaration
  name: (simple_identifier) @name
) @signature @kind

; Doc comments (/// style)
(comment) @doc

; Multiline comments (/** style)
(multiline_comment) @doc
`
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/parser/treesitter/languages/swift_test.go

**Imports:**
- `import "testing"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import tree_sitter_swift "github.com/indigo-net/Brf.it/pkg/parser/treesitter/grammars/swift"`

```go
func TestSwiftQueryLanguage(t *testing.T)
func TestSwiftQueryPattern(t *testing.T)
func TestSwiftQueryImportPattern(t *testing.T)
func TestSwiftQueryExtractFunction(t *testing.T)
func TestSwiftQueryExtractTypes(t *testing.T)
func TestSwiftQueryExtractProtocol(t *testing.T)
func TestSwiftQueryExtractExtension(t *testing.T)
func TestSwiftQueryExtractProperties(t *testing.T)
func TestSwiftQueryExtractInitDeinit(t *testing.T)
func TestSwiftQueryExtractSubscript(t *testing.T)
func TestSwiftQueryExtractImport(t *testing.T)
func TestSwiftQueryExtractGenerics(t *testing.T)
func TestSwiftQueryKindMapping(t *testing.T)
func TestSwiftQueryCaptures(t *testing.T)
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
) @import_path

; Export clause (barrel exports: export { foo, bar })
(export_statement
  (export_clause) @import_path
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
- `import "strings"`
- `import "sync"`
- `import sitter "github.com/tree-sitter/go-tree-sitter"`
- `import "github.com/indigo-net/Brf.it/pkg/parser"`
- `import "github.com/indigo-net/Brf.it/pkg/parser/treesitter/languages"`

```go
func init()
type queryType int
queryTypeSignature queryType = iota
queryTypeImport
type queryCacheKey struct {
	lang string
	typ  queryType
}
type TreeSitterParser struct {
	queries          map[string]LanguageQuery
	compiledQueries  sync.Map // map[queryCacheKey]*sitter.Query
	queryCacheMutex  sync.RWMutex
	parserPool       sync.Pool
	cursorPool       sync.Pool
}
func NewTreeSitterParser() *TreeSitterParser
func (p *TreeSitterParser) getOrCreateQuery(lang string, langQuery LanguageQuery, typ queryType) (*sitter.Query, error)
queryStr string
func (p *TreeSitterParser) Parse(content []byte, opts *parser.Options) (result *parser.ParseResult, err error)
rawImports []string
func (p *TreeSitterParser) Languages() []string
func (p *TreeSitterParser) extractSignatures(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
) ([]parser.Signature, error)
signatures []parser.Signature
kindNode *sitter.Node
func cleanComment(text string) string
func isExported(name, language string) bool
func stripBody(text, kind, language string) string
func stripGoBody(text, kind string) string
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
func stripRustBody(text, kind string) string
func findRustBodyStart(text string) int
func refineSwiftClassKind(text string) string
func stripSwiftBody(text, kind string) string
func findSwiftBodyStart(text string) int
func stripKotlinBody(text, kind string) string
func findKotlinBodyStart(text string) int
func refineKotlinClassKind(text string) string
func refineLuaFunctionKind(text string) string
func stripLuaBody(text, kind string) string
func stripPHPBody(text, kind string) string
func stripRubyBody(text, kind string) string
func stripShellBody(text, kind string) string
func findPHPBodyStart(text string) int
func stripCSharpBody(text, kind string) string
func findCSharpBodyStart(text string) int
func isExpressionBodied(text string) bool
func findCSharpArrowIndex(text string) int
func extractCSharpOperatorName(text string) string
func extractCSharpConversionOperatorName(text string) string
func (p *TreeSitterParser) extractImports(
	root *sitter.Node,
	content []byte,
	langQuery LanguageQuery,
	opts *parser.Options,
) ([]string, error)
imports []string
importNode *sitter.Node
func removeBlankLines(text string) string
buf strings.Builder
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
func TestTreeSitterParserParseRust(t *testing.T)
func TestRustSignatureOnlyExtraction(t *testing.T)
func TestRustImportExtraction(t *testing.T)
func TestRustAutoRegistration(t *testing.T)
func TestRustConstAndStaticExtraction(t *testing.T)
func TestRustMacroExtraction(t *testing.T)
func TestRustGenericsAndLifetimes(t *testing.T)
func TestTreeSitterParserParseSwift(t *testing.T)
func TestTreeSitterParserParseKotlin(t *testing.T)
func TestKotlinBodyStripping(t *testing.T)
func TestRefineKotlinClassKind(t *testing.T)
func TestKotlinAutoRegistration(t *testing.T)
func TestKotlinSignatureOnlyExtraction(t *testing.T)
func TestKotlinImportExtraction(t *testing.T)
func TestParsePanicRecovery(t *testing.T)
func TestTreeSitterParserParseCSharp(t *testing.T)
func TestCSharpSignatureOnlyExtraction(t *testing.T)
func TestCSharpOperatorNameSynthesis(t *testing.T)
func TestCSharpStaticFieldExtraction(t *testing.T)
func TestCSharpImportExtraction(t *testing.T)
func TestCSharpAutoRegistration(t *testing.T)
func TestCSharpBodyStripping(t *testing.T)
func TestFindCSharpBodyStart(t *testing.T)
func TestIsExpressionBodied(t *testing.T)
func TestExtractCSharpOperatorName(t *testing.T)
func TestExtractCSharpConversionOperatorName(t *testing.T)
func TestFindFunctionBodyStart(t *testing.T)
func TestParsePanicRecoveryMechanism(t *testing.T)
func TestTreeSitterParserParseLua(t *testing.T)
func TestTreeSitterParserParseLuaImports(t *testing.T)
func TestLuaBodyStripping(t *testing.T)
func TestRefineLuaFunctionKind(t *testing.T)
func TestIsPythonMethod(t *testing.T)
func TestTreeSitterParserParsePHP(t *testing.T)
foundClass, foundMethod, foundInterface, foundTrait, foundFunction, foundConst bool
func TestTreeSitterParserParsePHPImports(t *testing.T)
func TestPHPBodyStripping(t *testing.T)
func TestTreeSitterParserParseRuby(t *testing.T)
func TestTreeSitterParserParseRubyImports(t *testing.T)
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
CaptureLuaRequireFn = "_fn"
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
- `import "errors"`
- `import "fmt"`
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

	// Warnings contains non-fatal issues encountered during scanning.
	Warnings []string
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
func getBaseName(path string) string
type Scanner interface {
	// Scan performs the scan and returns scan results.
	Scan() (*ScanResult, error)
}
type FileScanner struct {
	opts             *ScanOptions
	ignorer          *ignore.GitIgnore
	ignorerErr       error
	ignorerErrWarned bool
	logger           *log.Logger
}
func NewFileScanner(opts *ScanOptions) (*FileScanner, error)
func (s *FileScanner) Scan() (*ScanResult, error)
warning string
func (s *FileScanner) checkFile(path string, info os.FileInfo) (FileEntry, bool)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/scanner/scanner_test.go

**Imports:**
- `import "bytes"`
- `import "log"`
- `import "os"`
- `import "path/filepath"`
- `import "strings"`
- `import "testing"`

```go
func TestGetBaseName(t *testing.T)
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
func TestScanGitignoreLoadFailureWarning(t *testing.T)
buf bytes.Buffer
buf bytes.Buffer
buf bytes.Buffer
func TestScanWalkDirPermissionDenied(t *testing.T)
buf bytes.Buffer
func TestScanSymlinkSkip(t *testing.T)
buf bytes.Buffer
func TestScanNestedDirectories(t *testing.T)
func TestLogOutputNoDoubleNewline(t *testing.T)
buf bytes.Buffer
func TestFilepathBaseEdgeCases(t *testing.T)
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

### /home/runner/work/Brf.it/Brf.it/website/docusaurus.config.ts

**Imports:**
- `import {themes as prismThemes} from 'prism-react-renderer';`
- `import type {Config} from '@docusaurus/types';`
- `import type * as Preset from '@docusaurus/preset-classic';`

```typescript
const config: Config = {
  title: 'Brf.it',
  tagline: 'Give AI the gist, not the bloat',
  favicon: 'img/favicon.ico',

  future: {
    v4: true,
  },

  url: 'https://indigo-net.github.io',
  baseUrl: '/Brf.it/',

  organizationName: 'indigo-net',
  projectName: 'Brf.it',
  trailingSlash: false,

  onBrokenLinks: 'throw',

  i18n: {
    defaultLocale: 'en',
    locales: ['en', 'ko', 'ja', 'de', 'hi'],
    localeConfigs: {
      en: { label: 'English', direction: 'ltr' },
      ko: { label: '한국어', direction: 'ltr' },
      ja: { label: '日本語', direction: 'ltr' },
      de: { label: 'Deutsch', direction: 'ltr' },
      hi: { label: 'हिन्दी', direction: 'ltr' },
    },
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/indigo-net/Brf.it/tree/main/website/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/docusaurus-social-card.jpg',
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'Brf.it',
      logo: {
        alt: 'Brf.it Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          type: 'localeDropdown',
          position: 'right',
          dropdownItemsBefore: [],
          dropdownItemsAfter: [],
        },
        {
          href: 'https://github.com/indigo-net/Brf.it',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/',
            },
            {
              label: 'CLI Reference',
              to: '/docs/cli-reference',
            },
            {
              label: 'Supported Languages',
              to: '/docs/languages/',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'GitHub Issues',
              href: 'https://github.com/indigo-net/Brf.it/issues',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/indigo-net/Brf.it',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} indigo-net. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['go', 'typescript', 'python', 'java', 'kotlin', 'rust', 'ruby', 'php', 'swift', 'scala'],
    },
  } satisfies Preset.ThemeConfig,
};
```

---

### /home/runner/work/Brf.it/Brf.it/website/sidebars.ts

**Imports:**
- `import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';`

```typescript
const sidebars: SidebarsConfig = {
  docsSidebar: [
    {
      type: 'doc',
      id: 'getting-started',
      label: 'Getting Started',
    },
    {
      type: 'doc',
      id: 'cli-reference',
      label: 'CLI Reference',
    },
    {
      type: 'category',
      label: 'Languages',
      items: [
        'languages/index',
        'languages/go',
        'languages/typescript',
        'languages/python',
        'languages/java',
        'languages/kotlin',
        'languages/rust',
        'languages/ruby',
        'languages/php',
        'languages/swift',
        'languages/scala',
        'languages/c-cpp',
      ],
    },
  ],
};
```

---

### /home/runner/work/Brf.it/Brf.it/website/src/components/FeaturesSection.tsx

**Imports:**
- `import React from 'react';`
- `import Translate from '@docusaurus/Translate';`

```typescript
interface Feature {
  title: string;
  description: string;
  icon: string;
}
```

---

### /home/runner/work/Brf.it/Brf.it/website/src/components/Hero.tsx

**Imports:**
- `import React from 'react';`
- `import Link from '@docusaurus/Link';`
- `import Translate, {translate} from '@docusaurus/Translate';`

**Exports:**
- `Hero`

```typescript
export default function Hero(): JSX.Element
```

---

### /home/runner/work/Brf.it/Brf.it/website/src/components/InstallSection.tsx

**Imports:**
- `import React, {useState} from 'react';`
- `import Translate from '@docusaurus/Translate';`

```typescript
type Platform = keyof typeof installCommands;
const copyToClipboard = (text: string)
```

---

### /home/runner/work/Brf.it/Brf.it/website/src/components/LanguageGrid.tsx

**Imports:**
- `import React from 'react';`
- `import Translate from '@docusaurus/Translate';`

**Exports:**
- `LanguageGrid`

```typescript
interface Language {
  name: string;
  icon: string;
}
const languages: Language[] = [
  { name: 'Go', icon: '🐹' },
  { name: 'TypeScript', icon: '📘' },
  { name: 'JavaScript', icon: '📒' },
  { name: 'Python', icon: '🐍' },
  { name: 'Java', icon: '☕' },
  { name: 'Kotlin', icon: '🟣' },
  { name: 'Rust', icon: '🦀' },
  { name: 'Ruby', icon: '💎' },
  { name: 'PHP', icon: '🐘' },
  { name: 'Swift', icon: '🍎' },
  { name: 'Scala', icon: '🔴' },
  { name: 'C/C++', icon: '⚙️' },
];
export default function LanguageGrid(): JSX.Element {
  return (
    <section className="section" style={{ background: 'var(--ifm-color-emphasis-100)' }}>
      <div className="container">
        <h2 className="section-title">
          <Translate id="languages.title">Supported Languages</Translate>
        </h2>
        <p className="section-subtitle">
          <Translate id="languages.subtitle">
            Tree-sitter powered parsing for accurate signature extraction
          </Translate>
        </p>

        <div className="language-grid">
          {languages.map(lang
```

---

### /home/runner/work/Brf.it/Brf.it/website/src/components/TokenComparison.tsx

**Imports:**
- `import React, {useState} from 'react';`
- `import Translate, {translate} from '@docusaurus/Translate';`

```typescript
interface CodeExample {
  language: string;
  label: string;
  before: string;
  after: string;
  beforeTokens: number;
  afterTokens: number;
  beforeLines: number;
  afterLines: number;
}
```

---

### /home/runner/work/Brf.it/Brf.it/website/src/pages/index.tsx

**Imports:**
- `import type {ReactNode} from 'react';`
- `import useDocusaurusContext from '@docusaurus/useDocusaurusContext';`
- `import Layout from '@theme/Layout';`
- `import Hero from '@site/src/components/Hero';`
- `import TokenComparison from '@site/src/components/TokenComparison';`
- `import FeaturesSection from '@site/src/components/FeaturesSection';`
- `import LanguageGrid from '@site/src/components/LanguageGrid';`
- `import InstallSection from '@site/src/components/InstallSection';`

**Exports:**
- `Home`

```typescript
export default function Home(): ReactNode
```

---

