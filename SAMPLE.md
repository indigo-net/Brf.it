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
func TestParseFlagsNoStdImports(t *testing.T)
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

	// NoStdImports excludes standard library imports from output.
	NoStdImports bool

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
func TestConfigToOptionsNoStdImports(t *testing.T)
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

	// NoStdImports excludes standard library imports from output.
	NoStdImports bool

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

	// NoStdImports excludes standard library imports from output.
	NoStdImports bool
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
func TestKindToTag(t *testing.T)
func TestXMLFormatterKindTags(t *testing.T)
func TestXMLFormatterNoStdImports(t *testing.T)
func TestXMLFormatterNoStdImportsAllFiltered(t *testing.T)
func TestMarkdownFormatterNoStdImports(t *testing.T)
func TestXMLFormatterNoStdImportsEmptyFile(t *testing.T)
func TestMarkdownFormatterNoStdImportsEmptyFile(t *testing.T)
func TestFormatterNoStdImportsDisabled(t *testing.T)
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/helpers.go

**Imports:**
- `import "strings"`

```go
func isStdLibImport(language, importPath string) bool
func isGoStdLib(importPath string) bool
func isPythonStdLib(importPath string) bool
func isJSStdLib(importPath string) bool
func isCStdLib(importPath string) bool
func isJavaStdLib(importPath string) bool
func isRustStdLib(importPath string) bool
pythonStdLibModules = map[string]struct{}{
	"abc": {}, "aifc": {}, "argparse": {}, "array": {}, "ast": {},
	"asynchat": {}, "asyncio": {}, "asyncore": {}, "atexit": {},
	"base64": {}, "bdb": {}, "binascii": {}, "binhex": {}, "bisect": {},
	"builtins": {}, "bz2": {},
	"calendar": {}, "cgi": {}, "cgitb": {}, "chunk": {}, "cmath": {},
	"cmd": {}, "code": {}, "codecs": {}, "codeop": {}, "collections": {},
	"colorsys": {}, "compileall": {}, "concurrent": {}, "configparser": {},
	"contextlib": {}, "contextvars": {}, "copy": {}, "copyreg": {},
	"cProfile": {}, "crypt": {}, "csv": {}, "ctypes": {}, "curses": {},
	"dataclasses": {}, "datetime": {}, "dbm": {}, "decimal": {}, "difflib": {},
	"dis": {}, "distutils": {},
	"email": {}, "encodings": {}, "enum": {}, "errno": {},
	"faulthandler": {}, "fcntl": {}, "filecmp": {}, "fileinput": {},
	"fnmatch": {}, "fractions": {}, "ftplib": {}, "functools": {},
	"gc": {}, "getopt": {}, "getpass": {}, "gettext": {}, "glob": {},
	"graphlib": {}, "grp": {}, "gzip": {},
	"hashlib": {}, "heapq": {}, "hmac": {}, "html": {}, "http": {},
	"idlelib": {}, "imaplib": {}, "imghdr": {}, "imp": {}, "importlib": {},
	"inspect": {}, "io": {}, "ipaddress": {}, "itertools": {},
	"json": {},
	"keyword": {},
	"lib2to3": {}, "linecache": {}, "locale": {}, "logging": {}, "lzma": {},
	"mailbox": {}, "mailcap": {}, "marshal": {}, "math": {}, "mimetypes": {},
	"mmap": {}, "modulefinder": {}, "multiprocessing": {},
	"netrc": {}, "nis": {}, "nntplib": {}, "numbers": {},
	"operator": {}, "optparse": {}, "os": {}, "ossaudiodev": {},
	"parser": {}, "pathlib": {}, "pdb": {}, "pickle": {}, "pickletools": {},
	"pipes": {}, "pkgutil": {}, "platform": {}, "plistlib": {}, "poplib": {},
	"posix": {}, "posixpath": {}, "pprint": {}, "profile": {}, "pstats": {},
	"pty": {}, "pwd": {}, "py_compile": {}, "pyclbr": {}, "pydoc": {},
	"queue": {}, "quopri": {},
	"random": {}, "re": {}, "readline": {}, "reprlib": {}, "resource": {},
	"rlcompleter": {}, "runpy": {},
	"sched": {}, "secrets": {}, "select": {}, "selectors": {}, "shelve": {},
	"shlex": {}, "shutil": {}, "signal": {}, "site": {}, "smtpd": {},
	"smtplib": {}, "sndhdr": {}, "socket": {}, "socketserver": {},
	"sqlite3": {}, "ssl": {}, "stat": {}, "statistics": {}, "string": {},
	"stringprep": {}, "struct": {}, "subprocess": {}, "sunau": {},
	"symtable": {}, "sys": {}, "sysconfig": {}, "syslog": {},
	"tabnanny": {}, "tarfile": {}, "telnetlib": {}, "tempfile": {},
	"termios": {}, "test": {}, "textwrap": {}, "threading": {}, "time": {},
	"timeit": {}, "tkinter": {}, "token": {}, "tokenize": {}, "tomllib": {},
	"trace": {}, "traceback": {}, "tracemalloc": {}, "tty": {}, "turtle": {},
	"turtledemo": {}, "types": {}, "typing": {},
	"unicodedata": {}, "unittest": {}, "urllib": {}, "uu": {}, "uuid": {},
	"venv": {},
	"warnings": {}, "wave": {}, "weakref": {}, "webbrowser": {},
	"winreg": {}, "winsound": {}, "wsgiref": {},
	"xdrlib": {}, "xml": {}, "xmlrpc": {},
	"zipapp": {}, "zipfile": {}, "zipimport": {}, "zlib": {},
	"zoneinfo": {},
	"_thread": {},
}
nodeBuiltinModules = map[string]struct{}{
	"assert": {}, "buffer": {}, "child_process": {}, "cluster": {},
	"console": {}, "constants": {}, "crypto": {}, "dgram": {},
	"diagnostics_channel": {}, "dns": {}, "domain": {}, "events": {},
	"fs": {}, "http": {}, "http2": {}, "https": {}, "inspector": {},
	"module": {}, "net": {}, "os": {}, "path": {}, "perf_hooks": {},
	"process": {}, "punycode": {}, "querystring": {}, "readline": {},
	"repl": {}, "stream": {}, "string_decoder": {}, "timers": {},
	"tls": {}, "trace_events": {}, "tty": {}, "url": {}, "util": {},
	"v8": {}, "vm": {}, "wasi": {}, "worker_threads": {}, "zlib": {},
}
func getEmptyComment(lang string) string
```

---

### /home/runner/work/Brf.it/Brf.it/pkg/formatter/helpers_test.go

**Imports:**
- `import "fmt"`
- `import "testing"`

```go
func TestIsStdLibImport(t *testing.T)
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
func (f *XMLFormatter) Format(data *PackageData) ([]byte, error)
buf bytes.Buffer
importLines []string
func escapeXML(s string) string
func kindToTag(kind string) string
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
func (m *MockParser) Parse(content string, opts *Options) (*ParseResult, error)
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
func (p *TreeSitterParser) Parse(content string, opts *parser.Options) (result *parser.ParseResult, err error)
imports []parser.ImportExport
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
func stripRustBody(text, kind string) string
func findRustBodyStart(text string) int
func refineSwiftClassKind(text string) string
func stripSwiftBody(text, kind string) string
func findSwiftBodyStart(text string) int
func stripKotlinBody(text, kind string) string
func findKotlinBodyStart(text string) int
func refineKotlinClassKind(text string) string
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
) ([]parser.ImportExport, error)
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

